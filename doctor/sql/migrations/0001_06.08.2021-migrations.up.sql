CREATE TABLE IF NOT EXISTS doctors (
   id uuid PRIMARY KEY,
   first_name varchar(25),
   last_name varchar(25),
   speciality varchar NOT NULL,
   start_at timestamp,
   end_at timestamp
);
