CREATE TABLE IF NOT EXISTS delegations
(
    id                  BIGSERIAL                NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    holder              TEXT                     NOT NULL,
    validator_id        DECIMAL(65, 0)           NOT NULL,
    amount              DECIMAL(65, 0)           NOT NULL,
    delegation_period   DECIMAL(65, 0)           NOT NULL,
    created             DECIMAL(65, 0)           NOT NULL,
    started             DECIMAL(65, 0)           NOT NULL,
    finished             DECIMAL(65, 0)           NOT NULL,
    info                TEXT                     NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_delegations_holder on delegations (holder);
CREATE index idx_delegations_validator_id on delegations (validator_id);