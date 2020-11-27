DROP TABLE IF EXISTS delegation_state_statistics;

CREATE TABLE IF NOT EXISTS delegation_state_statistics
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL,
    validator_id            DECIMAL(65, 0)           NOT NULL,
    status                  SMALLINT                 NOT NULL,
    amount                  DECIMAL(65, 0)           NOT NULL,
    UNIQUE (validator_id, status),
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_delegation_state_statistics_validator_id_and_status on delegations (validator_id, status);