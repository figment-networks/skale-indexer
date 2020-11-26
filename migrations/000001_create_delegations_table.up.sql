CREATE TABLE IF NOT EXISTS delegations
(
    id                  UUID DEFAULT   uuid_generate_v4(),
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    holder              NUMERIC(78)              NOT NULL,
    validator_id        DECIMAL(65, 0)           NOT NULL,
    amount              DECIMAL(65, 0)           NOT NULL,
    delegation_period   DECIMAL(65, 0)           NOT NULL,
    created             TIMESTAMP WITH TIME ZONE NOT NULL,
    started             TIMESTAMP WITH TIME ZONE NOT NULL,
    finished            TIMESTAMP WITH TIME ZONE NOT NULL,
    info                TEXT                     NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_delegations_holder on delegations (holder);
CREATE index idx_delegations_validator_id on delegations (validator_id);