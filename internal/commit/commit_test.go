package commit

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/poolkey"
)

func TestAssessmentAssignmentSpan(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	assignment := "11111111111111111111111111111111"
	root := filepath.Join(t.TempDir(), poolkey.AssignmentSegment(assignment, assignment))
	finish := beginCommitSpan(root)
	finish(0, commitMeasures{subject: "published", pathCount: 3, counted: true})
	spans, err := otelrecord.ReadSpans(home, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 1 || spans[0].Attributes[otelrecord.AttrAssignmentID] != assignment || spans[0].Attributes[otelrecord.AttrSubjectID] != "published" || spans[0].Attributes[otelrecord.AttrMeasurePathCount] != "3" {
		t.Fatalf("assignment provenance missing: %+v", spans)
	}
}

// The integrated command behavior (block-check, gate, flip, stage, commit, exit codes)
// belongs to the black-box command tests; this unit test pins only the pure argument
// parser, whose branch table those tests exercise only by outcome. The rendered reason text belongs to the shared grammar and its toon renderer,
// so the misuse cases assert the whole line rather than a local phrasing.
func TestParseArgs(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantMsg  string
		wantPath string // first positional path; "" means none expected
		wantN    int    // number of positional paths
		wantErr  string // substring of usageErr; "" means no error
	}{
		{name: "message and one path", args: []string{"-m", "msg", "a.txt"}, wantMsg: "msg", wantPath: "a.txt", wantN: 1},
		{name: "paths on both sides of flags", args: []string{"a.txt", "-m", "m", "b.txt"}, wantMsg: "m", wantPath: "a.txt", wantN: 2},
		{name: "-- makes a leading-dash token a path", args: []string{"-m", "m", "--", "-weird.txt"}, wantMsg: "m", wantPath: "-weird.txt", wantN: 1},
		{name: "no message", args: []string{"a.txt"}, wantErr: "-m <msg> is required"},
		{name: "no paths", args: []string{"-m", "m"}, wantErr: "at least one <path> is required"},
		{name: "unknown flag", args: []string{"-m", "m", "--nope", "a.txt"}, wantErr: "usage: bench commit (unknown argument: --nope)"},
		{name: "dangling -m", args: []string{"-m"}, wantErr: "usage: bench commit (missing argument: -m)"},
		{name: "the retired --spec is unknown", args: []string{"-m", "m", "--spec", "x", "a.txt"}, wantErr: "usage: bench commit (unknown argument: --spec)"},
		// The two shapes an unset or blank shell variable produces. An empty path is
		// caught by the shared grammar (it would otherwise resolve to the cwd and
		// widen the commit); an empty or blank message is caught here, so it is
		// reported like the missing -m rather than by a raw `git commit -m ""`.
		{name: "empty path", args: []string{"-m", "m", ""}, wantErr: `usage: bench commit (unknown argument: "")`},
		{name: "empty message", args: []string{"-m", "", "a.txt"}, wantErr: "-m <msg> must not be empty"},
		{name: "blank message", args: []string{"-m", "   \t\n", "a.txt"}, wantErr: "-m <msg> must not be empty"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, paths, _, help, usageErr := parseArgs(tc.args)
			if help != "" {
				t.Fatalf("unexpected help %q", help)
			}
			if tc.wantErr != "" {
				if usageErr == "" || !strings.Contains(usageErr, tc.wantErr) {
					t.Fatalf("usageErr = %q, want substring %q", usageErr, tc.wantErr)
				}
				return
			}
			if usageErr != "" {
				t.Fatalf("unexpected usageErr %q", usageErr)
			}
			if msg != tc.wantMsg || len(paths) != tc.wantN {
				t.Errorf("got msg=%q paths=%v, want msg=%q n=%d", msg, paths, tc.wantMsg, tc.wantN)
			}
			if len(paths) > 0 && paths[0] != tc.wantPath {
				t.Errorf("first path = %q, want %q", paths[0], tc.wantPath)
			}
		})
	}
}

// Help is a success the caller prints, never a misuse: all three spellings return the
// declared help text with no usage error.
func TestParseArgsHelpIsSuccess(t *testing.T) {
	for _, spelling := range []string{"help", "--help", "-h"} {
		t.Run(spelling, func(t *testing.T) {
			_, _, _, help, usageErr := parseArgs([]string{spelling})
			if usageErr != "" {
				t.Fatalf("usageErr = %q, want none", usageErr)
			}
			if !strings.HasPrefix(help, "usage: bench commit") {
				t.Fatalf("help = %q, want the declared help text", help)
			}
		})
	}
}

func TestCommandGrammarErrorPrintsFlatUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Command([]string{"--unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("exit = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if got, want := stderr.String(), "usage: bench commit (unknown argument: --unknown)\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}
