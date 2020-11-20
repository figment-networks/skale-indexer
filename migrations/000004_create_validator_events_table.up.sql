CREATE TABLE IF NOT EXISTS validator_events
(
    id                  UUID DEFAULT   uuid_generate_v4(),
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    validator_id        UUID                     NOT NULL,
    event_name          TEXT                     NOT NULL,
    event_time          TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_events_validator_id on validator_events (validator_id);
