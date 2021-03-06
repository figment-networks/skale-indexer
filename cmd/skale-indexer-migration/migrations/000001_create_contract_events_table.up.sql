CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--create types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'boundtype') THEN
        CREATE TYPE BOUNDTYPE AS ENUM ('validator', 'delegation', 'node', 'token');
    END IF;
    --more types here...
END$$;

CREATE TABLE IF NOT EXISTS contract_events
(
    id                      UUID                     DEFAULT   uuid_generate_v4(),
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    contract_name           VARCHAR(100)             NOT NULL,
    event_name              VARCHAR(50)              NOT NULL,
    contract_address        NUMERIC(78)              NOT NULL,
    block_height            DECIMAL(65, 0)           NOT NULL,
    time                    TIMESTAMP WITH TIME ZONE NOT NULL,
    transaction_hash        NUMERIC(125)             NOT NULL,
    params                  JSONB                    NOT NULL,
    removed                 BOOLEAN                  NOT NULL,
    bound_type              BOUNDTYPE                NOT NULL,
    bound_id                NUMERIC(78)[],
    bound_address           NUMERIC(78)[],
    PRIMARY KEY (id)
);

-- Indexes
CREATE INDEX idx_c_ev_time ON contract_events (time);
CREATE INDEX idx_c_ev_bound_type ON contract_events (bound_type);
CREATE INDEX idx_c_ev_bound_id ON contract_events USING GIN (bound_id);

CREATE UNIQUE INDEX idx_c_ev_unique ON contract_events (contract_address, event_name, block_height, transaction_hash, removed);
