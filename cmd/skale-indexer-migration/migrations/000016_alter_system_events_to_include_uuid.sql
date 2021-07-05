ALTER TABLE system_events DROP COLUMN IF EXISTS id;
ALTER TABLE system_events ADD COLUMN id uuid DEFAULT uuid_generate_v4();
