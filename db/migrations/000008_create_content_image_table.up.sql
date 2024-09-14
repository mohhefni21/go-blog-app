CREATE TABLE IF NOT EXISTS content_image (
    id_content_image SERIAL PRIMARY KEY,
    id_post INT REFERENCES posts(post_id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
