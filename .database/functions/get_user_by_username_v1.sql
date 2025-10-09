CREATE OR REPLACE FUNCTION root.get_user_by_username_v1(
	p_username VARCHAR(64)
)
RETURNS
	root.users AS $$
DECLARE
	record_user root.users;
BEGIN
	SELECT
		*
	INTO
		record_user
	FROM
		root.users
	WHERE
		username = p_username;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'user with given username not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN record_user;

END;
$$ LANGUAGE plpgsql;