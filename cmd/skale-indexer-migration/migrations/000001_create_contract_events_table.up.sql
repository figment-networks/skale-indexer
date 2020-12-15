CREATE TYPE BOUNDTYPE AS ENUM ('validator', 'delegation');

CREATE TABLE IF NOT EXISTS contract_events
(
    id                      UUID DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE,
    contract_name           VARCHAR(100)             NOT NULL,
    event_name              VARCHAR(50)              NOT NULL,
    contract_address        NUMERIC(78)[]            NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    time                    TIMESTAMP WITH TIME ZONE NOT NULL,
    transaction_hash        NUMERIC(125)[]           NOT NULL,
    params                  JSONB                    NOT NULL,
    removed                 BOOLEAN                  NOT NULL,
    bound_type              BOUNDTYPE                NOT NULL,
    bound_id                NUMERIC(78)[]            NOT NULL,
    bound_address           NUMERIC(78)[]            NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_contract_events_bound_id on contract_events USING GIN (bound_id);
