-- Table for storing raw Garmin .FIT files
CREATE TABLE IF NOT EXISTS garmin_fit_files (
    id SERIAL PRIMARY KEY,
    filename TEXT,
    data BYTEA NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT NOW(),
    user_id INTEGER
    -- Add more metadata columns as needed
);

