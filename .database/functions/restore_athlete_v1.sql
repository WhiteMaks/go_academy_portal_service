CREATE OR REPLACE FUNCTION root.restore_athlete_v1(p_id BIGINT)
RETURNS SETOF root.athletes AS $$
DECLARE
	current_athlete root.athletes;
BEGIN
	SELECT * INTO current_athlete
	FROM root.athletes
	WHERE id = p_id
	FOR UPDATE;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'athlete not found' USING ERRCODE = 'P0002';
	END IF;

	IF current_athlete.is_active = TRUE THEN
		RAISE EXCEPTION 'athlete is already active' USING ERRCODE = 'P4092';
	END IF;

	RETURN QUERY
	UPDATE root.athletes AS athlete
	SET is_active = TRUE,
	    updated_at = GREATEST(clock_timestamp(), athlete.updated_at + INTERVAL '1 microsecond')
	WHERE athlete.id = p_id
	RETURNING athlete.*;
END;
$$ LANGUAGE plpgsql VOLATILE;
