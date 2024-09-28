CREATE TABLE IF NOT EXISTS comments
(
    id         bigserial PRIMARY KEY,
    user_id    bigserial NOT NULL,
    post_id    bigserial NOT NULL,
    content    varchar   NOT NULL,
    created_at timestamptz default now()
);

CREATE INDEX ON "comments" ("content");

ALTER TABLE "comments"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments"
    ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
