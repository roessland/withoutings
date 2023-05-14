alter table raw_notification
    -- When the webhook request was receieved.
    add column if not exists received_at  timestamptz not null default now(),

    -- When the webhook request was processed and data was fetched from Withings API.
    add column if not exists processed_at timestamptz not null default now(),

    -- Process that owns a lease on this row.
    add column if not exists lease_owner  text        not null default '',

    -- When the lease expires.
    add column if not exists lease_expiry timestamptz not null default now();


-- Domain model for a notification.
create table if not exists notification
(
    -- Primary key
    notification_uuid     uuid primary key,

    -- Account that owns this notification.
    -- So we can delete everything when an account is deleted.
    account_uuid          uuid        not null,
    constraint fk_account foreign key (account_uuid) references account (account_uuid),

    -- When the raw notification was received
    received_at           timestamptz not null default now(),

    -- URLencoded query string from webhook POST parameters converted to JSON
    params                jsonb       not null,

    -- Data retrieved from Withings API using the webhook parameters.
    data                  jsonb       not null,

    -- When the data was fetched from Withings API
    fetched_at            timestamptz not null default now(),

    -- Reference to the raw notification that triggered this notification
    -- So that we can delete that too, when this notification is deleted.
    raw_notification_uuid uuid        not null,

    -- Source IP address of the webhook POST request
    source                text        not null
);