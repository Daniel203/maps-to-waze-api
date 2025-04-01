UPDATE request
SET request_type_id = NULL; -- Or another appropriate default if you have one

DELETE FROM request_type WHERE description IN ('Google Maps Places API', 'Geoapify Static Map API');
