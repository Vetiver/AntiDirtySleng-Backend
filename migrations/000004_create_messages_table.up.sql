-- up.sql
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    messagefromuser TEXT NOT NULL,
    chatid UUID REFERENCES chat(chatid) ON DELETE CASCADE,
    ownerid UUID REFERENCES users(userid) ON DELETE CASCADE
);