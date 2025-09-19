DROP INDEX IF EXISTS idx_chat_owner_id;
DROP INDEX IF EXISTS idx_chat_user_id_chat_id;
DROP INDEX IF EXISTS idx_message_chat_id;
DROP INDEX IF EXISTS idx_message_sender_id;
DROP TABLE IF EXISTS chat_user;
DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS chat;
DROP TYPE IF EXISTS chat_type;