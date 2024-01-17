CREATE TABLE IF NOT EXISTS phones (
    id BIGSERIAL PRIMARY KEY,
    phone_number text NOT NULL,
    provider text NOT NULL
);