CREATE TABLE root.users (
	id 				bigserial 		PRIMARY KEY,
	username 		VARCHAR(64) 	UNIQUE NOT NULL,
	password 		VARCHAR(255) 	NOT NULL,
	role 			root.user_role 	NOT NULL,
	is_active 		BOOLEAN 		NOT NULL DEFAULT FALSE,
	created_at 		TIMESTAMPTZ 	NOT NULL DEFAULT now(),
	updated_at 		TIMESTAMPTZ 	NOT NULL DEFAULT now()
);