CREATE TABLE IF NOT EXISTS validators
(
    id                          UUID DEFAULT   uuid_generate_v4(),
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    validator_id                DECIMAL(65, 0)           NOT NULL,
    name                        TEXT                     NOT NULL,
    validator_address           NUMERIC(78)              NOT NULL,
    requested_address           NUMERIC(78)              NOT NULL,
    description                 TEXT                     NOT NULL,
    fee_rate                    DECIMAL(65, 0)           NOT NULL,
    registration_time           TIMESTAMP WITH TIME ZONE NOT NULL,
    minimum_delegation_amount   DECIMAL(65, 0)           NOT NULL,
    accept_new_requests         BOOLEAN                  NOT NULL,
    authorized                  BOOLEAN                  NOT NULL,
    active_nodes                SMALLINT                 NOT NULL,
    linked_nodes                SMALLINT                 NOT NULL,
    staked                      DECIMAL(65, 0)           NOT NULL,
    pending                     DECIMAL(65, 0)           NOT NULL,
    rewards                     DECIMAL(65, 0)           NOT NULL,
    block_height                DECIMAL(65, 0)           NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE index idx_v_validator_id on validators (validator_id);
