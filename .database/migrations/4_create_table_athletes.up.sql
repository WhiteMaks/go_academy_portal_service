CREATE TYPE root.athlete_activity_status AS ENUM ('active', 'inactive');

CREATE TABLE root.athletes (
	id                  BIGSERIAL       PRIMARY KEY,
	last_name           VARCHAR(100)    NOT NULL CHECK (btrim(last_name) <> ''),
	first_name          VARCHAR(100)    NOT NULL CHECK (btrim(first_name) <> ''),
	patronymic          VARCHAR(100)    NULL,
	birth_date          DATE            NOT NULL,
	phone               VARCHAR(16)     NOT NULL,
	parent_last_name    VARCHAR(100)    NULL,
	parent_first_name   VARCHAR(100)    NULL,
	parent_patronymic   VARCHAR(100)    NULL,
	parent_birth_date   DATE            NULL,
	parent_phone        VARCHAR(16)     NULL,
	note                TEXT            NULL CHECK (note IS NULL OR char_length(note) <= 2000),
	activity_status     root.athlete_activity_status NOT NULL DEFAULT 'active',
	archived_at         TIMESTAMPTZ     NULL,
	created_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
	updated_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
	CONSTRAINT athletes_phone_e164_check CHECK (phone ~ '^\+[1-9]\d{7,14}$'),
	CONSTRAINT athletes_parent_phone_e164_check CHECK (
		parent_phone IS NULL OR parent_phone ~ '^\+[1-9]\d{7,14}$'
	)
);

CREATE UNIQUE INDEX athletes_phone_active_uidx ON root.athletes (phone) WHERE archived_at IS NULL;