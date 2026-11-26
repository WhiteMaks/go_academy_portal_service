CREATE OR REPLACE FUNCTION root.get_athlete_v1(p_id BIGINT)
RETURNS root.athletes AS $$
DECLARE
	found_athlete root.athletes;
BEGIN
	SELECT * INTO found_athlete
	FROM root.athletes
	WHERE id = p_id;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'athlete not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN found_athlete;
END;
$$ LANGUAGE plpgsql;