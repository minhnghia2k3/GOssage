CREATE TABLE IF NOT EXISTS user_invitation (
    user_id bigserial NOT NULL,
    token bytea NOT NULL,
    expiry timestamptz NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id)
);