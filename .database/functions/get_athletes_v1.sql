DROP FUNCTION IF EXISTS root.get_athletes_v1(BOOLEAN, TEXT, BOOLEAN, INTEGER, INTEGER);
DROP FUNCTION IF EXISTS root.get_athletes_v1(BOOLEAN, TEXT, BOOLEAN, INTEGER, INTEGER, BOOLEAN);
DROP FUNCTION IF EXISTS root.get_athletes_search_v1(BOOLEAN, TEXT[], INTEGER, INTEGER);
DROP FUNCTION IF EXISTS root.get_athletes_search_v1(BOOLEAN, TEXT[], TEXT[], INTEGER, INTEGER);

CREATE OR REPLACE FUNCTION root.get_athletes_search_v1(
	p_is_active BOOLEAN,
	p_patterns TEXT[],
	p_score_tokens TEXT[],
	p_limit INTEGER,
	p_offset INTEGER,
	p_is_coach BOOLEAN
)
RETURNS SETOF root.athletes AS $$
BEGIN
	IF cardinality(p_patterns) = 1 AND cardinality(p_score_tokens) = 1 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
		  )
		ORDER BY
			(translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\') DESC,
			(translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\') DESC,
			(translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\') DESC,
			athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
		LIMIT p_limit OFFSET p_offset;
		RETURN;
	END IF;
	IF cardinality(p_patterns) = 1 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		CROSS JOIN LATERAL (
			SELECT
				count(*) FILTER (WHERE strpos(translate(lower(athlete.second_name), 'ё', 'е'), search_token.token) > 0) AS second_name_matches,
				count(*) FILTER (WHERE strpos(translate(lower(athlete.first_name), 'ё', 'е'), search_token.token) > 0) AS first_name_matches,
				count(*) FILTER (WHERE strpos(translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е'), search_token.token) > 0) AS middle_name_matches
			FROM unnest(p_score_tokens) AS search_token(token)
		) AS relevance
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND (
			translate(lower(athlete.second_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(athlete.first_name), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
			OR translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е') LIKE p_patterns[1] ESCAPE E'\\'
		  )
		ORDER BY relevance.second_name_matches DESC, relevance.first_name_matches DESC, relevance.middle_name_matches DESC,
			athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
		LIMIT p_limit OFFSET p_offset;
		RETURN;
	ELSIF cardinality(p_patterns) = 2 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		CROSS JOIN LATERAL (
			SELECT
				count(*) FILTER (WHERE strpos(translate(lower(athlete.second_name), 'ё', 'е'), search_token.token) > 0) AS second_name_matches,
				count(*) FILTER (WHERE strpos(translate(lower(athlete.first_name), 'ё', 'е'), search_token.token) > 0) AS first_name_matches,
				count(*) FILTER (WHERE strpos(translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е'), search_token.token) > 0) AS middle_name_matches
			FROM unnest(p_score_tokens) AS search_token(token)
		) AS relevance
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
		ORDER BY relevance.second_name_matches DESC, relevance.first_name_matches DESC, relevance.middle_name_matches DESC,
			athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
		LIMIT p_limit OFFSET p_offset;
		RETURN;
	ELSIF cardinality(p_patterns) = 3 THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		CROSS JOIN LATERAL (
			SELECT
				count(*) FILTER (WHERE strpos(translate(lower(athlete.second_name), 'ё', 'е'), search_token.token) > 0) AS second_name_matches,
				count(*) FILTER (WHERE strpos(translate(lower(athlete.first_name), 'ё', 'е'), search_token.token) > 0) AS first_name_matches,
				count(*) FILTER (WHERE strpos(translate(lower(coalesce(athlete.middle_name, '')), 'ё', 'е'), search_token.token) > 0) AS middle_name_matches
			FROM unnest(p_score_tokens) AS search_token(token)
		) AS relevance
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
		  )
		ORDER BY relevance.second_name_matches DESC, relevance.first_name_matches DESC, relevance.middle_name_matches DESC,
			athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
		LIMIT p_limit OFFSET p_offset;
	END IF;
	RETURN;
END;
$$ LANGUAGE plpgsql STABLE SET work_mem = '16MB';

CREATE FUNCTION root.get_athletes_v1(
	p_is_active BOOLEAN,
	p_search_query TEXT,
	p_is_birth_date_query BOOLEAN,
	p_limit INTEGER,
	p_offset INTEGER,
	p_is_coach BOOLEAN
)
RETURNS SETOF root.athletes AS $$
DECLARE
	v_patterns TEXT[];
	v_score_tokens TEXT[];
BEGIN
	IF p_search_query = '' THEN
		IF p_is_coach IS NULL THEN
			RETURN QUERY
			SELECT athlete.*
			FROM root.athletes AS athlete
			WHERE athlete.is_active = p_is_active
			ORDER BY athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
			LIMIT p_limit OFFSET p_offset;
		ELSE
			RETURN QUERY
			SELECT athlete.*
			FROM root.athletes AS athlete
			WHERE athlete.is_active = p_is_active
			  AND athlete.is_coach = p_is_coach
			ORDER BY athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
			LIMIT p_limit OFFSET p_offset;
		END IF;
	ELSIF p_is_birth_date_query THEN
		RETURN QUERY
		SELECT athlete.*
		FROM root.athletes AS athlete
		WHERE athlete.is_active = p_is_active
		  AND (p_is_coach IS NULL OR athlete.is_coach = p_is_coach)
		  AND to_char(athlete.birth_date, 'DDMMYYYY') LIKE p_search_query || '%'
		ORDER BY athlete.second_name ASC, athlete.first_name ASC, athlete.middle_name ASC NULLS LAST, athlete.id ASC
		LIMIT p_limit OFFSET p_offset;
	ELSE
		v_score_tokens := regexp_split_to_array(translate(lower(p_search_query), 'ё', 'е'), '[[:space:]]+');
		IF cardinality(v_score_tokens) > 3 THEN
			RETURN;
		END IF;
		v_patterns := root.athlete_search_patterns_v1(v_score_tokens);
		RETURN QUERY
		SELECT * FROM root.get_athletes_search_v1(p_is_active, v_patterns, v_score_tokens, p_limit, p_offset, p_is_coach);
	END IF;
	RETURN;
END;
$$ LANGUAGE plpgsql STABLE;
