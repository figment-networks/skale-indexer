DROP TABLE IF EXISTS validator_statistics;

CREATE TABLE IF NOT EXISTS validator_statistics
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE,
    validator_id            DECIMAL(65, 0)           NOT NULL,
    amount                  DECIMAL(65, 0)           NOT NULL,
    eth_block_height        DECIMAL(65, 0)           NOT NULL,
    statistics_type         SMALLINT                 NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_statistics_statistics_type_and_validator_id on validator_statistics (statistics_type, validator_id);