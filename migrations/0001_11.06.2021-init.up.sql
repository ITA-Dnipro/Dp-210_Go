-- Here we will assume that book has only one author

CREATE TABLE IF NOT EXISTS users
(
    id 	           text NOT NULL PRIMARY KEY,
    name      	   text DEFAULT '',
    email     	   text DEFAULT '',
    role     	   text DEFAULT '',
    password_hash  text DEFAULT ''
);


