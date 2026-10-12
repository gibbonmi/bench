package intent

import (
	"reflect"
	"strings"
	"testing"
)

// buildRecords is one lease, one gap, and one override that the ledger accepts.
func buildRecords() CommitmentState {
	return CommitmentState{
		Leases:    []LandingLease{{Destination: "refs/heads/main", Assignment: strings.Repeat("b", 32)}},
		Gaps:      []ChunkGap{{Milestone: "m1", Outcome: "o1", Chunk: "A", Open: true}},
		Overrides: []LandingOverride{{Assignment: strings.Repeat("e", 32), Reason: "back out the release"}},
	}
}

// TestValidateCommitmentRefusesDuplicateBuildRecords is BW12. A hand edit that adds a
// second owner of one record, an invalid assignment, or an unreadable reason must not
// authorize anything, so the ledger refuses each one.
func TestValidateCommitmentRefusesDuplicateBuildRecords(t *testing.T) {
	valid := buildRecords()
	if err := ValidateCommitment(&valid); err != nil {
		t.Fatalf("one lease, gap, and override = %v, want them accepted", err)
	}
	for _, row := range []struct {
		name   string
		mutate func(*CommitmentState)
	}{
		{"second lease on one destination", func(s *CommitmentState) {
			s.Leases = append(s.Leases, LandingLease{Destination: "refs/heads/main", Assignment: strings.Repeat("c", 32)})
		}},
		{"second gap for one outcome", func(s *CommitmentState) {
			s.Gaps = append(s.Gaps, ChunkGap{Milestone: "m1", Outcome: "o1", Chunk: "B"})
		}},
		{"second override for one assignment", func(s *CommitmentState) {
			s.Overrides = append(s.Overrides, LandingOverride{Assignment: strings.Repeat("e", 32), Reason: "again"})
		}},
		{"invalid lease assignment", func(s *CommitmentState) { s.Leases[0].Assignment = "not an id" }},
		{"empty override reason", func(s *CommitmentState) { s.Overrides[0].Reason = " " }},
		{"override reason with a control byte", func(s *CommitmentState) { s.Overrides[0].Reason = "back out\x1b[2J" }},
	} {
		t.Run(row.name, func(t *testing.T) {
			state := buildRecords()
			row.mutate(&state)
			if err := ValidateCommitment(&state); err == nil {
				t.Fatalf("ValidateCommitment(%+v) = nil, want a refusal", state)
			}
		})
	}
}

// TestTolerantReadKeepsTheBuildRecords is BW13. A recovery pass reads the ledger in the
// tolerant mode and writes it back, so the lease, the gap, and the override must survive
// that write.
func TestTolerantReadKeepsTheBuildRecords(t *testing.T) {
	root, path := mutatorFixture(t)
	seed, err := readPath(path)
	if err != nil {
		t.Fatal(err)
	}
	records := buildRecords()
	seed.Commitment = &records
	if err := writePath(path, seed); err != nil {
		t.Fatal(err)
	}
	if err := Transact(root, TolerantRead, func(ledger Ledger) (Ledger, bool, error) { return ledger, true, nil }, nil); err != nil {
		t.Fatal(err)
	}
	got, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.Commitment == nil || !reflect.DeepEqual(*got.Commitment, buildRecords()) {
		t.Fatalf("commitment after a tolerant write = %+v, want %+v", got.Commitment, buildRecords())
	}
}
