CREATE TYPE root.coach_activity_status AS ENUM ('active', 'inactive');

CREATE TABLE root.coaches (
	id                  BIGSERIAL       PRIMARY KEY,
	last_name           VARCHAR(100)    NOT NULL CHECK (btrim(last_name) <> ''),
	first_name          VARCHAR(100)    NOT NULL CHECK (btrim(first_name) <> ''),
	patronymic          VARCHAR(100)    NULL,
	phone               VARCHAR(16)     NOT NULL,
	specialization      VARCHAR(200)    NOT NULL CHECK (btrim(specialization) <> ''),
	note                TEXT            NULL CHECK (note IS NULL OR char_length(note) <= 2000),
	activity_status     root.coach_activity_status NOT NULL DEFAULT 'active',
	archived_at         TIMESTAMPTZ     NULL,
	created_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
	updated_at          TIMESTAMPTZ     NOT NULL DEFAULT now(),
	CONSTRAINT coaches_phone_e164_check CHECK (phone ~ '^\+[1-9]\d{7,14}$')
);

CREATE UNIQUE INDEX coaches_phone_active_uidx ON root.coaches (phone) WHERE archived_at IS NULL;