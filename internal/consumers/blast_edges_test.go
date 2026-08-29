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

// An edit that only deletes a body line still touches the declaration it
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

// BL8 (story 19): blast rows come only from the frozen pair, so a dirty checkout refuses
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

// A pair that changed no Go file answers with the definitive empty table
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

// The escaped-path fixture is one repository whose Go file sits at a path carrying an
// escape byte. Git C-quotes such a header path whatever core.quotePath says, so the file
// reaches the blast only when the header parse unquotes it. A clean target file and a
// clean consumer file bracket the poisoned one: the poisoned file's own rows must drop,
// and every clean row must stand.
const (
	escapedDir  = "poi\x1bson"
	escapedFile = escapedDir + "/poison.go"

	escapedBaseTarget = "package target\n" + // 1
		"\n" + // 2
		"func Changed() int { return 1 }\n" + // 3
		"\n" + // 4
		"func Gone() {}\n" // 5
	escapedTipTarget = "package target\n" + // 1
		"\n" + // 2
		"func Changed() int { return 2 }\n" // 3
	escapedBasePoison = "package poison\n" + // 1
		"\n" + // 2
		"import \"example.com/fx/target\"\n" + // 3
		"\n" + // 4
		"func UseChanged() int { return target.Changed() }\n" + // 5
		"\n" + // 6
		"func Dropped() {}\n" + // 7
		"\n" + // 8
		"func Shifted() int { return 1 }\n" // 9
	escapedTipPoison = "package poison\n" + // 1
		"\n" + // 2
		"import \"example.com/fx/target\"\n" + // 3
		"\n" + // 4
		"func UseChanged() int { return target.Changed() }\n" + // 5
		"\n" + // 6
		"func Shifted() int { return 2 }\n" // 7
	escapedClean = "package clean\n" + // 1
		"\n" + // 2
		"import \"example.com/fx/poison\"\n" + // 3
		"\n" + // 4
		"func UseShifted() int { return poison.Shifted() }\n" // 5
)

// escapedPathRepo commits the escaped-path fixture's two revisions and returns the base
// sha. The tip both edits a declaration in the poisoned file and deletes one from it.
func escapedPathRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	base := commitTree(t, root, map[string]string{
		"target/target.go": escapedBaseTarget,
		escapedFile:        escapedBasePoison,
		"clean/clean.go":   escapedClean,
	}, "base")
	commitTree(t, root, map[string]string{
		"target/target.go": escapedTipTarget,
		escapedFile:        escapedTipPoison,
	}, "tip")
	t.Chdir(root)
	return base
}

// escapedTipFixture is the typed tip tree the stubbed loader returns for that repository.
func escapedTipFixture(root string) []fixturePkg {
	return []fixturePkg{
		{path: "example.com/fx/target", files: map[string]string{root + "/target/target.go": escapedTipTarget}},
		{path: "example.com/fx/poison", files: map[string]string{root + "/" + escapedFile: escapedTipPoison}},
		{path: "example.com/fx/clean", files: map[string]string{root + "/clean/clean.go": escapedClean}},
	}
}

// A declaration in a C-quoted path is still a changed declaration, so its clean consumer
// earns a row. Only the rows that name the poisoned path itself drop, and the answer says
// it truncated.
func TestQuotedHeaderPathStillContributesItsChangedDeclarations(t *testing.T) {
	base := escapedPathRepo(t)
	stubLoad(t, escapedTipFixture)
	out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "blast[1]{changed_symbol,file,line,touched}:\n  poison.Shifted,clean/clean.go,5,false\n") {
		t.Fatalf("stdout = %q, want the clean consumer row of the escaped path's edited declaration", out)
	}
	if !strings.Contains(out, "blast_deleted[1]{changed_symbol,base_file,base_line}:\n  target.Gone,target/target.go,5\n") {
		t.Fatalf("stdout = %q, want the clean deleted row kept and the escaped-path one dropped", out)
	}
	if strings.Contains(out, "\x1b") {
		t.Fatalf("stdout = %q, want no escape byte in the response", out)
	}
	const metaHead = "meta[1]{packages,files,matches,rows,truncated}:\n  "
	i := strings.Index(out, metaHead)
	if i < 0 {
		t.Fatalf("stdout = %q, want a meta row", out)
	}
	meta := out[i+len(metaHead):]
	meta = meta[:strings.Index(meta, "\n")]
	if !strings.HasSuffix(meta, ",true") {
		t.Fatalf("meta = %q, want truncated=true", meta)
	}
}

// BL9 (story 19): enumeration runs in this checkout, so a tip that is not HEAD refuses
// rather than positioning rows in a tree the pair does not name. The refusal names both
// commits, so the agent can see which one to check out.
func TestSourceTipThatIsNotHeadRefusesAndNamesBothCommits(t *testing.T) {
	base := blastRepo(t)
	tip := head(t)
	root, err := git.Root()
	if err != nil {
		t.Fatalf("git root: %v", err)
	}
	if _, err := git.Output("-C", root, "checkout", "--detach", base); err != nil {
		t.Fatalf("detach at the base: %v", err)
	}
	out, code := run(t, "--changed", "--base", base, "--source-tip", tip, "--full")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: ") {
		t.Fatalf("stdout = %q, want a structured refusal", out)
	}
	if !strings.Contains(out, tip) || !strings.Contains(out, base) {
		t.Fatalf("stdout = %q, want both the tip %s and the checkout HEAD %s named", out, tip, base)
	}
	if strings.Contains(out, "blast[") || strings.Contains(out, "citation[") {
		t.Fatalf("stdout = %q, want no result block and no citation row", out)
	}
}
