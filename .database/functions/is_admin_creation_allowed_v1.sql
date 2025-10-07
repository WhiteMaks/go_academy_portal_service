CREATE OR REPLACE FUNCTION root.is_admin_creation_allowed_v1()
RETURNS BOOLEAN AS $$
BEGIN
	RETURN NOT EXISTS (
		SELECT 1 FROM root.users WHERE role = 'admin'::root.user_role
	);
END;
$$ LANGUAGE plpgsql;