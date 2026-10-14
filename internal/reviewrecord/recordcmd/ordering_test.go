package recordcmd_test

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestChunkNeedsCurrentAmendment covers NR15 and NR16. A plan change, a prose edit
// included, refuses the next chunk entry before any record byte changes and names the
// amendment form; the amendment then admits the entry and keeps the prior one. The
// planned-cell refusal keeps its place ahead of the plan guard.
func TestChunkNeedsCurrentAmendment(t *testing.T) {
	for name, change := range map[string]func(*recordtest.Fixture){
		"plan":  func(f *recordtest.Fixture) { replan(f, "go test -count=1 ./...") },
		"prose": func(f *recordtest.Fixture) { f.Write(recordtest.Spec, specText(f)+"\nA later prose paragraph.\n") },
	} {
		t.Run(name, func(t *testing.T) {
			f := recorded(t, 2)
			prior := onlyChunk(t, f)
			change(f)
			base, tip := advance(f, "chunk 2")
			want := "; write it with bench record amendment '" + slug(t) + "' --source '" + tip + "' "
			if out := refuseArgs(t, f, "stale plan digest", chunkArgs(t, "2", base, tip)); !strings.Contains(out, want) {
				t.Fatalf("stale chunk refusal %q, want the amendment form with the slug and the tip %q", out, want)
			}
			amend(t, f)
			record(t, f, "2", base, tip)
			if got := read(t, f); len(got.Amendments) != 1 || len(got.Chunks) != 2 || !reflect.DeepEqual(got.Chunks[0], prior) {
				t.Fatalf("record after the amendment = %+v, want one amendment, chunk 1 unchanged, and chunk 2", got)
			}
		})
	}
	t.Run("planned cell first", func(t *testing.T) {
		f := recorded(t, 2)
		withMap(cited, planned)(f)
		base, tip := advance(f, "chunk 2")
		if out := refuseArgs(t, f, "planned seam cell", chunkArgs(t, "2", base, tip)); strings.Contains(out, "stale plan digest") {
			t.Fatalf("refusal %q, want the planned-cell refusal ahead of the plan guard", out)
		}
	})
}

// TestVerificationRepairNamesKnownValues: a verification refusal names the usage line of its
// repair form with the slug and each value that the refusing owner knows.
func TestVerificationRepairNamesKnownValues(t *testing.T) {
	undispatched := func(t *testing.T) *recordtest.Fixture {
		f, _ := recordtest.NewLinked(t, 1, recordtest.Delegate, recordtest.Undispatched)
		base, tip := advance(f, "chunk 1")
		record(t, f, "1", base, tip)
		return f
	}
	for _, tc := range []struct {
		name, want string
		fixture    func(*testing.T) *recordtest.Fixture
		edit       map[string]string
	}{
		{"missing chunk", "chunk '%s' --chunk '2' --base <commit> --tip <commit>", func(t *testing.T) *recordtest.Fixture { return recorded(t, 2) }, map[string]string{"--chunk": "2"}},
		{"undispatched ticket", "assignment '%s' --source <commit> --ticket '1.md' ", undispatched, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			refuseVerify(t, tc.fixture(t), "; write it with bench record "+fmt.Sprintf(tc.want, slug(t)), tc.edit)
		})
	}
}

// TestCompletionHelpOrdering covers NR22 and RC45: the completion help row states the record
// writes in the order that a completion needs, and each step names a record form.
func TestCompletionHelpOrdering(t *testing.T) {
	want := []string{"amendment", "chunk", "verification", "review", "confirmation", "fold", "final verification", "completion"}
	var description string
	for _, row := range recordcmd.HelpRows() {
		if strings.Fields(row.Suffix)[0] == "completion" {
			description = row.Description
		}
	}
	if _, order, ok := strings.Cut(description, "write in the order "); !ok || !slices.Equal(strings.Split(order, ", "), want) {
		t.Fatalf("completion help %q, want the order %q", description, want)
	}
	// A refusal names its repair form by that name, so each name, the assignment form
	// outside the completion order included, must answer its own help.
	for _, step := range append(want, "assignment") {
		form := strings.TrimPrefix(step, "final ")
		if out, code := recordcmd.Command("", []string{form, "--help"}); code != 0 || !strings.HasPrefix(out, "usage: bench record "+form+" ") {
			t.Errorf("%s --help = exit %d, output %q; want the form's usage", form, code, out)
		}
	}
}
