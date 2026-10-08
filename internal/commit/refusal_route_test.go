// Refusal route tests for the commit: each commit refusal prints the route of its face in
// the shared refusal-route registry, and that route clears the refusal.
package commit

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/worktree"
)

// TestPublishedUnreconciledRouteIsTheResetPlan is RR32 through RR34, the collision 8 restore
// repro. A commit exit 3 names the reset plan at the published commit for the commit's own
// root, and never the restore that the destructive-git guard denies. A root that is not
// line-safe prints the checkout placeholder, so no control byte reaches the record.
func TestPublishedUnreconciledRouteIsTheResetPlan(t *testing.T) {
	// requireResetPlan requires the reset plan at the published commit for checkout, the
	// root as the route prints it, in the next= field of the exit 3 record.
	requireResetPlan := func(t *testing.T, stdout, published, checkout string) {
		t.Helper()
		_, fields, _ := recordFields(t, stdout)
		next := fields[refusalroute.NextField]
		if want := "bench worktree reset --to " + sanitize.ShellQuote(published) + " " + checkout; next != want {
			t.Errorf("next = %q, want %q", next, want)
		}
		if strings.Contains(next, "git restore") {
			t.Errorf("next = %q, want no restore step", next)
		}
	}
	t.Run("line-safe root", func(t *testing.T) {
		root, stdout, published := publishedUnreconciled(t, "named")
		requireResetPlan(t, stdout, published, sanitize.ShellQuote(root))
	})
	// A test's first TempDir call makes the one parent in TMPDIR that all its temporary
	// directories share. So this test names the unsafe directory, and the subtest sets
	// TMPDIR to it before its first call: the commit's root then sits under the ESC byte.
	unsafe := filepath.Join(t.TempDir(), "esc\x1bdir")
	t.Run("root that is not line-safe", func(t *testing.T) {
		mustMkdirAll(t, unsafe)
		t.Setenv("TMPDIR", unsafe)
		root, stdout, published := publishedUnreconciled(t, "named")
		if sanitize.LineSafe(root) {
			t.Fatalf("root = %q, want a root that is not line-safe", root)
		}
		requireResetPlan(t, stdout, published, "<checkout>")
	})
}

// commitSet is one commit route fixture: the primary checkout, the Bench home that holds
// its assignments, and the checkout the caller commits in.
type commitSet struct{ primary, home, checkout string }

// commitFaceFixture produces exactly one commit face. build makes the set and breaks it so
// the commit prints that face. The registry walk requires a fixture for each commit face,
// so a face added with no fixture turns TestCommitFacesFollowTheirRoutes red.
type commitFaceFixture struct {
	// cause tells apart two fixtures of one face. flags are the caller's own flags ahead
	// of commitArgs.
	face, cause string
	flags       []string
	build       func(t *testing.T) commitSet
	// exit is the exit the face prints at: 1 for a refusal, 3 for a published commit whose
	// checkout did not reconcile.
	exit int
	// contains and suffix state what the printed route must hold and end with.
	contains, suffix string
	// clear removes the fixture's own fault scaffolding once the face printed, which no
	// operator's tree holds, so the route repairs only the cause the face names.
	clear func(t *testing.T, f commitSet)
	// carry carries out a printed instruction by the fixture's own means, keyed by the
	// step's index in the face's route.
	carry map[int]func(t *testing.T, f commitSet)
	// slots are the replacement pairs for the operator slots of a printed command step: the
	// values the operator holds.
	slots []string
	// after carries out what the route's last step leaves to the operator, and it returns
	// the checkout the commit then reruns in.
	after func(t *testing.T, f commitSet, last string) commitSet
	// settled grades the rerun. With none, the rerun publishes.
	settled func(t *testing.T, f commitSet, code int, stdout, stderr string)
}

// commitArgs are the caller's own arguments in every route fixture.
var commitArgs = []string{"-m", "m", "--", "a.txt"}

// callerCommit is the caller's own commit at the worktree that label addresses, as the
// route that re-runs it prints it: label and words are in their printed form.
func callerCommit(label string, words ...string) string {
	return strings.Join(append([]string{"bench commit --in", label}, words...), " ")
}

// TestCommitFacesFollowTheirRoutes is RR36 through RR38. The registry is the source of the
// commit's face set, so each commit face needs a producing fixture. Each fixture's printed
// route is carried out step by step: an instruction by the fixture's own means, and each
// command step verbatim through the verb's own entry. The commit then reruns out of the
// face: as the route's own last step, or after a route that ends elsewhere.
func TestCommitFacesFollowTheirRoutes(t *testing.T) {
	var keys [][2]string
	for _, fixture := range commitFaceFixtures() {
		keys = append(keys, [2]string{fixture.face, fixture.cause})
	}
	faces := routetest.Fixtures(t, refusalroute.Commit, keys)
	for _, fixture := range commitFaceFixtures() {
		t.Run(fixture.face+"/"+fixture.cause, func(t *testing.T) { followCommitFace(t, faces[fixture.face], fixture) })
	}
}

// followCommitFace drives one fixture: the commit prints the face, the walk carries out the
// printed route, and the commit reruns out of the face.
func followCommitFace(t *testing.T, face refusalroute.Face, fixture commitFaceFixture) {
	f := fixture.build(t)
	args := slices.Concat(fixture.flags, commitArgs)
	code, stdout, stderr := runCommand(t, f.checkout, args...)
	next, printed := printedNext(stderr)
	if fixture.exit == 3 {
		_, fields, _ := recordFields(t, stdout)
		next, printed = fields[refusalroute.NextField], true
	}
	if code != fixture.exit || !printed || !strings.Contains(next, fixture.contains) || !strings.HasSuffix(next, fixture.suffix) {
		t.Fatalf("face %s = (%d, %q, %q), want exit %d and a next= route that holds %q and ends with %q", face.Name, code, stdout, stderr, fixture.exit, fixture.contains, fixture.suffix)
	}
	steps := routetest.Steps(t, face, next, "", nil)
	if fixture.clear != nil {
		fixture.clear(t, f)
	}
	carry := func(index int) (func(), bool) {
		step, carried := fixture.carry[index]
		return func() { step(t, f) }, carried
	}
	last := routetest.Follow(t, face, steps, carry, func(step string) string {
		if fixture.slots != nil {
			step = strings.NewReplacer(fixture.slots...).Replace(step)
		}
		return runRouteStep(t, f, step)
	})
	// A route that ends with the caller's own commit has rerun it, and the step exited 0.
	if strings.HasPrefix(steps[len(steps)-1], "bench commit ") {
		return
	}
	if fixture.after != nil {
		f = fixture.after(t, f, last)
	}
	code, stdout, stderr = runCommand(t, f.checkout, args...)
	if fixture.settled != nil {
		fixture.settled(t, f, code, stdout, stderr)
	} else if code != 0 {
		t.Fatalf("%s rerun = (%d, %q, %q), want exit 0", face.Name, code, stdout, stderr)
	}
}

// printedNext is the route of a refusal's next= line on stderr.
func printedNext(stderr string) (string, bool) {
	for _, line := range strings.Split(stderr, "\n") {
		if next, ok := strings.CutPrefix(line, refusalroute.NextField+"="); ok {
			return next, true
		}
	}
	return "", false
}

// runRouteStep runs one printed command step in process, the way the CLI dispatches it: a
// commit at the worktree that its tree target names, and any other verb at the caller's
// checkout. The shared walk states the rules a printed step and its exit meet. It returns
// the step's stdout.
func runRouteStep(t *testing.T, f commitSet, step string) string {
	t.Helper()
	words := routetest.Words(t, step)
	var stdout, stderr bytes.Buffer
	code := 0
	switch strings.Join(words[:min(3, len(words))], " ") {
	case "bench commit --in":
		dir, err := worktree.TreeTarget(f.checkout, words[3])
		if err != nil {
			t.Fatalf("printed step %q names no worktree: %v", step, err)
		}
		var out, errOut string
		code, out, errOut = runCommand(t, dir, words[4:]...)
		stdout.WriteString(out)
		stderr.WriteString(errOut)
	case "bench worktree create":
		code = worktree.CreateCommand(f.checkout, f.home, words[3:], &stdout, &stderr)
	case "bench worktree reset":
		code = worktree.ResetCommand(f.checkout, f.home, words[3:], &stdout, &stderr)
	case "bench doctor":
		code = adopt.Doctor(words[2:], io.Discard, io.Discard, "fixture")
	default:
		t.Fatalf("printed step %q runs no verb that a commit route names", step)
	}
	if code != 0 && !routetest.Diagnostic(words) {
		t.Fatalf("printed step %q = (%d, %q, %q), want exit 0", step, code, stdout.String(), stderr.String())
	}
	return stdout.String()
}

// commitFixtureSpec is the delivery the fixture repository's policy approves.
const commitFixtureSpec = "specs/commit-fixture/spec.md"

// primaryCommitSet is the commit fixture whose caller commits in the primary checkout, with
// its change to a.txt uncommitted.
func primaryCommitSet(t *testing.T) commitSet {
	t.Helper()
	notTheKitRoot(t)
	f := commitSet{primary: t.TempDir(), home: t.TempDir()}
	t.Setenv("BENCH_HOME", f.home)
	initializeLandingRepo(t, f.primary, 0)
	f.checkout = f.primary
	mustWrite(t, filepath.Join(f.primary, "a.txt"), "changed\n", 0o644)
	return f
}

// assignedCommitSet is the commit fixture whose checkout is an assignment that the create
// verb made and the delivery admits, so a printed tree target resolves as the CLI resolves
// it. gate writes the checkout's gate, which its base commit holds. The caller's change to
// a.txt waits uncommitted.
func assignedCommitSet(t *testing.T, label string, gate func(t *testing.T, root string)) commitSet {
	t.Helper()
	f := primaryCommitSet(t)
	runRouteStep(t, f, "bench worktree create --request "+label+" --label "+label)
	f.checkout = admitCreated(t, f, label)
	prepareLandingCheckout(t, f.checkout, gate)
	runGit(t, f.checkout, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(f.checkout, "a.txt"), "changed\n", 0o644)
	return f
}

// admitCreated admits the assignment that a create of label made, as the operator's start
// of the delivery would, and returns its checkout.
func admitCreated(t *testing.T, f commitSet, label string) string {
	t.Helper()
	checkout, err := worktree.TreeTarget(f.primary, label)
	if err != nil {
		t.Fatalf("created assignment %q: %v", label, err)
	}
	commitmenttest.Admit(t, checkout, label, commitFixtureSpec)
	return checkout
}

// gateScript writes the checkout's gate as body.
func gateScript(body string) func(t *testing.T, root string) {
	return func(t *testing.T, root string) {
		t.Helper()
		mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), "#!/bin/sh\n"+body+"\n", 0o755)
	}
}

// adminPath is the path of name in the checkout's own Git directory.
func adminPath(t *testing.T, checkout, name string) string {
	t.Helper()
	dir, err := git.AdminDir(checkout)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, name)
}

func commitFaceFixtures() []commitFaceFixture {
	quoted := sanitize.ShellQuote
	// rerun is the fixtures' own commit at the worktree of label, with flags.
	rerun := func(label string, flags ...string) string {
		return callerCommit(quoted(label), slices.Concat(flags, []string{"-m", quoted(commitArgs[1]), "--", quoted(commitArgs[3])})...)
	}
	// red is a commit-red fixture: grade runs a check that reds on a.txt until the walk
	// carries out the repair, and the caller passes flags.
	red := func(cause string, flags []string, grade func(check string) func(*testing.T, string)) commitFaceFixture {
		return commitFaceFixture{
			face: faceRed, cause: cause, flags: flags, exit: 1,
			build: func(t *testing.T) commitSet {
				f := assignedCommitSet(t, faceRed, grade("! grep -qs red a.txt"))
				mustWrite(t, filepath.Join(f.checkout, "a.txt"), "red\n", 0o644)
				return f
			},
			suffix:   rerun(faceRed, flags...),
			contains: "repair each failure",
			carry: map[int]func(*testing.T, commitSet){0: func(t *testing.T, f commitSet) {
				mustWrite(t, filepath.Join(f.checkout, "a.txt"), "fixed\n", 0o644)
			}},
		}
	}
	// laneCheck declares check as the checkout's one lane check, behind a gate that passes.
	laneCheck := func(check string) func(*testing.T, string) {
		return func(t *testing.T, root string) {
			gateScript("exit 0")(t, root)
			mustWrite(t, filepath.Join(root, ".bench", "phases.json"), `{"lane":[{"name":"check","argv":["sh","-c",`+strconv.Quote(check)+`]}]}`, 0o644)
		}
	}
	return []commitFaceFixture{
		{
			// The gate takes the checkout's index lock, so the commit publishes and the
			// checkout cannot follow it.
			face: facePublishedUnreconciled,
			exit: 3,
			build: func(t *testing.T) commitSet {
				return assignedCommitSet(t, facePublishedUnreconciled, func(t *testing.T, root string) {
					gateScript(": > "+quoted(adminPath(t, root, "index.lock")))(t, root)
				})
			},
			contains: "bench worktree reset --to ",
			// The gate's stale lock is the fixture's fault, and the reset needs the index
			// lock too.
			clear: func(t *testing.T, f commitSet) { mustRemove(t, adminPath(t, f.checkout, "index.lock")) },
			// The reset plan prints its own apply, which the operator runs next.
			after: func(t *testing.T, f commitSet, plan string) commitSet {
				for _, field := range strings.Split(plan, ",") {
					if apply, ok := strings.CutPrefix(field, refusalroute.NextField+"="); ok {
						runRouteStep(t, f, apply)
						return f
					}
				}
				t.Fatalf("reset plan = %q, want a next= apply", plan)
				return f
			},
			// The published commit holds the change, so the checkout is clean at it and
			// the rerun has nothing left to publish.
			settled: func(t *testing.T, f commitSet, code int, stdout, stderr string) {
				if status := string(runGit(t, f.checkout, "status", "--porcelain")); status != "" || code == 3 {
					t.Fatalf("after the route the checkout status is %q and the rerun = (%d, %q, %q), want a clean checkout and no exit 3", status, code, stdout, stderr)
				}
			},
		},
		{
			face:     facePrimaryCheckout,
			exit:     1,
			build:    primaryCommitSet,
			contains: "bench worktree create --request ",
			// The operator names its own request and label in the printed create, then
			// starts the delivery there and carries its change over.
			slots: []string{"<opaque-id>", facePrimaryCheckout, "<work-item>", facePrimaryCheckout},
			after: func(t *testing.T, f commitSet, _ string) commitSet {
				f.checkout = admitCreated(t, f, facePrimaryCheckout)
				mustWrite(t, filepath.Join(f.checkout, "a.txt"), "changed\n", 0o644)
				return f
			},
		},
		red("", nil, gateScript),
		red("lane", nil, laneCheck),
		red("dry run", []string{"--dry-run"}, gateScript),
		{
			// The gate cannot open its lock until the fixture frees it.
			face:   faceInfrastructure,
			exit:   1,
			suffix: rerun(faceInfrastructure),
			build: func(t *testing.T) commitSet {
				f := assignedCommitSet(t, faceInfrastructure, gateScript("exit 0"))
				mustMkdirAll(t, adminPath(t, f.checkout, "bench-gate.lock"))
				return f
			},
			clear: func(t *testing.T, f commitSet) { mustRemove(t, adminPath(t, f.checkout, "bench-gate.lock")) },
		},
		{
			// The caller names a.txt before it writes the file, so the named path is absent.
			face:   faceNamedPath,
			exit:   1,
			suffix: rerun(faceNamedPath),
			build: func(t *testing.T) commitSet {
				f := assignedCommitSet(t, faceNamedPath, gateScript("exit 0"))
				mustRemove(t, filepath.Join(f.checkout, "a.txt"))
				return f
			},
			carry: map[int]func(*testing.T, commitSet){0: func(t *testing.T, f commitSet) {
				mustWrite(t, filepath.Join(f.checkout, "a.txt"), "changed\n", 0o644)
			}},
		},
		tipMovedFixture(rerun(faceTipMoved)),
		{
			// A lane declaration that the loader cannot read has no repair the commit can
			// name, so it hands back. The reviewer withdraws the declaration.
			face:   faceHandback,
			exit:   1,
			suffix: rerun(faceHandback),
			build: func(t *testing.T) commitSet {
				f := assignedCommitSet(t, faceHandback, gateScript("exit 0"))
				mustWrite(t, filepath.Join(f.checkout, ".bench", "phases.json"), `{"lane":[{"name":"fmt","argv":[]}]}`, 0o644)
				return f
			},
			carry: map[int]func(*testing.T, commitSet){0: func(t *testing.T, f commitSet) {
				mustRemove(t, filepath.Join(f.checkout, ".bench", "phases.json"))
			}},
		},
	}
}

func mustRemove(t *testing.T, path string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
}
