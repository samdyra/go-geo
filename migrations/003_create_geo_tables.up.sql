-- Create the geo_data_list table
CREATE TABLE IF NOT EXISTS geo_data_list (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL UNIQUE,
    coordinate POINT,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL
);

-- Create an index on the table_name for faster lookups
CREATE INDEX idx_geo_data_list_table_name ON geo_data_list(table_name);