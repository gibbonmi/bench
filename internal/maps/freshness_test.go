// Tests for the advisory Sources-freshness rows the maps command projects.
package maps

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

const (
	assetPath   = "decisions/alpha/assets/r.md"
	frontierRow = "  alpha,Which parser owns the map?,Research,frontier,\"\",decisions/alpha/tickets/1.md\n"
	staleRow    = "  alpha," + assetPath + ",source,stale,\"\"," + assetPath + "\n"
	shapeHelp   = "help[1]{cmd,why}:\n  /bench-shape-idea,\"shape alpha: Which parser owns the map?\"\n"
)

// citingIndex is splitIndex with one Sources record naming the map's own asset.
var citingIndex = strings.Replace(splitIndex, "## Sources\n",
	"## Sources\n\n- Path: "+assetPath+"\n  Supports: the parser choice\n  Drift: the asset changes\n", 1)

// freshnessRepo writes the citing map, its ticket, and its asset into a fresh
// repository and returns the root. Nothing is committed yet, so each test commits
// the order it grades.
func freshnessRepo(t *testing.T) string {
	t.Helper()
	root := gittest.Repo(t)
	commitGit(t, root, "config", "user.email", "fixture@example.invalid")
	commitGit(t, root, "config", "user.name", "Fixture")
	writeSplitMap(t, root, "decisions/alpha.md", citingIndex, map[string]string{"1.md": splitTicket})
	writeAsset(t, root, "first\n")
	return root
}

func writeAsset(t *testing.T, root, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(assetPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, out)
	}
}

// commitAt commits the named paths at a fixed date. The dates are fixed because two
// commits made in the same second would leave the compare with nothing to order.
func commitAt(t *testing.T, root, date string, paths ...string) {
	t.Helper()
	commitGit(t, root, append([]string{"add", "--"}, paths...)...)
	cmd := exec.Command("git", "-C", root, "commit", "-q", "-m", date)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+date, "GIT_COMMITTER_DATE="+date)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit %s: %v: %s", date, err, out)
	}
}

func TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex(t *testing.T) {
	for _, testCase := range []struct {
		name  string
		order []string
		want  string
	}{
		{"asset after index", []string{"2001-01-01T00:00:00Z", "2001-01-02T00:00:00Z"}, frontierRow + staleRow},
		{"index after asset", []string{"2001-01-02T00:00:00Z", "2001-01-01T00:00:00Z"}, frontierRow},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := freshnessRepo(t)
			commitAt(t, root, testCase.order[0], "decisions/alpha.md", "decisions/alpha/tickets")
			writeAsset(t, root, "second\n")
			commitAt(t, root, testCase.order[1], assetPath)
			t.Chdir(root)

			out, code := Command(nil)
			want := rendered(testCase.want) + shapeHelp
			if code != 0 || out != want {
				t.Fatalf("Command(nil) = (exit %d, %q), want (exit 0, %q)", code, out, want)
			}
		})
	}
}

// rendered wraps the expected body rows in the typed header the row count implies.
func rendered(body string) string {
	return "maps[" + strconv.Itoa(strings.Count(body, "\n")) + "]{map,title,type,state,blockers,path}:\n" + body
}

func TestCommandProjectsNoStaleRowWhenTheAssetOrTheIndexIsUncommitted(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		commits []string
	}{
		{"asset uncommitted", []string{"decisions/alpha.md", "decisions/alpha/tickets"}},
		{"index uncommitted", []string{assetPath}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := freshnessRepo(t)
			commitAt(t, root, "2001-01-01T00:00:00Z", testCase.commits...)
			t.Chdir(root)

			out, code := Command(nil)
			want := rendered(frontierRow) + shapeHelp
			if code != 0 || out != want {
				t.Fatalf("Command(nil) = (exit %d, %q), want (exit 0, %q)", code, out, want)
			}
			if count, code := Command([]string{"--count"}); code != 0 || count != "1\n" {
				t.Fatalf("Command(--count) = (exit %d, %q), want (exit 0, \"1\\n\")", code, count)
			}
		})
	}
}

func TestCommandProjectsNoStaleRowWhenTheAssetAndTheIndexShareOneCommit(t *testing.T) {
	root := freshnessRepo(t)
	commitAt(t, root, "2001-01-01T00:00:00Z", "decisions/alpha.md", "decisions/alpha")
	t.Chdir(root)

	out, code := Command(nil)
	want := rendered(frontierRow) + shapeHelp
	if code != 0 || out != want {
		t.Fatalf("Command(nil) = (exit %d, %q), want (exit 0, %q)", code, out, want)
	}
}

func TestStaleRowLeavesTheExitCodeAndTheCountUnchanged(t *testing.T) {
	root := freshnessRepo(t)
	commitAt(t, root, "2001-01-01T00:00:00Z", "decisions/alpha.md", "decisions/alpha/tickets")
	writeAsset(t, root, "second\n")
	commitAt(t, root, "2001-01-02T00:00:00Z", assetPath)
	t.Chdir(root)

	out, code := Command(nil)
	if code != 0 || !strings.Contains(out, staleRow) {
		t.Fatalf("Command(nil) = (exit %d, %q), want exit 0 with the stale row", code, out)
	}
	if count, code := Command([]string{"--count"}); code != 0 || count != "1\n" {
		t.Fatalf("Command(--count) = (exit %d, %q), want (exit 0, \"1\\n\")", code, count)
	}
}
