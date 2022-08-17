ALTER TABLE todo DROP COLUMN draft;
ALTER TABLE todo ADD COLUMN is_active bool DEFAULT false;