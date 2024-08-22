CREATE TABLE IF NOT EXISTS layer_layer_group (
    id SERIAL PRIMARY KEY,
    layer_id INTEGER REFERENCES layer(id),
    layer_group_id INTEGER REFERENCES layer_group(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) REFERENCES users(username),
    updated_by VARCHAR(50) REFERENCES users(username),
    UNIQUE(layer_id, layer_group_id)
);