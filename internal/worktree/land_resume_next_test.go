package worktree

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/landing"
)

// TestLandingResumeNextNamesTheRequestSlot pins both resume forms to the one request
// placeholder: the form that quotes its values and the form that a value with a control
// byte drops to.
func TestLandingResumeNextNamesTheRequestSlot(t *testing.T) {
	t.Parallel()
	cases := map[string]landing.ReviewedResult{
		"line-safe":     {Commit: "c1", SourceBase: "b1", SourceTip: "t1"},
		"not-line-safe": {Commit: "c1\x1b", SourceBase: "b1", SourceTip: "t1"},
	}
	for name, result := range cases {
		t.Run(name, func(t *testing.T) {
			got := landingResumeNext(result, "x", "/tmp/source", "a1")
			if want := " --request " + requestSlot + " "; !strings.Contains(got, want) {
				t.Fatalf("resume command %q lacks %q", got, want)
			}
		})
	}
}
