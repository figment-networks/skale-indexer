CREATE TABLE IF NOT EXISTS nodes
(
    id                              UUID DEFAULT   uuid_generate_v4(),
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    node_id                         DECIMAL(65, 0)           NOT NULL,
    name                            VARCHAR(100)             NOT NULL,
    ip                              cidr               NOT NULL,
    public_ip                       cidr              NOT NULL,
    port                            SMALLINT                 NOT NULL,
    start_block                     DECIMAL(65, 0)           NOT NULL,
    next_reward_date                TIMESTAMP WITH TIME ZONE NOT NULL,
    last_reward_date                TIMESTAMP WITH TIME ZONE NOT NULL,
    finish_time                     DECIMAL(65, 0)           NOT NULL,
    status                          VARCHAR(50)              NOT NULL,
    validator_id                    DECIMAL(65, 0)           NOT NULL,
    event_time                      TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_nodes_validator_id_and_start_block on nodes (validator_id, start_block);
