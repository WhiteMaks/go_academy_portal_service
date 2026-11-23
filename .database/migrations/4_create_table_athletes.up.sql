CREATE TABLE root.athletes (
	id BIGSERIAL PRIMARY KEY,
	second_name VARCHAR(100) NOT NULL CHECK (btrim(second_name) <> ''),
	first_name VARCHAR(100) NOT NULL CHECK (btrim(first_name) <> ''),
	middle_name VARCHAR(100),
	birth_date DATE NOT NULL,
	photo_url VARCHAR(2048),
	phone VARCHAR(32),
	parent_second_name VARCHAR(100),
	parent_first_name VARCHAR(100),
	parent_middle_name VARCHAR(100),
	parent_phone VARCHAR(32),
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
