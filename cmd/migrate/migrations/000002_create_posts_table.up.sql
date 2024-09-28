CREATE TABLE IF NOT EXISTS posts
(
    id         bigserial PRIMARY KEY,
    user_id    bigserial    NOT NULL,
    title      VARCHAR(255) NOT NULL,
    content    VARCHAR(255) NOT NULL,
    tags       VARCHAR(100)[],
    created_at timestamptz default now(),
    updated_at timestamptz default now()
)