CREATE TABLE IF NOT EXISTS events
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    smart_contract_address  NUMERIC(78)              NOT NULL,
    transaction_index       DECIMAL(65, 0)           NOT NULL,
    event_type              VARCHAR(50)              NOT NULL,
    event_name              TEXT                     NOT NULL,
    event_time              TIMESTAMP WITH TIME ZONE NOT NULL,
    event_info              JSONB                    NOT NULL,
    PRIMARY KEY (id)
);
