CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.migrations
(
    id         text        not null primary key,
    note       text,
    applied_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS public.accounts
(
    id            uuid        not null default gen_random_uuid(),

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
    slug          text        not null default gen_random_uuid()::text,

    name          text        not null,
    kind          text        not null, --help?!

    billing_email text        not null default 'selfhosted@supa-manager.io',

    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),

    primary key (slug)
);

CREATE TABLE IF NOT EXISTS public.organization_membership
(
    organization_slug text        not null,
    account_id        uuid        not null,

    role              text        not null, -- todo does this need a change?

    created_at        timestamptz not null default now(),
    updated_at        timestamptz not null default now(),

    primary key (organization_slug, account_id)
);

ALTER TABLE public.organization_membership
    ADD CONSTRAINT fk_membership_org FOREIGN KEY (organization_slug) REFERENCES organizations (slug);

ALTER TABLE public.organization_membership
    ADD CONSTRAINT fk_membership_account FOREIGN KEY (account_id) REFERENCES accounts (id);

CREATE TABLE IF NOT EXISTS public.project
(
    project_ref       text        not null,

    project_name      text        not null,
    organization_slug text        not null,

    status            text        not null, -- make this an enum

    cloud_provider    text        not null default 'k8s',
    region            text        not null default 'mars-1',

    jwt_secret        text        not null,

    created_at        timestamptz not null default now(),
    updated_at        timestamptz not null default now(),

    primary key (project_ref)
);

