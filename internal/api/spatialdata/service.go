package spatialdata

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/geojson"
	"github.com/samdyra/go-geo/internal/utils"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type SpatialDataService struct {
    db *sqlx.DB
}

func NewSpatialDataService(db *sqlx.DB) *SpatialDataService {
    return &SpatialDataService{db: db}
}
func (s *SpatialDataService) CreateSpatialData(spatial_data SpatialDataCreate, file io.Reader, username string) error {
    var existsInSchema, existsInSpatialData bool

    // Check in information_schema.tables
    err := s.db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", spatial_data.TableName).Scan(&existsInSchema)
    if err != nil {
        return fmt.Errorf("error checking schema: %w", errors.ErrInternalServer)
    }

    // Check in spatial_data table
    err = s.db.QueryRow("SELECT EXISTS (SELECT FROM spatial_data WHERE table_name = $1)", spatial_data.TableName).Scan(&existsInSpatialData)
    if err != nil {
        return fmt.Errorf("error checking spatial_data: %w", errors.ErrInternalServer)
    }

    if existsInSchema {
        return errors.ErrResourceAlreadyExists
    }

    if existsInSpatialData {
        return errors.ErrResourceAlreadyExists
    }
    
    tx, err := s.db.Beginx()
    if err != nil {
        return errors.ErrInternalServer
    }
    defer tx.Rollback()

    now := time.Now()

    _, err = tx.Exec(`
        INSERT INTO spatial_data (table_name, type, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, spatial_data.TableName, spatial_data.Type, now, now, username, username)
    if err != nil {
        return errors.ErrInternalServer
    }

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return errors.ErrInvalidInput
    }

    fc, err := geojson.UnmarshalFeatureCollection(fileBytes)
    if err != nil {
        return errors.ErrInvalidInput
    }

    propertyTypes := make(map[string]string)
    for _, feature := range fc.Features {
        for key, value := range feature.Properties {
            inferredType := utils.InferPostgresType(value)
            if existingType, ok := propertyTypes[key]; ok {
                propertyTypes[key] = utils.ReconcileTypes(existingType, inferredType)
            } else {
                propertyTypes[key] = inferredType
            }
        }
    }

    createTableSQL := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        geom GEOMETRY,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        created_by VARCHAR(255),
        updated_by VARCHAR(255)`, spatial_data.TableName)

    for propName, propType := range propertyTypes {
        createTableSQL += fmt.Sprintf(",\n        %s %s", propName, propType)
    }
    createTableSQL += "\n    )"

    _, err = tx.Exec(createTableSQL)
    if err != nil {
        return errors.ErrInternalServer
    }

    insertSQL := fmt.Sprintf("INSERT INTO %s (geom, created_by, updated_by", spatial_data.TableName)
    valueSQL := ") VALUES (ST_GeomFromWKB($1, 4326), $2, $3"
    paramCount := 4

    propertyNames := make([]string, 0, len(propertyTypes))
    for propName := range propertyTypes {
        propertyNames = append(propertyNames, propName)
        insertSQL += fmt.Sprintf(", %s", propName)
        valueSQL += fmt.Sprintf(", $%d", paramCount)
        paramCount++
    }
    insertSQL += valueSQL + ")"

    for _, feature := range fc.Features {
        geom := feature.Geometry
        wkbData, err := wkb.Marshal(geom)
        if err != nil {
            return errors.ErrInternalServer
        }

        params := []interface{}{wkbData, username, username}
        for _, propName := range propertyNames {
            if val, ok := feature.Properties[propName]; ok {
                convertedVal, err := utils.ConvertToType(val, propertyTypes[propName])
                if err != nil {
                    return errors.ErrInvalidInput
                }
                params = append(params, convertedVal)
            } else {
                params = append(params, nil)
            }
        }

        _, err = tx.Exec(insertSQL, params...)
        if err != nil {
            return errors.ErrInternalServer
        }
    }

    return tx.Commit()
}



func (s *SpatialDataService) GetSpatialDataList() ([]SpatialData, error) {
    query := `SELECT id, table_name, type, created_at, updated_at, created_by, updated_by FROM spatial_data`
    
    var spatialDataList []SpatialData
    err := s.db.Select(&spatialDataList, query)
    if err != nil {
        return nil, errors.ErrInternalServer
    }

    return spatialDataList, nil
}

func (s *SpatialDataService) DeleteSpatialData(tableName string) error {
    tx, err := s.db.Beginx()
    if err != nil {
        return errors.ErrInternalServer
    }
    defer tx.Rollback()

    // Delete from layer_layer_group
    _, err = tx.Exec("DELETE FROM layer_layer_group WHERE layer_id IN (SELECT id FROM layer WHERE spatial_data_id = (SELECT id FROM spatial_data WHERE table_name = $1))", tableName)
    if err != nil {
        return errors.ErrInternalServer
    }

    // Delete associated layers
    _, err = tx.Exec("DELETE FROM layer WHERE spatial_data_id = (SELECT id FROM spatial_data WHERE table_name = $1)", tableName)
    if err != nil {
        return errors.ErrInternalServer
    }

    // Delete from spatial_data table
    result, err := tx.Exec("DELETE FROM spatial_data WHERE table_name = $1", tableName)
    if err != nil {
        return errors.ErrInternalServer
    }

    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
        return errors.ErrNotFound
    }

    // Drop the spatial spatial_data table
    _, err = tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
    if err != nil {
        return errors.ErrInternalServer
    }

    return tx.Commit()
}

func (s *SpatialDataService) EditSpatialData(oldTableName string, spatial_data SpatialDataEdit, file io.Reader, username string) error {
    tx, err := s.db.Beginx()
    if err != nil {
        return errors.ErrInternalServer
    }
    defer tx.Rollback()

    // Update spatial_data table
    query := "UPDATE spatial_data SET updated_at = $1, updated_by = $2"
    params := []interface{}{time.Now(), username}
    paramCount := 3

    if spatial_data.TableName != nil {
        query += fmt.Sprintf(", table_name = $%d", paramCount)
        params = append(params, *spatial_data.TableName)
        paramCount++
    }

    query += " WHERE table_name = $" + fmt.Sprintf("%d", paramCount)
    params = append(params, oldTableName)

    result, err := tx.Exec(query, params...)
    if err != nil {
        return errors.ErrInternalServer
    }

    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
        return errors.ErrNotFound
    }

    // If table name was changed
    if spatial_data.TableName != nil && *spatial_data.TableName != oldTableName {
        // Rename the spatial spatial_data table
        _, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldTableName, *spatial_data.TableName))
        if err != nil {
            return errors.ErrInternalServer
        }
        oldTableName = *spatial_data.TableName
    }

    // If new file is provided
    if file != nil {
        // Clear existing spatial_data from the table
        _, err = tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", oldTableName))
        if err != nil {
            return errors.ErrInternalServer
        }

        // Read and parse the new GeoJSON file
        fileBytes, err := io.ReadAll(file)
        if err != nil {
            return errors.ErrInvalidInput
        }

        fc, err := geojson.UnmarshalFeatureCollection(fileBytes)
        if err != nil {
            return errors.ErrInvalidInput
        }

        // Insert new features into the table
        for _, feature := range fc.Features {
            geom := feature.Geometry
            properties, err := json.Marshal(feature.Properties)
            if err != nil {
                return errors.ErrInternalServer
            }

            wkbData, err := wkb.Marshal(geom)
            if err != nil {
                return errors.ErrInternalServer
            }

            _, err = tx.Exec(fmt.Sprintf(`
                INSERT INTO %s (geom, properties, created_by, updated_by)
                VALUES (ST_GeomFromWKB($1, 4326), $2, $3, $4)
            `, oldTableName), wkbData, properties, username, username)
            if err != nil {
                return errors.ErrInternalServer
            }
        }
    }

    return tx.Commit()
}