package worktree

// This file owns the verb call census. Every worktree test runs a verb through the
// verb runner, so the census reports each reference to a verb form or to a
// runner-private function outside the two runner files. The census parses; it never
// builds, so it decides by bare identifier, the same way as the serial helper edge.

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// verbRunnerFiles are the two test files that may name a verb form or a
// runner-private function. A runner file absent from dir exempts nothing.
var verbRunnerFiles = map[string]bool{"verb_runner_test.go": true, "verb_runner_check_test.go": true}

// verbFormNames returns the verb entries and the joins forms of the non-test files
// in dir. An entry is derived from its signature, and a joins form from the entry
// body, so a new or renamed verb needs no list edit.
func verbFormNames(dir string) (map[string]bool, error) {
	files, err := parseSourceFiles(dir)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, file := range files {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.FuncDecl)
			if !ok || decl.Recv != nil || decl.Body == nil || !decl.Name.IsExported() || !hasArgsParam(decl.Type.Params) {
				continue
			}
			names[decl.Name.Name] = true
			ast.Inspect(decl.Body, func(node ast.Node) bool {
				if call, ok := node.(*ast.CallExpr); ok && len(call.Args) > 0 {
					if first, ok := call.Args[0].(*ast.CallExpr); ok && calleeName(first) == "defaultJoins" && calleeName(call) != "" {
						names[calleeName(call)] = true
					}
				}
				return true
			})
		}
	}
	return names, nil
}

// hasArgsParam reports whether params holds a parameter named args of type []string.
func hasArgsParam(params *ast.FieldList) bool {
	for _, field := range params.List {
		array, ok := field.Type.(*ast.ArrayType)
		if !ok || array.Len != nil {
			continue
		}
		if elem, ok := array.Elt.(*ast.Ident); !ok || elem.Name != "string" {
			continue
		}
		for _, name := range field.Names {
			if name.Name == "args" {
				return true
			}
		}
	}
	return false
}

// addRunnerPrivateNames adds each top-level function of a runner file whose last
// result is an error. Such a function is a core reader or the call check, and a
// test outside the runner files reaches it only through a runner function that
// fails the test on that error.
func addRunnerPrivateNames(file *ast.File, names map[string]bool) {
	for _, node := range file.Decls {
		decl, ok := node.(*ast.FuncDecl)
		if !ok || decl.Recv != nil || decl.Type.Results == nil || len(decl.Type.Results.List) == 0 {
			continue
		}
		results := decl.Type.Results.List
		if last, ok := results[len(results)-1].Type.(*ast.Ident); ok && last.Name == "error" {
			names[decl.Name.Name] = true
		}
	}
}

// verbCallCensus reports each identifier in a test file of dir, other than the runner
// files, that names a verb form or a runner-private function. Each report names the
// file, the line, the enclosing declaration, and the name, and the reports are sorted.
func verbCallCensus(dir string) ([]string, error) {
	refused, err := verbFormNames(dir)
	if err != nil {
		return nil, err
	}
	files, fset, names, err := parseTestFiles(dir)
	if err != nil {
		return nil, err
	}
	for i, file := range files {
		if verbRunnerFiles[names[i]] {
			addRunnerPrivateNames(file, refused)
		}
	}
	var reports []string
	report := func(enclosing string, root ast.Node) {
		var walk func(ast.Node) bool
		walk = func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.SelectorExpr:
				// A field or a method named like a verb form is not the form.
				ast.Inspect(node.X, walk)
				return false
			case *ast.Ident:
				if refused[node.Name] {
					position := fset.Position(node.Pos())
					reports = append(reports, fmt.Sprintf("%s:%d: %s names %s", filepath.Base(position.Filename), position.Line, enclosing, node.Name))
				}
			}
			return true
		}
		ast.Inspect(root, walk)
	}
	for i, file := range files {
		if verbRunnerFiles[names[i]] {
			continue
		}
		for _, node := range file.Decls {
			switch decl := node.(type) {
			case *ast.FuncDecl:
				if decl.Body != nil {
					report(decl.Name.Name, decl.Body)
				}
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					report(specName(spec), spec)
				}
			}
		}
	}
	sort.Strings(reports)
	return reports, nil
}

// specName names the declaration a package-level spec makes, the first name of a
// value spec with several names standing for all of them.
func specName(spec ast.Spec) string {
	switch spec := spec.(type) {
	case *ast.ValueSpec:
		return spec.Names[0].Name
	case *ast.TypeSpec:
		return spec.Name.Name
	}
	return token.IMPORT.String()
}

// --- verb call census unit tests over synthetic file sets ---

// syntheticVerbSource is the non-test file of a synthetic file set. RunCommand is a
// verb entry, runWith is its joins form, and ExportedHelper is neither.
const syntheticVerbSource = `package worktree

func RunCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return runWith(defaultJoins(), root, home, args, stdout, stderr)
}

func runWith(j joins, root, home string, args []string, stdout, stderr io.Writer) int {
	return 0
}

func ExportedHelper(root string) string {
	return root
}
`

// syntheticCallerFile is a synthetic test file that holds body as the whole body of
// its one top-level test. The first line of body is line 6 of the file.
func syntheticCallerFile(body string) string {
	return `package worktree

import "testing"

func TestCaller(t *testing.T) {
` + body + `}
`
}

func verbCensusOf(t *testing.T, files map[string]string) []string {
	t.Helper()
	reports, err := verbCallCensus(plantTestFiles(t, files))
	if err != nil {
		t.Fatalf("verb call census: %v", err)
	}
	return reports
}

func wantVerbReports(t *testing.T, reports []string, want ...string) {
	t.Helper()
	if strings.Join(reports, "\n") != strings.Join(want, "\n") {
		t.Fatalf("verb call census = %q, want exactly %q", reports, want)
	}
}

func TestVerbCallCensusReportsAnEntryCall(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verbs.go":       syntheticVerbSource,
		"caller_test.go": syntheticCallerFile("\tRunCommand(\"\", \"\", nil, nil, nil)\n"),
	})
	wantVerbReports(t, reports, "caller_test.go:6: TestCaller names RunCommand")
}

func TestVerbCallCensusReportsAJoinsFormCall(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verbs.go":       syntheticVerbSource,
		"caller_test.go": syntheticCallerFile("\trunWith(defaultJoins(), \"\", \"\", nil, nil, nil)\n"),
	})
	wantVerbReports(t, reports, "caller_test.go:6: TestCaller names runWith")
}

func TestVerbCallCensusReportsAnEntryUsedAsAValue(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verbs.go":       syntheticVerbSource,
		"caller_test.go": syntheticCallerFile("\trun := RunCommand\n\t_ = run\n"),
	})
	wantVerbReports(t, reports, "caller_test.go:6: TestCaller names RunCommand")
}

// TestVerbCallCensusReportsACallInsideASubtest proves a closure is no shelter. The
// call sits in a subtest closure, and a second call sits in a method of a test type.
func TestVerbCallCensusReportsACallInsideASubtest(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verbs.go": syntheticVerbSource,
		"caller_test.go": syntheticCallerFile(`	t.Run("x", func(t *testing.T) {
		RunCommand("", "", nil, nil, nil)
	})
}

type caller struct{}

func (caller) run() {
	runWith(defaultJoins(), "", "", nil, nil, nil)
`),
	})
	wantVerbReports(t, reports,
		"caller_test.go:14: run names runWith",
		"caller_test.go:7: TestCaller names RunCommand")
}

// TestVerbCallCensusReportsAPackageLevelReference proves a table outside a function
// is no shelter: the report names the declared variable as the enclosing declaration.
func TestVerbCallCensusReportsAPackageLevelReference(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verbs.go": syntheticVerbSource,
		"table_test.go": `package worktree

var verbTable = []any{RunCommand, runWith}
`,
	})
	wantVerbReports(t, reports,
		"table_test.go:3: verbTable names RunCommand",
		"table_test.go:3: verbTable names runWith")
}

// TestVerbCallCensusAllowsTheRunnerFiles proves each runner file can name every verb
// form and every runner-private function. Both runner files hold the same text: the
// census only parses, so the duplicate declaration does no harm.
func TestVerbCallCensusAllowsTheRunnerFiles(t *testing.T) {
	t.Parallel()
	runnerFile := `package worktree

func checkVerbCall(call int) error {
	_ = runWith(defaultJoins(), "", "", nil, nil, nil)
	return nil
}

var runnerTable = []any{RunCommand, checkVerbCall}
`
	reports := verbCensusOf(t, map[string]string{
		"verbs.go":                  syntheticVerbSource,
		"verb_runner_test.go":       runnerFile,
		"verb_runner_check_test.go": runnerFile,
	})
	wantVerbReports(t, reports)
}

// TestVerbCallCensusDerivesEntriesFromTheSignature proves an exported function with
// an args parameter of type []string is an entry whatever its name.
func TestVerbCallCensusDerivesEntriesFromTheSignature(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"fresh.go": `package worktree

func Fresh(root string, args []string) int {
	return len(args)
}
`,
		"caller_test.go": syntheticCallerFile("\t_ = Fresh(\"\", nil)\n"),
	})
	wantVerbReports(t, reports, "caller_test.go:6: TestCaller names Fresh")
}

// TestVerbCallCensusIgnoresAnExportedHelper proves an exported function is not an
// entry without its args parameter, and that a method or field of the same name
// as an entry is not the entry.
func TestVerbCallCensusIgnoresAnExportedHelper(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verbs.go":       syntheticVerbSource,
		"caller_test.go": syntheticCallerFile("\t_ = ExportedHelper(\"\")\n\t_ = t.RunCommand\n"),
	})
	wantVerbReports(t, reports)
}

// TestVerbCallCensusDerivesTheJoinsForm proves the joins form is the function an
// entry body calls with defaultJoins() first, and that another callee of the entry
// is not a joins form.
func TestVerbCallCensusDerivesTheJoinsForm(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"fresh.go": `package worktree

func FreshCommand(args []string) int {
	return freshAt(defaultJoins(), parseFresh(args))
}
`,
		"caller_test.go": syntheticCallerFile("\t_ = parseFresh(nil)\n\t_ = freshAt(defaultJoins(), nil)\n"),
	})
	wantVerbReports(t, reports, "caller_test.go:7: TestCaller names freshAt")
}

// TestVerbCallCensusReportsACoreReaderCall proves a top-level runner-file function
// whose last result is an error is runner-private, and that another runner-file
// function is not.
func TestVerbCallCensusReportsACoreReaderCall(t *testing.T) {
	t.Parallel()
	reports := verbCensusOf(t, map[string]string{
		"verb_runner_test.go": `package worktree

func readFresh(stdout string) (string, error) {
	return stdout, nil
}

func formatFresh(stdout string) string {
	return stdout
}
`,
		"caller_test.go": syntheticCallerFile("\t_, _ = readFresh(\"\")\n\t_ = formatFresh(\"\")\n"),
	})
	wantVerbReports(t, reports, "caller_test.go:6: TestCaller names readFresh")
}

// --- live-tree pin ---

// TestVerbCallCensusOnTheLiveTree proves no test outside the runner files names a
// verb form or a runner-private function.
func TestVerbCallCensusOnTheLiveTree(t *testing.T) {
	t.Parallel()
	reports, err := verbCallCensus(".")
	if err != nil {
		t.Fatalf("verb call census the package: %v", err)
	}
	if len(reports) != 0 {
		t.Fatalf("the verb call census reports %d direct references:\n%s", len(reports), strings.Join(reports, "\n"))
	}
}
