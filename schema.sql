CREATE TABLE IF NOT EXISTS config (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace  TEXT NOT NULL,
    key        TEXT NOT NULL,
    value      TEXT NOT NULL,
    value_type TEXT NOT NULL,
    source     TEXT DEFAULT 'default',
    updated_at TEXT DEFAULT (datetime('now')),
    updated_by TEXT,
    UNIQUE(namespace, key)
);

CREATE TABLE IF NOT EXISTS config_schema (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace   TEXT NOT NULL,
    key         TEXT NOT NULL,
    value_type  TEXT NOT NULL,
    description TEXT,
    default_val TEXT,
    required    INTEGER DEFAULT 0,
    choices     TEXT,
    UNIQUE(namespace, key)
);

CREATE TABLE IF NOT EXISTS config_history (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    namespace  TEXT NOT NULL,
    key        TEXT NOT NULL,
    old_value  TEXT,
    new_value  TEXT,
    changed_by TEXT,
    changed_at TEXT DEFAULT (datetime('now'))
);
