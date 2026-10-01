package harness

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// The worktree-lifecycle hook passes each Claude worktree event to WorktreeCommand, and
// WorktreeCommand states that all event validation happens before the shared lifecycle
// is called. Each test here feeds one malformed event and reads the verdict: the exit
// code, the stderr verdict, and an empty stdout, because create prints a path only on
// success.

const (
	createFieldVerdict = "bench worktree-hook create: event requires session_id, cwd, and name"
	removeFieldVerdict = "bench worktree-hook remove: event requires session_id and worktree_path"
	removeRepoVerdict  = "bench worktree-hook remove: repository unavailable: "
	// projectDirEnv names the Claude project directory, the first remove hint.
	projectDirEnv = "CLAUDE_PROJECT_DIR"
	// eventLimit is the documented event size limit: one MiB.
	eventLimit = 1 << 20
)

var oversizeVerdict = fmt.Sprintf("invalid event JSON: event exceeds %d-byte limit", eventLimit)

// assertRefused calls WorktreeCommand with args and the event on stdin. It fails the
// test unless the hook refused with exit want and a stderr that holds verdict.
func assertRefused(t *testing.T, args []string, event string, want int, verdict string) {
	t.Helper()
	assertReaderRefused(t, args, strings.NewReader(event), want, verdict)
}

// assertReaderRefused is assertRefused with stdin as a reader, so a test can feed a
// stdin whose read fails.
func assertReaderRefused(t *testing.T, args []string, stdin io.Reader, want int, verdict string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := WorktreeCommand(args, stdin, &stdout, &stderr)
	if code != want || stdout.Len() != 0 || !strings.Contains(stderr.String(), verdict) {
		t.Errorf("WorktreeCommand(%q) = exit %d, stdout %q, stderr %q; want exit %d, empty stdout, stderr holding %q", args, code, stdout.String(), stderr.String(), want, verdict)
	}
}

// isolateRemove points the repository hint and the Bench home at fresh directories
// outside any repository. A remove event that passes validation then reaches no real
// repository and no real Bench home, even under a mutation of the validation.
func isolateRemove(t *testing.T) {
	t.Setenv(projectDirEnv, t.TempDir())
	t.Setenv(benchhome.Env, t.TempDir())
}

// padTo returns event followed by spaces up to size bytes. JSON permits the trailing
// whitespace, so only the size limit can refuse the padded event.
func padTo(event string, size int) string {
	return event + strings.Repeat(" ", size-len(event))
}

func TestWorktreeCommandRefusesAnUnknownAction(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{nil, {""}, {"delete"}, {"Create"}, {"create", "remove"}} {
		assertRefused(t, args, `{"session_id":"s"}`, 2, "usage: bench worktree-hook create|remove")
	}
}

// A stdin read that fails partway refuses even when the bytes read before the failure
// form a complete event. The event names only directories outside any repository, so
// a mutation that ignores the read error still reaches no repository.
func TestWorktreeCommandRefusesAnEventWhoseReadFails(t *testing.T) {
	isolateRemove(t)
	const readFailure = "event stdin closed"
	outside := t.TempDir()
	event := `{"session_id":"s","cwd":"` + outside + `","name":"n","worktree_path":"` + outside + `"}`
	for _, action := range []string{"create", "remove"} {
		stdin := io.MultiReader(strings.NewReader(event), iotest.ErrReader(errors.New(readFailure)))
		assertReaderRefused(t, []string{action}, stdin, 1, "bench worktree-hook "+action+": invalid event JSON: "+readFailure)
	}
}

func TestWorktreeCommandRefusesAnOversizeEvent(t *testing.T) {
	t.Parallel()
	event := padTo(`{"session_id":"s"}`, eventLimit+1)
	assertRefused(t, []string{"create"}, event, 1, "bench worktree-hook create: "+oversizeVerdict)
	assertRefused(t, []string{"remove"}, event, 1, "bench worktree-hook remove: "+oversizeVerdict)
}

func TestWorktreeCommandAcceptsAnEventAtTheSizeLimit(t *testing.T) {
	t.Parallel()
	assertRefused(t, []string{"create"}, padTo(`{"session_id":"s"}`, eventLimit), 1, createFieldVerdict)
}

func TestWorktreeCommandRefusesAnEmptyEvent(t *testing.T) {
	t.Parallel()
	for _, event := range []string{"", " \n\t "} {
		assertRefused(t, []string{"create"}, event, 1, "bench worktree-hook create: invalid event JSON: empty input")
		assertRefused(t, []string{"remove"}, event, 1, "bench worktree-hook remove: invalid event JSON: empty input")
	}
}

func TestWorktreeCommandRefusesAnEventThatIsNotJSON(t *testing.T) {
	t.Parallel()
	for _, event := range []string{
		"session_id=s",
		`{"session_id":"s"`,
		`{"session_id":"s","cwd":"/","name":"n"} {}`,
		`{"session_id":5,"cwd":"/","name":"n"}`,
		`["s","/","n"]`,
	} {
		assertRefused(t, []string{"create"}, event, 1, "bench worktree-hook create: invalid event JSON: ")
		assertRefused(t, []string{"remove"}, event, 1, "bench worktree-hook remove: invalid event JSON: ")
	}
}

func TestWorktreeCommandCreateRefusesAMissingField(t *testing.T) {
	t.Parallel()
	outside := t.TempDir()
	for _, event := range []string{
		`null`,
		`{}`,
		`{"cwd":"` + outside + `","name":"n"}`,
		`{"session_id":"s","name":"n"}`,
		`{"session_id":"s","cwd":"` + outside + `"}`,
		`{"session_id":"s","cwd":"` + outside + `","name":" \t"}`,
	} {
		assertRefused(t, []string{"create"}, event, 1, createFieldVerdict)
	}
}

func TestWorktreeCommandRemoveRefusesAMissingField(t *testing.T) {
	isolateRemove(t)
	outside := t.TempDir()
	for _, event := range []string{
		`null`,
		`{}`,
		`{"worktree_path":"` + outside + `"}`,
		`{"session_id":"s"}`,
		`{"session_id":"s","worktree_path":""}`,
	} {
		assertRefused(t, []string{"remove"}, event, 1, removeFieldVerdict)
	}
}

func TestWorktreeCommandCreateRefusesACwdOutsideARepository(t *testing.T) {
	t.Parallel()
	event := `{"session_id":"s","cwd":"` + t.TempDir() + `","name":"n"}`
	assertRefused(t, []string{"create"}, event, 1, "bench worktree-hook create: cwd is not in a Git repository: ")
}

// The remove hint is CLAUDE_PROJECT_DIR when it is set. A hint outside a repository
// refuses even when the event's worktree path is inside one.
func TestWorktreeCommandRemoveRefusesAProjectDirOutsideARepository(t *testing.T) {
	isolateRemove(t)
	event := `{"session_id":"s","worktree_path":"` + gittest.Repo(t) + `"}`
	assertRefused(t, []string{"remove"}, event, 1, removeRepoVerdict)
}

// With CLAUDE_PROJECT_DIR empty, the remove hint is the event's worktree path.
func TestWorktreeCommandRemoveRefusesAWorktreePathOutsideARepository(t *testing.T) {
	isolateRemove(t)
	t.Setenv(projectDirEnv, "")
	event := `{"session_id":"s","worktree_path":"` + t.TempDir() + `"}`
	assertRefused(t, []string{"remove"}, event, 1, removeRepoVerdict)
}

// An event that passes validation in a repository with no commit reaches the shared
// lifecycle, and the lifecycle cannot resolve a start for the worktree.
// Mutation that turns this red: delete the worktree.Create error return, so create
// prints an empty path and exits 0.
func TestWorktreeCommandCreateRefusesACreateError(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	event := `{"session_id":"s","cwd":"` + gittest.Repo(t) + `","name":"n"}`
	assertRefused(t, []string{"create"}, event, 1, "bench worktree-hook create: resolve assignment start: ")
}

// A symlinked worktrees admin directory makes the registration scan refuse, so remove
// has no registration to release through.
// Mutation that turns this red: delete the registration guard, so remove indexes an
// empty registration list.
func TestWorktreeCommandRemoveRefusesAnUnavailableRegistration(t *testing.T) {
	isolateRemove(t)
	root := gittest.Repo(t)
	if err := os.Symlink(t.TempDir(), filepath.Join(root, ".git", "worktrees")); err != nil {
		t.Fatal(err)
	}
	scan := git.ScanWorktreeAdmin(filepath.Join(root, ".git"))
	if scan == nil {
		t.Fatal("registration scan accepted a symlinked worktrees admin directory")
	}
	t.Setenv(projectDirEnv, root)
	event := `{"session_id":"s","worktree_path":"` + root + `"}`
	assertRefused(t, []string{"remove"}, event, 1, "bench worktree-hook remove: worktree registration unavailable: "+scan.Error())
}

// A remove event from a session that owns no assignment reaches the release through
// the worktree path, and the release refuses it.
func TestWorktreeCommandRemoveRefusesASessionThatOwnsNoAssignment(t *testing.T) {
	isolateRemove(t)
	t.Setenv(projectDirEnv, "")
	event := `{"session_id":"s","worktree_path":"` + gittest.Repo(t) + `"}`
	assertRefused(t, []string{"remove"}, event, 1, "request token matches no assignment")
}
