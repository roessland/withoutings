create table if not exists account
(
    account_id                   bigserial primary key,
    withings_user_id             varchar                  not null unique,
    withings_access_token        varchar                  not null,
    withings_refresh_token       varchar                  not null,
    withings_access_token_expiry timestamp with time zone not null,
    withings_scopes              varchar                  not null
);

create table if not exists notification_category
(
    appli       integer primary key,
    scope       varchar not null,
    description text not null
);

insert into notification_category (appli, scope, description)
values (1, 'user.metrics', 'New weight-related data'),
       (2, 'user.metrics', 'New temperature related data'),
       (4, 'user.metrics', 'New pressure related data'),
       (16, 'users.activity', 'New activity-related data'),
       (44, 'users.activity', 'New sleep-related data'),
       (46, 'user.info', 'New action on user profile'),
       (50, 'user.sleepevents', 'New bed in event'),
       (51, 'user.sleepevents', 'New bed out event'),
       (52, 'user.sleepevents', 'New inflate done event'),
       (53, 'n/a', 'No account associated'),
       (54, 'user.metrics', 'New ECG data'),
       (55, 'user.metrics', 'ECG measure failed event'),
       (58, 'user.metrics', 'New glucose data')
on conflict do nothing;



create table if not exists subscription
(
    subscription_id bigserial primary key,
    account_id      bigint  not null,
    appli           int     not null,
    callbackurl     varchar not null,
    webhook_secret  varchar not null,
    status          varchar not null,
    comment         varchar not null default '',
    constraint fk_account
        foreign key (account_id)
            references account (account_id)
);

create table if not exists raw_notification
(
    raw_notification_id bigserial primary key,
    source              varchar not null,
    status              varchar not null,
    data                varchar not null
);