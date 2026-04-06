package heimdall

// ChangeNotifier is called when a configuration value changes.
// Implementations can react to changes (e.g., reload config, update UI).
type ChangeNotifier interface {
	NotifyChange(namespace, key, value, oldValue string)
}

type noOpNotifier struct{}

func (noOpNotifier) NotifyChange(_, _, _, _ string) {}
