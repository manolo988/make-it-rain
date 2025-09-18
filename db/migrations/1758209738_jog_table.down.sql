DROP INDEX IF EXISTS idx_job_user_id_status;
DROP INDEX IF EXISTS idx_job_user_id_scheduled_time;
DROP TABLE IF EXISTS job;
DROP TYPE IF EXISTS status_enum;