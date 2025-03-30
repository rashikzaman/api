BEGIN;

CREATE TABLE user_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    latitude decimal NOT NULL,
    longitude decimal NOT NULL,
    location GEOMETRY (Point, 4326) NOT NULL,
    formatted_address TEXT,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_task FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

create trigger set_timestamp_user_locations before
update
    on tasks for each row execute procedure trigger_set_updated_at_timestamp();

COMMIT;