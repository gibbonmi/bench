package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// defaultBranchFuncRe matches a declaration or call of a DefaultBranch function. The
// trailing `\(` scopes the sweep to the function. RepoFacts.DefaultBranch is the
// surviving struct field. A selector read, like `gf.DefaultBranch`, carries no paren. A
// sweep on the bare identifier would report every honest reader of the field.
var defaultBranchFuncRe = regexp.MustCompile(`\bDefaultBranch\b\s*\(`)

// checkDefaultBranchSingleSource asserts that no shipped Go source declares or calls a
// DefaultBranch function. git.ResolvedDefault owns the default-branch fact and reports
// ok=false when nothing resolves. A helper that answers with a name instead lets one
// caller fabricate a default the repository does not have. Test files are exempt. A
// reintroduced helper is caught at its declaration. Test identifiers legitimately end in
// DefaultBranch.
func checkDefaultBranchSingleSource(root string) []string {
	var diags []string
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		switch {
		case err != nil:
			return nil
		case d.IsDir():
			switch d.Name() {
			case ".git", "node_modules", "tests", "dist":
				return filepath.SkipDir
			}
			return nil
		case !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go"):
			return nil
		case defaultBranchFuncRe.MatchString(readIfExists(path)):
			diags = append(diags, fmt.Sprintf("%s declares or calls a DefaultBranch function — git.ResolvedDefault is the single owner of the default-branch fact", slashRel(root, path)))
		}
		return nil
	})
	sort.Strings(diags)
	return diags
}

// TestDefaultBranchSweepScopesToTheFunction is the recorded scoping proof. The sweep must
// fire on a reintroduced function. The sweep must stay silent on RepoFacts.DefaultBranch,
// the field every renderer still reads.
func TestDefaultBranchSweepScopesToTheFunction(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("render.go", "package p\n\nfunc render(f Facts) string { return f.DefaultBranch }\n")
	if diags := checkDefaultBranchSingleSource(root); len(diags) != 0 {
		t.Fatalf("struct-field read reported: %v", diags)
	}

	write("legacy.go", "package p\n\nfunc DefaultBranch(root string) string { return \"main\" }\n")
	diags := checkDefaultBranchSingleSource(root)
	if len(diags) != 1 || !strings.Contains(diags[0], "legacy.go declares or calls a DefaultBranch function") {
		t.Fatalf("reintroduced function: want one legacy.go diagnostic, got %v", diags)
	}
}
