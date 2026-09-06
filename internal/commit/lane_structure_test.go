package commit

// The fast lane's growth ratchet at the worktree commit. The row drives the real command
// against a fixture repository whose lane declares a structure check, so the observed
// evidence is the verb's own stdout, exit code, and branch ref. The lane's other rows live
// in lane_test.go and share this package's fixture helpers.

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// structureCheck stands in for `bench structure --growth <base>` the way proseCheck stands
// in for `bench gate-prose`: it takes the base as its one operand and reds a Go file that
// exceeds the budget and also gained lines since that base. It reads the base through git,
// so an unreplaced base token names no revision and every file reads as new.
const structureCheck = "#!/bin/sh\nbase=\"$1\"\nmax=\"${BENCH_MAX_LINES:?BENCH_MAX_LINES unset}\"\nstatus=0\nfor f in $(find . -name '*.go' -type f | sed 's|^\\./||'); do\n  tip=$(wc -l < \"$f\")\n  was=$(git show \"$base:$f\" 2>/dev/null | wc -l)\n  [ \"$tip\" -gt \"$max\" ] || continue\n  [ \"$tip\" -gt \"$was\" ] || continue\n  echo \"FILE GREW       $tip lines, was $was (max $max)   $f\"\n  status=1\ndone\nexit $status\n"

// structureLaneManifest declares the controllable check beside the growth check. The
// growth check carries the base token, so the row observes the resolver's own substitution
// rather than a base the fixture spelled itself.
func structureLaneManifest(t *testing.T) string {
	t.Helper()
	doc := map[string]any{
		"phases": []any{map[string]any{"name": "build", "argv": []string{"go", "build", "./..."}}},
		"lane": []any{
			map[string]any{"name": "check", "argv": []string{"sh", ".bench/check.sh"}},
			map[string]any{"name": "structure", "argv": []string{"sh", ".bench/structure.sh", gate.LaneBaseToken}},
		},
	}
	encoded, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

// growthGoFile is a gofmt-clean Go file of exactly n lines, so a row states the count it
// wants rather than counting a literal.
func growthGoFile(n int) string {
	return "package fixture\n\nfunc Big() int {\n\tn := 0\n" +
		strings.Repeat("\tn++\n", n-6) + "\treturn n\n}\n"
}

// SR52: the fast lane runs the growth ratchet, so a commit that raises an over-budget Go
// file's count is refused under the structure check with the check's own `FILE GREW` line,
// and a commit that lowers that count passes. The base commit holds the file at 12 lines
// against a budget of 10, so existing debt alone never refuses the commit.
func TestLaneStructureRefusesACommitThatGrowsAnOverBudgetGoFile(t *testing.T) {
	for _, tc := range []struct {
		name     string
		lines    int
		wantExit int
	}{
		{name: "raised", lines: 15, wantExit: 1},
		{name: "lowered", lines: 11, wantExit: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("BENCH_MAX_LINES", "10")
			root, before := laneRepo(t, 0, func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, ".bench", "structure.sh"), structureCheck, 0o755)
				mustWrite(t, filepath.Join(root, ".bench", "phases.json"), structureLaneManifest(t), 0o644)
				mustWrite(t, filepath.Join(root, "big.go"), growthGoFile(12), 0o644)
			})
			runGit(t, root, "reset", "-q", "--hard", "HEAD")
			mustWrite(t, filepath.Join(root, "big.go"), growthGoFile(tc.lines), 0o644)

			code, stdout, stderr := runCommand(t, root, "-m", "m", "big.go")
			if code != tc.wantExit {
				t.Fatalf("exit = %d, want %d; stdout=%q stderr=%q", code, tc.wantExit, stdout, stderr)
			}
			if tc.wantExit == 0 {
				if !strings.Contains(stdout, "outcome=pass") {
					t.Fatalf("stdout = %q, want a lane pass: the file lost lines since the base", stdout)
				}
				return
			}
			if !strings.Contains(stdout, "lane{outcome=fail,check=structure}") {
				t.Errorf("stdout = %q, want the structure check named", stdout)
			}
			if !strings.Contains(stdout, "FILE GREW") {
				t.Errorf("stdout = %q, want the growth row the check printed", stdout)
			}
			if after := head(t, root); after != before {
				t.Fatalf("the branch ref moved from %s to %s on a growth refusal", before, after)
			}
		})
	}
}
