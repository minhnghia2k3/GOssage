CREATE TABLE IF NOT EXISTS comments
(
    id         bigserial PRIMARY KEY,
    user_id    bigserial NOT NULL,
    post_id    bigserial NOT NULL,
    content    varchar   NOT NULL,
    created_at timestamptz default now(),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE
);

CREATE INDEX ON "comments" ("content");
