ALTER TABLE request
ALTER COLUMN request_type_id SET NOT NULL;

ALTER TABLE request
ADD CONSTRAINT fk_request_request_type
FOREIGN KEY (request_type_id)
REFERENCES request_type(id);
