CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX athletes_second_name_trgm_idx
	ON root.athletes USING GIN (translate(lower(second_name), 'ё', 'е') gin_trgm_ops);

CREATE INDEX athletes_first_name_trgm_idx
	ON root.athletes USING GIN (translate(lower(first_name), 'ё', 'е') gin_trgm_ops);

CREATE INDEX athletes_middle_name_trgm_idx
	ON root.athletes USING GIN (translate(lower(coalesce(middle_name, '')), 'ё', 'е') gin_trgm_ops);


