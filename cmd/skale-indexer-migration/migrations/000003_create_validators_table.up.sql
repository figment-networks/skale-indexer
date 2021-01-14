CREATE TABLE IF NOT EXISTS validators
(
    id                          UUID DEFAULT   uuid_generate_v4(),
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    validator_id                DECIMAL(65, 0)           NOT NULL,
    validator_address           NUMERIC(78)              NOT NULL,
    requested_address           NUMERIC(78)              NOT NULL,
    fee_rate                    DECIMAL(65, 0)           NOT NULL DEFAULT 0,
    registration_time           TIMESTAMP WITH TIME ZONE NOT NULL,
    minimum_delegation_amount   DECIMAL(65, 0)           NOT NULL DEFAULT 0,
    accept_new_requests         BOOLEAN                  NOT NULL,
    authorized                  BOOLEAN                  NOT NULL,
    active_nodes                SMALLINT                 NOT NULL DEFAULT 0,
    linked_nodes                SMALLINT                 NOT NULL DEFAULT 0,
    staked                      DECIMAL(65, 0)           NOT NULL DEFAULT 0,
    pending                     DECIMAL(65, 0)           NOT NULL DEFAULT 0,
    rewards                     DECIMAL(65, 0)           NOT NULL DEFAULT 0,
    block_height                DECIMAL(65, 0)           NOT NULL,
    name                        TEXT                     ,
    description                 TEXT                     ,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_v_validator_id ON validators (validator_id);
CREATE INDEX idx_v_h ON validators (block_height);
