CREATE OR REPLACE FUNCTION root.create_athlete_v1(
	p_second_name VARCHAR(100),
	p_first_name VARCHAR(100),
	p_middle_name VARCHAR(100),
	p_birth_date DATE,
	p_photo_url VARCHAR(2048),
	p_phone VARCHAR(32),
	p_parent_second_name VARCHAR(100),
	p_parent_first_name VARCHAR(100),
	p_parent_middle_name VARCHAR(100),
	p_parent_phone VARCHAR(32)
)
RETURNS root.athletes AS $$
DECLARE
	new_athlete root.athletes;
BEGIN
	INSERT INTO root.athletes (
		second_name,
		first_name,
		middle_name,
		birth_date,
		photo_url,
		phone,
		parent_second_name,
		parent_first_name,
		parent_middle_name,
		parent_phone
	)
	VALUES (
		p_second_name,
		p_first_name,
		p_middle_name,
		p_birth_date,
		p_photo_url,
		p_phone,
		p_parent_second_name,
		p_parent_first_name,
		p_parent_middle_name,
		p_parent_phone
	)
	RETURNING * INTO new_athlete;

	RETURN new_athlete;
END;
$$ LANGUAGE plpgsql;
