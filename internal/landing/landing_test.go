package landing

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/spec"
)

func TestLandStreamsProspectiveGateOutputOnce(t *testing.T) {
	root := fixture(t)
	write(t, root, ".bench/gate.sh", "#!/bin/sh\necho 'gate stdout'\necho 'gate stderr' >&2\ntest ! -f fail\n")
	if err := os.Chmod(filepath.Join(root, ".bench", "gate.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	git(t, root, "add", ".bench")
	git(t, root, "commit", "-qm", "gate")
	write(t, root, "named", "changed")
	base := git(t, root, "rev-parse", "HEAD")
	var stdout, stderr bytes.Buffer

	if _, err := New().Land(t.Context(), Request{
		Root: root, Destination: "refs/heads/main", Expected: base, Message: "named", Paths: []string{"named"},
		Stdout: &stdout, Stderr: &stderr,
	}); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "gate stdout\n" || stderr.String() != "gate stderr\n" {
		t.Fatalf("gate output stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRuntimeIgnoredPathRejectsUnsafeOrCollapsedNames(t *testing.T) {
	for _, tc := range []struct {
		path string
		want bool
	}{
		{path: ".logs", want: true},
		{path: ".logs/", want: true},
		{path: ".logs/gate.jsonl", want: true},
		{path: ".logs-foreign/gate.jsonl"},
		{path: ".logs/../foreign"},
		{path: ".logs//gate.jsonl"},
		{path: ".logs/\x1b"},
	} {
		t.Run(strings.ReplaceAll(tc.path, "\x1b", "escape"), func(t *testing.T) {
			if got := RuntimeIgnoredPath(tc.path); got != tc.want {
				t.Fatalf("RuntimeIgnoredPath(%q) = %t, want %t", tc.path, got, tc.want)
			}
		})
	}
}

func TestLandComposesOnlyNamedPathsAndCASesExpectedBase(t *testing.T) {
	root := fixture(t)
	write(t, root, "named", "changed")
	write(t, root, "foreign", "staged")
	git(t, root, "add", "foreign")
	base := git(t, root, "rev-parse", "HEAD")
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "named", Paths: []string{"named", "named"}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Tree != git(t, root, "rev-parse", got.Commit+"^{tree}") {
		t.Fatal("landing commit did not use authorized tree")
	}
	if git(t, root, "show", got.Commit+":foreign") != "base" {
		t.Fatal("foreign staged content entered landing")
	}
	if git(t, root, "show", ":foreign") != "staged" {
		t.Fatal("foreign index content was not preserved")
	}
}

func TestLandReconcilesAuthorizedNamedSnapshotAfterAuthorizationMutation(t *testing.T) {
	root := raceFixture(t)
	authorized := []byte("authorized\n")
	write(t, root, "named", string(authorized))
	unnamed, before := dirtyUnnamedState(t, root)
	base := git(t, root, "rev-parse", "HEAD")

	o := New()
	var authorizedTree string
	o.authorize = func(_ context.Context, root, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		authorizedTree = tree
		if got := gitBytes(t, root, "show", tree+":named"); !bytes.Equal(got, authorized) {
			t.Fatalf("authorized named bytes = %q", got)
		}
		write(t, root, "named", "mutated-after-authorization\n")
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "named", Paths: []string{"named"}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Tree != authorizedTree || git(t, root, "rev-parse", got.Commit+"^{tree}") != authorizedTree {
		t.Fatal("winning commit did not retain the authorized tree")
	}
	for source, content := range map[string][]byte{
		"commit":   gitBytes(t, root, "show", got.Commit+":named"),
		"index":    gitBytes(t, root, "show", ":named"),
		"worktree": mustRead(t, filepath.Join(root, "named")),
	} {
		if !bytes.Equal(content, authorized) {
			t.Errorf("%s named bytes = %q", source, content)
		}
	}
	if status := git(t, root, "status", "--porcelain=v1", "--", "named"); status != "" {
		t.Fatalf("named status = %q", status)
	}
	if after := snapshotPaths(t, root, unnamed...); !reflect.DeepEqual(after, before) {
		t.Fatalf("unnamed state changed\nbefore: %#v\nafter:  %#v", before, after)
	}
}

func TestLandReconcilesAuthorizedSpecSnapshotAfterAuthorizationMutation(t *testing.T) {
	root := raceFixture(t)
	write(t, root, "named", "authorized work\n")
	specPath := filepath.Join(root, "specs", "x", "spec.md")
	write(t, root, "specs/x/spec.md", "Status: staged\nauthorized body\n")
	if err := os.Chmod(specPath, 0o750); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", "specs/x/spec.md")
	unnamed, before := dirtyUnnamedState(t, root)
	base := git(t, root, "rev-parse", "HEAD")
	authorized := []byte("Status: implemented\nauthorized body\n")

	o := New()
	var authorizedTree string
	o.authorize = func(_ context.Context, root, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		authorizedTree = tree
		if got := gitBytes(t, root, "show", tree+":specs/x/spec.md"); !bytes.Equal(got, authorized) {
			t.Fatalf("authorized spec bytes = %q", got)
		}
		if mode := gitMode(t, root, "ls-tree", tree, "--", "specs/x/spec.md"); mode != "100755" {
			t.Fatalf("authorized spec mode = %q", mode)
		}
		write(t, root, "specs/x/spec.md", "Status: staged\nmutated after authorization\n")
		if err := os.Chmod(specPath, 0o600); err != nil {
			t.Fatal(err)
		}
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "spec", Paths: []string{"named"}, Spec: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Tree != authorizedTree || git(t, root, "rev-parse", got.Commit+"^{tree}") != authorizedTree {
		t.Fatal("winning commit did not retain the authorized tree")
	}
	for source, content := range map[string][]byte{
		"commit":   gitBytes(t, root, "show", got.Commit+":specs/x/spec.md"),
		"index":    gitBytes(t, root, "show", ":specs/x/spec.md"),
		"worktree": mustRead(t, specPath),
	} {
		if !bytes.Equal(content, authorized) {
			t.Errorf("%s spec bytes = %q", source, content)
		}
	}
	if mode := gitMode(t, root, "ls-tree", got.Commit, "--", "specs/x/spec.md"); mode != "100755" {
		t.Errorf("published spec mode = %q", mode)
	}
	if mode := gitMode(t, root, "ls-files", "--stage", "--", "specs/x/spec.md"); mode != "100755" {
		t.Errorf("index spec mode = %q", mode)
	}
	info, statErr := os.Stat(specPath)
	if statErr != nil {
		t.Fatal(statErr)
	}
	if info.Mode().Perm() != 0o750 {
		t.Errorf("worktree spec permissions = %v", info.Mode().Perm())
	}
	if status := git(t, root, "status", "--porcelain=v1", "--", "named", "specs/x/spec.md"); status != "" {
		t.Errorf("named status = %q", status)
	}
	if after := snapshotPaths(t, root, unnamed...); !reflect.DeepEqual(after, before) {
		t.Fatalf("unnamed state changed\nbefore: %#v\nafter:  %#v", before, after)
	}
}

func TestLandReconcilesLateUntrackedDescendantUnderNamedDirectory(t *testing.T) {
	root := raceFixture(t)
	write(t, root, "owned/tracked", "base\n")
	git(t, root, "add", "owned/tracked")
	git(t, root, "commit", "-qm", "directory base")
	write(t, root, "owned/tracked", "authorized\n")
	unnamed, before := dirtyUnnamedState(t, root)
	base := git(t, root, "rev-parse", "HEAD")

	authorizing := make(chan string, 1)
	green := make(chan struct{})
	o := New()
	o.authorize = func(_ context.Context, _ string, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		authorizing <- tree
		<-green
		return authorization.Result{Kind: authorization.Green}
	}
	results := make(chan Result, 1)
	errors := make(chan error, 1)
	go func() {
		got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "directory", Paths: []string{"owned"}})
		results <- got
		errors <- err
	}()

	authorizedTree := <-authorizing
	if content := gitBytes(t, root, "show", authorizedTree+":owned/tracked"); !bytes.Equal(content, []byte("authorized\n")) {
		t.Fatalf("authorized tracked bytes = %q", content)
	}
	if names := git(t, root, "ls-tree", "-r", "--name-only", authorizedTree, "--", "owned"); names != "owned/tracked" {
		t.Fatalf("authorized directory entries = %q", names)
	}
	latePath := filepath.Join(root, "owned", "late")
	write(t, root, "owned/late", "created during authorization\n")
	close(green)

	got, err := <-results, <-errors
	if err != nil {
		t.Fatal(err)
	}
	if got.Tree != authorizedTree || git(t, root, "rev-parse", got.Commit+"^{tree}") != authorizedTree {
		t.Fatal("winning commit did not retain the authorized tree")
	}
	if _, statErr := os.Lstat(latePath); !os.IsNotExist(statErr) {
		t.Fatalf("late descendant remains: %v", statErr)
	}
	git(t, root, "diff", "--quiet", "--cached", got.Commit, "--", "owned")
	git(t, root, "diff", "--quiet", got.Commit, "--", "owned")
	if status := git(t, root, "status", "--porcelain=v1", "--", "owned"); status != "" {
		t.Fatalf("named directory status = %q", status)
	}
	if after := snapshotPaths(t, root, unnamed...); !reflect.DeepEqual(after, before) {
		t.Fatalf("unnamed state changed\nbefore: %#v\nafter:  %#v", before, after)
	}
}

func TestLandPreservesWinnerOnCASLoss(t *testing.T) {
	root := fixture(t)
	write(t, root, "named", "changed")
	base := git(t, root, "rev-parse", "HEAD")
	winner := git(t, root, "commit-tree", git(t, root, "rev-parse", base+"^{tree}"), "-p", base, "-m", "winner")
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	o.updateRef = func(root, ref, new, old string) error {
		git(t, root, "update-ref", ref, winner, old)
		return updateRef(root, ref, new, old)
	}
	_, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "loser", Paths: []string{"named"}})
	if err == nil || !strings.Contains(err.Error(), "compare-and-swap") {
		t.Fatalf("CAS loss = %v", err)
	}
	if git(t, root, "rev-parse", "HEAD") != winner {
		t.Fatal("loser overwrote winner")
	}
	if string(mustRead(t, filepath.Join(root, "named"))) != "changed" {
		t.Fatal("losing bytes were not preserved")
	}
}

func TestLandRetainsInfrastructureRefUpdateFailureWhenDestinationIsUnchanged(t *testing.T) {
	root := fixture(t)
	write(t, root, "named", "changed")
	base := git(t, root, "rev-parse", "HEAD")
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	o.updateRef = func(string, string, string, string) error { return os.ErrPermission }

	_, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "blocked", Paths: []string{"named"}})
	if err == nil || strings.Contains(err.Error(), "compare-and-swap") || !strings.Contains(err.Error(), os.ErrPermission.Error()) {
		t.Fatalf("unchanged destination update error = %v", err)
	}
	if got := git(t, root, "rev-parse", "HEAD"); got != base {
		t.Fatalf("destination changed to %s", got)
	}
}

func TestLandClassifiesHeldDestinationRefLockAsInfrastructure(t *testing.T) {
	root := fixture(t)
	write(t, root, "named", "changed")
	base := git(t, root, "rev-parse", "HEAD")
	lock := filepath.Join(root, ".git", "refs", "heads", "main.lock")
	if err := os.WriteFile(lock, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(lock) })
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}

	_, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "locked", Paths: []string{"named"}})
	if err == nil || strings.Contains(err.Error(), "compare-and-swap") {
		t.Fatalf("held destination ref lock error = %v", err)
	}
	if got := git(t, root, "rev-parse", "HEAD"); got != base {
		t.Fatalf("destination changed to %s", got)
	}
}

func TestLandRefusesEveryNonGreenAuthorizationWithoutMutation(t *testing.T) {
	for _, kind := range []authorization.Kind{authorization.Candidate, authorization.Inherited, authorization.Infrastructure} {
		t.Run(string(kind), func(t *testing.T) {
			root := fixture(t)
			write(t, root, "named", "changed")
			base := git(t, root, "rev-parse", "HEAD")
			before := git(t, root, "status", "--porcelain=v1")
			o := New()
			calls := 0
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				calls++
				return authorization.Result{Kind: kind}
			}
			if _, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "x", Paths: []string{"named"}}); err == nil {
				t.Fatal("Land succeeded")
			}
			if calls != 1 || git(t, root, "rev-parse", "HEAD") != base || git(t, root, "status", "--porcelain=v1") != before {
				t.Fatal("authorization refusal mutated checkout")
			}
		})
	}
}

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

func TestLandReviewedRechecksCheckoutFingerprintsAfterAuthorization(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		mutate     func(*testing.T, string, string)
	}{
		{name: "destination", want: "destination checkout changed", mutate: func(t *testing.T, root, _ string) { write(t, root, "late", "destination dirt") }},
		{name: "source", want: "source checkout changed", mutate: func(t *testing.T, _, source string) { write(t, source, "late", "source dirt") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
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
				tc.mutate(t, root, sourceWorktree)
				return authorization.Result{Kind: authorization.Green}
			}
			_, err := o.LandReviewed(context.Background(), ReviewedRequest{
				Root: root, Destination: "refs/heads/main", DestinationBase: destination,
				Source: "refs/heads/reviewed-source", SourceTip: source, ReviewBase: destination,
				SourceWorktree: sourceWorktree, SourceFingerprint: sourceFingerprint, DestinationFingerprint: destinationFingerprint,
				SpecPath: "specs/x/spec.md", SpecBytes: []byte("Status: staged\n"), SpecMode: 0o644, Message: "must refuse",
			})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("fingerprint mutation error = %v, want %q", err, tc.want)
			}
			if got := git(t, root, "rev-parse", "main"); got != destination {
				t.Fatalf("destination moved to %s", got)
			}
		})
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

func TestLandReportsPublishedCommitWhenReconciliationFails(t *testing.T) {
	root := fixture(t)
	write(t, root, "named", "changed")
	base := git(t, root, "rev-parse", "HEAD")
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	o.reconcile = func(Request, []string, composedSnapshot) error { return os.ErrPermission }
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "x", Paths: []string{"named"}})
	if err == nil || !strings.Contains(err.Error(), "landed-but-checkout-incomplete") {
		t.Fatalf("err = %v", err)
	}
	if got.Base != base || got.Commit == "" || got.Tree == "" || git(t, root, "rev-parse", "HEAD") != got.Commit {
		t.Fatal("published identity was not retained")
	}
}

func TestLandPublishesExecutableSpecModeAndReconcilesClean(t *testing.T) {
	root := fixture(t)
	git(t, root, "config", "core.filemode", "false")
	write(t, root, "work", "changed")
	specPath := filepath.Join(root, "specs", "x", "spec.md")
	write(t, root, "specs/x/spec.md", "Status: staged\nbody\n")
	if err := os.Chmod(specPath, 0o750); err != nil {
		t.Fatal(err)
	}
	base := git(t, root, "rev-parse", "HEAD")
	o := New()
	var authorizedContent, authorizedMode string
	o.authorize = func(_ context.Context, root, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		authorizedContent = git(t, root, "show", tree+":specs/x/spec.md")
		authorizedMode = gitMode(t, root, "ls-tree", tree, "--", "specs/x/spec.md")
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "spec", Paths: []string{"work"}, Spec: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if authorizedContent != "Status: implemented\nbody" || authorizedMode != "100755" {
		t.Fatalf("authorized content=%q mode=%q", authorizedContent, authorizedMode)
	}
	if mode := gitMode(t, root, "ls-tree", got.Commit, "--", "specs/x/spec.md"); mode != "100755" {
		t.Fatalf("published mode = %q", mode)
	}
	if mode := gitMode(t, root, "ls-files", "--stage", "--", "specs/x/spec.md"); mode != "100755" {
		t.Fatalf("index mode = %q", mode)
	}
	info, err := os.Stat(specPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o750 {
		t.Fatalf("worktree mode = %v", info.Mode().Perm())
	}
	if status := git(t, root, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("checkout status = %q", status)
	}
}

func TestLandSpecTransitionIsAuthorizedBeforePublicationAndPreservesModeOnLoss(t *testing.T) {
	root := fixture(t)
	write(t, root, "work", "changed")
	specPath := filepath.Join(root, "specs", "x", "spec.md")
	write(t, root, "specs/x/spec.md", "Status: staged\nbody\n")
	if err := os.Chmod(specPath, 0o600); err != nil {
		t.Fatal(err)
	}
	base := git(t, root, "rev-parse", "HEAD")
	seen, seenMode := "", ""
	o := New()
	o.authorize = func(_ context.Context, root, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		seen = git(t, root, "show", tree+":specs/x/spec.md")
		seenMode = gitMode(t, root, "ls-tree", tree, "--", "specs/x/spec.md")
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "spec", Paths: []string{"work"}, Spec: "x"})
	if err != nil || seen != "Status: implemented\nbody" || seenMode != "100644" || git(t, root, "show", got.Commit+":specs/x/spec.md") != seen {
		t.Fatalf("transition: content=%q mode=%q err=%v", seen, seenMode, err)
	}
	if mode := gitMode(t, root, "ls-tree", got.Commit, "--", "specs/x/spec.md"); mode != "100644" {
		t.Fatalf("published mode = %q", mode)
	}
	if mode := gitMode(t, root, "ls-files", "--stage", "--", "specs/x/spec.md"); mode != "100644" {
		t.Fatalf("index mode = %q", mode)
	}
	info, err := os.Stat(specPath)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatal("spec mode changed")
	}
	if string(mustRead(t, specPath)) != "Status: implemented\nbody\n" {
		t.Fatal("checkout did not reconcile")
	}
	if status := git(t, root, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("checkout status = %q", status)
	}
	root = fixture(t)
	write(t, root, "work", "changed")
	write(t, root, "specs/x/spec.md", "Status: staged\n")
	base = git(t, root, "rev-parse", "HEAD")
	winner := git(t, root, "commit-tree", git(t, root, "rev-parse", base+"^{tree}"), "-p", base, "-m", "winner")
	o = New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	o.updateRef = func(root, ref, new, old string) error {
		git(t, root, "update-ref", ref, winner, old)
		return updateRef(root, ref, new, old)
	}
	if _, err = o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "loss", Paths: []string{"work"}, Spec: "x"}); err == nil {
		t.Fatal("CAS loss succeeded")
	}
	if string(mustRead(t, filepath.Join(root, "specs/x/spec.md"))) != "Status: staged\n" {
		t.Fatal("CAS loss flipped checkout spec")
	}
}

func TestLandDetachedHEADHasOneExpectedParent(t *testing.T) {
	root := fixture(t)
	base := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "--detach", "-q")
	write(t, root, "named", "changed")
	o := New()
	var authorized string
	o.authorize = func(_ context.Context, _ string, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		authorized = tree
		return authorization.Result{Kind: authorization.Green}
	}
	got, err := o.Land(context.Background(), Request{Root: root, Destination: "HEAD", Expected: base, Message: "detached", Paths: []string{"named"}})
	if err != nil || got.Base != base || got.Tree != authorized || git(t, root, "rev-parse", "HEAD") != got.Commit || git(t, root, "rev-list", "--parents", "-n", "1", got.Commit) != got.Commit+" "+base {
		t.Fatalf("detached result %+v %v", got, err)
	}
}

func TestLandPreservesCompleteUnnamedState(t *testing.T) {
	root := fixture(t)
	write(t, root, "named", "changed")
	write(t, root, "foreign", "staged")
	git(t, root, "add", "foreign")
	write(t, root, "foreign", "staged-plus-unstaged")
	write(t, root, "other", "unstaged")
	write(t, root, "new", "untracked")
	write(t, root, ".gitignore", "ignored\n")
	write(t, root, "ignored", "ignored")
	blob := git(t, root, "show", ":foreign")
	base := git(t, root, "rev-parse", "HEAD")
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	if _, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "named", Paths: []string{"named"}}); err != nil {
		t.Fatal(err)
	}
	if git(t, root, "show", ":foreign") != blob || string(mustRead(t, filepath.Join(root, "foreign"))) != "staged-plus-unstaged" || string(mustRead(t, filepath.Join(root, "other"))) != "unstaged" || string(mustRead(t, filepath.Join(root, "new"))) != "untracked" || string(mustRead(t, filepath.Join(root, "ignored"))) != "ignored" {
		t.Fatal("unnamed state changed")
	}
}

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-q", "-b", "main")
	git(t, root, "config", "user.email", "a@b.c")
	git(t, root, "config", "user.name", "a")
	write(t, root, "named", "base")
	write(t, root, "foreign", "base")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "base")
	return root
}

func raceFixture(t *testing.T) string {
	t.Helper()
	root := fixture(t)
	write(t, root, "other", "base")
	write(t, root, ".gitignore", "ignored\n")
	git(t, root, "add", "other", ".gitignore")
	git(t, root, "commit", "-qm", "race fixture")
	return root
}

type pathSnapshot struct {
	Status, Index []byte
	Worktree      map[string][]byte
}

func dirtyUnnamedState(t *testing.T, root string) ([]string, pathSnapshot) {
	t.Helper()
	write(t, root, "foreign", "staged\n")
	git(t, root, "add", "foreign")
	write(t, root, "foreign", "staged-plus-unstaged\n")
	write(t, root, "other", "unstaged\n")
	write(t, root, "new", "untracked\n")
	write(t, root, "ignored", "ignored\n")
	paths := []string{"foreign", "other", "new", "ignored"}
	return paths, snapshotPaths(t, root, paths...)
}

func snapshotPaths(t *testing.T, root string, paths ...string) pathSnapshot {
	t.Helper()
	statusArgs := append([]string{"status", "--porcelain=v1", "--ignored", "--"}, paths...)
	indexArgs := append([]string{"ls-files", "--stage", "--"}, paths...)
	snapshot := pathSnapshot{
		Status:   gitBytes(t, root, statusArgs...),
		Index:    gitBytes(t, root, indexArgs...),
		Worktree: make(map[string][]byte, len(paths)),
	}
	for _, path := range paths {
		snapshot.Worktree[path] = mustRead(t, filepath.Join(root, path))
	}
	return snapshot
}

func write(t *testing.T, root, path, value string) {
	t.Helper()
	p := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(value), 0o644); err != nil {
		t.Fatal(err)
	}
}
func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	return strings.TrimSpace(string(gitBytes(t, root, args...)))
}

func gitBytes(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	b, e := c.CombinedOutput()
	if e != nil {
		t.Fatalf("git %v: %v: %s", args, e, b)
	}
	return b
}

func gitMode(t *testing.T, root string, args ...string) string {
	t.Helper()
	fields := strings.Fields(git(t, root, args...))
	if len(fields) == 0 {
		t.Fatal("git mode output is empty")
	}
	return fields[0]
}

// The reviewed-landing fixtures below vary one axis: which side edited the spec.
// Every published-tree row shares this builder so a fixture change cannot make two
// rows disagree about what a reviewed landing composes.
const (
	reviewedBaseSpec        = "# X\n\nStatus: staged\n\n## Stories\n\nstory one\nstory two\n"
	reviewedAmendedSpec     = "# X\n\nStatus: staged\n\n## Stories\n\nstory one, amended by the review\nstory two\n"
	reviewedHeadingSpec     = "# X, moved on the destination\n\nStatus: staged\n\n## Stories\n\nstory one\nstory two\n"
	reviewedOverlapSpec     = "# X\n\nStatus: staged\n\n## Stories\n\nstory one, rewritten on the destination\nstory two\n"
	specProvenanceRefusal   = "staged spec bytes are not the reviewed source tip's committed spec"
	reviewedFixtureSpecPath = "specs/x/spec.md"
)

type reviewedLanding struct {
	root, sourceWorktree      string
	base, destination, source string
}

// newReviewedLanding stages baseSpec at the shared review base, commits sourceSpec on
// the reviewed source branch, and advances the destination with destinationSpec. An
// empty spec string means that side wrote no spec at all.
func newReviewedLanding(t *testing.T, baseSpec, destinationSpec, sourceSpec string) reviewedLanding {
	t.Helper()
	root := fixture(t)
	if baseSpec != "" {
		write(t, root, reviewedFixtureSpecPath, baseSpec)
		git(t, root, "add", reviewedFixtureSpecPath)
	}
	git(t, root, "commit", "--allow-empty", "-qm", "review base")
	base := git(t, root, "rev-parse", "HEAD")
	sourceWorktree := filepath.Join(t.TempDir(), "source")
	git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, base)
	write(t, sourceWorktree, "reviewed", "source bytes\n")
	if sourceSpec != "" {
		write(t, sourceWorktree, reviewedFixtureSpecPath, sourceSpec)
	}
	git(t, sourceWorktree, "add", "-A")
	git(t, sourceWorktree, "commit", "-qm", "reviewed source")
	destination := base
	if destinationSpec != "" {
		write(t, root, reviewedFixtureSpecPath, destinationSpec)
		git(t, root, "add", reviewedFixtureSpecPath)
		git(t, root, "commit", "-qm", "destination spec movement")
		destination = git(t, root, "rev-parse", "HEAD")
	}
	return reviewedLanding{
		root: root, sourceWorktree: sourceWorktree, base: base,
		destination: destination, source: git(t, sourceWorktree, "rev-parse", "HEAD"),
	}
}

func (f reviewedLanding) request(t *testing.T, specBytes, message string) ReviewedRequest {
	t.Helper()
	sourceFingerprint, err := CheckoutFingerprint(f.sourceWorktree)
	if err != nil {
		t.Fatal(err)
	}
	destinationFingerprint, err := CheckoutFingerprint(f.root)
	if err != nil {
		t.Fatal(err)
	}
	return ReviewedRequest{
		Root: f.root, Destination: "refs/heads/main", DestinationBase: f.destination,
		Source: "refs/heads/reviewed-source", SourceTip: f.source, ReviewBase: f.base,
		SourceWorktree: f.sourceWorktree, SourceFingerprint: sourceFingerprint,
		DestinationFingerprint: destinationFingerprint,
		SpecPath:               reviewedFixtureSpecPath, SpecBytes: []byte(specBytes), SpecMode: 0o644,
		Message: message,
	}
}

func TestLandReviewedPublishesSourceSpecBytesOverAnyDestinationCopy(t *testing.T) {
	for _, tc := range []struct{ name, baseSpec, destinationSpec, sourceSpec string }{
		{name: "destination-copy-is-stale", baseSpec: reviewedBaseSpec, sourceSpec: reviewedAmendedSpec},
		{name: "destination-moved-after-the-review-base", baseSpec: reviewedBaseSpec, destinationSpec: reviewedHeadingSpec, sourceSpec: reviewedAmendedSpec},
		{name: "destination-never-carried-the-spec", sourceSpec: reviewedBaseSpec},
		{name: "destination-overlaps-the-amendment", baseSpec: reviewedBaseSpec, destinationSpec: reviewedOverlapSpec, sourceSpec: reviewedAmendedSpec},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newReviewedLanding(t, tc.baseSpec, tc.destinationSpec, tc.sourceSpec)
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
			f := newReviewedLanding(t, "", "", tc.sourceSpec)
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
