CREATE OR REPLACE FUNCTION root.athlete_filter_v1(
	p_is_active BOOLEAN,
	p_search_query TEXT,
	p_is_birth_date_query BOOLEAN
)
RETURNS SETOF root.athletes AS $$
DECLARE
	v_patterns TEXT[];
BEGIN
	IF p_search_query = '' THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active;
		RETURN;
	END IF;

	IF p_is_birth_date_query THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND to_char(athlete.birth_date, 'DDMMYYYY') LIKE p_search_query || '%';
		RETURN;
	END IF;

	SELECT ARRAY(
		SELECT '%' || replace(replace(replace(search_token.token, E'\\', E'\\\\'), '%', E'\\%'), '_', E'\\_') || '%'
		FROM unnest(regexp_split_to_array(p_search_query, '[[:space:]]+')) AS search_token(token)
	) INTO v_patterns;

	IF cardinality(v_patterns) = 1 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
		  );
	ELSIF cardinality(v_patterns) = 2 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
		  )
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE v_patterns[2] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE v_patterns[2] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE v_patterns[2] ESCAPE E'\\'
		  );
	ELSIF cardinality(v_patterns) = 3 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE v_patterns[1] ESCAPE E'\\'
		  )
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE v_patterns[2] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE v_patterns[2] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE v_patterns[2] ESCAPE E'\\'
		  )
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE v_patterns[3] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE v_patterns[3] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE v_patterns[3] ESCAPE E'\\'
		  );
	END IF;
	RETURN;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION root.athlete_search_patterns_v1(p_tokens TEXT[])
RETURNS TEXT[] AS $$
DECLARE
	v_patterns TEXT[] := ARRAY[]::TEXT[];
	v_token TEXT;
	v_other_token TEXT;
	v_subsumed BOOLEAN;
	i INTEGER;
	j INTEGER;
BEGIN
	FOR i IN 1..cardinality(p_tokens) LOOP
		v_token := translate(lower(p_tokens[i]), 'ё', 'е');
		v_subsumed := FALSE;
		FOR j IN 1..cardinality(p_tokens) LOOP
			v_other_token := translate(lower(p_tokens[j]), 'ё', 'е');
			IF i <> j AND (
				(length(v_other_token) > length(v_token) AND strpos(v_other_token, v_token) > 0)
				OR (j < i AND v_other_token = v_token)
			) THEN
				v_subsumed := TRUE;
				EXIT;
			END IF;
		END LOOP;
		IF NOT v_subsumed THEN
			v_patterns := array_append(v_patterns, '%' || replace(replace(replace(v_token, E'\\', E'\\\\'), '%', E'\\%'), '_', E'\\_') || '%');
		END IF;
	END LOOP;
	RETURN v_patterns;
END;
$$ LANGUAGE plpgsql IMMUTABLE;
