CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--create types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'boundkind') THEN
        CREATE TYPE BOUNDKIND AS ENUM ('validator', 'delegator');
    END IF;
    --more types here...
END$$;

CREATE TABLE IF NOT EXISTS accounts
(
    id                              UUID DEFAULT   uuid_generate_v4(),
    created_at                      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    address                         NUMERIC(78)              NOT NULL,
    bound_kind                      BOUNDKIND                NOT NULL,
    bound_id                        DECIMAL(65, 0)           NOT NULL,
    block_height                    DECIMAL(65, 0)           NOT NULL,
    UNIQUE(address, bound_kind, bound_id, block_height),
    PRIMARY KEY (id)
);
