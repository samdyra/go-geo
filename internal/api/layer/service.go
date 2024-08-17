package layer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samdyra/go-geo/internal/utils"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type Service struct {
    db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
    return &Service{db: db}
}

func (s *Service) CreateLayer(layer LayerCreate, userID int64) error {
    query := `INSERT INTO layer (spatial_data_id, layer_name, coordinate, color, created_at, updated_at, created_by, updated_by)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
    
    now := time.Now()
    _, err := s.db.Exec(query, layer.SpatialDataID, layer.LayerName, layer.Coordinate, layer.Color, now, now, userID, userID)
    if err != nil {
        return errors.ErrInternalServer
    }
    
    return nil
}


func (s *Service) UpdateLayer(id int64, update LayerUpdate, userID int64) error {
    query := "UPDATE layer SET updated_at = $1, updated_by = $2"
    args := []interface{}{time.Now(), userID}
    argCount := 3

    if update.LayerName != nil {
        query += fmt.Sprintf(", layer_name = $%d", argCount)
        args = append(args, *update.LayerName)
        argCount++
    }
    if update.Coordinate != nil {
        query += fmt.Sprintf(", coordinate = $%d", argCount)
        args = append(args, *update.Coordinate)
        argCount++
    }
    if update.Color != nil {
        query += fmt.Sprintf(", color = $%d", argCount)
        args = append(args, *update.Color)
        argCount++
    }

    query += fmt.Sprintf(" WHERE id = $%d", argCount)
    args = append(args, id)

    result, err := s.db.Exec(query, args...)
    if err != nil {
        return errors.ErrInternalServer
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.ErrInternalServer
    }
    if rowsAffected == 0 {
        return errors.ErrNotFound
    }

    return nil
}

func (s *Service) DeleteLayer(id int64) error {
    tx, err := s.db.Beginx()
    if err != nil {
        return errors.ErrInternalServer
    }
    defer tx.Rollback()

    // Delete from layer_layer_group
    _, err = tx.Exec("DELETE FROM layer_layer_group WHERE layer_id = $1", id)
    if err != nil {
        return errors.ErrInternalServer
    }

    // Delete from layer
    result, err := tx.Exec("DELETE FROM layer WHERE id = $1", id)
    if err != nil {
        return errors.ErrInternalServer
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.ErrInternalServer
    }
    if rowsAffected == 0 {
        return errors.ErrNotFound
    }

    return tx.Commit()
}

func (s *Service) GetAllFormattedLayers() ([]FormattedLayer, error) {
    query := `SELECT l.id, l.layer_name, l.coordinate, l.color, sd.table_name, sd.type 
              FROM layer l
              JOIN spatial_data sd ON l.spatial_data_id = sd.id`
    
    return s.queryFormattedLayers(query)
}

func (s *Service) GetFormattedLayers(ids []int64) ([]FormattedLayer, error) {
    query := `SELECT l.id, l.layer_name, l.coordinate, l.color, sd.table_name, sd.type 
              FROM layer l
              JOIN spatial_data sd ON l.spatial_data_id = sd.id
              WHERE l.id IN (?)`
    
    query, args, err := sqlx.In(query, ids)
    if err != nil {
        return nil, errors.ErrInternalServer
    }
    
    query = s.db.Rebind(query)
    return s.queryFormattedLayers(query, args...)
}

func (s *Service) queryFormattedLayers(query string, args ...interface{}) ([]FormattedLayer, error) {
    rows, err := s.db.Query(query, args...)
    if err != nil {
        return nil, errors.ErrInternalServer
    }
    defer rows.Close()

    var result []FormattedLayer
    for rows.Next() {
        var id int64
        var layerName, color, tableName, dataType string
        var coordinate []float64
        err := rows.Scan(&id, &layerName, &coordinate, &color, &tableName, &dataType)
        if err != nil {
            return nil, errors.ErrInternalServer
        }

        layerType := utils.GetLayerType(dataType)
        paint := utils.GetPaint(dataType, color)

        layer := map[string]interface{}{
            "id": tableName,
            "source": map[string]interface{}{
                "type":  "vector",
                "tiles": fmt.Sprintf("http://localhost:8080/mvt/%s/{z}/{x}/{y}", tableName),
            },
            "source-layer": tableName,
            "type":         layerType,
            "paint":        paint,
        }

        layerJSON, err := json.Marshal(layer)
        if err != nil {
            return nil, errors.ErrInternalServer
        }

        result = append(result, FormattedLayer{
            ID:         id,
            LayerName:  layerName,
            Coordinate: coordinate,
            Layer:      layerJSON,
        })
    }

    if err = rows.Err(); err != nil {
        return nil, errors.ErrInternalServer
    }

    if len(result) == 0 {
        return nil, errors.ErrNotFound
    }

    return result, nil
}