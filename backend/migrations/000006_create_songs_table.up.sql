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

CREATE UNIQUE INDEX idx_songs_track_number ON songs(track_number, album_id) WHERE album_id IS NOT NULL AND deleted_at IS NULL; 

CREATE MATERIALIZED VIEW active_songs AS SELECT 
  id,
  title,
  CASE WHEN s.main_artist_id IS NOT NULL AND ar.deleted_at IS NULL THEN s.main_artist_id ELSE NULL END AS main_artist_id,
  CASE WHEN s.album_id IS NOT NULL AND al.deleted_at IS NULL THEN s.album_id ELSE NULL END AS album_id,
  CASE WHEN s.album_id IS NOT NULL AND al.deleted_at IS NULL THEN s.track_number ELSE NULL END AS track_number,
  premiered_at,
  runtime_ms,
  CASE WHEN s.album_id IS NOT NULL AND al.deleted_at IS NULL THEN s.song_type ELSE 'orphaned'::song_type END AS song_type,
  upload_method,
  is_public,
  audio_key,
  cover_art_key,
  added_at,
  updated_at,
  deleted_at
FROM songs s
LEFT JOIN artists ar ON s.main_artist_id = ar.id
LEFT JOIN albums al ON s.album_id = al.id
WHERE deleted_at IS NULL;

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE UNIQUE INDEX idx_active_songs_id ON active_songs(id); --to allow parallel updates
CREATE INDEX idx_active_songs_title_fuzzy ON active_songs USING GIN(title gin_trgm_ops);
CREATE INDEX idx_active_songs_album ON active_songs(album_id);

CREATE TABLE plays (
  id UUID PRIMARY KEY,
  song_id UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  played_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_song_plays ON plays(song_id);

CREATE TABLE song_artists (
  song_id   UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
  artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE RESTRICT,
  PRIMARY KEY (song_id, artist_id)
);

CREATE INDEX idx_song_artists_artists ON song_artists(artist_id); 

CREATE TABLE song_contributors (
  song_id UUID REFERENCES songs(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  contributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (song_id, user_id)
);

CREATE INDEX idx_song_contributors ON song_contributors(user_id);

CREATE OR REPLACE FUNCTION orphan_songs_on_album_deletion()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.album_id IS NULL AND OLD.album_id IS NOT NULL THEN
    NEW.song_type = 'orphaned';
    NEW.track_number = NULL;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_orphan_songs ON songs;
CREATE TRIGGER trg_orphan_songs BEFORE UPDATE ON songs
FOR EACH ROW EXECUTE FUNCTION orphan_songs_on_album_deletion();

