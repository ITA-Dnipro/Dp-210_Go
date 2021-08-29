CREATE TABLE IF NOT EXISTS roles (
   name varchar(25) PRIMARY KEY,
   description varchar(50)
);

CREATE TABLE IF NOT EXISTS permisions (
   name varchar(25) PRIMARY KEY,
   description varchar(50),
   role_name varchar REFERENCES roles (name)
);

CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   name varchar(25) UNIQUE NOT NULL,
   email varchar(25) UNIQUE NOT NULL,
   password_hash text NOT NULL,
   role varchar REFERENCES roles (name)
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

INSERT INTO roles (name)
   VALUES ('admin'), ('operator'), ('viewer');

CREATE TABLE IF NOT EXISTS data_from_patients (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    first_name varchar(25) NOT NULL,
    last_name varchar(25) NOT NULL,
    email varchar(25) UNIQUE NOT NULL,
    gender varchar(25) NOT NULL,
    birthday_str varchar(10) NOT NULL,
    phone varchar(20) NOT NULL,
    address varchar(150) NOT NULL,
    job_info varchar(150) NOT NULL,
    disability BOOLEAN NOT NULL,
    allergies BOOLEAN NOT NULL,
    reg_day varchar(10) NOT NULL,
    patient_role INT NOT NULL
);

