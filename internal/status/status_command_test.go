// Tests for the status command's argument handling and route output.
package status

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/maps"
)

// Command rejects an unknown argument with a usage line and exit 2. It prints usage on
// -h, and accepts --all as the one added token. --all plus junk and near-misses stay
// usage errors, so a typo never silently prints the default board.
func TestCommandArgs(t *testing.T) {
	if r, c := Command([]string{"--bogus"}); c != 2 || !strings.Contains(r, "usage:") {
		t.Errorf("unknown arg: report %q exit %d", r, c)
	}
	if r, c := Command([]string{"-h"}); c != 0 || !strings.Contains(r, "usage: bench status") {
		t.Errorf("help: report %q exit %d", r, c)
	}
	// HC21: the help line advertises the record's four names with claude first.
	if r, c := Command([]string{"-h"}); !strings.Contains(r, "[--route [--harness claude|codex|none|opencode]]") {
		t.Errorf("help usage should advertise route grammar, got %q exit %d", r, c)
	}
	if r, c := Command([]string{"--all"}); c != 0 {
		t.Errorf("--all should be accepted with exit 0, got report %q exit %d", r, c)
	}
	for _, bad := range [][]string{{"--all", "extra"}, {"--allx"}, {"-a"}, {"--route", "--all"}, {"--harness", "codex"}, {"--route", "--harness", "cursor"}, {"--route", "extra"}} {
		if r, c := Command(bad); c != 2 || !strings.Contains(r, "usage:") {
			t.Errorf("args %q: report %q exit %d, want usage exit 2", bad, r, c)
		}
	}
}

func TestCommandMalformedRoutePrintsExactGrammar(t *testing.T) {
	want := grammar.Help + "\n"
	for _, args := range [][]string{
		{"--route", "--all"},
		{"--harness", "codex"},
		// HC23: an unrecorded harness prints the grammar and exits 2.
		{"--route", "--harness", "cursor"},
		{"--route", "extra"},
	} {
		if got, code := Command(args); code != 2 || got != want {
			t.Errorf("Command(%q) = (%q, %d), want (%q, 2)", args, got, code, want)
		}
	}
}

func TestCommandRouteOutsideRepositoryReturnsStructuredError(t *testing.T) {
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	if got, code := Command([]string{"--route"}); code != 1 || got != "error: not in a git repository — run inside a Bench-linked repo\n" {
		t.Fatalf("Command(--route) = (%q, %d)", got, code)
	}
}

func TestCommandRoutePrintsLeadAndRunnersUp(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "capture", "IDEAS.md"), []byte("- 2026-08-18 pending\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte("{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[]}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte("#!/bin/sh\nexit 7\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "tracked.txt")
	if result := gate.Execute(context.Background(), root, io.Discard, io.Discard); result.ActionExit != 7 {
		t.Fatalf("gate exit = %d, want red 7", result.ActionExit)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	want := "next[1]{state,why,command}:\n" +
		"  gate,red,/bench-debug\n" +
		"also: git (1 dirty path) → /bench-final-check; drain (1 idea(s), 0 open learning(s), 0 pending retro(s)) → /bench-drain\n"
	if got, code := Command([]string{"--route"}); code != 0 || got != want {
		t.Fatalf("Command(--route) = (%q, %d), want (%q, 0)", got, code, want)
	}
}

func TestCommandRouteEscapesControlBytesInProducerPaths(t *testing.T) {
	root := initRepo(t)
	writeFile := func(path, body string, mode os.FileMode) {
		t.Helper()
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(".bench/gate-inputs.json", "{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[]}\n", 0o644)
	writeFile(".bench/gate.sh", "#!/bin/sh\nexit 7\n", 0o755)
	writeFile("specs/my [draft]\x1b/spec.md", "Status: staged\n", 0o644)
	ready := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	ready = strings.Replace(ready, "Status: shaping", "Status: ready", 1)
	writeFile("decisions/my * map\x07.md", ready, 0o644)
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	if result := gate.Execute(context.Background(), root, io.Discard, io.Discard); result.ActionExit != 7 {
		t.Fatalf("gate exit = %d, want red 7", result.ActionExit)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	want := "next[1]{state,why,command}:\n" +
		"  gate,red,/bench-debug\n" +
		"also: specs (1 staged spec(s)) → /bench-implement-spec specs/my [draft]\\u001b/spec.md; decisions (1 ready map(s)) → /bench-write-spec decisions/my * map\\u0007.md\n"
	if got, code := Command([]string{"--route"}); code != 0 || got != want {
		t.Fatalf("Command(--route) = (%q, %d), want (%q, 0)", got, code, want)
	}
}

func TestCommandRouteEscapesControlByteInLeadStagedSpecPath(t *testing.T) {
	want := "next[1]{state,why,command}:\n" +
		"  specs,1 staged spec(s),\"/bench-implement-spec specs/lead\\\\u001b/spec.md\"\n" +
		"also: none\n"
	assertLeadControlRoute(t, "specs/lead\x1b/spec.md", "Status: staged\n", want)
}

func TestCommandRouteEscapesControlByteInLeadReadyMapPath(t *testing.T) {
	ready := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	ready = strings.Replace(ready, "Status: shaping", "Status: ready", 1)
	want := "next[1]{state,why,command}:\n" +
		"  decisions,1 ready map(s),\"/bench-write-spec decisions/lead\\\\u0007.md\"\n" +
		"also: none\n"
	assertLeadControlRoute(t, "decisions/lead\x07.md", ready, want)
}

func assertLeadControlRoute(t *testing.T, relativePath, body, want string) {
	t.Helper()
	root := initRepo(t)
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	if got, code := Command([]string{"--route"}); code != 0 || got != want {
		t.Fatalf("Command(--route) = (%q, %d), want (%q, 0)", got, code, want)
	}
}

// HC50. A harness with a phase form keeps the translation through the command surface.
func TestCommandRoutePrintsTheCodexPhaseForm(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "capture", "IDEAS.md"), []byte("- 2026-08-18 pending\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	got, code := Command([]string{"--route", "--harness", "codex"})
	if code != 0 || !strings.Contains(got, ",$bench-drain\n") {
		t.Fatalf("Command(--route --harness codex) = (%q, %d), want the $bench- form", got, code)
	}
}
