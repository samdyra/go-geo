CREATE TABLE IF NOT EXISTS layer_layer_group (
    id SERIAL PRIMARY KEY,
    layer_id INTEGER REFERENCES layer(id),
    layer_group_id INTEGER REFERENCES layer_group(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id),
    UNIQUE(layer_id, layer_group_id)
);