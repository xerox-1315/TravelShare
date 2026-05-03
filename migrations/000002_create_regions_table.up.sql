CREATE TABLE regions (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    polygon GEOGRAPHY(POLYGON, 4326) NOT NULL
);