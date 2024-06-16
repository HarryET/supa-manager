CREATE TABLE IF NOT EXISTS versions
(
    id         uuid      not null primary key default gen_random_uuid(),
    service_id text      not null,
    image      text      not null,
    tag        text      not null,
    created_at timestamp not null             default now(),
    created_by text      not null
);