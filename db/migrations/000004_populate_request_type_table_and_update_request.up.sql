INSERT INTO request_type (description) VALUES ('Google Maps Places API');
INSERT INTO request_type (description) VALUES ('Geoapify Static Map API');

UPDATE request
SET request_type_id = (SELECT id FROM request_type WHERE description = 'Google Maps Places API' LIMIT 1);
