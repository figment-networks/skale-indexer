CREATE TABLE IF NOT EXISTS delegations
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delegation_id           DECIMAL(65, 0)           NOT NULL,
    holder                  NUMERIC(78)              NOT NULL,
    validator_id            DECIMAL(65, 0)           NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    amount                  DECIMAL(65, 0)           NOT NULL,
    delegation_period       DECIMAL(65, 0)           NOT NULL,
    created                 TIMESTAMP WITH TIME ZONE NOT NULL,
    started                 DECIMAL(65, 0)           NOT NULL,
    finished                DECIMAL(65, 0)           NOT NULL,
    info                    TEXT                     NOT NULL,
    state                   SMALLINT                 NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_delegations_holder on delegations (holder);
CREATE UNIQUE index idx_delegations_delegation_id_and_block_height on delegations (delegation_id, block_height);
