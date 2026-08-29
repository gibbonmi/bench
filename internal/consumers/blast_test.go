package consumers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// The blast fixture is one repository with two commits. Every file's two revisions are
// named here once and are used twice: the commit writes them, and the stubbed loader
// type-checks the tip revision in process. So the git pair and the typed tip can never
// describe different trees.
//
// Between the two commits, target.Changed and target.Quiet gain new bodies, target.Gone
// is deleted, and inside/ is edited so it sits inside the diff while outside/ does not.
// Quiet's only consumer is the edited inside/ file, so it is the planted symbol whose
// consumers are all touched.
const (
	baseTarget = "package target\n" +
		"\n" +
		"func Keep() {}\n" +
		"\n" +
		"func Gone() {}\n" +
		"\n" +
		"func Changed() int { return 1 }\n" +
		"\n" +
		"func Quiet() {}\n"
	tipTarget = "package target\n" +
		"\n" +
		"func Keep() {}\n" +
		"\n" +
		"func Changed() int { return 2 }\n" +
		"\n" +
		"func Quiet() int { return 3 }\n"
	baseInside = "package inside\n" +
		"\n" +
		"import \"example.com/fx/target\"\n" +
		"\n" +
		"func UseChanged() int { return target.Changed() }\n" +
		"\n" +
		"func UseGone() { target.Gone() }\n" +
		"\n" +
		"func UseQuiet() { target.Quiet() }\n"
	tipInside = "package inside\n" +
		"\n" +
		"import \"example.com/fx/target\"\n" +
		"\n" +
		"func UseChanged() int { return target.Changed() }\n" +
		"\n" +
		"func UseGone() {}\n" +
		"\n" +
		"func UseQuiet() int { return target.Quiet() }\n"
	sameOutside = "package outside\n" +
		"\n" +
		"import \"example.com/fx/target\"\n" +
		"\n" +
		"func UseChanged() int { return target.Changed() }\n"
)

// blastRepo commits the fixture's two revisions and makes the tip the process working
// directory. It returns the base sha; the tip is HEAD, which is the only revision the
// blast mode enumerates in place.
func blastRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	base := commitTree(t, root, map[string]string{
		"target/target.go":   baseTarget,
		"inside/inside.go":   baseInside,
		"outside/outside.go": sameOutside,
	}, "base")
	commitTree(t, root, map[string]string{
		"target/target.go": tipTarget,
		"inside/inside.go": tipInside,
	}, "tip")
	t.Chdir(root)
	return base
}

// commitTree writes one revision of a fixture repository and commits it, returning the new
// sha. Every blast fixture builds its pair through this one helper.
func commitTree(t *testing.T, root string, files map[string]string, message string) string {
	t.Helper()
	for rel, body := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	if _, err := git.Output("-C", root, "add", "-A"); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if _, err := git.Output("-C", root, "commit", "-q", "-m", message); err != nil {
		t.Fatalf("git commit: %v", err)
	}
	sha, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}
	return sha
}

// tipFixture is the typed tip tree the stubbed loader returns for the blast fixture.
func tipFixture(root string) []fixturePkg {
	return []fixturePkg{
		{path: "example.com/fx/target", files: map[string]string{root + "/target/target.go": tipTarget}},
		{path: "example.com/fx/inside", files: map[string]string{root + "/inside/inside.go": tipInside}},
		{path: "example.com/fx/outside", files: map[string]string{root + "/outside/outside.go": sameOutside}},
	}
}

// head reports the tip sha of the repository the test is running in.
func head(t *testing.T) string {
	t.Helper()
	sha, err := git.Output("rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("rev-parse HEAD: %v", err)
	}
	return sha
}

// BL1 (story 16): the blast enumerates the consumers of every touched declaration, in
// both the edited package and the untouched one, and names no declaration the diff left
// alone.
func TestBlastEnumeratesConsumersOfEveryTouchedDeclaration(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	for _, want := range []string{
		"target.Changed,inside/inside.go,5,true\n",
		"target.Changed,outside/outside.go,5,false\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout = %q, want the row %q", out, want)
		}
	}
	if strings.Contains(out, "target.Keep") {
		t.Fatalf("stdout = %q, want no row for the untouched declaration target.Keep", out)
	}
}

// BL2 (story 17): a deleted declaration is reported in its own table at its base
// position, and it enumerates nothing.
func TestDeletedDeclarationEmitsOneDeletedRowAndNoConsumerRows(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "blast_deleted[1]{changed_symbol,base_file,base_line}:\n  target.Gone,target/target.go,5\n") {
		t.Fatalf("stdout = %q, want exactly one blast_deleted row for target.Gone", out)
	}
	if strings.Contains(out, "target.Gone,inside") {
		t.Fatalf("stdout = %q, want no consumer rows for the deleted declaration", out)
	}
}

// BL3 (story 18): a pair with no diff answers with the definitive empty table at exit 0,
// and it is a terminal read.
func TestIdenticalPairEmitsDefinitiveEmptyBlastTable(t *testing.T) {
	blastRepo(t)
	stubLoad(t, tipFixture)
	tip := head(t)
	out, code := run(t, "--changed", "--base", tip, "--source-tip", tip)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "blast[0]{changed_symbol,file,line,touched}:\n") {
		t.Fatalf("stdout = %q, want the definitive empty blast table", out)
	}
	if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("stdout = %q, want a terminal empty help envelope", out)
	}
}

// BL4 (story 19): the pair is frozen, so a recomputation at review time byte-matches.
func TestTwoRunsOverOneFrozenPairAreByteEqual(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	tip := head(t)
	first, code := run(t, "--changed", "--base", base, "--source-tip", tip, "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, first)
	}
	second, _ := run(t, "--changed", "--base", base, "--source-tip", tip, "--full")
	if first != second {
		t.Fatalf("second run = %q, want byte-equal to the first %q", second, first)
	}
}

// BL5 (story 20): the revision grammar is the `bench test` grammar, so a source tip
// without a base is a usage error rather than a pair against a defaulted base.
func TestSourceTipWithoutBaseExitsTwoWithUsage(t *testing.T) {
	out, code := run(t, "--changed", "--source-tip", "HEAD")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "usage: bench consumers") {
		t.Fatalf("stdout = %q, want a usage line", out)
	}
}

// BL6 (story 31): touched separates the consumers already inside the diff from the ones
// the review must still walk.
func TestTouchedMarksConsumersInsideTheDiffOnly(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, ",inside/inside.go,5,true\n") {
		t.Fatalf("stdout = %q, want the inside consumer marked touched", out)
	}
	if !strings.Contains(out, ",outside/outside.go,5,false\n") {
		t.Fatalf("stdout = %q, want the outside consumer marked untouched", out)
	}
}

// BL7 (story 31): the envelope offers one walk per symbol with an untouched consumer, and
// none for target.Quiet, whose only consumer the diff already edited.
func TestBlastOffersOneFullActionPerSymbolWithAnUntouchedConsumer(t *testing.T) {
	base := blastRepo(t)
	stubLoad(t, tipFixture)
	out, code := run(t, "--changed", "--base", base, "--source-tip", head(t), "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.HasSuffix(out, "help[1]{cmd,why}:\n  bench consumers target.Changed --full,walk the consumers outside the diff\n") {
		t.Fatalf("stdout = %q, want exactly one per-symbol walk action", out)
	}
}
