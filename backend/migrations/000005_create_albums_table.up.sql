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
  premiered_at   TIMESTAMPTZ,
  uploader_id    UUID REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE album_artists (
  album_id  UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
  artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  PRIMARY KEY (album_id, artist_id)
);

CREATE TABLE album_contributors (
  album_id UUID REFERENCES albums(id) ON DELETE CASCADE,
  user_id  UUID REFERENCES users(id) ON DELETE CASCADE,
  contributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (album_id, user_id)
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_albums_name ON albums(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_albums_name_fuzzy ON albums USING GIN(name gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX idx_albums_album_artists ON album_artists(artist_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_albums_album_contributors ON album_contributors(user_id) WHERE deleted_at IS NULL;
