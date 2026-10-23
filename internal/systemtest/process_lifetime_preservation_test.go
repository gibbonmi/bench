//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// The differential rows run each production entry over canned producer shapes and compare
// the projection with the production baseline. A shape names the exit status of the child,
// whether the child leaves a normal descendant in its group, and whether the child ends on
// its own TERM instead of an exit.
type preservationShape struct {
	name       string
	exit       int
	descendant bool
	signal     bool
}

var preservationShapes = []preservationShape{
	{name: "complete zero exit", exit: 0},
	{name: "complete exit seven", exit: 7},
	{name: "surviving normal descendant", exit: 0, descendant: true},
}

// signalShape is the child that a TERM from inside its own body ends.
var signalShape = preservationShape{name: "terminating signal", signal: true}

// shapeScript is the canned child body of shape. It writes one line to each stream, publishes
// its group in mark, and leaves a descendant in that group when the shape asks for one. sleep
// and mv name the tools as the caller resolves them.
func shapeScript(shape preservationShape, sleep, mv, mark string) string {
	script := "echo canned-stdout\necho canned-stderr >&2\n"
	if shape.descendant {
		script += sleep + " 300 </dev/null >/dev/null 2>&1 &\n"
	}
	quoted := "'" + mark + "'"
	end := "exit " + strconv.Itoa(shape.exit) + "\n"
	if shape.signal {
		end = "kill -TERM \"$$\"\n"
	}
	return script + "echo \"$$\" > " + quoted + ".partial && " + mv + " " + quoted + ".partial " + quoted + "\n" + end
}

// jsonShapeScript is the canned child body of shape for a caller that decodes the child's
// stdout. It writes body to stdout first, and the canned lines of shape go to stderr.
func jsonShapeScript(shape preservationShape, body, mark string) string {
	return "printf '%s' '" + body + "'\n{\n" + shapeScript(shape, "sleep", "mv", mark) + "} >&2\n"
}

// groupOf waits for the group that the child of a shape published in mark, and kills that
// group at cleanup, so a red row leaves no descendant behind.
func groupOf(t *testing.T, mark string) int {
	t.Helper()
	pgid := processtest.WaitPID(t, mark)
	t.Cleanup(func() { _ = syscall.Kill(-pgid, syscall.SIGKILL) })
	return pgid
}

// TestProcessLifetimeGateDifferential is PL92. The gate keeps its exit projection of the
// gate script, and it drains a surviving normal descendant before its success.
func TestProcessLifetimeGateDifferential(t *testing.T) {
	t.Parallel()
	for _, shape := range preservationShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			mark := filepath.Join(t.TempDir(), "group")
			root := gittest.RepoOnBranch(t, "main")
			f := testrepo.NewGateFixture(t.TempDir(), "HOME")
			f.MustWrite(t, root, shapeScript(shape, f.Command("sleep"), f.Command("mv"), mark), "")
			gittest.Output(t, root, "add", ".")
			gittest.Output(t, root, "commit", "-qm", "fixture")
			env := []string{"HOME=" + t.TempDir(), "TMPDIR=" + t.TempDir(), "TEST_TELEMETRY_DIR=" + os.DevNull}
			result := systemSelected(t, root, env, "gate", "--fresh")
			pgid := groupOf(t, mark)
			// The baseline gate passes the script's status and both canned streams through
			// below its tree header.
			if result.code != shape.exit || !strings.HasSuffix(result.stdout, "\ncanned-stdout\n") || result.stderr != "canned-stderr\n" {
				t.Fatalf("gate over exit %d = (%d, %q, %q), want the script's status and streams", shape.exit, result.code, result.stdout, result.stderr)
			}
			if shape.descendant {
				// The drain completes before the gate returns, so the group is gone at once.
				if state := protection.ObserveGroup(pgid).State; state != protection.GroupAbsent {
					t.Fatalf("gate script group after the green gate = %v, want drained", state)
				}
			}
		})
	}
}

// TestProcessLifetimeExecDifferential is PL100. Exec passes the canned stdin, streams, exit
// status, environment, and assignment directory through as a direct run of the same child
// does, and it adds no normal-exit signal to a surviving descendant.
func TestProcessLifetimeExecDifferential(t *testing.T) {
	t.Parallel()
	fixture, source := execBoundFixture(t)
	for _, shape := range preservationShapes {
		t.Run(shape.name, func(t *testing.T) {
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			const inputs = "cat\nprintf '%s|%s\\n' \"$PRESERVE_VALUE\" \"$(pwd -P)\"\n"
			env := append(fixture.environment(fixture.home), "PRESERVE_VALUE=a b*ü")
			// The direct run is the baseline projection. Its own descendant would outlive
			// it in no group the test owns, and the descendant writes nothing, so the
			// baseline runs the same shape without it.
			plain := preservationShape{name: shape.name, exit: shape.exit}
			direct := owner.runWithInput(source.path, env, "canned stdin\n", "sh", "-c", inputs+shapeScript(plain, "sleep", "mv", filepath.Join(dir, "direct")))
			if err := owner.observeSelected(); err != nil {
				t.Fatal(err)
			}
			routed := owner.runWithInput(fixture.root, env, "canned stdin\n", owner.selected.path, execBoundArgs(source, "sh", "-c", inputs+shapeScript(shape, "sleep", "mv", mark))...)
			pgid := groupOf(t, mark)
			if routed.code != direct.code || routed.stdout != direct.stdout || !strings.Contains(routed.stderr, direct.stderr) {
				t.Fatalf("exec = (%d, %q, %q), want the direct run (%d, %q, %q)", routed.code, routed.stdout, routed.stderr, direct.code, direct.stdout, direct.stderr)
			}
			if shape.descendant {
				if state := protection.ObserveGroup(pgid).State; state != protection.GroupPresent {
					t.Fatalf("exec child group after its normal exit = %v, want the preserved descendant", state)
				}
				_ = syscall.Kill(-pgid, syscall.SIGKILL)
				processtest.WaitAbsent(t, pgid)
			}
		})
	}
}
