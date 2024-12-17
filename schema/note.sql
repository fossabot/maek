CREATE TABLE IF NOT EXISTS "note"
(
    id               BIGSERIAL PRIMARY KEY,
    uuid             VARCHAR(100) NOT NULL UNIQUE,
    content          TEXT         NOT NULL,
    favorite         BOOLEAN      NOT NULL DEFAULT FALSE,
    deleted          BOOLEAN      NOT NULL DEFAULT FALSE,
    trashed          BOOLEAN      NOT NULL DEFAULT FALSE,
    has_content      BOOLEAN      NOT NULL DEFAULT FALSE,
    has_images       BOOLEAN      NOT NULL DEFAULT FALSE,
    has_videos       BOOLEAN      NOT NULL DEFAULT FALSE,
    has_open_tasks   BOOLEAN      NOT NULL DEFAULT FALSE,
    has_closed_tasks BOOLEAN      NOT NULL DEFAULT FALSE,
    has_code         BOOLEAN      NOT NULL DEFAULT FALSE,
    has_audios       BOOLEAN      NOT NULL DEFAULT FALSE,
    has_links        BOOLEAN      NOT NULL DEFAULT FALSE,
    has_files        BOOLEAN      NOT NULL DEFAULT FALSE,
    has_quotes       BOOLEAN      NOT NULL DEFAULT FALSE,
    has_tables       BOOLEAN      NOT NULL DEFAULT FALSE,
    workspace_id     BIGINT       NOT NULL,
    created          BIGINT       NOT NULL DEFAULT 0,
    updated          BIGINT       NOT NULL DEFAULT 0,
    created_by_id    BIGINT       NOT NULL,
    updated_by_id    BIGINT       NOT NULL
);

-- Index for UUID lookups
CREATE INDEX idx_note_uuid ON "note" (uuid);

-- Composite index for workspace queries
CREATE INDEX idx_note_workspace_created ON "note" (workspace_id, created DESC);

-- Partial indexes for common filters
CREATE INDEX idx_note_favorite
    ON "note" (workspace_id, created DESC)
    WHERE favorite = TRUE AND deleted = FALSE AND trashed = FALSE;

CREATE INDEX idx_note_active
    ON "note" (workspace_id, created DESC)
    WHERE deleted = FALSE AND trashed = FALSE;
