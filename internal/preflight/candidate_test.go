package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// redRow returns the rendered red row of check, and fails the test when the verdict holds
// no such row.
func redRow(t *testing.T, out, check string) string {
	t.Helper()
	row, rendered := rowOf(t, out, check)
	if !rendered || !strings.HasPrefix(strings.TrimSpace(row), check+",red,") {
		t.Fatalf("verdict holds no red %s row:\n%s", check, out)
	}
	return row
}

// TestBuildVerdictAtGradesTheCandidate covers PP01, PP02, PP03, and PP06: the grade reads the
// source tip, the completion plan, the decision maps, and an explicit range at a candidate
// that HEAD does not name.
func TestBuildVerdictAtGradesTheCandidate(t *testing.T) {
	t.Run("PP01 green candidate pins the tip", func(t *testing.T) {
		f := seedGreenCandidate(t)
		out, code := BuildVerdictAt(f.root, f.slug, "", f.candidate)
		if code != 0 || summaryCounts(t, out).total() != bareRowCount(modeBuild)+1 {
			t.Fatalf("green candidate = (%d), want 0 with %d rows:\n%s", code, bareRowCount(modeBuild)+1, out)
		}
	})
	t.Run("PP02 completion plan at the candidate", func(t *testing.T) {
		f := seedPlanCandidate(t)
		out, code := BuildVerdictAt(f.root, f.slug, "", f.candidate)
		if code != 1 || !strings.Contains(redRow(t, out, "completion-plan"), f.candidate) {
			t.Fatalf("plan candidate = (%d), want 1 with completion-plan red at the candidate:\n%s", code, out)
		}
	})
	t.Run("PP03 decision maps at the candidate", func(t *testing.T) {
		f := seedMapCandidate(t)
		out, code := BuildVerdictAt(f.root, f.slug, "", f.candidate)
		if code != 1 || !strings.Contains(redRow(t, out, "map-sources"), "decisions/m.md: Sources Path docs/kept.md"+missingFile) {
			t.Fatalf("map candidate = (%d), want 1 with map-sources red:\n%s", code, out)
		}
	})
	t.Run("PP06 explicit base ranges to the candidate", func(t *testing.T) {
		f := seedGreenCandidate(t)
		base := preflighttest.RunGit(t, "rev-parse", "main")
		out, _ := BuildVerdictAt(f.root, f.slug, base, f.candidate)
		rows := preflighttest.ResponseTable(t, out, "source")
		if len(rows) != 1 || rows[0]["base"] != base || rows[0]["tip"] != f.candidate {
			t.Fatalf("source rows = %v, want base %s and tip %s:\n%s", rows, base, f.candidate, out)
		}
	})
}

// TestBuildVerdictAtRefusesBadInput covers PP05, PP07, and PP08: an invalid coverage map, a
// candidate that names no commit, and a base that is not an ancestor of the candidate each
// refuse with exit 1 and print no verdict.
func TestBuildVerdictAtRefusesBadInput(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		run        func(t *testing.T) (string, int)
	}{
		{name: "PP05 two-predicate coverage row", want: "error: coverage map invalid — ",
			run: func(t *testing.T) (string, int) {
				f := seedTwoPredicateCandidate(t)
				return BuildVerdictAt(f.root, f.slug, "", f.candidate)
			}},
		{name: "PP07 candidate names no commit", want: "error: cannot resolve candidate — 'missing'",
			run: func(t *testing.T) (string, int) {
				f := seedGreenCandidate(t)
				return BuildVerdictAt(f.root, f.slug, "", "missing")
			}},
		{name: "PP08 base not an ancestor of the candidate", want: "error: --base is not an ancestor — ",
			run: func(t *testing.T) (string, int) {
				root, slug := preflighttest.SeedAdmitted(t)
				head := preflighttest.RunGit(t, "rev-parse", "HEAD")
				editFenced(t, slug)
				f := writeCandidate(t, root, slug, "main")
				return BuildVerdictAt(f.root, f.slug, head, f.candidate)
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, code := tc.run(t)
			if code != 1 || !strings.HasPrefix(out, tc.want) || strings.Contains(out, "checks{") {
				t.Fatalf("refusal = (%d), want 1 with prefix %q and no verdict:\n%s", code, tc.want, out)
			}
		})
	}
}

// TestBuildVerdictAtMatchesThePublishedVerdict covers PP04, PP33, and PP34. Each fixture
// varies one HEAD read between HEAD and the candidate. The candidate verdict must equal, byte
// for byte, the --source-tip verdict that the published command prints after HEAD moves to
// the candidate. The published run derives the expectation, so no grade output is authored.
func TestBuildVerdictAtMatchesThePublishedVerdict(t *testing.T) {
	for _, tc := range []struct {
		name  string
		seed  func(t *testing.T) candidateFixture
		green string
	}{
		{name: "PP04 green source tip", seed: seedGreenCandidate},
		{name: "PP04 completion plan", seed: seedPlanCandidate},
		{name: "PP04 decision maps", seed: seedMapCandidate},
		{name: "PP33 merge candidate", seed: seedMergeCandidate, green: "base-current"},
		{name: "PP34 default-equal path", seed: seedDefaultEqualCandidate, green: "paths-authorized"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.seed(t)
			got, gotCode := BuildVerdictAt(f.root, f.slug, "", f.candidate)
			f.publish(t)
			want, wantCode := Command([]string{"build", f.slug, "--source-tip", f.candidate})
			if got != want || gotCode != wantCode {
				t.Fatalf("candidate verdict (%d):\n%s\npublished verdict (%d):\n%s", gotCode, got, wantCode, want)
			}
			if tc.green == "" {
				return
			}
			if _, red := rowOf(t, got, tc.green); red {
				t.Fatalf("%s is red, want green:\n%s", tc.green, got)
			}
		})
	}
}
