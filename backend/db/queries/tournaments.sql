-- name: GetTournament :one
SELECT * FROM tournaments WHERE id = $1;

-- name: ListTournaments :many
SELECT * FROM tournaments 
ORDER BY tournament_date DESC 
LIMIT $1 OFFSET $2;

-- name: CreateTournament :one
INSERT INTO tournaments (
    startgg_id,
    tournament_name,
    game,
    placement,
    expected_seed,
    total_entrants,
    tournament_date,
    location,
    event_name,
    bracket_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateTournament :exec
UPDATE tournaments 
SET 
    placement = $2,
    total_entrants = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteTournament :exec
DELETE FROM tournaments WHERE id = $1;

-- name: GetUpcomingTournaments :many
SELECT * FROM tournaments 
WHERE tournament_date >= CURRENT_DATE 
ORDER BY tournament_date ASC;
