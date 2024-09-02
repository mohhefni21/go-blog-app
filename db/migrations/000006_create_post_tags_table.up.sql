CREATE TABLE IF NOT EXISTS posts_tags (
    postTagsId SERIAL PRIMARY KEY,
    tagId INT NOT NULL,
    postId INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tagId)
        REFERENCES tags(tagId)
        ON DELETE CASCADE,
    FOREIGN KEY (postId)
        REFERENCES posts(postId)
        ON DELETE CASCADE
);