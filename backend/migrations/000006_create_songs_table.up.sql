CREATE TYPE upload_source AS ENUM (
  'manual_upload',
  'youtube_scrape',
  'spotify_scrape',
  'apple_music_scrape',
  'other'
);

CREATE TYPE song_type AS ENUM (
  'album_track',
  'single',
  'orphaned'
);

CREATE TABLE songs (
  id             UUID PRIMARY KEY,
  title          TEXT NOT NULL,
  main_artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  album_id       UUID REFERENCES albums(id) ON DELETE SET NULL,
  track_number   INT,
  premiered_at   TIMESTAMPTZ,
  runtime_ms     INT NOT NULL,
  song_type song_type NOT NULL,
  upload_method  upload_source NOT NULL,
  is_public      BOOL NOT NULL DEFAULT FALSE,
  audio_key      UUID NOT NULL,
  cover_art_key  UUID,
  added_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at     TIMESTAMPTZ
);

CREATE TABLE plays (
  id UUID PRIMARY KEY,
  song_id UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  played_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE song_artists (
  song_id   UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
  artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  PRIMARY KEY (song_id, artist_id)
);

CREATE TABLE song_contributors (
  song_id UUID REFERENCES songs(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  contributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (song_id, user_id)
);

CREATE OR REPLACE FUNCTION orphan_songs_on_album_deletion()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.album_id IS NULL AND OLD.album_id IS NOT NULL THEN
    NEW.song_type = 'orphaned';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_orphan_songs ON songs;
CREATE TRIGGER trg_orphan_songs BEFORE UPDATE ON songs
FOR EACH ROW EXECUTE FUNCTION orphan_songs_on_album_deletion();

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_songs_title ON songs(title) WHERE deleted_at IS NULL;
CREATE INDEX idx_songs_title_fuzzy ON songs USING GIN(title gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX idx_songs_main_artist ON songs(main_artist_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_songs_album ON songs(album_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_song_artists_artists ON song_artists(artist_id); 
CREATE INDEX idx_song_plays ON plays(songs_id);

CREATE UNIQUE INDEX idx_songs_track_number ON songs(track_number, album_id) WHERE album_id IS NOT NULL AND deleted_at IS NULL; 
