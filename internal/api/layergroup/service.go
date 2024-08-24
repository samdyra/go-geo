package layergroup

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateGroup(group LayerGroupCreate, username string) error {
	query := `INSERT INTO layer_group (group_name, created_at, updated_at, created_by, updated_by)
              VALUES ($1, $2, $3, $4, $5)`
	
	now := time.Now()
	_, err := s.db.Exec(query, group.GroupName, now, now, username, username)
	if err != nil {
		return errors.ErrInternalServer
	}
	
	return nil
}

func (s *Service) AddLayerToGroup(connection LayerToGroup, username string) error {
	query := `INSERT INTO layer_layer_group (layer_id, layer_group_id, created_at, updated_at, created_by, updated_by)
              VALUES ($1, $2, $3, $4, $5, $6)`
	
	now := time.Now()
	_, err := s.db.Exec(query, connection.LayerID, connection.GroupID, now, now, username, username)
	if err != nil {
		return errors.ErrInternalServer
	}
	
	return nil
}

func (s *Service) GetGroupsWithLayers() ([]GroupWithLayers, error) {
	query := `
		SELECT 
			lg.id AS group_id, 
			lg.group_name, 
			COALESCE(json_agg(
				json_build_object(
					'layer_id', l.id, 
					'layer_name', l.layer_name,
					'coordinate', l.coordinate,
					'color', l.color
				)
			) FILTER (WHERE l.id IS NOT NULL), '[]'::json) AS layers
		FROM 
			layer_group lg
		LEFT JOIN 
			layer_layer_group llg ON lg.id = llg.layer_group_id
		LEFT JOIN 
			layer l ON llg.layer_id = l.id
		GROUP BY 
			lg.id, lg.group_name
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	var results []GroupWithLayers
	for rows.Next() {
		var group GroupWithLayers
		var layersJSON []byte
		err := rows.Scan(&group.GroupID, &group.GroupName, &layersJSON)
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		err = json.Unmarshal(layersJSON, &group.Layers)
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		results = append(results, group)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ErrInternalServer
	}
	
	return results, nil
}

func (s *Service) RemoveLayerFromGroup(layerID, groupID int64) error {
	query := `DELETE FROM layer_layer_group
              WHERE layer_id = $1 AND layer_group_id = $2`
	
	result, err := s.db.Exec(query, layerID, groupID)
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

func (s *Service) DeleteGroup(groupID int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return errors.ErrInternalServer
	}
	defer tx.Rollback()

	// Delete related entries in layer_layer_group
	_, err = tx.Exec("DELETE FROM layer_layer_group WHERE layer_group_id = $1", groupID)
	if err != nil {
		return errors.ErrInternalServer
	}

	// Delete the group
	result, err := tx.Exec("DELETE FROM layer_group WHERE id = $1", groupID)
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

	if err := tx.Commit(); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}