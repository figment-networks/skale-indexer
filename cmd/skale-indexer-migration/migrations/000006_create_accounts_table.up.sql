CREATE TABLE IF NOT EXISTS accounts
(
    id                              UUID DEFAULT   uuid_generate_v4(),
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    address                         NUMERIC(78)              NOT NULL,
    account_type                    VARCHAR(50)              NOT NULL,
    UNIQUE(address),
    PRIMARY KEY (id)
);
