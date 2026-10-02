package worktree

// This file checks the route that a verb result reports. A not-called assertion on a joins
// stub holds only when the run took the joins form, so the result must tell the two
// routes apart.

import "testing"

func TestVerbResultReportsTheJoinsRoute(t *testing.T) {
	t.Parallel()
	call := verbCall{root: t.TempDir(), home: t.TempDir()}
	entry := runVerb(t, verbClean, call)
	j := defaultJoins()
	call.joins = &j
	runVerb(t, verbClean, call).mustViaJoins(t)
	var recorder verbRecorder
	entry.mustViaJoins(&recorder)
	if len(recorder.failures) != 1 {
		t.Fatalf("mustViaJoins on a public-entry run recorded %q, want one failure", recorder.failures)
	}
}
