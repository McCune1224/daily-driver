-- Activities table for fitness tracking
CREATE TABLE activities (
    id BIGSERIAL PRIMARY KEY,
    activity_type VARCHAR(50) NOT NULL,
    distance_meters DECIMAL(10,2),
    duration_seconds INTEGER,
    avg_pace_min_per_km DECIMAL(5,2),
    avg_heart_rate INTEGER,
    max_heart_rate INTEGER,
    calories INTEGER,
    elevation_gain_meters DECIMAL(8,2),
    activity_date TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activities_date ON activities(activity_date DESC);
CREATE INDEX idx_activities_type ON activities(activity_type);

-- Tournaments table for Start.GG data
CREATE TABLE tournaments (
    id BIGSERIAL PRIMARY KEY,
    startgg_id BIGINT UNIQUE,
    tournament_name VARCHAR(255) NOT NULL,
    game VARCHAR(100),
    placement INTEGER,
    expected_seed INTEGER,
    total_entrants INTEGER,
    tournament_date DATE,
    location VARCHAR(255),
    event_name VARCHAR(255),
    bracket_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tournaments_date ON tournaments(tournament_date DESC);
CREATE INDEX idx_tournaments_game ON tournaments(game);

-- Art pieces cache from Art Institute of Chicago API
CREATE TABLE art_pieces (
    id BIGSERIAL PRIMARY KEY,
    api_id INTEGER UNIQUE,
    title VARCHAR(500),
    artist VARCHAR(255),
    date_display VARCHAR(100),
    image_id VARCHAR(255),
    image_url TEXT,
    description TEXT,
    department VARCHAR(255),
    artwork_type VARCHAR(100),
    last_fetched TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_art_pieces_artist ON art_pieces(artist);
CREATE INDEX idx_art_pieces_fetched ON art_pieces(last_fetched DESC);
