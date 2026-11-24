CREATE OR REPLACE FUNCTION root.get_athlete_by_id_v1(p_id BIGINT)
RETURNS SETOF root.athletes AS $$
BEGIN
	RETURN QUERY
	SELECT *
	FROM root.athletes
	WHERE id = p_id;
END;
$$ LANGUAGE plpgsql STABLE;
