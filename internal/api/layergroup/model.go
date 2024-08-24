package layergroup

import "time"

type LayerGroup struct {
    ID        int64     `db:"id" json:"id"`
    GroupName string    `db:"group_name" json:"group_name"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
    CreatedBy int64     `db:"created_by" json:"created_by"`
    UpdatedBy int64     `db:"updated_by" json:"updated_by"`
}

type LayerGroupCreate struct {
    GroupName string `json:"group_name" binding:"required"`
}

type LayerToGroup struct {
    LayerID int64 `json:"layer_id" binding:"required"`
    GroupID int64 `json:"group_id" binding:"required"`
}

type GroupWithLayers struct {
    GroupID   int64         `db:"group_id" json:"group_id"`
    GroupName string        `db:"group_name" json:"group_name"`
    Layers    []LayerDetail `db:"layers" json:"layers"`
}

type LayerDetail struct {
    LayerID   int64  `json:"layer_id"`
    LayerName string `json:"layer_name"`
    Coordinate []float64 `json:"coordinate"`
    Color      string    `json:"color"`
}