-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    status TEXT NOT NULL,
    image_url TEXT,
    video_url TEXT,
    published_at TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    saved_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_news_published_at
ON news (published_at DESC);

CREATE INDEX IF NOT EXISTS idx_news_saved_at
ON news (saved_at DESC);

CREATE INDEX IF NOT EXISTS idx_news_title_tgrm
ON news USING gin (title gin_trgm_ops);

CREATE TABLE IF NOT EXISTS subscribers (
    id UUID PRIMARY KEY,
    chat_id BIGINT NOT NULL UNIQUE,
    username TEXT,
    first_name TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_subscribers_active
ON subscribers (is_active);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_news_published_at;
DROP INDEX IF EXISTS idx_news_saved_at;
DROP INDEX IF EXISTS idx_news_title_tg;
DROP TABLE IF EXISTS news;
DROP INDEX IF EXISTS idx_subscribers_active;
DROP TABLE IF EXISTS subscribers;
-- +goose StatementEnd
