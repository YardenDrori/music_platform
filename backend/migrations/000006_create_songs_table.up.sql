CREATE TYPE upload_source AS ENUM (
  'manual_upload',
  'youtube_scrape',
  'spotify_scrape',
  'apple_music_scrape',
  'other'
);

CREATE TABLE songs (
  id             UUID PRIMARY KEY,
  title          TEXT NOT NULL,
  main_artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  album_id       UUID REFERENCES albums(id) ON DELETE SET NULL,
  uploader_id    UUID REFERENCES users(id) ON DELETE SET NULL,
  is_public      BOOL NOT NULL DEFAULT FALSE,
  deleted_at     TIMESTAMPTZ,
  runtime_ms     INT NOT NULL,
  upload_method  upload_source NOT NULL,
  total_plays    INT NOT NULL DEFAULT 0,
  audio_key      UUID NOT NULL,
  cover_art_key  UUID,
  added_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  premiered_at   TIMESTAMPTZ
);

CREATE TABLE song_artists (
  song_id   UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
  artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  PRIMARY KEY (song_id, artist_id)
);
