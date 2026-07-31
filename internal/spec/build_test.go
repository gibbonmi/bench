package spec

import (
	"strings"
	"testing"
)

func TestParseBuildExposesExactlyEightOperations(t *testing.T) {
	cases := []struct {
		op   string
		args []string
	}{
		{"start", nil}, {"assign", []string{"--ticket", "one.md", "--request", "request"}},
		{"checkpoint", []string{"--assignment", "a", "--evidence", "/receipt"}},
		{"integrate", []string{"--assignment", "a"}}, {"review", []string{"--evidence", "/review"}},
		{"status", []string{"--full"}}, {"promote", nil}, {"abandon", []string{"--apply", "fingerprint"}},
	}
	for _, tc := range cases {
		t.Run(tc.op, func(t *testing.T) {
			got, out, code := ParseBuild(append([]string{tc.op, "slug"}, tc.args...))
			if code != 0 || out != "" || got.Operation != tc.op || got.Slug != "slug" {
				t.Fatalf("ParseBuild = %+v, %q, %d", got, out, code)
			}
		})
	}
	if _, out, code := ParseBuild([]string{"ninth", "slug"}); code != 2 || !strings.Contains(out, "unknown argument: ninth") {
		t.Fatalf("unknown operation = %q, %d", out, code)
	}
}

func TestParseBuildKeepsFlagValuesAndTerminatedTextLiteral(t *testing.T) {
	got, out, code := ParseBuild([]string{"assign", "slug", "--ticket", "promote", "--request", "--assignment"})
	if code != 0 || out != "" || got.Flags["--ticket"] != "promote" || got.Flags["--request"] != "--assignment" {
		t.Fatalf("flag values = %+v, %q, %d", got, out, code)
	}
	for _, args := range [][]string{{"start", "--", "slug"}, {"assign", "slug", "--", "--ticket"}} {
		if _, out, code := ParseBuild(args); code != 2 || !strings.HasPrefix(out, "usage:") {
			t.Fatalf("terminated text %v = %q, %d", args, out, code)
		}
	}
}
