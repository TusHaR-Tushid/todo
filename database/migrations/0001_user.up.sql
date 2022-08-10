CREATE TABLE users(
    id SERIAL PRIMARY KEY ,
    name varchar(50) NOT NULL,
    archived_at timestamp with time zone
)