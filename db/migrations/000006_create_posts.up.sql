BEGIN;

CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    title TEXT NOT NULL CHECK (title <> ''),
    description TEXT NOT NULL CHECK (description <> ''),
    required_volunteers_count INT DEFAULT 0,
    required_skills TEXT [],
    latitude decimal NOT NULL,
    longitude decimal NOT NULL,
    location GEOMETRY (Point, 4326) NOT NULL,
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_post FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_category_post FOREIGN KEY (category_id) REFERENCES categories (id) ON UPDATE CASCADE ON DELETE CASCADE
);

create trigger set_timestamp_posts before
update
    on posts for each row execute procedure trigger_set_updated_at_timestamp();

COMMIT;