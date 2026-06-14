DROP MATERIALIZED VIEW IF EXISTS active_songs;
DROP TABLE IF EXISTS plays;
DROP TABLE IF EXISTS song_artists;
DROP TABLE IF EXISTS song_contributors;
DROP TABLE IF EXISTS songs;
DROP FUNCTION IF EXISTS orphan_songs_on_album_deletion();
DROP TYPE IF EXISTS upload_source;
DROP TYPE IF EXISTS song_type;
