package staleness_test

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// Each test of this package drives the actual preflight command through committed fixture
// objects, except the committed tip inventory case, whose input no command reaches past the
// clean-source refusal. The tests state the expected drift rows, table headers, and refusal
// texts independently of the audit renderer, and each expected object ID comes from git
// rev-parse, not from the collector's tree listing.

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
	spec := specPath(slug)
	body, err := os.ReadFile(spec)
	if err != nil {
		t.Fatal(err)
	}
	preflighttest.RunGit(t, "rm", "-q", spec)
	base = commitAll(t, "baseline without the spec")
	write(t, map[string]string{spec: string(body)})
	return base, commitAll(t, "spec returns")
}

// specPath is the repository path of the fixture spec.
func specPath(slug string) string {
	return "specs/" + slug + "/spec.md"
}

// itemLocation is the inventory location of an item that the fixture spec's title line holds.
func itemLocation(slug string) string {
	return `"location":{"path":"` + specPath(slug) + `","line":1,"needle":"# ` + slug + `","source":"spec"}`
}

// seamClaim is one inventory test seam whose path is the claim source path.
func seamClaim(slug, path string) string {
	return `{"id":"claim",` + itemLocation(slug) + `,"path":"` + path + `","package":"example.com/outside","function":"TestClaim","state":"existing","rows":["PF1"]}`
}

// inventorySpec is the fixture spec whose inventory holds each given premise and seam.
func inventorySpec(slug string, premises, seams []string) string {
	return preflighttest.SpecBodyNoInventory(slug) + inventorytest.Section(inventorytest.Manifest("declared-complete", premises, seams, nil, nil))
}

// TestDriftObjects covers SS1, SS2, SS3, SS4, SS5, SS6, SS10, and SS11.
func TestDriftObjects(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	// Git quotes this path in every listing form without NUL framing.
	const quoted = "internal/example/q\"b\\s\tü.txt"
	write(t, map[string]string{
		"internal/example/mod.txt": "one\n", "internal/example/gone.txt": "gone\n", "internal/example/old.txt": "moved bytes\n",
		"internal/example/run.sh": "#!/bin/sh\n", "internal/example/.hidden": "dot one\n", "internal/example2/collide.txt": "outside one\n",
		quoted: "quoted one\n",
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
		"internal/example/.hidden": "dot two\n", "internal/example2/collide.txt": "outside two\n", quoted: "quoted two\n",
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
		same(quoted, "modified"),
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
	spec := specPath(slug)

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

	// The reused-owner premise names its owner path only inside a nested subject object.
	owner := `{"id":"owner",` + itemLocation(slug) + `,"kind":"reused-owner","state":"current",` +
		`"subject":{"owner":{"path":"outside/owner.go","line":1,"needle":"package outside","source":"tip"}},"expected":{"required":true}}`
	write(t, map[string]string{
		spec:                    inventorySpec(slug, []string{owner}, []string{seamClaim(slug, "outside/claim_test.go")}),
		"outside/claim_test.go": "package outside\n",
		"outside/owner.go":      "package outside\n",
	})
	before = commitAll(t, "claims outside the fence")
	write(t, map[string]string{"outside/claim_test.go": "package outside\n\n// changed seam\n", "outside/owner.go": "package outside\n\n// changed premise owner\n"})
	tip = commitAll(t, "claim sources change")
	modified := func(path string) string {
		return path + "|" + path + "|modified|" + object(t, before, path) + "|" + object(t, tip, path)
	}
	claims := sorted(modified("outside/claim_test.go"), modified("outside/owner.go"))
	if out, _ := audit(t, slug, before, tip); !slices.Equal(driftRows(t, out), claims) {
		t.Fatalf("SS12 claim drift = %q, want %q:\n%s", driftRows(t, out), claims, out)
	}

	// The tip drops both claims while their sources change, so only the baseline names them.
	before = tip
	write(t, map[string]string{
		spec:                    inventorySpec(slug, nil, nil),
		"outside/claim_test.go": "package outside\n\n// dropped seam\n",
		"outside/owner.go":      "package outside\n\n// dropped premise owner\n",
	})
	tip = commitAll(t, "claims drop while their sources change")
	dropped := sorted(modified(spec), modified("outside/claim_test.go"), modified("outside/owner.go"))
	if out, _ := audit(t, slug, before, tip); !slices.Equal(driftRows(t, out), dropped) {
		t.Fatalf("dropped claim drift = %q, want %q:\n%s", driftRows(t, out), dropped, out)
	}

	base, tip = dropBaselineSpec(t, slug)
	if out, code := audit(t, slug, base, tip); code != 1 || !strings.Contains(out, "baseline spec refused") || strings.Contains(out, "staleness[") {
		t.Fatalf("SS8 absent baseline spec = (%d), want a refusal and no staleness result:\n%s", code, out)
	}
}
