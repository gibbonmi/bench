// Stack reader tests: the strict stack file, its member bound, the landing grammar per member, and the overrides.
package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// stackMemberJSON is member i of a valid stack file, with the 1-based index in each value.
func stackMemberJSON(i int) string {
	return fmt.Sprintf(`{"path":"wt%d","request":"r%d","base":"b%d","source_tip":"t%d","message":"m%d"}`, i, i, i, i, i)
}

// stackFileJSON is a valid stack document of n members and no final newline.
func stackFileJSON(n int) string {
	members := make([]string, n)
	for i := range members {
		members[i] = stackMemberJSON(i + 1)
	}
	return `{"members":[` + strings.Join(members, ",") + `]}`
}

// writeStackFile writes body as stack.json in a new directory and returns that directory.
func writeStackFile(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	mustNoError(t, os.WriteFile(filepath.Join(dir, "stack.json"), []byte(body), 0o600))
	return dir
}

// requireStackRefusal fails unless err is the reader's refusal for member (zero for a
// refusal of the whole file or an override) and the input is empty.
func requireStackRefusal(t *testing.T, input StackInput, err error, member int) *StackInputError {
	t.Helper()
	var refusal *StackInputError
	requireTest(t, errors.As(err, &refusal), "error = %v, want a stack input refusal", err)
	requireTest(t, refusal.Member == member, "refusal member = %d (%v), want %d", refusal.Member, err, member)
	requireTest(t, reflect.DeepEqual(input, StackInput{}), "refused read returned input %+v", input)
	return refusal
}

func TestStackReaderRefusesAMalformedFile(t *testing.T) {
	t.Parallel()
	missingMessage := strings.Replace(stackFileJSON(3), `,"message":"m2"`, "", 1)
	for _, tc := range []struct {
		name   string
		body   string
		absent bool
		member int
	}{
		{name: "absent (LS54)", absent: true},
		{name: "empty (LS55)", body: ""},
		{name: "duplicate key (LS56)", body: strings.Replace(stackFileJSON(2), `]}`, `],"members":[`+stackMemberJSON(3)+","+stackMemberJSON(4)+`]}`, 1)},
		{name: "unknown field (LS57)", body: strings.Replace(stackFileJSON(2), `"message":"m1"`, `"message":"m1","mesage":"m1"`, 1)},
		{name: "content after the object (LS58)", body: stackFileJSON(2) + stackFileJSON(2)},
		{name: "one member (LS59)", body: stackFileJSON(1)},
		{name: "seventeen members (LS60)", body: stackFileJSON(17)},
		{name: "member without message (LS61)", body: missingMessage, member: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			if !tc.absent {
				dir = writeStackFile(t, tc.body)
			}
			input, err := ReadStack("stack.json", dir, nil)
			refusal := requireStackRefusal(t, input, err, tc.member)
			if tc.member != 0 {
				requireTest(t, strings.Contains(refusal.Error(), fmt.Sprintf("member %d", tc.member)), "refusal %q does not name member %d", refusal.Error(), tc.member)
			}
		})
	}
}

func TestStackReaderAcceptsSixteenMembers(t *testing.T) {
	t.Parallel()
	body := strings.Replace(stackFileJSON(16), `"message":"m3"`, `"spec":"s3","message":"m3"`, 1)
	input, err := ReadStack("stack.json", writeStackFile(t, body), nil)
	mustNoError(t, err)
	requireTest(t, len(input.Members) == 16 && len(input.Overrides) == 0, "input = %+v, want sixteen members and no override", input)
	for i, member := range input.Members {
		n := i + 1
		want := StackMember{
			Path: fmt.Sprintf("wt%d", n), Request: fmt.Sprintf("r%d", n), Base: fmt.Sprintf("b%d", n),
			SourceTip: fmt.Sprintf("t%d", n), Message: fmt.Sprintf("m%d", n), EffectiveTip: fmt.Sprintf("t%d", n),
		}
		if n == 3 {
			want.Spec = "s3"
		}
		requireTest(t, member == want, "member %d = %+v, want %+v", n, member, want)
	}
}

func TestStackReaderParsesEachMemberWithTheLandingGrammar(t *testing.T) {
	t.Parallel()
	body := strings.Replace(stackFileJSON(3), `"message":"m2"`, `"message":""`, 1)
	input, err := ReadStack("stack.json", writeStackFile(t, body), nil)
	refusal := requireStackRefusal(t, input, err, 2)
	requireTest(t, strings.Contains(refusal.Error(), landGrammar.Help) && strings.Contains(refusal.Error(), "member 2"),
		"refusal %q does not name member 2 and the landing grammar line", refusal.Error())
}

func TestStackReaderRefusesABadOverride(t *testing.T) {
	t.Parallel()
	dir := writeStackFile(t, stackFileJSON(3))
	for _, tc := range []struct {
		name      string
		overrides []string
	}{
		{name: "no member (LS96)", overrides: []string{"r9=o9"}},
		{name: "two for one request (LS97)", overrides: []string{"r1=o1", "r1=o2"}},
		{name: "no equals (LS98)", overrides: []string{"r1"}},
		{name: "empty request (LS98)", overrides: []string{"=o1"}},
		{name: "empty commit (LS98)", overrides: []string{"r1="}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input, err := ReadStack("stack.json", dir, tc.overrides)
			requireStackRefusal(t, input, err, 0)
		})
	}
	t.Run("valid override sets the effective tip", func(t *testing.T) {
		t.Parallel()
		input, err := ReadStack("stack.json", dir, []string{"r3=o3", "r2=o2"})
		mustNoError(t, err)
		requireTest(t, reflect.DeepEqual(input.Overrides, []StackOverride{{Request: "r3", Commit: "o3"}, {Request: "r2", Commit: "o2"}}),
			"overrides = %+v, want r3 then r2 in argument order", input.Overrides)
		for i, want := range []string{"t1", "o2", "o3"} {
			member := input.Members[i]
			requireTest(t, member.EffectiveTip == want && member.SourceTip == fmt.Sprintf("t%d", i+1),
				"member %d tips = (file %q, effective %q), want (t%d, %s)", i+1, member.SourceTip, member.EffectiveTip, i+1, want)
		}
	})
}

func TestStackReaderRefusesASpecialFile(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "stack.json")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	type result struct {
		input StackInput
		err   error
	}
	done := make(chan result, 1)
	go func() {
		input, err := ReadStack(path, "/", nil)
		done <- result{input, err}
	}()
	select {
	case got := <-done:
		refusal := requireStackRefusal(t, got.input, got.err, 0)
		requireTest(t, strings.Contains(refusal.Error(), "not a regular file"), "refusal %q does not name the special file", refusal.Error())
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("the stack reader blocked on a FIFO with no writer")
	}
}
