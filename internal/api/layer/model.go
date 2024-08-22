package layer

import (
	"encoding/json"
	"time"
)

type Layer struct {
    ID             int64     `db:"id" json:"id"`
    SpatialDataID  int64     `db:"spatial_data_id" json:"spatial_data_id"`
    LayerName      string    `db:"layer_name" json:"layer_name"`
    Coordinate     []float64 `db:"coordinate" json:"coordinate"`
    Color          string    `db:"color" json:"color"`
    CreatedAt      time.Time `db:"created_at" json:"created_at"`
    UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
    CreatedBy      string     `db:"created_by" json:"created_by"`
    UpdatedBy      string     `db:"updated_by" json:"updated_by"`
}

type LayerCreate struct {
    SpatialDataID int64     `json:"spatial_data_id" binding:"required"`
    LayerName     string    `json:"layer_name" binding:"required"`
    Coordinate    []float64 `json:"coordinate" binding:"required"`
    Color         string    `json:"color" binding:"required"`
}

type LayerUpdate struct {
    LayerName  *string    `json:"layer_name"`
    Coordinate *[]float64 `json:"coordinate"`
    Color      *string    `json:"color"`
}

type FormattedLayer struct {
    ID         int64           `json:"id"`
    LayerName  string          `json:"layer_name"`
    Coordinate []float64       `json:"coordinate"`
    Layer      json.RawMessage `json:"layer"`
}