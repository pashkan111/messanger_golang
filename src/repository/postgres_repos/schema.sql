CREATE TYPE message_type AS ENUM ('TEXT', 'IMAGE', 'FILE');

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(12) UNIQUE NOT NULL,
    chats INTEGER[] DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS dialog (
    dialog_id BIGSERIAL UNIQUE,
    deleted_for INTEGER[],
    is_deleted BOOLEAN DEFAULT false,
    creator_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE,
    receiver_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE,
    CONSTRAINT dialog_creator_participant_pk PRIMARY KEY (creator_id, receiver_id)
);

CREATE TABLE IF NOT EXISTS dialog_message (
    dialog_message_id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    message_type message_type DEFAULT 'TEXT',
    link VARCHAR(255),
    dialog_id INTEGER REFERENCES dialog (dialog_id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS chat (
--     chat_id SERIAL PRIMARY KEY,
--     creator_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE,
--     type chat_type NOT NULL DEFAULT 'dialog',
--     name VARCHAR(255),
--     participants INTEGER[] DEFAULT '{}',
--     deleted BOOLEAN DEFAULT false
-- );