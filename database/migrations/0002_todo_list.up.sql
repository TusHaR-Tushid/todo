CREATE TABLE IF NOT EXISTS todo(
    id SERIAL PRIMARY KEY ,
    created_by INTEGER REFERENCES users(id) NOT NULL ,
    title varchar(50) ,
    description varchar(50) ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    expiring_at TIMESTAMP WITH TIME ZONE,
    is_completed bool DEFAULT FALSE NOT NULL
)