package staleness_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/preflight/staleness"
	"github.com/gibbonmi/bench/internal/responsebound"
)

// TestAuditPrecedence covers SS35: the checkout refusal wins over a staleness refusal.
func TestAuditPrecedence(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	base, tip := dropBaselineSpec(t, slug)
	write(t, map[string]string{"dirty.txt": "dirty\n"})
	out, code := audit(t, slug, base, tip)
	if code != 1 || !strings.Contains(out, "source not clean") || strings.Contains(out, "baseline spec") || strings.Contains(out, "staleness[") {
		t.Fatalf("dirty source = (%d), want the checkout refusal before any staleness result:\n%s", code, out)
	}
}

// TestAuditMovement covers SS36 and SS37.
func TestAuditMovement(t *testing.T) {
	t.Run("SS36 persistent movement", func(t *testing.T) {
		_, slug := preflighttest.SeedAdmitted(t)
		base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
		// Each read ends with a new commit, so the checkout stays clean and only the source moves.
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			preflighttest.RunGit(t, "commit", "-q", "--allow-empty", "-m", fmt.Sprintf("move %d", calls))
		})
		defer restore()
		out, code := audit(t, slug, base, tip)
		if code != 1 || calls != 2 || !strings.Contains(out, "snapshot drift") || strings.Contains(out, "staleness[") || strings.Contains(out, "drift[") {
			t.Fatalf("persistent movement = (%d, %d calls), want the drift refusal and no audit result:\n%s", code, calls, out)
		}
	})
	t.Run("SS37 source-tip mismatch", func(t *testing.T) {
		// The pair holds no drift, so only the readiness red can make the audit exit 1.
		_, slug := preflighttest.SeedAdmitted(t)
		out, code := audit(t, slug, preflighttest.RunGit(t, "rev-parse", "HEAD"), preflighttest.RunGit(t, "rev-parse", "main"))
		if code != 1 || !strings.Contains(out, "tip-current,red") || table(t, out, "staleness")[0]["manual_required"] != true {
			t.Fatalf("stale pin = (%d), want the red tip-current row and manual review:\n%s", code, out)
		}
	})
}

// TestAuditTrust covers SS38, SS56, and the committed tip inventory.
func TestAuditTrust(t *testing.T) {
	t.Run("SS38 executable seal", func(t *testing.T) {
		root, slug := preflighttest.SeedAdmitted(t)
		write(t, map[string]string{".gitignore": "dist/\n"})
		tip := commitAll(t, "ignore the published binary")
		executable := freshness.PublishedExecutable(root)
		if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(root, "absent-bench"), executable); err != nil {
			t.Fatal(err)
		}
		out, code := audit(t, slug, tip, tip)
		if code != 1 || !strings.Contains(out, "binary-seal,red") || table(t, out, "staleness")[0]["manual_required"] != true {
			t.Fatalf("refused seal = (%d), want the red binary-seal row and manual review:\n%s", code, out)
		}
	})
	t.Run("SS56 current-action binding", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		prepared, code := preflight.Command(args)
		if code != 0 {
			t.Fatalf("prepare = (%d):\n%s", code, prepared)
		}
		identity := table(t, prepared, "prepared")[0]["evidence"].(string)
		check := []string{"evidence", identity, "--check-current"}
		if out, _ := audit(t, slug, args[6], args[8]); strings.Contains(out, identity) || len(preflighttest.PublishedPacks(t, root)) != 1 {
			t.Fatalf("audit touched the prepared evidence:\n%s", out)
		}
		if out, code := preflight.Command(check); code != 0 {
			t.Fatalf("binding after the audit = (%d):\n%s", code, out)
		}
		write(t, map[string]string{".agents/skills/bench-craft-delegate/SKILL.md": "# Delegation skill, changed\n"})
		tip := commitAll(t, "required source changes")
		out, _ := audit(t, slug, args[6], tip)
		if !slices.ContainsFunc(driftRows(t, out), func(row string) bool { return strings.HasPrefix(row, ".agents/skills/bench-craft-delegate/SKILL.md|") }) {
			t.Fatalf("audit does not report the changed required source:\n%s", out)
		}
		if out, code := preflight.Command(check); code != 1 || !strings.Contains(out, "source required") {
			t.Fatalf("binding after the source moved = (%d), want the source refusal:\n%s", code, out)
		}
	})
	t.Run("committed tip inventory", func(t *testing.T) {
		root, slug := preflighttest.SeedAdmitted(t)
		write(t, map[string]string{"outside/claim_test.go": "package outside\n"})
		base := commitAll(t, "unclaimed source")
		write(t, map[string]string{"outside/claim_test.go": "package outside\n\n// changed\n"})
		tip := commitAll(t, "unclaimed source changes")
		// Only this inventory names the source, and no committed spec holds it.
		checkout := ownership.ParseManifest([]byte(inventorySpec(slug, nil, []string{seamClaim(slug, "outside/claim_test.go")})))
		facts, err := staleness.Collect(root, ownership.Source{Base: base, Tip: tip, SpecPath: specPath(slug)}, checkout, ownership.Snapshot{})
		claimed := slices.ContainsFunc(facts.Paths, func(state staleness.PathState) bool { return state.Path == "outside/claim_test.go" })
		held := slices.ContainsFunc(facts.Diagnostics, func(d staleness.Diagnostic) bool { return d.Subject == ownership.MetadataLabel })
		if err != nil || claimed || !held {
			t.Fatalf("checkout inventory = (%v), scope holds its claim %t, diagnostic %t; want the committed tip scope and a %s diagnostic", err, claimed, held, ownership.MetadataLabel)
		}
	})
}

// TestAuditGrammar covers SS39 and SS40.
func TestAuditGrammar(t *testing.T) {
	_, slug := preflighttest.SeedConformant(t)
	base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
	for _, selector := range []string{"--charge", "--plan-only", "--propose-writes", "--breakdown", "--check-current", "--verify"} {
		out, code := audit(t, slug, base, tip, selector)
		if code != 2 || !strings.Contains(out, "cannot be combined") || strings.Contains(out, "staleness[") {
			t.Errorf("--audit-spec %s = (%d), want exit 2 and no collection:\n%s", selector, code, out)
		}
	}
	for _, args := range [][]string{
		{"build", slug, "--audit-spec", "--source-tip", tip},
		{"build", slug, "--audit-spec", "--base", base},
		{"build", slug, "--audit-spec"},
	} {
		out, code := preflight.Command(args)
		if code != 2 || !strings.Contains(out, "--audit-spec requires build, --base, and --source-tip") {
			t.Errorf("%v = (%d), want exit 2 naming both identities:\n%s", args, code, out)
		}
	}
}

// TestAuditSafety covers SS51: a link in scope at either identity refuses, and its outside
// target stays unread.
func TestAuditSafety(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	base := preflighttest.RunGit(t, "rev-parse", "HEAD")
	const sentinel = "outside-sentinel-bytes"
	target := filepath.Join(t.TempDir(), "sentinel")
	if err := os.WriteFile(target, []byte(sentinel), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, "internal/example/link"); err != nil {
		t.Fatal(err)
	}
	refuses := func(side, base, tip string) {
		t.Helper()
		out, code := audit(t, slug, base, tip)
		if code != 1 || !strings.Contains(out, "special source refused") || strings.Contains(out, sentinel) || strings.Contains(out, "staleness[") {
			t.Errorf("special %s source = (%d), want a refusal that follows nothing:\n%s", side, code, out)
		}
	}
	linked := commitAll(t, "link in scope")
	refuses("after", base, linked)
	preflighttest.RunGit(t, "rm", "-q", "internal/example/link")
	refuses("before", linked, commitAll(t, "link leaves scope"))
}

// TestAuditBound covers SS52: the shared response bound spills a long audit whole.
func TestAuditBound(t *testing.T) {
	root, slug := preflighttest.SeedAdmitted(t)
	base := preflighttest.RunGit(t, "rev-parse", "HEAD")
	var want []string
	for i := range 40 {
		path := fmt.Sprintf("internal/example/bulk/file%02d.txt", i)
		write(t, map[string]string{path: path + "\n"})
		want = append(want, path)
	}
	tip := commitAll(t, "many fenced members")
	restore := evidencecmd.SetResponseLimitForTest(4096)
	defer restore()
	out, _ := audit(t, slug, base, tip)

	var stdout, stderr bytes.Buffer
	owner := responsebound.New(t.TempDir(), &stdout, &stderr, func() string { return root }, false)
	if _, err := io.WriteString(owner.Stdout(), out); err != nil {
		t.Fatal(err)
	}
	owner.Finish()
	frame, _, err := responsebound.First(stdout.Bytes())
	if err != nil || frame.Path == "" {
		t.Fatalf("long audit printed no spill frame (%v):\n%s", err, stdout.String())
	}
	spilled, err := os.ReadFile(frame.Path)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, row := range driftRows(t, string(spilled)) {
		got = append(got, strings.SplitN(row, "|", 3)[1])
	}
	if !slices.Equal(got, sorted(want...)) {
		t.Fatalf("spilled drift after paths = %q, want every one of %q", got, want)
	}
}
