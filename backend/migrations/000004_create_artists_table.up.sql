CREATE TABLE artists(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  is_band BOOLEAN NOT NULL,
  artist_image_key UUID,
  artist_banner_key UUID,
  link_to_youtube TEXT,
  link_to_spotify TEXT,
  link_to_apple_music TEXT,
  origin_date DATE,
  origin_place CHAR(2),
  added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE artist_aliases (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
  alias TEXT NOT NULL,
  UNIQUE(artist_id, alias)
);

CREATE TABLE artist_contributors (
  artist_id UUID REFERENCES artists(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  contributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (artist_id, user_id)
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_artists_name ON artists(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_artists_fuzzy ON artists USING GIN(name gin_trgm_ops) WHERE deleted_at IS NULL;
CREATE INDEX idx_artists_aliases ON artist_aliases USING GIN(alias gin_trgm_ops) WHERE deleted_at IS NULL;
