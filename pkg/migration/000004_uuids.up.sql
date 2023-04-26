create extension if not exists "uuid-ossp";

alter table account
    add column if not exists account_uuid uuid not null default gen_random_uuid() unique;
alter table subscription
    add column if not exists subscription_uuid uuid not null default gen_random_uuid() unique;
alter table subscription
    add column if not exists account_uuid uuid not null default uuid_nil() ;
alter table raw_notification
    add column if not exists raw_notification_uuid uuid not null default gen_random_uuid() unique;

update subscription
set account_uuid = (select account.account_uuid from account where account.account_id = subscription.account_id)
where subscription.account_uuid = uuid_nil();

alter table subscription
    add constraint fk_account_uuid
        foreign key (account_uuid)
            references account (account_uuid);

alter table subscription
    drop constraint fk_account;

alter table subscription
    drop column account_id;

alter table account add unique (withings_user_id)