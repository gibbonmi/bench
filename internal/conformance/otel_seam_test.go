package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// The seam registry in internal/otelrecord is the one source for the instrumented set.
// This check reads it and grades each named symbol's own source: a registered seam whose
// symbol opens no span is an advertisement without instrumentation, and it would
// otherwise stay silently green. The registry is the advertisement and the Go source is
// the evidence, so the two halves stay independent.
//
// The registry is a kit-owned policy, so the check rides the kit-compliance subject like
// the footprint check beside it. A subject that carries no source for a registered
// package is not the kit, and it carries nothing to grade.

const (
	missingSeamSymbolMessage = "registered seam names a symbol the package does not declare"
	unstartedSeamMessage     = "registered seam symbol starts no span"
	unparsedSeamPackage      = "registered seam package does not parse"
)

// checkOtelSeamSpans grades the kit's own instrumentation against the seam registry.
func checkOtelSeamSpans(kitRoot string) []string {
	return otelSeamDiags(kitRoot, otelrecord.Registry)
}

func otelSeamDiags(root string, entries []otelrecord.SeamEntry) []string {
	var diags []string
	for _, entry := range entries {
		dir := filepath.Join(root, filepath.FromSlash(entry.Package))
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			continue
		}
		body, found, err := seamFunctionBody(dir, entry.Function)
		switch {
		case err != nil:
			diags = append(diags, unparsedSeamPackage+": "+entry.Seam+" ("+entry.Package+"): "+err.Error())
		case !found:
			diags = append(diags, missingSeamSymbolMessage+": "+entry.Seam+" ("+entry.Package+"."+entry.Function+")")
		case !startsSpan(body):
			diags = append(diags, unstartedSeamMessage+": "+entry.Seam+" ("+entry.Package+"."+entry.Function+")")
		}
	}
	return diags
}

// seamFunctionBody parses the package directory's shipped Go files and returns the named
// function's body. Test files are excluded, because a seam instrumented only under test
// is not instrumented.
func seamFunctionBody(dir, function string) (*ast.BlockStmt, bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, false, err
	}
	fset := token.NewFileSet()
	for _, item := range entries {
		name := item.Name()
		if item.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			return nil, false, err
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Name.Name == function && fn.Body != nil {
				return fn.Body, true, nil
			}
		}
	}
	return nil, false, nil
}

// spanOpeners are the calls that open a Bench span: the OpenTelemetry tracer's own Start
// method, and the otelrecord primitives that own the open-and-close protocol above it. A
// seam that reaches one of these has instrumentation a reader of the source can see.
var spanOpeners = map[string]bool{"Start": true, "Begin": true, "BeginIn": true}

// startsSpan reports whether the body opens a span. The selector name is the evidence,
// so the check needs no build and no type information.
func startsSpan(body *ast.BlockStmt) bool {
	started := false
	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if selector, ok := call.Fun.(*ast.SelectorExpr); ok && spanOpeners[selector.Sel.Name] {
			started = true
			return false
		}
		return true
	})
	return started
}

const otelSeamFixtureInstrumented = `package sample

func beginSampleSpan(ctx context.Context) {
	_, span := tracerFrom(ctx).Start(ctx, "sample")
	span.End()
}

func beginPrimitiveSpan(ctx context.Context) {
	_, span, finish := otelrecord.Begin("", "/root", "sample.primitive")
	span.End()
	finish()
}

func beginQuietSpan(ctx context.Context) {
	log("sample")
}
`

func TestOtelSeamCheckBites(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "internal", "sample")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sample.go"), []byte(otelSeamFixtureInstrumented), 0o644); err != nil {
		t.Fatal(err)
	}
	instrumented := otelrecord.SeamEntry{Seam: "sample", Package: "internal/sample", Function: "beginSampleSpan"}

	if diags := otelSeamDiags(root, []otelrecord.SeamEntry{instrumented}); len(diags) != 0 {
		t.Fatalf("the instrumented fixture is not clean: %v", diags)
	}

	primitive := otelrecord.SeamEntry{Seam: "sample.primitive", Package: "internal/sample", Function: "beginPrimitiveSpan"}
	if diags := otelSeamDiags(root, []otelrecord.SeamEntry{primitive}); len(diags) != 0 {
		t.Fatalf("a seam opened through the otelrecord primitive reads uninstrumented: %v", diags)
	}

	quiet := otelrecord.SeamEntry{Seam: "sample.quiet", Package: "internal/sample", Function: "beginQuietSpan"}
	diags := otelSeamDiags(root, []otelrecord.SeamEntry{instrumented, quiet})
	if len(diags) != 1 {
		t.Fatalf("want one diagnostic for the span-less seam, got %v", diags)
	}
	if !strings.Contains(diags[0], unstartedSeamMessage) || !strings.Contains(diags[0], "sample.quiet") {
		t.Errorf("the diagnostic does not name the span-less seam: %q", diags[0])
	}

	absent := otelrecord.SeamEntry{Seam: "sample.absent", Package: "internal/sample", Function: "beginAbsentSpan"}
	if diags := otelSeamDiags(root, []otelrecord.SeamEntry{absent}); len(diags) != 1 ||
		!strings.Contains(diags[0], missingSeamSymbolMessage) {
		t.Errorf("want a missing-symbol diagnostic, got %v", diags)
	}

	if diags := otelSeamDiags(root, []otelrecord.SeamEntry{{Seam: "gone", Package: "internal/gone", Function: "f"}}); len(diags) != 0 {
		t.Errorf("a subject without the registered package carries nothing to grade, got %v", diags)
	}
}

func TestOtelSeamRegistryIsNotEmpty(t *testing.T) {
	if len(otelrecord.Registry) == 0 {
		t.Fatal("the seam registry is empty, so the check grades nothing")
	}
}
