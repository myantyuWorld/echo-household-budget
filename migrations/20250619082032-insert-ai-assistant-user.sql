-- +migrate Up
-- AIアシスタントユーザーを登録
INSERT INTO
  user_accounts (id, user_id, name)
VALUES
  (
    0,
    -- システム用のID
    'ai-assistant',
    'AI家計アシスタント'
  );

-- +migrate Down
delete from
  user_accounts
where
  userID = 0;
