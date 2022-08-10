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
);
CREATE TABLE IF NOT EXISTS todo(
                                   id SERIAL PRIMARY KEY ,
                                   created_by INTEGER REFERENCES users(id) NOT NULL ,
                                   title varchar(50) ,
                                   description varchar(50) ,
                                   updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
                                   created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
                                   expiring_at TIMESTAMP WITH TIME ZONE,
                                   is_completed bool DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
                                        id uuid primary key default gen_random_uuid() not null ,
                                        user_id INTEGER REFERENCES users(id) NOT NULL ,
                                        expires_at TIMESTAMP


) ;