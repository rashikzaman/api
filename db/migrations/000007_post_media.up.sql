BEGIN;

CREATE TABLE post_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    mime_type VARCHAR,
    link VARCHAR,
    post_id UUID,
    user_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_post_post_media FOREIGN KEY (post_id) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_user_post_media FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

COMMIT;