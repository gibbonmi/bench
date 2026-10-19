package preflight

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// TestWritesResolveEarlierChunkDeletion covers AP48. A build commit before the chunk
// --base deleted the owned path, and the review record holds a chunk, so the deletion
// lies in the build's landing range: both the unmarked and the (deleted) entry resolve.
// A never-tracked path stays red. A record without a chunk and a landing range that does
// not resolve give no deletion facts, so the deleted path stays red. A bare preflight
// reads the landing range from the review base.
// TestCommandFrozenRangeWithoutDeletionEvidenceReds holds the same range without a
// review record red.
func TestWritesResolveEarlierChunkDeletion(t *testing.T) {
	const gone, typo = "internal/example/gone.go", "internal/example/typo.go"
	for _, tc := range []struct {
		name, missing, marker, wantRed string
		noChunk, noDefault, bare       bool
	}{
		{name: "unmarked entry"},
		{name: "deleted entry", marker: " (deleted)"},
		{name: "never-tracked path", missing: typo, wantRed: typo},
		{name: "record without a chunk", noChunk: true, wantRed: gone},
		{name: "landing range does not resolve", noDefault: true, wantRed: gone},
		{name: "bare preflight at the review base", bare: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, slug, base, deletion := seedDeletionReview(t, tc.missing, true)
			ticket := "specs/" + slug + "/tickets/one.md"
			doc, err := os.ReadFile(ticket)
			if err != nil {
				t.Fatal(err)
			}
			preflighttest.MustWriteFile(t, ticket, strings.Replace(string(doc), gone, gone+tc.marker, 1))
			chunks := []reviewrecord.Chunk{{ID: "C1", Base: base, Tip: deletion, PlanDigest: "sha256:plan", SourceDigest: deletion}}
			if tc.noChunk {
				chunks = []reviewrecord.Chunk{}
			}
			preflighttest.MustWriteFile(t, "reviews/"+slug+".md", recordFence(t, reviewrecord.Record{
				Version: 2, Spec: "specs/" + slug + "/spec.md", PlanDigest: "sha256:plan", Chunks: chunks,
			}))
			preflighttest.RunGit(t, "add", ".")
			preflighttest.RunGit(t, "commit", "-q", "-m", "record the chunk")
			if tc.noDefault {
				preflighttest.RunGit(t, "branch", "-m", "main", "trunk")
			}
			tip := preflighttest.RunGit(t, "rev-parse", "HEAD")
			args := []string{"build", slug, "--base", tip, "--source-tip", tip}
			if tc.bare {
				args = []string{"build", slug}
			}
			out, _ := Command(args)
			red := strings.Contains(out, "writes-resolve,red")
			if red != (tc.wantRed != "") || (red && (!strings.Contains(out, "one.md: "+tc.wantRed) || (tc.wantRed != gone && strings.Contains(out, "one.md: "+gone)))) {
				t.Fatalf("writes-resolve red = %t, want the red to name only %q:\n%s", red, tc.wantRed, out)
			}
		})
	}
}
