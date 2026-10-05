package commitment

import (
	"reflect"
	"testing"
)

// The bound sources and their identity do not depend on the order of the input.
func TestBoundSourcesIgnoresInputOrder(t *testing.T) {
	first := SourceBinding{ID: "FT1", Path: "roadmap/FT1.md", Identity: Identity([]byte("FT1"))}
	ninth := SourceBinding{ID: "FT9", Path: "roadmap/FT9.md", Identity: Identity([]byte("FT9"))}
	descending, descendingIdentity := boundSources([]SourceBinding{ninth, first})
	ascending, ascendingIdentity := boundSources([]SourceBinding{first, ninth})
	if want := []SourceBinding{first, ninth}; !reflect.DeepEqual(descending, want) || !reflect.DeepEqual(ascending, want) {
		t.Fatalf("boundSources lists = %v and %v, want %v", descending, ascending, want)
	}
	if descendingIdentity != ascendingIdentity {
		t.Fatalf("boundSources identities = %s and %s, want one identity", descendingIdentity, ascendingIdentity)
	}
}
