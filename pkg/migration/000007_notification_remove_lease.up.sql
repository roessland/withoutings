alter table raw_notification
    -- Should be null if the notification has not been processed yet.
    alter column processed_at drop not null,
    alter column processed_at drop default,

    -- Don't need these, going for a queue instead.
    drop column if exists lease_owner,
    drop column if exists lease_expiry;