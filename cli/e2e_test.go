// Package cli_test runs revive as a separate process and checks its exit code and output.
package cli_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mgechev/revive/cli"
)

func TestMain(m *testing.M) {
	if os.Getenv("REVIVE_E2E_RUN_REVIVE") == "1" {
		cli.RunRevive()
		return
	}

	m.Run()
}

type reviveResult struct {
	exitCode int
	stdout   string
	stderr   string
}

func runRevive(t *testing.T, homeConfig, xdgConfig string, args ...string) reviveResult {
	t.Helper()

	self, err := os.Executable()
	if err != nil {
		t.Fatalf("getting executable: %v", err)
	}

	cmd := exec.CommandContext(t.Context(), self, args...)
	cmd.Dir = "testdata"
	home, xdg := t.TempDir(), t.TempDir()
	writeConfig(t, home, homeConfig)
	writeConfig(t, xdg, xdgConfig)
	cmd.Env = []string{"REVIVE_E2E_RUN_REVIVE=1", "HOME=" + home, "USERPROFILE=" + home, "XDG_CONFIG_HOME=" + xdg}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return reviveResult{
			exitCode: exitErr.ExitCode(),
			stdout:   stdoutBuf.String(),
			stderr:   stderrBuf.String(),
		}
	}
	if err != nil {
		t.Fatalf("running revive: %v", err)
	}
	return reviveResult{
		stdout: stdoutBuf.String(),
		stderr: stderrBuf.String(),
	}
}

func writeConfig(t *testing.T, dir, content string) {
	t.Helper()

	if content == "" {
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "revive.toml"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRevive(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		homeConfig   string // revive.toml in $HOME
		xdgConfig    string // revive.toml in $XDG_CONFIG_HOME
		wantExitCode int
		wantStdout   []string // substrings; no entry means stdout must be empty
		wantStderr   []string // substrings
	}{
		{
			name:       "version",
			args:       []string{"-version"},
			wantStdout: []string{"version"},
		},
		{
			name: "clean file",
			args: []string{"clean.go"},
		},
		{
			name: "clean file with exit status",
			args: []string{"-set_exit_status", "clean.go"},
		},
		{
			name:       "failures are reported but exit code is zero",
			args:       []string{"bad.go"},
			wantStdout: []string{"bad.go:1:1: should have a package comment", "bad.go:3:1: exported function Bad should have comment or be unexported"},
		},
		{
			name:         "failures with exit status",
			args:         []string{"-set_exit_status", "bad.go"},
			wantExitCode: 1,
			wantStdout:   []string{"bad.go:3:1: exported function Bad should have comment or be unexported"},
		},
		{
			name:       "config enables only one rule",
			args:       []string{"-config", "exported.toml", "bad.go"},
			wantStdout: []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name: "config from $XDG_CONFIG_HOME",
			xdgConfig: `error-code = 5
		[rule.exported]
severity = "error"`,
			args:         []string{"bad.go"},
			wantExitCode: 5,
			wantStdout:   []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name: "config from $HOME",
			homeConfig: `warning-code = 3
[rule.exported]`,
			args:         []string{"bad.go"},
			wantExitCode: 3,
			wantStdout:   []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name: "config from $XDG_CONFIG_HOME wins over the one from $HOME",
			homeConfig: `warning-code = 3
[rule.exported]`,
			xdgConfig: `error-code = 5
		[rule.exported]
severity = "error"`,
			args:         []string{"bad.go"},
			wantExitCode: 5,
			wantStdout:   []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name:       "no config file in $XDG_CONFIG_HOME or $HOME enables the default rules",
			args:       []string{"bad.go"},
			wantStdout: []string{"bad.go:1:1: should have a package comment"},
		},
		{
			name: "excluded file",
			args: []string{"-set_exit_status", "-exclude", "bad.go", "./..."},
		},
		{
			name:         "warning-code",
			args:         []string{"-config", "warning_code.toml", "bad.go"},
			wantExitCode: 3,
			wantStdout:   []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name:         "error-code",
			args:         []string{"-config", "error_code.toml", "bad.go"},
			wantExitCode: 5,
			wantStdout:   []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name:         "error-code wins over warning-code",
			args:         []string{"-config", "both_codes.toml", "bad.go"},
			wantExitCode: 5,
			wantStdout:   []string{"should have a package comment", "exported function Bad should have comment or be unexported"},
		},
		{
			name: "no failures ignore error-code and warning-code",
			args: []string{"-config", "both_codes.toml", "clean.go"},
		},
		{
			name:         "set_exit_status overrides error-code and warning-code",
			args:         []string{"-set_exit_status", "-config", "both_codes.toml", "bad.go"},
			wantExitCode: 1,
			wantStdout:   []string{"exported function Bad should have comment or be unexported"},
		},
		{
			name: "multiple excludes",
			args: []string{"-set_exit_status", "-exclude", "bad.go", "-exclude", "clean.go", "./..."},
		},
		{
			name:       "max open files",
			args:       []string{"-max_open_files", "1", "bad.go"},
			wantStdout: []string{"bad.go:3:1: exported function Bad should have comment or be unexported"},
		},
		{
			name:         "invalid max open files",
			args:         []string{"-max_open_files", "many", "bad.go"},
			wantExitCode: 2,
			wantStdout:   []string{"Example:"}, // a flag error prints the banner on stdout and the usage on stderr
			wantStderr:   []string{`invalid value "many" for flag -max_open_files`},
		},
		{
			name:       "help",
			args:       []string{"-help"},
			wantStdout: []string{"Example:"},
			wantStderr: []string{"-config", "-exclude", "-formatter", "-max_open_files", "-set_exit_status", "-version"},
		},
		{
			name:       "json formatter",
			args:       []string{"-formatter", "json", "bad.go"},
			wantStdout: []string{`"RuleName":"exported","Category":"comments"`},
		},
		{
			name:         "missing config file",
			args:         []string{"-config", "missing.toml", "clean.go"},
			wantExitCode: 1,
			wantStderr:   []string{"cannot read the config file"},
		},
		{
			name:         "invalid config file",
			args:         []string{"-config", "invalid.toml", "clean.go"},
			wantExitCode: 1,
		},
		{
			name:         "unknown formatter",
			args:         []string{"-formatter", "unknown", "bad.go"},
			wantExitCode: 1,
			wantStderr:   []string{"unknown formatter unknown"},
		},
		{
			name:         "unknown flag",
			args:         []string{"-unknown"},
			wantExitCode: 2,
			wantStdout:   []string{"Example:"},
			wantStderr:   []string{"flag provided but not defined: -unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runRevive(t, tt.homeConfig, tt.xdgConfig, tt.args...)

			if got.exitCode != tt.wantExitCode {
				t.Errorf("exit code = %d, want %d\nstdout: %s\nstderr: %s", got.exitCode, tt.wantExitCode, got.stdout, got.stderr)
			}
			if len(tt.wantStdout) == 0 && got.stdout != "" {
				t.Errorf("stdout = %q, want empty", got.stdout)
			}
			for _, w := range tt.wantStdout {
				if !strings.Contains(got.stdout, w) {
					t.Errorf("stdout = %q, want it to contain %q", got.stdout, w)
				}
			}
			for _, w := range tt.wantStderr {
				if !strings.Contains(got.stderr, w) {
					t.Errorf("stderr = %q, want it to contain %q", got.stderr, w)
				}
			}
		})
	}
}
