CREATE OR REPLACE FUNCTION root.set_athlete_coach_v1(
	p_id BIGINT,
	p_is_coach BOOLEAN
)
RETURNS root.athletes AS $$
DECLARE
	updated_athlete root.athletes;
BEGIN
	UPDATE root.athletes AS athlete
	SET is_coach = p_is_coach,
		updated_at = CASE
			WHEN athlete.is_coach IS DISTINCT FROM p_is_coach THEN GREATEST(clock_timestamp(), athlete.updated_at + INTERVAL '1 microsecond')
			ELSE athlete.updated_at
		END
	WHERE athlete.id = p_id
	RETURNING * INTO updated_athlete;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'athlete not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN updated_athlete;
END;
$$ LANGUAGE plpgsql;
