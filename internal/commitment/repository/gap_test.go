package repository_test

import (
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// TestGapFailsOnADefaultBranchFault pins that OpenGap and CloseGap return the default-branch
// error and write no gap when the default-branch tip does not resolve.
func TestGapFailsOnADefaultBranchFault(t *testing.T) {
	plan := reviewrecord.Plan{Chunks: []reviewrecord.PlannedChunk{{ID: "C1", Tickets: []string{"1.md"}}, {ID: "C2", Tickets: []string{"2.md"}}}}
	root := commitmenttest.SeedMilestone(t)
	build := commitmenttest.BindBuild(t, root, "build", commitmenttest.MilestoneSpec)
	gittest.Output(t, root, "branch", "-m", "main", "trunk")
	store := commitrepo.Store{Root: build}
	for _, row := range []struct {
		name string
		run  func() error
	}{
		{name: "OpenGap", run: func() error { return store.OpenGap(plan, "C1") }},
		{name: "CloseGap", run: func() error { return store.CloseGap(plan, "1.md") }},
	} {
		t.Run(row.name, func(t *testing.T) {
			commitmenttest.RequireLedgerKept(t, root, func() {
				if err := row.run(); !commitrepo.DefaultUnresolved(err) {
					t.Fatalf("%s = %v, want the unresolved default branch error", row.name, err)
				}
			})
			if gaps, _ := commitmenttest.ClaimedGaps(t, root); len(gaps) != 0 {
				t.Fatalf("gap records after %s = %+v, want none", row.name, gaps)
			}
		})
	}
}
