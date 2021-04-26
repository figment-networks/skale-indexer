ALTER TABLE delegations ADD COLUMN until TIMESTAMP WITH TIME ZONE;
UPDATE delegations SET until = (date_trunc('month', created::date) + make_interval(MONTHS => 1+delegation_period::INTEGER) - interval '1 day')::TIMESTAMP;
CREATE INDEX idx_d_until ON delegations (until);
