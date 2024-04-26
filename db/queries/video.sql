-- name: GetVideo :one
SELECT hash,
       path,
       date_time,
       iso,
       exposure_time,
       x_dimension,
       y_dimension,
       model,
       f_number,
       orientation
FROM videos
WHERE hash = ?;

-- name: GetThumbnail :one
SELECT thumbnail
FROM thumbnails
WHERE hash = ?
  AND width = ?
  AND height = ?;

-- name: CreateOrReplaceThumbnail :exec
REPLACE INTO thumbnails (hash, height, width, thumbnail)
VALUES (?, ?, ?, ?);

-- name: CreateOrReplaceVideo :exec
REPLACE INTO videos (hash, path, date_time, iso, exposure_time, x_dimension, y_dimension, model, f_number, orientation)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CountVideoByHash :one
SELECT COUNT(*)
FROM videos
WHERE hash = ?;

-- name: CountVideos :one
SELECT COUNT(*)
FROM videos;

-- name: List :many
SELECT hash,
       path,
       date_time,
       iso,
       exposure_time,
       x_dimension,
       y_dimension,
       model,
       f_number,
       orientation
FROM videos
ORDER BY date_time DESC
LIMIT ? OFFSET ?;

-- name: DeleteVideoByPath :exec
DELETE
FROM videos
WHERE path = ?;

-- name: DeleteAllVideoInPath :exec
DELETE
FROM videos
WHERE path LIKE ?;

-- name: DeleteAllVideos :exec
DELETE
FROM videos
WHERE 1
