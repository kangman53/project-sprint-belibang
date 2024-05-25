CREATE UNIQUE INDEX IF NOT EXISTS unique_email 
    ON users(email, role);
CREATE UNIQUE INDEX IF NOT EXISTS unique_username 
    ON users(username);
CREATE INDEX IF NOT EXISTS index_users_id
    ON users (id);
CREATE INDEX IF NOT EXISTS index_users_name
    ON users USING HASH(lower(username));
CREATE INDEX IF NOT EXISTS index_users_role
    ON users (role);