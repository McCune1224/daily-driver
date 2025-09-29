-- name: GetGarminFitFileByID :one
select *
from garmin_fit_files
where id = $1
;


-- name: GetGarminFitFileByFilename :one
select *
from garmin_fit_files
where filename = $1
;

-- name: SearchGarminFitFilesByFilename :many
select *
from garmin_fit_files
where filename ilike '%' || $1 || '%'
order by uploaded_at desc
;

-- name: ListGarminFitFilesByDateRange :many
select *
from garmin_fit_files
where uploaded_at between $1 and $2
order by uploaded_at desc
;

-- name: ListGarminFilenames :many
select filename
from garmin_fit_files
order by uploaded_at desc
;

-- name: ListGarminFilesByFileCategory :many
select *
from garmin_fit_files
where file_category = $1
;


-- name: ListGarminFilesPaginated :many
select *
from garmin_fit_files
order by uploaded_at desc
limit $1 offset $2
;

-- name: CountGarminFiles :one
select count(*) from garmin_fit_files
;


-- name: InsertGarminFitFile :one
insert into garmin_fit_files 
(filename, data, file_category)
values ($1, $2, $3)
returning *
;

