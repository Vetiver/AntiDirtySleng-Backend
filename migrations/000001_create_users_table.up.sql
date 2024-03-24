-- up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    userid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(225) NOT NULL,
    email VARCHAR(225) NOT NULL UNIQUE,
    password VARCHAR(225) NOT NULL,
    isadmin BOOLEAN DEFAULT false,
    description VARCHAR(255) DEFAULT 'Кто я?',
    avatar VARCHAR(255)
);

INSERT INTO users (username, email, password) VALUES ('Попкинс', 'popkins@gmail.com', '$2a$10$gSaf3tjZ6jJiS3rP/pGVl.zPJzEaLtrsEavXjYeJ6xzu5GOaeDkRW');--пароль 11111111
INSERT INTO users (username, email, password) VALUES ('Ванька', 'vanek@gmail.com', '$2a$10$gSaf3tjZ6jJiS3rP/pGVl.zPJzEaLtrsEavXjYeJ6xzu5GOaeDkRW');--пароль 11111111
INSERT INTO users (username, email, password) VALUES ('Дима', 'vinek@gmail.com', '$2a$10$gSaf3tjZ6jJiS3rP/pGVl.zPJzEaLtrsEavXjYeJ6xzu5GOaeDkRW');--пароль 11111111
INSERT INTO users (username, email, password) VALUES ('Тимофей', 'vonek@gmail.com', '$2a$10$gSaf3tjZ6jJiS3rP/pGVl.zPJzEaLtrsEavXjYeJ6xzu5GOaeDkRW');--пароль 11111111
INSERT INTO users (username, email, password) VALUES ('Санек Морозов', 'sashamoroz0412@gmail.com', '$2a$10$gSaf3tjZ6jJiS3rP/pGVl.zPJzEaLtrsEavXjYeJ6xzu5GOaeDkRW');--пароль 11111111
INSERT INTO users (username, email, password) VALUES ('Бага', 'vznek@gmail.com', '$2a$10$gSaf3tjZ6jJiS3rP/pGVl.zPJzEaLtrsEavXjYeJ6xzu5GOaeDkRW');--пароль 11111111