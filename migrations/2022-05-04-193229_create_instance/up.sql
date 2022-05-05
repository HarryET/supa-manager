CREATE TABLE IF NOT EXISTS public.instances
(
    id                         uuid primary key not null default gen_random_uuid(),
    nickname                   varchar(50)      null,
    hostname                   varchar(50)      not null,

    studio_container_id        text             null,
    studio_enabled             bool             not null default false,

    kong_container_id          text             not null,
    database_container_id      text             not null,
    gotrue_container_id        text             not null,
    realtime_container_id      text             not null,
    postgrest_container_id      text             not null,
    postgres_meta_container_id text not null
);