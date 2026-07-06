package commit

import (
	"strings"
	"testing"
)

// The integrated command behavior (block-check, gate, flip, stage, commit, exit codes)
// is gate-observed through the CLI in internal/contract/runtime; this unit test pins only
// the pure argument parser, whose branch table the black-box tests exercise only by
// outcome.
func TestParseArgs(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantMsg  string
		wantSpec string
		wantN    int    // number of positional paths
		wantErr  string // substring of usageErr; "" means no error
	}{
		{"message and one path", []string{"-m", "msg", "a.txt"}, "msg", "", 1, ""},
		{"spec flag consumes its value", []string{"-m", "m", "--spec", "feature", "a.txt"}, "m", "feature", 1, ""},
		{"paths on both sides of flags", []string{"a.txt", "-m", "m", "b.txt"}, "m", "", 2, ""},
		{"no message", []string{"a.txt"}, "", "", 0, "-m <msg> is required"},
		{"no paths", []string{"-m", "m"}, "", "", 0, "at least one <path> is required"},
		{"unknown flag", []string{"-m", "m", "--nope", "a.txt"}, "", "", 0, "unknown flag"},
		{"dangling -m", []string{"-m"}, "", "", 0, "-m needs a message"},
		{"dangling --spec", []string{"-m", "m", "--spec"}, "", "", 0, "--spec needs a slug"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, specSlug, paths, usageErr := parseArgs(tc.args)
			if tc.wantErr != "" {
				if usageErr == "" || !strings.Contains(usageErr, tc.wantErr) {
					t.Fatalf("usageErr = %q, want substring %q", usageErr, tc.wantErr)
				}
				return
			}
			if usageErr != "" {
				t.Fatalf("unexpected usageErr %q", usageErr)
			}
			if msg != tc.wantMsg || specSlug != tc.wantSpec || len(paths) != tc.wantN {
				t.Errorf("got msg=%q spec=%q paths=%v, want msg=%q spec=%q n=%d", msg, specSlug, paths, tc.wantMsg, tc.wantSpec, tc.wantN)
			}
		})
	}
}
