CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(12) UNIQUE NOT NULL,
    chats INTEGER[] DEFAULT '{}'
);

CREATE TABLE chat (
    chat_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    participants INTEGER[] DEFAULT '{}'
);

CREATE TABLE message (
    message_id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    chat_id INTEGER REFERENCES chat (chat_id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE
);
