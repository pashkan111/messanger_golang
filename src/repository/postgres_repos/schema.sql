create TYPE chat_type as ENUM (
    'group',
    'dialog',
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(12) UNIQUE NOT NULL,
    chats INTEGER[] DEFAULT '{}'
);

CREATE TABLE chat (
    chat_id SERIAL PRIMARY KEY,
    creator_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE,
    type chat_type NOT NULL DEFAULT 'dialog',
    name VARCHAR(255),
    participants INTEGER[] DEFAULT '{}',
    deleted BOOLEAN DEFAULT false
);

CREATE TABLE message (
    message_id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    chat_id INTEGER REFERENCES chat (chat_id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE dialog (
    chat_id int,
    creator_id INTEGER REFERENCES chat (chat_id) ON DELETE CASCADE,
    participant_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE,
    CONSTRAINT dialog_creator_participant_pk PRIMARY KEY (creator_id, participant_id)
);