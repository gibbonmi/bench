package main

import (
	"bytes"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestCommandRunsVersionInProcess(t *testing.T) {
	var stdout, stderr, observation bytes.Buffer
	command := Command{Stdout: &stdout, Stderr: &stderr, Executable: "/selected/bench", Observe: &observation}
	if code := command.Run([]string{"version"}); code != 0 {
		t.Fatalf("version exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	want := versionLine(version, runtime.GOOS, runtime.GOARCH) + "\n"
	if stdout.String() != want {
		t.Fatalf("version stdout = %q, want %q", stdout.String(), want)
	}
	if want := "command-registry:version\n"; observation.String() != want {
		t.Fatalf("version implementation observation = %q, want %q", observation.String(), want)
	}
}

func TestCommandDispositionsAreComplete(t *testing.T) {
	want := map[processAttachment][]string{
		attachmentDirect: {"check-agent-line", "commit", "gate-go", "guard-git", "resume-clean", "session-inspect", "shift", "spec", "version", "worktree"},
		attachmentSystem: {"canary", "doctor", "freshness-check", "freshness-publish", "gate", "gate-phases", "gate-pin", "gate-run", "init", "link", "setup", "stop-verdict", "unlink", "upgrade", "worktree-hook"},
		attachmentShip:   {"prep-release", "release", "release-preflight"},
	}
	got := map[processAttachment][]string{}
	seen := map[string]bool{}
	for _, disposition := range commandDispositions() {
		if seen[disposition.Name] {
			t.Fatalf("command disposition repeats %q", disposition.Name)
		}
		seen[disposition.Name] = true
		got[disposition.Attachment] = append(got[disposition.Attachment], disposition.Name)
	}
	for attachment := range got {
		sort.Strings(got[attachment])
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command dispositions = %#v, want %#v", got, want)
	}
}

// keptRoutes is the surface a removal may not take with it, written down rather than
// derived. Every other routing check reads commandRegistry, so deleting a route deletes it
// from the expectation in the same edit and the check stays green; an enumeration authored
// against the reviewer's keep list is the only thing an over-broad deletion turns red. Each
// entry drives the real dispatcher and asserts the route answers its own grammar at exit 0,
// so a surviving-but-misrouted verb fails as loudly as a deleted one.
var keptRoutes = []struct {
	argv []string
	help string
}{
	{[]string{"worktree", "--help"}, "usage: bench worktree"},
	{[]string{"gate", "--help"}, "usage: bench gate"},
	{[]string{"commit", "--help"}, "usage: bench commit"},
	{[]string{"status", "--help"}, "usage: bench status"},
	{[]string{"guards", "--help"}, "usage: bench guards"},
	{[]string{"idea", "--help"}, "usage: bench idea"},
	{[]string{"roadmap", "--help"}, "usage: bench roadmap"},
	{[]string{"spec", "implemented", "--help"}, "usage: bench spec implemented"},
	{[]string{"spec", "retire", "--help"}, "usage: bench spec retire"},
	{[]string{"spec", "history", "--help"}, "usage: bench spec history"},
}

// keptWorktreeGrammars are the pool operations the worktree family help has to keep naming.
// The family route surviving says nothing about the operations under it, and each one is
// reached only through that dispatcher, so the grammar line is where their survival shows.
var keptWorktreeGrammars = []string{
	"bench worktree create",
	"bench worktree path",
	"bench worktree exec",
	"bench worktree release",
	"bench worktree clean",
}

func TestKeptRoutesAnswerTheirOwnHelp(t *testing.T) {
	for _, kept := range keptRoutes {
		name := strings.Join(kept.argv, " ")
		t.Run(name, func(t *testing.T) {
			out, code := runKeptRoute(kept.argv)
			if code != 0 {
				t.Fatalf("%s exit = %d, want 0; output=%q", name, code, out)
			}
			if !strings.Contains(out, kept.help) {
				t.Fatalf("%s output = %q, want it to name %q", name, out, kept.help)
			}
		})
	}
}

func TestKeptWorktreeOperationsKeepTheirGrammar(t *testing.T) {
	out, code := runKeptRoute([]string{"worktree", "--help"})
	if code != 0 {
		t.Fatalf("worktree --help exit = %d, want 0; output=%q", code, out)
	}
	for _, grammar := range keptWorktreeGrammars {
		if !strings.Contains(out, grammar) {
			t.Errorf("worktree help = %q, want it to name %q", out, grammar)
		}
	}
}

// runKeptRoute joins both sinks: help lands on stdout for some grammars and stderr for
// others, and which sink a route picked is not what these two tests are grading.
func runKeptRoute(argv []string) (string, int) {
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	return stdout.String() + stderr.String(), code
}
