CREATE TABLE IF NOT EXISTS followers
(
    user_id    bigserial NOT NULL,
    follow_id  bigserial NOT NULL,
    created_at timestamptz DEFAULT NOW(),

    PRIMARY KEY (user_id, follow_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (follow_id) REFERENCES users (id) ON DELETE CASCADE
);





