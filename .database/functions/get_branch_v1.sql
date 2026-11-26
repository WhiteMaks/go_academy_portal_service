CREATE OR REPLACE FUNCTION root.get_branch_v1(p_id BIGINT)
RETURNS root.branches AS $$
DECLARE
	found_branch root.branches;
BEGIN
	SELECT * INTO found_branch
	FROM root.branches
	WHERE id = p_id AND archived_at IS NULL;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'branch not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN found_branch;
END;
$$ LANGUAGE plpgsql;