DROP INDEX IF EXISTS root.athletes_middle_name_trgm_idx;
DROP INDEX IF EXISTS root.athletes_first_name_trgm_idx;
DROP INDEX IF EXISTS root.athletes_second_name_trgm_idx;

-- Keep pg_trgm installed: it may have existed before this migration or be used elsewhere.
