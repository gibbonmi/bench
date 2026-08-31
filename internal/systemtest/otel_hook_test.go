//go:build system

package systemtest

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestOtelHookRecordJourney pipes a hook envelope to the built binary's hook plumbing
// dispatch and reads the record the run left under a private BENCH_HOME. The git guard is
// the hook verb that answers both ways from the same envelope shape, so one verb states
// both exits the record must tell apart.
//
// Row OT21.
func TestOtelHookRecordJourney(t *testing.T) {
	fixture := scaffoldRecordedPublicationRepo(t)

	for _, leg := range []struct {
		name, envelope, outcome string
		code                    int
	}{
		{
			name:     "allowed call",
			envelope: `{"tool_name":"Bash","tool_input":{"command":"rg guard AGENTS.md"}}`,
			outcome:  "green",
		},
		{
			name:     "refused call",
			envelope: `{"tool_name":"Bash","tool_input":{"command":"git reset --hard HEAD~1"}}`,
			outcome:  "red",
			code:     2,
		},
	} {
		t.Run(leg.name, func(t *testing.T) {
			// The reader globs for the record, so the home carries no path character a glob
			// pattern would read as a character class.
			home := filepath.Join(t.TempDir(), "bench-home")
			if err := owner.observeSelected(); err != nil {
				t.Fatal(err)
			}
			result := owner.runWithInput(fixture.root, fixture.environment(home), leg.envelope,
				owner.selected.path, "guard-git")
			if result.code != leg.code {
				t.Fatalf("bench guard-git = (%d, %q, %q), want exit %d",
					result.code, result.stdout, result.stderr, leg.code)
			}

			spans := readSeamRecord(t, home)
			// OT21: the dispatch names the hook verb it ran and the exit the harness saw.
			hook, present := spans["hook.guard-git"]
			if !present {
				t.Fatalf("the record has no hook span; it holds %v", spanNames(spans))
			}
			if hook.Attributes[otelrecord.AttrSeam] != "hook.guard-git" {
				t.Fatalf("hook seam = %q, want the verb's own seam", hook.Attributes[otelrecord.AttrSeam])
			}
			if hook.Attributes[otelrecord.AttrOutcome] != leg.outcome {
				t.Fatalf("hook outcome = %q, want %q for exit %d",
					hook.Attributes[otelrecord.AttrOutcome], leg.outcome, leg.code)
			}
			// Story 19: the dispatch carries no attribute outside the declared set.
			for name, span := range spans {
				for key := range span.Attributes {
					if !slices.Contains(otelrecord.DeclaredAttributes, key) {
						t.Fatalf("span %s carries the undeclared attribute %q", name, key)
					}
				}
			}
		})
	}
}
