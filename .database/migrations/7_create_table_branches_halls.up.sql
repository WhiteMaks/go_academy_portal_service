CREATE TYPE root.branch_activity_status AS ENUM ('active', 'inactive');

CREATE TABLE root.branches (
	id                  BIGSERIAL       PRIMARY KEY,
	name                VARCHAR(120)    NOT NULL CHECK (btrim(name) <> ''),
	address             VARCHAR(300)    NOT NULL CHECK (btrim(address) <> ''),
	note                TEXT            NULL CHECK (note IS NULL OR char_length(note) <= 2000),
	activity_status     root.branch_activity_status NOT NULL DEFAULT 'active',
	archived_at         TIMESTAMPTZ     NULL,
	created_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
	updated_at          TIMESTAMPTZ     NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX branches_name_active_uidx
	ON root.branches (lower(btrim(name)))
	WHERE archived_at IS NULL;

CREATE TYPE root.hall_activity_status AS ENUM ('active', 'inactive');

CREATE TABLE root.halls (
	id                  BIGSERIAL       PRIMARY KEY,
	branch_id           BIGINT          NOT NULL REFERENCES root.branches (id),
	name                VARCHAR(120)    NOT NULL CHECK (btrim(name) <> ''),
	subtitle            VARCHAR(200)    NULL,
	capacity            INTEGER         NOT NULL CHECK (capacity BETWEEN 1 AND 999),
	note                TEXT            NULL CHECK (note IS NULL OR char_length(note) <= 2000),
	activity_status     root.hall_activity_status NOT NULL DEFAULT 'active',
	archived_at         TIMESTAMPTZ     NULL,
	created_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
	updated_at          TIMESTAMPTZ     NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX halls_branch_name_active_uidx
	ON root.halls (branch_id, lower(btrim(name)))
	WHERE archived_at IS NULL;