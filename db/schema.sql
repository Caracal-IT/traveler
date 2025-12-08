-- Safe to re-run; uses IF NOT EXISTS guards

-- Example table to verify DB creation; replace/extend as needed
CREATE TABLE IF NOT EXISTS app_metadata (
  key TEXT PRIMARY KEY,
  value TEXT,
  updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO app_metadata(key, value)
SELECT 'schema_version', '1'
WHERE NOT EXISTS (SELECT 1 FROM app_metadata WHERE key = 'schema_version');
