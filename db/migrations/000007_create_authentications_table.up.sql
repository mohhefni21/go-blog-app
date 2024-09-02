CREATE TABLE IF NOT EXISTS authentications (
    authenticationId SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    accessToken TEXT NOT NULL,
    refreshToken TEXT NOT NULL,
    accessTokenExpiresAt TIMESTAMP NOT NULL,
    refreshTokenExpiresAt TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (userId)
        REFERENCES users(userId)
        ON DELETE CASCADE
);
