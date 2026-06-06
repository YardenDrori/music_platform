CREATE TABLE artists(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  artist_image_key UUID,
  artist_banner_key UUID,
  link_to_youtube TEXT,
  link_to_spotify TEXT,
  link_to_apple_music TEXT,
  birth_date DATE,
  birth_placeCHAR(2),
  added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  uploader_id UUID REFERENCES users(id) ON DELETE SET NULL
);
