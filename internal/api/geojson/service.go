package geojson

import (
	"fmt"

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
					'properties', properties
				)
			)
		)::text
		FROM %s;
	`, tableName)

	var geojson []byte
	err := s.db.Get(&geojson, query)
	if err != nil {
		return nil, err
	}

	return geojson, nil
}