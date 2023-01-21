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

create table if not exists notification_category
(
    appli integer primary key,
    scope varchar,
    description  text
);



insert into notification_category  (appli, scope, description)
values
    (1, 'user.metrics', 'New weight-related data'),
    (2, 'user.metrics', 'New temperature related data'),
    (4, 'user.metrics', 'New pressure related data'),
    (16, 'user.metrics', 'New activity-related data'),
    (44, 'user.metrics', 'New sleep-related data'),
    (46, 'user.metrics', 'New action on user profile'),
    (50, 'user.metrics', 'New bed in event'),
    (51, 'user.metrics', 'New bed out event'),
    (52, 'user.metrics', 'New inflate done event'),
    (53, 'user.metrics', 'No account associated'),
    (54, 'user.metrics', 'New ECG data'),
    (55, 'user.metrics', 'ECG measure failed event'),
    (58, 'user.metrics', 'New glucose data')
on conflict do nothing;



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