package testrepo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/sanitize"
)

// GateFixture owns gate scripts and their declared external inputs.
type GateFixture struct {
	Environment []string
	Paths       []string
	tools       map[string]bool
	commandPath string
}

// NewGateFixture uses commandPath as the fixture's private command directory.
// The caller owns that directory's lifetime.
func NewGateFixture(commandPath string, environment ...string) *GateFixture {
	return &GateFixture{Environment: environment, commandPath: commandPath, tools: map[string]bool{}}
}

// Command returns a shell token and declares its ambient executable.
func (f *GateFixture) Command(name string) string {
	f.tools[name] = true
	return sanitize.ShellQuote(name)
}

// MustWrite materializes a fixture or fails the calling test.
func (f *GateFixture) MustWrite(t testing.TB, root, ordinary, prospective string) {
	t.Helper()
	if err := f.Write(root, ordinary, prospective); err != nil {
		t.Fatal(err)
	}
}

// Script gives a companion script the gate's private command path.
func (f *GateFixture) Script(body string) string {
	return "#!/bin/sh\nPATH=" + sanitize.ShellQuote(f.commandPath) + "\nexport PATH\n" + body
}

// Write writes each non-empty script body and one canonical input manifest.
func (f *GateFixture) Write(root, ordinary, prospective string) error {
	tools := make([]string, 0, len(f.tools))
	for name := range f.tools {
		tools = append(tools, name)
	}
	sort.Strings(tools)
	if err := os.MkdirAll(f.commandPath, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(f.commandPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !f.tools[entry.Name()] {
			return fmt.Errorf("private command path contains undeclared command %q", entry.Name())
		}
	}
	for _, name := range tools {
		if name == "" || strings.ContainsAny(name, "/\\\x00\r\n") || name == "." || name == ".." {
			return fmt.Errorf("invalid ambient command %q", name)
		}
		target, err := exec.LookPath(name)
		if err != nil {
			return fmt.Errorf("ambient command %q: %w", name, err)
		}
		target, err = filepath.Abs(target)
		if err != nil {
			return err
		}
		link := filepath.Join(f.commandPath, name)
		if old, err := os.Readlink(link); err == nil && old == target {
			continue
		}
		if err := os.Symlink(target, link); err != nil {
			return err
		}
	}
	dir := filepath.Join(root, ".bench")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, body := range map[string]string{
		"gate.sh": ordinary, "gate-prospective.sh": prospective,
	} {
		if body != "" {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(f.Script(body)), 0o755); err != nil {
				return err
			}
		}
	}
	inputs := struct {
		Schema      int      `json:"schema"`
		Closure     string   `json:"closure"`
		Environment []string `json:"environment"`
		Paths       []string `json:"paths"`
		Tools       []string `json:"tools"`
	}{1, "local", append([]string{}, f.Environment...), append([]string{}, f.Paths...), tools}
	data, err := json.Marshal(inputs)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "gate-inputs.json"), append(data, '\n'), 0o644)
}
