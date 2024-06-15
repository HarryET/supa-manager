CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.migrations
(
    id         text        not null primary key,
    note       text,
    applied_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS public.accounts
(
    id            serial      not null,
    gotrue_id     text        not null default gen_random_uuid()::text,

    email         text        not null,
    password_hash text        not null,

    username      text        not null,

    first_name    text,
    last_name     text,

    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),

    primary key (id)
);

CREATE TABLE IF NOT EXISTS public.organizations
(
    id            serial      not null,
    slug          text        not null default gen_random_uuid()::text,

    name          text        not null,

    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),

    primary key (id)
);

CREATE TABLE IF NOT EXISTS public.organization_membership
(
    organization_id int         not null,
    account_id      int         not null,

    role            text        not null, -- todo does this need a change?

    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),

    primary key (organization_id, account_id)
);

ALTER TABLE public.organization_membership
    ADD CONSTRAINT fk_membership_org FOREIGN KEY (organization_id) REFERENCES organizations (id);

ALTER TABLE public.organization_membership
    ADD CONSTRAINT fk_membership_account FOREIGN KEY (account_id) REFERENCES accounts (id);

CREATE TABLE IF NOT EXISTS public.project
(
    id              serial      not null,
    project_ref     text        not null,

    project_name    text        not null,
    organization_id int         not null,

    status          text        not null, -- make this an enum

    cloud_provider  text        not null default 'k8s',
    region          text        not null default 'mars-1',

    jwt_secret      text        not null,

    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),

    primary key (id)
);

ALTER TABLE public.project
    ADD CONSTRAINT fk_project_org FOREIGN KEY (organization_id) REFERENCES organizations (id);
