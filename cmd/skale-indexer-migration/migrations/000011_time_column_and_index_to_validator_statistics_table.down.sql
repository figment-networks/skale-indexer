DROP INDEX IF EXISTS idx_v_s_time;
ALTER TABLE validator_statistics DROP COLUMN IF EXISTS time;
