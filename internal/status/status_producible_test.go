// Table test proving every producible board action is invocable or deliberately empty.
package status

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/worktree"
)

func TestAllProducibleBoardActionsAreInvocableOrEmpty(t *testing.T) {
	write := func(t *testing.T, root, path, body string, mode os.FileMode) {
		t.Helper()
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	commit := func(t *testing.T, root string) {
		t.Helper()
		gitRun(t, root, "add", "-A")
		gitRun(t, root, "commit", "-m", "fixture")
	}
	cleanRepo := func(t *testing.T) string {
		t.Helper()
		root := initRepo(t)
		write(t, root, "tracked.txt", "base\n", 0o644)
		commit(t, root)
		gitRun(t, root, "branch", "-M", "main")
		return root
	}
	gateCache := func(t *testing.T, root string) string {
		t.Helper()
		return filepath.Join(gitRun(t, root, "rev-parse", "--absolute-git-dir"), git.GateCacheFile)
	}
	writePendingGate := func(t *testing.T, root string) {
		t.Helper()
		record := fmt.Sprintf(`{"schema":1,"state":"pending","tree":%q,"oracle":%q,"started_at":%q,"owner_pid":999999}`+"\n",
			gitRun(t, root, "write-tree"), strings.Repeat("0", 64), time.Now().UTC().Add(-time.Minute).Truncate(time.Second).Format(time.RFC3339))
		if err := os.WriteFile(gateCache(t, root), []byte(record), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	readyGate := func(t *testing.T, status string) string {
		t.Helper()
		root := cleanRepo(t)
		write(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n", 0o644)
		write(t, root, ".bench/gate.sh", "#!/bin/sh\nexit 0\n", 0o755)
		commit(t, root)
		if result := gate.Execute(context.Background(), root, io.Discard, io.Discard); result.ActionExit != 0 {
			t.Fatalf("seed gate exit = %d", result.ActionExit)
		}
		path := gateCache(t, root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var record map[string]any
		if err := json.Unmarshal(data, &record); err != nil {
			t.Fatal(err)
		}
		record["status"] = status
		data, err = json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
		return root
	}
	readyMap := func() string {
		body := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
		return strings.Replace(body, "Status: shaping", "Status: ready", 1)
	}

	type fixture struct {
		name, signal, detail string
		count                int
		setup                func(*testing.T) (string, Query)
		exact                []Signal
		board                string
		route                *RouteResult
	}
	cases := []fixture{
		{name: "setup", signal: "setup", detail: "no .bench/", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			if err := os.RemoveAll(filepath.Join(root, ".bench")); err != nil {
				t.Fatal(err)
			}
			return root, Query{}
		}},
		{name: "gate red", signal: "gate", detail: "red", setup: func(t *testing.T) (string, Query) {
			return readyGate(t, "red"), Query{}
		}},
		{name: "gate timeout", signal: "gate", detail: "timeout", setup: func(t *testing.T) (string, Query) {
			return readyGate(t, "timeout"), Query{}
		}},
		{name: "gate invalid", signal: "gate", detail: "invalid verdict", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			if err := os.WriteFile(gateCache(t, root), []byte("{}\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			return root, Query{}
		}},
		{name: "gate unavailable", signal: "gate", detail: "verdict unavailable", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			gittest.StubGit(t, root, "fail-rev-parse", filepath.Join(t.TempDir(), "argv"))
			return root, Query{}
		}},
		{name: "gate drifted", signal: "gate", detail: "stale (gated tree", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			writeFullGateCache(t, root, gitRun(t, root, "write-tree"), "green")
			write(t, root, "tracked.txt", "moved\n", 0o644)
			gitRun(t, root, "add", "tracked.txt")
			return root, Query{}
		}},
		{name: "gate partial", signal: "gate", detail: "partial green", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			writePartialGateCache(t, root, gitRun(t, root, "write-tree"), "docs")
			return root, Query{}
		}},
		{name: "gate interrupted", signal: "gate", detail: "interrupted-pending", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			writePendingGate(t, root)
			return root, Query{}
		}},
		{name: "gate locked", signal: "gate", detail: "locked-pending", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			started, release := filepath.Join(t.TempDir(), "started"), filepath.Join(t.TempDir(), "release")
			write(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n", 0o644)
			write(t, root, ".bench/gate.sh", fmt.Sprintf("#!/bin/sh\nset -eu\n: > %q\nwhile [ ! -e %q ]; do sleep 0.01; done\n", started, release), 0o755)
			commit(t, root)
			done := make(chan struct{})
			go func() {
				gate.Execute(context.Background(), root, io.Discard, io.Discard)
				close(done)
			}()
			deadline := time.Now().Add(5 * time.Second)
			for {
				if _, err := os.Stat(started); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("gate fixture did not acquire its lock")
				}
				time.Sleep(10 * time.Millisecond)
			}
			t.Cleanup(func() {
				_ = os.WriteFile(release, nil, 0o600)
				select {
				case <-done:
				case <-time.After(5 * time.Second):
					t.Error("gate fixture did not release its lock")
				}
			})
			return root, Query{}
		}},
		{name: "git unavailable", signal: "git", detail: "git state unavailable", setup: func(t *testing.T) (string, Query) {
			return t.TempDir(), Query{}
		}, exact: []Signal{testSignal(1, "git", "git state unavailable", "git status")}},
		{name: "git dirty", signal: "git", detail: "dirty path", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "tracked.txt", "dirty\n", 0o644)
			return root, Query{}
		}},
		{name: "git unpushed", signal: "git", detail: "unpushed commit", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			branch := gitRun(t, root, "rev-parse", "--abbrev-ref", "HEAD")
			gitRun(t, root, "remote", "add", "origin", root)
			gitRun(t, root, "update-ref", "refs/remotes/origin/"+branch, "HEAD")
			gitRun(t, root, "config", "branch."+branch+".remote", "origin")
			gitRun(t, root, "config", "branch."+branch+".merge", "refs/heads/"+branch)
			write(t, root, "ahead.txt", "ahead\n", 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(1, "git", "1 unpushed commit", "git push")}},
		{name: "git unique branch", signal: "git", detail: "unique branch", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			gitRun(t, root, "checkout", "-b", "feature")
			write(t, root, "feature.txt", "feature\n", 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(1, "git", "1 unique branch", "git push")}},
		{name: "git unclaimed assignment branch", signal: "git", detail: "unclaimed assignment branch", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			gitRun(t, root, "checkout", "-b", "bench/assign/orphan/ref")
			write(t, root, "orphan.txt", "orphan\n", 0o644)
			commit(t, root)
			gitRun(t, root, "checkout", "main")
			return root, Query{}
		}, exact: []Signal{testSignal(1, "git", "1 unclaimed assignment branch", "bench worktree clean --discard-branch --unclaimed --apply-current")}},
		{name: "git mixed unclaimed assignment and feature branches", signal: "git", detail: "1 unclaimed assignment branch, 1 unique branch", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			gitRun(t, root, "checkout", "-b", "bench/assign/orphan/mixed")
			write(t, root, "orphan-mixed.txt", "orphan\n", 0o644)
			commit(t, root)
			gitRun(t, root, "checkout", "main")
			gitRun(t, root, "checkout", "-b", "feature-mixed")
			write(t, root, "feature-mixed.txt", "feature\n", 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(1, "git", "1 unclaimed assignment branch, 1 unique branch", "bench worktree clean --discard-branch --unclaimed --apply-current")}},
		{name: "worktree leased and out of pool", signal: "worktree", count: 2, setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			t.Setenv("BENCH_HOME", t.TempDir())
			leased := filepath.Join(worktree.Pool(root), "leased")
			if err := os.MkdirAll(filepath.Dir(leased), 0o755); err != nil {
				t.Fatal(err)
			}
			gitRun(t, root, "worktree", "add", "-q", "--detach", leased, "HEAD")
			lease, err := worktree.LeaseFile(leased)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(lease, []byte("leased\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			outside := filepath.Join(t.TempDir(), "outside")
			gitRun(t, root, "worktree", "add", "-q", "--detach", outside, "HEAD")
			return root, Query{}
		}, exact: []Signal{
			testSignal(2, "worktree", "1 out-of-pool worktree", "bench worktree clean <path>"),
			testSignal(2, "worktree", "1 leased pool worktree", "bench worktree list"),
		}},
		{name: "worktree typed failure", signal: "worktree", detail: "rev-parse", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			gittest.StubGit(t, root, "fail-rev-parse", filepath.Join(t.TempDir(), "argv"))
			return root, Query{}
		}},
		{name: "worktree porcelain failure", signal: "worktree", detail: "git worktree list failed", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			gittest.StubGit(t, root, "fail-worktree", filepath.Join(t.TempDir(), "argv"))
			return root, Query{}
		}},
		{name: "intent live", signal: "intent", detail: "correlated", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			if err := intent.Upsert(root, intent.Entry{Key: "live", Kind: intent.KindWorktree, CreatedAt: time.Unix(1, 0), Worktree: root}); err != nil {
				t.Fatal(err)
			}
			return root, Query{}
		}, exact: []Signal{testSignal(2, "intent", "1 correlated, 0 uncorrelated; oldest: live", "bench status --all")}},
		{name: "intent unavailable", signal: "intent", detail: "intent ledger unavailable", setup: func(t *testing.T) (string, Query) {
			return t.TempDir(), Query{}
		}},
		{name: "guards", signal: "guards", detail: "pre-push missing", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, ".bench/lines.env", "BENCH_CODEX_MID=test\n", 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "one staged spec", signal: "specs", detail: "1 staged spec", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "specs/my spec/spec.md", "Status: staged\n", 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "multiple staged specs", signal: "specs", detail: "2 staged spec", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "specs/one/spec.md", "Status: staged\n", 0o644)
			write(t, root, "specs/two/spec.md", "Status: staged\n", 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "drain", signal: "drain", detail: "1 idea", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "capture/IDEAS.md", "- 2026-08-18  pending\n", 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "structure", signal: "structure", detail: "1 issue", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "long.go", strings.Repeat("x\n", 401), 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(5, "structure", "1 issue(s)", "bench structure")}},
		{name: "unresolved map", signal: "decisions", detail: "unresolved map", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "decisions/shaping.md", maps.DecisionMapTemplate(), 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(6, "decisions", "1 unresolved map(s)", "/bench-shape-idea")},
			board: "▶ /bench-shape-idea  (decisions)\n  decisions  1 unresolved map(s)            → /bench-shape-idea\n",
			route: &RouteResult{Lead: testSignal(6, "decisions", "1 unresolved map(s)", "/bench-shape-idea")}},
		{name: "one ready map", signal: "decisions", detail: "1 ready map", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "decisions/my map.md", readyMap(), 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(6, "decisions", "1 ready map(s)", "/bench-write-spec decisions/my map.md")},
			board: "▶ /bench-write-spec decisions/my map.md  (decisions)\n  decisions  1 ready map(s)                 → /bench-write-spec decisions/my map.md\n",
			route: &RouteResult{Lead: testSignal(6, "decisions", "1 ready map(s)", "/bench-write-spec decisions/my map.md")}},
		{name: "multiple ready maps", signal: "decisions", detail: "2 ready map", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "decisions/one.md", readyMap(), 0o644)
			write(t, root, "decisions/two.md", readyMap(), 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "map scan failure", signal: "decisions", detail: "unknown", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			if err := syscall.Mkfifo(filepath.Join(root, "decisions"), 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
			return root, Query{}
		}, exact: []Signal{testSignal(6, "decisions", "unknown (decisions is wrong-type)", "bench maps")}},
		{name: "implemented spec", signal: "specs", detail: "awaiting retirement", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "specs/implemented/spec.md", "Status: implemented\n", 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "orphaned review", signal: "reviews", detail: "orphaned review", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "reviews/orphan.md", "pickup\n", 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{{Severity: 9, Name: "reviews", Detail: "1 orphaned review pickup", Action: ""}}},
		{name: "tickets-only residue", signal: "specs", detail: "tickets-only spec folder", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, "specs/landed-ticket/tickets/one.md", "ticket\n", 0o644)
			commit(t, root)
			return root, Query{}
		}, exact: []Signal{testSignal(11, "specs", "1 tickets-only spec folder", "bench worktree land --spec <slug>")}},
		{name: "roadmap reconcile", signal: "roadmap", detail: "retired spec", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			write(t, root, roadmap.RoadmapFile, "specs/retired/spec.md\n", 0o644)
			commit(t, root)
			return root, Query{}
		}},
		{name: "roadmap scan failure", signal: "roadmap", detail: "unknown", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			if err := syscall.Mkfifo(filepath.Join(root, roadmap.RoadmapFile), 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
			return root, Query{}
		}},
		{name: "handoff", signal: "handoff", detail: "commit behind", setup: func(t *testing.T) (string, Query) {
			root := cleanRepo(t)
			commitHandoff(t, root, "current")
			commitFile(t, root, "after.txt")
			return root, Query{}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, query := tc.setup(t)
			matched := 0
			var exact []Signal
			producedSignals := SignalsWith(root, query)
			for _, produced := range producedSignals {
				if produced.Name == tc.signal {
					exact = append(exact, produced)
				}
				if produced.Name == tc.signal && (tc.detail == "" || strings.Contains(produced.Detail, tc.detail)) {
					matched++
				}
				if produced.Action == "" {
					if produced.actionID != noAction || produced.invocable() {
						t.Errorf("%s empty board action is not explicitly advisory", produced.Name)
					}
					continue
				}
				if !produced.invocable() || !IsInvocable(produced.Action) {
					t.Errorf("%s board action %q is not typed and parser-invocable", produced.Name, produced.Action)
				}
			}
			want := tc.count
			if want == 0 {
				want = 1
			}
			if matched != want {
				t.Fatalf("fixture produced %d matching %s row(s), want %d", matched, tc.signal, want)
			}
			if tc.exact != nil && !reflect.DeepEqual(exact, tc.exact) {
				t.Fatalf("%s rows = %#v, want %#v", tc.signal, exact, tc.exact)
			}
			if tc.board != "" {
				if got := render(root, false); got != tc.board {
					t.Fatalf("board = %q, want %q", got, tc.board)
				}
			}
			if tc.route != nil {
				if got := RouteFor(root, producedSignals, HarnessClaude); !reflect.DeepEqual(got, *tc.route) {
					t.Fatalf("route = %#v, want %#v", got, *tc.route)
				}
			}
		})
	}
}

func TestGitSignalIgnoresLandedUnclaimedAssignmentBranch(t *testing.T) {
	t.Parallel()
	root := initRepo(t)
	gitRun(t, root, "commit", "--allow-empty", "-m", "base")
	gitRun(t, root, "branch", "-M", "main")
	gitRun(t, root, "branch", "bench/assign/orphan/landed")
	for _, signal := range SignalsWith(root, Query{}) {
		if signal.Name == "git" {
			t.Fatalf("landed unclaimed assignment branch produced git signal %#v", signal)
		}
	}
}
