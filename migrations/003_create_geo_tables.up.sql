-- Create the geo_data_list table
CREATE TABLE IF NOT EXISTS geo_data_list (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL UNIQUE,
    coordinate VARCHAR(255)[2] NOT NULL,
    type VARCHAR(50) NOT NULL,
    color VARCHAR(7) NOT NULL DEFAULT '#000000',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    CONSTRAINT check_type CHECK (type IN ('POINT', 'LINESTRING', 'POLYGON'))
);

-- Create an index on the table_name for faster lookups
CREATE INDEX idx_geo_data_list_table_name ON geo_data_list(table_name);