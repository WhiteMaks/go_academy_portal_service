CREATE OR REPLACE FUNCTION root.update_athlete_v1(
	p_id BIGINT,
	p_last_name TEXT,
	p_first_name TEXT,
	p_patronymic TEXT,
	p_birth_date DATE,
	p_phone TEXT,
	p_parent_last_name TEXT,
	p_parent_first_name TEXT,
	p_parent_patronymic TEXT,
	p_parent_birth_date DATE,
	p_parent_phone TEXT,
	p_note TEXT,
	p_activity_status TEXT
)
RETURNS root.athletes AS $$
DECLARE
	v_last_name TEXT;
	v_first_name TEXT;
	v_patronymic TEXT;
	v_phone TEXT;
	v_parent_last_name TEXT;
	v_parent_first_name TEXT;
	v_parent_patronymic TEXT;
	v_parent_phone TEXT;
	v_note TEXT;
	v_activity_status TEXT;
	v_parent_core_set BOOLEAN;
	v_parent_partial BOOLEAN;
	current_athlete root.athletes;
	updated_athlete root.athletes;
BEGIN
	SELECT * INTO current_athlete
	FROM root.athletes
	WHERE id = p_id AND archived_at IS NULL
	FOR UPDATE;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'athlete not found' USING ERRCODE = 'P0002';
	END IF;

	v_last_name := btrim(coalesce(p_last_name, ''));
	v_first_name := btrim(coalesce(p_first_name, ''));
	IF v_last_name = '' OR v_first_name = '' THEN
		RAISE EXCEPTION 'last_name and first_name are required' USING ERRCODE = '22023';
	END IF;
	IF char_length(v_last_name) > 100 OR char_length(v_first_name) > 100 THEN
		RAISE EXCEPTION 'name value too long' USING ERRCODE = '22023';
	END IF;
	IF p_birth_date IS NULL THEN
		RAISE EXCEPTION 'birth_date is required' USING ERRCODE = '22023';
	END IF;
	IF p_birth_date > CURRENT_DATE THEN
		RAISE EXCEPTION 'birth_date cannot be in the future' USING ERRCODE = '22023';
	END IF;

	v_patronymic := NULLIF(btrim(coalesce(p_patronymic, '')), '');
	IF v_patronymic IS NOT NULL AND char_length(v_patronymic) > 100 THEN
		RAISE EXCEPTION 'patronymic value too long' USING ERRCODE = '22023';
	END IF;
	v_phone := root.normalize_phone_e164_v1(p_phone);

	v_parent_last_name := NULLIF(btrim(coalesce(p_parent_last_name, '')), '');
	v_parent_first_name := NULLIF(btrim(coalesce(p_parent_first_name, '')), '');
	v_parent_patronymic := NULLIF(btrim(coalesce(p_parent_patronymic, '')), '');
	v_parent_phone := CASE
		WHEN p_parent_phone IS NULL OR btrim(p_parent_phone) = '' THEN NULL
		ELSE root.normalize_phone_e164_v1(p_parent_phone)
	END;
	v_note := NULLIF(btrim(coalesce(p_note, '')), '');
	IF v_note IS NOT NULL AND char_length(v_note) > 2000 THEN
		RAISE EXCEPTION 'note value too long' USING ERRCODE = '22023';
	END IF;

	IF v_parent_last_name IS NOT NULL AND char_length(v_parent_last_name) > 100 THEN
		RAISE EXCEPTION 'parent_last_name value too long' USING ERRCODE = '22023';
	END IF;
	IF v_parent_first_name IS NOT NULL AND char_length(v_parent_first_name) > 100 THEN
		RAISE EXCEPTION 'parent_first_name value too long' USING ERRCODE = '22023';
	END IF;
	IF v_parent_patronymic IS NOT NULL AND char_length(v_parent_patronymic) > 100 THEN
		RAISE EXCEPTION 'parent_patronymic value too long' USING ERRCODE = '22023';
	END IF;

	v_parent_core_set := (
		v_parent_last_name IS NOT NULL
		AND v_parent_first_name IS NOT NULL
		AND p_parent_birth_date IS NOT NULL
		AND v_parent_phone IS NOT NULL
	);
	v_parent_partial := (
		v_parent_last_name IS NOT NULL
		OR v_parent_first_name IS NOT NULL
		OR v_parent_patronymic IS NOT NULL
		OR p_parent_birth_date IS NOT NULL
		OR v_parent_phone IS NOT NULL
	);
	IF v_parent_partial AND NOT v_parent_core_set THEN
		RAISE EXCEPTION 'fill all parent fields or leave the parent block empty' USING ERRCODE = '22023';
	END IF;
	IF v_parent_core_set AND p_parent_birth_date > CURRENT_DATE THEN
		RAISE EXCEPTION 'parent_birth_date cannot be in the future' USING ERRCODE = '22023';
	END IF;

	v_activity_status := lower(btrim(coalesce(p_activity_status, '')));
	IF v_activity_status NOT IN ('active', 'inactive') THEN
		RAISE EXCEPTION 'invalid activity_status' USING ERRCODE = '22023';
	END IF;

	BEGIN
		UPDATE root.athletes
		SET
			last_name = v_last_name,
			first_name = v_first_name,
			patronymic = v_patronymic,
			birth_date = p_birth_date,
			phone = v_phone,
			parent_last_name = CASE WHEN v_parent_core_set THEN v_parent_last_name END,
			parent_first_name = CASE WHEN v_parent_core_set THEN v_parent_first_name END,
			parent_patronymic = CASE WHEN v_parent_core_set THEN v_parent_patronymic END,
			parent_birth_date = CASE WHEN v_parent_core_set THEN p_parent_birth_date END,
			parent_phone = CASE WHEN v_parent_core_set THEN v_parent_phone END,
			note = v_note,
			activity_status = v_activity_status::root.athlete_activity_status,
			updated_at = now()
		WHERE id = p_id
		RETURNING * INTO updated_athlete;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'athlete with the same phone already exists' USING ERRCODE = '23505';
	END;

	RETURN updated_athlete;
END;
$$ LANGUAGE plpgsql;