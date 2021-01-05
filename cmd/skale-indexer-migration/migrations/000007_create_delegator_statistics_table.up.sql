CREATE TABLE IF NOT EXISTS delegator_statistics
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    holder                  NUMERIC(78)              NOT NULL,
    amount                  DECIMAL(65, 0)           NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    statistics_type         SMALLINT                 NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_ds_statistics_type_and_holder_and_block_height ON delegator_statistics (statistics_type, holder, block_height);