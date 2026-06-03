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
