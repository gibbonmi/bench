package releasepreflight

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// The release identity checks grade one fixture repository. Its tagged commit, its
// package version, its go.mod toolchain line, and the binary version agree, so each
// refusal case changes one fact and only that fact can cause the refusal.

const (
	fixtureVersion = "1.2.3"
	fixtureTag     = "v" + fixtureVersion
	fixtureRef     = "refs/tags/" + fixtureTag
)

// identityFixture returns a runner whose root is a committed repository. The commit
// carries the tag, origin/main points at it, and the runner's source commit is HEAD.
// It sets the preflight ref to the exact tag.
func identityFixture(t *testing.T) *runner {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	goVersion, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		t.Fatal(err)
	}
	toolchain := strings.TrimPrefix(strings.TrimSpace(string(goVersion)), "go")
	writeFixture(t, root, "package.json", `{"version":"`+fixtureVersion+`"}`)
	writeFixture(t, root, "go.mod", "module example.invalid/fixture\n\ngo 1.25\n\ntoolchain go"+toolchain+"\n")
	gittest.Output(t, root, "add", "-A")
	gittest.Output(t, root, "commit", "-qm", "release")
	gittest.Output(t, root, "tag", fixtureTag)
	gittest.Output(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	head := gittest.Output(t, root, "rev-parse", "HEAD")
	t.Setenv("BENCH_PREFLIGHT_REF", fixtureRef)
	return &runner{root: root, binaryVersion: fixtureVersion, identity: Identity{SourceCommit: &head}}
}

func writeFixture(t *testing.T, root, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// assertIdentityRefusal fails the test unless err is an identity refusal whose message
// is want.
func assertIdentityRefusal(t *testing.T, err error, want string) {
	t.Helper()
	var failure commandFailure
	if !errors.As(err, &failure) || failure.failure.Kind != "identity" || failure.failure.Message != want {
		t.Fatalf("err = %v, want identity refusal %q", err, want)
	}
}

func TestCheckIdentityAcceptsAnAgreeingRelease(t *testing.T) {
	r := identityFixture(t)
	if err := r.checkIdentity(context.Background()); err != nil {
		t.Fatalf("checkIdentity = %v, want nil", err)
	}
	if r.identity.Tag == nil || *r.identity.Tag != fixtureTag || r.identity.PackageVersion == nil || *r.identity.PackageVersion != fixtureVersion {
		t.Fatalf("identity = %#v, want tag %s and package %s", r.identity, fixtureTag, fixtureVersion)
	}
}

// Each case names the production mutation in checkIdentity that turns it red.
func TestCheckIdentityRefusesADisagreement(t *testing.T) {
	for _, tc := range []struct {
		name     string
		mutation string
		change   func(*testing.T, *runner)
		want     string
	}{
		{
			name:     "ref is not an exact tag",
			mutation: "delete the exactTag refusal",
			change:   func(t *testing.T, _ *runner) { t.Setenv("BENCH_PREFLIGHT_REF", "refs/heads/"+fixtureTag) },
			want:     "publish requires exact GITHUB_REF refs/tags/vMAJOR.MINOR.PATCH",
		},
		{
			name:     "tag does not resolve to HEAD",
			mutation: "drop the tagCommit != commit comparison",
			change: func(t *testing.T, r *runner) {
				gittest.Output(t, r.root, "commit", "-qm", "after the tag", "--allow-empty")
				head := gittest.Output(t, r.root, "rev-parse", "HEAD")
				r.identity.SourceCommit = &head
			},
			want: "tag does not resolve exactly to HEAD",
		},
		{
			name:     "package version disagrees",
			mutation: "drop the pkg != version comparison",
			change:   func(t *testing.T, r *runner) { writeFixture(t, r.root, "package.json", `{"version":"1.2.4"}`) },
			want:     "tag, package version, and binary version must agree",
		},
		{
			name:     "binary version disagrees",
			mutation: "drop the r.binaryVersion != version comparison",
			change:   func(_ *testing.T, r *runner) { r.binaryVersion = "1.2.4" },
			want:     "tag, package version, and binary version must agree",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := identityFixture(t)
			tc.change(t, r)
			err := r.checkIdentity(context.Background())
			if err == nil {
				t.Fatalf("checkIdentity accepted; red under mutation: %s", tc.mutation)
			}
			assertIdentityRefusal(t, err, tc.want)
		})
	}
}

func TestCheckAncestryAcceptsAHeadThatOriginMainContains(t *testing.T) {
	r := identityFixture(t)
	if err := r.checkAncestry(context.Background()); err != nil {
		t.Fatalf("checkAncestry = %v, want nil", err)
	}
}

// Mutation that turns this red: swap the merge-base operands to origin/main HEAD.
func TestCheckAncestryRefusesAHeadThatOriginMainDoesNotContain(t *testing.T) {
	r := identityFixture(t)
	gittest.Output(t, r.root, "commit", "-qm", "unpublished", "--allow-empty")
	assertIdentityRefusal(t, r.checkAncestry(context.Background()), "tagged HEAD ancestry to origin/main could not be proven")
}

// changelogFixture returns a runner whose identity holds the fixture tag and whose root
// holds a changelog with the given release heading.
func changelogFixture(t *testing.T, heading string) *runner {
	t.Helper()
	root := t.TempDir()
	writeFixture(t, root, "CHANGELOG.md", "# Changelog\n\n## [Unreleased]\n\n"+heading+"\n\n- Released.\n")
	tag := fixtureTag
	return &runner{root: root, identity: Identity{Tag: &tag}}
}

func TestCheckChangelogAcceptsOneDatedReleaseHeading(t *testing.T) {
	r := changelogFixture(t, "## ["+fixtureVersion+"] - 2026-09-30")
	if err := r.checkChangelog(); err != nil {
		t.Fatalf("checkChangelog = %v, want nil", err)
	}
	if r.identity.ChangelogHeading == nil || *r.identity.ChangelogHeading != "## ["+fixtureVersion+"] - 2026-09-30" {
		t.Fatalf("changelog heading = %v", r.identity.ChangelogHeading)
	}
}

// Mutation that turns this red: delete the len(matches) != 1 refusal.
func TestCheckChangelogRefusesAMissingReleaseHeading(t *testing.T) {
	r := changelogFixture(t, "## [1.2.2] - 2026-09-29")
	assertIdentityRefusal(t, r.checkChangelog(), "CHANGELOG.md must contain exactly one matching release heading")
}

// Mutation that turns this red: weaken len(matches) != 1 to len(matches) == 0.
func TestCheckChangelogRefusesADuplicateReleaseHeading(t *testing.T) {
	heading := "## [" + fixtureVersion + "] - 2026-09-30"
	r := changelogFixture(t, heading+"\n\n"+heading)
	assertIdentityRefusal(t, r.checkChangelog(), "CHANGELOG.md must contain exactly one matching release heading")
}
