package heimdall

// ConfigEntry represents a stored configuration value.
type ConfigEntry struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Type      string `json:"value_type"`
	Source    string `json:"source"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}

// ConfigSchema defines validation rules for a config key.
type ConfigSchema struct {
	Namespace   string   `json:"namespace"`
	Key         string   `json:"key"`
	Type        string   `json:"value_type"`
	Description string   `json:"description"`
	DefaultVal  string   `json:"default_val"`
	Required    bool     `json:"required"`
	Choices     []string `json:"choices,omitempty"`
}
