package axi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestAXIOutlineContracts drives the built `bench outline` in throwaway fixture repos:
// the walk, path scoping, TOON shape, empty state, exit codes, and hostile inputs. The
// pure per-language parser is unit-tested in internal/outline; these are the black-box
// command-surface rows of the acceptance coverage map.
func TestAXIOutlineContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI outline no-arg header and rows", testAXIOutlineNoArg)
	contract.RunParallel(t, "AXI outline definitive empty table", testAXIOutlineEmpty)
	contract.RunParallel(t, "AXI outline not-in-repo error", testAXIOutlineNotInRepo)
	contract.RunParallel(t, "AXI outline usage/exit-2", testAXIOutlineUsage)
	contract.RunParallel(t, "AXI outline file-path scoping", testAXIOutlineFileScope)
	contract.RunParallel(t, "AXI outline directory-path scoping", testAXIOutlineDirScope)
	contract.RunParallel(t, "AXI outline subdir root-relative", testAXIOutlineSubdir)
	contract.RunParallel(t, "AXI outline untracked and binary skipped", testAXIOutlineUntrackedBinary)
	contract.RunParallel(t, "AXI outline tracked symlink skipped", testAXIOutlineSymlinkSkipped)
	contract.RunParallel(t, "AXI outline control-byte drops one row", testAXIOutlineControlByte)
	contract.RunParallel(t, "AXI outline file-then-line ordering", testAXIOutlineOrdering)
	contract.RunParallel(t, "AXI outline help promise", testAXIOutlineHelpPromise)
	contract.RunParallel(t, "AXI outline literal spaced/glob path", testAXIOutlineLiteralPath)
	contract.RunParallel(t, "AXI outline empty tracked file", testAXIOutlineEmptyFile)
	contract.RunParallel(t, "AXI outline git failure", testAXIOutlineGitFailure)
	contract.RunParallel(t, "AXI outline re-run idempotency", testAXIOutlineIdempotent)
}

// stageAll stages every current file so `git ls-files` sees it; outline reads the
// index, so an unstaged file is invisible to the walk.
func stageAll(f contract.Fixture) {
	f.Git("add", "-A")
}

func testAXIOutlineNoArg(t *testing.T) {
	contract.NoteContractFailure(t, "AXI outline no-arg contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("pkg/one.go", "package pkg\n\nfunc Alpha() {}\n")
	stageAll(f)

	out := f.Bench("outline")
	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "outline[1]{file,line,kind,name}:")
	requireAXILine(t, out.Stdout, "  pkg/one.go,\"3\",func,Alpha")
}

func testAXIOutlineEmpty(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("notes.txt", "no declarations here\n")
	stageAll(f)

	out := f.Bench("outline")
	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "outline[0]{file,line,kind,name}:")
}

func testAXIOutlineNotInRepo(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())

	out := f.Bench("outline")
	out.RequireExit(1)
	requireContainsFold(t, out.Stdout, "not in a git repository")
}

func testAXIOutlineUsage(t *testing.T) {
	f := contract.NewFixture(t)

	flag := f.Bench("outline", "--bogus")
	flag.RequireExit(2)
	requireContainsFold(t, flag.Stdout, "usage")

	second := f.Bench("outline", "a", "b")
	second.RequireExit(2)
	requireContainsFold(t, second.Stdout, "usage")
}

func testAXIOutlineFileScope(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("pkg/one.go", "package pkg\nfunc One() {}\n")
	f.WriteFile("pkg/two.go", "package pkg\nfunc Two() {}\n")
	stageAll(f)

	out := f.Bench("outline", "pkg/one.go")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  pkg/one.go,\"2\",func,One")
	requireNoAXILineMatching(t, out.Stdout, `two\.go`)
}

func testAXIOutlineDirScope(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("pkg/one.go", "package pkg\nfunc One() {}\n")
	f.WriteFile("top.go", "package top\nfunc Top() {}\n")
	stageAll(f)

	out := f.Bench("outline", "pkg")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  pkg/one.go,\"2\",func,One")
	requireNoAXILineMatching(t, out.Stdout, `^  top\.go,`)
}

func testAXIOutlineSubdir(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("pkg/one.go", "package pkg\nfunc One() {}\n")
	stageAll(f)
	sub := filepath.Join(f.Root, "pkg")

	out := runBenchInDir(t, f, sub, "outline")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  pkg/one.go,\"2\",func,One")
}

func testAXIOutlineUntrackedBinary(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("clean.go", "package x\nfunc Clean() {}\n")
	f.WriteFile("bin.go", "package x\nfunc Binny() {}\n\x00\n") // NUL → binary, tracked
	stageAll(f)
	f.WriteFile("loose.go", "package x\nfunc Loose() {}\n") // never staged → untracked

	out := f.Bench("outline")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  clean.go,\"2\",func,Clean")
	requireNoAXILineMatching(t, out.Stdout, `Binny`)
	requireNoAXILineMatching(t, out.Stdout, `Loose`)
}

// A tracked symlink's git content is its target string, not the target's
// declarations — following it would index the target's symbols under the symlink's
// path, emitting file:line anchors that don't hold. The walk must skip non-regular
// entries so the target is indexed once, under its own path.
func testAXIOutlineSymlinkSkipped(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("real.go", "package pkg\n\ntype Bar struct{}\n")
	if err := os.Symlink("real.go", filepath.Join(f.Root, "link.go")); err != nil {
		t.Skipf("symlinks unavailable on this filesystem: %v", err)
	}
	stageAll(f)

	out := f.Bench("outline")
	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "outline[1]{file,line,kind,name}:")
	requireAXILine(t, out.Stdout, "  real.go,\"3\",type,Bar")
	if strings.Contains(out.Stdout, "  link.go,") && !strings.Contains(out.Stdout, "  link.go,nonregular") {
		t.Fatalf("tracked symlink was indexed as its target:\n%s", out.Stdout)
	}
	requireAXILine(t, out.Stdout, "  link.go,nonregular")
}

func testAXIOutlineControlByte(t *testing.T) {
	contract.NoteContractFailure(t, "AXI outline control-byte contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("doc.md", "# good heading\n## bad\x1bheading\n")
	stageAll(f)

	out := f.Bench("outline")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  doc.md,\"1\",heading,good heading")
	requireNoAXILineMatching(t, out.Stdout, `bad`)
}

func testAXIOutlineOrdering(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("b.go", "package x\nfunc Bee() {}\n")
	f.WriteFile("a.go", "package x\n\nfunc First() {}\nfunc Second() {}\n")
	stageAll(f)

	out := f.Bench("outline")
	out.RequireExit(0)
	requireLineOrder(t, out.Stdout,
		"  a.go,\"3\",func,First",
		"  a.go,\"4\",func,Second",
		"  b.go,\"2\",func,Bee",
	)
}

func testAXIOutlineHelpPromise(t *testing.T) {
	f := contract.NewFixture(t)

	out := f.Bench("outline", "-h")
	out.RequireExit(0)
	requireContainsFold(t, out.Stdout, "locates")
	requireContainsFold(t, out.Stdout, "does not identify")
	requireContainsFold(t, out.Stdout, "blessed")
}

func testAXIOutlineLiteralPath(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("weird [dir]/inside.go", "package w\nfunc Inside() {}\n")
	f.WriteFile("other.go", "package o\nfunc Other() {}\n")
	stageAll(f)

	out := f.Bench("outline", "weird [dir]")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  \"weird [dir]/inside.go\",\"2\",func,Inside")
	requireNoAXILineMatching(t, out.Stdout, `Other`)
}

func testAXIOutlineEmptyFile(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("empty.go", "")
	f.WriteFile("real.go", "package x\nfunc Real() {}\n")
	stageAll(f)

	out := f.Bench("outline")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  real.go,\"2\",func,Real")
	requireNoAXILineMatching(t, out.Stdout, `^  empty\.go,`)
}

func testAXIOutlineGitFailure(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("x.go", "package x\nfunc X() {}\n")
	stageAll(f)
	f.WriteFile(".git/index", "this is not a valid git index file at all\x00\x01\x02")

	out := f.Bench("outline")
	out.RequireExit(1)
	requireContainsFold(t, out.Stdout, "error:")
	requireContainsFold(t, out.Stdout, "ls-files")
}

func testAXIOutlineIdempotent(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("a.go", "package x\nfunc A() {}\n")
	f.WriteFile("b.md", "# Heading\n")
	stageAll(f)

	first := f.Bench("outline")
	first.RequireExit(0)
	second := f.Bench("outline")
	second.RequireExit(0)
	if first.Stdout != second.Stdout {
		t.Fatalf("outline output not byte-identical across runs\nfirst:\n%s\nsecond:\n%s", first.Stdout, second.Stdout)
	}
}

// requireLineOrder asserts the given lines each appear in stdout in the given relative
// order (not necessarily adjacent), pinning the file-then-line determinism.
func requireLineOrder(t *testing.T, out string, lines ...string) {
	t.Helper()
	body := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	idx := func(want string) int {
		for i, got := range body {
			if got == want {
				return i
			}
		}
		return -1
	}
	prev := -1
	for _, want := range lines {
		at := idx(want)
		if at < 0 {
			t.Fatalf("missing ordered line %q\noutput:\n%s", want, out)
		}
		if at <= prev {
			t.Fatalf("line %q out of order (index %d after %d)\noutput:\n%s", want, at, prev, out)
		}
		prev = at
	}
}
