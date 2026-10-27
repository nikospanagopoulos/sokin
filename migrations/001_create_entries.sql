CREATE TABLE entries (
    id         BIGSERIAL PRIMARY KEY,
    entry_date DATE        NOT NULL UNIQUE,
    did_today  TEXT        NOT NULL DEFAULT '',
    learned    TEXT        NOT NULL DEFAULT '',
    mistakes   TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);