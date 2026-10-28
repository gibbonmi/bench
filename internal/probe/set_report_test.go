package probe

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/probe/probetest"
)

// TestMutantSetSuccessOmitsTheReports covers BC53, BC54, BC55, and BC57. A set of four mutants
// that each bit and restored prints its three tables with trailing_bytes=0 in at most 10 lines,
// and no report. A silent row keeps the report of each started mutant, and a failed restore
// keeps the preserved table.
func TestMutantSetSuccessOmitsTheReports(t *testing.T) {
	zeroSwap := []string{"--swap", "return 0", "--with", "return 1"}
	for _, tc := range []struct {
		name      string
		entries   [][]string
		mutated   []probetest.Start
		fault     map[durablefiletest.Fault]map[string]error
		command   int
		started   []int
		preserved bool
	}{
		{name: "four bites", entries: [][]string{guardSwap, returnSwap, {"--omit", "return 0"}, zeroSwap}, mutated: []probetest.Start{bite, bite, bite, bite}},
		{name: "one silent row", entries: [][]string{guardSwap, returnSwap}, mutated: []probetest.Start{bite, pass}, command: 1, started: []int{1, 2}},
		{name: "failed restore", entries: [][]string{guardSwap, returnSwap}, mutated: []probetest.Start{bite, bite},
			fault: map[durablefiletest.Fault]map[string]error{restoreReplacement: {"rename": os.ErrPermission}}, command: 2, started: []int{1}, preserved: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := clampFixture(t)
			installSet(t, f, tc.mutated...)
			router := subjectFaults()
			if tc.fault != nil {
				router.Fail = tc.fault
			}
			out, code := runSetWith(t, f, router.Replacer(), setArgs(mutantsFile(t, tc.entries...)))
			if tc.started == nil {
				want := setRows(t, biteRow(1, "swap"), biteRow(2, "swap"), biteRow(3, "omit"), biteRow(4, "swap")) + setSelection(t, "passed", 0) + setExecution(t, 0, 0)
				if out != want || code != 0 || strings.Count(out, "\n") > 10 {
					t.Fatalf("stdout = (%q, %d), want (%q, 0) in at most 10 lines", out, code, want)
				}
				return
			}
			if code != tc.command {
				t.Fatalf("exit = %d, want %d\n%s", code, tc.command, out)
			}
			for _, n := range tc.started {
				if heading := "mutant_report[1]{mutant}:\n  " + strconv.Itoa(n) + "\n"; !strings.Contains(out, heading) {
					t.Errorf("stdout = %q, want the report of mutant %d", out, n)
				}
			}
			requireTrailingBytes(t, out, tc.command)
			if tc.preserved {
				requirePreserved(t, f, out)
			}
		})
	}
}
