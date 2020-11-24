CREATE TABLE IF NOT EXISTS events
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    smart_contract_address  TEXT                     NOT NULL,
    transaction_index       DECIMAL(65, 0)           NOT NULL,
    event_type              TEXT                     NOT NULL,
    event_name              TEXT                     NOT NULL,
    event_time              TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);
