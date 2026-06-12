package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// redirectConfigDir points the config base at a temp dir via XDG_CONFIG_HOME,
// which configBase honors on every platform.
func redirectConfigDir(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
}

func TestAppDirCreatesUnderConfigDir(t *testing.T) {
	redirectConfigDir(t)

	appDir, err := AppDir()
	if err != nil {
		t.Fatalf("AppDir: %v", err)
	}
	if filepath.Base(appDir) != appDirName {
		t.Errorf("AppDir = %q, want base %q", appDir, appDirName)
	}
	if parent := filepath.Base(filepath.Dir(appDir)); parent != appOrgName {
		t.Errorf("AppDir = %q, want parent %q", appDir, appOrgName)
	}
	if info, err := os.Stat(appDir); err != nil || !info.IsDir() {
		t.Errorf("AppDir %q was not created as a directory (err=%v)", appDir, err)
	}
}

func TestDefaultPathJoinsName(t *testing.T) {
	redirectConfigDir(t)

	appDir, err := AppDir()
	if err != nil {
		t.Fatalf("AppDir: %v", err)
	}
	if got, want := DefaultPath("config.json"), filepath.Join(appDir, "config.json"); got != want {
		t.Errorf("DefaultPath = %q, want %q", got, want)
	}
}

// TestAppDirFallsBackToDotConfig checks that, with XDG_CONFIG_HOME unset, Unix
// platforms (Linux and macOS) use ~/.config rather than os.UserConfigDir's
// macOS default of ~/Library/Application Support.
func TestAppDirFallsBackToDotConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows uses %AppData%, not ~/.config")
	}
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	appDir, err := AppDir()
	if err != nil {
		t.Fatalf("AppDir: %v", err)
	}
	if want := filepath.Join(home, ".config", appOrgName, appDirName); appDir != want {
		t.Errorf("AppDir = %q, want %q", appDir, want)
	}
}
