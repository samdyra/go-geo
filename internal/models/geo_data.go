package models

import (
	"time"
)

type GeoDataList struct {
    ID          int64     `db:"id" json:"id"`
    TableName   string    `db:"table_name" json:"table_name"`
    Coordinate  *string   `db:"coordinate" json:"coordinate,omitempty"`
    Type        string    `db:"type" json:"type"`
    Color       string    `db:"color" json:"color"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
    CreatedBy   string    `db:"created_by" json:"created_by"`
    UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

type GeoDataUpload struct {
    TableName   string  `form:"table_name" binding:"required"`
    Type        string  `form:"type" binding:"required,oneof=POINT LINESTRING POLYGON"`
    Color       string  `form:"color" binding:"required,hexcolor"`
    Coordinate  *string `form:"coordinate"`
}