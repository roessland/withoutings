alter table raw_notification
    alter column processed_at set not null,
    alter column processed_at set default now(),

    add column if not exists lease_owner  text        not null default '',
    add column if not exists lease_expiry timestamptz not null default now();
