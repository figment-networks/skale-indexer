CREATE TABLE IF NOT EXISTS delegations
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delegation_id           DECIMAL(65, 0)           NOT NULL,
    holder                  NUMERIC(78)              NOT NULL,
    validator_id            DECIMAL(65, 0)           NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    transaction_hash        NUMERIC(125)             NOT NULL,
    amount                  DECIMAL(65, 0)           NOT NULL,
    delegation_period       DECIMAL(65, 0)           NOT NULL,
    created                 TIMESTAMP WITH TIME ZONE NOT NULL,
    started                 DECIMAL(65, 0)           NOT NULL,
    finished                DECIMAL(65, 0)           NOT NULL,
    info                    TEXT                     NOT NULL,
    state                   SMALLINT                 NOT NULL,
    PRIMARY KEY (id)
);

-- TODO: unique constraints?
-- Indexes
CREATE INDEX idx_del_h ON delegations (holder);
CREATE INDEX idx_del_v_id_bl_height ON delegations (validator_id, block_height);
CREATE UNIQUE INDEX idx_del_unique ON delegations (delegation_id, transaction_hash);
CREATE INDEX idx_del_created ON delegations (created);
