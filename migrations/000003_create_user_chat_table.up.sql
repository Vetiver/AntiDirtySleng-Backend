-- up.sql

CREATE TABLE IF NOT EXISTS user_chat (
    userid UUID REFERENCES users(userid) ON DELETE CASCADE,
    chatid UUID REFERENCES chat(chatid) ON DELETE CASCADE,
    PRIMARY KEY (userid, chatid)
);