CREATE INDEX IF NOT EXISTS athletes_active_catalog_order_idx
	ON root.athletes (is_active, second_name, first_name, middle_name, id);
