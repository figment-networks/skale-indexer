CREATE TABLE IF NOT EXISTS validators
(
    id                          UUID DEFAULT   uuid_generate_v4(),
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL,
    name                        TEXT                     NOT NULL,
    address                     TEXT                     NOT NULL,
    requested_address           TEXT                     NOT NULL,
    description                 TEXT                     NOT NULL,
    fee_rate                    DECIMAL(65, 0)           NOT NULL,
    registration_time           TIMESTAMP WITH TIME ZONE NOT NULL,
    minimum_delegation_amount   DECIMAL(65, 0)           NOT NULL,
    accept_new_requests         BOOLEAN                  NOT NULL,
    trusted                     BOOLEAN                  NOT NULL,
    PRIMARY KEY (id)
);
-- Indexes
CREATE index idx_validators_address on validators (address);
CREATE index idx_validators_requested_address on validators (requested_address);