package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/geojson"
	"github.com/samdyra/go-geo/internal/models"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type GeoService struct {
    db *sqlx.DB
}

func NewGeoService(db *sqlx.DB) *GeoService {
    return &GeoService{db: db}
}

func (s *GeoService) CreateGeoData(data models.GeoDataUpload, file io.Reader) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return errors.ErrInternalServer
	}
	defer tx.Rollback()

	// Check if table already exists
	var exists bool
	err = tx.Get(&exists, "SELECT EXISTS (SELECT FROM geo_data_list WHERE table_name = $1)", data.TableName)
	if err != nil {
		return errors.ErrInternalServer
	}
	if exists {
		return errors.ErrUserAlreadyExists
	}

	// Insert into geo_data_list
	_, err = tx.Exec(`
		INSERT INTO geo_data_list (table_name, coordinate, type)
		VALUES ($1, $2, $3)
	`, data.TableName, data.Coordinate, data.Type)
	if err != nil {
		return errors.ErrInternalServer
	}

	// Create new table for the geo data
	_, err = tx.Exec(fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        geom GEOMETRY,
        properties JSONB
    )
	`, data.TableName))
	if err != nil {
		log.Printf("Error creating table %s: %v", data.TableName, err)
		return errors.ErrInternalServer
	}

	// Read the entire file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return errors.ErrInvalidInput
	}

	// Parse GeoJSON
	fc, err := geojson.UnmarshalFeatureCollection(fileBytes)
	if err != nil {
		return errors.ErrInvalidInput
	}

	// Insert features into the new table
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
			INSERT INTO %s (geom, properties)
			VALUES (ST_GeomFromWKB($1, 4326), $2)
		`, data.TableName), wkbData, properties)
		if err != nil {
			return errors.ErrInternalServer
		}
	}

	return tx.Commit()
}

func (s *GeoService) DeleteGeoData(tableName string) error {
    tx, err := s.db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Delete from geo_data_list
    _, err = tx.Exec("DELETE FROM geo_data_list WHERE table_name = $1", tableName)
    if err != nil {
        return err
    }

    // Drop the geo data table
    _, err = tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
    if err != nil {
        return err
    }

    return tx.Commit()
}