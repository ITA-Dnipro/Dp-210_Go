CREATE TABLE IF NOT EXISTS roles (
   id serial PRIMARY KEY,
   name varchar NOT NULL,
   description varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS permisions (
   id serial PRIMARY KEY,
   name varchar NOT NULL,
   description varchar NOT NULL,
   role_id int REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   name varchar NOT NULL,
   email varchar NOT NULL,
   password_hash varchar NOT NULL,
   role_id int REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS cards (
   id uuid PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS patients (
   id uuid PRIMARY KEY,
   full_name varchar(50),
   user_id uuid UNIQUE REFERENCES users (id),
   card_id uuid UNIQUE REFERENCES cards (id)
);

CREATE TABLE IF NOT EXISTS schedules (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   name varchar,
   description varchar NOT NULL,
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
   full_name varchar(50),
   speciality varchar NOT NULL,
   schedule_id uuid REFERENCES schedules (id)
);

CREATE TABLE IF NOT EXISTS appointments (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   doctor_id uuid REFERENCES doctors (id),
   patient_id uuid REFERENCES patients (id),
   reason varchar,
   result varchar,
   timeRange tsrange
);
