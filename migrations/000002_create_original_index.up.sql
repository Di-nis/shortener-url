DROP INDEX IF EXISTS idx_short;
CREATE INDEX idx_original ON urls(original);