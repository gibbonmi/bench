package diff

import (
	"reflect"
	"testing"
)

func TestParseNameStatusZ(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want [][]string
	}{
		{"empty diff", "", nil},
		{"single", "A\x00work.txt\x00", [][]string{{"A", "work.txt"}}},
		{"space and non-ascii raw", "A\x00a b.txt\x00A\x00café.txt\x00A\x00a\"q.txt\x00",
			[][]string{{"A", "a b.txt"}, {"A", "café.txt"}, {"A", `a"q.txt`}}},
		{"multiple statuses", "A\x00a\x00M\x00b\x00D\x00c\x00",
			[][]string{{"A", "a"}, {"M", "b"}, {"D", "c"}},
		},
	}
	for _, c := range cases {
		if got := parseNameStatusZ([]byte(c.in)); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: parseNameStatusZ = %v, want %v", c.name, got, c.want)
		}
	}
}
