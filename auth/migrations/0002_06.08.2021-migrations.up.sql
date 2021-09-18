CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   email varchar(25) UNIQUE NOT NULL,
   password_hash text NOT NULL
);
