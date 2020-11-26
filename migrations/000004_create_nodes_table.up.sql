CREATE TABLE IF NOT EXISTS nodes
(
    id                          UUID DEFAULT   uuid_generate_v4(),
    created_at                  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL,
    name                        TEXT                     NOT NULL,
    ip                          TEXT                     NOT NULL,
    public_ip                   TEXT                     NOT NULL,
    port                        SMALLINT                 NOT NULL,
    public_key                  TEXT                     NOT NULL,
    start_block                 DECIMAL(65, 0)           NOT NULL,
    last_reward_date            TIMESTAMP WITH TIME ZONE NOT NULL,
    finish_time                 TIMESTAMP WITH TIME ZONE NOT NULL,
    status                      TEXT                     NOT NULL,
    validator_id                DECIMAL(65, 0)           NOT NULL,
    PRIMARY KEY (id)
);
-- Indexes
CREATE index idx_nodes_validator_id on nodes (validator_id);

