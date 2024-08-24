package geojson

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type GeoJSONService struct {
	db *sqlx.DB
}

func NewGeoJSONService(db *sqlx.DB) *GeoJSONService {
	return &GeoJSONService{db: db}
}

func (s *GeoJSONService) GenerateGeoJSON(tableName string) ([]byte, error) {
	query := fmt.Sprintf(`
		SELECT json_build_object(
			'type', 'FeatureCollection',
			'features', json_agg(
				json_build_object(
					'type', 'Feature',
					'geometry', ST_AsGeoJSON(geom)::json,
					'properties', json_build_object(
						'wadmkc', wadmkc,
						'wadmkk', wadmkk,
						'wadmpr', wadmpr
					)
				)
			)
		)::text
		FROM %s;
	`, tableName)

	log.Printf("Executing query: %s", query)

	var geojson []byte
	err := s.db.Get(&geojson, query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}

	log.Printf("Successfully generated GeoJSON for table: %s", tableName)
	return geojson, nil
}