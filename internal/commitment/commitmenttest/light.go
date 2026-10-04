package commitmenttest

import (
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TicketsSlug names the tickets-only folder that WriteTickets writes.
const TicketsSlug = "t"

// TicketsFolder is the repository path of that folder.
const TicketsFolder = "specs/" + TicketsSlug

// WriteTickets writes the light-path tickets-only folder below root and returns the
// identity that approves it. The identity is the folder's Git tree, so the folder stays
// staged in the index. The caller commits.
func WriteTickets(t testing.TB, root string) string {
	t.Helper()
	Write(t, root, TicketsFolder+"/tickets/one.md", "Light path ticket.\n")
	gittest.Output(t, root, "add", "--", TicketsFolder)
	return commitrepo.TreeIdentity(gittest.Output(t, root, "write-tree", "--prefix="+TicketsFolder+"/"))
}

// SeedTicketsOnly writes the SeedClosure policy and the tickets-only folder, and approves
// that folder as the complete delivery of FT1. The spec at deliverable stays approved for
// the delivery outcome with no obligation, so an assignment can start from it and rebind
// to the folder. The caller commits.
func SeedTicketsOnly(t testing.TB, root, deliverable string, residual ...string) {
	t.Helper()
	SeedClosure(t, root, deliverable, residual...)
	identity := WriteTickets(t, root)
	EditPolicy(t, root, func(policy *commitment.Policy) {
		delivery := &policy.Milestones[0].Outcomes[0]
		delivery.Deliverables[0].Obligations = nil
		delivery.Deliverables = append(delivery.Deliverables, commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: "tickets", Path: TicketsFolder, Identity: identity}, Obligations: []string{"FT1"}})
	})
}

// RowlessIndex is the board index that SeedRowless writes. Only outcome B owns a row.
const RowlessIndex = "# Roadmap\n\n## Parked\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. " + DeliveryOutcome + "\n2. B\n"

// DeliveredIndex is the board after the delivery outcome is delivered from the index of
// SeedClosure or of SeedRowless: only outcome B and its row remain.
const DeliveredIndex = "# Roadmap\n\n## Parked\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. B\n"

// ResidualIndex is the board after a delivery closes FT1 from the SeedClosure index with
// the residual row FT3. The delivery outcome keeps FT3 and its sequence entry.
const ResidualIndex = "# Roadmap\n\n## Parked\n\n**FT3 — " + DeliveryOutcome + "**\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. " + DeliveryOutcome + "\n2. B\n"

// SeedRowless writes an active milestone whose delivery outcome owns no source and
// approves the spec at deliverable and the tickets-only folder. Outcome B owns FT2 and
// follows it on the board. The caller commits.
func SeedRowless(t testing.TB, root, deliverable string) {
	t.Helper()
	SeedAdmission(t, root, deliverable)
	identity := WriteTickets(t, root)
	EditPolicy(t, root, func(policy *commitment.Policy) {
		milestone := &policy.Milestones[0]
		milestone.Outcomes[0].Deliverables = append(milestone.Outcomes[0].Deliverables, commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: "tickets", Path: TicketsFolder, Identity: identity}})
		milestone.Outcomes = append(milestone.Outcomes, outcome("B", writeRow(t, root, "FT2", "B")))
	})
	Write(t, root, "ROADMAP.md", RowlessIndex)
}
