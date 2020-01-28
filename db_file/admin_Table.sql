CREATE TABLE admin (

    AID VARCHAR(12),
    Firstname varchar(50),
    Lastname varchar(50) DEFAULT NULL,
    Password varchar(255) NOT NULL,
    Phonenumber varchar(13) DEFAULT NULL,
    Email varchar(50) NOT NULL
    
)