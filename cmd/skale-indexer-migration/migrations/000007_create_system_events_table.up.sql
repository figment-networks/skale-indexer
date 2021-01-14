CREATE TABLE IF NOT EXISTS system_events
(
    height                  DECIMAL(65, 0)           NOT NULL,
    kind                    SMALLINT                 NOT NULL,
    time                    TIMESTAMP WITH TIME ZONE NOT NULL,

    sender                  DECIMAL(65, 0)           NOT NULL,
    recipient               DECIMAL(65, 0)           NOT NULL,
    sender_id               DECIMAL(65, 0)           NOT NULL,
    recipient_id            DECIMAL(65, 0)           NOT NULL,

    before                  DECIMAL(65, 0)           NOT NULL,
    after                   DECIMAL(65, 0)           NOT NULL,
    change                  DECIMAL(65, 0)           NOT NULL
);

CREATE UNIQUE INDEX idx_sys_evt_unique ON system_events ( height, kind, sender, sender_id, recipient, recipient_id);
