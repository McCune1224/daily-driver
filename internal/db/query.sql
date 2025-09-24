-- name: GetGarminFitFileByID :one
SELECT * FROM garmin_fit_files
WHERE id = $1;

-- name: ListGarminFitFilesByUser :many
SELECT * FROM garmin_fit_files
WHERE user_id = $1
ORDER BY uploaded_at DESC;

-- name: GetGarminFitFileByFilename :one
SELECT * FROM garmin_fit_files
WHERE filename = $1;

-- name: SearchGarminFitFilesByFilename :many
SELECT * FROM garmin_fit_files
WHERE filename ILIKE '%' || $1 || '%'
ORDER BY uploaded_at DESC;

-- name: ListGarminFitFilesByDateRange :many
SELECT * FROM garmin_fit_files
WHERE uploaded_at BETWEEN $1 AND $2
ORDER BY uploaded_at DESC;

-- name: ListGarminFilenames :many
SELECT filename FROM garmin_fit_files
ORDER BY uploaded_at DESC;
