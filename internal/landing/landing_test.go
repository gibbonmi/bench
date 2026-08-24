package landing

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
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

func TestCheckoutFingerprintIgnoresRuntimeRecords(t *testing.T) {
	root := fixture(t)
	write(t, root, ".gitignore", ".logs/\n")
	git(t, root, "add", ".gitignore")
	git(t, root, "commit", "-qm", "ignore runtime logs")
	before, err := CheckoutFingerprint(root)
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, ".logs/gate.jsonl", "record\n")
	after, err := CheckoutFingerprint(root)
	if err != nil || after != before {
		t.Fatalf("runtime record fingerprint = (%q, %v), want %q", after, err, before)
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
	var remainder *PublishedUnreconciledError
	if !errors.As(err, &remainder) || !strings.Contains(err.Error(), "landed-but-checkout-incomplete") || !errors.Is(err, os.ErrPermission) {
		t.Fatalf("err = %v", err)
	}
	if got.Base != base || got.Commit == "" || got.Tree == "" || git(t, root, "rev-parse", "HEAD") != got.Commit {
		t.Fatal("published identity was not retained")
	}
	if remainder.Commit != got.Commit || !reflect.DeepEqual(remainder.Paths, []string{"named"}) {
		t.Fatalf("remainder = %+v, want the published commit and its named paths", remainder)
	}
}

// FB3: the reconcile walks the named paths in order, and the path it fails on is the one
// the publication boundary reports. A report that named the first path, or no path, would
// send the repair at a path that already reconciled.
func TestPublicationBoundaryNamesThePathTheReconcileFailedOn(t *testing.T) {
	root := fixture(t)
	write(t, root, "a-first", "changed")
	write(t, root, "b-second", "changed")
	base := git(t, root, "rev-parse", "HEAD")
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	o.reconcile = func(_ Request, paths []string, _ composedSnapshot) error {
		for _, path := range paths {
			if path == "b-second" {
				return &ReconcileError{Path: path, Err: os.ErrPermission}
			}
		}
		return nil
	}
	_, err := o.Land(context.Background(), Request{
		Root: root, Destination: "refs/heads/main", Expected: base, Message: "x",
		Paths: []string{"b-second", "a-first"},
	})
	var remainder *PublishedUnreconciledError
	if !errors.As(err, &remainder) {
		t.Fatalf("err = %v, want the publication boundary", err)
	}
	if remainder.Path != "b-second" {
		t.Fatalf("Path = %q, want the second named path", remainder.Path)
	}
	if !reflect.DeepEqual(remainder.Paths, []string{"a-first", "b-second"}) {
		t.Fatalf("Paths = %v, want the sorted named paths", remainder.Paths)
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
