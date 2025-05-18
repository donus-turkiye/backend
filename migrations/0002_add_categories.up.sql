DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM category WHERE category_id = 1) THEN
        INSERT INTO category (category_id, waste_type, unit_type) VALUES
            (1, 'Plastic', 'kg'),
            (2, 'Paper', 'kg'),
            (3, 'Glass', 'piece');
    END IF;
END
$$;