CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

CREATE TABLE IF NOT EXISTS roles (
    role_id INT PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    tel_no VARCHAR(11) UNIQUE NOT NULL,
    mail VARCHAR UNIQUE NOT NULL,
    role_id INT REFERENCES roles(role_id) ON DELETE SET NULL,
    adress VARCHAR,
    wallet INT DEFAULT 0,
    total_recycle_count INT DEFAULT 0,
    coordinate VARCHAR
);

CREATE TABLE IF NOT EXISTS category (
    category_id INT PRIMARY KEY,
    waste_type VARCHAR NOT NULL,
    unit_type VARCHAR NOT NULL
);

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

CREATE TABLE IF NOT EXISTS recycling_points (
    recycling_points_id INT PRIMARY KEY,
    name VARCHAR NOT NULL,
    adress VARCHAR NOT NULL,
    coordinate VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS calendar (
    calendar_id INT PRIMARY KEY,
    is_available BOOLEAN DEFAULT TRUE,
    date DATE NOT NULL,
    hour TIME NOT NULL
);

CREATE TABLE IF NOT EXISTS appointment (
    appointment_id INT PRIMARY KEY,
    note VARCHAR,
    photo BYTEA,
    is_approved BOOLEAN DEFAULT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS appointment_calendar (
    id SERIAL PRIMARY KEY,
    appointment_id INT NOT NULL REFERENCES appointment(appointment_id) ON DELETE CASCADE,
    calendar_id INT NOT NULL REFERENCES calendar(calendar_id) ON DELETE CASCADE,
    UNIQUE (appointment_id, calendar_id)
);

CREATE TABLE IF NOT EXISTS appointment_waste (
    appointment_waste_id INT PRIMARY KEY,
    category_id INT REFERENCES category(category_id) ON DELETE CASCADE,
    amount INT NOT NULL,
    appointment_id INT REFERENCES appointment(appointment_id) ON DELETE CASCADE
);

-- Idempotent INSERT: only insert if role_id does not already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM roles WHERE role_id = 1) THEN
        INSERT INTO roles (role_id, name) VALUES
            (1, 'admin'),
            (2, 'customer'),
            (3, 'staff');
    END IF;
END
$$;