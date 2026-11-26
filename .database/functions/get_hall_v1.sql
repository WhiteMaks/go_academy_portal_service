CREATE OR REPLACE FUNCTION root.get_hall_v1(p_id BIGINT)
RETURNS root.halls AS $$
DECLARE
	found_hall root.halls;
BEGIN
	SELECT h.* INTO found_hall
	FROM root.halls h
	INNER JOIN root.branches b ON b.id = h.branch_id
	WHERE h.id = p_id
		AND h.archived_at IS NULL
		AND b.archived_at IS NULL;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'hall not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN found_hall;
END;
$$ LANGUAGE plpgsql;