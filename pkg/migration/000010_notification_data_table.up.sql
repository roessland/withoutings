create table if not exists notification_data
(
    notification_data_uuid uuid primary key,

    account_uuid           uuid        not null,
    constraint fk_account foreign key (account_uuid) references account (account_uuid),

    notification_uuid      uuid        not null,
    constraint fk_notification foreign key (notification_uuid) references notification (notification_uuid),

    service                text        not null,

    data                   jsonb       not null,

    fetched_at             timestamptz not null default now()
);

alter table notification_data add unique (notification_uuid, service);

alter table notification drop column if exists data;
