CREATE TABLE users (
    uuid       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nickname   TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    surname    TEXT NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT ''
);

CREATE TABLE posts (
    uuid               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_uuid        UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ,
    content_text       TEXT NOT NULL DEFAULT '',
    content_image_urls TEXT[] NOT NULL DEFAULT '{}'
);

CREATE INDEX idx__posts__created_at ON posts(created_at DESC);
CREATE INDEX idx__posts__author_uuid ON posts(author_uuid);
