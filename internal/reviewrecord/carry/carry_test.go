package carry_test

import (
	"fmt"
	"reflect"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/carry"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// The fixture packages: a delta changes a test file of one, and a requirement runs the other.
const pkgA, pkgB = "internal/a", "internal/b"

const (
	systemRun      = "bench test --check system"
	conformanceRun = "bench test --check root-conformance"
)

func testIn(pkg string) string     { return pkg + "/carry_test.go" }
func packageRun(pkg string) string { return "bench test --package ./" + pkg }

// carryCase is one owed probe requirement of a row and its expected reason, empty for a carry.
// A chunk requirement is always a probe, so the reason reads the command and the probe entry.
type carryCase struct {
	command, reason string
}

// carryProbe is the planned probe of every owed requirement.
const carryProbe = "omit the check"

// carryRow is one fixture: delta commits the change after the prior entries, edit reshapes the
// chunk entry before the record saves, and each case is decided against the saved chunk.
type carryRow struct {
	name  string
	delta func(*recordtest.Fixture)
	edit  func(*recordtest.Fixture, *rr.Chunk)
	cases []carryCase
}

func testOnlyIn(pkg string) func(*recordtest.Fixture) {
	return func(f *recordtest.Fixture) { f.CommitDelta(testIn(pkg)) }
}

// decideRows decides each case of each row. A carry must answer the last entry of its
// requirement in the saved chunk, and a refusal must answer no candidate.
func decideRows(t *testing.T, rows []carryRow) {
	t.Helper()
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			owed := make([]rr.Requirement, len(row.cases))
			for i, c := range row.cases {
				owed[i] = rr.Requirement{ID: fmt.Sprintf("owed-%d", i), Command: c.command, Probe: carryProbe}
			}
			f := recordtest.Carrying(t, owed)
			if row.delta != nil {
				row.delta(f)
			}
			if row.edit != nil {
				row.edit(f, &f.Record.Chunks[0])
			}
			chunk := f.SavedChunk()
			for i, requirement := range f.Plan.Chunks[0].Verification {
				candidate, reason, err := carry.Decide(f.Root, recordtest.Spec, chunk, requirement)
				if err != nil {
					t.Fatalf("%s: %v", requirement.Command, err)
				}
				var want rr.Verification
				if row.cases[i].reason == "" {
					want = lastEntry(t, chunk, requirement.ID)
				}
				if reason != row.cases[i].reason || !reflect.DeepEqual(candidate, want) {
					t.Errorf("%s: reason %q, candidate %+v; want reason %q, candidate %+v", requirement.Command, reason, candidate, row.cases[i].reason, want)
				}
			}
		})
	}
}

func lastEntry(t *testing.T, chunk rr.Chunk, id string) rr.Verification {
	t.Helper()
	for i := len(chunk.Verification) - 1; i >= 0; i-- {
		if chunk.Verification[i].Requirement == id {
			return chunk.Verification[i]
		}
	}
	t.Fatalf("the fixture chunk holds no entry of %s", id)
	return rr.Verification{}
}

// TestCarryDecision covers the carry of an unreached requirement and the refusal of each
// candidate that is not a current passing entry.
func TestCarryDecision(t *testing.T) {
	const noEntry = "no current passing entry"
	failed := func(entry *rr.Verification) {
		one := 1
		entry.ID, entry.Outcome, entry.ExitCode = entry.ID+"-failed", "fail", &one
	}
	decideRows(t, []carryRow{
		{name: "a test-only delta in another package", delta: testOnlyIn(pkgB), cases: []carryCase{
			{command: systemRun},
			{command: systemRun + " --run TestA"},
			{command: packageRun(pkgA)},
		}},
		{name: "an assignment-only plan commit and a test change", delta: func(f *recordtest.Fixture) {
			f.Reassign("1.md")
			f.CommitDelta(testIn(pkgB))
		}, cases: []carryCase{{command: packageRun(pkgA)}}},
		{name: "an unchanged source", cases: []carryCase{{command: conformanceRun}}},
		{name: "no prior entry", delta: testOnlyIn(pkgB), edit: func(_ *recordtest.Fixture, chunk *rr.Chunk) {
			chunk.Verification = nil
		}, cases: []carryCase{{command: packageRun(pkgA), reason: noEntry}}},
		{name: "a failed last entry", delta: testOnlyIn(pkgB), edit: func(_ *recordtest.Fixture, chunk *rr.Chunk) {
			failed(&chunk.Verification[0])
		}, cases: []carryCase{{command: packageRun(pkgA), reason: noEntry}}},
		{name: "an earlier pass and a failed last entry", delta: testOnlyIn(pkgB), edit: func(_ *recordtest.Fixture, chunk *rr.Chunk) {
			last := chunk.Verification[0]
			failed(&last)
			chunk.Verification = append(chunk.Verification, last)
		}, cases: []carryCase{{command: packageRun(pkgA), reason: noEntry}}},
		{name: "a legacy entry with no source commit", delta: testOnlyIn(pkgB), edit: func(f *recordtest.Fixture, chunk *rr.Chunk) {
			f.Legacy(chunk.Verification[0].Requirement)
		}, cases: []carryCase{{command: packageRun(pkgA), reason: noEntry}}},
		{name: "an entry off the chunk line", delta: testOnlyIn(pkgB), edit: func(f *recordtest.Fixture, chunk *rr.Chunk) {
			f.Bind(&chunk.Verification[0], f.OffLine())
		}, cases: []carryCase{{command: packageRun(pkgA), reason: noEntry}}},
		{name: "a probe that did not bite", delta: testOnlyIn(pkgB), edit: func(_ *recordtest.Fixture, chunk *rr.Chunk) {
			chunk.Verification[0].Probe.Outcome = "silent"
		}, cases: []carryCase{{command: packageRun(pkgA), reason: noEntry}}},
	})
}
