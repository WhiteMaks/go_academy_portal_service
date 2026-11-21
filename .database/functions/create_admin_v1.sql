-- VOLATILE ensures the post-lock existence check can see a concurrent winner's commit.
CREATE OR REPLACE FUNCTION root.create_admin_v1(
	p_username VARCHAR(64),
	p_password VARCHAR(255)
)
RETURNS root.users AS $$
DECLARE
	new_admin root.users;
BEGIN
	-- The two-key advisory lock namespace (1094926660, 1) is reserved for the
	-- one-time bootstrap administrator creation operation.
	PERFORM pg_advisory_xact_lock(1094926660, 1);

	IF EXISTS (
		SELECT 1
		FROM root.users
		WHERE role = 'admin'::root.user_role
	) THEN
		RAISE EXCEPTION 'admin creation is not allowed' USING ERRCODE = 'P0002';
	END IF;

	new_admin := root.create_user_v1(
		p_username,
		p_password,
		'admin'::root.user_role,
		TRUE
	);

	RETURN new_admin;
END;
$$ LANGUAGE plpgsql VOLATILE;
