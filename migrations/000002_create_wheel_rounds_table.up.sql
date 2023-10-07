CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    IF NOT EXISTS wheel_rounds (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        outcome_idx INT NOT NULL,
        round_start_time TIMESTAMP NOT NULL DEFAULT NOW (),
        spin_start_time TIMESTAMP NOT NULL,
        round_end_time TIMESTAMP NOT NULL,
        server_seed VARCHAR(64) NOT NULL,
        client_seed VARCHAR(255) NOT NULL,
        blinded_server_seed VARCHAR(64) NOT NULL,
        status VARCHAR(10) NOT NULL,
        nonce INT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS wheel_round_entries (
        round_id UUID,
        user_id VARCHAR(255),
        wager FLOAT NOT NULL,
        pick VARCHAR(255) NOT NULL,
        entered_at TIMESTAMP NOT NULL DEFAULT NOW (),
        PRIMARY KEY (round_id, user_id),
        FOREIGN KEY (round_id) REFERENCES wheel_rounds (id) ON UPDATE CASCADE ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
    );