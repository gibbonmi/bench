package ownership_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// repin points the charge arguments at the current main and HEAD after a commit.
func repin(t *testing.T, args []string) []string {
	t.Helper()
	args[6], args[8] = preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
	return args
}

// chargeWithin runs the charge and fails the test when it does not answer in time, so a
// special file that the gather opens fails the case rather than hanging the suite.
func chargeWithin(t *testing.T, args []string) (string, int) {
	t.Helper()
	type answer struct {
		out  string
		code int
	}
	done := make(chan answer, 1)
	go func() {
		out, code := preflight.Command(args)
		done <- answer{out, code}
	}()
	select {
	case got := <-done:
		return got.out, got.code
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("charge blocked on its metadata source")
		return "", 0
	}
}

// TestSourceSafety covers each special spec source through the actual charge.
func TestSourceSafety(t *testing.T) {
	t.Run("OC71 symlink", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		outside := filepath.Join(t.TempDir(), "spec.md")
		data, err := os.ReadFile(specPath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(outside, append(data, "\n```bench-planning-facts\n"+completeEmpty+"\n```\n"...), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(specPath); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, specPath); err != nil {
			t.Fatal(err)
		}
		commit(t, "linked spec")
		assertChargeRefused(t, root, repin(t, args), "spec not readable")
	})
	t.Run("OC72 FIFO", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		if err := os.Remove(specPath); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(specPath, 0o600); err != nil {
			t.Fatal(err)
		}
		if out, code := chargeWithin(t, args); code == 0 || !strings.Contains(out, "spec") {
			t.Fatalf("charge over a FIFO spec = (%d), want a refusal naming the spec:\n%s", code, out)
		}
		preflighttest.AssertNothingPublished(t, root)
	})
	t.Run("OC73 directory", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		if err := os.Remove(specPath); err != nil {
			t.Fatal(err)
		}
		preflighttest.MustWriteFile(t, filepath.Join(specPath, "inner.md"), "# inner\n")
		commit(t, "directory spec")
		assertChargeRefused(t, root, repin(t, args), "spec not found")
	})
	t.Run("OC74 oversized", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		data, err := os.ReadFile(specPath)
		if err != nil {
			t.Fatal(err)
		}
		filler := "\n<!-- " + strings.Repeat("x", 2<<20) + " -->\n"
		preflighttest.MustWriteFile(t, specPath, string(data)+filler+"\n```bench-planning-facts\n"+completeEmpty+"\n```\n")
		commit(t, "oversized spec")
		assertChargeRefused(t, root, repin(t, args), "spec not readable")
	})
}

const (
	slug     = "example"
	specPath = "specs/" + slug + "/spec.md"
)

// completeEmpty is the declared-complete planning inventory with no item. Each grammar case
// edits one part of it, so the cases share one valid baseline.
var completeEmpty = inventorytest.CompleteEmpty

// fenced is one bench-planning-facts section that holds payload.
var fenced = inventorytest.Section

// commit records every change in the working tree as the new base on main, then adds one
// authorized change. A later collection reads the change at HEAD, and no fence needs to
// authorize it.
func commit(t *testing.T, message string) {
	t.Helper()
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", message)
	preflighttest.RunGit(t, "update-ref", "refs/heads/main", "HEAD")
	preflighttest.MustWriteFile(t, "internal/example/after.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/example/after.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "advance "+message)
}

// emptyTree is a linked repository with no anchor registry and no fixture inventory, whose
// ticket writes a guidance path and a source file. tail closes its spec and states its whole
// planning inventory, so an empty tail declares none.
func emptyTree(tail string) ownershiptest.Repo {
	return ownershiptest.Repo{
		Writes: []string{".agents/x/SKILL.md", "internal/example/pinned.go"},
		Files: map[string]string{
			".agents/x/SKILL.md":         "# X\n",
			"internal/example/pinned.go": "package example\n",
		},
		OwnInventory: true,
		Tail:         tail,
	}
}

// coveredTree is a conformant charge repository, so only a source or checkout refusal can
// stop it.
func coveredTree() ownershiptest.Repo {
	return ownershiptest.Repo{
		Writes: []string{"internal/example/pinned.go", "tests/canary/fam/pin-fixture"},
		Files: map[string]string{
			"internal/example/pinned.go":          "package example\n",
			"tests/canary/fam/pin-fixture/BASE":   "internal/example/pinned.go\n",
			"tests/canary/fam/pin-fixture/EXPECT": "red\n",
		},
	}
}

// TestSourcePrecedence covers the retained refusal order and the charge source binding.
func TestSourcePrecedence(t *testing.T) {
	t.Run("OC79 dirty checkout first", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		args[4] = "missing.md"
		preflighttest.MustWriteFile(t, "dirty.txt", "dirty\n")
		out, code := preflight.Command(args)
		if code == 0 || !strings.Contains(out, "checkout") || strings.Contains(out, "selected ticket") {
			t.Fatalf("dirty checkout with a missing ticket = (%d), want the checkout refusal first:\n%s", code, out)
		}
	})
	t.Run("OC90 evidence source mismatch", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		preflighttest.RunGit(t, "update-index", "--assume-unchanged", chargesource.BuildPhase)
		preflighttest.MustWriteFile(t, chargesource.BuildPhase, "# live drift\n")
		out, code := preflight.Command(args)
		if code == 0 || !strings.Contains(out, chargesource.BuildPhase+" does not match source tip") {
			t.Fatalf("charge over a drifted source = (%d), want the source refusal:\n%s", code, out)
		}
		preflighttest.AssertNothingPublished(t, root)
	})
}

// TestSourceMovement covers the retained movement and source-tip refusals.
func TestSourceMovement(t *testing.T) {
	t.Run("OC80 persistent movement", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			preflighttest.MustWriteFile(t, "moved.txt", strings.Repeat("x", calls)+"\n")
			preflighttest.RunGit(t, "add", "moved.txt")
			preflighttest.RunGit(t, "commit", "-q", "-m", "move")
		})
		t.Cleanup(restore)
		out, code := preflight.Command(args)
		restore()
		if calls != 2 || code == 0 || !strings.Contains(out, "snapshot drift") {
			t.Fatalf("charge under persistent movement = (%d) after %d callbacks, want the drift refusal after 2:\n%s", code, calls, out)
		}
		preflighttest.AssertNothingPublished(t, root)
	})
	t.Run("OC81 stale source tip", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		args[8] = preflighttest.RunGit(t, "rev-parse", "main")
		assertChargeRefused(t, root, args, "tip-current")
	})
}

// TestSubjectPaths covers the subject paths that each premise kind names, nested paths
// included, and the refusal of a kind that no collector reads.
func TestSubjectPaths(t *testing.T) {
	owner := func(dir string) string { return ownerJSON("internal/"+dir+"/a.go", dir+".A", "func") }
	at := func(path string) string { return `{"path":"` + path + `","line":1,"needle":"a","source":"tip"}` }
	probe := func(path string) string {
		return `{"path":"` + path + `","function":"TestA","record_spec":"r","source_commit":"s","evidence_commit":"e","chunk":"c","verification":"v","mutation":"m"}`
	}
	for _, row := range []struct{ kind, subject, want string }{
		{"lines", `{"path":"internal/lines/a.go","source":"tip"}`, "internal/lines/a.go"},
		{"glob-members", `{"pattern":"internal/glob/*.go","scope":"internal/glob/","source":"tip"}`, "internal/glob/"},
		{"headroom", `{"directory":"internal/room","files":{"create":[],"delete":[]}}`, "internal/room"},
		{"symbol", owner("symbol"), "internal/symbol/a.go"},
		{"bound", owner("bound"), "internal/bound/a.go"},
		{"import", `{"caller":"example.com/caller","target":"example.com/target","owner":` + owner("imported") + `}`, "internal/imported/a.go"},
		{"anchor-kind", `{"owner":` + at("internal/anchor/a.go") + `,"kind":"FORBID"}`, "internal/anchor/a.go"},
		{"sentinel", `{"owner":` + owner("sentinel") + `,"producers":[` + owner("producer") + `]}`, "internal/sentinel/a.go internal/producer/a.go"},
		{"error", `{"owner":` + owner("failure") + `,"producers":[` + owner("raiser") + `]}`, "internal/failure/a.go internal/raiser/a.go"},
		{"reused-test", probe("internal/reused/a_test.go"), "internal/reused/a_test.go"},
		{"reused-owner", `{"owner":` + at("internal/rule/a.go") + `}`, "internal/rule/a.go"},
		{"retirement", `{"owner":` + owner("retired") + `,"consumers":[` + owner("consumer") + `]}`, "internal/retired/a.go internal/consumer/a.go"},
		{"posture", `{"owner":` + owner("posture") + `,"callers":[` + owner("caller") + `],"writer_probe":` + probe("internal/writer/a_test.go") + `}`,
			"internal/posture/a.go internal/caller/a.go internal/writer/a_test.go"},
		{"printed-field", `{"field":` + owner("field") + `,"producers":[` + owner("writer") + `]}`, "internal/field/a.go internal/writer/a.go"},
		{"fence-exclusion", `{"path":"internal/excluded"}`, "refused"},
	} {
		paths, err := ownership.SubjectPaths(ownership.PremiseItem{Kind: row.kind, Subject: json.RawMessage(row.subject)})
		got := strings.Join(paths, " ")
		if err != nil {
			got = "refused"
		}
		if got != row.want {
			t.Errorf("%s subject paths = %q (%v), want %q", row.kind, got, err, row.want)
		}
	}
}
