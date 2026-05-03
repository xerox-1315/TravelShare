CREATE TABLE votes_lifehacks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    lifehack_id INT REFERENCES lifehacks(id) ON DELETE CASCADE,
    type VARCHAR(50) -- like / dislike
);