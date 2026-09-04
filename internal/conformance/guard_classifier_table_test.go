package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
)

// guardClassifierRow is one command and the Bench-call verdict both classifiers owe it.
type guardClassifierRow struct {
	command string
	invokes bool
}

// guardClassifierTable is the one source both Bench-call classifiers answer to. The
// follow-on shim cannot call benchguard.InvokesBench, because the shim reaches for its
// own word test exactly when the binary that holds the Go classifier cannot answer. The
// two derivations are therefore real, and this table is what keeps them from drifting.
//
// Every row is resolver-independent: no head resolves through PATH or a symlink, so the
// verdict never depends on what the running machine has installed. The shell word test
// reads no wrapper string, so a `bash -c 'bench gate'` row belongs to the Go classifier
// alone and stays out.
//
// checkGuardClassifierTable parses these rows out of this file's source in whatever tree
// it grades, so the table is data for the shell half as much as it is code for the Go
// half. A row's shape is therefore load-bearing: one string literal, one bool.
var guardClassifierTable = []guardClassifierRow{
	{"bench", true},
	{"bench gate", true},
	{"bench gate 2>&1", true},
	{"./bin/bench.sh help", true},
	{"/opt/kit/bin/bench.sh gate", true},
	{"X=1 bench help", true},
	{"env X=1 bench help", true},
	{"env -i -- X=1 bench help", true},
	{"env -u X bench help", true},
	{"command -- bench help", true},
	{"nohup -- bench help", true},
	{"timeout -- 5 bench help", true},
	{"timeout -s KILL -k 1 5 bench help", true},
	{"xargs -n 1 bench help", true},
	{"xargs -- bench help", true},
	{"cat a && echo x; bench maps", true},
	{"cat file; bench help", true},
	// One row per quoting form the head word can carry. The shell half splits words
	// before it reads them, so it saw the quote characters as part of the name and
	// answered false on every one of these while the Go lexer answered true. Each row
	// is resolver-independent: the head resolves to the literal name `bench`, never
	// through PATH.
	{`"bench" gate`, true},
	{`'bench' gate`, true},
	{`\bench gate`, true},
	{`be"nch" gate`, true},
	{`ls; "bench" gate`, true},
	{`env "X=1" bench help`, true},
	// The false direction for the class. Quote removal concatenates rather than
	// splits, so the head is one word `benchgate`, which is not Bench.
	{`'bench'gate`, false},
	// A quote that spans the whitespace the shell half split on. Both halves arrive
	// half-quoted, so neither folds, and the head keeps its quote character. The
	// tokenizer reads one word `bench gate`, which is not Bench either.
	{`bench" "gate`, false},
	// One row per control operator that can stand between the command's head and a Bench
	// call, each with bench in a later segment. Without them the shell splitter's arms go
	// unpinned: a splitter narrowed to `;` alone still answers every other row correctly,
	// while letting `rg foo | bench help` pass a stale core as a non-Bench call.
	{"rg foo | bench help", true},
	{"cat a |& bench help", true},
	{"sleep 1 & bench help", true},
	{"cat a && bench help", true},
	{"cat a || bench help", true},
	// The false direction on a pipeline, so the operator rows are not one-directional. A
	// bench that is an argument rather than a head stays a non-Bench call.
	{"rg bench | wc -l", false},
	{"command -v bench", false},
	{"command -V bench", false},
	{"rg bench AGENTS.md", false},
	{"echo bench", false},
	{"printf hi > bench", false},
	{"benchtool run", false},
	{"go test ./...", false},
	{"git -C /tmp log", false},
	{"ls", false},
}

// TestGuardClassifierTableIsResolverIndependent proves the table's own premise. A row
// whose verdict moves when the resolver is taken away is a row the shell word test could
// never answer, and it would make the agreement check vacuous.
func TestGuardClassifierTableIsResolverIndependent(t *testing.T) {
	for _, row := range guardClassifierTable {
		if withResolver, without := benchguard.InvokesBench(row.command, benchguard.DefaultResolver()), benchguard.InvokesBench(row.command, benchguard.Resolver{}); withResolver != without {
			t.Errorf("InvokesBench(%q) = %t with a resolver and %t without; the row is not resolver-independent", row.command, withResolver, without)
		}
	}
}

// TestGoGuardClassifierAnswersTheSharedTable grades the Go half of the agreement. The
// shell half is checkGuardClassifierTable, which runs at the gate over a graded root.
func TestGoGuardClassifierAnswersTheSharedTable(t *testing.T) {
	for _, row := range guardClassifierTable {
		if got := benchguard.InvokesBench(row.command, benchguard.DefaultResolver()); got != row.invokes {
			t.Errorf("InvokesBench(%q) = %t, want %t", row.command, got, row.invokes)
		}
	}
}

// guardClassifierTableRelPath and guardClassifierLibRelPath anchor
// checkGuardClassifierTable. The shim's word test lives in the shared library rather
// than in the shim, because the shim reaches for it only when the Go classifier's binary
// cannot answer, and a sourceable function is the only seam a check can drive.
const (
	guardClassifierTableRelPath = "internal/conformance/guard_classifier_table_test.go"
	guardClassifierLibRelPath   = ".bench/lib/resolve-bench.sh"
	guardClassifierTableVar     = "guardClassifierTable"
	guardClassifierShellFunc    = "bench_invokes_bench"
)

// checkGuardClassifierTable grades the shell Bench-call word test against the shared
// table. It parses the row literals out of the graded root's own table file, sources that
// root's shared library, and asks bench_invokes_bench for each row. A row whose shell
// verdict disagrees with its declared verdict reds by name.
//
// The check applies only where internal/conformance/ exists, which is the kit source
// tree. A consumer repository carries the shared library but no table, and grading it
// there would refuse every consumer build. Inside that tree the check reds honestly
// rather than passing vacuously: an absent table file, a table with no rows, a renamed
// table, an absent library, a renamed word test, and an unrunnable shell each produce a
// diagnostic, so no amputation of the subject can quietly turn the check green.
func checkGuardClassifierTable(root string) []string {
	if info, err := os.Stat(filepath.Join(root, "internal", "conformance")); err != nil || !info.IsDir() {
		return nil
	}
	tablePath := filepath.Join(root, filepath.FromSlash(guardClassifierTableRelPath))
	rows, err := parseGuardClassifierRows(tablePath)
	if err != nil {
		return []string{fmt.Sprintf("guard classifier table check: %s: %v", guardClassifierTableRelPath, err)}
	}
	if len(rows) == 0 {
		return []string{fmt.Sprintf("guard classifier table check: %s declares no %s rows; the shell word test would pass vacuously", guardClassifierTableRelPath, guardClassifierTableVar)}
	}
	libPath := filepath.Join(root, filepath.FromSlash(guardClassifierLibRelPath))
	libText, err := os.ReadFile(libPath)
	if err != nil {
		return []string{fmt.Sprintf("guard classifier table check: %s is missing; cannot verify the shell word test", guardClassifierLibRelPath)}
	}
	if !strings.Contains(string(libText), guardClassifierShellFunc+"()") {
		return []string{fmt.Sprintf("guard classifier table check: %s has no %s() function to anchor on", guardClassifierLibRelPath, guardClassifierShellFunc)}
	}
	verdicts, err := runGuardClassifierWordTest(libPath, rows)
	if err != nil {
		return []string{fmt.Sprintf("guard classifier table check: cannot run the %s word test: %v", guardClassifierLibRelPath, err)}
	}
	var diags []string
	for index, row := range rows {
		if verdicts[index] != row.invokes {
			diags = append(diags, fmt.Sprintf("guard classifier row %q drifts: %s answers %t, the shared table declares %t", row.command, guardClassifierLibRelPath, verdicts[index], row.invokes))
		}
	}
	return diags
}

// parseGuardClassifierRows reads the table's row literals from Go source. Each row is one
// string literal and one bool, so a row that changes shape refuses rather than silently
// dropping out of the graded set.
func parseGuardClassifierRows(path string) ([]guardClassifierRow, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	var literal *ast.CompositeLit
	ast.Inspect(file, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if !ok || len(spec.Values) != 1 {
			return true
		}
		for _, name := range spec.Names {
			if name.Name != guardClassifierTableVar {
				continue
			}
			if composite, ok := spec.Values[0].(*ast.CompositeLit); ok {
				literal = composite
			}
		}
		return true
	})
	if literal == nil {
		return nil, fmt.Errorf("no %s composite literal to anchor on", guardClassifierTableVar)
	}
	rows := make([]guardClassifierRow, 0, len(literal.Elts))
	for _, element := range literal.Elts {
		row, ok := element.(*ast.CompositeLit)
		if !ok || len(row.Elts) != 2 {
			return nil, fmt.Errorf("a %s row is not a command-and-verdict pair", guardClassifierTableVar)
		}
		command, ok := row.Elts[0].(*ast.BasicLit)
		if !ok || command.Kind != token.STRING {
			return nil, fmt.Errorf("a %s row does not open with a command literal", guardClassifierTableVar)
		}
		text, err := strconv.Unquote(command.Value)
		if err != nil {
			return nil, err
		}
		verdict, ok := row.Elts[1].(*ast.Ident)
		if !ok || (verdict.Name != "true" && verdict.Name != "false") {
			return nil, fmt.Errorf("row %q does not declare a bool verdict", text)
		}
		if strings.ContainsAny(text, "\n\r") {
			return nil, fmt.Errorf("row %q holds a line break; the shell driver reads one row per line", text)
		}
		rows = append(rows, guardClassifierRow{command: text, invokes: verdict.Name == "true"})
	}
	return rows, nil
}

// guardClassifierDriver sources the shared library and answers one line per row.
const guardClassifierDriver = `. "$1"
while IFS= read -r line; do
  if bench_invokes_bench "$line"; then echo true; else echo false; fi
done`

// runGuardClassifierWordTest returns the shell word test's verdict for each row, in order.
func runGuardClassifierWordTest(libPath string, rows []guardClassifierRow) ([]bool, error) {
	shell, err := exec.LookPath("bash")
	if err != nil {
		return nil, err
	}
	commands := make([]string, 0, len(rows))
	for _, row := range rows {
		commands = append(commands, row.command)
	}
	cmd := exec.Command(shell, "-c", guardClassifierDriver, "bash", libPath)
	cmd.Stdin = strings.NewReader(strings.Join(commands, "\n") + "\n")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	answers := strings.Fields(string(out))
	if len(answers) != len(rows) {
		return nil, fmt.Errorf("the word test answered %d rows, want %d", len(answers), len(rows))
	}
	verdicts := make([]bool, len(rows))
	for index, answer := range answers {
		if answer != "true" && answer != "false" {
			return nil, fmt.Errorf("the word test answered %q for row %q", answer, rows[index].command)
		}
		verdicts[index] = answer == "true"
	}
	return verdicts, nil
}

// writeGuardClassifierTree plants a graded root holding a table and a library. The table
// is written from rows, so a test states the drift it means to prove rather than pasting
// a second copy of the file.
func writeGuardClassifierTree(t *testing.T, root string, rows []guardClassifierRow, libSource string) {
	t.Helper()
	var body strings.Builder
	body.WriteString("package conformance\n\nvar " + guardClassifierTableVar + " = []guardClassifierRow{\n")
	for _, row := range rows {
		body.WriteString(fmt.Sprintf("\t{%q, %t},\n", row.command, row.invokes))
	}
	body.WriteString("}\n")
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardClassifierTableRelPath)), body.String())
	if libSource != "" {
		writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardClassifierLibRelPath)), libSource)
	}
}

// liveGuardClassifierLib is the graded tree's own shared library, which the tests below
// plant unchanged or with one named mutation. Reading it beats pasting a second copy,
// because a paste would keep passing after the real word test changed.
func liveGuardClassifierLib(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(guardClassifierLibRelPath)))
	if err != nil {
		t.Fatalf("read the shared library: %v", err)
	}
	return string(data)
}

func TestCheckGuardClassifierTableAcceptsAgreeingClassifiers(t *testing.T) {
	root := t.TempDir()
	writeGuardClassifierTree(t, root, guardClassifierTable, liveGuardClassifierLib(t))
	if diags := checkGuardClassifierTable(root); len(diags) != 0 {
		t.Fatalf("agreeing classifiers produced diagnostics: %v", diags)
	}
}

// TestCheckGuardClassifierTableDetectsOneSidedDrift covers the drift class the check
// exists to catch: a word test that loses a routine prefix while the table still claims
// the row. The diagnostic must name the row, so a reader repairs one side rather than
// hunting the whole table.
func TestCheckGuardClassifierTableDetectsOneSidedDrift(t *testing.T) {
	root := t.TempDir()
	live := liveGuardClassifierLib(t)
	drifted := strings.Replace(live, "      xargs)", "      xargs-disabled)", 1)
	if drifted == live {
		t.Fatal("the xargs prefix arm was not found; the drift mutation anchors on it")
	}
	writeGuardClassifierTree(t, root, guardClassifierTable, drifted)
	diags := checkGuardClassifierTable(root)
	if !anyContains(diags, `guard classifier row "xargs -- bench help" drifts`) {
		t.Fatalf("a dropped xargs prefix did not produce the drift diagnostic: %v", diags)
	}
}

// TestCheckGuardClassifierTableRedsRatherThanPassingVacuously is the anti-amputation
// proof. Every way the subject can go missing has to red, or a check that parses nothing
// reports green over a classifier nobody graded.
func TestCheckGuardClassifierTableRedsRatherThanPassingVacuously(t *testing.T) {
	live := liveGuardClassifierLib(t)
	for _, test := range []struct {
		name, want string
		plant      func(t *testing.T, root string)
	}{
		{
			name: "absent table",
			want: guardClassifierTableRelPath,
			plant: func(t *testing.T, root string) {
				writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardClassifierLibRelPath)), live)
				if err := os.MkdirAll(filepath.Join(root, "internal", "conformance"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:  "table with no rows",
			want:  "declares no " + guardClassifierTableVar + " rows",
			plant: func(t *testing.T, root string) { writeGuardClassifierTree(t, root, nil, live) },
		},
		{
			name: "renamed table variable",
			want: "no " + guardClassifierTableVar + " composite literal to anchor on",
			plant: func(t *testing.T, root string) {
				writeGuardClassifierTree(t, root, guardClassifierTable, live)
				path := filepath.Join(root, filepath.FromSlash(guardClassifierTableRelPath))
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				writeFixtureFile(t, path, strings.Replace(string(data), guardClassifierTableVar, "guardClassifierTableV2", 1))
			},
		},
		{
			name:  "absent library",
			want:  guardClassifierLibRelPath + " is missing",
			plant: func(t *testing.T, root string) { writeGuardClassifierTree(t, root, guardClassifierTable, "") },
		},
		{
			name: "renamed word test",
			want: "has no " + guardClassifierShellFunc + "() function to anchor on",
			plant: func(t *testing.T, root string) {
				writeGuardClassifierTree(t, root, guardClassifierTable, strings.ReplaceAll(live, guardClassifierShellFunc, "bench_invokes_bench_v2"))
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			test.plant(t, root)
			if diags := checkGuardClassifierTable(root); !anyContains(diags, test.want) {
				t.Fatalf("%s did not red with %q: %v", test.name, test.want, diags)
			}
		})
	}
}

// TestCheckGuardClassifierTableSkipsATreeWithoutTheConformancePackage keeps a consumer
// repository out of the check's reach. Such a tree carries the shared library and no
// table, and refusing it would block every consumer build.
func TestCheckGuardClassifierTableSkipsATreeWithoutTheConformancePackage(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardClassifierLibRelPath)), liveGuardClassifierLib(t))
	if diags := checkGuardClassifierTable(root); len(diags) != 0 {
		t.Fatalf("a tree with no conformance package produced diagnostics: %v", diags)
	}
}
