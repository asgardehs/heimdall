package heimdall

import "fmt"

// Defaults contains the initial config values for all ecosystem tools.
var Defaults = []ConfigEntry{
	// Odin
	{Namespace: "odin", Key: "database_path", Value: "", Type: "path"},
	{Namespace: "odin", Key: "backup_path", Value: "", Type: "path"},
	{Namespace: "odin", Key: "theme", Value: "system", Type: "string"},
	{Namespace: "odin", Key: "auto_backup", Value: "true", Type: "boolean"},

	// Muninn
	{Namespace: "muninn", Key: "vault_path", Value: "", Type: "path"},
	{Namespace: "muninn", Key: "database_path", Value: "", Type: "path"},
	{Namespace: "muninn", Key: "model_name", Value: "nomic-ai/nomic-embed-text-v1.5", Type: "string"},

	// Huginn
	{Namespace: "huginn", Key: "default_theme", Value: "odin_compliance", Type: "string"},
	{Namespace: "huginn", Key: "output_dir", Value: "", Type: "path"},

	// AI (cross-cutting, opt-in — off by default)
	{Namespace: "ai", Key: "enabled", Value: "false", Type: "boolean"},
	{Namespace: "ai", Key: "provider", Value: "", Type: "enum"},
	{Namespace: "ai", Key: "api_key", Value: "", Type: "secret"},
	{Namespace: "ai", Key: "odin_access", Value: "[]", Type: "array"},
	{Namespace: "ai", Key: "muninn_access", Value: "[]", Type: "array"},
	{Namespace: "ai", Key: "huginn_access", Value: "[]", Type: "array"},
}

// DefaultSchemas contains the schema definitions corresponding to Defaults.
var DefaultSchemas = []ConfigSchema{
	// Odin
	{Namespace: "odin", Key: "database_path", Type: "path", Description: "Path to Odin SQLite database", DefaultVal: ""},
	{Namespace: "odin", Key: "backup_path", Type: "path", Description: "Path for database backups", DefaultVal: ""},
	{Namespace: "odin", Key: "theme", Type: "string", Description: "UI theme (system, light, dark)", DefaultVal: "system"},
	{Namespace: "odin", Key: "auto_backup", Type: "boolean", Description: "Enable automatic backups", DefaultVal: "true"},

	// Muninn
	{Namespace: "muninn", Key: "vault_path", Type: "path", Description: "Path to Muninn note vault", DefaultVal: ""},
	{Namespace: "muninn", Key: "database_path", Type: "path", Description: "Path to Muninn SQLite database", DefaultVal: ""},
	{Namespace: "muninn", Key: "model_name", Type: "string", Description: "Embedding model identifier", DefaultVal: "nomic-ai/nomic-embed-text-v1.5"},

	// Huginn
	{Namespace: "huginn", Key: "default_theme", Type: "string", Description: "Default rendering theme", DefaultVal: "odin_compliance"},
	{Namespace: "huginn", Key: "output_dir", Type: "path", Description: "Output directory for rendered files", DefaultVal: ""},

	// AI (cross-cutting, opt-in)
	{Namespace: "ai", Key: "enabled", Type: "boolean", Description: "Master AI enable/disable switch (off by default)", DefaultVal: "false"},
	{Namespace: "ai", Key: "provider", Type: "enum", Description: "AI provider", DefaultVal: "", Choices: []string{"anthropic", "openai"}},
	{Namespace: "ai", Key: "api_key", Type: "secret", Description: "AI provider API key", DefaultVal: ""},
	{Namespace: "ai", Key: "odin_access", Type: "array", Description: "Odin modules AI can access (e.g. incidents, chemicals, training, compliance, reports)", DefaultVal: "[]"},
	{Namespace: "ai", Key: "muninn_access", Type: "array", Description: "Muninn features AI can access (e.g. search, notes, snippets)", DefaultVal: "[]"},
	{Namespace: "ai", Key: "huginn_access", Type: "array", Description: "Huginn features AI can access", DefaultVal: "[]"},
}

// seedDefaults registers schemas and inserts default config values for any
// keys that don't already exist.
func seedDefaults(s *store) error {
	if err := s.registerSchema(DefaultSchemas); err != nil {
		return fmt.Errorf("registering schemas: %w", err)
	}

	for _, d := range Defaults {
		if _, err := s.get(d.Namespace, d.Key); err == nil {
			continue // already exists
		}
		if _, err := s.set(d.Namespace, d.Key, d.Value, "default", "heimdall"); err != nil {
			return fmt.Errorf("seeding default %s.%s: %w", d.Namespace, d.Key, err)
		}
	}
	return nil
}
