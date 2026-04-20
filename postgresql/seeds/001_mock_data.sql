-- Seed mock data for development/testing.
-- Run: psql -h localhost -U blog -d blog -f postgresql/seeds/001_mock_data.sql

INSERT INTO users (uuid, nickname, name, surname, avatar_url) VALUES
    ('11111111-1111-1111-1111-111111111111', 'john_doe', 'John', 'Doe', 'https://example.com/avatars/john.png'),
    ('22222222-2222-2222-2222-222222222222', 'jane_smith', 'Jane', 'Smith', 'https://example.com/avatars/jane.png'),
    ('33333333-3333-3333-3333-333333333333', 'bob_wilson', 'Bob', 'Wilson', 'https://example.com/avatars/bob.png')
ON CONFLICT (uuid) DO NOTHING;

INSERT INTO posts (uuid, author_uuid, content_text, content_image_urls) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111',
     'Hello everyone! This is my first blog post about Go programming.',
     '{"https://example.com/images/go-gopher.png"}'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222',
     'Just finished reading a great book on microservices architecture!',
     '{}'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333',
     'Check out this beautiful sunset I captured!',
     '{"https://example.com/images/sunset1.png","https://example.com/images/sunset2.png"}')
ON CONFLICT (uuid) DO NOTHING;
