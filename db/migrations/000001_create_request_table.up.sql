CREATE TABLE IF NOT EXISTS request (
    id SERIAL PRIMARY KEY,
    http_request_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT now()  
);
