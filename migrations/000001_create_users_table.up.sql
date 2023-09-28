CREATE TABLE
    IF NOT EXISTS users (
        id VARCHAR(255) PRIMARY KEY,
        unique_id VARCHAR(255) UNIQUE.
        nickname VARCHAR(255) NOT NULL,
        points INT NOT NULL DEFAULT 0,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );