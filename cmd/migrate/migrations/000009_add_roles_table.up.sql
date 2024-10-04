CREATE TABLE IF NOT EXISTS roles
(
    id          bigserial    NOT NULL PRIMARY KEY,
    name        varchar(255) NOT NULL,
    level       int          NOT NULL DEFAULT 0,
    description text
);

INSERT INTO roles(name, level, description)
VALUES ('user', 1, 'A user can create posts and comments');

INSERT INTO roles(name, level, description)
VALUES ('moderator', 2, 'A moderator can update other user posts');

INSERT INTO roles(name, level, description)
VALUES ('admin', 3, 'An admin can update and delete other user posts')