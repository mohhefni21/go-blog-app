CREATE TABLE IF NOT EXISTS posts (
    postId SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    cover VARCHAR(100),
    title VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    excerpt TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (userId)
        REFERENCES users(userId)
        ON DELETE CASCADE
);