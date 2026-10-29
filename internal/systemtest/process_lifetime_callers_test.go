//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// callerTool writes body as the executable name in dir and returns its path.
func callerTool(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// writeTree writes each slash-spelled file of files below root.
func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// toolPath is a PATH value that finds the tools in dir first.
func toolPath(dir string) string {
	return "PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// requireGroup fails unless the group of a shape's child is in state after its caller returned.
func requireGroup(t *testing.T, pgid int, state protection.GroupState, caller string) {
	t.Helper()
	if got := protection.ObserveGroup(pgid).State; got != state {
		t.Fatalf("%s child group after the caller returned = %v, want %v", caller, got, state)
	}
}

// TestProcessLifetimeShellDifferential is PL101. The shell verb exits with its release code
// whatever status the shell ends with, so a zero, a nonzero, and a signal shell each exit
// zero with the assignment released. The approved change adds no shell status repair. A
// surviving descendant gets no normal-exit signal, and its protection keeps the checkout.
func TestProcessLifetimeShellDifferential(t *testing.T) {
	t.Parallel()
	for _, shape := range append([]preservationShape{signalShape}, preservationShapes...) {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			fixture := scaffoldRecordedPublicationRepo(t)
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			shell := callerTool(t, dir, "shell", shapeScript(shape, "sleep", "mv", mark))
			result := systemSelected(t, fixture.root, append(fixture.environment(fixture.home), "SHELL="+shell, "TMPDIR="+t.TempDir()), "worktree", "shell")
			pgid := groupOf(t, mark)
			_, rest, _ := strings.Cut(result.stderr, "worktree: ")
			checkout, _, _ := strings.Cut(rest, "  (exit to release)")
			_, statErr := os.Stat(checkout)
			if !strings.Contains(result.stdout, "canned-stdout\n") || !strings.Contains(result.stderr, "canned-stderr\n") || checkout == "" {
				t.Fatalf("shell = (%d, %q, %q), want both canned streams and the checkout line", result.code, result.stdout, result.stderr)
			}
			if shape.descendant {
				requireGroup(t, pgid, protection.GroupPresent, "shell")
				if result.code == 0 || statErr != nil {
					t.Fatalf("shell with a live descendant = (%d, %q), checkout %v; want the retained checkout and its release code", result.code, result.stderr, statErr)
				}
				return
			}
			if result.code != 0 || !os.IsNotExist(statErr) {
				t.Fatalf("shell over %s = (%d, %q), checkout %v; want the release code zero and the released checkout", shape.name, result.code, result.stderr, statErr)
			}
		})
	}
}

// TestProcessLifetimeBoundsDifferential is PL94. The bounds runner projects a complete codex
// as the live inventory, an exit as the bundled fallback, a failed start as the unavailable
// row, and a codex past the provider deadline as the timeout row with its group stopped. A
// surviving descendant gets no normal-exit signal. No models caller cancels its codex, so
// the canceled projection has no reachable fixture.
func TestProcessLifetimeBoundsDifferential(t *testing.T) {
	t.Parallel()
	const inventory = `{"models":[{"slug":"fake-model","visibility":"list"}]}`
	bundled := "case \"$*\" in *--bundled*) printf '%s' '" + inventory + "'; exit 0 ;; esac\n"
	cases := []struct {
		shape preservationShape
		start bool
		hold  bool
		want  string
	}{
		{shape: preservationShapes[0], want: "codex,live,fake-model"},
		{shape: preservationShapes[1], want: "codex,bundled,fake-model"},
		{shape: preservationShapes[2], want: "codex,live,fake-model"},
		{shape: preservationShape{name: "start failure"}, start: true, want: "codex debug models unavailable"},
		{shape: preservationShape{name: "provider timeout"}, hold: true, want: "codex,timeout,timeout"},
	}
	for _, c := range cases {
		t.Run(c.shape.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			if c.start {
				// A missing interpreter fails the exec itself, so no child starts.
				if err := os.WriteFile(filepath.Join(dir, "codex"), []byte("#!"+filepath.Join(dir, "absent")+"\n"), 0o755); err != nil {
					t.Fatal(err)
				}
			} else if c.hold {
				callerTool(t, dir, "codex", publishGroup("mv", mark)+"exec sleep 300\n")
			} else {
				// The merged output is the inventory, so the canned lines go nowhere.
				callerTool(t, dir, "codex", bundled+"printf '%s' '"+inventory+"'\n{\n"+shapeScript(c.shape, "sleep", "mv", mark)+"} >/dev/null 2>&1\n")
			}
			env := []string{toolPath(dir), "TMPDIR=" + t.TempDir(), "OPENAI_API_KEY", "ANTHROPIC_API_KEY", "BENCH_OFFLINE"}
			result := systemSelected(t, t.TempDir(), env, "models")
			if !strings.Contains(result.stdout, c.want) || result.code != 0 {
				t.Fatalf("models over %s = (%d, %q, %q), want %q", c.shape.name, result.code, result.stdout, result.stderr, c.want)
			}
			if c.start {
				if strings.Contains(result.stdout, "fake-model") {
					t.Fatalf("models with no codex start = %q, want no model row", result.stdout)
				}
				return
			}
			pgid := groupOf(t, mark)
			if c.shape.descendant {
				requireGroup(t, pgid, protection.GroupPresent, "codex")
			}
			if c.hold {
				requireGroup(t, pgid, protection.GroupAbsent, "codex past its deadline")
			}
		})
	}
}

// kitRootEnv is the environment of a focused run at the kit root with no inherited selection,
// so the run builds its own private executable. The fake go in tools answers the test run.
func kitRootEnv(home, root, tools, tmp string) []string {
	return []string{"BENCH_COMMAND_OBSERVE=1", "BENCH_HOME=" + home, "BENCH_RUN_BINARY", "BENCH_KIT=" + root,
		"LAND_BASELINE_BENCH=" + owner.selected.path, toolPath(tools), "TMPDIR=" + tmp}
}

// realGo is the go executable that a fake go passes every other subcommand to.
func realGo(t *testing.T) string {
	t.Helper()
	path, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// goEvents is a canned go test event stream. Each package passes one test, and the packages
// in terminal also end with their package pass.
func goEvents(packages []string, terminal ...string) string {
	events := ""
	for _, pkg := range packages {
		events += `{"Action":"start","Package":"` + pkg + `"}` + "\n" + `{"Action":"run","Package":"` + pkg + `","Test":"TestA"}` + "\n" + `{"Action":"pass","Package":"` + pkg + `","Test":"TestA","Elapsed":0}` + "\n"
	}
	for _, pkg := range terminal {
		events += `{"Action":"pass","Package":"` + pkg + `","Elapsed":0}` + "\n"
	}
	return events
}

// fakeGo writes a go into dir that answers each arguments pattern with its script and passes
// every other invocation to the real go, so the freshness check keeps its own go list.
func fakeGo(t *testing.T, dir string, arms map[string]string) {
	t.Helper()
	body := "case \"$*\" in\n"
	for pattern, script := range arms {
		body += pattern + ")\n" + script + ";;\n"
	}
	callerTool(t, dir, "go", body+"*) exec '"+realGo(t)+"' \"$@\" ;;\nesac\n")
}

// goTestPattern is the fakeGo pattern of a go test run.
const goTestPattern = "'test '*"

// jsonCase is one decoder input: the producer shape, the canned stdout, and the line that the
// caller's projection of that input names.
type jsonCase struct {
	shape preservationShape
	body  string
	want  string
	code  int
}

// TestProcessLifetimeFocusedDifferential is PL98. The focused run keeps its report of a
// complete terminal stream, marks a nonzero child red, and refuses malformed and incomplete
// streams. A surviving descendant gets no normal-exit signal, and its protected group
// refuses the report.
func TestProcessLifetimeFocusedDifferential(t *testing.T) {
	t.Parallel()
	pkg, other := "example.invalid/a", "example.invalid/b"
	valid := goEvents([]string{pkg}, pkg)
	cases := []jsonCase{
		{shape: preservationShapes[0], body: valid, want: pkg, code: 0},
		{shape: preservationShapes[1], body: valid, want: pkg, code: 1},
		{shape: preservationShapes[2], body: valid, want: "go test cleanup incomplete", code: 1},
		{shape: preservationShape{name: "malformed JSON"}, body: "{not json\n", want: "go test output malformed", code: 1},
		{shape: preservationShape{name: "incomplete JSON"}, body: goEvents([]string{pkg, other}, pkg), want: "go test reported incomplete packages", code: 1},
	}
	for _, c := range cases {
		t.Run(c.shape.name, func(t *testing.T) {
			t.Parallel()
			fixture := scaffoldRecordedPublicationRepo(t)
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			fakeGo(t, dir, map[string]string{goTestPattern: jsonShapeScript(c.shape, c.body, mark)})
			result := systemSelected(t, fixture.root, append(fixture.environment(fixture.home), toolPath(dir), "TMPDIR="+t.TempDir()), "test", "--package", "./...")
			pgid := groupOf(t, mark)
			if out := result.stdout + result.stderr; result.code != c.code || !strings.Contains(out, c.want) {
				t.Fatalf("focused test over %s = (%d, %q), want %d and %q", c.shape.name, result.code, out, c.code, c.want)
			}
			if c.shape.descendant {
				requireGroup(t, pgid, protection.GroupPresent, "go test")
			}
		})
	}
}

// TestProcessLifetimeListDifferential is PL99. The changed selection keeps the package that
// go list names for a changed file and refuses a failed, malformed, or outside listing. The
// list registers no resource, so a surviving descendant gets no normal-exit signal and the
// selection stands.
func TestProcessLifetimeListDifferential(t *testing.T) {
	t.Parallel()
	pkg := "example.invalid/a"
	listing := func(dir string) string {
		return `{"Dir":"` + dir + `","ImportPath":"` + pkg + `","Match":["./..."]}` + "\n"
	}
	cases := []struct {
		jsonCase
		outside bool
	}{
		{jsonCase: jsonCase{shape: preservationShapes[0], want: pkg, code: 0}},
		{jsonCase: jsonCase{shape: preservationShapes[1], want: "go list failed", code: 1}},
		{jsonCase: jsonCase{shape: preservationShapes[2], want: pkg, code: 0}},
		{jsonCase: jsonCase{shape: preservationShape{name: "malformed JSON"}, body: "{not json\n", want: "go list output malformed", code: 1}},
		{jsonCase: jsonCase{shape: preservationShape{name: "package outside"}, want: "go list returned a package outside the repository", code: 1}, outside: true},
	}
	for _, c := range cases {
		t.Run(c.shape.name, func(t *testing.T) {
			t.Parallel()
			fixture := scaffoldRecordedPublicationRepo(t)
			writeTree(t, fixture.root, map[string]string{"go.mod": "module example.invalid\n\ngo 1.25\n", "a/a.go": "package a\n"})
			systemGit(t, fixture.root, "add", ".")
			systemGit(t, fixture.root, "commit", "-qm", "list base")
			base := systemGitOutput(t, fixture.root, "rev-parse", "HEAD")
			if err := os.WriteFile(filepath.Join(fixture.root, "a", "a.go"), []byte("package a\n\nconst A = 1\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			systemGit(t, fixture.root, "commit", "-qam", "list change")
			body := c.body
			if body == "" {
				body = listing(filepath.Join(systemGitOutput(t, fixture.root, "rev-parse", "--show-toplevel"), "a"))
			}
			if c.outside {
				body = listing(t.TempDir())
			}
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			fakeGo(t, dir, map[string]string{"'list -buildvcs=false -json -test ./...'": jsonShapeScript(c.shape, body, mark), goTestPattern: "printf '%s' '" + goEvents([]string{pkg}, pkg) + "'\n"})
			result := systemSelected(t, fixture.root, append(fixture.environment(fixture.home), toolPath(dir), "TMPDIR="+t.TempDir()), "test", "--changed", "--base", base, "--source-tip", "HEAD")
			pgid := groupOf(t, mark)
			if out := result.stdout + result.stderr; result.code != c.code || !strings.Contains(out, c.want) {
				t.Fatalf("changed test over %s = (%d, %q), want %d and %q", c.shape.name, result.code, out, c.code, c.want)
			}
			if c.shape.descendant {
				requireGroup(t, pgid, protection.GroupPresent, "go list")
			}
		})
	}
}

// TestProcessLifetimeBuilderDifferential is PL93. The private build validates and uses the
// executable of a zero build script, and it names the script's output for a failed build.
// The builder drains a surviving descendant before its success.
func TestProcessLifetimeBuilderDifferential(t *testing.T) {
	t.Parallel()
	for _, shape := range preservationShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			fixture := scaffoldRecordedPublicationRepo(t)
			configureArtifactLandingFixture(t, fixture.root)
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			script := filepath.Join(fixture.root, "scripts", "go-build.sh")
			stub, err := os.ReadFile(script)
			if err != nil {
				t.Fatal(err)
			}
			// The shape runs after the stub publishes its executable. The stub's PATH holds no
			// shell tools, so the shape names them absolutely.
			sleep, sleepErr := exec.LookPath("sleep")
			mv, mvErr := exec.LookPath("mv")
			if sleepErr != nil || mvErr != nil {
				t.Fatal(sleepErr, mvErr)
			}
			if err := os.WriteFile(script, append(stub, shapeScript(shape, sleep, mv, mark)...), 0o755); err != nil {
				t.Fatal(err)
			}
			systemGit(t, fixture.root, systemCommitArgs("-qam", "shaped builder")...)
			fakeGo(t, dir, map[string]string{goTestPattern: "printf '%s' '" + goEvents([]string{"landingrace"}, "landingrace") + "'\n"})
			result := systemSelected(t, fixture.root, kitRootEnv(fixture.home, fixture.root, dir, t.TempDir()), "test", "--package", "./...")
			out := result.stdout + result.stderr
			// The CLI has returned, so a build script that ran has published its group.
			if _, err := os.Stat(mark); err != nil {
				t.Fatalf("test with no inherited selection = (%d, %q), want the build script's run", result.code, out)
			}
			pgid := groupOf(t, mark)
			if shape.exit != 0 {
				if result.code != 1 || !strings.Contains(out, "build private Bench executable") || !strings.Contains(out, "canned-stdout") || !strings.Contains(out, "canned-stderr") {
					t.Fatalf("test over a failed build = (%d, %q), want the build failure with the script's output", result.code, out)
				}
				return
			}
			if result.code != 0 || !strings.Contains(out, "landingrace") {
				t.Fatalf("test over a built executable = (%d, %q), want the validated selection's report", result.code, out)
			}
			// The drain completes before the builder returns, so the group is gone at once.
			requireGroup(t, pgid, protection.GroupAbsent, "builder")
		})
	}
}
