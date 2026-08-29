package consumers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// otherPkg is the Go side of the non-Go fixture: one declaration whose name is not the
// queried one, so Go resolution honestly finds nothing and the sweep decides the answer.
func otherPkg(root string) []fixturePkg {
	return []fixturePkg{{path: "example.com/target", files: map[string]string{
		root + "/target/target.go": "package target\n\nfunc Other() {}\n"}}}
}

// nonGoRepo commits one TypeScript file declaring name and makes the repository the
// process working directory, so the sweep reads a tracked file rather than a temp path.
func nonGoRepo(t *testing.T, rel, content string) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture file: %v", err)
	}
	if _, err := git.Output("-C", root, "add", rel); err != nil {
		t.Fatalf("track fixture file: %v", err)
	}
	if _, err := git.Output("-C", root, "commit", "-q", "-m", "one"); err != nil {
		t.Fatalf("commit fixture file: %v", err)
	}
	t.Chdir(root)
	return root
}

// CS7 (story 6): a name the Go resolver cannot see, but a tracked non-Go file declares,
// refuses with the language named. An empty table here would state a false "no callers".
func TestNonGoDeclarationRefusesNamingTheLanguage(t *testing.T) {
	nonGoRepo(t, "web/app.ts", "export function Symbol() {}\n")
	stubLoad(t, otherPkg)
	out, code := run(t, "Symbol")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: ") {
		t.Fatalf("stdout = %q, want a structured error line", out)
	}
	if !strings.Contains(out, "TypeScript") {
		t.Fatalf("stdout = %q, want the language named", out)
	}
	if !strings.Contains(out, "web/app.ts") {
		t.Fatalf("stdout = %q, want the declaring file named", out)
	}
	if strings.Contains(out, "consumers[") || strings.Contains(out, "citation[") {
		t.Fatalf("stdout = %q, want no result block and no citation row", out)
	}
}

// A name no file declares in any language keeps the plain unresolved refusal, so the
// sweep adds a class of answer rather than replacing the existing one.
func TestUnknownNameKeepsThePlainUnresolvedRefusal(t *testing.T) {
	nonGoRepo(t, "web/app.ts", "export function Symbol() {}\n")
	stubLoad(t, otherPkg)
	out, code := run(t, "Absent")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; out=%q", code, out)
	}
	if !strings.Contains(out, "no declaration named") {
		t.Fatalf("stdout = %q, want the unresolved refusal", out)
	}
}

// CS22 support: the help text states that the three refusals exist and carry no citation.
func TestHelpNamesTheThreeRefusals(t *testing.T) {
	out, code := run(t, "--help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "refuse at exit 1") {
		t.Fatalf("help = %q, want the refusal line", out)
	}
}
