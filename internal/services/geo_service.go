package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/geojson"
	"github.com/samdyra/go-geo/internal/models"
	"github.com/samdyra/go-geo/internal/utils"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type GeoService struct {
    db *sqlx.DB
}

func NewGeoService(db *sqlx.DB) *GeoService {
    return &GeoService{db: db}
}

type FormattedGeoData struct {
	Name       string          `json:"name"`
	Coordinate []float64       `json:"coordinate"`
	Layer      json.RawMessage `json:"layer"`
}

func (s *GeoService) GetFormattedGeoData() ([]FormattedGeoData, error) {
	log.Println("Starting GetFormattedGeoData")
	
	query := `SELECT table_name, type, color, coordinate FROM geo_data_list`
	log.Printf("Executing query: %s", query)
	
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, errors.ErrInternalServer
	}
	defer rows.Close()

	log.Println("Query executed successfully")

	var result []FormattedGeoData
	for rows.Next() {
		var tableName, dataType, color string
		var coordinateStr []string
		err := rows.Scan(&tableName, &dataType, &color, pq.Array(&coordinateStr))
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, errors.ErrInternalServer
		}
		log.Printf("Scanned row: table_name=%s, type=%s, color=%s, coordinate=%v", tableName, dataType, color, coordinateStr)

		coordinate, err := utils.ParseCoordinate(coordinateStr)
		if err != nil {
			log.Printf("Error parsing coordinate: %v", err)
			return nil, errors.ErrInvalidInput
		}
		log.Printf("Parsed coordinate: %v", coordinate)

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
			log.Printf("Error marshaling layer to JSON: %v", err)
			return nil, errors.ErrInternalServer
		}

		result = append(result, FormattedGeoData{
			Name:       utils.FormatTableName(tableName),
			Coordinate: coordinate,
			Layer:      layerJSON,
		})
		log.Printf("Added formatted data for table: %s", tableName)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		return nil, errors.ErrInternalServer
	}

	if len(result) == 0 {
		log.Println("No results found")
		return nil, errors.ErrNotFound
	}

	log.Printf("Returning %d formatted geo data entries", len(result))
	return result, nil
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
        INSERT INTO geo_data_list (table_name, coordinate, type, color, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, data.TableName, pq.Array(data.Coordinate), data.Type, data.Color, now, now, username, username)
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

func (s *GeoService) EditGeoData(oldTableName string, data models.GeoDataEdit, file io.Reader, username string) error {
    tx, err := s.db.Beginx()
    if err != nil {
        return errors.ErrInternalServer
    }
    defer tx.Rollback()

    // Prepare the update query for geo_data_list
    query := "UPDATE geo_data_list SET updated_at = $1, updated_by = $2"
    params := []interface{}{time.Now(), username}
    paramCount := 3

    if data.Type != nil {
        query += fmt.Sprintf(", type = $%d", paramCount)
        params = append(params, *data.Type)
        paramCount++
    }

    if data.Color != nil {
        query += fmt.Sprintf(", color = $%d", paramCount)
        params = append(params, *data.Color)
        paramCount++
    }

    if len(data.Coordinate) > 0 {
        query += fmt.Sprintf(", coordinate = $%d", paramCount)
        params = append(params, pq.Array(data.Coordinate))
        paramCount++
    }

    query += " WHERE table_name = $" + fmt.Sprintf("%d", paramCount)
    params = append(params, oldTableName)

    // Execute the update on geo_data_list
    result, err := tx.Exec(query, params...)
    if err != nil {
        log.Printf("Error updating geo_data_list: %v", err)
        return errors.ErrInternalServer
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected: %v", err)
        return errors.ErrInternalServer
    }

    if rowsAffected == 0 {
        return errors.ErrNotFound
    }

    // If table name was changed and file is provided, rename the table and update its contents
    if data.TableName != nil && *data.TableName != oldTableName && file != nil {
        // Check if the new table name already exists
        var exists bool
        err = tx.Get(&exists, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", *data.TableName)
        if err != nil {
            log.Printf("Error checking if new table name exists: %v", err)
            return errors.ErrInternalServer
        }
        if exists {
            return errors.ErrTableAlreadyExists
        }

        // Rename the spatial data table
        _, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldTableName, *data.TableName))
        if err != nil {
            log.Printf("Error renaming table from %s to %s: %v", oldTableName, *data.TableName, err)
            return errors.ErrInternalServer
        }

        log.Printf("Renamed spatial data table from %s to %s", oldTableName, *data.TableName)

        // Clear existing data from the renamed table
        _, err = tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", *data.TableName))
        if err != nil {
            log.Printf("Error clearing data from table %s: %v", *data.TableName, err)
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

        // Insert new features into the renamed table
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
            `, *data.TableName), wkbData, properties, username, username)
            if err != nil {
                return errors.ErrInternalServer
            }
        }
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Error committing transaction: %v", err)
        return errors.ErrInternalServer
    }

    return nil
}