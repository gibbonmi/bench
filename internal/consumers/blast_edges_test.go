package consumers

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// The deletion fixture is one repository whose tip only removes lines. Git reports such a
// hunk with a zero-count tip side, which is the one shape a naive span rule reads as no tip
// edit at all. The line numbers are load-bearing. The first deletion sits inside `Changed`,
// which spans tip lines 3 to 6. The second removes the blank line before `Later`, whose tip
// line 7 is the second of the two lines that removed run sat between.
const (
	deletionBaseTarget = "package target\n" + // 1
		"\n" + // 2
		"func Changed() int {\n" + // 3
		"\tn := 1\n" + // 4
		"\tn++\n" + // 5
		"\treturn n\n" + // 6
		"}\n" + // 7
		"\n" + // 8
		"func Later() {}\n" // 9
	deletionTipTarget = "package target\n" + // 1
		"\n" + // 2
		"func Changed() int {\n" + // 3
		"\tn := 1\n" + // 4
		"\treturn n\n" + // 5
		"}\n" + // 6
		"func Later() {}\n" // 7
	deletionConsumer = "package outside\n" + // 1
		"\n" + // 2
		"import \"example.com/fx/target\"\n" + // 3
		"\n" + // 4
		"func UseChanged() int { return target.Changed() }\n" + // 5
		"\n" + // 6
		"func UseLater() { target.Later() }\n" // 7
)

// deletionRepo commits the deletion fixture's two revisions and returns the base sha.
func deletionRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	base := commitTree(t, root, map[string]string{
		"target/target.go":   deletionBaseTarget,
		"outside/outside.go": deletionConsumer,
	}, "base")
	commitTree(t, root, map[string]string{"target/target.go": deletionTipTarget}, "tip")
	t.Chdir(root)
	return base
}

// deletionTipFixture is the typed tip tree the stubbed loader returns for that repository.
func deletionTipFixture(root string) []fixturePkg {
	return []fixturePkg{
		{path: "example.com/fx/target", files: map[string]string{root + "/target/target.go": deletionTipTarget}},
		{path: "example.com/fx/outside", files: map[string]string{root + "/outside/outside.go": deletionConsumer}},
	}
}

// R1 (story 16): an edit that only deletes a body line still touches the declaration it
// edited, so the blast enumerates that declaration's consumers.
func TestDeletionOnlyEditTouchesTheSurvivingDeclaration(t *testing.T) {
	base := deletionRepo(t)
	stubLoad(t, deletionTipFixture)
	out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	for _, want := range []string{
		"target.Changed,outside/outside.go,5,false\n",
		"target.Later,outside/outside.go,7,false\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout = %q, want the deletion-only row %q", out, want)
		}
	}
}

// R3 (story 19): blast rows come only from the frozen pair, so a dirty checkout refuses
// rather than positioning rows in working-tree bytes the pair never froze.
func TestDirtyCheckoutRefusesTheChangedQuery(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	tip := head(t)
	out, code := run(t, "--changed", "--base", base, "--source-tip", tip, "--full")
	if code != 0 {
		t.Fatalf("clean run exit = %d, want 0; out=%q", code, out)
	}
	root, err := git.Root()
	if err != nil {
		t.Fatalf("git root: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "scratch.txt"), []byte("probe\n"), 0o644); err != nil {
		t.Fatalf("dirty the checkout: %v", err)
	}
	out, code = run(t, "--changed", "--base", base, "--source-tip", tip, "--full")
	if code != 1 {
		t.Fatalf("dirty run exit = %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: ") || !strings.Contains(out, "dirty") {
		t.Fatalf("stdout = %q, want a structured refusal naming the dirty checkout", out)
	}
	if !strings.Contains(out, "commit or clean the checkout") {
		t.Fatalf("stdout = %q, want the remedy named", out)
	}
	if strings.Contains(out, "blast[") || strings.Contains(out, "citation[") {
		t.Fatalf("stdout = %q, want no result block and no citation row", out)
	}
}

// R4 (story 18): a pair that changed no Go file answers with the definitive empty table
// before the loader runs, so a tree the loader cannot load still gets that answer.
func TestIdenticalPairAnswersWithoutLoadingPackages(t *testing.T) {
	blastRepo(t)
	original := load
	load = func(string, ...string) ([]*Package, error) {
		return nil, errors.New("loader must not run for an empty changed set")
	}
	t.Cleanup(func() { load = original })
	tip := head(t)
	out, code := run(t, "--changed", "--base", tip, "--source-tip", tip)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "blast[0]{changed_symbol,file,line,touched}:\n") {
		t.Fatalf("stdout = %q, want the definitive empty blast table", out)
	}
	if !strings.Contains(out, "meta[1]{packages,files,matches,rows,truncated}:\n  0,0,0,0,false\n") {
		t.Fatalf("stdout = %q, want meta packages=0", out)
	}
}
