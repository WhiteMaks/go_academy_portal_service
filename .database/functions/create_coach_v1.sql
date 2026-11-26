CREATE OR REPLACE FUNCTION root.create_coach_v1(
	p_last_name TEXT,
	p_first_name TEXT,
	p_patronymic TEXT,
	p_phone TEXT,
	p_specialization TEXT,
	p_note TEXT,
	p_birth_date DATE
)
RETURNS root.coaches AS $$
DECLARE
	v_last_name TEXT;
	v_first_name TEXT;
	v_patronymic TEXT;
	v_phone TEXT;
	v_specialization TEXT;
	v_note TEXT;
	new_coach root.coaches;
BEGIN
	v_last_name := btrim(coalesce(p_last_name, ''));
	v_first_name := btrim(coalesce(p_first_name, ''));
	IF v_last_name = '' OR v_first_name = '' THEN
		RAISE EXCEPTION 'last_name and first_name are required' USING ERRCODE = '22023';
	END IF;
	IF char_length(v_last_name) > 100 OR char_length(v_first_name) > 100 THEN
		RAISE EXCEPTION 'name value too long' USING ERRCODE = '22023';
	END IF;

	v_patronymic := NULLIF(btrim(coalesce(p_patronymic, '')), '');
	IF v_patronymic IS NOT NULL AND char_length(v_patronymic) > 100 THEN
		RAISE EXCEPTION 'patronymic value too long' USING ERRCODE = '22023';
	END IF;

	v_phone := root.normalize_phone_e164_v1(p_phone);

	v_specialization := btrim(coalesce(p_specialization, ''));
	IF v_specialization = '' THEN
		RAISE EXCEPTION 'specialization is required' USING ERRCODE = '22023';
	END IF;
	IF char_length(v_specialization) > 200 THEN
		RAISE EXCEPTION 'specialization value too long' USING ERRCODE = '22023';
	END IF;

	v_note := NULLIF(btrim(coalesce(p_note, '')), '');
	IF v_note IS NOT NULL AND char_length(v_note) > 2000 THEN
		RAISE EXCEPTION 'note value too long' USING ERRCODE = '22023';
	END IF;

	IF p_birth_date IS NULL THEN
		RAISE EXCEPTION 'birth_date is required' USING ERRCODE = '22023';
	END IF;
	IF p_birth_date > CURRENT_DATE THEN
		RAISE EXCEPTION 'birth_date cannot be in the future' USING ERRCODE = '22023';
	END IF;

	BEGIN
		INSERT INTO root.coaches (
			last_name,
			first_name,
			patronymic,
			phone,
			specialization,
			note,
			activity_status,
			birth_date
		)
		VALUES (
			v_last_name,
			v_first_name,
			v_patronymic,
			v_phone,
			v_specialization,
			v_note,
			'active',
			p_birth_date
		)
		RETURNING * INTO new_coach;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'coach with the same phone already exists' USING ERRCODE = '23505';
	END;

	RETURN new_coach;
END;
$$ LANGUAGE plpgsql;