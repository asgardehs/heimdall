package heimdall

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// store is the SQLite-backed configuration store.
type store struct {
	db     *sql.DB
	logger *slog.Logger
}

func openStore(dbPath string, logger *slog.Logger) (*store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("initializing schema: %w", err)
	}

	// Migrate: add choices column if it doesn't exist (idempotent).
	_, _ = db.Exec(`ALTER TABLE config_schema ADD COLUMN choices TEXT`)

	logger.Debug("heimdall store opened", "path", dbPath)
	return &store{db: db, logger: logger}, nil
}

func (s *store) get(namespace, key string) (*ConfigEntry, error) {
	row := s.db.QueryRow(
		`SELECT namespace, key, value, value_type, source, updated_at, COALESCE(updated_by, '')
		 FROM config WHERE namespace = ? AND key = ?`,
		namespace, key,
	)

	var e ConfigEntry
	err := row.Scan(&e.Namespace, &e.Key, &e.Value, &e.Type, &e.Source, &e.UpdatedAt, &e.UpdatedBy)
	if err == sql.ErrNoRows {
		return nil, &ConfigError{
			Code:    ErrConfigNotFound,
			Message: fmt.Sprintf("config not found: %s.%s", namespace, key),
		}
	}
	if err != nil {
		return nil, fmt.Errorf("querying config: %w", err)
	}
	return &e, nil
}

func (s *store) set(namespace, key, value, source, updatedBy string) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current value for history (may not exist).
	var oldValue string
	err = tx.QueryRow(
		`SELECT value FROM config WHERE namespace = ? AND key = ?`,
		namespace, key,
	).Scan(&oldValue)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("querying old value: %w", err)
	}

	// Look up the value_type from schema, fall back to "string".
	var valueType string
	err = tx.QueryRow(
		`SELECT value_type FROM config_schema WHERE namespace = ? AND key = ?`,
		namespace, key,
	).Scan(&valueType)
	if err != nil {
		valueType = "string"
	}

	// Upsert config.
	_, err = tx.Exec(
		`INSERT INTO config (namespace, key, value, value_type, source, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
		 ON CONFLICT(namespace, key) DO UPDATE SET
		     value = excluded.value,
		     source = excluded.source,
		     updated_by = excluded.updated_by,
		     updated_at = excluded.updated_at`,
		namespace, key, value, valueType, source, updatedBy,
	)
	if err != nil {
		return "", fmt.Errorf("upserting config: %w", err)
	}

	// Record history.
	_, err = tx.Exec(
		`INSERT INTO config_history (namespace, key, old_value, new_value, changed_by)
		 VALUES (?, ?, ?, ?, ?)`,
		namespace, key, oldValue, value, updatedBy,
	)
	if err != nil {
		return "", fmt.Errorf("recording history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("committing transaction: %w", err)
	}
	return oldValue, nil
}

func (s *store) list(namespace string) ([]ConfigEntry, error) {
	rows, err := s.db.Query(
		`SELECT namespace, key, value, value_type, source, updated_at, COALESCE(updated_by, '')
		 FROM config WHERE namespace = ? ORDER BY key`,
		namespace,
	)
	if err != nil {
		return nil, fmt.Errorf("listing config: %w", err)
	}
	defer rows.Close()

	var entries []ConfigEntry
	for rows.Next() {
		var e ConfigEntry
		if err := rows.Scan(&e.Namespace, &e.Key, &e.Value, &e.Type, &e.Source, &e.UpdatedAt, &e.UpdatedBy); err != nil {
			return nil, fmt.Errorf("scanning config row: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (s *store) reset(namespace, key string) (string, error) {
	var defaultVal string
	err := s.db.QueryRow(
		`SELECT default_val FROM config_schema WHERE namespace = ? AND key = ?`,
		namespace, key,
	).Scan(&defaultVal)
	if err == sql.ErrNoRows {
		return "", &ConfigError{
			Code:    ErrConfigNotFound,
			Message: fmt.Sprintf("no schema default for: %s.%s", namespace, key),
		}
	}
	if err != nil {
		return "", fmt.Errorf("querying default: %w", err)
	}

	return s.set(namespace, key, defaultVal, "default", "heimdall")
}

func (s *store) getSchema(namespace string) ([]ConfigSchema, error) {
	rows, err := s.db.Query(
		`SELECT namespace, key, value_type, COALESCE(description, ''), COALESCE(default_val, ''), required, choices
		 FROM config_schema WHERE namespace = ? ORDER BY key`,
		namespace,
	)
	if err != nil {
		return nil, fmt.Errorf("listing schema: %w", err)
	}
	defer rows.Close()

	var schemas []ConfigSchema
	for rows.Next() {
		var cs ConfigSchema
		var req int
		var choicesJSON sql.NullString
		if err := rows.Scan(&cs.Namespace, &cs.Key, &cs.Type, &cs.Description, &cs.DefaultVal, &req, &choicesJSON); err != nil {
			return nil, fmt.Errorf("scanning schema row: %w", err)
		}
		cs.Required = req != 0
		if choicesJSON.Valid && choicesJSON.String != "" {
			json.Unmarshal([]byte(choicesJSON.String), &cs.Choices)
		}
		schemas = append(schemas, cs)
	}
	return schemas, rows.Err()
}

func (s *store) registerSchema(entries []ConfigSchema) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO config_schema (namespace, key, value_type, description, default_val, required, choices)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, e := range entries {
		req := 0
		if e.Required {
			req = 1
		}
		var choicesJSON *string
		if len(e.Choices) > 0 {
			data, _ := json.Marshal(e.Choices)
			s := string(data)
			choicesJSON = &s
		}
		if _, err := stmt.Exec(e.Namespace, e.Key, e.Type, e.Description, e.DefaultVal, req, choicesJSON); err != nil {
			return fmt.Errorf("inserting schema %s.%s: %w", e.Namespace, e.Key, err)
		}
	}

	return tx.Commit()
}

func (s *store) close() error {
	s.logger.Debug("heimdall store closed")
	return s.db.Close()
}
