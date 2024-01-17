CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email text NOT NULL,
    name text NULL,
    password text NOT NULL
);