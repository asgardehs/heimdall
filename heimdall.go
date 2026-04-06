package heimdall

import (
	"fmt"
	"log/slog"
)

// Heimdall is the configuration manager for the Asgard ecosystem.
type Heimdall struct {
	store    *store
	notifier ChangeNotifier
	logger   *slog.Logger
}

// Option configures a Heimdall instance.
type Option func(*Heimdall)

// WithNotifier sets a change notifier that is called when config values change.
func WithNotifier(n ChangeNotifier) Option {
	return func(h *Heimdall) { h.notifier = n }
}

// Open creates a Heimdall instance using the platform default database path.
func Open(logger *slog.Logger, opts ...Option) (*Heimdall, error) {
	return OpenPath(DefaultDBPath(), logger, opts...)
}

// OpenPath creates a Heimdall instance using a specific database path.
func OpenPath(dbPath string, logger *slog.Logger, opts ...Option) (*Heimdall, error) {
	s, err := openStore(dbPath, logger)
	if err != nil {
		return nil, fmt.Errorf("opening store: %w", err)
	}

	h := &Heimdall{
		store:    s,
		notifier: noOpNotifier{},
		logger:   logger,
	}
	for _, opt := range opts {
		opt(h)
	}

	if err := seedDefaults(s); err != nil {
		s.close()
		return nil, fmt.Errorf("seeding defaults: %w", err)
	}

	logger.Info("heimdall initialized", "path", dbPath)
	return h, nil
}

// Get retrieves a single config entry.
func (h *Heimdall) Get(namespace, key string) (*ConfigEntry, error) {
	return h.store.get(namespace, key)
}

// Set writes a config value. It validates against the schema, records history,
// and notifies the change notifier.
func (h *Heimdall) Set(namespace, key, value string) error {
	// Validate against schema if one exists.
	schemas, _ := h.store.getSchema(namespace)
	for _, s := range schemas {
		if s.Key == key {
			if err := ValidateValue(s, value); err != nil {
				return err
			}
			break
		}
	}

	oldValue, err := h.store.set(namespace, key, value, "user", "user")
	if err != nil {
		return err
	}

	h.notifier.NotifyChange(namespace, key, value, oldValue)
	return nil
}

// List returns all config entries for a namespace.
func (h *Heimdall) List(namespace string) ([]ConfigEntry, error) {
	return h.store.list(namespace)
}

// Reset sets a config key back to its schema default.
func (h *Heimdall) Reset(namespace, key string) error {
	oldEntry, _ := h.store.get(namespace, key)
	oldValue := ""
	if oldEntry != nil {
		oldValue = oldEntry.Value
	}

	_, err := h.store.reset(namespace, key)
	if err != nil {
		return err
	}

	entry, _ := h.store.get(namespace, key)
	if entry != nil {
		h.notifier.NotifyChange(namespace, key, entry.Value, oldValue)
	}
	return nil
}

// Schema returns all schema entries for a namespace.
func (h *Heimdall) Schema(namespace string) ([]ConfigSchema, error) {
	return h.store.getSchema(namespace)
}

// Close closes the underlying database.
func (h *Heimdall) Close() error {
	return h.store.close()
}
