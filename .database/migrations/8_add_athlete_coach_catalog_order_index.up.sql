CREATE INDEX IF NOT EXISTS athletes_coach_catalog_order_idx
	ON root.athletes (is_active, is_coach, second_name, first_name, middle_name, id);
