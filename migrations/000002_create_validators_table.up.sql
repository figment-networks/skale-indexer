CREATE TABLE IF NOT EXISTS validators
(
    id                          BIGSERIAL                NOT NULL,
    created_at                  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL,
    name                        TEXT                     NOT NULL,
    validator_address           TEXT                     NOT NULL,
    requested_address           TEXT                     NOT NULL,
    description                 TEXT                     NOT NULL,
    fee_rate                    DECIMAL(65, 0)           NOT NULL,
    registration_time           DECIMAL(65, 0)           NOT NULL,
    minimum_delegation_amount   DECIMAL(65, 0)           NOT NULL,
    accept_new_requests         BOOLEAN                  NOT NULL,
    PRIMARY KEY (id)
);
-- Indexes
CREATE index idx_validators_validator_address on validators (validator_address);
CREATE index idx_validators_requested_address on validators (requested_address);