CREATE TABLE IF NOT EXISTS delegation_events
(
    id                  UUID DEFAULT   uuid_generate_v4(),
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    delegation_id       UUID                     NOT NULL,
    event_name          TEXT                     NOT NULL,
    event_time          TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_delegation_events_delegation_id on delegation_events (delegation_id);
