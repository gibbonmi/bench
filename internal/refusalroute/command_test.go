package refusalroute

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// TestRecoveryListsEveryFace runs `bench recovery` and decodes its whole stdout. The
// header and the rows derive from the declared faces: one row for each face, in registry
// order, whose route cell is the face's route over no facts, so each slot prints its
// placeholder. A matrix kept apart from the registry, or a face it drops, turns this red.
// It covers RR51 and RR52.
func TestRecoveryListsEveryFace(t *testing.T) {
	faces := Inventory()
	out, code := RecoveryCommand(nil)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q", code, out)
	}
	header := fmt.Sprintf("recovery[%d]{verb,face,authority,route}:\n", len(faces))
	if !strings.HasPrefix(out, header) {
		t.Fatalf("stdout = %q, want it to start with %q", out, header)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatalf("stdout = %q, want one TOON document: %v", out, err)
	}
	if want := []string{"recovery"}; !reflect.DeepEqual(document.Blocks, want) {
		t.Fatalf("stdout blocks = %q, want %q", document.Blocks, want)
	}
	rows, err := document.Rows("recovery")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != len(faces) {
		t.Fatalf("matrix has %d rows, want one for each of the %d registered faces", len(rows), len(faces))
	}
	for i, face := range faces {
		want := map[string]any{
			"verb":      string(face.Verb),
			"face":      face.Name,
			"authority": string(face.Authority),
			"route":     face.Render(Facts{}),
		}
		if !reflect.DeepEqual(rows[i], want) {
			t.Errorf("row %d = %#v, want %#v", i, rows[i], want)
		}
	}
}
