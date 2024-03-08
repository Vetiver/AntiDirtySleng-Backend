-- up.sql

CREATE TABLE IF NOT EXISTS user_chat (
    userid UUID REFERENCES users(userid) ON DELETE CASCADE,
    chatid UUID REFERENCES chat(chatid) ON DELETE CASCADE,
    PRIMARY KEY (userid, chatid)
);

INSERT INTO user_chat (userid, chatid) VALUES ((SELECT userid FROM users WHERE username = 'Попкинс'), (SELECT chatid FROM chat WHERE chatName = 'Чат Попкинса и Ванька'));
INSERT INTO user_chat (userid, chatid) VALUES ((SELECT userid FROM users WHERE username = 'Ванька'), (SELECT chatid FROM chat WHERE chatName = 'Чат Попкинса и Ванька'));