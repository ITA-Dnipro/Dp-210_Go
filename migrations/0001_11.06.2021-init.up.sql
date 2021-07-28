-- Here we will assume that book has only one author

CREATE TABLE IF NOT EXISTS users
(
    id        text NOT NULL PRIMARY KEY,
    first_name text DEFAULT '',
    last_name  text DEFAULT ''
);


