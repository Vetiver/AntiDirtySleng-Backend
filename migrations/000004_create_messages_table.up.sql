-- up.sql
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    messagefromuser TEXT NOT NULL,
    chatid INT REFERENCES chat(chatid) ON DELETE CASCADE,
    ownerid INT REFERENCES users(userid) ON DELETE CASCADE
);