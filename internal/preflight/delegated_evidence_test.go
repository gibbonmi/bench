package preflight

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
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

// The prepared review metadata describes both valid record versions, and it keeps naming
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
			root, _, args := seedReviewEvidence(t, false)
			preflighttest.MustWriteFile(t, "reviews/example.md", tc.body(t))
			preflighttest.RunGit(t, "add", ".")
			preflighttest.RunGit(t, "commit", "-q", "-m", "retain the record under test")
			// Retaining the record moved the tip, so the preparation pins the tip it
			// now reads rather than the one the seed computed.
			for i, arg := range args {
				if arg == "--source-tip" {
					args[i+1] = preflighttest.RunGit(t, "rev-parse", "HEAD")
				}
			}
			out, code := Command(args)
			if code != 0 {
				t.Fatalf("review preparation = (%d):\n%s", code, out)
			}
			completion := preflighttest.PublishedPack(t, root, preparedIdentity(t, out)).Metadata().Completion
			if len(completion) != 1 {
				t.Fatalf("completion rows = %d, want one", len(completion))
			}
			row := completion[0]
			if !strings.Contains(row.RecordState+" "+row.Detail, tc.want) {
				t.Fatalf("the completion row does not describe %q: %#v", tc.want, row)
			}
		})
	}
}
