CREATE TABLE artists(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  aliases TEXT[],
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

CREATE TABLE artist_contributors (
  artist_id UUID REFERENCES artists(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  contributed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (artist_id, user_id)
);
