package worktree

// This file checks the route that a verb result reports. A not-called assertion on a joins
// stub holds only when the run took the joins form with the call's joins value, so the
// result must tell that route apart from every other route.

import "testing"

func TestVerbResultReportsTheJoinsRoute(t *testing.T) {
	t.Parallel()
	call := verbCall{root: t.TempDir(), home: t.TempDir()}
	entry := runVerb(t, verbClean, call)
	kitOnly := call
	kitOnly.kit = t.TempDir()
	defaulted := runVerb(t, verbClean, kitOnly)
	j := defaultJoins()
	call.joins = &j
	runVerb(t, verbClean, call).mustViaJoins(t)
	for name, r := range map[string]verbResult{"public-entry": entry, "kit-only": defaulted} {
		var recorder verbRecorder
		r.mustViaJoins(&recorder)
		if len(recorder.failures) != 1 {
			t.Fatalf("mustViaJoins on a %s run recorded %q, want one failure", name, recorder.failures)
		}
	}
}
