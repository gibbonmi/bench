package repairtest

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/gittest"
	toonlib "github.com/toon-format/toon-go"
	"os"
	"path/filepath"
	"testing"
)

type session struct{ root, home string }

func newSession(t *testing.T) session {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := gittest.Repo(t)
	home := t.TempDir()
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(cwd, "..", "..", "..")))
	t.Setenv("HOME", home)
	t.Setenv("BENCH_HOME", filepath.Join(home, "bench-state"))
	t.Setenv("CODEX_HOME", filepath.Join(home, "personal-codex"))
	t.Chdir(root)
	return session{root, home}
}

func linkedSession(t *testing.T) session {
	t.Helper()
	s := newSession(t)
	s.link(t, 0)
	return s
}

func (s session) link(t *testing.T, want int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := adopt.Link(nil, &stdout, &stderr, "1.0.0"); code != want {
		t.Fatalf("fixture link = %d, want %d: %s %s", code, want, &stdout, &stderr)
	}
}

func (s session) doctor(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := adopt.Doctor(append([]string{"--compat", "codex-cli"}, args...), &stdout, &stderr, "1.0.0")
	return code, stdout.String(), stderr.String()
}

func repairID(t *testing.T, output string) string {
	t.Helper()
	value, err := toonlib.DecodeString(output)
	if err != nil {
		t.Fatal(err)
	}
	document, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("repair output is not a document: %T", value)
	}
	rows, ok := document["repair"].([]any)
	if !ok || len(rows) != 1 {
		t.Fatalf("repair output has no single result: %#v", document["repair"])
	}
	row, ok := rows[0].(map[string]any)
	if !ok {
		t.Fatalf("repair result is not a row: %#v", rows[0])
	}
	id, ok := row["id"].(string)
	if !ok || id == "" {
		t.Fatalf("repair result has no recovery identifier: %#v", row)
	}
	return id
}

func (s session) track(t *testing.T) {
	t.Helper()
	gittest.Output(t, s.root, "add", "--all")
	gittest.Output(t, s.root, "-c", "user.name=Bench fixture", "-c", "user.email=bench@example.invalid", "commit", "-qm", "fixture")
}
