package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// gitAdapterOwner is the one package allowed to spell a Git administration flag.
const gitAdapterOwner = "internal/git"

// gitAdapterHarness is the fixture harness that stubs git. It sits under the skipped
// tree because it cannot import the adapter without a cycle.
const gitAdapterHarness = "internal/gittest"

// gitPlumbingFlags are the four Git administration flags the adapter alone may spell.
var gitPlumbingFlags = []string{"--git-dir", "--absolute-git-dir", "--git-path", "--git-common-dir"}

// checkGitPlumbingOwner grades the single-source rule for the two repository facts the
// Git adapter owns: the checkout administration directory, and a named file inside it. A
// production function outside internal/git that spells one of the four flags inside a
// `rev-parse` call re-derives that spelling beside the owner, and the two copies drift on
// the failure posture. Test files are exempt, so a fixture can plant a flag. The harness
// and the two guards' option tables are exempt: the harness cannot import the adapter
// without a cycle, and the guards parse Git's global options with no `rev-parse` call.
//
// The unit is the call expression, not the file, so a file that spells the flag once
// inside a rev-parse call and once inside an unrelated call is graded on each occurrence
// independently.
func checkGitPlumbingOwner(root string) []string {
	var diags []string
	for _, top := range []string{"cmd", "internal"} {
		base := filepath.Join(root, top)
		_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel := slashRel(root, path)
			if strings.HasPrefix(rel, gitAdapterOwner+"/") || strings.HasPrefix(rel, gitAdapterHarness+"/") {
				return nil
			}
			diags = append(diags, gitPlumbingDerivations(path, rel)...)
			return nil
		})
	}
	diags = append(diags, gitPlumbingRootFiles(root)...)
	return uniqueSorted(diags)
}

// gitPlumbingRootFiles grades the non-test Go files directly at the module root
// (depth one), alongside cmd/ and internal/. Every other directory stays skipped —
// the canary trees under tests/ are fixtures, not production code the rule owns.
func gitPlumbingRootFiles(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var diags []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(root, entry.Name())
		diags = append(diags, gitPlumbingDerivations(path, slashRel(root, path))...)
	}
	return diags
}

func gitPlumbingDerivations(path, rel string) []string {
	if readIfExists(path) == "" {
		return nil
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for git-plumbing ownership: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || !callHoldsRevParse(call) {
			return true
		}
		for _, arg := range call.Args {
			if flag, matched := gitPlumbingFlagLiteral(arg); matched {
				diags = append(diags, fmt.Sprintf("%s spells the Git administration flag %s outside internal/git", rel, flag))
			}
		}
		return true
	})
	return diags
}

// callHoldsRevParse reports whether one of call's arguments is the string literal
// "rev-parse". A flag literal outside such a call passes, because the two guards parse
// Git's global options there with no rev-parse call.
func callHoldsRevParse(call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		if value, ok := stringLiteralValue(arg); ok && value == "rev-parse" {
			return true
		}
	}
	return false
}

// gitPlumbingFlagLiteral reports whether arg is a string literal equal to one of the
// four Git administration flags. A flag embedded inside a longer literal does not match.
func gitPlumbingFlagLiteral(arg ast.Expr) (string, bool) {
	value, ok := stringLiteralValue(arg)
	if !ok {
		return "", false
	}
	for _, flag := range gitPlumbingFlags {
		if value == flag {
			return flag, true
		}
	}
	return "", false
}

func stringLiteralValue(arg ast.Expr) (string, bool) {
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return value, true
}

// gitPlumbingProductionSnippet renders one production Go file whose call embeds flag
// inside a rev-parse invocation, matching the shape a migration ticket left behind.
func gitPlumbingProductionSnippet(flag string) string {
	return fmt.Sprintf(`package example

import "os/exec"

func adminDir(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", %q).Output()
	return string(out), err
}
`, flag)
}

// gitPlumbingFlagsIndependentExpectation is this test's own list of the four Git
// administration flags, kept independent of the production gitPlumbingFlags variable
// under AGENTS.md's one independent-expectation exception. An expectation that read the
// production list instead would grade the check against itself: a mutation that drops
// one flag from gitPlumbingFlags would silently narrow both the check and this test's
// range together, and the omission would never turn the gate red. The mutation this
// independence catches is exactly that drop — removing "--git-path" from
// gitPlumbingFlags leaves TestGitPlumbingOwnerRedsARetypedFlag/--git-path red, because
// this list still ranges over the flag the mutated check no longer recognizes.
var gitPlumbingFlagsIndependentExpectation = []string{"--git-dir", "--absolute-git-dir", "--git-path", "--git-common-dir"}

// TestGitPlumbingOwnerRedsARetypedFlag is GR27. Each of the four flags, planted inside a
// rev-parse call in a non-test Go file outside internal/git, reds with one diagnostic
// naming the file and the flag. A check hard-coded to one flag would pass the other
// three cases. The range is this test's own independent list, not the production
// gitPlumbingFlags variable, so a flag the production list drops still reds here.
func TestGitPlumbingOwnerRedsARetypedFlag(t *testing.T) {
	for _, flag := range gitPlumbingFlagsIndependentExpectation {
		t.Run(flag, func(t *testing.T) {
			root := throwawayRoot{
				files: map[string]string{
					"internal/example/example.go": gitPlumbingProductionSnippet(flag),
				},
			}.build(t)
			diags := checkGitPlumbingOwner(root)
			if len(diags) != 1 {
				t.Fatalf("diagnostics = %v, want exactly one", diags)
			}
			if !strings.Contains(diags[0], "internal/example/example.go") || !strings.Contains(diags[0], flag) {
				t.Fatalf("diagnostic = %q, want it to name the file and the flag %s", diags[0], flag)
			}
		})
	}

	t.Run("root-level file", func(t *testing.T) {
		root := throwawayRoot{
			files: map[string]string{
				"payload.go": gitPlumbingProductionSnippet("--git-dir"),
			},
		}.build(t)
		diags := checkGitPlumbingOwner(root)
		if len(diags) != 1 {
			t.Fatalf("diagnostics = %v, want exactly one", diags)
		}
		if !strings.Contains(diags[0], "payload.go") || !strings.Contains(diags[0], "--git-dir") {
			t.Fatalf("diagnostic = %q, want it to name the file and the flag --git-dir", diags[0])
		}
	})
}

// TestGitPlumbingOwnerSkipsTestsAndTheAdapter is GR28. A rev-parse call with
// --git-common-dir in a _test.go file, in a file under internal/git, and in a file under
// internal/gittest, produces no diagnostic. A check that grades every Go file would red
// the adapter or the harness.
func TestGitPlumbingOwnerSkipsTestsAndTheAdapter(t *testing.T) {
	snippet := gitPlumbingProductionSnippet("--git-common-dir")
	root := throwawayRoot{
		files: map[string]string{
			"internal/example/example_test.go":     snippet,
			"internal/git/worktree_admin_extra.go": snippet,
			"internal/gittest/gittest_extra.go":    snippet,
		},
	}.build(t)
	if diags := checkGitPlumbingOwner(root); len(diags) != 0 {
		t.Fatalf("diagnostics = %v, want none", diags)
	}
}

// TestGitPlumbingOwnerToleratesEmbeddedFlag is GR29. A flag that sits inside a longer
// literal, such as the doctor's shell snippet, passes. A substring check would red it.
func TestGitPlumbingOwnerToleratesEmbeddedFlag(t *testing.T) {
	root := throwawayRoot{
		files: map[string]string{
			"internal/example/example.go": `package example

func snippet() string {
	return "run: git rev-parse --git-dir to find the admin directory"
}
`,
		},
	}.build(t)
	if diags := checkGitPlumbingOwner(root); len(diags) != 0 {
		t.Fatalf("diagnostics = %v, want none", diags)
	}
}

// TestGitPlumbingOwnerIgnoresCallsWithoutRevParse is GR37. A flag literal in a map
// literal, and a flag literal in a call whose arguments hold no "rev-parse" literal,
// produce no diagnostic. A bare-literal check would red the guard option tables.
func TestGitPlumbingOwnerIgnoresCallsWithoutRevParse(t *testing.T) {
	root := throwawayRoot{
		files: map[string]string{
			"internal/example/example.go": `package example

import "os/exec"

var gitGlobalFlags = map[string]bool{
	"--git-dir": true,
}

func status(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "--git-dir", "status").Output()
	return string(out), err
}
`,
		},
	}.build(t)
	if diags := checkGitPlumbingOwner(root); len(diags) != 0 {
		t.Fatalf("diagnostics = %v, want none", diags)
	}
}
