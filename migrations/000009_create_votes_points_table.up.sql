CREATE TABLE votes_points (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    point_id INT REFERENCES points(id) ON DELETE CASCADE,
    type VARCHAR(50) -- like / dislike
);