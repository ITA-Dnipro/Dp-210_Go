-- +migrate Down
DROP TABLE IF EXISTS users cascade;

DROP TABLE IF EXISTS roles cascade;

DROP TABLE IF EXISTS users cascade;

DROP TABLE IF EXISTS cards cascade;

DROP TABLE IF EXISTS patients cascade;

DROP TABLE IF EXISTS schedules cascade;

DROP TABLE IF EXISTS doctors cascade;

DROP TABLE IF EXISTS appointments cascade;