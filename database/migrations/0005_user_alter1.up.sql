CREATE TABLE IF NOT EXISTS users(
                      id SERIAL PRIMARY KEY ,
                      name varchar(50) NOT NULL,
                      email TEXT NOT NULL ,
                      password TEXT NOT NULL ,
                      age INTEGER,
                      gender TEXT ,
                      address TEXT,
                      created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
                      updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
                      archived_at timestamp with time zone
)
