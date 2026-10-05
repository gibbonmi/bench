package landing

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// recordingAdmission appends each decision the admitted route asks for to calls.
type recordingAdmission struct {
	calls          *[]string
	check, publish error
}

func (a recordingAdmission) Check(string) error {
	*a.calls = append(*a.calls, "check")
	return a.check
}

func (a recordingAdmission) Publish(_ string, publish func() error) error {
	*a.calls = append(*a.calls, "publish")
	if a.publish != nil {
		return a.publish
	}
	return publish()
}

// The admitted route decides the exact composed tree before the gate and again around the
// ref update. A missing admission refuses before composition, and a refusal at either
// decision leaves the destination where it was.
func TestLandAdmittedDecidesAroundTheGate(t *testing.T) {
	refused := errors.New("admission refused")
	for _, row := range []struct {
		name      string
		admission func(*[]string) Admission
		calls     string
		published bool
	}{
		{name: "missing", admission: func(*[]string) Admission { return nil }},
		{name: "check-refuses", admission: func(calls *[]string) Admission { return recordingAdmission{calls: calls, check: refused} }, calls: "check"},
		{name: "publish-refuses", admission: func(calls *[]string) Admission { return recordingAdmission{calls: calls, publish: refused} }, calls: "check,gate,publish"},
		{name: "admitted", admission: func(calls *[]string) Admission { return recordingAdmission{calls: calls} }, calls: "check,gate,publish", published: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := fixture(t)
			destination := git(t, root, "rev-parse", "HEAD")
			sourceWorktree := filepath.Join(t.TempDir(), "source")
			git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, destination)
			write(t, sourceWorktree, "reviewed", "source bytes\n")
			git(t, sourceWorktree, "add", "reviewed")
			git(t, sourceWorktree, "commit", "-qm", "reviewed work")
			source := git(t, sourceWorktree, "rev-parse", "HEAD")
			sourceFingerprint, err := CheckoutFingerprint(sourceWorktree)
			if err != nil {
				t.Fatal(err)
			}
			destinationFingerprint, err := CheckoutFingerprint(root)
			if err != nil {
				t.Fatal(err)
			}
			var calls []string
			o := New()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				calls = append(calls, "gate")
				return authorization.Result{Kind: authorization.Green}
			}
			got, err := o.LandAdmitted(context.Background(), ReviewedRequest{
				Root: root, Destination: "refs/heads/main", DestinationBase: destination,
				Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
				SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
				Message: "land the admitted source",
			}, row.admission(&calls))
			if strings.Join(calls, ",") != row.calls {
				t.Fatalf("decisions = %q, want %q", strings.Join(calls, ","), row.calls)
			}
			main := git(t, root, "rev-parse", "main")
			if row.published {
				if err != nil || main != got.Commit {
					t.Fatalf("admitted landing = (%+v, %v), main=%s", got, err, main)
				}
				return
			}
			if err == nil || main != destination {
				t.Fatalf("refused landing = (%+v, %v), main=%s, want an error and main=%s", got, err, main, destination)
			}
		})
	}
}
