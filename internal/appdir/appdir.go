package appdir

import (
	"os"
	"path/filepath"
)

// Dir returns the yamp configuration directory (~/.config/yamp).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "yamp"), nil
}
