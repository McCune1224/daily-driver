-- name: GetActivity :one
SELECT * FROM activities WHERE id = $1;

-- name: ListActivities :many
SELECT * FROM activities 
ORDER BY activity_date DESC 
LIMIT $1 OFFSET $2;

-- name: CreateActivity :one
INSERT INTO activities (
    activity_type, 
    distance_meters, 
    duration_seconds, 
    avg_pace_min_per_km, 
    avg_heart_rate,
    max_heart_rate,
    calories,
    elevation_gain_meters,
    activity_date,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateActivity :exec
UPDATE activities 
SET 
    activity_type = $2,
    distance_meters = $3,
    duration_seconds = $4,
    avg_pace_min_per_km = $5,
    notes = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteActivity :exec
DELETE FROM activities WHERE id = $1;

-- name: GetActivityStats :one
SELECT 
    COUNT(*) as total_activities,
    SUM(distance_meters) as total_distance,
    SUM(duration_seconds) as total_duration,
    AVG(avg_pace_min_per_km) as avg_pace
FROM activities
WHERE activity_type = $1;
