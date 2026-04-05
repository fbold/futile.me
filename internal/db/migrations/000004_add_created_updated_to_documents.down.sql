ALTER TABLE documents
DROP COLUMN created,
DROP COLUMN updated;

DROP TRIGGER trigger_update_updated ON documents;
