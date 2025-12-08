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

-- Specials table to store offering specials
CREATE TABLE IF NOT EXISTS specials (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  price REAL NOT NULL,
  currency TEXT NOT NULL DEFAULT 'USD',
  active INTEGER NOT NULL DEFAULT 1,
  starts_at TEXT,
  ends_at TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed initial specials if they don't exist
INSERT INTO specials(id, name, price, currency)
SELECT 'sp-1001', 'Winter Escape', 799.0, 'USD'
WHERE NOT EXISTS (SELECT 1 FROM specials WHERE id = 'sp-1001');

INSERT INTO specials(id, name, price, currency)
SELECT 'sp-1002', 'City Break Deluxe', 499.0, 'USD'
WHERE NOT EXISTS (SELECT 1 FROM specials WHERE id = 'sp-1002');
