CREATE OR REPLACE FUNCTION root.create_branch_v1(
	p_name TEXT,
	p_address TEXT,
	p_note TEXT
)
RETURNS root.branches AS $$
DECLARE
	v_name TEXT;
	v_address TEXT;
	v_note TEXT;
	new_branch root.branches;
BEGIN
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

	BEGIN
		INSERT INTO root.branches (name, address, note, activity_status)
		VALUES (v_name, v_address, v_note, 'active')
		RETURNING * INTO new_branch;
	EXCEPTION
		WHEN unique_violation THEN
			RAISE EXCEPTION 'branch with the same name already exists' USING ERRCODE = '23505';
	END;

	RETURN new_branch;
END;
$$ LANGUAGE plpgsql;