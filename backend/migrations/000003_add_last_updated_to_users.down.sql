DROP TRIGGER IF EXISTS users_set_last_updated ON users;
DROP FUNCTION IF EXISTS set_last_updated;
ALTER TABLE users DROP COLUMN last_updated;
