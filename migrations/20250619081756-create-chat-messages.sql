-- +migrate Up
CREATE TYPE message_type AS ENUM ('user', 'ai');

CREATE TABLE chat_messages (
  id SERIAL PRIMARY KEY,
  household_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  -- null不可に変更
  content TEXT NOT NULL,
  message_type message_type NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (household_id) REFERENCES household_books(id),
  FOREIGN KEY (user_id) REFERENCES user_accounts(id)
);

CREATE INDEX idx_chat_messages_household_id ON chat_messages(household_id);

CREATE INDEX idx_chat_messages_user_id ON chat_messages(user_id);

CREATE INDEX idx_chat_messages_message_type ON chat_messages(message_type);

-- +migrate Down
DROP TABLE chat_messages;

DROP TYPE message_type;
