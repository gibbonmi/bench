package preflight

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// recordFence renders the bench-review-record section the charge packet reads.
func recordFence(t *testing.T, record reviewrecord.Record) string {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	return "# Review outcomes\n\n```bench-review-record\n" + string(data) + "\n```\n"
}

// The charge packet describes both valid record versions, and it keeps naming
// an invalid record as invalid. The projection reports the state the reader
// found; it does not decide acceptance.
func TestDelegatedEvidenceProjection(t *testing.T) {
	cases := []struct {
		name, want string
		body       func(*testing.T) string
	}{
		{"version 1", "parsed", func(t *testing.T) string {
			return recordFence(t, reviewrecord.Record{Version: 1, Spec: "specs/example/spec.md", PlanDigest: "sha256:x", ImplementationSession: "author"})
		}},
		{"version 2", "parsed", func(t *testing.T) string {
			return recordFence(t, reviewrecord.Record{Version: 2, Spec: "specs/example/spec.md", PlanDigest: "sha256:x"})
		}},
		{"version 2 smuggling an identity", "invalid implementation session", func(t *testing.T) string {
			return recordFence(t, reviewrecord.Record{Version: 2, Spec: "specs/example/spec.md", PlanDigest: "sha256:x", ImplementationSession: "smuggled"})
		}},
		{"unknown version", "unsupported version", func(t *testing.T) string {
			return recordFence(t, reviewrecord.Record{Version: 9, Spec: "specs/example/spec.md", PlanDigest: "sha256:x", ImplementationSession: "author"})
		}},
		{"unterminated fence", "invalid unterminated bench-review-record fence", func(*testing.T) string {
			return "# Review outcomes\n\n```bench-review-record\n{\n"
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, args := seedReviewEvidence(t, false)
			mustWriteFile(t, "reviews/example.md", tc.body(t))
			runGit(t, "add", ".")
			runGit(t, "commit", "-q", "-m", "retain the record under test")
			// Retaining the record moved the tip, so the charge pins the tip it
			// now reads rather than the one the seed computed.
			for i, arg := range args {
				if arg == "--source-tip" {
					args[i+1] = runGit(t, "rev-parse", "HEAD")
				}
			}
			out, code := Command(args)
			if code != 0 || !strings.Contains(out, "completion_evidence[1]") {
				t.Fatalf("the charge packet lost its completion evidence row (%d):\n%s", code, out)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("the completion evidence row does not describe %q:\n%s", tc.want, out)
			}
		})
	}
}
