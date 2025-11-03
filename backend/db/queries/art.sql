-- name: GetArtPiece :one
SELECT * FROM art_pieces WHERE id = $1;

-- name: GetArtByApiId :one
SELECT * FROM art_pieces WHERE api_id = $1;

-- name: ListArtPieces :many
SELECT * FROM art_pieces 
ORDER BY last_fetched DESC 
LIMIT $1 OFFSET $2;

-- name: CreateArtPiece :one
INSERT INTO art_pieces (
    api_id,
    title,
    artist,
    date_display,
    image_id,
    image_url,
    description,
    department,
    artwork_type
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: UpdateArtPieceFetchTime :exec
UPDATE art_pieces 
SET last_fetched = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetRandomArtPiece :one
SELECT * FROM art_pieces 
ORDER BY RANDOM() 
LIMIT 1;
