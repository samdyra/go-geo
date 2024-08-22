package spatialdata

import "time"

type SpatialDataCreate struct {
    TableName string `form:"table_name" binding:"required"`
    Type      string `form:"type" binding:"required"`
}

type SpatialData struct {
    ID        int64     `db:"id" json:"id"`
    TableName string    `db:"table_name" json:"table_name"`
    Type      string    `db:"type" json:"type"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
    CreatedBy string     `db:"created_by" json:"created_by"`
    UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type SpatialDataEdit struct {
    TableName *string `json:"table_name"`
}