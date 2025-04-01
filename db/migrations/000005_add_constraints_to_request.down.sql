ALTER TABLE request
DROP CONSTRAINT IF EXISTS fk_request_request_type;

ALTER TABLE request
ALTER COLUMN request_type_id DROP NOT NULL;
