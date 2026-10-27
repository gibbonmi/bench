package carry_test

import (
	"path"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestCarryReach covers each delta that reaches its requirement. The system suite and the
// race table each own a whole-only package, and root conformance reads the conformance
// package, so those three paths stand here as independent expectations.
func TestCarryReach(t *testing.T) {
	const (
		reached    = "the delta reaches the command"
		notTest    = "the delta is not test-only"
		wholeOnly  = "the delta changes a whole-only test file"
		systemPkg  = "internal/systemtest"
		racePkg    = "internal/worktree"
		conformPkg = "internal/conformance"
	)
	each := func(reason string) []carryCase {
		return []carryCase{{command: systemRun, reason: reason}, {command: packageRun(pkgA), reason: reason}, {command: conformanceRun, reason: reason}}
	}
	decideRows(t, []carryRow{
		{name: "a test change in the run package", delta: testOnlyIn(pkgA), cases: []carryCase{
			{command: packageRun(pkgA), reason: reached},
			{command: "bench test --package ./" + path.Dir(pkgA) + "/...", reason: reached},
		}},
		{name: "a doc change", delta: func(f *recordtest.Fixture) { f.CommitDelta("README.md") }, cases: each(notTest)},
		{name: "a ticket change", delta: func(f *recordtest.Fixture) { f.CommitDelta(f.TicketPath("1.md")) }, cases: each(notTest)},
		{name: "a system-suite test change", delta: testOnlyIn(systemPkg), cases: []carryCase{{command: systemRun, reason: wholeOnly}}},
		{name: "a race-table test change", delta: testOnlyIn(racePkg), cases: each(wholeOnly)},
		{name: "a test-only delta in another package", delta: testOnlyIn(pkgB), cases: []carryCase{
			{command: conformanceRun, reason: reached},
			{command: "bench test --check prose", reason: reached},
			{command: packageRun(conformPkg), reason: reached},
			{command: "go test ./" + pkgA, reason: reached},
			{command: packageRun(pkgA) + " && true", reason: reached},
			{command: systemRun + " --fixtures", reason: reached},
			{command: packageRun(pkgA) + "/", reason: reached},
		}},
		{name: "an assignment-only plan commit", delta: func(f *recordtest.Fixture) {
			f.Reassign("1.md")
			f.Refreeze()
		}, cases: []carryCase{{command: conformanceRun, reason: reached}}},
	})
}
