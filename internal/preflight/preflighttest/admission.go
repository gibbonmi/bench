package preflighttest

import (
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"testing"
)

func admitFixture(t *testing.T, root string) {
	t.Helper()
	policy, exists, err := (commitrepo.Store{Root: root}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	if exists && policy.ActiveMilestone == "fixture" {
		commitmenttest.Admit(t, root, "preflight-assignment-target", policy.Milestones[0].Outcomes[0].Deliverables[0].Source.Path)
	}
}

// SeedAdmitted adds a current delivery binding to the conformant command fixture.
func SeedAdmitted(t *testing.T) (root, slug string) {
	t.Helper()
	root, slug = SeedConformant(t)
	ActiveAssignment(t, root, root)
	return root, slug
}
