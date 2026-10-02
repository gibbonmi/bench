package worktree

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/census"
)

// BO60: the landing prints the output breakdown on stderr beside the head breakdown. The
// row fixes the line for two canned records of one head.
func TestLandingPrintsOutputBreakdown(t *testing.T) {
	t.Parallel()
	request := "census-landed-output"
	f := publicLandingFixture(t, request, "", "")
	for _, size := range []int64{16938, 12040} {
		output := census.Output{Head: "bench worktree list", Lines: 55, Bytes: size, Spilled: true}
		if err := census.RecordOutput(f.home, f.root, f.creation.Assignment.ID, output, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
	var stdout, stderr bytes.Buffer
	code := LandCommand(f.root, f.home, landArgs(request, f.base, f.tip, f.creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stderr.String(), "census output{bench worktree list=2/28978}\n") {
		t.Fatalf("landing evidence = (%d, %q, %q), want the output breakdown on stderr", code, stdout.String(), stderr.String())
	}
	if !strings.HasSuffix(stdout.String(), ",census=0}\n") {
		t.Fatalf("landed record = %q, want the raw-call count unchanged by output records", stdout.String())
	}
}
