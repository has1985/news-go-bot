BEGIN;

CREATE TABLE IF NOT EXISTS sources (
    id serial primary key,
    name varchar,
    feed_url varchar,
    priority int,
    created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS articles (
    id serial primary key,
    source_id int,
    title varchar,
    link varchar unique,
    summary text,
    published_at timestamptz,
    created_at timestamptz DEFAULT current_timestamp,
    posted_at timestamptz
);

COMMIT;