CREATE TABLE lifehacks (
    id SERIAL PRIMARY KEY,
    coords GEOGRAPHY(POINT, 4326) NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    region_id INT REFERENCES regions(id) ON DELETE SET NULL,
    season_id INT REFERENCES seasons(id) ON DELETE SET NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    image TEXT 
);