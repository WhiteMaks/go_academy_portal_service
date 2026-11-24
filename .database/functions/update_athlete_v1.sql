DROP FUNCTION IF EXISTS root.update_athlete_v1(
	BIGINT, TIMESTAMPTZ, VARCHAR(100), VARCHAR(100), VARCHAR(100), DATE,
	VARCHAR(32), VARCHAR(100), VARCHAR(100), VARCHAR(100), VARCHAR(32)
);

CREATE OR REPLACE FUNCTION root.update_athlete_v1(
	p_id BIGINT,
	p_second_name VARCHAR(100),
	p_first_name VARCHAR(100),
	p_middle_name VARCHAR(100),
	p_birth_date DATE,
	p_phone VARCHAR(32),
	p_parent_second_name VARCHAR(100),
	p_parent_first_name VARCHAR(100),
	p_parent_middle_name VARCHAR(100),
	p_parent_phone VARCHAR(32)
)
RETURNS SETOF root.athletes AS $$
BEGIN
	RETURN QUERY
	UPDATE root.athletes AS athlete
	SET second_name = p_second_name,
	    first_name = p_first_name,
	    middle_name = p_middle_name,
	    birth_date = p_birth_date,
	    phone = p_phone,
	    parent_second_name = p_parent_second_name,
	    parent_first_name = p_parent_first_name,
	    parent_middle_name = p_parent_middle_name,
	    parent_phone = p_parent_phone,
	    updated_at = GREATEST(clock_timestamp(), athlete.updated_at + INTERVAL '1 microsecond')
	WHERE athlete.id = p_id
	RETURNING athlete.*;
END;
$$ LANGUAGE plpgsql VOLATILE;
