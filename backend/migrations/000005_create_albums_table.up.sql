CREATE TABLE albums (
  id             UUID PRIMARY KEY,
  name           TEXT NOT NULL,
  description    TEXT,
  main_artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  album_art_key  UUID,
  has_all_tracks BOOL NOT NULL DEFAULT FALSE,
  added_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at     TIMESTAMPTZ,
  premiered_at   TIMESTAMPTZ
);

CREATE TABLE album_artists (
  album_id  UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
  artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  PRIMARY KEY (album_id, artist_id)
);

CREATE INDEX idx_albums_album_artists ON album_artists(artist_id);

CREATE TABLE album_contributors (
  album_id UUID REFERENCES albums(id) ON DELETE CASCADE,
  user_id  UUID REFERENCES users(id) ON DELETE CASCADE,
  contributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (album_id, user_id)
);

CREATE INDEX idx_albums_album_contributors ON album_contributors(user_id);

CREATE MATERIALIZED VIEW active_albums AS
SELECT
    al.id,
    al.name,
    al.description,
    CASE WHEN ar.deleted_at IS NULL THEN al.main_artist_id ELSE NULL END AS main_artist_id,
    al.album_art_key,
    al.has_all_tracks,
    al.added_at,
    al.updated_at,
    al.deleted_at,
    al.premiered_at
FROM albums al
LEFT JOIN artists ar ON ar.id = al.main_artist_id
WHERE al.deleted_at IS NULL;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE UNIQUE INDEX idx_active_albums_id ON active_albums(id);
CREATE INDEX idx_active_albums_name_fuzzy ON active_albums USING GIN(name gin_trgm_ops);
