package landing

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/spec"
)

type completionFixture struct {
	evidence   *recordtest.Fixture
	root, base string
}

func newCompletionFixture(t *testing.T) completionFixture {
	return attachedCompletionFixture(t, recordtest.Attach)
}

// attachedCompletionFixture builds a two-chunk landing fixture from either
// record form: the oracle, the ignore rule, the evidence worktree, and the
// chunk sequence. The version 1 and delegated landings share this one
// sequence, so a change to the oracle or the build order cannot drift between
// them.
func attachedCompletionFixture(t *testing.T, attach func(testing.TB, string, int) *recordtest.Fixture) completionFixture {
	t.Helper()
	f := attach(t, fixture(t), 2)
	f.Write(".gitignore", ".logs/\n")
	f.Write(".bench/gate.sh", "#!/bin/sh\ngrep -q '^Status: implemented$' specs/example/spec.md\n")
	if err := os.Chmod(filepath.Join(f.Root, ".bench/gate.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	f.Write(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.Commit("landing oracle")
	root, base := f.Root, f.Tip()
	source := filepath.Join(t.TempDir(), "source")
	f.Git("worktree", "add", "-qb", "evidence-source", source, base)
	f.Root = source
	f.AddChunk()
	f.Save()
	f.Commit("retain first chunk")
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("retain completion")
	return completionFixture{f, root, base}
}

func (f completionFixture) request(t *testing.T) ReviewedRequest {
	t.Helper()
	sourceFingerprint, err := CheckoutFingerprint(f.evidence.Root)
	if err != nil {
		t.Fatal(err)
	}
	destinationFingerprint, err := CheckoutFingerprint(f.root)
	if err != nil {
		t.Fatal(err)
	}
	return ReviewedRequest{
		Root: f.root, Destination: "refs/heads/main", DestinationBase: git(t, f.root, "rev-parse", "main"),
		Source: "refs/heads/evidence-source", SourceTip: f.evidence.Tip(), ReviewBase: f.base,
		SourceWorktree: f.evidence.Root, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
		SpecPath: f.evidence.Record.Spec, SpecBytes: mustRead(t, filepath.Join(f.evidence.Root, f.evidence.Record.Spec)), SpecMode: 0644,
		Message: "land complete evidence", Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{},
	}
}

func TestLandingCompletionEvidence(t *testing.T) {
	cases := []struct {
		name, reason string
		change       func(*testing.T, completionFixture, *Owner)
	}{
		{"complete", "", func(*testing.T, completionFixture, *Owner) {}},
		{"missing reconciliation", "reconciliation E1", func(t *testing.T, f completionFixture, _ *Owner) {
			delete(f.evidence.Record.Completion.Reconciliation, "E1")
			f.evidence.Save()
		}},
		{"missing record", "retain a valid native result record", func(t *testing.T, f completionFixture, _ *Owner) {
			if err := os.Remove(filepath.Join(f.evidence.Root, "reviews/example.md")); err != nil {
				t.Fatal(err)
			}
		}},
		{"missing integration", "verification integration", func(t *testing.T, f completionFixture, _ *Owner) {
			f.evidence.Record.Completion.Verification = f.evidence.Record.Completion.Verification[:1]
			f.evidence.Save()
		}},
		{"failed integration", "verification integration", func(t *testing.T, f completionFixture, _ *Owner) {
			one := 1
			f.evidence.Record.Completion.Verification[1].ExitCode = &one
			f.evidence.Save()
		}},
		{"stale integration", "verification integration", func(t *testing.T, f completionFixture, _ *Owner) {
			f.evidence.Record.Completion.Verification[1].SourceDigest = f.base
			f.evidence.Save()
		}},
		{"pending integration", "verification integration", func(t *testing.T, f completionFixture, _ *Owner) {
			f.evidence.Record.Completion.Verification[1].State = "pending"
			f.evidence.Save()
		}},
		{"wrong final verifier", "verification integration", func(t *testing.T, f completionFixture, _ *Owner) {
			f.evidence.Record.Completion.Verification[1].Performer = "someone-else"
			f.evidence.Save()
		}},
		{"destination addition", "completion composition changes added.txt", func(t *testing.T, f completionFixture, _ *Owner) {
			write(t, f.root, "added.txt", "unreviewed destination content\n")
			git(t, f.root, "add", "added.txt")
			git(t, f.root, "commit", "-qm", "destination changed")
		}},
		{"destination deletion", "completion composition changes foreign", func(t *testing.T, f completionFixture, _ *Owner) {
			git(t, f.root, "rm", "foreign")
			git(t, f.root, "commit", "-qm", "destination removed content")
		}},
		{"extra spec bytes", "exact status transform", func(t *testing.T, f completionFixture, owner *Owner) {
			owner.authorize = func(ctx context.Context, root, tree string, stdout, stderr io.Writer) authorization.Result {
				changed := append(gitBytes(t, root, "show", tree+":"+f.evidence.Record.Spec), []byte("unreviewed acceptance change\n")...)
				wrong, err := replaceTreeFile(root, tree, f.evidence.Record.Spec, changed, 0644)
				if err != nil {
					t.Fatal(err)
				}
				return authorization.AuthorizeWithWriters(ctx, root, wrong, stdout, stderr)
			}
		}},
		{"spec mode change", "exact status transform", func(t *testing.T, f completionFixture, owner *Owner) {
			owner.authorize = func(ctx context.Context, root, tree string, stdout, stderr io.Writer) authorization.Result {
				wrong, err := replaceTreeFile(root, tree, f.evidence.Record.Spec, gitBytes(t, root, "show", tree+":"+f.evidence.Record.Spec), 0755)
				if err != nil {
					t.Fatal(err)
				}
				return authorization.AuthorizeWithWriters(ctx, root, wrong, stdout, stderr)
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newCompletionFixture(t)
			owner := New()
			tc.change(t, f, &owner)
			if f.evidence.Git("status", "--porcelain") != "" {
				f.evidence.Commit("evidence case")
			}
			request := f.request(t)
			result, err := owner.LandReviewed(t.Context(), request)
			if tc.reason == "" {
				if err != nil {
					t.Fatalf("complete evidence refused: %v", err)
				}
				want, err := spec.Implemented(request.SpecBytes)
				if err != nil {
					t.Fatal(err)
				}
				if got := gitBytes(t, f.root, "show", result.Commit+":"+request.SpecPath); !bytes.Equal(got, want) {
					t.Fatalf("published status transform = %q", got)
				}
				if got := git(t, f.root, "rev-parse", "main"); got != result.Commit {
					t.Fatalf("publication missing: %s", got)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("incomplete evidence published or lost reason %q: %v", tc.reason, err)
			}
			if got := git(t, f.root, "rev-parse", "main"); got != request.DestinationBase {
				t.Fatalf("ref moved before proof: %s", got)
			}
		})
	}
}
