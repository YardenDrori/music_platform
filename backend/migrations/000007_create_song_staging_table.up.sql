CREATE TABLE songs_staging (
  id UUID PRIMARY KEY,
  object_key TEXT UNIQUE,
  upload_id TEXT UNIQUE,
  uploading_user UUID REFERENCES users ON DELETE CASCADE,
  added_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
