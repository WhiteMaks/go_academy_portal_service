CREATE OR REPLACE FUNCTION root.get_athletes_v1()
RETURNS SETOF root.athletes AS $$
BEGIN
	RETURN QUERY
	SELECT *
	FROM root.athletes
	WHERE is_active = TRUE
	ORDER BY
		second_name ASC,
		first_name ASC,
		middle_name ASC NULLS LAST,
		id ASC;
END;
$$ LANGUAGE plpgsql STABLE;
