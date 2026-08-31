//go:build system

package systemtest

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestOtelWorktreeVerbRecordJourney drives the built binary's worktree verbs against one
// scaffolded repository and reads the record each verb left behind. The four verbs share
// one Bench home, because the record is addressed by repository and the run is one
// repository's traffic.
//
// Row OT20.
func TestOtelWorktreeVerbRecordJourney(t *testing.T) {
	fixture := scaffoldRecordedPublicationRepo(t)
	tip := systemGitOutput(t, fixture.root, "rev-parse", "main")
	source := systemCreateLandingWorktree(t, fixture.root, fixture.home, "otel-verbs", "otel verbs")

	executed := systemSelected(t, fixture.root, fixture.environment(fixture.home),
		"worktree", "exec", source.assignment, "--", "git", "rev-parse", "HEAD")
	if executed.code != 0 {
		t.Fatalf("bench worktree exec = (%d, %q, %q)", executed.code, executed.stdout, executed.stderr)
	}
	// The merge target already holds the incoming commit, so the verb publishes nothing
	// and still records its own invocation.
	merged := systemSelected(t, fixture.root, fixture.environment(fixture.home),
		"worktree", "merge", "--from", tip, source.assignment)
	if merged.code != 0 {
		t.Fatalf("bench worktree merge = (%d, %q, %q)", merged.code, merged.stdout, merged.stderr)
	}
	released := systemSelected(t, fixture.root, fixture.environment(fixture.home),
		"worktree", "release", "--request", source.request, source.path)
	if released.code != 0 {
		t.Fatalf("bench worktree release = (%d, %q, %q)", released.code, released.stdout, released.stderr)
	}

	spans := readSeamRecord(t, fixture.home)
	// OT20: each verb of the run left one span that names the verb and the assignment it
	// acted on.
	for _, seam := range []string{"worktree.create", "worktree.exec", "worktree.merge", "worktree.release"} {
		span, present := spans[seam]
		if !present {
			t.Fatalf("the record has no %s span; it holds %v", seam, spanNames(spans))
		}
		if span.Attributes[otelrecord.AttrSeam] != seam {
			t.Errorf("%s span seam = %q, want the verb's own seam", seam, span.Attributes[otelrecord.AttrSeam])
		}
		if span.Attributes[otelrecord.AttrSubjectID] != source.assignment {
			t.Errorf("%s subject id = %q, want the assignment %s",
				seam, span.Attributes[otelrecord.AttrSubjectID], source.assignment)
		}
		if span.Attributes[otelrecord.AttrOutcome] != "green" {
			t.Errorf("%s outcome = %q, want the verb's own exit", seam, span.Attributes[otelrecord.AttrOutcome])
		}
		// A verb span states the verb, the assignment, and the exit, and nothing outside
		// the declared set reaches the record.
		for key := range span.Attributes {
			if !slices.Contains(otelrecord.DeclaredAttributes, key) {
				t.Errorf("%s span carries the undeclared attribute %q", seam, key)
			}
		}
	}
}
