package skillsindex

import (
	"os/exec"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// Command is the operator's only regenerator now that the retired shell script is gone:
// --check (the default) reports drift and exits 1, --write clears it and exits 0, and
// a following --check on the freshly written tree reports nothing.
func TestCommand(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
	run(t, root, "add", ".")
	run(t, root, "commit", "-q", "-m", "fixture base")
	t.Chdir(root)

	wantCheck := "skills index missing entry for skill 'alpha' (regenerate: bench skills-index --write)\n"
	report, code := Command(nil)
	if report != wantCheck || code != 1 {
		t.Fatalf("no-args (default --check) on drifted tree = %q, %d, want %q, 1", report, code, wantCheck)
	}

	report, code = Command([]string{"--check"})
	if report != wantCheck || code != 1 {
		t.Fatalf("--check on drifted tree = %q, %d, want %q, 1", report, code, wantCheck)
	}

	report, code = Command([]string{"--write"})
	if report != "" || code != 0 {
		t.Fatalf("--write = %q, %d, want \"\", 0", report, code)
	}

	report, code = Command([]string{"--check"})
	if report != "" || code != 0 {
		t.Fatalf("--check after --write = %q, %d, want \"\", 0", report, code)
	}

	report, code = Command([]string{"--help"})
	if code != 0 || report != grammar.Help+"\n" {
		t.Fatalf("--help = %q, %d, want %q, 0", report, code, grammar.Help+"\n")
	}

	report, code = Command([]string{"--bogus"})
	if code != 2 {
		t.Fatalf("--bogus code = %d, want 2 (report %q)", code, report)
	}
}

// run is the fixture git driver shared by this file's Command test.
func run(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
