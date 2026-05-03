CREATE TABLE comments_lifehacks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    lifehack_id INT REFERENCES lifehacks(id) ON DELETE CASCADE,
    text TEXT,
    image TEXT
);