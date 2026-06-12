package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
)

// appOrgName and appDirName form the per-user subdirectory (emprcl/signls)
// under the config directory where signls stores its config and bank files.
const (
	appOrgName = "emprcl"
	appDirName = "signls"
)

// configBase returns the base config directory for the current OS. Unlike
// os.UserConfigDir (which is ~/Library/Application Support on macOS), terminal
// apps conventionally use ~/.config on both Linux and macOS, so we use that —
// honoring XDG_CONFIG_HOME — and fall back to %AppData% on Windows.
func configBase() (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir, nil
	}
	if runtime.GOOS == "windows" {
		return os.UserConfigDir()
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

// AppDir returns the per-user signls config directory, creating it if needed.
// It resolves to $XDG_CONFIG_HOME/emprcl/signls or ~/.config/emprcl/signls on
// Linux and macOS, and %AppData%\emprcl\signls on Windows.
func AppDir() (string, error) {
	base, err := configBase()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(base, appOrgName, appDirName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return appDir, nil
}

// DefaultPath resolves name inside AppDir. If the config directory can't be
// determined, it falls back to the bare name (the current working directory),
// preserving the previous behavior.
func DefaultPath(name string) string {
	dir, err := AppDir()
	if err != nil {
		return name
	}
	return filepath.Join(dir, name)
}
