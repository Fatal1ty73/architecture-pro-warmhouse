CREATE DATABASE smarthome;
\c smarthome;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS houses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sensor_types (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS device_types (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sensors (
    id BIGSERIAL PRIMARY KEY,
    house_id BIGINT NOT NULL DEFAULT 1 REFERENCES houses(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type_id BIGINT NOT NULL DEFAULT 1 REFERENCES sensor_types(id),
    -- Legacy column for monolith compatibility.
    type VARCHAR(50) NOT NULL DEFAULT 'temperature',
    location VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL DEFAULT 'legacy',
    value FLOAT DEFAULT 0,
    unit VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,
    house_id BIGINT NOT NULL DEFAULT 1 REFERENCES houses(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type_id BIGINT NOT NULL DEFAULT 1 REFERENCES device_types(id),
    location VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO users (id, name, email, login, password)
VALUES (1, 'Admin', 'admin@warmhouse.local', 'admin', 'admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO houses (id, user_id, address)
VALUES (1, 1, 'ул.Центральная д.1')
ON CONFLICT (id) DO NOTHING;

INSERT INTO sensor_types (id, name)
VALUES (1, 'temperature')
ON CONFLICT (id) DO NOTHING;

INSERT INTO device_types (id, name)
VALUES (1, 'heating')
ON CONFLICT (id) DO NOTHING;

INSERT INTO sensors (
    id, house_id, name, type_id, type, location, serial_number, value, unit, status
)
VALUES (
    1, 1, 'Датчик температуры улицы', 1, 'temperature', 'Kitchen', 'temp-001', 22.5, 'C', 'active'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO devices (
    id, house_id, name, type_id, location, serial_number, status
)
VALUES (
    1, 1, 'Отопление', 1, 'Boiler Room', 'heating-001', 'inactive'
)
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM users), 1));
SELECT setval('houses_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM houses), 1));
SELECT setval('sensor_types_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM sensor_types), 1));
SELECT setval('device_types_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM device_types), 1));
SELECT setval('sensors_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM sensors), 1));
SELECT setval('devices_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM devices), 1));
