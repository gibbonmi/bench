// The background row of the board: the latest background outcome of a worktree is a signal
// beside the gate verdict, and never the verdict itself.
package status

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/gittest"
)

// writeBackgroundRecord installs a background record with outcome in the git directory of the
// checkout root, and it answers the log path that the record names.
func writeBackgroundRecord(t *testing.T, root string, outcome backgroundstate.Outcome) string {
	t.Helper()
	gitdir := gittest.Output(t, root, "rev-parse", "--absolute-git-dir")
	log := backgroundstate.LogPath(gitdir)
	record := backgroundstate.Record{Outcome: outcome, Key: "key-" + string(outcome), Tip: strings.Repeat("c", 40), Log: log}
	if err := backgroundstate.WriteRecord(durablefile.Native(), gitdir, record); err != nil {
		t.Fatal(err)
	}
	return log
}

// backgroundSignals answers the background rows of signals.
func backgroundSignals(signals []Signal) []Signal {
	var out []Signal
	for _, s := range signals {
		if s.Name == backgroundSignal {
			out = append(out, s)
		}
	}
	return out
}

// TestStatusShowsTheBackgroundOutcome is BG84. The board of each worktree shows the outcome of
// its own latest background grade and the log path of that grade, for each outcome. A linked
// worktree and the primary checkout each read the record of their own git directory.
func TestStatusShowsTheBackgroundOutcome(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	for _, outcome := range []backgroundstate.Outcome{backgroundstate.Green, backgroundstate.Red, backgroundstate.Aborted, backgroundstate.Stopped} {
		t.Run(string(outcome), func(t *testing.T) {
			primary := initRepo(t)
			gittest.MustRun(t, primary, "commit", "-q", "--allow-empty", "-m", "base")
			linked := filepath.Join(t.TempDir(), "linked")
			gittest.MustRun(t, primary, "worktree", "add", "-q", "-b", "topic", linked)
			other := backgroundstate.Green
			if outcome == backgroundstate.Green {
				other = backgroundstate.Red
			}
			worktrees := []struct {
				root, log string
				outcome   backgroundstate.Outcome
			}{{root: linked, outcome: outcome}, {root: primary, outcome: other}}
			for i := range worktrees {
				worktrees[i].log = writeBackgroundRecord(t, worktrees[i].root, worktrees[i].outcome)
			}
			severity := map[backgroundstate.Outcome]int{}
			for _, tc := range worktrees {
				rows := backgroundSignals(Signals(tc.root))
				if len(rows) != 1 || !strings.HasPrefix(rows[0].Detail, string(tc.outcome)+" ") || !strings.Contains(rows[0].Action, tc.log) {
					t.Fatalf("background rows of %s = %#v, want one %s row that names %s", tc.root, rows, tc.outcome, tc.log)
				}
				if board := render(tc.root, true); !strings.Contains(board, rows[0].Detail) {
					t.Fatalf("board of %s = %q, want the background row %q", tc.root, board, rows[0].Detail)
				}
				severity[tc.outcome] = rows[0].Severity
			}
			// A red asks the author to act, so it ranks above each other outcome.
			if red, ok := severity[backgroundstate.Red]; ok && len(severity) == 2 {
				for outcome, rank := range severity {
					if outcome != backgroundstate.Red && red >= rank {
						t.Fatalf("background ranks = %v, want the red above %s", severity, outcome)
					}
				}
			}
		})
	}
}

// TestStatusShowsNoBackgroundRowWithoutARecord keeps a worktree with no background grade free
// of a background row, so the row spends the five-row budget only on a grade that ran.
func TestStatusShowsNoBackgroundRowWithoutARecord(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := initRepo(t)
	if rows := backgroundSignals(Signals(root)); len(rows) != 0 {
		t.Fatalf("background rows = %#v, want none without a record", rows)
	}
}

// TestStatusShowsAnUnreadableBackgroundRecord keeps a broken record from hiding a red: the board
// shows a background row that names the record, at the rank of a red.
func TestStatusShowsAnUnreadableBackgroundRecord(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := initRepo(t)
	path := backgroundstate.RecordPath(gittest.Output(t, root, "rev-parse", "--absolute-git-dir"))
	if err := os.WriteFile(path, []byte("{\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	redRoot := initRepo(t)
	writeBackgroundRecord(t, redRoot, backgroundstate.Red)
	rows, red := backgroundSignals(Signals(root)), backgroundSignals(Signals(redRoot))
	if len(rows) != 1 || len(red) != 1 || !strings.Contains(rows[0].Action, path) || rows[0].Severity != red[0].Severity {
		t.Fatalf("background rows = %#v, want one row that names %s at the rank of the red %#v", rows, path, red)
	}
}

// TestStatusKeepsTheGateVerdictBesideABackgroundRed is BG99. A background record never stands
// in for the gate cache: a red background outcome beside a green gate cache, and a green
// background outcome beside a red gate cache, each leave the gate verdict and every other row
// of the board as the gate cache alone gives them.
func TestStatusKeepsTheGateVerdictBesideABackgroundRed(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	for _, tc := range []struct {
		gate       string
		background backgroundstate.Outcome
	}{{"green", backgroundstate.Red}, {"red", backgroundstate.Green}} {
		t.Run(tc.gate+" gate", func(t *testing.T) {
			root := initRepo(t)
			tree := treeOf(t, root, map[string]string{"f.txt": "x\n"})
			gittest.MustRun(t, root, "commit", "-q", "--allow-empty", "-m", "base")
			writeFullGateCache(t, root, tree, tc.gate)
			verdict, board := GateVerdict(root), Signals(root)
			if !verdict.Present || verdict.Status != tc.gate {
				t.Fatalf("verdict = %#v, want the recorded %s gate verdict", verdict, tc.gate)
			}
			writeBackgroundRecord(t, root, tc.background)
			if after := GateVerdict(root); !reflect.DeepEqual(after, verdict) {
				t.Fatalf("verdict beside a background %s = %#v, want the gate verdict %#v", tc.background, after, verdict)
			}
			after := Signals(root)
			rest := slices.DeleteFunc(slices.Clone(after), func(s Signal) bool { return s.Name == backgroundSignal })
			if !reflect.DeepEqual(rest, board) || len(after) != len(board)+1 {
				t.Fatalf("board beside a background %s = %#v, want %#v and one background row", tc.background, after, board)
			}
		})
	}
}
