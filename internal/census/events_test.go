package census

import (
	"path/filepath"
	"testing"
)

// TestReadEventsRefusesAForeignName proves the event reader opens only an assignment's
// record file, so a stray file in the census directory, or a name that climbs out of
// it, never becomes an event source.
func TestReadEventsRefusesAForeignName(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	writeRecordFile(t, home, root, "README", "t\tsed\n")
	for _, name := range []string{"README", filepath.Join("..", "census")} {
		events, problems, err := ReadEvents(home, root, name)
		if err == nil || events != nil || problems != nil {
			t.Errorf("ReadEvents(%q) = %v, %v, %v; want a refusal", name, events, problems, err)
		}
	}
}
