package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestXDGConfigDirUnsetUsesHomeBeforeCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.Mkdir(home, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "revive.toml"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "revive.toml"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Chdir(dir)
	if got, want := buildDefaultConfigPath(), filepath.Join(home, "revive.toml"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
