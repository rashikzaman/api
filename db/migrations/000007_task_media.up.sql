BEGIN;

CREATE TABLE task_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    mime_type VARCHAR,
    link VARCHAR,
    task_id UUID,
    user_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_task_task_media FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_user_task_media FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

COMMIT;