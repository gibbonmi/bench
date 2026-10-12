// Override verb tests. The reviewer grants one landing an override from a terminal in the
// primary checkout, and each refusal writes nothing.
package commitcmd

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
)

// overrideSet is a primary checkout, and the checkout and id of the active assignment that
// an override names.
type overrideSet struct {
	root, checkout, id string
}

// backOut is the request and the label of the assignment that each override names.
const backOut = "back-out"

// overrideFixture is a primary checkout whose ledger holds one active assignment, the back-out.
// Its policy lets outcomes A and B start in parallel.
func overrideFixture(t *testing.T) overrideSet {
	t.Helper()
	t.Setenv("BENCH_HOME", t.TempDir())
	root := commitmenttest.StagedInParallel(t, "A", "B")
	checkout := commitmenttest.Assignment(t, root, backOut)
	owner, _ := intent.AssignmentForWorktree(checkout)
	return overrideSet{root: root, checkout: checkout, id: owner.ID}
}

// overrideReason is the reason of every override that the tests grant.
const overrideReason = "back out the broken landing"

// overrideArgs are the verb arguments that override the assignment id with reason.
func overrideArgs(id, reason string) []string {
	return []string{overrideForm, "--assignment", id, "--reason", reason}
}

// pipeWith is the read end of a pipe that holds input. A pipe is no character device, so
// it stands for the stdin of an agent's tool call.
func pipeWith(t *testing.T, input string) *os.File {
	t.Helper()
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(write, input); err != nil {
		t.Fatal(err)
	}
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = read.Close() })
	return read
}

// overridden runs the verb in dir over stdin with the production terminal predicate.
func overridden(dir string, stdin io.Reader, args []string) (string, int) {
	var out strings.Builder
	code := Run(dir, args, stdin, &out)
	return out.String(), code
}

// requireOverrideRefusal fails unless out is the override refusal for sentence, at exit 1.
func requireOverrideRefusal(t *testing.T, out string, code int, sentence string) {
	t.Helper()
	if want := toon.Errorf("bench commitment override refused", sentence) + "\n"; code != 1 || !strings.HasPrefix(out, want) {
		t.Fatalf("override = (%d, %q), want exit 1 and the refusal %q", code, out, want)
	}
}

// atTerminal runs the verb in dir with typed as the reviewer's input at a terminal.
func atTerminal(dir, typed string, args []string) (string, int) {
	var out strings.Builder
	code := run(dir, args, console{stdin: strings.NewReader(typed), stdout: &out, isTerminal: func(io.Reader) bool { return true }})
	return out.String(), code
}

// BW71: the reviewer confirms the row that names the assignment, its label, the open build
// that refuses its landing, and the reason. The verb then writes one record. After the
// back-out binds outcome B, its own build admits its landing, so the row names only build A.
func TestOverrideWritesOneRecord(t *testing.T) {
	f := overrideFixture(t)
	if err := (commitrepo.Store{Root: commitmenttest.Assignment(t, f.root, "build")}).Start("A", "build", "specs/A/spec.md"); err != nil {
		t.Fatal(err)
	}
	grant := func(reason string) {
		t.Helper()
		out, code := atTerminal(f.root, f.id+"\n", overrideArgs(f.id, reason))
		if want := confirmation(t, f.id, "A", reason) + "overridden{assignment=" + f.id + "}\n"; code != 0 || out != want {
			t.Fatalf("override = (%d, %q), want exit 0 and %q", code, out, want)
		}
		got, err := commitmenttest.Overrides(f.root)
		if want := []intent.LandingOverride{{Assignment: f.id, Reason: reason}}; err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("overrides = %+v (%v), want %+v", got, err, want)
		}
	}
	grant(overrideReason)
	if err := (commitrepo.Store{Root: f.checkout}).Start("B", backOut, "specs/B/spec.md"); err != nil {
		t.Fatal(err)
	}
	grant("back out the bound landing")
}

// notConfirmed is the refusal of a confirmation that does not type the id.
const notConfirmed = "commitment override not confirmed"

// BW77 and BW78: only the exact id on its own line confirms. The null device passes the
// terminal test, so its EOF meets the confirmation. A typed line that differs in one byte,
// a trailing CR or NBSP among them, refuses at an injected terminal.
func TestOverrideNeedsTheTypedAssignment(t *testing.T) {
	f := overrideFixture(t)
	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = null.Close() })
	prompt := confirmation(t, f.id, "", overrideReason)
	commitmenttest.RequireLedgerKept(t, f.root, func() {
		out, code := overridden(f.root, null, overrideArgs(f.id, overrideReason))
		if !strings.HasPrefix(out, prompt) {
			t.Fatalf("null device: override = (%d, %q), want the prompt %q first", code, out, prompt)
		}
		requireOverrideRefusal(t, strings.TrimPrefix(out, prompt), code, notConfirmed)
		for _, typed := range []string{f.id + "\r\n", f.id + "\u00a0\n", " " + f.id + "\n", f.id[:31] + "\n", f.id + "0\n", f.id} {
			out, code := atTerminal(f.root, typed, overrideArgs(f.id, overrideReason))
			if !strings.HasPrefix(out, prompt) {
				t.Fatalf("typed %q: override = (%d, %q), want the prompt %q first", typed, code, out, prompt)
			}
			requireOverrideRefusal(t, strings.TrimPrefix(out, prompt), code, notConfirmed)
		}
	})
}

// confirmation is what the verb prints before it reads the typed id. It is the row that
// names the back-out assignment, the refusing builds, and reason, and then the prompt.
func confirmation(t *testing.T, id, builds, reason string) string {
	t.Helper()
	row, err := toon.Table(overrideTable, overrideFields, [][]string{{id, backOut, builds, reason}})
	if err != nil {
		t.Fatal(err)
	}
	return row + overridePrompt + "\n"
}

// BW80: an id that names no active assignment refuses before the prompt, so a typo grants
// nothing. An unknown id and an id outside the identity pattern refuse the same way.
func TestOverrideRefusesAnInactiveAssignment(t *testing.T) {
	f := overrideFixture(t)
	for _, id := range []string{"0123456789abcdef0123456789abcdef", "back out *"} {
		commitmenttest.RequireLedgerKept(t, f.root, func() {
			out, code := atTerminal(f.root, id+"\n", overrideArgs(id, overrideReason))
			requireOverrideRefusal(t, out, code, "commitment override assignment "+id+" is not active")
		})
	}
}

// BW81: the override reason follows the reason rule of the block verb. A blank reason and a
// reason with a newline, ESC, or BEL refuse with that rule's sentence before the prompt. The
// grammar refuses an empty value for both verbs before that rule.
func TestOverrideRefusesABadReason(t *testing.T) {
	f := overrideFixture(t)
	commitmenttest.RequireLedgerKept(t, f.root, func() {
		for _, reason := range []string{" ", "two\nlines", "esc\x1b", "bel\a"} {
			out, code := atTerminal(f.root, f.id+"\n", overrideArgs(f.id, reason))
			blocked, blockCode := Command(f.root, []string{"block", "--outcome", "A", "--reason", reason})
			sentence := strings.TrimPrefix(strings.SplitN(blocked, "\n", 2)[0], toon.Errorf("bench commitment block refused", ""))
			if blockCode != 1 || sentence == blocked {
				t.Fatalf("block with reason %q = (%d, %q), want the reason refusal", reason, blockCode, blocked)
			}
			requireOverrideRefusal(t, out, code, sentence)
		}
		_, blockCode := Command(f.root, []string{"block", "--outcome", "A", "--reason", ""})
		if out, code := atTerminal(f.root, f.id+"\n", overrideArgs(f.id, "")); blockCode != 2 || code != 2 {
			t.Fatalf("empty reason = (%d, %q), block = %d; want the grammar refusal of both", code, out, blockCode)
		}
	})
}

// BW79: a Bench worktree's own executable never writes the override. The checkout test
// runs first, so a worktree with a pipe on stdin meets it, and not the terminal test.
func TestOverrideRefusesOutsideThePrimaryCheckout(t *testing.T) {
	f := overrideFixture(t)
	worktree := commitmenttest.Assignment(t, f.root, "worktree")
	commitmenttest.RequireLedgerKept(t, f.root, func() {
		out, code := overridden(worktree, pipeWith(t, f.id+"\n"), overrideArgs(f.id, overrideReason))
		requireOverrideRefusal(t, out, code, "commitment override runs from the primary checkout")
	})
}

// BW76: a pipe on stdin is no terminal, so the verb refuses before it reads the typed id,
// although the pipe holds that id.
func TestOverrideRefusesWithoutATerminal(t *testing.T) {
	f := overrideFixture(t)
	commitmenttest.RequireLedgerKept(t, f.root, func() {
		out, code := overridden(f.root, pipeWith(t, f.id+"\n"), overrideArgs(f.id, overrideReason))
		requireOverrideRefusal(t, out, code, "commitment override needs the reviewer at a terminal; run it in your own terminal")
	})
}
