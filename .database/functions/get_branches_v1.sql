CREATE OR REPLACE FUNCTION root.branches_search_matches_v1(p_row root.branches, p_q TEXT)
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
		OR p_row.address ILIKE '%' || v_q || '%';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION root.count_branches_v1(p_q TEXT)
RETURNS BIGINT AS $$
DECLARE
	v_q TEXT;
BEGIN
	v_q := NULLIF(btrim(coalesce(p_q, '')), '');
	RETURN (
		SELECT count(*)::BIGINT
		FROM root.branches b
		WHERE b.archived_at IS NULL
			AND root.branches_search_matches_v1(b, coalesce(v_q, ''))
	);
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION root.get_branches_v1(
	p_q TEXT,
	p_limit INTEGER,
	p_offset INTEGER
)
RETURNS SETOF root.branches AS $$
DECLARE
	v_q TEXT;
	v_limit INTEGER;
	v_offset INTEGER;
BEGIN
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
	SELECT b.*
	FROM root.branches b
	WHERE b.archived_at IS NULL
		AND root.branches_search_matches_v1(b, coalesce(v_q, ''))
	ORDER BY b.name, b.id
	LIMIT v_limit
	OFFSET v_offset;
END;
$$ LANGUAGE plpgsql STABLE;