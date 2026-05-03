CREATE TABLE comments_points (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    point_id INT REFERENCES points(id) ON DELETE CASCADE,
    text TEXT,
    image TEXT --base64 format
);