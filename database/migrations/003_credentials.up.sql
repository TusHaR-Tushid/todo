CREATE TABLE credentials(
    cred_id INTEGER REFERENCES users(id),
    user_name varchar(50) NOT NULL ,
    password varchar(50) NOT NULL ,
    expiry TIMESTAMP

)