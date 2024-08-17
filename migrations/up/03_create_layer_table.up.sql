CREATE TABLE IF NOT EXISTS layer (
    id SERIAL PRIMARY KEY,
    spatial_data_id INTEGER REFERENCES spatial_data(id),
    layer_name VARCHAR(100) NOT NULL,
    coordinate POINT,
    color VARCHAR(7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id)
);