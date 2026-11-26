CREATE OR REPLACE FUNCTION root.halls_search_matches_v1(p_row root.halls, p_q TEXT)
RETURNS BOOLEAN AS $$
DECLARE
	v_q TEXT;
BEGIN
	v_q := btrim(coalesce(p_q, ''));
	IF v_q = '' THEN
		RETURN TRUE;
	END IF;

	RETURN
		p_row.name ILIKE '%' || v_q || '%'
		OR (p_row.subtitle IS NOT NULL AND p_row.subtitle ILIKE '%' || v_q || '%');
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION root.count_branch_halls_v1(p_branch_id BIGINT, p_q TEXT)
RETURNS BIGINT AS $$
DECLARE
	v_q TEXT;
BEGIN
	v_q := NULLIF(btrim(coalesce(p_q, '')), '');
	RETURN (
		SELECT count(*)::BIGINT
		FROM root.halls h
		WHERE h.branch_id = p_branch_id
			AND h.archived_at IS NULL
			AND root.halls_search_matches_v1(h, coalesce(v_q, ''))
	);
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION root.get_branch_halls_v1(
	p_branch_id BIGINT,
	p_q TEXT,
	p_limit INTEGER,
	p_offset INTEGER
)
RETURNS SETOF root.halls AS $$
DECLARE
	v_q TEXT;
	v_limit INTEGER;
	v_offset INTEGER;
	parent_branch root.branches;
BEGIN
	SELECT * INTO parent_branch
	FROM root.branches
	WHERE id = p_branch_id AND archived_at IS NULL;

	IF NOT FOUND THEN
		RAISE EXCEPTION 'branch not found' USING ERRCODE = 'P0002';
	END IF;

	v_q := NULLIF(btrim(coalesce(p_q, '')), '');
	-- default/max must match Go model.ListQueryDefaultLimit / ListQueryMaxLimit (20)
	v_limit := coalesce(p_limit, 20);
	IF v_limit <= 0 THEN
		v_limit := 20;
	ELSIF v_limit > 20 THEN
		RAISE EXCEPTION 'invalid limit' USING ERRCODE = '22023';
	END IF;
	v_offset := greatest(coalesce(p_offset, 0), 0);

	RETURN QUERY
	SELECT h.*
	FROM root.halls h
	WHERE h.branch_id = p_branch_id
		AND h.archived_at IS NULL
		AND root.halls_search_matches_v1(h, coalesce(v_q, ''))
	ORDER BY h.name, h.id
	LIMIT v_limit
	OFFSET v_offset;
END;
$$ LANGUAGE plpgsql STABLE;