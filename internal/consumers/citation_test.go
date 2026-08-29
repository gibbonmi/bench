package consumers

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// citationHeader is the disclosure block a success response prints. It is pinned as a
// literal, not read from citationFields, because a silently reshaped row is exactly the
// failure CS10 exists to catch.
const citationHeader = "citation[1]{sha,state,version,cmd,hash}:\n"

// citationRepo builds a one-commit repository and makes it the process working
// directory, so the command resolves its root, its HEAD, and its state there instead of
// in the developer's own checkout. The command's loader stays stubbed, so the repository
// needs no Go source.
func citationRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	if _, err := git.Output("-C", root, "commit", "-q", "--allow-empty", "-m", "one"); err != nil {
		t.Fatalf("seed commit: %v", err)
	}
	t.Chdir(root)
	return root
}

// citationCells reads the printed citation row into its five cells, with the output
// bytes that precede the row. No cell the tests assert on carries a comma, so the split
// is the whole parse.
func citationCells(t *testing.T, out string) (cells []string, before string) {
	t.Helper()
	index := strings.Index(out, citationHeader)
	if index < 0 {
		t.Fatalf("stdout = %q, want a citation row", out)
	}
	before = out[:index]
	rest := out[index+len(citationHeader):]
	row, _, _ := strings.Cut(rest, "\n")
	cells = strings.Split(strings.TrimPrefix(row, "  "), ",")
	if len(cells) != 5 {
		t.Fatalf("citation row = %q, want five cells", row)
	}
	// A hex sha or hash that reads as a number with a leading zero is a quoted cell per
	// spec-TOON, so the parse unquotes before it compares a cell to a raw value.
	for i, cell := range cells {
		if len(cell) > 1 && strings.HasPrefix(cell, `"`) && strings.HasSuffix(cell, `"`) {
			cells[i] = cell[1 : len(cell)-1]
		}
	}
	return cells, before
}

// CS10 (story 9): every success response carries one citation row, and it sits
// immediately before the terminal help envelope.
func TestEverySuccessResponseCarriesCitationBeforeHelp(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fixture func(string) []fixturePkg
		arg     string
	}{
		{"symbol result", referencesFixture(3), "target.Symbol"},
		{"candidates result", ambiguousFixture, "Symbol"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			citationRepo(t)
			stubLoad(t, tc.fixture)
			out, code := run(t, tc.arg)
			if code != 0 {
				t.Fatalf("exit = %d, want 0; out=%q", code, out)
			}
			if strings.Count(out, citationHeader) != 1 {
				t.Fatalf("stdout = %q, want exactly one %q", out, citationHeader)
			}
			row := out[strings.Index(out, citationHeader):]
			_, after, _ := strings.Cut(row, "\n")
			_, afterRow, _ := strings.Cut(after, "\n")
			if !strings.HasPrefix(afterRow, "help[") {
				t.Fatalf("stdout = %q, want the citation row immediately before the help envelope", out)
			}
		})
	}
}

// CS11 (story 10): a dirty checkout says so, and a clean one is not marked dirty.
func TestCitationStateReportsCheckoutDirtiness(t *testing.T) {
	root := citationRepo(t)
	stubLoad(t, referencesFixture(1))
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if cells, _ := citationCells(t, out); cells[1] != "clean" {
		t.Fatalf("state = %q on a clean checkout, want clean", cells[1])
	}
	if err := os.WriteFile(filepath.Join(root, "untracked.txt"), []byte("probe\n"), 0o644); err != nil {
		t.Fatalf("dirty the checkout: %v", err)
	}
	out, code = run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if cells, _ := citationCells(t, out); cells[1] != "dirty" {
		t.Fatalf("state = %q on a dirty checkout, want dirty", cells[1])
	}
}

// CS12 (story 11): two runs at one SHA are byte-equal, so a replay of the citation
// reproduces the answer.
func TestTwoRunsAtOneShaAreByteEqual(t *testing.T) {
	citationRepo(t)
	stubLoad(t, splitFixture(3, 2))
	first, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, first)
	}
	second, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, second)
	}
	if first != second {
		t.Fatalf("second run = %q, want the first run's bytes %q", second, first)
	}
}

// CS22 (story 9): the printed hash is the sha256 of every output byte before the
// citation row, so a constant or a partial hash cannot pass.
func TestPrintedHashCoversEveryByteBeforeTheCitation(t *testing.T) {
	citationRepo(t)
	stubLoad(t, referencesFixture(2))
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	cells, before := citationCells(t, out)
	sum := sha256.Sum256([]byte(before))
	if want := hex.EncodeToString(sum[:]); cells[4] != want {
		t.Fatalf("hash = %q, want the recomputed %q", cells[4], want)
	}
}

// CS10 (story 9): the citation names the run that produced it — the checkout's HEAD, the
// bench version, and the replay spelling of the argv.
func TestCitationNamesHeadVersionAndReplaySpelling(t *testing.T) {
	root := citationRepo(t)
	stubLoad(t, referencesFixture(1))
	out, code := run(t, "target.Symbol", "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	head, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("read HEAD: %v", err)
	}
	cells, _ := citationCells(t, out)
	if cells[0] != head {
		t.Fatalf("sha = %q, want HEAD %q", cells[0], head)
	}
	if cells[2] != testVersion {
		t.Fatalf("version = %q, want %q", cells[2], testVersion)
	}
	if cells[3] != "bench consumers target.Symbol --full" {
		t.Fatalf("cmd = %q, want the replay spelling", cells[3])
	}
}

// shellSplit recovers the argv a POSIX shell reads from one printed command line. It is
// the same recovery internal/axi applies to a help row, so a citation's replay spelling is
// graded against a real shell rather than against a re-spelling.
func shellSplit(t *testing.T, command string) []string {
	t.Helper()
	out, err := exec.Command("sh", "-c", "set -- "+command+`; for argument do printf '%s\000' "$argument"; done`).Output()
	if err != nil {
		t.Fatalf("shell split %q: %v", command, err)
	}
	fields := strings.Split(string(out), "\x00")
	return fields[:len(fields)-1]
}

// A revision name carrying a shell metacharacter keeps its token boundary in
// the citation's replay spelling, so the printed line replays as the argv that answered.
func TestCitationCmdRoundTripsThroughAPOSIXShellSplit(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	root, err := git.Root()
	if err != nil {
		t.Fatalf("git root: %v", err)
	}
	const ref = "feature;echo"
	if _, err := git.Output("-C", root, "branch", ref, base); err != nil {
		t.Fatalf("create the fixture ref: %v", err)
	}
	args := []string{"--changed", "--base", ref, "--source-tip", head(t), "--full"}
	out, code := run(t, args...)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	cells, _ := citationCells(t, out)
	want := append([]string{"bench", "consumers"}, args...)
	if got := shellSplit(t, cells[3]); !slices.Equal(got, want) {
		t.Fatalf("shell argv = %q, want %q (cmd cell %q)", got, want, cells[3])
	}
}
