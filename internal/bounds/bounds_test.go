package bounds

import (
	"bytes"
	"context"
	"errors"
	"go/ast"
	"go/constant"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestProductionPolicyValues(t *testing.T) {
	if ProviderTimeout != 10*time.Second || GitRefreshTimeout != 30*time.Second || GuardScanTimeout != 5*time.Second || GateTimeout != 45*time.Minute {
		t.Fatalf("duration policy changed: provider=%s refresh=%s guard=%s gate=%s", ProviderTimeout, GitRefreshTimeout, GuardScanTimeout, GateTimeout)
	}
	if ModelReadLimit != 5<<20 || OutlineFileLimit != 2<<20 || OutlineRowLimit != 200 {
		t.Fatalf("read/output policy changed: model=%d outline_file=%d outline_rows=%d", ModelReadLimit, OutlineFileLimit, OutlineRowLimit)
	}
	if IterationMin != 1 || IterationMax != 100 || MainIterationsDefault != 12 || RefactorIterationsDefault != 4 || MaxWall != 24*time.Hour {
		t.Fatalf("shift policy changed: range=[%d,%d] defaults=%d/%d max_wall=%s", IterationMin, IterationMax, MainIterationsDefault, RefactorIterationsDefault, MaxWall)
	}
}

func TestDeadlineExceedsEveryRegistryDuration(t *testing.T) {
	registry := registryDurations(t)
	if len(registry) == 0 {
		t.Fatal("policy registry yielded no duration bounds to grade")
	}
	for name, bound := range registry {
		if got := TestDeadline(bound); got <= bound {
			t.Errorf("TestDeadline(%s = %s) = %s, want strictly greater than the bound it contains", name, bound, got)
		}
	}
}

// registryDurations type-checks the policy registry's own declarations and returns every
// time.Duration entry by name. It reads the const block instead of restating a sample of
// it, so a duration bound added later is graded here without anyone remembering to add it.
func registryDurations(t *testing.T) map[string]time.Duration {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "bounds.go", nil, 0)
	if err != nil {
		t.Fatalf("parse policy registry: %v", err)
	}
	source := "package bounds\n\nimport \"time\"\n"
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST || !scalarConstBlock(gen) {
			continue
		}
		var block bytes.Buffer
		if err := format.Node(&block, fset, gen); err != nil {
			t.Fatalf("render const block: %v", err)
		}
		source += "\n" + block.String() + "\n"
	}
	synthetic := token.NewFileSet()
	parsed, err := parser.ParseFile(synthetic, "registry.go", source, 0)
	if err != nil {
		t.Fatalf("parse extracted registry: %v", err)
	}
	config := types.Config{Importer: importer.ForCompiler(synthetic, "source", nil)}
	pkg, err := config.Check("bounds", synthetic, []*ast.File{parsed}, nil)
	if err != nil {
		t.Fatalf("type-check extracted registry: %v", err)
	}
	durations := map[string]time.Duration{}
	for _, name := range pkg.Scope().Names() {
		entry, ok := pkg.Scope().Lookup(name).(*types.Const)
		if !ok || entry.Type().String() != "time.Duration" {
			continue
		}
		value, exact := constant.Int64Val(entry.Val())
		if !exact {
			t.Fatalf("registry bound %s is not an exact integer duration", name)
		}
		durations[name] = time.Duration(value)
	}
	return durations
}

// scalarConstBlock admits the policy registry and rejects the package's string-enum const
// blocks: a bound is declared with no type or a predeclared one, an enum with its own.
func scalarConstBlock(gen *ast.GenDecl) bool {
	for _, item := range gen.Specs {
		spec, ok := item.(*ast.ValueSpec)
		if !ok {
			return false
		}
		if spec.Type == nil {
			continue
		}
		named, ok := spec.Type.(*ast.Ident)
		if !ok || types.Universe.Lookup(named.Name) == nil {
			return false
		}
	}
	return true
}

func TestRunClassifiesProcessOutcomes(t *testing.T) {
	tests := []struct {
		name   string
		parent func() (context.Context, context.CancelFunc)
		limit  time.Duration
		cmd    func(context.Context) *exec.Cmd
		want   ProcessStatus
		exit   int
	}{
		{name: "complete", parent: liveContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd { return exec.CommandContext(ctx, "sh", "-c", "printf ok") }, want: ProcessComplete},
		{name: "timeout", parent: liveContext, limit: 20 * time.Millisecond, cmd: func(ctx context.Context) *exec.Cmd { return exec.CommandContext(ctx, "sh", "-c", "sleep 5") }, want: ProcessTimeout},
		{name: "parent cancellation", parent: canceledContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd { return exec.CommandContext(ctx, "sh", "-c", "sleep 5") }, want: ProcessCanceled},
		{name: "nonzero exit", parent: liveContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd {
			return exec.CommandContext(ctx, "sh", "-c", "printf bad >&2; exit 23")
		}, want: ProcessExit, exit: 23},
		{name: "start failure", parent: liveContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd {
			return exec.CommandContext(ctx, "/definitely/missing/bench-command")
		}, want: ProcessStart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.parent()
			defer cancel()
			got := Run(ctx, tt.limit, tt.cmd(ctx))
			if got.Status != tt.want || got.Exit != tt.exit {
				t.Fatalf("Run status/exit = %q/%d, want %q/%d (err=%v output=%q)", got.Status, got.Exit, tt.want, tt.exit, got.Err, got.Output)
			}
		})
	}
}

func TestReadClassifiesExactLimitAndLimitPlusOne(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		want ReadStatus
	}{
		{name: "exact", body: "12345", want: ReadComplete},
		{name: "plus one", body: "123456", want: ReadOversized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Read(strings.NewReader(tt.body), 5)
			if got.Status != tt.want {
				t.Fatalf("Read status = %q, want %q", got.Status, tt.want)
			}
			if got.Status == ReadComplete && string(got.Data) != tt.body {
				t.Fatalf("Read data = %q, want %q", got.Data, tt.body)
			}
		})
	}
	failing := Read(errorReader{}, 5)
	if failing.Status != ReadFailed || failing.Err == nil {
		t.Fatalf("failing reader = %#v, want failed with error", failing)
	}
}

func liveContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func canceledContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx, cancel
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }
