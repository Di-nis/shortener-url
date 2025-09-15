DROP INDEX IF EXISTS idx_original;
CREATE INDEX idx_short ON urls(short);