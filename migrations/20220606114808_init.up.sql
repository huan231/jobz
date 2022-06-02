CREATE TABLE IF NOT EXISTS job
(
    id           TEXT PRIMARY KEY,
    cron_job_id  TEXT    NOT NULL,
    status       INTEGER NOT NULL,
    created_at   INTEGER NOT NULL,
    updated_at   INTEGER NOT NULL,
    completed_at INTEGER
);
