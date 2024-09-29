CREATE TABLE IF NOT EXISTS users
(
    id         bigserial PRIMARY KEY,
    email      citext UNIQUE       NOT NULL, -- Case insensitive data type
    username   varchar(255) UNIQUE NOT NULL,
    password   bytea               NOT NULL,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);
