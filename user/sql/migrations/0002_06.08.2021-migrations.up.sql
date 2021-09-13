CREATE TABLE IF NOT EXISTS roles (
   name varchar(25) PRIMARY KEY,
   description varchar(50)
);

CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   name varchar(25) UNIQUE NOT NULL,
   email varchar(25) UNIQUE NOT NULL,
   password_hash text NOT NULL,
   role varchar REFERENCES roles (name)
);

INSERT INTO roles (name)
   VALUES ('admin'), ('operator'), ('viewer'), ('doctor'), ('patient')
