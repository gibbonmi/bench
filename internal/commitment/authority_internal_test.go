package commitment

import (
	"reflect"
	"testing"
)

// The bound sources and their identity do not depend on the order of the input. In byte
// order, the lowercase identifier a sorts after FT9. Two revisions of row FT9 share an
// identifier and a path, so they order by identity: the digest of "FT9 revised" sorts before
// the digest of "FT9".
func TestBoundSourcesIgnoresInputOrder(t *testing.T) {
	row := func(id, content string) SourceBinding {
		return SourceBinding{ID: id, Path: "roadmap/" + id + ".md", Identity: Identity([]byte(content))}
	}
	first, ninth, revised, lower := row("FT1", "FT1"), row("FT9", "FT9"), row("FT9", "FT9 revised"), row("a", "a")
	descending, descendingIdentity := boundSources([]SourceBinding{lower, ninth, revised, first})
	ascending, ascendingIdentity := boundSources([]SourceBinding{first, revised, ninth, lower})
	if want := []SourceBinding{first, revised, ninth, lower}; !reflect.DeepEqual(descending, want) || !reflect.DeepEqual(ascending, want) {
		t.Fatalf("boundSources lists = %v and %v, want %v", descending, ascending, want)
	}
	if descendingIdentity != ascendingIdentity {
		t.Fatalf("boundSources identities = %s and %s, want one identity", descendingIdentity, ascendingIdentity)
	}
}
