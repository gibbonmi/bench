package preflight

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// foldedPath is the unfenced path main adds after the source base. Only the fold of main
// brings it into the source range.
const foldedPath = "outside/main.go"

// seedFoldedMain moves main past the conformant source base and folds main into the source,
// as a build does before its completion landing. It returns the pre-fold source base.
func seedFoldedMain(t *testing.T) (root, slug, base string) {
	t.Helper()
	root, slug = preflighttest.SeedConformant(t)
	base = preflighttest.RunGit(t, "rev-parse", "main")
	preflighttest.RunGit(t, "checkout", "-q", "main")
	preflighttest.MustWriteFile(t, foldedPath, "package outside\n")
	preflighttest.RunGit(t, "add", foldedPath)
	preflighttest.RunGit(t, "commit", "-q", "-m", "main moves")
	preflighttest.RunGit(t, "checkout", "-q", "feature")
	preflighttest.RunGit(t, "merge", "-q", "--no-edit", "main")
	return root, slug, base
}

// pathsAuthorizedRow is the paths-authorized row the preflight renders for mode over the
// explicit range from base.
func pathsAuthorizedRow(t *testing.T, root, mode, slug, base string) (CheckResult, Facts) {
	t.Helper()
	facts, failure := Gather(root, mode, slug, base)
	if failure != nil {
		t.Fatalf("Gather(%s) = %s: %s", mode, failure.Kind, failure.Hint)
	}
	for _, row := range Decide(facts).Checks {
		if row.Check == "paths-authorized" {
			return row, facts
		}
	}
	t.Fatalf("Decide(%s) renders no paths-authorized row", mode)
	return CheckResult{}, facts
}

func TestPathsAuthorizedAdmitsWhatAFoldedMainBringsInUnchanged(t *testing.T) {
	root, slug, base := seedFoldedMain(t)
	for _, mode := range []string{modeReview, modeBuild} {
		row, facts := pathsAuthorizedRow(t, root, mode, slug, base)
		if !slices.Contains(facts.ChangedPaths, foldedPath) {
			t.Fatalf("%s range %v lacks the folded path %s", mode, facts.ChangedPaths, foldedPath)
		}
		if row.Verdict != verdictGreen {
			t.Fatalf("%s paths-authorized over a folded main = %+v, want green", mode, row)
		}
	}
	source, err := AuthorizeReviewedSource(root, slug, base)
	if err != nil {
		t.Fatalf("landing authorization over a folded main: %v", err)
	}
	if !slices.Contains(source.CommittedPaths, foldedPath) {
		t.Fatalf("authorized range %v lacks the folded path %s", source.CommittedPaths, foldedPath)
	}
}

func TestPathsAuthorizedKeepsTheFenceRuleForBuildWritesBesideAFold(t *testing.T) {
	for _, test := range []struct {
		name, unauthorized string
		write              func(t *testing.T)
	}{
		{
			name:         "an unfenced path the build adds",
			unauthorized: "unfenced/build.go",
			write: func(t *testing.T) {
				preflighttest.MustWriteFile(t, "unfenced/build.go", "package unfenced\n")
				preflighttest.RunGit(t, "add", "unfenced/build.go")
			},
		},
		{
			name:         "a main path the build changes again",
			unauthorized: foldedPath,
			write: func(t *testing.T) {
				preflighttest.MustWriteFile(t, foldedPath, "package outside\n// build\n")
				preflighttest.RunGit(t, "add", foldedPath)
			},
		},
		{
			name:         "a main path the build changes only in mode",
			unauthorized: foldedPath,
			write: func(t *testing.T) {
				if err := os.Chmod(foldedPath, 0o755); err != nil {
					t.Fatal(err)
				}
				preflighttest.RunGit(t, "update-index", "--chmod=+x", foldedPath)
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug, base := seedFoldedMain(t)
			test.write(t)
			preflighttest.RunGit(t, "commit", "-q", "-m", "build write")
			want := "not authorized by any ownership fence: " + test.unauthorized
			for _, mode := range []string{modeReview, modeBuild} {
				if row, _ := pathsAuthorizedRow(t, root, mode, slug, base); row.Verdict != verdictRed || row.Detail != want {
					t.Fatalf("%s paths-authorized = %+v, want red with detail %q", mode, row, want)
				}
			}
			_, err := AuthorizeReviewedSource(root, slug, base)
			var unfenced UnauthorizedPathsError
			if !errors.As(err, &unfenced) || !slices.Equal(unfenced.Paths, []string{test.unauthorized}) {
				t.Fatalf("landing authorization = %v, want only %s unfenced", err, test.unauthorized)
			}
		})
	}
}

// An untracked path is absent at both tips, so the tree comparison alone would read it as
// unchanged. A dirty source therefore keeps the fence rule for every path.
func TestPathsAuthorizedKeepsTheFenceRuleInADirtyFoldedSource(t *testing.T) {
	root, slug, base := seedFoldedMain(t)
	const untracked = "unfenced/untracked.go"
	preflighttest.MustWriteFile(t, untracked, "package unfenced\n")
	row, facts := pathsAuthorizedRow(t, root, modeBuild, slug, base)
	if !slices.Contains(facts.ChangedPaths, untracked) {
		t.Fatalf("dirty build range %v lacks the untracked path %s", facts.ChangedPaths, untracked)
	}
	if row.Verdict != verdictRed || !strings.Contains(row.Detail, untracked) {
		t.Fatalf("dirty build paths-authorized = %+v, want red naming %s", row, untracked)
	}
}
