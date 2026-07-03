package learnings

import (
	"reflect"
	"testing"
)

func TestRows(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want [][]string
	}{
		{"em-dash separator", "## 2026-01-01 — first learning  [open]\n- body\n", [][]string{{"2026-01-01", "first learning"}}},
		{"ascii-hyphen separator", "## 2026-01-01 - ascii title  [open]\n", [][]string{{"2026-01-01", "ascii title"}}},
		{"comma/quote title", `## 2026-03-03 — a, "b"  [open]`, [][]string{{"2026-03-03", `a, "b"`}}},
		{"CRLF heading", "## 2026-04-04 — crlf  [open]\r\n", [][]string{{"2026-04-04", "crlf"}}},
		{"no trailing newline", "## 2026-05-05 — tail  [open]", [][]string{{"2026-05-05", "tail"}}},
		{"template example is not an entry", "## <date> — <short title>  [open]\n", nil},
		{"resolved heading ignored", "## 2026-01-01 — done  [resolved]\n", nil},
		{"empty", "", nil},
	}
	for _, c := range cases {
		if got := Rows([]byte(c.in)); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: Rows = %v, want %v", c.name, got, c.want)
		}
	}
}
