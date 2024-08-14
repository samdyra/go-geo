package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

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

func (s *GeoService) CreateGeoData(data models.GeoDataUpload, file io.Reader, username string) error {
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

    now := time.Now()

    // Insert into geo_data_list
    _, err = tx.Exec(`
        INSERT INTO geo_data_list (table_name, coordinate, type, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, data.TableName, data.Coordinate, data.Type, now, now, username, username)
    if err != nil {
        return errors.ErrInternalServer
    }

    // Create new table for the geo data
    _, err = tx.Exec(fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id SERIAL PRIMARY KEY,
        geom GEOMETRY,
        properties JSONB,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        created_by VARCHAR(255),
        updated_by VARCHAR(255)
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

        // @TODO: Delete created_by and updated_by from the table
        _, err = tx.Exec(fmt.Sprintf(`
            INSERT INTO %s (geom, properties, created_by, updated_by)
            VALUES (ST_GeomFromWKB($1, 4326), $2, $3, $4)
        `, data.TableName), wkbData, properties, username, username)
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