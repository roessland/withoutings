alter table raw_notification
    drop column if exists received_at,
    drop column if exists processed_at,
    drop column if exists lease_owner,
    drop column if exists lease_expiry;

drop table if exists notification;