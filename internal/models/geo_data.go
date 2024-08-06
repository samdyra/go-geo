package models

import (
	"time"
)

type GeoDataList struct {
    ID          int64      `db:"id" json:"id"`
    TableName   string     `db:"table_name" json:"table_name"`
    Coordinate  *string    `db:"coordinate" json:"coordinate,omitempty"`
    Type        string     `db:"type" json:"type"`
    CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

type GeoDataUpload struct {
    TableName   string  `form:"table_name" binding:"required"`
    Type        string  `form:"type" binding:"required"`
    Coordinate  *string `form:"coordinate"`
}