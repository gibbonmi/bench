package runtime

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/handoff"
	"github.com/gibbonmi/bench/internal/status"
)

// handoffFixtureOnMain pins HEAD to main. The unpushed-commit count resolves the default
// branch, so a fixture that left the initial branch to the local git default would report a
// different fact on a machine whose init.defaultBranch differs.
func handoffFixtureOnMain(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	return f
}

// handoffShortSHA is the seven-byte HEAD prefix the pin block renders. It is read off the
// fixture rather than restated, so the assertion cannot pass against a constant.
func handoffShortSHA(t *testing.T, f contract.Fixture) string {
	t.Helper()
	return handoffShort(t, "HEAD", strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout))
}

// handoffShort takes the seven-byte prefix of a hash the subject reported. Every such
// slice goes through here because a subject that answers nothing yields an empty string,
// and slicing that panics — which fails the whole test binary rather than one test, so
// every later contract in this package stops reporting. The canary reads those reports:
// a fixture proving some other check still bites goes quiet, and the gate loses a
// tripwire to what looks like an unrelated failure.
func handoffShort(t *testing.T, what, value string) string {
	t.Helper()
	if len(value) < 7 {
		t.Fatalf("%s: want a hash from the subject, got %q", what, value)
	}
	return value[:7]
}

func handoffGitRoot(t *testing.T, f contract.Fixture) string {
	t.Helper()
	return strings.TrimSpace(f.Git("rev-parse", "--show-toplevel").Stdout)
}

func TestHandoffWritesAndPrints(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")

	out := f.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "# Session handoff")
	contract.RequireContains(t, out.Stdout, "## State")
	contract.RequireContains(t, out.Stdout, "## Next command")
	if !f.Exists(status.HandoffFile) {
		t.Fatal("handoff exited zero without writing " + status.HandoffFile)
	}
	// stdout is the pin block; the file is the pin block plus the Shape section. One
	// derivation feeding both sinks means the file must start with exactly what printed.
	written := f.ReadFile(status.HandoffFile)
	contract.RequireContains(t, written, strings.TrimRight(out.Stdout, "\n"))
	contract.RequireContains(t, written, "## Shape")
	contract.RequireNotContains(t, out.Stdout, "## Shape")
}

func TestHandoffNamesIdentity(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	for _, origin := range []string{"https://example.test/one.git", "git@example.test:two/three.git"} {
		f := handoffFixtureOnMain(t)
		f.WriteFile("tracked.txt", "base\n")
		f.CommitAll("base")
		f.Git("remote", "add", "origin", origin)

		out := f.Bench("handoff")
		out.RequireExit(0)
		want := fmt.Sprintf("Repository: `%s` (origin `%s`)", filepath.Base(handoffGitRoot(t, f)), origin)
		contract.RequireContains(t, out.Stdout, want)
	}
}

func TestHandoffEmitsDerivedPath(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	root := handoffGitRoot(t, f)
	if strings.Contains(root, "workspace") {
		t.Fatalf("fixture root %q sits under a workspace directory; the row proves no layout constant", root)
	}

	// The fixture's isolated HOME lives inside the repo, so the root is outside $HOME and
	// the path renders absolute.
	outside := f.Bench("handoff")
	outside.RequireExit(0)
	contract.RequireContains(t, outside.Stdout, "Path: `"+root+"`")

	// The same root under a $HOME that contains it renders abbreviated, and the absolute
	// form is gone — which a command printing os.Getwd() or a constant cannot do.
	inside := f.BenchEnv(map[string]string{"HOME": filepath.Dir(root)}, "handoff")
	inside.RequireExit(0)
	contract.RequireContains(t, inside.Stdout, "Path: `~/"+filepath.Base(root)+"`")
	contract.RequireNotContains(t, inside.Stdout, "Path: `"+root+"`")
}

func TestHandoffCarriesTreeFacts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	// Dirty with an unpushed commit: an upstream is configured and then outrun.
	dirty := handoffFixtureOnMain(t)
	dirty.WriteFile("tracked.txt", "base\n")
	dirty.CommitAll("base")
	bare := filepath.Join(t.TempDir(), "origin.git")
	dirty.Run("git", "init", "-q", "--bare", bare).RequireExit(0)
	dirty.Git("remote", "add", "origin", bare)
	dirty.Git("push", "-q", "-u", "origin", "main")
	dirty.WriteFile("tracked.txt", "second\n")
	dirty.CommitAll("second")
	dirty.WriteFile("scratch.txt", "uncommitted\n")

	out := dirty.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Branch: `main` — HEAD `"+handoffShortSHA(t, dirty)+"`")
	contract.RequireContains(t, out.Stdout, "1 dirty path,")
	contract.RequireContains(t, out.Stdout, "1 unpushed commit")

	// Clean with none, on a differently named branch: the hardcoded twin of either value
	// above fails here.
	clean := handoffFixtureOnMain(t)
	clean.WriteFile("tracked.txt", "base\n")
	clean.CommitAll("base")
	clean.Git("checkout", "-q", "-b", "topic")

	out = clean.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Branch: `topic` — HEAD `"+handoffShortSHA(t, clean)+"`")
	contract.RequireContains(t, out.Stdout, "clean tree,")
	contract.RequireContains(t, out.Stdout, "0 unpushed commits")
	contract.RequireNotContains(t, out.Stdout, "dirty path")
}

func TestHandoffNamesStagedSpec(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	staged := handoffFixtureOnMain(t)
	staged.WriteFile("specs/alpha/spec.md", "# Alpha\n\nStatus: staged\n")
	staged.CommitAll("staged spec")
	out := staged.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Spec: `specs/alpha/spec.md` (Status: staged)")

	flat := handoffFixtureOnMain(t)
	flat.WriteFile("specs/transitional.md", "# Transitional\n\nStatus: staged\n")
	flat.CommitAll("transitional flat spec")
	out = flat.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Spec: none staged.")
	contract.RequireNotContains(t, out.Stdout, "specs/transitional.md")

	// A second, different Status value: a constant `Status: staged` fails here. The
	// status-less spec alongside it is malformed rather than staged, so the field passes
	// over it instead of naming it with a status the file does not carry.
	drafted := handoffFixtureOnMain(t)
	drafted.WriteFile("specs/beta/spec.md", "# Beta\n\nStatus: drafting\n")
	drafted.WriteFile("specs/gamma/spec.md", "# Gamma\n\nno status line at all\n")
	drafted.CommitAll("drafted spec")
	out = drafted.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Spec: `specs/beta/spec.md` (Status: drafting)")
	contract.RequireNotContains(t, out.Stdout, "specs/gamma/spec.md")

	// An absent specs/ directory states absence rather than failing — the third fixture
	// row 6 names, which no case exercised while a specs/ holding an implemented spec
	// stood in for it.
	bare := handoffFixtureOnMain(t)
	bare.WriteFile("tracked.txt", "base\n")
	bare.CommitAll("no specs directory")
	out = bare.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Spec: none staged.")

	// No live spec: an implemented one is finished work, so the field states absence
	// rather than naming it.
	none := handoffFixtureOnMain(t)
	none.WriteFile("specs/done/spec.md", "# Done\n\nStatus: implemented\n")
	none.CommitAll("no staged spec")
	out = none.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Spec: none staged.")
	contract.RequireNotContains(t, out.Stdout, "specs/done/spec.md")
}

func TestHandoffDoesNotBlockOnSpecialSpec(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	f := handoffFixtureOnMain(t)
	f.WriteFile("README.md", "# fixture\n")
	f.CommitAll("init")
	f.WriteFifo("specs/hang/spec.md")

	out := f.BenchDeadlined("handoff")
	if out.TimedOut {
		t.Fatal("bench handoff blocked on a FIFO at specs/hang/spec.md, so its live-spec facts opened the path before classifying it")
	}
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Spec: none staged.")
}

func TestHandoffGateFieldIsStaleAware(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	// Green and current: a real gate run, so the cached tree is the work tree.
	current := handoffFixtureOnMain(t)
	current.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	current.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	current.CommitAll("gated base")
	current.Bench("gate").RequireExit(0)
	tree := handoffShort(t, "tree-hash", strings.TrimSpace(current.Bench("tree-hash").Stdout))
	out := current.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Gate: green at `"+tree+"` — current")

	// Green but stale: the cached verdict describes a tree that has moved. A bare
	// `Gate: green` reads this as a statement about the inherited tree.
	stale := handoffFixtureOnMain(t)
	stale.WriteFile("tracked.txt", "base\n")
	stale.CommitAll("base")
	cached := strings.Repeat("d", 40)
	work := strings.TrimSpace(stale.Bench("tree-hash").Stdout)
	writeGateCache(t, stale, strings.Replace(gateCacheLine(t, stale, "green"), work, cached, 1))
	out = stale.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Gate: green at `"+cached[:7]+"` — stale, work tree `"+handoffShort(t, "tree-hash", work)+"`")

	// Absent: no cache at all is its own answer, not a verdict.
	absent := handoffFixtureOnMain(t)
	absent.WriteFile("tracked.txt", "base\n")
	absent.CommitAll("base")
	out = absent.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Gate: no gate has run.")

	// A record that never parsed: the state is known but no cached tree survives behind
	// it, so the field names the work tree the reader is actually on. An implementation
	// that reaches for the cached tree regardless renders an empty inline-code span,
	// which reads as a deliberate blank rather than as a fact that does not exist.
	corrupt := handoffFixtureOnMain(t)
	corrupt.WriteFile("tracked.txt", "base\n")
	corrupt.CommitAll("base")
	writeGateCache(t, corrupt, "{not a verdict record}\n")
	work = strings.TrimSpace(corrupt.Bench("tree-hash").Stdout)
	out = corrupt.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Gate: invalid — no cached tree survives, work tree `"+handoffShort(t, "tree-hash", work)+"`")
	contract.RequireNotContains(t, out.Stdout, "``")
}

// handoffBoardRow is one `bench status --all` row reduced to the two fields the next-command
// field derives from: the signal's name and the action it carries.
type handoffBoardRow struct{ name, action string }

// handoffBoard parses every row `bench status --all` printed, in the board's own order.
// `--all` rather than the default board because the five-row budget can drop the very row
// the field selects, and the lead line — which repeats one of the rows — is skipped so a
// signal is never counted twice.
func handoffBoard(t *testing.T, f contract.Fixture) []handoffBoardRow {
	t.Helper()
	board := f.Bench("status", "--all")
	board.RequireExit(0)
	var rows []handoffBoardRow
	for _, line := range strings.Split(strings.TrimRight(board.Stdout, "\n"), "\n") {
		body, indented := strings.CutPrefix(line, "  ")
		if !indented {
			continue
		}
		detail, action, ok := strings.Cut(body, " → ")
		if !ok {
			t.Fatalf("board row unparseable: %q", line)
		}
		rows = append(rows, handoffBoardRow{strings.Fields(detail)[0], action})
	}
	return rows
}

// handoffBoardNext is the expectation, derived from the printed board rather than from the
// subject: the first row carrying a command a session could type. What qualifies comes from
// handoff.IsInvocable, the rule's one source — a copy spelled out here would keep asserting
// whatever it was written against, which is how this expectation came to demand the
// compound actions the rule had already stopped accepting.
//
// The independence that matters is in the *walk*, not the predicate: this takes the first
// qualifying row of the board as printed, so an implementation that re-ranked the board or
// took the leading row whatever its action said still disagrees with it.
func handoffBoardNext(rows []handoffBoardRow) (handoffBoardRow, bool) {
	for _, r := range rows {
		if handoff.IsInvocable(r.action) {
			return r, true
		}
	}
	return handoffBoardRow{}, false
}

func TestHandoffNextMatchesStatus(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	// Two fixtures whose first invocable board action differs, each with a higher-severity
	// prose row ahead of it, so neither a constant nor an unconditional `signals[0]` passes.
	drain := handoffFixtureOnMain(t)
	drain.WriteFile("capture/IDEAS.md", "- 2026-07-05  parked idea\n")
	drain.CommitAll("parked")
	drain.WriteFile("scratch.txt", "uncommitted\n")

	specs := handoffFixtureOnMain(t)
	specs.WriteFile("specs/done/spec.md", "# Done\n\nStatus: implemented\n")
	specs.CommitAll("merged spec awaiting retirement")

	found := map[string]handoffBoardRow{}
	for name, f := range map[string]contract.Fixture{"drain": drain, "specs": specs} {
		// The board is read first: `bench handoff` writes an untracked file, which is
		// itself a git signal, so a board read afterwards would answer a different tree.
		rows := handoffBoard(t, f)
		want, ok := handoffBoardNext(rows)
		if !ok {
			t.Fatalf("%s: no invocable action on the board %v; the fixture no longer exercises the row", name, rows)
		}
		found[name] = want

		out := f.Bench("handoff")
		out.RequireExit(0)
		contract.RequireContains(t, out.Stdout, "`"+want.action+"` — the board's leading invocable signal (`"+want.name+"`).")
		if rows[0] != want {
			// The prose row the selection stepped over must not have been rendered as
			// though it were a command.
			contract.RequireNotContains(t, out.Stdout, "`"+rows[0].action+"`")
		}
	}
	if found["drain"].action == found["specs"].action {
		t.Fatalf("both fixtures selected %q; the row needs two differing actions", found["drain"].action)
	}

	// A board whose only signal carries prose says so and points at the override. Reading
	// this as a clean board would tell a session with work waiting that nothing is pending.
	prose := handoffFixtureOnMain(t)
	prose.WriteFile("tracked.txt", "base\n")
	prose.CommitAll("base")
	prose.WriteFile("scratch.txt", "uncommitted\n")
	rows := handoffBoard(t, prose)
	if len(rows) == 0 {
		t.Fatal("prose fixture: the board is empty, so it exercises the clean case instead")
	}
	if want, ok := handoffBoardNext(rows); ok {
		t.Fatalf("prose fixture: the board offers %q; it no longer exercises the all-prose case", want.action)
	}
	out := prose.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "No invocable command derives from the board; name the next step with `--next`.")
	contract.RequireNotContains(t, out.Stdout, "Nothing pending")
	contract.RequireNotContains(t, out.Stdout, "`"+rows[0].action+"`")
}

func TestHandoffTrackedCaptureIsNeutral(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	f := handoffFixtureOnMain(t)
	f.Bench("handoff").RequireExit(0)
	f.CommitAll("track handoff")

	// The first rewrite makes the tracked handoff dirty. The next run must describe the
	// inherited checkout, not promote its own capture write into work to commit.
	f.Bench("handoff").RequireExit(0)
	out := f.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "clean tree")
	contract.RequireContains(t, out.Stdout, "Nothing pending — the board is clean.")
	contract.RequireNotContains(t, out.Stdout, "commit on green")

	f.WriteFile("scratch.txt", "uncommitted\n")
	signals := status.SignalsWith(f.Root, status.Query{ExcludeDirtyPaths: []string{status.HandoffFile}})
	var gitSignal *status.Signal
	for i := range signals {
		if signals[i].Name == "git" {
			gitSignal = &signals[i]
			break
		}
	}
	if gitSignal == nil {
		t.Fatalf("tracked handoff plus scratch path suppressed the git row: %v", signals)
	}
	if gitSignal.Detail != "1 dirty path" || gitSignal.Action != "commit on green" {
		t.Fatalf("git signal = %#v, want one dirty path with commit-on-green action", *gitSignal)
	}

	out = f.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Branch: `main` — HEAD `"+handoffShortSHA(t, f)+"`, 1 dirty path, 0 unpushed commits")
}
