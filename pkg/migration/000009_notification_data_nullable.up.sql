ALTER TABLE notification
    ALTER COLUMN fetched_at DROP NOT NULL,
    ALTER COLUMN data DROP NOT NULL,
    ADD COLUMN data_status TEXT NOT NULL DEFAULT 'awaiting_fetch';

UPDATE notification SET fetched_at = NULL;