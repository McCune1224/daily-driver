CREATE TYPE garmin_file_category AS ENUM ('workout', 'activity', 'other');

-- Table for storing raw Garmin .FIT files
CREATE TABLE IF NOT EXISTS garmin_fit_files (
    id SERIAL PRIMARY KEY,
    filename TEXT UNIQUE NOT NULL,
    data BYTEA NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT NOW(),
    file_category garmin_file_category DEFAULT 'other'
);

