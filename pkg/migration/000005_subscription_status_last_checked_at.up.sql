
alter table subscription
    add column if not exists status_last_checked_at timestamptz not null default now() - '25 hours'::interval;