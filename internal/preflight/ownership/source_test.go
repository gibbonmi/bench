package ownership_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// refused fails the test unless out is a refusal that names want and publishes no table.
func refused(t *testing.T, out string, code int, want string) {
	t.Helper()
	if code == 0 || !strings.Contains(out, want) || strings.Contains(out, "writes_proposal[") || strings.Contains(out, "closure_evidence[") {
		t.Fatalf("proposal = (%d), want a refusal naming %q and no proposal table:\n%s", code, want, out)
	}
}

// table decodes one table of a proposal response through the shared TOON decoder.
func table(t *testing.T, out, name string) []map[string]any {
	t.Helper()
	return preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), name)
}

func cell(row map[string]any, key string) string { return fmt.Sprint(row[key]) }

// matching is every row of one table whose named cells equal want.
func matching(t *testing.T, out, name string, want map[string]string) []map[string]any {
	t.Helper()
	var found []map[string]any
	for _, row := range table(t, out, name) {
		equal := true
		for key, value := range want {
			if cell(row, key) != value {
				equal = false
			}
		}
		if equal {
			found = append(found, row)
		}
	}
	return found
}

// one fails the test unless exactly one row of the table matches want, and answers it.
func one(t *testing.T, out, name string, want map[string]string) map[string]any {
	t.Helper()
	found := matching(t, out, name, want)
	if len(found) != 1 {
		t.Fatalf("%s rows matching %v = %d, want 1:\n%s", name, want, len(found), out)
	}
	return found[0]
}

// repin points the proposal arguments at the current main and HEAD after a commit.
func repin(t *testing.T, args []string) []string {
	t.Helper()
	args[6], args[8] = preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
	return args
}

// proposeWithin runs the proposal and fails the test when it does not answer in time, so
// a special file that the collector opens fails the case rather than hanging the suite.
func proposeWithin(t *testing.T, args []string) (string, int) {
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
		t.Fatal("proposal blocked on its metadata source")
		return "", 0
	}
}

// TestSourceSafety covers each special metadata source and each unsafe ownership path.
func TestSourceSafety(t *testing.T) {
	t.Run("OC71 symlink", func(t *testing.T) {
		args := proposalArgs(t, closureTree().Seed(t))
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
		out, code := preflight.Command(repin(t, args))
		refused(t, out, code, "spec not readable")
	})
	t.Run("OC72 FIFO", func(t *testing.T) {
		root := closureTree().Seed(t)
		args := proposalArgs(t, root)
		if err := os.Remove(specPath); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(specPath, 0o600); err != nil {
			t.Fatal(err)
		}
		out, code := proposeWithin(t, args)
		refused(t, out, code, "spec")
	})
	t.Run("OC73 directory", func(t *testing.T) {
		args := proposalArgs(t, closureTree().Seed(t))
		if err := os.Remove(specPath); err != nil {
			t.Fatal(err)
		}
		preflighttest.MustWriteFile(t, filepath.Join(specPath, "inner.md"), "# inner\n")
		commit(t, "directory spec")
		out, code := preflight.Command(repin(t, args))
		refused(t, out, code, "spec not found")
	})
	t.Run("OC74 oversized", func(t *testing.T) {
		args := proposalArgs(t, closureTree().Seed(t))
		data, err := os.ReadFile(specPath)
		if err != nil {
			t.Fatal(err)
		}
		filler := "\n<!-- " + strings.Repeat("x", 2<<20) + " -->\n"
		preflighttest.MustWriteFile(t, specPath, string(data)+filler+"\n```bench-planning-facts\n"+completeEmpty+"\n```\n")
		commit(t, "oversized spec")
		out, code := preflight.Command(repin(t, args))
		refused(t, out, code, "spec not readable")
	})
	t.Run("OC75 control byte", func(t *testing.T) {
		tr := closureTree()
		tr.Files["tests/canary/fam/bad\x1bfixture/BASE"] = "internal/example/pinned.go\n"
		tr.Files["tests/canary/fam/bad\x1bfixture/EXPECT"] = "red\n"
		out, code := propose(t, tr.Seed(t))
		refused(t, out, code, "represent")
	})
	t.Run("OC76 escaping path", func(t *testing.T) {
		for _, entry := range []string{"internal/example/../../../outside.go (new)", "internal/example/../../outside.go (new)"} {
			tr := closureTree()
			tr.Writes = append(tr.Writes, entry)
			tr.Fence = []string{"- `internal/example/pinned.go`", "- `internal/toon/toon.go`", "- `internal/example/cyclic.go`", "- `internal/example/included.go`"}
			out, code := propose(t, tr.Seed(t))
			refused(t, out, code, "escapes the repository")
		}
	})
	t.Run("OC77 segment prefix", func(t *testing.T) {
		tr := closureTree()
		tr.Fence = []string{"- `internal/example/pinned.go`", "- `internal/toon/toon.go`", "- `internal/example/cyclic.go`", "- `internal/example/included.go`", "- `internal/conformanc`"}
		out := mustPropose(t, tr.Seed(t))
		one(t, out, "writes_proposal", map[string]string{"path": "internal/conformance/data_handling_test.go", "fence": "spec fence expansion required"})
	})
}

// manyPins is a tree whose one written path thirty fixtures pin, so its closure evidence
// passes the response bound.
func manyPins() ownershiptest.ProposalRepo {
	tr := ownershiptest.ProposalRepo{
		Writes: []string{"internal/example/pinned.go"},
		Files:  map[string]string{"internal/example/pinned.go": "package example\n"},
	}
	for i := range 30 {
		dir := "tests/canary/fam/pin-" + string(rune('a'+i/26)) + string(rune('a'+i%26))
		tr.Files[dir+"/BASE"] = "internal/example/pinned.go\n"
		tr.Files[dir+"/EXPECT"] = "red\n"
	}
	return tr
}

// TestProposalBound grades the bounded route: the response owner spills the long proposal, and the
// complete spill file still holds every evidence reason.
func TestProposalBound(t *testing.T) {
	root := manyPins().Seed(t)
	out := mustPropose(t, root)
	var stdout, stderr bytes.Buffer
	owner := responsebound.New(t.TempDir(), &stdout, &stderr, func() string { return root }, false)
	if _, err := owner.Stdout().Write([]byte(out)); err != nil {
		t.Fatal(err)
	}
	owner.Finish()
	spilled, err := os.ReadFile(responseboundtest.Path(t, stdout.String()))
	if err != nil {
		t.Fatal(err)
	}
	rows := table(t, string(spilled), "closure_evidence")
	if len(rows) != 30 {
		t.Fatalf("spilled evidence rows = %d, want 30:\n%s", len(rows), spilled)
	}
	for _, row := range rows {
		if cell(row, "source") == "" || !strings.Contains(cell(row, "reason"), "internal/example/pinned.go") {
			t.Errorf("spilled evidence row lost its reason: %v", row)
		}
	}
}

// coveredTree owns its whole closure, so only a source or checkout refusal can stop it.
func coveredTree() ownershiptest.ProposalRepo {
	return ownershiptest.ProposalRepo{
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
		args := proposalArgs(t, root)
		args[4] = "missing.md"
		preflighttest.MustWriteFile(t, "dirty.txt", "dirty\n")
		out, code := preflight.Command(args)
		if code == 0 || !strings.Contains(out, "checkout") || strings.Contains(out, "selected ticket") {
			t.Fatalf("dirty checkout with a missing ticket = (%d), want the checkout refusal first:\n%s", code, out)
		}
	})
	t.Run("OC90 evidence source mismatch", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := proposalArgs(t, root)
		args[2] = "--charge"
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
		args := proposalArgs(t, root)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			preflighttest.MustWriteFile(t, "moved.txt", strings.Repeat("x", calls)+"\n")
			preflighttest.RunGit(t, "add", "moved.txt")
			preflighttest.RunGit(t, "commit", "-q", "-m", "move")
		})
		out, code := preflight.Command(args)
		restore()
		if calls != 2 {
			t.Fatalf("movement callbacks = %d, want 2", calls)
		}
		refused(t, out, code, "snapshot drift")
	})
	t.Run("OC81 stale source tip", func(t *testing.T) {
		root := coveredTree().Seed(t)
		args := proposalArgs(t, root)
		args[8] = preflighttest.RunGit(t, "rev-parse", "main")
		out, code := preflight.Command(args)
		refused(t, out, code, "tip-current")
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
