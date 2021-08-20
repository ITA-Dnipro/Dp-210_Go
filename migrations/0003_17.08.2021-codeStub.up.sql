CREATE TABLE IF NOT EXISTS password_codes (
    email varchar(25) UNIQUE NOT NULL,
    code varchar(6) NOT NULL)
