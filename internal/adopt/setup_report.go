package adopt

import (
	"fmt"
	"io"
)

// finishSetup is the exit contract every convergeSetup path ends through. It invokes the
// one doctor rendering source, reportDoctorRows, with no second renderer, then prints the
// harness reload instruction and the exact next action. priorPartial carries a partial
// verdict convergeSetup already knows about, such as a written conflict block or the
// zero-signal stub. This makes the exit code answer for the whole run, never just the
// doctor rows in isolation. Setup exits zero only when everything, the transaction and
// doctor, is green.
func finishSetup(stdout io.Writer, facts setupFacts, priorPartial bool) int {
	red := reportDoctorRows(stdout)
	fmt.Fprintln(stdout, "reload your harness session so it picks up the converged AGENTS.md / CLAUDE.md instructions")
	if facts.zeroSignal {
		fmt.Fprintln(stdout, "next: configure .bench/gate.sh - replace the "+SentinelMarker+" sentinel with your project's real checks, then re-run bench setup")
	} else {
		fmt.Fprintln(stdout, "next: continue the /bench-setup-repo conversation to fill in the gate command, seams, and lines")
	}
	if red || priorPartial {
		return 3
	}
	return 0
}
