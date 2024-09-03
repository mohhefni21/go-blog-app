CREATE TABLE IF NOT EXISTS posts_tags (
    postTags_id SERIAL PRIMARY KEY,
    tag_id INT NOT NULL,
    post_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tag_id)
        REFERENCES tags(tag_id)
        ON DELETE CASCADE,
    FOREIGN KEY (post_id)
        REFERENCES posts(post_id)
        ON DELETE CASCADE
);