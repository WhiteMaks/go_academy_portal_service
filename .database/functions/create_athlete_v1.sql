CREATE OR REPLACE FUNCTION root.normalize_phone_e164_v1(p_phone TEXT)
RETURNS TEXT AS $$
DECLARE
	v_normalized TEXT;
BEGIN
	v_normalized := regexp_replace(btrim(coalesce(p_phone, '')), '[\s\-]', '', 'g');
	IF v_normalized !~ '^\+[1-9]\d{7,14}$' THEN
		RAISE EXCEPTION 'invalid phone format' USING ERRCODE = '22023';
	END IF;
	RETURN v_normalized;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION root.create_athlete_v1(
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
	p_note TEXT
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
	v_parent_core_set BOOLEAN;
	v_parent_partial BOOLEAN;
	new_athlete root.athletes;
BEGIN
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

	BEGIN
		INSERT INTO root.athletes (
			last_name,
			first_name,
			patronymic,
			birth_date,
			phone,
			parent_last_name,
			parent_first_name,
			parent_patronymic,
			parent_birth_date,
			parent_phone,
			note,
			activity_status
		)
		VALUES (
			v_last_name,
			v_first_name,
			v_patronymic,
			p_birth_date,
			v_phone,
			CASE WHEN v_parent_core_set THEN v_parent_last_name END,
			CASE WHEN v_parent_core_set THEN v_parent_first_name END,
			CASE WHEN v_parent_core_set THEN v_parent_patronymic END,
			CASE WHEN v_parent_core_set THEN p_parent_birth_date END,
			CASE WHEN v_parent_core_set THEN v_parent_phone END,
			v_note,
			'active'
		)
		RETURNING * INTO new_athlete;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'athlete with the same phone already exists' USING ERRCODE = '23505';
	END;

	RETURN new_athlete;
END;
$$ LANGUAGE plpgsql;