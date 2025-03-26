BEGIN;

CREATE TABLE user_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    task_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_user_tasks FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_task_user_tasks FOREIGN KEY (task_id) REFERENCES tasks (id) ON UPDATE CASCADE ON DELETE CASCADE
);

create trigger set_timestamp_user_tasks before
update
    on tasks for each row execute procedure trigger_set_updated_at_timestamp();

COMMIT;