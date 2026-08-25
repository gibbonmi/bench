// Identity component tests: one producing fixture per registry entry, the same six
// sentences under a resumed landing, and the source fence that keeps a merged sentence
// from returning.
package worktree

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// identityComponentFixture produces exactly one component's refusal. mutate breaks one
// identity dimension of the landing fixture, and want renders the whole refused record
// the verb must print for it. The registry walk requires one fixture per entry, so a
// new component with no fixture turns this file red.
type identityComponentFixture struct {
	component string
	request   func(request string) string
	mutate    func(t *testing.T, root string, creation Creation)
	want      func(creation Creation, base, tip string) string
}

func identityComponentFixtures() []identityComponentFixture {
	unchanged := func(request string) string { return request }
	return []identityComponentFixture{
		{
			component: componentRequest,
			request:   func(string) string { return "unknown-request" },
			mutate:    func(*testing.T, string, Creation) {},
			want: func(creation Creation, base, tip string) string {
				next := "bench worktree reauthorize --assignment " + creation.Assignment.ID +
					" --request <new-request> --base '" + base + "' --source-tip '" + tip + "' '" + creation.Path + "'"
				return "detail=request token matches no assignment,observed=assignment:" + creation.Assignment.ID + ",next=" + next
			},
		},
		{
			component: componentAssignmentState,
			request:   unchanged,
			mutate: func(t *testing.T, root string, creation Creation) {
				a := creation.Assignment
				a.State = intent.StateComplete
				if err := intent.PutAssignment(root, a); err != nil {
					t.Fatal(err)
				}
			},
			want: func(creation Creation, _, _ string) string {
				return "detail=assignment " + creation.Assignment.ID + " is not active,observed=complete,wanted=active"
			},
		},
		{
			component: componentAssignmentPath,
			request:   unchanged,
			mutate: func(t *testing.T, root string, creation Creation) {
				a := creation.Assignment
				a.Worktree = otherWorktreePath(creation.Path)
				if err := intent.PutAssignment(root, a); err != nil {
					t.Fatal(err)
				}
			},
			want: func(creation Creation, _, _ string) string {
				return "detail=assignment " + creation.Assignment.ID + " owns another worktree,observed=" +
					otherWorktreePath(creation.Path) + ",wanted=" + creation.Path
			},
		},
		{
			component: componentOwnerMarker,
			request:   unchanged,
			mutate: func(t *testing.T, _ string, creation Creation) {
				rewriteMarkerOwner(t, creation.Path, strings.Repeat("a", 32))
			},
			want: func(creation Creation, _, _ string) string {
				return "detail=owner marker does not match assignment " + creation.Assignment.ID
			},
		},
		{
			component: componentRegistration,
			request:   unchanged,
			mutate: func(t *testing.T, _ string, creation Creation) {
				gitRun(t, creation.Path, "symbolic-ref", "HEAD", "refs/heads/re-pointed")
			},
			want: func(creation Creation, _, _ string) string {
				return "detail=worktree registration does not match assignment " + creation.Assignment.ID
			},
		},
		{
			component: componentLock,
			request:   unchanged,
			mutate: func(t *testing.T, root string, creation Creation) {
				gitRun(t, root, "worktree", "unlock", creation.Path)
			},
			want: func(creation Creation, _, _ string) string {
				return "detail=Bench lock does not match assignment " + creation.Assignment.ID
			},
		},
	}
}

// otherWorktreePath names a tree the assignment could own but the caller did not target.
// It is derived from the target so the two paths differ in the record the operator reads.
func otherWorktreePath(target string) string { return target + "-elsewhere" }

func rewriteMarkerOwner(t *testing.T, path, owner string) {
	t.Helper()
	markerFile, err := markerPath(path)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatal(err)
	}
	marker, err := decodeMarker(data)
	if err != nil {
		t.Fatal(err)
	}
	marker.OwnerID = owner
	rewritten, err := json.Marshal(marker)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, markerFile, append(rewritten, '\n'), 0o600)
}

// TestLandCommandNamesEachIdentityComponent is LR01 and LR04 through LR08: one broken
// dimension per case, and the whole refused record the operator reads for it.
func TestLandCommandNamesEachIdentityComponent(t *testing.T) {
	t.Parallel()
	for _, fixture := range identityComponentFixtures() {
		t.Run(fixture.component, func(t *testing.T) {
			request := "land-component-" + fixture.component
			root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
			fixture.mutate(t, root, creation)
			var stdout, stderr bytes.Buffer
			code := LandCommand(root, home, "", landArgs(fixture.request(request), base, tip, creation.Path), &stdout, &stderr)
			want := "refused{" + fixture.want(creation, base, tip) + "}\n"
			if code != 1 || stdout.String() != want {
				t.Fatalf("%s landing = (%d, %q, %q), want exit 1 and %q", fixture.component, code, stdout.String(), stderr.String(), want)
			}
		})
	}
}

// TestResumeLandCommandNamesEachIdentityComponent is LR09: a resumed landing reads like
// a first landing, so the same six mutations print the same six sentences.
func TestResumeLandCommandNamesEachIdentityComponent(t *testing.T) {
	for _, fixture := range identityComponentFixtures() {
		t.Run(fixture.component, func(t *testing.T) {
			request := "resume-component-" + fixture.component
			root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
			published := interruptLandingAtMarker(t, root, creation, request, base, tip)
			fixture.mutate(t, root, creation)
			var stdout, stderr bytes.Buffer
			args := []string{"--resume", published, "--request", fixture.request(request), "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
			code := LandCommand(root, home, "", args, &stdout, &stderr)
			want := "refused{" + fixture.want(creation, base, tip) + "}\n"
			if code != 1 || stdout.String() != want {
				t.Fatalf("%s resume = (%d, %q, %q), want exit 1 and %q", fixture.component, code, stdout.String(), stderr.String(), want)
			}
		})
	}
}

// interruptLandingAtMarker publishes the landing and then fails its marker step, which
// is the state a resume exists to finish. It returns the published commit.
func interruptLandingAtMarker(t *testing.T, root string, creation Creation, request, base, tip string) string {
	t.Helper()
	old := advanceLandingMarker
	advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, Home(), "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	advanceLandingMarker = old
	if code != 3 {
		t.Fatalf("interrupted landing = (%d, %q, %q), want exit 3", code, stdout.String(), stderr.String())
	}
	return gitOutput(t, root, "rev-parse", "main")
}

// TestIdentityComponentRegistryHasAProducingFixture is LR16. The registry is the source
// of the component set, so an entry added without a fixture is a red here rather than an
// unproven sentence in production.
func TestIdentityComponentRegistryHasAProducingFixture(t *testing.T) {
	t.Parallel()
	produced := map[string]bool{}
	for _, fixture := range identityComponentFixtures() {
		if produced[fixture.component] {
			t.Fatalf("component %q has two producing fixtures", fixture.component)
		}
		produced[fixture.component] = true
	}
	for _, component := range identityComponents {
		if !produced[component.name] {
			t.Errorf("registry component %q has no producing fixture", component.name)
		}
		delete(produced, component.name)
	}
	for name := range produced {
		t.Errorf("fixture %q produces no registered component", name)
	}
	if identityComponentByName(componentRequest).recovers == identityComponentByName(componentLock).recovers {
		t.Errorf("the request component is the only one that carries a recovery command")
	}
}

// mergedDetailWords are the words the registry's own component names are built from. A
// production sentence that puts two of them on opposite sides of " or " is the merged
// string this spec retires. `branch` is deliberately absent: it names no component, and
// the reauthorize verb's own proofs stay out of this registry.
var mergedDetailWords = []string{"request", "assignment", "state", "path", "owner", "marker", "registration", "lock"}

// deferredMergedDetails are the merged sentences a later ticket of this spec retires.
// The test requires each one to still be present, so the allowance cannot outlive the
// sentence it covers.
var deferredMergedDetails []string

// TestPackageDetailSentencesNameOneComponent is LR17. It reads the package's own source
// rather than its output, so a merged sentence at a site with no fixture is still caught.
func TestPackageDetailSentencesNameOneComponent(t *testing.T) {
	t.Parallel()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	deferred := map[string]bool{}
	for _, name := range entries {
		if !strings.HasSuffix(name.Name(), ".go") || strings.HasSuffix(name.Name(), "_test.go") {
			continue
		}
		for _, literal := range stringLiterals(t, name.Name()) {
			if !mergesTwoComponentWords(literal) {
				continue
			}
			if slices.Contains(deferredMergedDetails, literal) {
				deferred[literal] = true
				continue
			}
			t.Errorf("%s: detail %q names two components; name the one that failed", name.Name(), literal)
		}
	}
	for _, sentence := range deferredMergedDetails {
		if !deferred[sentence] {
			t.Errorf("deferred merged detail %q is gone; drop it from deferredMergedDetails", sentence)
		}
	}
}

func stringLiterals(t *testing.T, name string) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var literals []string
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}
		if value, err := strconv.Unquote(literal.Value); err == nil {
			literals = append(literals, value)
		}
		return true
	})
	return literals
}

func mergesTwoComponentWords(literal string) bool {
	before, after, found := strings.Cut(literal, " or ")
	return found && namesAComponentWord(before) && namesAComponentWord(after)
}

func namesAComponentWord(text string) bool {
	for _, field := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return r < 'a' || r > 'z' }) {
		if slices.Contains(mergedDetailWords, field) {
			return true
		}
	}
	return false
}

// identityComponentFixtureFor returns the one fixture that produces the named component.
// A double-fault case composes two of these mutations, so the mutation logic stays in the
// fixture table and no case copies it.
func identityComponentFixtureFor(t *testing.T, component string) identityComponentFixture {
	t.Helper()
	for _, fixture := range identityComponentFixtures() {
		if fixture.component == component {
			return fixture
		}
	}
	t.Fatalf("no producing fixture for component %q", component)
	return identityComponentFixture{}
}

// TestLandCommandNamesTheEarlierComponentOfTwo covers the edge inventory row for two
// components that fail inside one bundle. The registration precedes the lock in the
// registry, so the registration sentence is the one the operator reads.
func TestLandCommandNamesTheEarlierComponentOfTwo(t *testing.T) {
	t.Parallel()
	request := "land-component-registration-and-lock"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	registration := identityComponentFixtureFor(t, componentRegistration)
	registration.mutate(t, root, creation)
	identityComponentFixtureFor(t, componentLock).mutate(t, root, creation)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	want := "refused{" + registration.want(creation, base, tip) + "}\n"
	if code != 1 || stdout.String() != want {
		t.Fatalf("double-fault landing = (%d, %q, %q), want exit 1 and %q", code, stdout.String(), stderr.String(), want)
	}
}

// TestTargetVerbNamesTheOwnerMarkerBeforeTheBranch pins the bundle validator's precedence.
// A wrong owner id and a detached HEAD fail together, and the marker is the earlier
// component, so the branch sentence must not win.
func TestTargetVerbNamesTheOwnerMarkerBeforeTheBranch(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "marker-before-branch")
	rewriteMarkerOwner(t, creation.Path, strings.Repeat("a", 32))
	gitRun(t, creation.Path, "checkout", "--detach")
	chdir(t, root)
	var stdout, stderr bytes.Buffer
	code := PathCommand(root, home, []string{creation.Assignment.Label}, &stdout, &stderr)
	want := "bench worktree path: owner marker does not match assignment " + creation.Assignment.ID + "\n"
	if code != 1 || stderr.String() != want {
		t.Fatalf("double-fault path = (%d, %q, %q), want exit 1 and stderr %q", code, stdout.String(), stderr.String(), want)
	}
}
