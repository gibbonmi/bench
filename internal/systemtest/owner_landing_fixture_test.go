//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/testrepo"
)

func systemLandingRaceFixture(t *testing.T) (root, home, tally, trees, ready, release string) {
	t.Helper()
	var err error
	root, err = os.MkdirTemp(owner.root, "landing-race [journey]-")
	if err != nil {
		t.Fatal(err)
	}
	if result := owner.runAt(root, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init landing race = (%d, %q)", result.code, result.stderr)
	}
	// The landing creates its own commit in this repository, so the identity belongs in
	// the config. A per-command -c leaves the product's commit without an author.
	for _, identity := range [][]string{{"user.email", "bench@local"}, {"user.name", "bench"}} {
		if result := owner.runAt(root, nil, "git", "config", identity[0], identity[1]); result.code != 0 {
			t.Fatalf("git config %s = (%d, %q)", identity[0], result.code, result.stderr)
		}
	}
	home, err = os.MkdirTemp(owner.root, "landing-race [home]-")
	if err != nil {
		t.Fatal(err)
	}
	tally = filepath.Join(home, "gate-tally")
	trees = filepath.Join(home, "gate-trees")
	ready = filepath.Join(home, "loser-ready")
	release = filepath.Join(home, "loser-release")
	if err := syscall.Mkfifo(ready, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(release, 0o600); err != nil {
		t.Fatal(err)
	}
	f := testrepo.NewGateFixture(t.TempDir())
	f.Environment = []string{"LAND_GATE_TALLY", "LAND_GATE_TREES", "LAND_RACE_READY", "LAND_RACE_RELEASE"}
	gate := "set -eu\nruntime=$1\n" + f.Command("grep") + " -q '^Status: implemented$' \"$runtime/specs/x/spec.md\"\ntree=$(" + f.Command("git") + " -C \"$runtime\" write-tree)\nprintf '%s\\n' \"$tree\" >> \"$LAND_GATE_TREES\"\nif [ -f \"$runtime/loser.txt\" ]; then\n  printf l >> \"$LAND_GATE_TALLY\"\n  if [ ! -f \"$runtime/winner.txt\" ]; then\n    printf r > \"$LAND_RACE_READY\"\n    IFS= read -r _ < \"$LAND_RACE_RELEASE\"\n  fi\nelse\n  printf w >> \"$LAND_GATE_TALLY\"\nfi\n"
	if err := f.Write(root, gate, gate); err != nil {
		t.Fatal(err)
	}
	specBody := "# x\n\nStatus: staged\n\n## User stories\n1. Land source.\n\n### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n| E1 | 1 | lands | command | catches failure |\n\n## Ownership fences\n\n- `loser.txt`\n- `winner.txt`\n- `reviews/x.md`\n"
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("retained-output\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	recordtest.Prepare(t, root, 1, "specs/x/spec.md", specBody)
	// The race writes the fenced files, so the ticket writes them too and the fence stays union-exact.
	ticket := filepath.Join(root, "specs", "x", "tickets", "1.md")
	data, err := os.ReadFile(ticket)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ticket, []byte(strings.Replace(string(data), "Writes: source.txt", "Writes: loser.txt (new), winner.txt (new)", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "commit", "-q", "-am", "race ticket writes")
	base := systemGitOutput(t, root, "rev-parse", "HEAD")
	systemGit(t, root, "update-ref", "refs/bench/green/main", base)
	return root, home, tally, trees, ready, release
}

func configureArtifactLandingFixture(t *testing.T, root string) {
	t.Helper()
	f := testrepo.NewGateFixture(t.TempDir())
	git := f.Command("git")
	// The selected executable resolves Go build inputs before the phase starts.
	f.Command("go")
	gate := "set -eu\nroot=${1:-$(" + git + " rev-parse --show-toplevel)}\nkit=${BENCH_KIT:?}\nbench=${BENCH_RUN_BINARY:?}\nexec " + f.Command("env") + " BENCH_KIT=\"$kit\" BENCH_RUN_BINARY=\"$bench\" \"$bench\" gate-phases \"$root\"\n"
	phase := "#!/bin/sh\nset -eu\nruntime=$(" + git + " rev-parse --show-toplevel)\n" + f.Command("grep") + " -q '^Status: implemented$' \"$runtime/specs/x/spec.md\"\ntree=$(" + git + " -C \"$runtime\" write-tree)\nprintf '%s\\n' \"$tree\" >> \"$LAND_GATE_TREES\"\nif [ -f \"$runtime/loser.txt\" ]; then\n  printf l >> \"$LAND_GATE_TALLY\"\nelse\n  printf w >> \"$LAND_GATE_TALLY\"\nfi\nprintf r > \"$LAND_RACE_READY\"\nIFS= read -r _ < \"$LAND_RACE_RELEASE\"\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "landing-race-phase.sh"), []byte(f.Script(phase)), 0o755); err != nil {
		t.Fatal(err)
	}
	phases := "{\"phases\":[{\"name\":\"landing-race\",\"argv\":[\".bench/landing-race-phase.sh\"]}]}\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "phases.json"), []byte(phases), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "cmd", "bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module landingrace\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The stub reads the two positionals off the end and creates the output directory, as
	// the real go-build.sh does; that script also takes options ahead of the positionals.
	dirname := f.Command("dirname")
	build := "set -eu\nshift $(($# - 2))\nroot=$1\nout=$2\nstaged=$out.staged\n" + f.Command("mkdir") + " -p \"$(" + dirname + " \"$out\")\"\n" + f.Command("cp") + " \"$LAND_BASELINE_BENCH\" \"$staged\"\n" + f.Command("chmod") + " 0700 \"$staged\"\n\"$staged\" freshness-publish \"$root\" \"$out\" \"$(" + dirname + " \"$out\")\" 1.2.3\n"
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.sh"), []byte(f.Script(build)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.inputs"), []byte("build_script=scripts/go-build.sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.Environment = []string{"HOME", "LAND_BASELINE_BENCH", "LAND_GATE_TALLY", "LAND_GATE_TREES", "LAND_RACE_READY", "LAND_RACE_RELEASE"}
	if err := f.Write(root, gate, gate); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "add", ".")
	systemGit(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "artifact recovery fixture")
	base := systemGitOutput(t, root, "rev-parse", "HEAD")
	systemGit(t, root, "update-ref", "refs/bench/green/main", base)
}
