CREATE OR REPLACE FUNCTION root.update_hall_v1(
	p_id BIGINT,
	p_name TEXT,
	p_subtitle TEXT,
	p_capacity INTEGER,
	p_note TEXT,
	p_activity_status TEXT
)
RETURNS root.halls AS $$
DECLARE
	v_name TEXT;
	v_subtitle TEXT;
	v_note TEXT;
	v_capacity INTEGER;
	v_activity_status TEXT;
	current_hall root.halls;
	updated_hall root.halls;
BEGIN
	SELECT * INTO current_hall
	FROM root.halls
	WHERE id = p_id AND archived_at IS NULL
	FOR UPDATE;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'hall not found' USING ERRCODE = 'P0002';
	END IF;

	PERFORM root.assert_hall_parent_branch_v1(current_hall.branch_id);

	v_name := btrim(coalesce(p_name, ''));
	IF v_name = '' THEN
		RAISE EXCEPTION 'name is required' USING ERRCODE = '22023';
	END IF;
	IF char_length(v_name) > 120 THEN
		RAISE EXCEPTION 'name value too long' USING ERRCODE = '22023';
	END IF;

	v_subtitle := NULLIF(btrim(coalesce(p_subtitle, '')), '');
	IF v_subtitle IS NOT NULL AND char_length(v_subtitle) > 200 THEN
		RAISE EXCEPTION 'subtitle value too long' USING ERRCODE = '22023';
	END IF;

	v_capacity := p_capacity;
	IF v_capacity IS NULL OR v_capacity < 1 OR v_capacity > 999 THEN
		RAISE EXCEPTION 'invalid capacity' USING ERRCODE = '22023';
	END IF;

	v_note := NULLIF(btrim(coalesce(p_note, '')), '');
	IF v_note IS NOT NULL AND char_length(v_note) > 2000 THEN
		RAISE EXCEPTION 'note value too long' USING ERRCODE = '22023';
	END IF;

	v_activity_status := lower(btrim(coalesce(p_activity_status, '')));
	IF v_activity_status NOT IN ('active', 'inactive') THEN
		RAISE EXCEPTION 'invalid activity_status' USING ERRCODE = '22023';
	END IF;

	PERFORM root.assert_hall_activity_allowed_v1(current_hall.branch_id, v_activity_status);

	BEGIN
		UPDATE root.halls
		SET
			name = v_name,
			subtitle = v_subtitle,
			capacity = v_capacity,
			note = v_note,
			activity_status = v_activity_status::root.hall_activity_status,
			updated_at = now()
		WHERE id = p_id
		RETURNING * INTO updated_hall;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'hall with the same name already exists in this branch' USING ERRCODE = '23505';
	END;

	RETURN updated_hall;
END;
$$ LANGUAGE plpgsql;