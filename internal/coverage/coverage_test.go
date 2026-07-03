package coverage

import (
	"reflect"
	"strings"
	"testing"
)

const stories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n"
const hdr = "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|\n"

func spec(body string) parsed { return parse([]byte(body)) }

func TestStateAndRows(t *testing.T) {
	p := spec("# t\n\n" + stories + "\n### Acceptance coverage map\n" + hdr +
		"| 2–3 | does x \\| y | cli seam | cmd fails, loudly | catches z |\n" +
		"| edge of 1 | edge case | gate | already covered | catches w |\n")
	if State(p) != "mapped" {
		t.Fatalf("state = %q", State(p))
	}
	want := [][]string{{"2–3", "cli seam", "cmd fails, loudly"}, {"edge of 1", "gate", "already covered"}}
	if got := Rows(p); !reflect.DeepEqual(got, want) {
		t.Errorf("Rows = %v, want %v", got, want)
	}

	if State(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")) != "historical" {
		t.Error("historical state not detected")
	}
	if State(spec("# n\nprose only\n")) != "no-map" {
		t.Error("no-map state not detected")
	}
}

// Every validation phrasing is matched by substring downstream; pin each here.
func TestCheck(t *testing.T) {
	valid := spec("# v\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1, 2–3 | b | s | r | w |\n")
	if v := Check(valid); len(v) != 0 {
		t.Errorf("valid map violations = %v", v)
	}
	// Historical opts out of validation.
	if v := Check(spec("# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")); v != nil {
		t.Errorf("historical Check = %v, want nil", v)
	}
	cases := []struct{ body, want string }{
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n", "coverage map missing the canonical header"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr, "coverage map has no data rows"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b | s | r |\n", "coverage map row 1 has 4 cells (want 5)"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 1 | b |  | r | w |\n", "coverage map row 1 has an empty 'seam' cell"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| 9 | b | s | r | w |\n", "references story 9 but the spec numbers only 3"},
		{"# b\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + "| x | b | s | r | w |\n", "has an unrecognized story reference 'x'"},
	}
	for _, c := range cases {
		v := Check(spec(c.body))
		if len(v) == 0 || !strings.Contains(strings.Join(v, "\n"), c.want) {
			t.Errorf("Check violations %v do not contain %q", v, c.want)
		}
	}
}
