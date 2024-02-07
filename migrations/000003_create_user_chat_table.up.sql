-- up.sql
CREATE TABLE IF NOT EXISTS user_chat (
    userid INT REFERENCES users(userid) ON DELETE CASCADE,
    chatid INT REFERENCES chat(chatid) ON DELETE CASCADE,
    PRIMARY KEY (userid, chatid)
);