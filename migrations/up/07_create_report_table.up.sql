CREATE TABLE report (
    id SERIAL PRIMARY KEY,
    reporter_name VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    data_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);