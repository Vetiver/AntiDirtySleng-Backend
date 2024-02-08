-- up.sql
CREATE TABLE IF NOT EXISTS chat (
    chatid SERIAL PRIMARY KEY,
    chatName VARCHAR(30) NOT NULL,
    owner INT,
    FOREIGN KEY (owner) REFERENCES users(userid) ON DELETE SET NULL
);