CREATE TABLE IF NOT EXISTS appointments (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
   doctor_id uuid,
   patient_id uuid,
   reason varchar(150),
   time_range tstzrange,
   EXCLUDE USING gist (doctor_id WITH =, time_range WITH &&)
);
