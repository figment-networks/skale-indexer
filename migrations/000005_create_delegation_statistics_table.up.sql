DROP TABLE IF EXISTS delegation_statistics;

CREATE TABLE IF NOT EXISTS delegation_statistics
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL,
    validator_id            DECIMAL(65, 0)           NOT NULL,
    status                  SMALLINT                 NOT NULL,
    amount                  DECIMAL(65, 0)           NOT NULL,
    statistics_type         SMALLINT                 NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_delegation_statistics_statistics_type_and_validator_id_and_status on delegation_statistics (statistics_type, validator_id, status);