package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	os.Unsetenv("XDG_CONFIG_HOME")
	AppFs = afero.NewMemMapFs()
	m.Run()
}

func TestXDGConfigDirIsPreferredFirst(t *testing.T) {
	t.Cleanup(func() {
		// reset fs after test
		AppFs = afero.NewMemMapFs()
	})

	xdgDirPath := filepath.FromSlash("/tmp-iofs/xdg/config")
	homeDirPath := filepath.FromSlash("/tmp-iofs/home/tester")
	if err := AppFs.MkdirAll(xdgDirPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := AppFs.MkdirAll(homeDirPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := afero.WriteFile(AppFs, filepath.Join(xdgDirPath, "revive.toml"), []byte("\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	if err := afero.WriteFile(AppFs, filepath.Join(homeDirPath, "revive.toml"), []byte("\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	setHome(t, homeDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(xdgDirPath, "revive.toml")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestHomeConfigDir(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	homeDirPath := filepath.FromSlash("/tmp-iofs/home/tester")
	if err := AppFs.MkdirAll(homeDirPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := afero.WriteFile(AppFs, filepath.Join(homeDirPath, "revive.toml"), []byte("\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	setHome(t, homeDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(homeDirPath, "revive.toml")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func setHome(t *testing.T, dir string) {
	t.Helper()

	homeEnv := "HOME"
	if runtime.GOOS == "windows" {
		homeEnv = "USERPROFILE"
	}
	t.Setenv(homeEnv, dir)
}

func TestXDGConfigDir(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	xdgDirPath := filepath.FromSlash("/tmp-iofs/xdg/config")
	if err := AppFs.MkdirAll(xdgDirPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := afero.WriteFile(AppFs, filepath.Join(xdgDirPath, "revive.toml"), []byte("\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(xdgDirPath, "revive.toml")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestXDGConfigDirNoFile(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	xdgDirPath := filepath.FromSlash("/tmp-iofs/xdg/config")
	t.Setenv("XDG_CONFIG_HOME", xdgDirPath)

	got := buildDefaultConfigPath()
	want := ""

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestXDGConfigDirUnsetUsesHomeBeforeCurrentDirectory(t *testing.T) {
	t.Cleanup(func() { AppFs = afero.NewMemMapFs() })
	t.Setenv("XDG_CONFIG_HOME", "")

	homeDirPath := filepath.FromSlash("/tmp-iofs/home/tester")
	if err := AppFs.MkdirAll(homeDirPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(AppFs, filepath.Join(homeDirPath, "revive.toml"), []byte("home config"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(AppFs, "revive.toml", []byte("current directory config"), 0o644); err != nil {
		t.Fatal(err)
	}
	setHome(t, homeDirPath)

	got := buildDefaultConfigPath()
	want := filepath.Join(homeDirPath, "revive.toml")
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetReleaseVersion(t *testing.T) {
	got := getVersion("builder", "2024-11-15 10:52 UTC", "7ee4500e125e2d1b12653b2c8e140fec380919b4", "v1.5.0-12-g7ee4500-dev")
	want := `Version:	v1.5.0-12-g7ee4500-dev
Commit:		7ee4500e125e2d1b12653b2c8e140fec380919b4
Built		2024-11-15 10:52 UTC by builder
`

	if got != want {
		t.Errorf("getVersion() = %q, want %q", got, want)
	}
}

func TestGetDevelopmentVersion(t *testing.T) {
	got := getVersion(defaultBuilder, defaultDate, defaultCommit, defaultVersion)
	want := "version (devel)\n"

	if got != want {
		t.Errorf("getVersion() = %q, want %q", got, want)
	}
}
