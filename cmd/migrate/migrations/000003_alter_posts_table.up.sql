CREATE INDEX ON "posts" ("title");

CREATE INDEX ON "posts" ("content");

CREATE INDEX ON "posts" ("tags");

ALTER TABLE posts
    ADD CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users (id);

