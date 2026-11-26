DROP FUNCTION IF EXISTS root.assert_hall_parent_branch_active_v1(BIGINT);

CREATE OR REPLACE FUNCTION root.assert_hall_parent_branch_v1(p_branch_id BIGINT)
RETURNS root.branches AS $$
DECLARE
	parent_branch root.branches;
BEGIN
	SELECT * INTO parent_branch
	FROM root.branches
	WHERE id = p_branch_id;

	IF NOT FOUND OR parent_branch.archived_at IS NOT NULL THEN
		RAISE EXCEPTION 'branch not found' USING ERRCODE = 'P0002';
	END IF;

	RETURN parent_branch;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION root.assert_hall_activity_allowed_v1(p_branch_id BIGINT, p_activity_status TEXT)
RETURNS void AS $$
DECLARE
	parent_branch root.branches;
	v_activity_status TEXT;
BEGIN
	parent_branch := root.assert_hall_parent_branch_v1(p_branch_id);
	v_activity_status := lower(btrim(coalesce(p_activity_status, '')));
	IF parent_branch.activity_status::text = 'inactive' AND v_activity_status = 'active' THEN
		RAISE EXCEPTION 'cannot activate hall while branch is inactive' USING ERRCODE = '22023';
	END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION root.create_hall_v1(
	p_branch_id BIGINT,
	p_name TEXT,
	p_subtitle TEXT,
	p_capacity INTEGER,
	p_note TEXT
)
RETURNS root.halls AS $$
DECLARE
	v_name TEXT;
	v_subtitle TEXT;
	v_note TEXT;
	v_capacity INTEGER;
	parent_branch root.branches;
	v_hall_activity_status root.hall_activity_status;
	new_hall root.halls;
BEGIN

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

	parent_branch := root.assert_hall_parent_branch_v1(p_branch_id);
	IF parent_branch.activity_status::text = 'inactive' THEN
		v_hall_activity_status := 'inactive';
	ELSE
		v_hall_activity_status := 'active';
	END IF;

	BEGIN
		INSERT INTO root.halls (branch_id, name, subtitle, capacity, note, activity_status)
		VALUES (p_branch_id, v_name, v_subtitle, v_capacity, v_note, v_hall_activity_status)
		RETURNING * INTO new_hall;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'hall with the same name already exists in this branch' USING ERRCODE = '23505';
	END;

	RETURN new_hall;
END;
$$ LANGUAGE plpgsql;