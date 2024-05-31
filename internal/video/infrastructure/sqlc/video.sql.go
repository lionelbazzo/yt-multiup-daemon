// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: video.sql

package sqlc

import (
	"context"
	"database/sql"
)

const countVideoByHash = `-- name: CountVideoByHash :one
SELECT COUNT(*)
FROM videos
WHERE hash = ?
`

func (q *Queries) CountVideoByHash(ctx context.Context, db DBTX, hash string) (int64, error) {
	row := db.QueryRowContext(ctx, countVideoByHash, hash)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countVideos = `-- name: CountVideos :one
SELECT COUNT(*)
FROM videos
`

func (q *Queries) CountVideos(ctx context.Context, db DBTX) (int64, error) {
	row := db.QueryRowContext(ctx, countVideos)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createOrReplaceThumbnail = `-- name: CreateOrReplaceThumbnail :exec
REPLACE INTO thumbnails (hash, height, width, thumbnail)
VALUES (?, ?, ?, ?)
`

type CreateOrReplaceThumbnailParams struct {
	Hash      string `db:"hash"`
	Height    int64  `db:"height"`
	Width     int64  `db:"width"`
	Thumbnail []byte `db:"thumbnail"`
}

func (q *Queries) CreateOrReplaceThumbnail(ctx context.Context, db DBTX, arg CreateOrReplaceThumbnailParams) error {
	_, err := db.ExecContext(ctx, createOrReplaceThumbnail,
		arg.Hash,
		arg.Height,
		arg.Width,
		arg.Thumbnail,
	)
	return err
}

const createOrReplaceVideo = `-- name: CreateOrReplaceVideo :exec
REPLACE INTO videos (hash, path, date_time, iso, exposure_time, x_dimension, y_dimension, model, f_number, orientation)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type CreateOrReplaceVideoParams struct {
	Hash         string         `db:"hash"`
	Path         string         `db:"path"`
	DateTime     sql.NullString `db:"date_time"`
	Iso          sql.NullInt64  `db:"iso"`
	ExposureTime sql.NullString `db:"exposure_time"`
	XDimension   sql.NullInt64  `db:"x_dimension"`
	YDimension   sql.NullInt64  `db:"y_dimension"`
	Model        sql.NullString `db:"model"`
	FNumber      sql.NullString `db:"f_number"`
	Orientation  sql.NullInt64  `db:"orientation"`
}

func (q *Queries) CreateOrReplaceVideo(ctx context.Context, db DBTX, arg CreateOrReplaceVideoParams) error {
	_, err := db.ExecContext(ctx, createOrReplaceVideo,
		arg.Hash,
		arg.Path,
		arg.DateTime,
		arg.Iso,
		arg.ExposureTime,
		arg.XDimension,
		arg.YDimension,
		arg.Model,
		arg.FNumber,
		arg.Orientation,
	)
	return err
}

const deleteAllVideoInPath = `-- name: DeleteAllVideoInPath :exec
DELETE
FROM videos
WHERE path LIKE ?
`

func (q *Queries) DeleteAllVideoInPath(ctx context.Context, db DBTX, path string) error {
	_, err := db.ExecContext(ctx, deleteAllVideoInPath, path)
	return err
}

const deleteAllVideos = `-- name: DeleteAllVideos :exec
DELETE
FROM videos
WHERE 1
`

func (q *Queries) DeleteAllVideos(ctx context.Context, db DBTX) error {
	_, err := db.ExecContext(ctx, deleteAllVideos)
	return err
}

const deleteVideoByPath = `-- name: DeleteVideoByPath :exec
DELETE
FROM videos
WHERE path = ?
`

func (q *Queries) DeleteVideoByPath(ctx context.Context, db DBTX, path string) error {
	_, err := db.ExecContext(ctx, deleteVideoByPath, path)
	return err
}

const getThumbnail = `-- name: GetThumbnail :one
SELECT thumbnail
FROM thumbnails
WHERE hash = ?
  AND width = ?
  AND height = ?
`

type GetThumbnailParams struct {
	Hash   string `db:"hash"`
	Width  int64  `db:"width"`
	Height int64  `db:"height"`
}

func (q *Queries) GetThumbnail(ctx context.Context, db DBTX, arg GetThumbnailParams) ([]byte, error) {
	row := db.QueryRowContext(ctx, getThumbnail, arg.Hash, arg.Width, arg.Height)
	var thumbnail []byte
	err := row.Scan(&thumbnail)
	return thumbnail, err
}

const getVideo = `-- name: GetVideo :one
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
WHERE hash = ?
`

func (q *Queries) GetVideo(ctx context.Context, db DBTX, hash string) (Video, error) {
	row := db.QueryRowContext(ctx, getVideo, hash)
	var i Video
	err := row.Scan(
		&i.Hash,
		&i.Path,
		&i.DateTime,
		&i.Iso,
		&i.ExposureTime,
		&i.XDimension,
		&i.YDimension,
		&i.Model,
		&i.FNumber,
		&i.Orientation,
	)
	return i, err
}

const list = `-- name: List :many
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
LIMIT ? OFFSET ?
`

type ListParams struct {
	Limit  int64 `db:"limit"`
	Offset int64 `db:"offset"`
}

func (q *Queries) List(ctx context.Context, db DBTX, arg ListParams) ([]Video, error) {
	rows, err := db.QueryContext(ctx, list, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Video{}
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.Hash,
			&i.Path,
			&i.DateTime,
			&i.Iso,
			&i.ExposureTime,
			&i.XDimension,
			&i.YDimension,
			&i.Model,
			&i.FNumber,
			&i.Orientation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}