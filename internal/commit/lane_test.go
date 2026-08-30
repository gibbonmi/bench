package commit

// The worktree commit's fast-lane route. Each test here drives the real command against
// a fixture repository whose phase manifest declares a lane, so the observed evidence is
// the verb's own stdout, exit code, and branch ref.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// laneScripts are the fixture's controllable checks. Each one stands in for a
// Bench-owned check of the kit's built-in lane, so a row can pick an outcome without
// paying for a run-binary build.
const (
	// passingCheck and failingCheck give a row direct control of the lane's outcome.
	passingCheck = "#!/bin/sh\nexit 0\n"
	failingCheck = "#!/bin/sh\necho 'the controllable check refused'\necho 'a later line'\nexit 1\n"
	// gofmtCheck reds on any misformatted Go file in the graded checkout.
	gofmtCheck = "#!/bin/sh\nout=$(gofmt -l .)\n[ -z \"$out\" ] || { echo \"gofmt: $out\"; exit 1; }\n"
	// proseCheck grades only the paths it is given, and reds on a line of 27 words or
	// more. It stands in for `bench gate-prose`, whose subject the commit supplies the
	// same way.
	proseCheck = "#!/bin/sh\nstatus=0\nfor f in \"$@\"; do\n  [ -f \"$f\" ] || continue\n  out=$(awk -v f=\"$f\" 'NF>=27 {print f \":\" NR \": sentence of \" NF \" words\"; exit}' \"$f\")\n  if [ -n \"$out\" ]; then echo \"$out\"; status=1; fi\ndone\nexit $status\n"
)

// longSentence is one line of 27 words: one more than the prose bound allows.
const longSentence = "one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twentyone twentytwo twentythree twentyfour twentyfive twentysix twentyseven\n"

// laneManifest is the fixture's declared lane. The checks run from the private checkout
// of the composed snapshot, so every argv is relative to it.
func laneManifest(t *testing.T) string {
	t.Helper()
	doc := map[string]any{
		"phases": []any{map[string]any{"name": "build", "argv": []string{"go", "build", "./..."}}},
		"lane": []any{
			map[string]any{"name": "check", "argv": []string{"sh", ".bench/check.sh"}},
			map[string]any{"name": "gofmt", "argv": []string{"sh", ".bench/gofmt.sh"}},
			map[string]any{"name": "prose", "argv": []string{"sh", ".bench/prose.sh", gate.LaneNamedMarkdownToken}},
			map[string]any{"name": "build", "argv": []string{"go", "build", "./..."}},
		},
	}
	encoded, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

// laneRepo is the commit landing fixture with a declared lane. checkExit picks the one
// controllable check's outcome; the other three checks grade the composed checkout.
func laneRepo(t *testing.T, checkExit int, write func(t *testing.T, root string)) (root, before string) {
	t.Helper()
	return landingRepo(t, 0, func(t *testing.T, root string) {
		bench := filepath.Join(root, ".bench")
		mustMkdirAll(t, bench)
		check := passingCheck
		if checkExit != 0 {
			check = failingCheck
		}
		mustWrite(t, filepath.Join(bench, "check.sh"), check, 0o755)
		mustWrite(t, filepath.Join(bench, "gofmt.sh"), gofmtCheck, 0o755)
		mustWrite(t, filepath.Join(bench, "prose.sh"), proseCheck, 0o755)
		mustWrite(t, filepath.Join(bench, "phases.json"), laneManifest(t), 0o644)
		mustWrite(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.24\n", 0o644)
		mustWrite(t, filepath.Join(root, "x.go"), "package fixture\n", 0o644)
		write(t, root)
	})
}

func noWrite(t *testing.T, root string) {}

// headTree is the tree the branch ref's commit points at.
func headTree(t *testing.T, root string) string {
	t.Helper()
	return strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD^{tree}")))
}

func head(t *testing.T, root string) string {
	t.Helper()
	return strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
}

// OG01, OG12, OG03: a declared lane grades the commit. The stdout states the lane's
// outcome and its checks, it names no gate phase, it never says green, and the branch
// ref moves to a commit whose tree is the composed snapshot.
func TestLanePassStatesItsOutcomeAndPublishesTheComposedSnapshot(t *testing.T) {
	root, before := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "lane{outcome=pass,checks=check,gofmt,prose,build}") {
		t.Errorf("stdout = %q, want the lane pass record naming its checks in order", stdout)
	}
	if strings.Contains(stdout, "phase ") {
		t.Errorf("stdout = %q, want no gate phase line: the lane replaced the gate", stdout)
	}
	if strings.Contains(stdout, "green") {
		t.Errorf("stdout = %q, want no green token: a lane pass is not a graded green", stdout)
	}
	if head(t, root) == before {
		t.Fatal("the branch ref did not move after a lane pass")
	}
	// OG03: the published tree is the composed snapshot, which carries the named path.
	if !headHasPrefix(t, root, "a.txt") {
		t.Fatalf("published tree does not carry the named path: %v", headPaths(t, root))
	}
	if tree := headTree(t, root); tree == "" {
		t.Fatal("the published commit has no tree")
	}
}

// OG02: the lane grades the composed snapshot, so an unnamed file that does not compile
// never reaches the lane's build check.
func TestLaneIgnoresAnUnnamedCompileError(t *testing.T) {
	root, _ := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)
	mustWrite(t, filepath.Join(root, "broken.go"), "package fixture\nfunc (\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "outcome=pass") {
		t.Fatalf("stdout = %q, want a lane pass: the unnamed file is not in the composed tree", stdout)
	}
}

// OG04, OG05: a check that exits nonzero refuses the commit, names the check with its
// first diagnostic line, and leaves the branch ref where it was.
func TestLaneFailRefusesNamingTheCheckAndItsFirstDiagnostic(t *testing.T) {
	root, before := laneRepo(t, 1, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "lane{outcome=fail,check=check}") {
		t.Errorf("stdout = %q, want the lane fail record naming the failing check", stdout)
	}
	if !strings.Contains(stdout, "the controllable check refused") {
		t.Errorf("stdout = %q, want the check's first diagnostic line", stdout)
	}
	if after := head(t, root); after != before {
		t.Fatalf("the branch ref moved from %s to %s on a lane fail", before, after)
	}
}

// OG06, OG07: the prose check grades the Markdown the commit names, and nothing else.
func TestLaneProseGradesOnlyTheNamedMarkdown(t *testing.T) {
	for _, tc := range []struct {
		name     string
		named    bool
		wantExit int
	}{
		{name: "named", named: true, wantExit: 1},
		{name: "unnamed", named: false, wantExit: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, _ := laneRepo(t, 0, noWrite)
			runGit(t, root, "reset", "-q", "--hard", "HEAD")
			mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)
			mustWrite(t, filepath.Join(root, "note.md"), "# Note\n\n"+longSentence, 0o644)

			args := []string{"-m", "m", "a.txt"}
			if tc.named {
				args = append(args, "note.md")
			}
			code, stdout, stderr := runCommand(t, root, args...)
			if code != tc.wantExit {
				t.Fatalf("exit = %d, want %d; stdout=%q stderr=%q", code, tc.wantExit, stdout, stderr)
			}
			if !tc.named {
				if !strings.Contains(stdout, "outcome=pass") {
					t.Fatalf("stdout = %q, want a lane pass: the file is not named", stdout)
				}
				return
			}
			if !strings.Contains(stdout, "lane{outcome=fail,check=prose}") {
				t.Errorf("stdout = %q, want the prose check named", stdout)
			}
			if !strings.Contains(stdout, "note.md:3:") {
				t.Errorf("stdout = %q, want the file and the line of the long sentence", stdout)
			}
		})
	}
}

// OG09: the commit rewrites the named Go files before the lane runs, so the lane's gofmt
// check grades the rewritten bytes and passes.
func TestLaneGofmtPassesOnAGoFileTheCommitRewrote(t *testing.T) {
	root, _ := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "y.go"), "package fixture\nfunc  Y( ) {\n}\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "y.go")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "formatted Go paths: y.go") {
		t.Errorf("stdout = %q, want the rewrite reported before the lane", stdout)
	}
	if !strings.Contains(stdout, "lane{outcome=pass") {
		t.Fatalf("stdout = %q, want the lane's gofmt check to pass on the rewritten file", stdout)
	}
}

// OG10, OG11: a dry run takes the same authority and stops before publication.
func TestLaneDryRunStatesTheOutcomeAndPublishesNothing(t *testing.T) {
	root, before := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	// PL23: a dry run takes the same lane the real run takes, check for check.
	if !strings.Contains(stdout, "lane{outcome=pass,checks=check,gofmt,prose,build}") {
		t.Errorf("stdout = %q, want the lane record line naming every declared check", stdout)
	}
	if strings.Contains(stdout, "classes=") {
		t.Errorf("stdout = %q, want no classes cell on a manifest lane", stdout)
	}
	if strings.Contains(stdout, "phase ") {
		t.Errorf("stdout = %q, want no gate phase line", stdout)
	}
	if strings.Contains(stdout, "green") {
		t.Errorf("stdout = %q, want no green token on a lane dry run", stdout)
	}
	if after := head(t, root); after != before {
		t.Fatalf("the branch ref moved from %s to %s on a dry run", before, after)
	}
}

// OG04, OG10, OG11: a dry run takes the lane's fail authority the same as a real run. A
// failing check exits 1, names the check, and leaves the branch ref unchanged.
func TestLaneDryRunFailRefusesNamingTheCheckAndLeavesTheRefUnchanged(t *testing.T) {
	root, before := laneRepo(t, 1, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "lane{outcome=fail,check=check}") {
		t.Errorf("stdout = %q, want the lane fail record naming the failing check", stdout)
	}
	if after := head(t, root); after != before {
		t.Fatalf("the branch ref moved from %s to %s on a dry-run lane fail", before, after)
	}
}

// OG23: a root that declares no lane keeps today's whole-project gate at the commit. The
// evidence is the tally the gate script writes and the published commit, which the gate
// authority publishes only on green. The fixture's gate is a script, and the script route
// prints no Bench-owned verdict line, so the row reads the tally instead.
func TestNoDeclaredLaneKeepsTheGateCommit(t *testing.T) {
	tally := ""
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {
		tally = filepath.Join(root, ".bench", "gate-tally")
		mustWrite(t, filepath.Join(root, ".bench", "gate.sh"),
			"#!/bin/sh\nprintf g >> "+sanitize.ShellQuote(tally)+"\nexit 0\n", 0o755)
	})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if strings.Contains(stdout, "lane{") {
		t.Errorf("stdout = %q, want no lane record: the root declares no lane", stdout)
	}
	if recorded, err := os.ReadFile(tally); err != nil || string(recorded) != "g" {
		t.Fatalf("gate tally = %q, %v; want one whole-project gate run", recorded, err)
	}
	if head(t, root) == before {
		t.Fatal("the branch ref did not move on a green gate")
	}
}

// OG24: a lane entry the loader cannot read refuses the commit and names the defect on
// stderr, before anything is graded or published.
func TestMalformedLaneEntryRefusesNamingTheDefect(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {
		mustMkdirAll(t, filepath.Join(root, ".bench"))
		mustWrite(t, filepath.Join(root, ".bench", "phases.json"),
			`{"phases":[{"name":"build","argv":["go","build","./..."]}],"lane":[{"name":"fmt","argv":[]}]}`, 0o644)
	})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "empty argv: fmt") {
		t.Fatalf("stderr = %q, want the loader's defect diagnostic", stderr)
	}
	if after := head(t, root); after != before {
		t.Fatalf("the branch ref moved from %s to %s on a malformed lane", before, after)
	}
}

// PL25: the lane grades the composed tree, so a commit that names the directory `docs`
// reaches the prose check with the Markdown under it. A named-path filter drops the
// directory and lets the long sentence through.
func TestLaneProseGradesAMarkdownFileUnderANamedDirectory(t *testing.T) {
	root, _ := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustMkdirAll(t, filepath.Join(root, "docs"))
	mustWrite(t, filepath.Join(root, "docs", "note.md"), "# Note\n\n"+longSentence, 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "docs")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "lane{outcome=fail,check=prose}") {
		t.Errorf("stdout = %q, want the prose check named", stdout)
	}
	if !strings.Contains(stdout, "docs/note.md:3:") {
		t.Errorf("stdout = %q, want the file and the line of the long sentence", stdout)
	}
}

// PL6: attribution refuses a special named path before the lane runs. The evidence that
// no check ran is the absent lane record: a lane that ran before attribution would grade
// the FIFO or block on it, and it would leave a record either way.
func TestLaneRefusesASpecialNamedPathBeforeAnyCheckRuns(t *testing.T) {
	root, before := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	fifo := filepath.Join(root, "fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		capability.Capability(t, capability.Fifo, "FIFOs unavailable: "+err.Error())
	}

	code, stdout, stderr := runCommand(t, root, "-m", "m", "fifo")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, `special file "fifo" is not attributable`) {
		t.Fatalf("stderr = %q, want the attribution refusal naming the special file", stderr)
	}
	record := filepath.Join(strings.TrimSpace(string(runGit(t, root, "rev-parse", "--absolute-git-dir"))), "bench-last-lane")
	if _, err := os.Stat(record); !os.IsNotExist(err) {
		t.Fatalf("lane record at %s: err = %v, want it absent because no check ran", record, err)
	}
	if after := head(t, root); after != before {
		t.Fatalf("the branch ref moved from %s to %s on an unattributable path", before, after)
	}
}

// TestManifestLaneRunsAsDeclared is PL21. A linked project's declared lane keeps the
// meaning its manifest gives it: every check runs, and the line names no class. A
// selection applied here would drop three of the four checks, because the manifest's
// names match the kit's.
func TestManifestLaneRunsAsDeclared(t *testing.T) {
	root, _ := laneRepo(t, 0, noWrite)
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "note.md"), "# Note\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "note.md")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "lane{outcome=pass,checks=check,gofmt,prose,build}") {
		t.Errorf("stdout = %q, want every declared check named", stdout)
	}
	if strings.Contains(stdout, "classes=") {
		t.Errorf("stdout = %q, want no classes cell on a manifest lane", stdout)
	}
}
