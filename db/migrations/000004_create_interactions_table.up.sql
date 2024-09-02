CREATE TABLE IF NOT EXISTS interactions (
    interactionId SERIAL PRIMARY KEY,
    postId INT NOT NULL,
    userId INT NOT NULL,
    type INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (postId)
        REFERENCES posts(postId)
        ON DELETE CASCADE,
    FOREIGN KEY (userId)
        REFERENCES users(userId)
        ON DELETE CASCADE
);