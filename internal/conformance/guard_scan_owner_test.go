package conformance

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// guardScanOwnerDir holds the production Go source that owns the guard scan, and
// guardScanOwnerFile is the file in it that holds Scan.
const (
	guardScanOwnerDir  = "internal/guards"
	guardScanOwnerFile = guardScanOwnerDir + "/guards.go"
)

// guardScanOwnerDiagPrefix starts each diagnostic of checkGuardScanLoopOwner. The other three
// constants name the shared diagnostic pieces for the check and its tests.
const (
	guardScanOwnerDiagPrefix = "guard scan loop owner "
	guardScanOwnerUnreadable = guardOwnerUnreadable
	guardScanOwnerIrregular  = guardOwnerIrregular
	guardScanOwnerUnparsed   = guardOwnerUnparsed
)

// checkGuardScanLoopOwner keeps the guard scan one synchronous loop. It parses each production
// Go file beneath internal/guards, unused function bodies included, and reports each go
// statement, select statement, channel construction, channel send, and channel receive. It reads
// syntax nodes, so a comment or a string literal never matches. The source walk skips test files
// and testdata directories, so a test goroutine stays legal. A root without internal/guards,
// such as a linked consumer, holds no scan source and gets no diagnostic; any other stat error is
// a diagnostic. Inside that directory, a missing or unreadable guards.go or subdirectory is a
// diagnostic, because an absent owner would otherwise read as a clean pass.
func checkGuardScanLoopOwner(root string) []string {
	dir := filepath.Join(root, filepath.FromSlash(guardScanOwnerDir))
	info, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) || (err == nil && !info.IsDir()) {
		return nil
	}
	if err != nil {
		return []string{guardScanOwnerDiagPrefix + guardScanOwnerDir + guardScanOwnerUnreadable}
	}
	production, diags := guardOwnerSources(root, guardScanOwnerDir, guardScanOwnerDiagPrefix)
	if !slices.Contains(production, guardScanOwnerFile) {
		return append(diags, guardScanOwnerDiagPrefix+guardScanOwnerFile+guardScanOwnerUnreadable)
	}
	return append(diags, guardOwnerParse(root, guardScanOwnerDiagPrefix, production, func(fset *token.FileSet, rel string, file *ast.File) []string {
		var found []string
		ast.Inspect(file, func(node ast.Node) bool {
			if construct := scanPipelineConstruct(node); construct != "" {
				found = append(found, fmt.Sprintf("%s%s:%d: %s", guardScanOwnerDiagPrefix, rel, fset.Position(node.Pos()).Line, construct))
			}
			return true
		})
		return found
	})...)
}

// scanPipelineConstruct names the concurrent pipeline construct that node is, or gives "".
// Channel construction is a make call whose type argument is a channel type; a type argument
// cannot reach a function a file declares, so the call names the builtin.
func scanPipelineConstruct(node ast.Node) string {
	switch n := node.(type) {
	case *ast.GoStmt:
		return "go statement"
	case *ast.SelectStmt:
		return "select statement"
	case *ast.SendStmt:
		return "channel send"
	case *ast.UnaryExpr:
		if n.Op == token.ARROW {
			return "channel receive"
		}
	case *ast.CallExpr:
		if ident, ok := n.Fun.(*ast.Ident); ok && ident.Name == "make" && len(n.Args) > 0 {
			if _, chanType := ast.Unparen(n.Args[0]).(*ast.ChanType); chanType {
				return "channel construction"
			}
		}
	}
	return ""
}

// guardScanOwnerDiags runs the registered package-core-guard binding on root and keeps only the
// scan loop owner diagnostics, so each test grades the subcheck through its registered owner.
func guardScanOwnerDiags(t *testing.T, root string) []string {
	t.Helper()
	return guardOwnerDiags(t, root, guardScanOwnerDiagPrefix)
}

// guardScanPipelineCopy is the retired worker pipeline of Scan under a new name. No caller
// reaches it, so only a check that grades every function body sees it.
const guardScanPipelineCopy = `package guards

import "context"

func scanPipeline(ctx context.Context, root string) ScanResult {
	type enumeration struct {
		candidates []candidate
		err        error
	}
	enumerated := make(chan enumeration, 1)
	go func() {
		candidates, err := enumerateGuards(ctx, root)
		enumerated <- enumeration{candidates: candidates, err: err}
	}()
	var candidates []candidate
	select {
	case result := <-enumerated:
		candidates = result.candidates
	case <-ctx.Done():
		<-enumerated
		return ScanResult{Status: "incomplete"}
	}
	var rows [][]string
	for _, c := range candidates {
		finished := make(chan [][]string, 1)
		go func(candidate candidate) { finished <- inspectGuard(ctx, root, candidate) }(c)
		rows = append(rows, <-finished...)
	}
	return ScanResult{Rows: rows, Status: "complete"}
}
`

// TestGuardScanLoopOwnerBitesOnPrivateCopy plants a renamed, unused copy of the retired
// pipeline beside a direct loop. The registered package-core-guard check must report each
// construct kind in the copy. A goroutine and a channel in a test file, and pipeline text in a
// comment or a string literal, must not get a diagnostic.
func TestGuardScanLoopOwnerBitesOnPrivateCopy(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)), "package guards\n\n// Scan once ran go func() { ch <- rows }() and a select on ctx.Done().\nfunc Scan() string { return \"go func() { <-done }(); select {}\" }\n")
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/pipeline.go")), guardScanPipelineCopy)
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/guards_test.go")), "package guards\n\nfunc waitInTest() {\n\tdone := make(chan struct{})\n\tgo func() { done <- struct{}{} }()\n\tselect {\n\tcase <-done:\n\t}\n}\n")
	diags := guardScanOwnerDiags(t, root)
	for _, construct := range []string{"go statement", "select statement", "channel construction", "channel send", "channel receive"} {
		if !containsDiagnostic(diags, ": "+construct) {
			t.Errorf("want a %s diagnostic on the private pipeline copy, got %v", construct, diags)
		}
	}
	for _, diag := range diags {
		if !strings.Contains(diag, guardScanOwnerDir+"/pipeline.go:") {
			t.Errorf("want diagnostics only on the private pipeline copy, got %q", diag)
		}
	}
}

// TestGuardScanLoopOwnerOptionalSurface grades a linked consumer root. It holds Go source with
// a goroutine outside internal/guards and no internal/guards directory, so the subcheck has no
// scan source to grade and reports nothing.
func TestGuardScanLoopOwnerOptionalSurface(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "internal", "worker", "worker.go"), "package worker\n\nfunc Start(f func()) { go f() }\n")
	if diags := guardScanOwnerDiags(t, root); len(diags) != 0 {
		t.Fatalf("want no scan loop owner diagnostic for a root without %s, got %v", guardScanOwnerDir, diags)
	}
}

// TestGuardScanLoopOwnerMissingOwner amputates guards.go from a root that holds internal/guards.
// An absent file, a directory in its place, and a file that does not parse each get a named
// diagnostic, so a missing owner never reads as a clean pass.
func TestGuardScanLoopOwnerMissingOwner(t *testing.T) {
	missing := guardScanOwnerDiagPrefix + guardScanOwnerFile + guardScanOwnerUnreadable
	for _, tc := range []struct {
		name, want string
		plant      func(t *testing.T, path string)
	}{
		{"absent", missing, func(t *testing.T, path string) {}},
		{"directory", missing, func(t *testing.T, path string) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{"unparseable", guardScanOwnerDiagPrefix + guardScanOwnerUnparsed + guardScanOwnerFile, func(t *testing.T, path string) {
			writeFixtureFile(t, path, "package guards\n\nfunc Scan( {\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/guards_test.go")), "package guards\n")
			tc.plant(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)))
			if diags := guardScanOwnerDiags(t, root); !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags)
			}
		})
	}
}

// TestGuardScanLoopOwnerGradesNestedProductionFiles plants a goroutine in a production file of a
// subdirectory of internal/guards and in a testdata file. Only the nested production file gets a
// diagnostic, because testdata holds fixtures and no production source.
func TestGuardScanLoopOwnerGradesNestedProductionFiles(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)), "package guards\n")
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/inner/inner.go")), "package inner\n\nfunc start(f func()) { go f() }\n")
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/testdata/fixture.go")), "package fixture\n\nfunc start(f func()) { go f() }\n")
	diags := guardScanOwnerDiags(t, root)
	want := guardScanOwnerDiagPrefix + guardScanOwnerDir + "/inner/inner.go:3: go statement"
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("want only %q, got %v", want, diags)
	}
}

// TestGuardScanLoopOwnerRefusesASpecialGoFile plants a FIFO named x.go in internal/guards. The
// check must name it before any open, because an open of a FIFO with no writer blocks the gate.
// The test calls the subcheck itself: a sibling package-core-guard subcheck still opens the FIFO.
func TestGuardScanLoopOwnerRefusesASpecialGoFile(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)), "package guards\n")
	rel := guardScanOwnerDir + "/x.go"
	if err := syscall.Mkfifo(filepath.Join(root, filepath.FromSlash(rel)), 0o644); err != nil {
		capability.Unavailable(t, capability.Fifo, err.Error())
	}
	result := make(chan []string, 1)
	go func() { result <- checkGuardScanLoopOwner(root) }()
	select {
	case diags := <-result:
		if want := guardScanOwnerDiagPrefix + rel + guardScanOwnerIrregular; !slices.Contains(diags, want) {
			t.Fatalf("want %q, got %v", want, diags)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatalf("the check blocked on the FIFO %s", rel)
	}
}

// TestGuardScanLoopOwnerReportsAnUnstatableOwnerDir strips the permissions of the parent of
// internal/guards. A stat error other than not-exist must get a named diagnostic, because only an
// absent directory marks a root with no scan source.
func TestGuardScanLoopOwnerReportsAnUnstatableOwnerDir(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)), "package guards\n")
	parent := filepath.Join(root, "internal")
	// The restore is registered before the strip, so TempDir's removal can descend again.
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	if err := os.Chmod(parent, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(guardScanOwnerDir))); err == nil {
		capability.Capability(t, capability.Privilege, "mode 0o000 directory is still searchable by this user")
	}
	want := guardScanOwnerDiagPrefix + guardScanOwnerDir + guardScanOwnerUnreadable
	if diags := guardScanOwnerDiags(t, root); !slices.Contains(diags, want) {
		t.Fatalf("want %q, got %v", want, diags)
	}
}

// TestGuardScanLoopOwnerReportsAnUnreadableSubdir strips the permissions of a subdirectory of
// internal/guards. The walk cannot list it, so it must get a named diagnostic, because an unlisted
// subdirectory would otherwise hide its production files from the check.
func TestGuardScanLoopOwnerReportsAnUnreadableSubdir(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)), "package guards\n")
	rel := guardScanOwnerDir + "/inner"
	nested := filepath.Join(root, filepath.FromSlash(rel))
	writeFixtureFile(t, filepath.Join(nested, "inner.go"), "package inner\n")
	// The restore is registered before the strip, so TempDir's removal can descend again.
	t.Cleanup(func() { _ = os.Chmod(nested, 0o755) })
	if err := os.Chmod(nested, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
	}
	if _, err := os.ReadDir(nested); err == nil {
		capability.Capability(t, capability.Privilege, "mode 0o000 directory is still readable by this user")
	}
	want := guardScanOwnerDiagPrefix + rel + guardScanOwnerUnreadable
	if diags := guardScanOwnerDiags(t, root); !slices.Contains(diags, want) {
		t.Fatalf("want %q, got %v", want, diags)
	}
}

// TestGuardScanLoopOwnerHoldsOnTheLiveTree runs the registered package-core-guard binding on the
// kit root, which holds the guard scan source, and requires no scan loop owner diagnostic.
func TestGuardScanLoopOwnerHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	root := h.KitRoot
	if !exists(filepath.Join(root, filepath.FromSlash(guardScanOwnerFile))) {
		t.Fatalf("the kit root holds no %s, so the live-tree check would grade nothing", guardScanOwnerFile)
	}
	if diags := guardScanOwnerDiags(t, root); len(diags) != 0 {
		t.Fatalf("the kit guard scan is not one synchronous loop:\n%s", strings.Join(diags, "\n"))
	}
}
