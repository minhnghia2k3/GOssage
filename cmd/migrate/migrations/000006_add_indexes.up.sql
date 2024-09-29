-- Create the extension and indexes for full-text search
-- Check article: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Using trigram search bc this is text field
CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);

-- Structure data -> not using trigram search
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

-- Dont want fuzzy search (partial search) on username
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- Integer-based fields
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);