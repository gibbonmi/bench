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
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/responsebound"
)

// Each test drives the actual preflight command through committed fixture objects. The
// expected drift rows, table headers, and refusal texts are stated here independently of the
// audit renderer, and each expected object ID comes from git rev-parse, not from the
// collector's tree listing. A dropped member, a collapsed kind, or a renamed column turns a
// case red.

// absent is the cell the audit prints for the side of a drift entry that holds no member.
const absent = "absent"

// audit runs the audit-spec form over one explicit pair, with any extra arguments after it.
func audit(t *testing.T, slug, base, tip string, extra ...string) (string, int) {
	t.Helper()
	args := append([]string{"build", slug, "--audit-spec", "--base", base, "--source-tip", tip}, extra...)
	return preflight.CommandWithVersion("staleness-test")(args)
}

// write writes each file under the fixture's working directory.
func write(t *testing.T, files map[string]string) {
	t.Helper()
	for path, body := range files {
		preflighttest.MustWriteFile(t, path, body)
	}
}

// commitAll commits every change of the working tree and returns the new head.
func commitAll(t *testing.T, message string) string {
	t.Helper()
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", message)
	return preflighttest.RunGit(t, "rev-parse", "HEAD")
}

// object is the object ID that commit holds at path.
func object(t *testing.T, commit, path string) string {
	t.Helper()
	return preflighttest.RunGit(t, "rev-parse", commit+":"+path)
}

// table decodes one table of a response. A table opens at the start of a line, and its rows
// are indented, so the first unindented line after the header closes it.
func table(t *testing.T, out, name string) []map[string]any {
	t.Helper()
	var block strings.Builder
	for _, line := range strings.SplitAfter(out, "\n") {
		if block.Len() == 0 && strings.HasPrefix(line, name+"[") || block.Len() > 0 && strings.HasPrefix(line, "  ") {
			block.WriteString(line)
		} else if block.Len() > 0 {
			break
		}
	}
	if block.Len() == 0 {
		t.Fatalf("response holds no %s table:\n%s", name, out)
	}
	return preflighttest.TableRows(t, preflighttest.DecodeMap(t, block.String()), name)
}

// driftRows is each drift row as before|after|kind|before_object|after_object, sorted.
func driftRows(t *testing.T, out string) []string {
	t.Helper()
	var rows []string
	for _, row := range table(t, out, "drift") {
		rows = append(rows, fmt.Sprintf("%v|%v|%v|%v|%v", row["before"], row["after"], row["kind"], row["before_object"], row["after_object"]))
	}
	slices.Sort(rows)
	return rows
}

// sorted returns rows in sorted order.
func sorted(rows ...string) []string {
	slices.Sort(rows)
	return rows
}

// dropBaselineSpec commits the spec's removal and then its return. The first commit is a
// baseline that holds no spec, and the second is a tip that holds the same spec bytes.
func dropBaselineSpec(t *testing.T, slug string) (base, tip string) {
	t.Helper()
	spec := "specs/" + slug + "/spec.md"
	body, err := os.ReadFile(spec)
	if err != nil {
		t.Fatal(err)
	}
	preflighttest.RunGit(t, "rm", "-q", spec)
	base = commitAll(t, "baseline without the spec")
	write(t, map[string]string{spec: string(body)})
	return base, commitAll(t, "spec returns")
}

// TestDriftObjects covers SS1, SS2, SS3, SS4, SS5, SS6, SS10, and SS11.
func TestDriftObjects(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	write(t, map[string]string{
		"internal/example/mod.txt": "one\n", "internal/example/gone.txt": "gone\n", "internal/example/old.txt": "moved bytes\n",
		"internal/example/run.sh": "#!/bin/sh\n", "internal/example/.hidden": "dot one\n", "internal/example2/collide.txt": "outside one\n",
	})
	base := commitAll(t, "baseline members")

	out, code := audit(t, slug, base, base)
	if code != 0 || !strings.Contains(out, "drift[0]{before,after,kind,before_object,after_object}:\n") || !strings.Contains(out, "spec-drift,green") {
		t.Fatalf("SS10 unchanged members = (%d), want a green spec-drift row and a complete empty drift table:\n%s", code, out)
	}
	// The symbol and test collectors do not run, so a clean drift result is still no complete audit.
	if row := table(t, out, "staleness")[0]; row["complete"] != false || row["manual_required"] != true {
		t.Fatalf("SS10 staleness = %v, want complete false and manual review while collectors are missing:\n%s", row, out)
	}

	write(t, map[string]string{
		"internal/example/mod.txt": "two\n", "internal/example/new.txt": "new\n",
		"internal/example/.hidden": "dot two\n", "internal/example2/collide.txt": "outside two\n",
	})
	preflighttest.RunGit(t, "rm", "-q", "internal/example/gone.txt")
	preflighttest.RunGit(t, "mv", "internal/example/old.txt", "internal/example/moved.txt")
	if err := os.Chmod("internal/example/run.sh", 0o755); err != nil {
		t.Fatal(err)
	}
	tip := commitAll(t, "drifted members")
	same := func(path, kind string) string {
		return path + "|" + path + "|" + kind + "|" + object(t, base, path) + "|" + object(t, tip, path)
	}
	want := sorted(
		same("internal/example/mod.txt", "modified"),
		same("internal/example/.hidden", "modified"),
		same("internal/example/run.sh", "mode 100644->100755"),
		"internal/example/gone.txt|"+absent+"|deleted|"+object(t, base, "internal/example/gone.txt")+"|"+absent,
		absent+"|internal/example/new.txt|added|"+absent+"|"+object(t, tip, "internal/example/new.txt"),
		"internal/example/old.txt|internal/example/moved.txt|renamed|"+object(t, base, "internal/example/old.txt")+"|"+object(t, tip, "internal/example/moved.txt"),
	)
	out, code = audit(t, slug, base, tip)
	if got := driftRows(t, out); code != 1 || !slices.Equal(got, want) || !strings.Contains(out, "spec-drift,red") {
		t.Fatalf("drifted members = (%d)\n got %q\nwant %q\n%s", code, got, want, out)
	}
}

// TestDriftIdentity covers SS7, SS8, SS9, and SS12.
func TestDriftIdentity(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	spec := "specs/" + slug + "/spec.md"

	base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
	out, _ := audit(t, slug, "main", "HEAD")
	source, staleness := table(t, out, "source")[0], table(t, out, "staleness")[0]
	if source["base"] != base || source["tip"] != tip || staleness["baseline"] != base || staleness["tip"] != tip {
		t.Fatalf("SS7 identities = %v and %v, want the resolved %s and %s:\n%s", source, staleness, base, tip, out)
	}

	before := tip
	write(t, map[string]string{spec: preflighttest.SpecBody(slug) + "\nA prose correction.\n"})
	tip = commitAll(t, "prose correction")
	want := spec + "|" + spec + "|modified|" + object(t, before, spec) + "|" + object(t, tip, spec)
	if out, code := audit(t, slug, before, tip); code != 1 || !strings.Contains(out, ",red=0}\n") || !slices.Equal(driftRows(t, out), []string{want}) {
		t.Fatalf("SS9 spec drift = (%d) %q, want the audit red alone and %q:\n%s", code, driftRows(t, out), want, out)
	}

	claim := `{"id":"claim","location":{"path":"` + spec + `","line":1,"needle":"# ` + slug + `","source":"spec"},` +
		`"path":"outside/claim_test.go","package":"example.com/outside","function":"TestClaim","state":"existing","rows":["PF1"]}`
	write(t, map[string]string{
		spec:                    preflighttest.SpecBodyNoInventory(slug) + inventorytest.Section(inventorytest.Manifest("declared-complete", nil, []string{claim}, nil, nil)),
		"outside/claim_test.go": "package outside\n",
	})
	before = commitAll(t, "claim outside the fence")
	write(t, map[string]string{"outside/claim_test.go": "package outside\n\n// changed premise owner\n"})
	tip = commitAll(t, "claim source changes")
	want = "outside/claim_test.go|outside/claim_test.go|modified|" + object(t, before, "outside/claim_test.go") + "|" + object(t, tip, "outside/claim_test.go")
	if out, _ := audit(t, slug, before, tip); !slices.Equal(driftRows(t, out), []string{want}) {
		t.Fatalf("SS12 claim drift = %q, want %q:\n%s", driftRows(t, out), want, out)
	}

	base, tip = dropBaselineSpec(t, slug)
	if out, code := audit(t, slug, base, tip); code != 1 || !strings.Contains(out, "baseline spec refused") || strings.Contains(out, "staleness[") {
		t.Fatalf("SS8 absent baseline spec = (%d), want a refusal and no staleness result:\n%s", code, out)
	}
}

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

// TestAuditTrust covers SS38 and SS56.
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

// TestAuditSafety covers SS51: a link in scope refuses, and its outside target stays unread.
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
	tip := commitAll(t, "link in scope")
	out, code := audit(t, slug, base, tip)
	if code != 1 || !strings.Contains(out, "special source refused") || strings.Contains(out, sentinel) || strings.Contains(out, "staleness[") {
		t.Fatalf("special source = (%d), want a refusal that follows nothing:\n%s", code, out)
	}
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
