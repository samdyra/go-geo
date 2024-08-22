package mvt

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MVTService struct {
	db *sqlx.DB
}

func NewMVTService(db *sqlx.DB) *MVTService {
	return &MVTService{db: db}
}

func (s *MVTService) GenerateMVT(tableName string, z, x, y int) ([]byte, error) {
	query := fmt.Sprintf(`
		WITH mvt_geom AS (
			SELECT ST_AsMVTGeom(
				ST_Transform(geom, 3857),
				ST_TileEnvelope(%d, %d, %d)
			) AS geom,
			properties
			FROM %s
			WHERE ST_Intersects(
				geom,
				ST_Transform(ST_TileEnvelope(%d, %d, %d), 4326)
			)
		)
		SELECT ST_AsMVT(mvt_geom.*, '%s', 4096, 'geom') FROM mvt_geom;
	`, z, x, y, tableName, z, x, y, tableName)

	var mvt []byte
	err := s.db.Get(&mvt, query)
	if err != nil {
		return nil, err
	}

	return mvt, nil
}