-- Supports GetNotificationDataByAccountUUIDAndService — turns the per-page
-- (account_uuid, service) lookup from a sequential scan into an index range.
create index if not exists notification_data_account_service_fetched_idx
    on notification_data (account_uuid, service, fetched_at desc);

-- Supports JSONB containment on body.series, used by the sleep detail page
-- to find a Sleep v2 - Getsummary row by an exact session startdate without
-- pulling and parsing every stored summary in Go. jsonb_path_ops only
-- supports @> which is what the lookup uses.
create index if not exists notification_data_body_series_idx
    on notification_data using gin ((data -> 'body' -> 'series') jsonb_path_ops);
