CREATE OR REPLACE FUNCTION root.athletes_search_matches_v1(p_row root.athletes, p_q TEXT)
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
		OR (v_digits_q <> '' AND v_phone_digits LIKE '%' || v_digits_q || '%')
		OR to_char(p_row.birth_date, 'DD.MM.YYYY') ILIKE '%' || v_q || '%'
		OR to_char(p_row.birth_date, 'YYYY-MM-DD') ILIKE '%' || v_q || '%'
		OR to_char(p_row.birth_date, 'YYYY') ILIKE '%' || v_q || '%';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION root.count_athletes_v1(p_q TEXT)
RETURNS BIGINT AS $$
DECLARE
	v_q TEXT;
BEGIN
	v_q := NULLIF(btrim(coalesce(p_q, '')), '');
	RETURN (
		SELECT count(*)::BIGINT
		FROM root.athletes a
		WHERE a.archived_at IS NULL
			AND root.athletes_search_matches_v1(a, coalesce(v_q, ''))
	);
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION root.get_athletes_v1(
	p_q TEXT,
	p_limit INTEGER,
	p_offset INTEGER
)
RETURNS SETOF root.athletes AS $$
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
	SELECT a.*
	FROM root.athletes a
	WHERE a.archived_at IS NULL
		AND root.athletes_search_matches_v1(a, coalesce(v_q, ''))
	ORDER BY a.last_name, a.first_name, a.id
	LIMIT v_limit
	OFFSET v_offset;
END;
$$ LANGUAGE plpgsql STABLE;