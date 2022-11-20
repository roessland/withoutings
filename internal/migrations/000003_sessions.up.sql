create extension pgcrypto;

create table if not exists session
(
    session_id uuid primary key default gen_random_uuid(),
    data       jsonb                    not null default '{}'::jsonb,
    created    timestamp with time zone not null default now(),
    expiry     timestamp with time zone not null default now() + interval '1' month
);