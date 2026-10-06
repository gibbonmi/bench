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
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/testrepo"
)

type completionFixture struct {
	evidence   *recordtest.Fixture
	root, base string
}

func newCompletionFixture(t *testing.T, prepare ...func(*recordtest.Fixture)) completionFixture {
	return attachedCompletionFixture(t, recordtest.Attach, prepare...)
}

// attachedCompletionFixture builds a two-chunk landing fixture from either
// record form: the oracle, the ignore rule, the evidence worktree, and the
// chunk sequence. The version 1 and delegated landings share this one
// sequence, so a change to the oracle or the build order cannot drift between
// them.
func attachedCompletionFixture(t *testing.T, attach func(testing.TB, string, int) *recordtest.Fixture, prepare ...func(*recordtest.Fixture)) completionFixture {
	t.Helper()
	f := attach(t, fixture(t), 2)
	f.Write(".gitignore", ".logs/\n")
	g := testrepo.NewGateFixture(t.TempDir())
	if err := g.Write(f.Root, g.Command("grep")+" -q '^Status: implemented$' specs/example/spec.md\n", ""); err != nil {
		t.Fatal(err)
	}
	f.Commit("landing oracle")
	root, base := f.Root, f.Tip()
	source := filepath.Join(t.TempDir(), "source")
	f.Git("worktree", "add", "-qb", "evidence-source", source, base)
	f.Root = source
	for _, apply := range prepare {
		apply(f)
	}
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

// The reviewed landing binds the gate's completion oracle to a tickets-only close. The
// real gate grades the exact close green, and it refuses an authorized tree that keeps
// the closed folder, so the landing does not publish that tree.
func TestLandingTicketsOnlyCompletion(t *testing.T) {
	const closed = "specs/t"
	for _, row := range []struct {
		name, reason string
		keep         bool
	}{
		{name: "exact-close"},
		{name: "kept-folder", reason: "completion keeps closed " + closed + "/tickets/one.md", keep: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := fixture(t)
			write(t, root, closed+"/tickets/one.md", "# One\n")
			g := testrepo.NewGateFixture(t.TempDir())
			if err := g.Write(root, "exit 0\n", ""); err != nil {
				t.Fatal(err)
			}
			git(t, root, "add", ".")
			git(t, root, "commit", "-qm", "stage tickets and the landing oracle")
			destination := git(t, root, "rev-parse", "HEAD")
			sourceWorktree := filepath.Join(t.TempDir(), "source")
			git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, destination)
			write(t, sourceWorktree, "reviewed", "source bytes\n")
			git(t, sourceWorktree, "add", "reviewed")
			git(t, sourceWorktree, "commit", "-qm", "reviewed work")
			sourceFingerprint, err := CheckoutFingerprint(sourceWorktree)
			if err != nil {
				t.Fatal(err)
			}
			destinationFingerprint, err := CheckoutFingerprint(root)
			if err != nil {
				t.Fatal(err)
			}
			owner := New()
			if row.keep {
				owner.authorize = func(ctx context.Context, root, tree string, stdout, stderr io.Writer) authorization.Result {
					kept, err := replaceTreeFile(root, tree, closed+"/tickets/one.md", []byte("# One\n"), 0o644)
					if err != nil {
						t.Fatal(err)
					}
					return authorization.AuthorizeWithWriters(ctx, root, kept, stdout, stderr)
				}
			}
			_, err = owner.LandReviewed(t.Context(), ReviewedRequest{
				Root: root, Destination: "refs/heads/main", DestinationBase: destination,
				Source: "refs/heads/reviewed-source", SourceTip: git(t, sourceWorktree, "rev-parse", "HEAD"), ReviewBase: destination,
				SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
				ClosePath: closed, Message: "land the tickets-only close", Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{},
			})
			if row.reason == "" {
				if err != nil || git(t, root, "rev-parse", "main") == destination {
					t.Fatalf("exact tickets-only close = %v, want a publication", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.reason) {
				t.Fatalf("kept folder = %v, want a refusal naming %q", err, row.reason)
			}
			if got := git(t, root, "rev-parse", "main"); got != destination {
				t.Fatalf("refused close moved main to %s", got)
			}
		})
	}
}

func TestLandingCompletionEvidence(t *testing.T) {
	cases := []struct {
		name, reason string
		change       func(*testing.T, completionFixture, *Owner)
	}{
		{"complete", "", func(*testing.T, completionFixture, *Owner) {}},
		{"comment-only gap", "", func(t *testing.T, f completionFixture, _ *Owner) {
			f.evidence.Write("comment_gap_test.go", "package fixture\nvar value = 1 // after\n")
			f.evidence.Commit("correct reviewed comment")
			digest, err := rr.SourceDigest(f.evidence.Root, f.evidence.Tree(), recordtest.Spec)
			if err != nil {
				t.Fatal(err)
			}
			f.evidence.Record.Completion.SourceDigest = digest
			f.evidence.Record.Completion.Verification = f.evidence.Verification("corrected-final", digest, f.evidence.Plan.FinalVerification)
			f.evidence.Save()
		}},
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
			var prepare []func(*recordtest.Fixture)
			if tc.name == "comment-only gap" {
				prepare = append(prepare, func(f *recordtest.Fixture) {
					f.Write("comment_gap_test.go", "package fixture\nvar value = 1 // before\n")
				})
			}
			f := newCompletionFixture(t, prepare...)
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
