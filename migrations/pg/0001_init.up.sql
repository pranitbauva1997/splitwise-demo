CREATE TABLE IF NOT EXISTS users
(
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    is_deleted BOOL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS bills
(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    is_deleted BOOL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS transactions
(
    id BIGSERIAL PRIMARY KEY,
    bill_id BIGINT NOT NULL,
    owed_to BIGINT NOT NULL,
    owes BIGINT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    is_deleted BOOL DEFAULT FALSE
);
