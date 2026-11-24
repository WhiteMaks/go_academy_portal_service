DROP FUNCTION IF EXISTS root.count_athletes_search_v1(BOOLEAN, TEXT[]);
DROP FUNCTION IF EXISTS root.count_athletes_v1(BOOLEAN, TEXT, BOOLEAN);

CREATE OR REPLACE FUNCTION root.count_athletes_search_v1(
	p_is_active BOOLEAN,
	p_patterns TEXT[],
	p_is_coach BOOLEAN
)
RETURNS BIGINT AS $$
DECLARE
	v_count BIGINT;
BEGIN
	IF cardinality(p_patterns) = 1 THEN
		SELECT count(*) INTO v_count
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
		  );
		RETURN v_count;
	ELSIF cardinality(p_patterns) = 2 THEN
		SELECT count(*) INTO v_count
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
		  )
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[2] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[2] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[2] ESCAPE E'\\'
		  );
		RETURN v_count;
	ELSIF cardinality(p_patterns) = 3 THEN
		SELECT count(*) INTO v_count
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
		  )
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[2] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[2] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[2] ESCAPE E'\\'
		  )
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[3] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[3] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[3] ESCAPE E'\\'
		  );
		RETURN v_count;
	END IF;
	RETURN 0;
END;
$$ LANGUAGE plpgsql STABLE SET work_mem = '16MB';

CREATE OR REPLACE FUNCTION root.count_athletes_v1(
	p_is_active BOOLEAN,
	p_search_query TEXT,
	p_is_birth_date_query BOOLEAN,
	p_is_coach BOOLEAN
)
RETURNS BIGINT AS $$
DECLARE
	v_patterns TEXT[];
	v_tokens TEXT[];
	v_count BIGINT;
BEGIN
	IF p_search_query = '' THEN
		IF p_is_coach IS NULL THEN
			SELECT count(*) INTO v_count
			FROM root.athletes AS athlete
			WHERE athlete.is_active = p_is_active;
		ELSE
			SELECT count(*) INTO v_count
			FROM root.athletes AS athlete
			WHERE athlete.is_active = p_is_active
			  AND athlete.is_coach = p_is_coach;
		END IF;
		RETURN v_count;
	ELSIF p_is_birth_date_query THEN
		SELECT count(*) INTO v_count
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND to_char(athlete.birth_date, 'DDMMYYYY') LIKE p_search_query || '%';
		RETURN v_count;
	END IF;

	v_tokens := regexp_split_to_array(translate(lower(p_search_query), 'ё', 'е'), '[[:space:]]+');
	IF cardinality(v_tokens) > 3 THEN
		RETURN 0;
	END IF;
	v_patterns := root.athlete_search_patterns_v1(v_tokens);
	RETURN root.count_athletes_search_v1(p_is_active, v_patterns, p_is_coach);
END;
$$ LANGUAGE plpgsql STABLE;
