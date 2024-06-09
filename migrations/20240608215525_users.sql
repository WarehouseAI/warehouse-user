-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.users (
  id public.xid NOT NULL DEFAULT xid(),
  firstname VARCHAR(80) NOT NULL,
  lastname VARCHAR(80) NOT NULL,
  username VARCHAR(80) NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  hash TEXT NOT NULL,
  role INTEGER NOT NULL DEFAULT 1,
  verified BOOLEAN NOT NULL DEFAULT false,
  created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
  updated_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
);
ALTER TABLE public.user
ADD CONSTRAINT user_pkey PRIMARY KEY (user_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE public.users;