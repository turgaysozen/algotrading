DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'algotrading') THEN
        PERFORM dblink_connect('dbname=postgres');
        PERFORM dblink_exec('CREATE DATABASE algotrading');
        PERFORM dblink_disconnect();
    END IF;
END $$;

\c algotrading;

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE IF NOT EXISTS order_books (
    id SERIAL,
    event_type TEXT,
    symbol TEXT,
    event_time BIGINT NOT NULL,
    best_bid FLOAT NOT NULL,
    best_ask FLOAT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, event_time)
);

SELECT create_hypertable('order_books', 'event_time');

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL,
    price NUMERIC,
    quantity NUMERIC,
    status TEXT,
    order_type TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
);

SELECT create_hypertable('orders', 'created_at');

CREATE TABLE IF NOT EXISTS signals (
    id SERIAL,
    type TEXT,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    price NUMERIC,
    short_sma NUMERIC,
    long_sma NUMERIC,
    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

SELECT create_hypertable('signals', 'timestamp');


-- print all created tables to make sure they are created
SELECT * FROM timescaledb_information.hypertables