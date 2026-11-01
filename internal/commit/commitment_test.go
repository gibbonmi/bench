package commit

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// landingLabel is the label of the assignment that owns the worktree of landingRepo.
const landingLabel = "commit-fixture"

// landingRepo builds the minimal linked worktree `bench commit` lands into. It returns
// that worktree and its pre-landing HEAD.
func landingRepo(t *testing.T, gateExit int, write func(t *testing.T, root string)) (root, before string) {
	t.Helper()
	notTheKitRoot(t)
	primary := t.TempDir()
	initializeLandingRepo(t, primary, gateExit)
	root = filepath.Join(t.TempDir(), "linked")
	gittest.MustRun(t, primary, "worktree", "add", "-q", "-b", "topic", root)
	commitmenttest.Register(t, root, landingLabel)
	commitmenttest.Admit(t, root, landingLabel, "specs/commit-fixture/spec.md")
	return prepareLandingCheckout(t, root, write)
}

func primaryLandingRepo(t *testing.T, gateExit int, write func(t *testing.T, root string)) (root, before string) {
	t.Helper()
	notTheKitRoot(t)
	root = t.TempDir()
	initializeLandingRepo(t, root, gateExit)
	return prepareLandingCheckout(t, root, write)
}

// notTheKitRoot points the kit-root selection away from the fixture. Every fixture here
// is a linked project, not the Bench kit, so it declares its lane in a phase manifest or
// declares none. Without this the selection would answer the fixture itself and hand it
// the kit's built-in lane.
func notTheKitRoot(t *testing.T) {
	t.Helper()
	t.Setenv("BENCH_KIT", t.TempDir())
}

func initializeLandingRepo(t *testing.T, root string, gateExit int) {
	t.Helper()
	gittest.MustRun(t, root, "init", "-q", "-b", "main")
	for _, args := range gittest.IdentityConfig() {
		gittest.MustRun(t, root, args...)
	}
	mustMkdirAll(t, filepath.Join(root, ".bench"))
	mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), "#!/bin/sh\nexit "+strconv.Itoa(gateExit)+"\n", 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	mustWrite(t, filepath.Join(root, "tracked.txt"), "base\n", 0o644)
	commitmenttest.Write(t, root, "specs/commit-fixture/spec.md", "# Commit fixture\n\nStatus: staged\n")
	commitmenttest.SeedAdmission(t, root, "specs/commit-fixture/spec.md")
	gittest.MustRun(t, root, "add", "-A")
	gittest.Commit(t, root, "bootstrap")
}

func prepareLandingCheckout(t *testing.T, root string, write func(t *testing.T, root string)) (string, string) {
	t.Helper()
	write(t, root)
	gittest.MustRun(t, root, "add", "-A")
	gittest.Commit(t, root, "base", "--allow-empty")
	mustWrite(t, filepath.Join(root, "tracked.txt"), "changed\n", 0o644)
	gittest.MustRun(t, root, "add", "tracked.txt")
	return root, gittest.Output(t, root, "rev-parse", "HEAD")
}

func planningCommitRepo(t *testing.T) string {
	t.Helper()
	notTheKitRoot(t)
	primary := t.TempDir()
	initializeLandingRepo(t, primary, 0)
	commitmenttest.SeedProtected(t, primary)
	commitmenttest.Commit(t, primary, "approve protected obligations")
	return commitmenttest.Planning(t, primary)
}

func refuseCommit(t *testing.T, root string, paths ...string) string {
	t.Helper()
	before := gittest.Output(t, root, "rev-parse", "HEAD")
	args := append([]string{"-m", "candidate", "--"}, paths...)
	code, out, errOut := runCommand(t, root, args...)
	if code != 1 || !strings.Contains(errOut, "commitment") {
		t.Fatalf("commit = %d %s %s", code, out, errOut)
	}
	if after := gittest.Output(t, root, "rev-parse", "HEAD"); after != before {
		t.Fatal("refused commit moved HEAD")
	}
	return errOut
}

func TestCommitmentProposalIdentity(t *testing.T) {
	root := planningCommitRepo(t)
	store := commitrepo.Store{Root: root}
	policy, _, err := store.Policy()
	if err != nil {
		t.Fatal(err)
	}
	items := policy.Milestones[0].Outcomes
	items[0], items[1] = items[1], items[0]
	data, err := commitment.Bytes(policy)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.Plan(data)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Approve(plan.ID, "reviewer direction", []string{"A"}, nil); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := runCommand(t, root, "--dry-run", "-m", "approved transition", "--", ".bench/commitment.json", "ROADMAP.md")
	if code != 0 {
		t.Fatalf("exact approved candidate = %d %s %s", code, out, errOut)
	}
	policy.Milestones[0].Outcomes[0].Criteria[0].Text = "A different criterion."
	commitmenttest.WritePolicy(t, root, policy)
	if out := refuseCommit(t, root, ".bench/commitment.json", "ROADMAP.md"); !strings.Contains(out, "exact approval") {
		t.Fatal(out)
	}
}

func TestCommitmentUncommittedAppend(t *testing.T) {
	root := planningCommitRepo(t)
	body, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.Write(t, root, "ROADMAP.md", strings.Replace(string(body), "## Recommended sequence", "**FT3 — C**\n\n## Recommended sequence", 1))
	commitmenttest.Write(t, root, "roadmap/FT3.md", "**FT3 — C**\n\nAn uncommitted finding.\n")
	code, out, errOut := runCommand(t, root, "-m", "append discovery", "--", "ROADMAP.md", "roadmap/FT3.md")
	if code != 0 {
		t.Fatalf("append = %d %s %s", code, out, errOut)
	}
	if got := gittest.Output(t, root, "show", "HEAD:ROADMAP.md"); !strings.Contains(got, "1. A\n2. B") {
		t.Fatal(got)
	}
}

func TestCommitmentProtectedDeletion(t *testing.T) {
	for _, kind := range []string{"detail", "index", "rename"} {
		t.Run(kind, func(t *testing.T) {
			root := planningCommitRepo(t)
			if kind != "index" {
				if err := os.Remove(filepath.Join(root, "roadmap/FT1.md")); err != nil {
					t.Fatal(err)
				}
			}
			if kind != "detail" {
				body, _ := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
				replacement := ""
				if kind == "rename" {
					replacement = "**FT9 — A**"
					commitmenttest.Write(t, root, "roadmap/FT9.md", "**FT9 — A**\n\nKeep the A obligation.\n")
				}
				commitmenttest.Write(t, root, "ROADMAP.md", strings.Replace(string(body), "**FT1 — A**", replacement, 1))
			}
			refuseCommit(t, root, "ROADMAP.md", "roadmap")
		})
	}
}

func TestCommitmentProtectedScope(t *testing.T) {
	for _, kind := range []string{"requirement", "sequence", "event-label"} {
		t.Run(kind, func(t *testing.T) {
			root := planningCommitRepo(t)
			path := "roadmap/FT1.md"
			body := "**FT1 — A**\n\nReplace the A obligation.\n"
			if kind == "event-label" {
				body = "**FT1 — A**\n\nOccurrences: replacement requirement prose\n"
			}
			if kind == "sequence" {
				path = "ROADMAP.md"
				data, _ := os.ReadFile(filepath.Join(root, path))
				body = strings.Replace(string(data), "1. A\n2. B", "1. C\n2. B", 1)
			}
			commitmenttest.Write(t, root, path, body)
			refuseCommit(t, root, path)
		})
	}
}

// An approved plan that changes the identity of FT1 commits with the policy, the board, and
// the row edit in one commit. The commit refuses when it names no row edit, and when the row
// bytes changed after the approval, and a refused commit leaves HEAD in place.
func TestCommitmentApprovedRowEdit(t *testing.T) {
	const path = "roadmap/FT1.md"
	const approved = "**FT1 — A**\n\nRecord the settled A decision.\n"
	const rewritten = "**FT1 — A**\n\nAn edit after the approval.\n"
	approvedPaths := []string{commitment.PolicyPath, "ROADMAP.md", path}
	for _, row := range []struct {
		name, after, got string
		paths            []string
	}{
		{name: "approved", paths: approvedPaths},
		{name: "policy-only", paths: approvedPaths[:2]},
		{name: "rewritten", after: rewritten, got: rewritten, paths: approvedPaths},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := planningCommitRepo(t)
			published, err := os.ReadFile(filepath.Join(root, path))
			if err != nil {
				t.Fatal(err)
			}
			store := commitrepo.Store{Root: root}
			policy, _, err := store.Policy()
			if err != nil {
				t.Fatal(err)
			}
			commitmenttest.Write(t, root, path, approved)
			policy.Milestones[0].Outcomes[0].Sources[0].Identity = commitment.Identity([]byte(approved))
			data, err := commitment.Bytes(policy)
			if err != nil {
				t.Fatal(err)
			}
			plan, err := store.Plan(data)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := store.Approve(plan.ID, "reviewer direction", plan.Effects.Delayed, plan.Effects.Removed); err != nil {
				t.Fatal(err)
			}
			if row.after != "" {
				commitmenttest.Write(t, root, path, row.after)
			}
			if row.name != "approved" {
				got := row.got
				if got == "" {
					got = string(published)
				}
				want := commitmenttest.ChangedFrom(policy.Milestones[0].Outcomes[0].Sources[0], got)
				// The approval check refuses before the protected-commitment check, so the
				// refusal keeps the face of the approval check.
				const protected = "candidate changes protected commitment"
				if errOut := refuseCommit(t, root, row.paths...); !strings.Contains(errOut, want) || strings.Contains(errOut, protected) {
					t.Fatalf("commit stderr = %q, want %q from the approval check, not %q", errOut, want, protected)
				}
				return
			}
			before := gittest.Output(t, root, "rev-parse", "HEAD")
			code, out, errOut := runCommand(t, root, append([]string{"-m", "record the FT1 decision", "--"}, row.paths...)...)
			if code != 0 {
				t.Fatalf("approved row edit commit = %d %s %s", code, out, errOut)
			}
			if parent := gittest.Output(t, root, "rev-parse", "HEAD^"); parent != before {
				t.Fatalf("HEAD^ = %s, want one commit on %s", parent, before)
			}
			if got := gittest.Output(t, root, "show", "HEAD:"+path); got+"\n" != approved {
				t.Fatalf("committed %s = %q, want %q", path, got, approved)
			}
		})
	}
}

func TestCommitmentCommitBeforeEffects(t *testing.T) {
	for _, prior := range []bool{false, true} {
		t.Run(strconv.FormatBool(prior), func(t *testing.T) {
			root := planningCommitRepo(t)
			body := []byte("package example\nfunc Value( )int{return 1}\n")
			commitmenttest.Write(t, root, "change.go", string(body))
			paths := []string{"change.go"}
			if prior {
				commitmenttest.Commit(t, root, "earlier unbound source")
				commitmenttest.Write(t, root, "capture/note.md", "# Planning note\n")
				paths = []string{"capture/note.md"}
			}
			refuseCommit(t, root, paths...)
			after, err := os.ReadFile(filepath.Join(root, "change.go"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(body, after) {
				t.Fatal("formatter changed an unadmitted candidate")
			}
		})
	}
}

// An unbound commit lands a production path inside the Writes line of its one light-path
// ticket. A path outside that line refuses with the path, the ticket, and the route that
// widens the ticket's Writes line, and it names no commitment start.
func TestCommitmentLightPathCommit(t *testing.T) {
	root := planningCommitRepo(t)
	const ticket = "specs/lp/tickets/one.md"
	commitmenttest.WriteLightTicket(t, root, ticket, "change.go")
	commitmenttest.Commit(t, root, "light-path ticket")
	commitmenttest.Write(t, root, "change.go", "package example\n")
	if code, out, errOut := runCommand(t, root, "-m", "light-path change", "--", "change.go"); code != 0 {
		t.Fatalf("covered light-path commit = %d %s %s", code, out, errOut)
	}
	commitmenttest.Write(t, root, "other.go", "package example\n")
	errOut := refuseCommit(t, root, "other.go")
	for _, want := range []string{`"other.go"`, `"` + ticket + `"`, "next=add the path to the Writes: line of '" + ticket + "'"} {
		if !strings.Contains(errOut, want) {
			t.Fatalf("uncovered light-path commit stderr = %q, want %q", errOut, want)
		}
	}
	if strings.Contains(errOut, "bench commitment start") {
		t.Fatalf("uncovered light-path commit stderr = %q, want no commitment start", errOut)
	}
}

func TestCommitmentPlanningOccurrence(t *testing.T) {
	root := planningCommitRepo(t)
	path := "roadmap/FT1.md"
	before, _ := os.ReadFile(filepath.Join(root, path))
	commitmenttest.Write(t, root, path, string(before)+"Occurrences: incident-1\n")
	code, out, errOut := runCommand(t, root, "-m", "record occurrence", "--", path)
	if code != 0 {
		t.Fatalf("occurrence = %d %s %s", code, out, errOut)
	}
}
