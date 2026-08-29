package consumers

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
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

// A repository the Go loader cannot load still refuses with the language
// named, because a tree with no go.mod says nothing about the query. This test drives the
// real loader, so the sweep is what decides the answer.
func TestUnloadableTreeRefusesWithTheNonGoLanguage(t *testing.T) {
	nonGoRepo(t, "web/app.ts", "export function Symbol() {}\n")
	out, code := run(t, "Symbol")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; out=%q", code, out)
	}
	if !strings.Contains(out, "TypeScript") {
		t.Fatalf("stdout = %q, want the language named", out)
	}
	if !strings.Contains(out, "web/app.ts") {
		t.Fatalf("stdout = %q, want the declaring file named", out)
	}
}

// The sweep reads a candidate the way `bench outline` does. A tracked path replaced on
// disk by a symlink is not followed, and one replaced by a FIFO does not block, so the
// command still returns its refusal. Both names sort before the declaring file, so the
// sweep reaches them.
func TestSweepSkipsSymlinkAndFifoCandidates(t *testing.T) {
	root := nonGoRepo(t, "web/app.ts", "export function Symbol() {}\n")
	fifo := filepath.Join(root, "web", "aaa-pipe.txt")
	link := filepath.Join(root, "web", "aab-link.txt")
	for _, path := range []string{fifo, link} {
		if err := os.WriteFile(path, []byte("placeholder\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	if _, err := git.Output("-C", root, "add", "-A"); err != nil {
		t.Fatalf("track the candidates: %v", err)
	}
	if _, err := git.Output("-C", root, "commit", "-q", "-m", "candidates"); err != nil {
		t.Fatalf("commit the candidates: %v", err)
	}
	if err := os.Remove(fifo); err != nil {
		t.Fatalf("remove the placeholder: %v", err)
	}
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		t.Fatalf("mkfifo: %v", err)
	}
	if err := os.Remove(link); err != nil {
		t.Fatalf("remove the placeholder: %v", err)
	}
	if err := os.Symlink(fifo, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	stubLoad(t, otherPkg)

	type answer struct {
		out  string
		code int
	}
	done := make(chan answer, 1)
	go func() {
		out, code := CommandWithVersion(testVersion)([]string{"Symbol"})
		done <- answer{out, code}
	}()
	select {
	case got := <-done:
		if got.code != 1 || !strings.Contains(got.out, "TypeScript") {
			t.Fatalf("stdout = %q exit %d, want the language refusal at 1", got.out, got.code)
		}
	case <-time.After(bounds.TestDeadline(bounds.GuardScanTimeout)):
		t.Fatal("the sweep blocked on a nonregular candidate")
	}
}
