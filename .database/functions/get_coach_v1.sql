CREATE OR REPLACE FUNCTION root.get_coach_v1(p_id BIGINT)
RETURNS root.coaches AS $$
DECLARE
	found_coach root.coaches;
BEGIN
	SELECT * INTO found_coach
	FROM root.coaches
	WHERE id = p_id AND archived_at IS NULL;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'coach not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN found_coach;
END;
$$ LANGUAGE plpgsql;