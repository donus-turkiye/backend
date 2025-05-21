-- 1.1: Make the 'is_approved' column NULL (just remove the default value)
DO $$
BEGIN
    -- Check if column exists and has a default
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'appointment'
          AND column_name = 'is_approved'
          AND column_default IS NOT NULL
    ) THEN
        ALTER TABLE appointment ALTER COLUMN is_approved DROP DEFAULT;
    END IF;
END;
$$;

-- 1.2: Remove the 'date_time' column
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'appointment'
          AND column_name = 'date_time'
    ) THEN
        ALTER TABLE appointment DROP COLUMN date_time;
    END IF;
END;
$$;

CREATE TABLE IF NOT EXISTS appointment_calendar (
    id SERIAL PRIMARY KEY,
    appointment_id INT NOT NULL REFERENCES appointment(appointment_id) ON DELETE CASCADE,
    calendar_id INT NOT NULL REFERENCES calendar(calendar_id) ON DELETE CASCADE,
    UNIQUE (appointment_id, calendar_id)
);