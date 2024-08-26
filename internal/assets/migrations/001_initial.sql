-- +migrate Up

CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    tx_hash  VARCHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    amount NUMERIC(78, 0) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- +migrate Down
