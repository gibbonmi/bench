package preflight

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// seedBuildFresh builds the build-mode fresh fixture. A base commit on
// main carries a bootstrap-conformant spec with no tickets/ directory at
// all, then a feature branch adds one authorized change. B1 answers
// not-applicable ticket rows over this tree.
func seedBuildFresh(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug))
	commitmenttest.SeedAdmission(t, root, "specs/"+slug+"/spec.md")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/"+slug+"/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	preflighttest.ActiveAssignment(t, root, root)
	commitmenttest.Admit(t, root, "preflight-assignment-target", "specs/"+slug+"/spec.md")
	return root, slug
}

func unbindBuild(t *testing.T, root string) []byte {
	t.Helper()
	if err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		ledger.Commitment = nil
		return ledger, true, nil
	}, nil); err != nil {
		t.Fatal(err)
	}
	path, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertBuildLedger(t *testing.T, root string, before []byte) {
	t.Helper()
	path, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("preflight changed delivery claims")
	}
}

func TestCommitmentBuildPreflight(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	args := preflighttest.ChargeArgs(t, root, slug)
	before := unbindBuild(t, root)
	for _, args := range [][]string{{"build", slug}, args} {
		out, code := Command(args)
		if code != 1 || !strings.Contains(out, "bench commitment start") {
			t.Fatalf("unbound preflight = %d %s", code, out)
		}
		assertBuildLedger(t, root, before)
		preflighttest.AssertNothingPublished(t, root)
	}
}

func TestCommitmentPlanOnly(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	before := unbindBuild(t, root)
	out, code := Command([]string{"build", slug, "--plan-only"})
	if code != 0 {
		t.Fatalf("plan-only = %d %s", code, out)
	}
	assertBuildLedger(t, root, before)
	for _, args := range [][]string{{"build", slug, "--plan-only", "--charge"}, {"review", slug, "--plan-only"}} {
		out, code = Command(args)
		if code != 2 {
			t.Fatalf("invalid planning form = %d %s", code, out)
		}
	}
	preflighttest.AssertNothingPublished(t, root)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.TicketDoc("One", "PF1"))
	out, code = Command([]string{"build", slug, "--plan-only"})
	if code != 1 || !strings.Contains(out, "rows-owned,red") {
		t.Fatalf("invalid graph = %d %s", code, out)
	}
	assertBuildLedger(t, root, before)
}
