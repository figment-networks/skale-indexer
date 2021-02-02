ALTER TABLE validator_statistics ADD COLUMN time TIMESTAMP WITH TIME ZONE NOT NULL;
CREATE INDEX idx_v_s_time ON validator_statistics (time);
