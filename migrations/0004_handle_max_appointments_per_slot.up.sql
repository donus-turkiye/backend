
CREATE TABLE IF NOT EXISTS system_settings (
    setting_key TEXT PRIMARY KEY,
    setting_value INTEGER NOT NULL
);

-- Update the 'max_appointments_per_slot' setting if it exists, otherwise add it
INSERT INTO system_settings (setting_key, setting_value)
VALUES ('max_appointments_per_slot', 3)
ON CONFLICT (setting_key) DO UPDATE
SET setting_value = EXCLUDED.setting_value;

-- 3. limit_appointments_per_slot trigger function
CREATE OR REPLACE FUNCTION limit_appointments_per_slot()
RETURNS TRIGGER AS $$
DECLARE
    count_appointments INTEGER;
    max_allowed INTEGER;
BEGIN
    SELECT setting_value INTO max_allowed
    FROM system_settings
    WHERE setting_key = 'max_appointments_per_slot';

    SELECT COUNT(*) INTO count_appointments
    FROM appointment_calendar
    WHERE calendar_id = NEW.calendar_id;

    IF count_appointments = max_allowed THEN
        UPDATE calendar SET is_available = false WHERE calendar_id = NEW.calendar_id;
    ELSIF count_appointments > max_allowed THEN
        RAISE EXCEPTION 'A maximum of % appointment can be assigned to this calendar slot.', max_allowed;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- handle_disapproval trigger function
CREATE OR REPLACE FUNCTION handle_disapproval()
RETURNS TRIGGER AS $$
DECLARE
    related_calendar_id INTEGER;
    count_appointments INTEGER;
    max_allowed INTEGER;
BEGIN
    IF NEW.is_approved = false AND OLD.is_approved IS DISTINCT FROM false THEN
        SELECT calendar_id INTO related_calendar_id
        FROM appointment_calendar
        WHERE appointment_id = NEW.appointment_id
        LIMIT 1;

        DELETE FROM appointment_calendar
        WHERE appointment_id = NEW.appointment_id;

        SELECT setting_value INTO max_allowed
        FROM system_settings
        WHERE setting_key = 'max_appointments_per_slot';

        SELECT COUNT(*) INTO count_appointments
        FROM appointment_calendar
        WHERE calendar_id = related_calendar_id;

        IF count_appointments < max_allowed THEN
            UPDATE calendar SET is_available = true WHERE calendar_id = related_calendar_id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Assigning triggers (if any, they may not have been deleted before)
DROP TRIGGER IF EXISTS trg_limit_appointments ON appointment_calendar;
CREATE TRIGGER trg_limit_appointments
AFTER INSERT ON appointment_calendar
FOR EACH ROW
EXECUTE FUNCTION limit_appointments_per_slot();

DROP TRIGGER IF EXISTS trg_handle_disapproval ON appointment;
CREATE TRIGGER trg_handle_disapproval
AFTER UPDATE OF is_approved ON appointment
FOR EACH ROW
EXECUTE FUNCTION handle_disapproval();