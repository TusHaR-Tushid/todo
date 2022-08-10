

CREATE TABLE IF NOT EXISTS sessions (
    id uuid primary key default gen_random_uuid() not null ,
    user_id INTEGER REFERENCES users(id) NOT NULL ,
    expires_at TIMESTAMP


) ;