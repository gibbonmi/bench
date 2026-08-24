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
	"github.com/gibbonmi/bench/internal/spec"
)

func TestLandReviewedPublishesExactSourceParentAndProspectiveSpec(t *testing.T) {
	root := fixture(t)
	write(t, root, "specs/x/spec.md", "Status: staged\nbody\n")
	git(t, root, "add", "specs/x/spec.md")
	git(t, root, "commit", "-qm", "stage spec")
	destination := git(t, root, "rev-parse", "HEAD")
	git(t, root, "switch", "-q", "-c", "reviewed-source")
	write(t, root, "reviewed", "source bytes\n")
	git(t, root, "add", "reviewed")
	git(t, root, "commit", "-qm", "reviewed work")
	source := git(t, root, "rev-parse", "HEAD")
	git(t, root, "switch", "-q", "main")
	sourceWorktree := filepath.Join(t.TempDir(), "source")
	git(t, root, "worktree", "add", "-q", sourceWorktree, "reviewed-source")
	sourceFingerprint, err := CheckoutFingerprint(sourceWorktree)
	if err != nil {
		t.Fatal(err)
	}
	destinationFingerprint, err := CheckoutFingerprint(root)
	if err != nil {
		t.Fatal(err)
	}

	o := New()
	o.authorize = func(_ context.Context, root, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		if got := git(t, root, "show", tree+":specs/x/spec.md"); got != "Status: implemented\nbody" {
			t.Fatalf("authorized spec = %q", got)
		}
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.LandReviewed(context.Background(), ReviewedRequest{
		Root: root, Destination: "refs/heads/main", DestinationBase: destination,
		Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
		SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
		SpecPath: "specs/x/spec.md", SpecBytes: []byte("Status: staged\nbody\n"), SpecMode: 0o644,
		Message: "land reviewed source",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.SourceTip != source || got.DestinationBase != destination || git(t, root, "rev-list", "--parents", "-n", "1", got.Commit) != got.Commit+" "+destination+" "+source {
		t.Fatalf("published identity = %+v", got)
	}
	if got.Tree != git(t, root, "rev-parse", got.Commit+"^{tree}") || git(t, root, "show", got.Commit+":specs/x/spec.md") != "Status: implemented\nbody" {
		t.Fatalf("published prospective tree was lost: %+v", got)
	}
}

// A spec-less reviewed landing publishes the composition unchanged: no staged spec is
// proven, neutralized, or transitioned, and the two reviewed parents still bind.
func TestLandReviewedWithoutASpecPublishesTheCompositionUntransitioned(t *testing.T) {
	root := fixture(t)
	write(t, root, "specs/x/spec.md", "Status: staged\nbody\n")
	git(t, root, "add", "specs/x/spec.md")
	git(t, root, "commit", "-qm", "stage spec")
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

	o := New()
	authorized := ""
	o.authorize = func(_ context.Context, _, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		authorized = tree
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.LandReviewed(context.Background(), ReviewedRequest{
		Root: root, Destination: "refs/heads/main", DestinationBase: destination,
		Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
		SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
		Message: "land the spec-less source",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Tree != authorized || got.Tree != git(t, root, "rev-parse", source+"^{tree}") {
		t.Fatalf("published tree = %q, authorized %q, want the composition %q", got.Tree, authorized, git(t, root, "rev-parse", source+"^{tree}"))
	}
	if git(t, root, "rev-list", "--parents", "-n", "1", got.Commit) != got.Commit+" "+destination+" "+source {
		t.Fatalf("published parents = %q", git(t, root, "rev-list", "--parents", "-n", "1", got.Commit))
	}
	if published := git(t, root, "show", got.Commit+":specs/x/spec.md"); published != "Status: staged\nbody" {
		t.Fatalf("spec-less landing transitioned the spec: %q", published)
	}
}

// FA6: the reviewed landing's close consumes the tickets-only folder from the tree it
// publishes. The folder carries a space and a glob character, so a literal removal is
// the only one that reaches it.
func TestLandReviewedClosePathPublishesATreeWithoutTheFolder(t *testing.T) {
	const closed = "specs/ft900 tickets [x]*"
	root := fixture(t)
	write(t, root, closed+"/tickets/one.md", "# One\n")
	git(t, root, "add", "--", closed)
	git(t, root, "commit", "-qm", "stage tickets")
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
	if git(t, root, "ls-tree", "-r", "--name-only", source, "--", closed) == "" {
		t.Fatal("fixture folder is untracked; the assertion below could not distinguish a close step")
	}

	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.LandReviewed(context.Background(), ReviewedRequest{
		Root: root, Destination: "refs/heads/main", DestinationBase: destination,
		Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
		SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
		ClosePath: closed, Message: "land and close the tickets-only folder",
	})
	if err != nil {
		t.Fatal(err)
	}
	if listed := git(t, root, "ls-tree", "-r", "--name-only", got.Commit, "--", closed); listed != "" {
		t.Fatalf("published tree still carries the closed folder: %q", listed)
	}
	if git(t, root, "show", got.Commit+":reviewed") != "source bytes" {
		t.Fatal("the close dropped the reviewed source content")
	}
}

func TestLandReviewedRefusesTreeEquivalentMovedSourceTip(t *testing.T) {
	root := fixture(t)
	write(t, root, "specs/x/spec.md", "Status: staged\n")
	git(t, root, "add", "specs/x/spec.md")
	git(t, root, "commit", "-qm", "stage spec")
	destination := git(t, root, "rev-parse", "HEAD")
	git(t, root, "switch", "-q", "-c", "reviewed-source")
	write(t, root, "reviewed", "source bytes\n")
	git(t, root, "add", "reviewed")
	git(t, root, "commit", "-qm", "reviewed work")
	reviewed := git(t, root, "rev-parse", "HEAD")
	git(t, root, "commit", "--allow-empty", "-qm", "tree-equivalent drift")
	git(t, root, "switch", "-q", "main")
	sourceWorktree := filepath.Join(t.TempDir(), "source")
	git(t, root, "worktree", "add", "-q", sourceWorktree, "reviewed-source")
	sourceFingerprint, err := CheckoutFingerprint(sourceWorktree)
	if err != nil {
		t.Fatal(err)
	}
	destinationFingerprint, err := CheckoutFingerprint(root)
	if err != nil {
		t.Fatal(err)
	}

	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	_, err = o.LandReviewed(context.Background(), ReviewedRequest{
		Root: root, Destination: "refs/heads/main", DestinationBase: destination,
		Source: "refs/heads/reviewed-source", SourceTip: reviewed, ReviewBase: destination,
		SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
		SpecPath: "specs/x/spec.md", SpecBytes: []byte("Status: staged\n"), SpecMode: 0o644, Message: "must not land",
	})
	if err == nil || !strings.Contains(err.Error(), "source tip moved") {
		t.Fatalf("tree-equivalent moved source = %v", err)
	}
	if got := git(t, root, "rev-parse", "refs/heads/main"); got != destination {
		t.Fatalf("destination moved to %s", got)
	}
}
func TestLandReviewedAuthorizationKindTablePreservesState(t *testing.T) {
	for _, kind := range []authorization.Kind{authorization.Green, authorization.Candidate, authorization.Inherited, authorization.Infrastructure} {
		t.Run(string(kind), func(t *testing.T) {
			root := fixture(t)
			write(t, root, "specs/x/spec.md", "Status: staged\n")
			git(t, root, "add", "specs/x/spec.md")
			git(t, root, "commit", "-qm", "stage spec")
			destination := git(t, root, "rev-parse", "HEAD")
			sourceWorktree := filepath.Join(t.TempDir(), "source")
			git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, destination)
			write(t, sourceWorktree, "reviewed", "source bytes\n")
			git(t, sourceWorktree, "add", "reviewed")
			git(t, sourceWorktree, "commit", "-qm", "reviewed")
			source := git(t, sourceWorktree, "rev-parse", "HEAD")
			sourceFingerprint, _ := CheckoutFingerprint(sourceWorktree)
			destinationFingerprint, _ := CheckoutFingerprint(root)
			o := New()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				return authorization.Result{Kind: kind}
			}
			got, err := o.LandReviewed(context.Background(), ReviewedRequest{
				Root: root, Destination: "refs/heads/main", DestinationBase: destination,
				Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
				SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
				SpecPath: "specs/x/spec.md", SpecBytes: []byte("Status: staged\n"), SpecMode: 0o644, Message: "authorize",
			})
			if kind == authorization.Green {
				if err != nil || got.Commit == "" {
					t.Fatalf("green authorization = %+v, %v", got, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), string(kind)) {
				t.Fatalf("%s authorization error = %v", kind, err)
			}
			if git(t, root, "rev-parse", "main") != destination || git(t, root, "status", "--porcelain=v1") != "" || git(t, sourceWorktree, "status", "--porcelain=v1") != "" {
				t.Fatalf("%s authorization changed repository state", kind)
			}
		})
	}
}

func TestLandReviewedDestinationCASLossThenRecomposition(t *testing.T) {
	root := fixture(t)
	write(t, root, "specs/x/spec.md", "Status: staged\n")
	git(t, root, "add", "specs/x/spec.md")
	git(t, root, "commit", "-qm", "stage spec")
	destination := git(t, root, "rev-parse", "HEAD")
	sourceWorktree := filepath.Join(t.TempDir(), "source")
	git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, destination)
	write(t, sourceWorktree, "reviewed", "source bytes\n")
	git(t, sourceWorktree, "add", "reviewed")
	git(t, sourceWorktree, "commit", "-qm", "reviewed")
	source := git(t, sourceWorktree, "rev-parse", "HEAD")

	winnerTree := git(t, root, "rev-parse", destination+"^{tree}")
	winner := git(t, root, "commit-tree", winnerTree, "-p", destination, "-m", "winner")
	calls := 0
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		calls++
		return authorization.Result{Kind: authorization.Green}
	}
	o.updateRef = func(gotRoot, ref, next, expected string) error {
		if calls == 1 {
			git(t, gotRoot, "update-ref", ref, winner, expected)
		}
		return updateRef(gotRoot, ref, next, expected)
	}
	request := func(base string) ReviewedRequest {
		sourceFingerprint, _ := CheckoutFingerprint(sourceWorktree)
		destinationFingerprint, _ := CheckoutFingerprint(root)
		return ReviewedRequest{
			Root: root, Destination: "refs/heads/main", DestinationBase: base,
			Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
			SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
			SpecPath: "specs/x/spec.md", SpecBytes: []byte("Status: staged\n"), SpecMode: 0o644, Message: "land",
		}
	}
	if _, err := o.LandReviewed(context.Background(), request(destination)); err == nil || !strings.Contains(err.Error(), "compare-and-swap") {
		t.Fatalf("CAS loser error = %v", err)
	}
	if got := git(t, root, "rev-parse", "main"); got != winner {
		t.Fatalf("winner was overwritten: %s", got)
	}
	git(t, root, "reset", "--hard", winner)
	o.updateRef = func(string, string, string, string) error { return os.ErrPermission }
	if _, err := o.LandReviewed(context.Background(), request(winner)); err == nil || strings.Contains(err.Error(), "compare-and-swap") || !strings.Contains(err.Error(), os.ErrPermission.Error()) {
		t.Fatalf("unchanged destination update error = %v", err)
	}
	o.updateRef = updateRef
	got, err := o.LandReviewed(context.Background(), request(winner))
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 || git(t, root, "rev-list", "--parents", "-n", "1", got.Commit) != got.Commit+" "+winner+" "+source {
		t.Fatalf("recomposition = %+v, authorization calls=%d", got, calls)
	}
}
func TestLandReviewedPublishesSourceSpecBytesOverAnyDestinationCopy(t *testing.T) {
	for _, tc := range []struct {
		name, baseSpec, destinationSpec, sourceSpec string
		specMode                                    os.FileMode
	}{
		{name: "destination-copy-is-stale", baseSpec: reviewedBaseSpec, sourceSpec: reviewedAmendedSpec},
		{name: "destination-moved-after-the-review-base", baseSpec: reviewedBaseSpec, destinationSpec: reviewedHeadingSpec, sourceSpec: reviewedAmendedSpec},
		{name: "destination-never-carried-the-spec", sourceSpec: reviewedBaseSpec},
		{name: "destination-overlaps-the-amendment", baseSpec: reviewedBaseSpec, destinationSpec: reviewedOverlapSpec, sourceSpec: reviewedAmendedSpec},
		{name: "source-spec-is-executable", baseSpec: reviewedBaseSpec, destinationSpec: reviewedHeadingSpec, sourceSpec: reviewedAmendedSpec, specMode: 0o755},
		{name: "source-spec-is-executable-and-the-destination-never-carried-it", sourceSpec: reviewedBaseSpec, specMode: 0o755},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newReviewedLanding(t, tc.baseSpec, tc.destinationSpec, tc.sourceSpec, tc.specMode)
			o := New()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				return authorization.Result{Kind: authorization.Green}
			}
			got, err := o.LandReviewed(context.Background(), f.request(t, tc.sourceSpec, "land the amended source"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := spec.Implemented([]byte(tc.sourceSpec))
			if err != nil {
				t.Fatal(err)
			}
			if published := gitBytes(t, f.root, "show", got.Commit+":"+reviewedFixtureSpecPath); !bytes.Equal(published, want) {
				t.Fatalf("published spec = %q, want %q", published, want)
			}
			wantMode := "100644"
			if f.specMode.Perm()&0o111 != 0 {
				wantMode = "100755"
			}
			if mode := gitMode(t, f.root, "ls-tree", got.Commit, "--", reviewedFixtureSpecPath); mode != wantMode {
				t.Fatalf("published spec mode = %q, want %q", mode, wantMode)
			}
			if got.DestinationBase != f.destination {
				t.Fatalf("destination base = %s, want the real destination %s", got.DestinationBase, f.destination)
			}
			if parents := git(t, f.root, "rev-list", "--parents", "-n", "1", got.Commit); parents != got.Commit+" "+f.destination+" "+f.source {
				t.Fatalf("published parents = %q, want %s and %s", parents, f.destination, f.source)
			}
		})
	}
}

func TestLandReviewedSpecStateTableRefusesBeforeAuthorization(t *testing.T) {
	for _, tc := range []struct{ name, sourceSpec, suppliedSpec, want string }{
		{name: "two-staged-lines", sourceSpec: "Status: staged\nStatus: staged\n", want: "more than one Status: staged line"},
		{name: "absent-at-the-source-tip", want: specProvenanceRefusal},
		{name: "already-implemented", sourceSpec: "Status: implemented\n", want: "no Status: staged line"},
		{name: "supplied-bytes-no-commit-carries", sourceSpec: "Status: staged\nsource\n", suppliedSpec: "Status: staged\nnever committed\n", want: specProvenanceRefusal},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newReviewedLanding(t, "", "", tc.sourceSpec, 0)
			supplied := tc.suppliedSpec
			if supplied == "" {
				supplied = tc.sourceSpec
			}
			calls := 0
			o := New()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				calls++
				return authorization.Result{Kind: authorization.Green}
			}
			_, err := o.LandReviewed(context.Background(), f.request(t, supplied, "must refuse"))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("spec state %s = %v, want %q", tc.name, err, tc.want)
			}
			if calls != 0 || git(t, f.root, "rev-parse", "main") != f.destination {
				t.Fatalf("spec state %s authorized %d times or moved the destination", tc.name, calls)
			}
		})
	}
}
