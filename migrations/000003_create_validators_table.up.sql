CREATE TABLE IF NOT EXISTS validators
(
    id                          UUID DEFAULT   uuid_generate_v4(),
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL,
    name                        TEXT                     NOT NULL,
    address                     []NUMERIC(78)            NOT NULL,
    description                 TEXT                     NOT NULL,
    fee_rate                    DECIMAL(65, 0)           NOT NULL,
    active                      BOOLEAN                  NOT NULL,
    active_nodes                SMALLINT                 NOT NULL,
    staked                      DECIMAL(65, 0)           NOT NULL,
    pending                     DECIMAL(65, 0)           NOT NULL,
    rewards                     DECIMAL(65, 0)           NOT NULL,
    data                        JSONB                    NOT NULL,
    PRIMARY KEY (id)
);
-- Indexes
CREATE index idx_validators_address on validators (address);
