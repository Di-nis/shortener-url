CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original VARCHAR(255) UNIQUE NOT NULL,
    short VARCHAR(8) UNIQUE NOT NULL,
);

CREATE INDEX idx_short ON urls(short);