CREATE TABLE videos
(
    hash          TEXT PRIMARY KEY,
    path          TEXT NOT NULL,
    date_time     TEXT,
    iso           INTEGER,
    exposure_time TEXT,
    x_dimension   INTEGER,
    y_dimension   INTEGER,
    model         TEXT,
    f_number      TEXT,
    orientation   INTEGER
);

CREATE TABLE thumbnails
(
    hash      TEXT PRIMARY KEY,
    height    INTEGER NOT NULL,
    width     INTEGER NOT NULL,
    thumbnail BLOB,

    CONSTRAINT hash_fk FOREIGN KEY (hash) REFERENCES videos (hash),
    CONSTRAINT thumbnails_unique UNIQUE (hash, height, width)
)
