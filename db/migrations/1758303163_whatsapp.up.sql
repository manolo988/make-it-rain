-- ENUM chat_type
CREATE TYPE chat_type AS ENUM ('group', 'private');

-- TABLE Chat
CREATE TABLE IF NOT EXISTS chat (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    chat_type chat_type NOT NULL DEFAULT 'private',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_chat_owner_id ON chat(owner_id);

-- TABLE chat_user
CREATE TABLE IF NOT EXISTS chat_user (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chat(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_chat_user UNIQUE (chat_id, user_id)
);
CREATE INDEX idx_chat_user_id_chat_id ON chat_user(user_id, chat_id);

-- TABLE Message
CREATE TABLE IF NOT EXISTS message (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chat(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_message_chat_id ON message(chat_id, created_at);
CREATE INDEX idx_message_sender_id ON message(sender_id);