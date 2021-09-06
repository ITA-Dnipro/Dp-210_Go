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

CREATE TABLE IF NOT EXISTS patients (
   id uuid PRIMARY KEY REFERENCES users (id),
   first_name varchar(25),
   last_name varchar(25)
);

CREATE TABLE IF NOT EXISTS doctors (
   id uuid PRIMARY KEY REFERENCES users (id),
   first_name varchar(25),
   last_name varchar(25),
   speciality varchar NOT NULL,
   start_at timestamp,
   end_at timestamp
);

CREATE TABLE IF NOT EXISTS appointments (
   appointment_id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   doctor_id uuid REFERENCES doctors (id),
   patient_id uuid REFERENCES patients (id),
   reason varchar(150),
   time_range tstzrange,
   EXCLUDE USING gist (doctor_id WITH =, patient_id WITH =, time_range WITH &&)
);

INSERT INTO roles (name)
   VALUES ('admin'), ('operator'), ('viewer'), ('doctor'), ('patient')
