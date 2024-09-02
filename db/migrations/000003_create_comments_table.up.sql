CREATE TABLE IF NOT EXISTS comments (
    commentId SERIAL PRIMARY KEY,
    postId INT NOT NULL,
    userId INT NOT NULL,
    parentId INT,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (postId)
        REFERENCES posts(postId)
        ON DELETE CASCADE,
    FOREIGN KEY (userId)
        REFERENCES users(userId)
        ON DELETE CASCADE
);