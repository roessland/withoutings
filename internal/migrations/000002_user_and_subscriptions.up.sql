begin transaction;

create table if not exists account
(
    account_id                   bigserial primary key,
    withings_user_id             varchar not null,
    withings_access_token        varchar not null,
    withings_refresh_token       varchar not null,
    withings_access_token_expiry timestamp with time zone not null,
    withings_scopes              varchar not null
);

create table if not exists subscription
(
    subscription_id serial primary key,
    account_id      bigint,
    appli           int     not null,
    callbackurl     varchar not null,
    comment         varchar not null default '',
    constraint fk_account
        foreign key (account_id)
            references account (account_id)
);

commit;