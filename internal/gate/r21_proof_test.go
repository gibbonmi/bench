package gate

import (
	"reflect"
	"sort"
	"testing"
)

func TestR21DeterministicFaultProofRegistryCompleteness(t *testing.T) {
	t.Parallel()
	got := make([]string, 0, len(r21ProofRegistry))
	seen := map[string]bool{}
	for _, proof := range r21ProofRegistry {
		if proof.id == "" || proof.driver == nil || seen[proof.id] {
			t.Fatalf("%s: invalid or duplicate registration %q", r21CompletenessFailure, proof.id)
		}
		seen[proof.id] = true
		got = append(got, proof.id)
	}
	want := append([]string(nil), r21ExpectedProofIDs...)
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: got IDs %v, want %v", r21CompletenessFailure, got, want)
	}
}

func TestR21DeterministicFaultEngine(t *testing.T) {
	t.Parallel()
	for _, proof := range r21ProofRegistry {
		proof := proof
		t.Run(proof.id, proof.driver)
	}
}
