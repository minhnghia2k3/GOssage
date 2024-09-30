CREATE TABLE IF NOT EXISTS followers
(
    user_id     bigserial NOT NULL,
    follower_id bigserial NOT NULL,
    created_at  timestamptz DEFAULT NOW(),

    PRIMARY KEY (user_id, follower_id),

    -- Prevent follower followed itself
    CONSTRAINT chk_user_follow_self CHECK (user_id <> follower_id),

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,

    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
);





