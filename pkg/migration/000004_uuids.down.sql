alter table subscription
    drop constraint if exists fk_account_uuid;

alter table account
    drop column if exists account_uuid;

alter table subscription
    drop column if exists subscription_uuid;

alter table subscription
    drop column if exists account_uuid;

alter table raw_notification
    drop column if exists raw_notification_uuid;
