CREATE TABLE IF NOT EXISTS validators
(
    id                          UUID DEFAULT   uuid_generate_v4(),
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL,
    validator_id                DECIMAL(65, 0)           NOT NULL,
    name                        TEXT                     NOT NULL,
    validator_address           NUMERIC(78)              NOT NULL,
    requested_address           NUMERIC(78)              NOT NULL,
    description                 TEXT                     NOT NULL,
    fee_rate                    DECIMAL(65, 0)           NOT NULL,
    active                      BOOLEAN                  NOT NULL,
    active_nodes                SMALLINT                 NOT NULL,
    linked_nodes                SMALLINT                 NOT NULL,
    staked                      DECIMAL(65, 0)           NOT NULL,
    pending                     DECIMAL(65, 0)           NOT NULL,
    rewards                     DECIMAL(65, 0)           NOT NULL,
    data                        JSONB                    NOT NULL,
    UNIQUE(validator_id),
    PRIMARY KEY (id)
);
-- Indexes
CREATE index idx_validators_validator_id on validators  (validator_id);
