
CREATE TYPE status_enum AS ENUM ('submitted', 'running', 'completed', 'failing', 'failed');

CREATE TABLE IF NOT EXISTS job (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    scheduled_time TIMESTAMP NOT NULL,
    metadata JSONB,
    partition_id TEXT,
    job_status status_enum DEFAULT 'submitted',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_job_user_id_status ON job(user_id, job_status);
CREATE INDEX idx_job_user_id_scheduled_time ON job(user_id, scheduled_time);