CREATE OR REPLACE FUNCTION root.create_user_v1(
	p_username VARCHAR(64),
	p_email VARCHAR(254),
	p_password VARCHAR(255),
	p_role root.user_role,
	p_is_active BOOLEAN
)
RETURNS
	root.users AS $$
DECLARE
	new_user root.users;
BEGIN
	BEGIN
		INSERT INTO
			root.users (username, email, password, role, is_active)
		VALUES
			(p_username, p_email, p_password, p_role, p_is_active)
		RETURNING * INTO new_user;
	EXCEPTION
		WHEN
			unique_violation
		THEN
			RAISE EXCEPTION 'user with the same username or email already exists' USING ERRCODE = '23505';
	END;

	RETURN new_user;
END;
$$ LANGUAGE plpgsql;