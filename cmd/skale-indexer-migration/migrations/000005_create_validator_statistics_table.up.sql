CREATE TABLE IF NOT EXISTS validator_statistics
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    validator_id            DECIMAL(65, 0)           NOT NULL,
    amount                  NUMERIC(125)             NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    time                    TIMESTAMP WITH TIME ZONE NOT NULL,
    statistic_type          SMALLINT                 NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_v_s_unique_st_vid_bh ON validator_statistics ( validator_id, block_height, statistic_type);
