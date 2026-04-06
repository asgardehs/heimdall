package heimdall

import (
	"os"
	"path/filepath"
	"runtime"
)

// DefaultDBPath returns the platform-specific default path for the Heimdall
// SQLite database.
func DefaultDBPath() string {
	return filepath.Join(dataDir(), "heimdall.db")
}

func dataDir() string {
	if p := os.Getenv("HEIMDALL_DATA_DIR"); p != "" {
		return p
	}
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "heimdall")
	case "windows":
		if dir := os.Getenv("APPDATA"); dir != "" {
			return filepath.Join(dir, "heimdall")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Roaming", "heimdall")
	default:
		if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
			return filepath.Join(dir, "heimdall")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local", "share", "heimdall")
	}
}
