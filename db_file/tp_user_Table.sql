CREATE TABLE tp_users (
    uid varchar(12) PRIMARY KEY,
    email varchar(50) NOT NULL,
    password varchar(255) NOT NULL,
    origin varchar(50) NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(uid)
);