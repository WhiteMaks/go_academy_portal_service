CREATE OR REPLACE FUNCTION root.update_branch_v1(
	p_id BIGINT,
	p_name TEXT,
	p_address TEXT,
	p_note TEXT,
	p_activity_status TEXT
)
RETURNS root.branches AS $$
DECLARE
	v_name TEXT;
	v_address TEXT;
	v_note TEXT;
	v_activity_status TEXT;
	current_branch root.branches;
	updated_branch root.branches;
BEGIN
	SELECT * INTO current_branch
	FROM root.branches
	WHERE id = p_id AND archived_at IS NULL
	FOR UPDATE;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'branch not found' USING ERRCODE = 'P0002';
	END IF;

	v_name := btrim(coalesce(p_name, ''));
	v_address := btrim(coalesce(p_address, ''));
	IF v_name = '' OR v_address = '' THEN
		RAISE EXCEPTION 'name and address are required' USING ERRCODE = '22023';
	END IF;
	IF char_length(v_name) > 120 OR char_length(v_address) > 300 THEN
		RAISE EXCEPTION 'field value too long' USING ERRCODE = '22023';
	END IF;

	v_note := NULLIF(btrim(coalesce(p_note, '')), '');
	IF v_note IS NOT NULL AND char_length(v_note) > 2000 THEN
		RAISE EXCEPTION 'note value too long' USING ERRCODE = '22023';
	END IF;

	v_activity_status := lower(btrim(coalesce(p_activity_status, '')));
	IF v_activity_status NOT IN ('active', 'inactive') THEN
		RAISE EXCEPTION 'invalid activity_status' USING ERRCODE = '22023';
	END IF;

	BEGIN
		UPDATE root.branches
		SET
			name = v_name,
			address = v_address,
			note = v_note,
			activity_status = v_activity_status::root.branch_activity_status,
			updated_at = now()
		WHERE id = p_id
		RETURNING * INTO updated_branch;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'branch with the same name already exists' USING ERRCODE = '23505';
	END;

	IF v_activity_status = 'inactive'
		AND current_branch.activity_status::text <> 'inactive'
	THEN
		UPDATE root.halls
		SET
			activity_status = 'inactive'::root.hall_activity_status,
			updated_at = now()
		WHERE branch_id = p_id
			AND archived_at IS NULL;
	END IF;

	RETURN updated_branch;
END;
$$ LANGUAGE plpgsql;