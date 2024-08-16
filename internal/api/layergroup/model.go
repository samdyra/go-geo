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

type LayerLayerGroup struct {
    ID           int64     `db:"id" json:"id"`
    LayerID      int64     `db:"layer_id" json:"layer_id"`
    LayerGroupID int64     `db:"layer_group_id" json:"layer_group_id"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
    CreatedBy    int64     `db:"created_by" json:"created_by"`
    UpdatedBy    int64     `db:"updated_by" json:"updated_by"`
}