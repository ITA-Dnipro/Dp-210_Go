CREATE TABLE IF NOT EXISTS roles (
   id serial PRIMARY KEY,
   name varchar(25) NOT NULL,
   description varchar(50)
);

CREATE TABLE IF NOT EXISTS permisions (
   id serial PRIMARY KEY,
   name varchar(25) NOT NULL,
   description varchar(50),
   role_id int REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   name varchar(25) UNIQUE NOT NULL,
   email varchar(25) UNIQUE NOT NULL,
   password_hash text NOT NULL,
   role_id int REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS cards (
   id uuid PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS patients (
   id uuid PRIMARY KEY,
   first_name varchar(25),
   last_name varchar(25),
   user_id uuid UNIQUE REFERENCES users (id),
   card_id uuid UNIQUE REFERENCES cards (id)
);

CREATE TABLE IF NOT EXISTS schedules (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   name varchar(25) NOT NULL,
   description varchar(150),
   sun tsrange,
   mon tsrange,
   tue tsrange,
   wed tsrange,
   thu tsrange,
   fri tsrange,
   sat tsrange
);

CREATE TABLE IF NOT EXISTS doctors (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   user_id uuid REFERENCES users (id),
   first_name varchar(25),
   last_name varchar(25),
   speciality varchar NOT NULL,
   schedule_id uuid REFERENCES schedules (id)
);

CREATE TABLE IF NOT EXISTS appointments (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   doctor_id uuid REFERENCES doctors (id),
   patient_id uuid REFERENCES patients (id),
   reason varchar(150),
   result varchar,
   timeRange tsrange
);
