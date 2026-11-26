CREATE OR REPLACE FUNCTION root.coaches_search_matches_v1(p_row root.coaches, p_q TEXT)
RETURNS BOOLEAN AS $$
DECLARE
	v_q TEXT;
	v_digits_q TEXT;
	v_phone_digits TEXT;
BEGIN
	v_q := btrim(coalesce(p_q, ''));
	IF v_q = '' THEN
		RETURN TRUE;
	END IF;

	v_digits_q := regexp_replace(v_q, '[^\d]', '', 'g');
	v_phone_digits := regexp_replace(p_row.phone, '[^\d]', '', 'g');

	RETURN
		p_row.last_name ILIKE '%' || v_q || '%'
		OR p_row.first_name ILIKE '%' || v_q || '%'
		OR (p_row.patronymic IS NOT NULL AND p_row.patronymic ILIKE '%' || v_q || '%')
		OR p_row.specialization ILIKE '%' || v_q || '%'
		OR (v_digits_q <> '' AND v_phone_digits LIKE '%' || v_digits_q || '%')
		OR to_char(p_row.birth_date, 'DD.MM.YYYY') ILIKE '%' || v_q || '%'
		OR to_char(p_row.birth_date, 'YYYY-MM-DD') ILIKE '%' || v_q || '%'
		OR to_char(p_row.birth_date, 'YYYY') ILIKE '%' || v_q || '%';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION root.count_coaches_v1(p_q TEXT)
RETURNS BIGINT AS $$
DECLARE
	v_q TEXT;
BEGIN
	v_q := NULLIF(btrim(coalesce(p_q, '')), '');
	RETURN (
		SELECT count(*)::BIGINT
		FROM root.coaches c
		WHERE c.archived_at IS NULL
			AND root.coaches_search_matches_v1(c, coalesce(v_q, ''))
	);
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION root.get_coaches_v1(
	p_q TEXT,
	p_limit INTEGER,
	p_offset INTEGER
)
RETURNS SETOF root.coaches AS $$
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
	SELECT c.*
	FROM root.coaches c
	WHERE c.archived_at IS NULL
		AND root.coaches_search_matches_v1(c, coalesce(v_q, ''))
	ORDER BY c.last_name, c.first_name, c.id
	LIMIT v_limit
	OFFSET v_offset;
END;
$$ LANGUAGE plpgsql STABLE;