package adopt

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/learnings"
)

func TestRewriteAgentsBlockEdges(t *testing.T) {
	block := BenchAgentsBlock()
	cases := []struct {
		name    string
		in      string
		wantHas []string
		wantErr string
	}{
		{
			name:    "append to project content without trailing newline",
			in:      "PROJECT",
			wantHas: []string{"PROJECT\n", "<!-- bench:start -->", "<!-- bench:end -->"},
		},
		{
			name: "replace only the unfenced managed block",
			in: strings.Join([]string{
				"before",
				"<!-- bench:start -->",
				"old managed",
				"<!-- bench:end -->",
				"after",
				"",
			}, "\n"),
			wantHas: []string{"before\n", block, "\nafter\n"},
		},
		{
			name: "preserve fenced marker examples",
			in: strings.Join([]string{
				"# Project",
				"```",
				"<!-- bench:start -->",
				"example",
				"<!-- bench:end -->",
				"```",
				"KEEP",
				"",
			}, "\n"),
			wantHas: []string{"example", "KEEP", "<!-- bench:start -->"},
		},
		{
			name: "reject reversed markers",
			in: strings.Join([]string{
				"PROJECT BEFORE",
				"<!-- bench:end -->",
				"PROJECT MIDDLE",
				"<!-- bench:start -->",
				"PROJECT AFTER",
				"",
			}, "\n"),
			wantErr: "malformed",
		},
		{
			name: "reject unclosed fence around markers",
			in: strings.Join([]string{
				"Broken docs:",
				"```",
				"<!-- bench:start -->",
				"<!-- bench:end -->",
				"KEEP",
				"",
			}, "\n"),
			wantErr: "fence",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := RewriteAgentsBlock(c.in)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want substring %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("RewriteAgentsBlock err = %v", err)
			}
			for _, want := range c.wantHas {
				if !strings.Contains(got, want) {
					t.Fatalf("output missing %q:\n%s", want, got)
				}
			}
			if strings.Count(got, "## Bench") != 1 {
				t.Fatalf("managed block not exactly once:\n%s", got)
			}
		})
	}
}

func TestManifestParseEdges(t *testing.T) {
	dir := t.TempDir()
	missing, err := ReadManifest(filepath.Join(dir, "missing.tsv"))
	if err != nil {
		t.Fatalf("missing manifest err = %v", err)
	}
	if missing.Hash(".bench/BENCH.md") != "" || missing.KitVersion != "" {
		t.Fatalf("missing manifest = %+v, want empty", missing)
	}

	path := filepath.Join(dir, "link-manifest.tsv")
	data := strings.Join([]string{
		"#kit\t0.2.0",
		".bench/BENCH.md\told",
		"#comment\tignored",
		".bench/BENCH.md\tnew",
		".agents/commands/bench-implement-spec.md\tcmdhash",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadManifest(path)
	if err != nil {
		t.Fatalf("ReadManifest err = %v", err)
	}
	if got.KitVersion != "0.2.0" {
		t.Fatalf("KitVersion = %q, want 0.2.0", got.KitVersion)
	}
	if got.Hash(".bench/BENCH.md") != "new" {
		t.Fatalf("duplicate rel did not use last value: %+v", got)
	}
	if got.Hash("#kit") != "" {
		t.Fatalf("#kit parsed as a file row: %+v", got)
	}
}

func TestPromoteAllRollsBackOnDestinationSyncFailure(t *testing.T) {
	root := t.TempDir()
	stage := t.TempDir()
	existing := filepath.Join(root, "existing")
	if err := os.WriteFile(existing, []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldStage := filepath.Join(stage, "existing")
	freshStage := filepath.Join(stage, "fresh")
	if err := os.WriteFile(oldStage, []byte("after\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(freshStage, []byte("fresh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldSync := syncDirectory
	syncDirectory = func(string) error { return os.ErrPermission }
	t.Cleanup(func() { syncDirectory = oldSync })

	err := promoteAll(root, []stagedChange{
		{rel: "existing", stage: oldStage, backup: filepath.Join(stage, "existing.backup")},
		{rel: "fresh", stage: freshStage, backup: filepath.Join(stage, "fresh.backup")},
	})
	if err == nil {
		t.Fatal("promoteAll succeeded despite destination sync failure")
	}
	if got, readErr := os.ReadFile(existing); readErr != nil || string(got) != "before\n" {
		t.Fatalf("existing destination after sync failure = %q, %v; want original bytes", got, readErr)
	}
	if _, statErr := os.Lstat(filepath.Join(root, "fresh")); !os.IsNotExist(statErr) {
		t.Fatalf("fresh destination survived sync failure: %v", statErr)
	}
}

func TestAdapterTarget(t *testing.T) {
	cases := map[string]string{
		".claude/commands/bench-implement-spec.md":      "../../.agents/commands/bench-implement-spec.md",
		".claude/skills/bench-craft-seams/SKILL.md":     "../../../.agents/skills/bench-craft-seams/SKILL.md",
		".claude/skills/bench-craft-seams/refs/deep.md": "../../../../.agents/skills/bench-craft-seams/refs/deep.md",
	}
	for rel, want := range cases {
		got, ok := AdapterTarget(rel)
		if !ok {
			t.Fatalf("AdapterTarget(%q) not ok", rel)
		}
		if got != want {
			t.Fatalf("AdapterTarget(%q) = %q, want %q", rel, got, want)
		}
	}
	if _, ok := AdapterTarget(".bench/BENCH.md"); ok {
		t.Fatalf("non-adapter rel returned ok")
	}
}

// TestLinkOutsideGitRepoNamesGitRepository pins link's own remedy. Link's job is to
// create the linkage. The shared AXI query-command wording, "run inside a Bench-linked
// repo", is nonsensical here. The message must name a git repository or a git init
// command instead, and the shared toon.NotInRepo() phrasing must not appear.
func TestLinkOutsideGitRepoNamesGitRepository(t *testing.T) {
	t.Chdir(t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Link(nil, &stdout, &stderr, "1.0.0")
	if code != 1 {
		t.Fatalf("Link exit = %d, want 1", code)
	}
	got := stderr.String()
	if !strings.Contains(got, "git init") && !strings.Contains(strings.ToLower(got), "git repository") {
		t.Fatalf("stderr = %q, want it to name a git repository / git init", got)
	}
	if strings.Contains(got, "Bench-linked repo") {
		t.Fatalf("stderr = %q, must not use the shared AXI NotInRepo() phrasing", got)
	}
}

// TestLinkInKitSourceCheckoutRefuses pins link's refusal in the kit's own source tree.
// The kit authors the managed block and bin/bench.sh, so a link here writes a tracked
// launcher copy the shim then prefers. The refusal names the kit source checkout and
// sends the reader to the kit-side route, and it writes no destination file.
func TestLinkInKitSourceCheckoutRefuses(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	t.Setenv("BENCH_KIT", root)
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	code := Link(nil, &stdout, &stderr, "1.0.0")
	if code != 1 {
		t.Fatalf("Link in the kit source checkout = %d, want 1\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	got := stderr.String()
	if !strings.Contains(got, "kit source checkout") {
		t.Fatalf("stderr = %q, want it to name the kit source checkout", got)
	}
	if !strings.Contains(got, "bench doctor --fix") {
		t.Fatalf("stderr = %q, want it to name bench doctor --fix", got)
	}
	if _, err := os.Lstat(filepath.Join(root, ".bench", "bin", "bench.sh")); !os.IsNotExist(err) {
		t.Fatalf("launcher copy after the refusal = %v, want absent", err)
	}
}

// TestLinkRelinkStaysGreenInAConsumerRepo is the over-broad guard on the kit-checkout
// refusal. The first link writes the consumer copy of .bench/BENCH.md, so a predicate
// that tests for that marker instead of the kit checkout refuses the second run. A
// consumer repo relinks as often as it likes, whatever assets the first run left.
func TestLinkRelinkStaysGreenInAConsumerRepo(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(wd, "..", "..")))
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	if code := Link(nil, &stdout, &stderr, "1.0.0"); code != 0 {
		t.Fatalf("first Link = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if _, statErr := os.Stat(filepath.Join(root, ".bench", "BENCH.md")); statErr != nil {
		t.Fatalf("first Link left no .bench/BENCH.md: %v", statErr)
	}

	stdout.Reset()
	stderr.Reset()
	if code := Link(nil, &stdout, &stderr, "1.0.0"); code != 0 {
		t.Fatalf("relink = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "linked Bench into "+root) {
		t.Fatalf("relink stdout = %q, want the linked message for %q", got, root)
	}
}

func TestReportDoctorRowsReportsMalformedWorktreeAdmin(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	fifo := gittest.FIFOWorktreeAdmin(t, root, "doctor-fifo")
	t.Chdir(root)

	var stdout bytes.Buffer
	if !reportDoctorRows(&stdout) {
		t.Fatal("reportDoctorRows did not report the malformed worktree admin entry")
	}
	if got := stdout.String(); !strings.Contains(got, "red:") || !strings.Contains(got, "worktrees/doctor-fifo/gitdir") {
		t.Fatalf("doctor rows = %q, want red worktrees/doctor-fifo/gitdir row", got)
	}
	if info, err := os.Lstat(fifo); err != nil || info.Mode()&os.ModeNamedPipe == 0 {
		t.Fatalf("FIFO after doctor = %v, %v; want extant FIFO", info, err)
	}
}

func TestReportDoctorRowsLeavesHealthyWorktreeAdminSilent(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	var stdout bytes.Buffer
	_ = reportDoctorRows(&stdout)
	// The match is the row's own name, not the word "worktrees" anywhere in the output. An
	// unrelated row that prints an absolute path inside a Bench worktree pool otherwise
	// reads as this row firing.
	if got := stdout.String(); strings.Contains(got, "worktree admin") {
		t.Fatalf("healthy doctor rows include worktree-admin row: %q", got)
	}
}

func runAdoptGit(t *testing.T, root string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func TestDoctorDirSelectionAndShimRoundTrip(t *testing.T) {
	home := t.TempDir()
	nvm := filepath.Join(home, ".nvm")
	manager := filepath.Join(nvm, "versions", "node", "v22", "bin")
	blocked := filepath.Join(home, "blocked")
	plain := filepath.Join(home, "plain bin")
	if err := os.MkdirAll(manager, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(blocked, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(blocked, 0o755)
	if err := os.MkdirAll(plain, 0o755); err != nil {
		t.Fatal(err)
	}
	env := DoctorEnv{
		Home:   home,
		Path:   strings.Join([]string{manager, blocked, plain}, string(os.PathListSeparator)),
		NVMDir: nvm,
	}
	chosen, create := SelectDoctorDir(env)
	if create {
		t.Fatalf("SelectDoctorDir create = true, want false")
	}
	if chosen != plain {
		t.Fatalf("SelectDoctorDir = %q, want %q", chosen, plain)
	}

	target := filepath.Join(home, "kit path [x]", "bin", "bench.sh")
	content := ShimContent(target)
	if !strings.Contains(content, "# bench-target: "+target) {
		t.Fatalf("shim missing literal target comment:\n%s", content)
	}
	if !strings.Contains(content, "target=") || !strings.Contains(content, `exec "$target" "$@"`) {
		t.Fatalf("shim missing executable target assignment:\n%s", content)
	}
	if got := ShimTarget(content); got != target {
		t.Fatalf("ShimTarget = %q, want %q", got, target)
	}

	env.Path = manager
	fallback, create := SelectDoctorDir(env)
	if fallback != filepath.Join(home, ".local", "bin") || !create {
		t.Fatalf("fallback = %q create=%v, want ~/.local/bin create", fallback, create)
	}
}

func TestDoctorParsesHelpAndUnknownFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Doctor([]string{"--help"}, &stdout, &stderr, "1.0.0"); code != 0 {
		t.Fatalf("Doctor --help = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if got, want := stdout.String(), "usage: bench doctor [--fix]\n"; got != want {
		t.Fatalf("Doctor --help stdout = %q, want %q", got, want)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("Doctor --help stderr = %q, want empty", got)
	}

	stdout.Reset()
	stderr.Reset()
	if code := Doctor([]string{"--bogus"}, &stdout, &stderr, "1.0.0"); code != 2 {
		t.Fatalf("Doctor --bogus = %d, want 2\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("Doctor --bogus stdout = %q, want empty", got)
	}
	if got, want := stderr.String(), "usage: bench doctor (unknown argument: --bogus)\n"; got != want {
		t.Fatalf("Doctor --bogus stderr = %q, want %q", got, want)
	}
}

// guardedCanaryCall is the exact guarded inventory call scaffoldGate must emit and
// bench init must write: rooted at $root, never relative.
const guardedCanaryCall = `if [ -d "$root/tests/canary" ]; then
  "$bench" canary "$root" || err "canary inventory validation failed"
fi`

func TestScaffoldGateUsesCanarySubcommand(t *testing.T) {
	gate := scaffoldGate()
	mustContain := []string{
		"BENCH_SENTINEL",
		"bench=\"$(dirname \"$0\")/bin/bench.sh\"; [ -x \"$bench\" ] || bench=bench",
		guardedCanaryCall,
	}
	for _, want := range mustContain {
		if !strings.Contains(gate, want) {
			t.Fatalf("scaffold gate missing %q:\n%s", want, gate)
		}
	}
	for _, forbidden := range []string{". \"$gate_dir/lib/canary-run.sh\"", "canary runner missing", "BENCH_CANARY_INNER"} {
		if strings.Contains(gate, forbidden) {
			t.Fatalf("scaffold gate still contains retired sourcing API %q:\n%s", forbidden, gate)
		}
	}
}

func TestInitWritesGuardedCanaryCall(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(wd, "..", "..")))
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	if code := Init(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("Init exit = %d, stderr:\n%s", code, stderr.String())
	}
	gate := readFile(t, filepath.Join(root, ".bench", "gate.sh"))
	if !strings.Contains(gate, guardedCanaryCall) {
		t.Fatalf("bench init gate.sh missing guarded canary call:\n%s", gate)
	}
	if !strings.Contains(gate, SentinelMarker) {
		t.Fatalf("bench init gate.sh missing %s sentinel:\n%s", SentinelMarker, gate)
	}
}

func TestInitDoesNotSeedLinkedProof(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("git", "-C", root, "init", "-q")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(wd, "..", "..")))
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	if code := Init(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("Init exit = %d, stderr:\n%s", code, stderr.String())
	}
	for _, path := range []string{"tests/canary/example/example/EXPECT", "tests/canary/example/example/files/DO-NOT-SHIP"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Fatalf("retired seed path %q exists or stat failed unexpectedly: %v", path, err)
		}
	}
	if strings.Contains(stdout.String(), "proved") || strings.Contains(stdout.String(), "seed canary") {
		t.Fatalf("Init output makes a planted-proof claim: %s", stdout.String())
	}
}

func TestInitReentryPreservesProjectInventoryInSpecialPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repo [x].*")
	if err := os.MkdirAll(filepath.Join(root, "tests", "canary", "project", "family", "files"), 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "canary", "project", "family", "EXPECT"), []byte("project check\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "canary", "project", "family", "files", "owned"), []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	kit := filepath.Clean(filepath.Join(mustGetwd(t), "..", ".."))
	t.Setenv("BENCH_KIT", kit)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if code := Init(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("Init run 1 exit = %d, stderr: %s", code, stderr.String())
	}
	projectGate := "#!/usr/bin/env bash\nproject-owned-check\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte(projectGate), 0o755); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := Init(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("Init run 2 exit = %d, stderr: %s", code, stderr.String())
	}
	if got := readFile(t, filepath.Join(root, ".bench", "gate.sh")); got != projectGate {
		t.Fatalf("project gate = %q, want preserved project-owned check", got)
	}
	if got := readFile(t, filepath.Join(root, "tests", "canary", "project", "family", "EXPECT")); got != "project check\n" {
		t.Fatalf("project EXPECT = %q", got)
	}
	if got := readFile(t, filepath.Join(root, "tests", "canary", "project", "family", "files", "owned")); got != "keep\n" {
		t.Fatalf("project fixture = %q", got)
	}
	if _, err := os.Stat(filepath.Join(root, "tests", "canary", "example")); !os.IsNotExist(err) {
		t.Fatalf("retired seed directory exists or stat failed unexpectedly: %v", err)
	}
}

func mustGetwd(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return wd
}

// TestCompareKitVersions pins the ordering `bench upgrade` decides on. The prerelease
// rows are the point. A repo pinned at a prerelease of the installed release must read as
// behind it, not as equal. Otherwise the applying run returns success without relinking,
// and the manifest stays stamped at the prerelease forever.
func TestCompareKitVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.4", "1.2.3", 1},
		{"1.2.3", "1.2.4", -1},
		{"1.3.0", "1.2.9", 1},
		{"2.0.0", "1.9.9", 1},
		{"1.2.3", "1.2.3-rc1", 1},
		{"1.2.3-rc1", "1.2.3", -1},
		{"1.2.3-rc2", "1.2.3-rc1", 1},
		{"1.2.3-rc1", "1.2.3-rc2", -1},
		{"1.2.4-rc1", "1.2.3", 1},
		{"1.2.3+build9", "1.2.3+build1", 0},
		{"1.2.3-rc1+build9", "1.2.3-rc1", 0},
		// An unparseable stamp on either side never reports a downgrade, so a
		// hand-edited header cannot make upgrade refuse an install it should do.
		{"dev", "1.2.3", 1},
		{"1.2.3", "dev", 1},
	}
	for _, tc := range cases {
		if got := compareKitVersions(tc.a, tc.b); got != tc.want {
			t.Errorf("compareKitVersions(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	return string(data)
}

// TestScaffoldLearningsMirrorsTheParserBoundary proves the scaffold and the parser stay
// one source. The bytes a fresh repo receives parse with zero malformed records. The
// scaffold's boundary line is the parser's own exported marker rather than a second
// literal. The journal fixture internal/learnings checks in is a mirror of this scaffold
// rather than a hand-copied second source of the same bytes.
func TestScaffoldLearningsMirrorsTheParserBoundary(t *testing.T) {
	got := scaffoldLearnings()
	if _, malformed := learnings.Parse([]byte(got)); len(malformed) != 0 {
		t.Fatalf("scaffolded journal parses with %d malformed records, want 0: %#v", len(malformed), malformed)
	}
	if !strings.Contains(got, "\n"+learnings.JournalEntriesMarker+"\n") {
		t.Fatalf("scaffold does not carry learnings.JournalEntriesMarker on its own line: %q", got)
	}
	mirror, err := os.ReadFile(filepath.Join("..", "learnings", "testdata", "scaffold-learnings.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(mirror) != got {
		t.Fatalf("internal/learnings/testdata/scaffold-learnings.md has drifted from scaffoldLearnings()\nfixture = %q\nscaffold = %q", mirror, got)
	}
}

// doctorGreenRepo builds a git repo whose doctor rows and shim row are all green, with a
// bench-managed pre-push hook in place. The caller sets BENCH_KIT to decide whether the
// repo is the kit source checkout or a consumer repo. BENCH_WRAPPER keeps doctor --fix's
// shim and broker-manifest writes inside the fixture.
func doctorGreenRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	for _, dir := range [][]string{{".bench"}, {"projects"}, {".agents", "commands"}} {
		if err := os.MkdirAll(filepath.Join(append([]string{root}, dir...)...), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFixtureFile(t, filepath.Join(root, ".bench", "gate.sh"), "#!/bin/sh\nexit 0\n", 0o755)
	writeFixtureFile(t, filepath.Join(root, "CLAUDE.md"), strings.Join(claudeImportLines(), "\n")+"\n", 0o644)
	writeFixtureFile(t, filepath.Join(root, "projects", "kit.md"), "# kit\n", 0o644)
	writeFixtureFile(t, filepath.Join(root, ".agents", "commands", "bench-setup-repo.md"), "# setup\n", 0o644)
	writeHook(t, filepath.Join(root, ".git", "hooks", "pre-push"), fallbackProtectedBranch)

	bin := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(bin, "bench-wrapper")
	writeFixtureFile(t, target, "#!/bin/sh\n", 0o755)
	writeFixtureFile(t, filepath.Join(bin, "bench"), ShimContent(target)+"\n", 0o755)
	t.Setenv("BENCH_WRAPPER", target)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Chdir(root)
	return root
}

func writeFixtureFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

// TestDoctorKitSourceCheckoutRowsAreGreen pins the kit source checkout's two ok rows. The
// kit repo is the source of the managed block and of bin/bench.sh, so neither asset can
// exist in the shape a consumer repo carries, and doctor must say so instead of sending
// the reader to bench link.
func TestDoctorKitSourceCheckoutRowsAreGreen(t *testing.T) {
	root := doctorGreenRepo(t)
	t.Setenv("BENCH_KIT", root)

	var stdout, stderr bytes.Buffer
	if code := Doctor(nil, &stdout, &stderr, "1.0.0"); code != 0 {
		t.Fatalf("Doctor = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	got := stdout.String()
	for _, want := range []string{
		"ok: kit source checkout - AGENTS.md is the source agreement; no managed block applies",
		"ok: kit source checkout - the launcher is bin/bench.sh; no .bench/bin copy applies",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("doctor output missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "run bench link") {
		t.Fatalf("kit source checkout doctor output names bench link:\n%s", got)
	}
}

// TestDoctorKitSourceCheckoutAbsentPrePushRoutesToFix pins the one remedy the kit checkout
// can act on. bench link is not that remedy here, so the absent-hook row names the fix
// that installs the managed hook, and the fix installs it.
func TestDoctorKitSourceCheckoutAbsentPrePushRoutesToFix(t *testing.T) {
	root := doctorGreenRepo(t)
	t.Setenv("BENCH_KIT", root)
	hook := filepath.Join(root, ".git", "hooks", "pre-push")
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := Doctor(nil, &stdout, &stderr, "1.0.0"); code != 1 {
		t.Fatalf("Doctor with absent hook = %d, want 1\n%s", code, stdout.String())
	}
	if got := stdout.String(); !strings.Contains(got, "hooks); run bench doctor --fix") {
		t.Fatalf("absent pre-push row does not name bench doctor --fix:\n%s", got)
	}

	stdout.Reset()
	stderr.Reset()
	if code := Doctor([]string{"--fix"}, &stdout, &stderr, "1.0.0"); code != 0 {
		t.Fatalf("Doctor --fix = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if health := InspectPrePush(root); health.State != PrePushManaged || health.Currency != PrePushCurrent {
		t.Fatalf("pre-push after --fix = %#v, want managed current", health)
	}

	stdout.Reset()
	stderr.Reset()
	if code := Doctor(nil, &stdout, &stderr, "1.0.0"); code != 0 {
		t.Fatalf("Doctor after --fix = %d, want 0\n%s", code, stdout.String())
	}
}

// TestDoctorConsumerRepoKeepsLinkRemedies is the regression side: a repo whose kit lives
// elsewhere keeps every current red row and keeps bench link as the remedy, and --fix
// still leaves an absent hook alone there.
func TestDoctorConsumerRepoKeepsLinkRemedies(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), "# project\n", 0o644)
	t.Setenv("BENCH_KIT", t.TempDir())
	t.Chdir(root)

	var stdout bytes.Buffer
	if !reportDoctorRows(&stdout) {
		t.Fatal("consumer doctor rows report no red")
	}
	got := stdout.String()
	for _, want := range []string{
		"red: AGENTS.md has no Bench managed block (run bench link)",
		"red: .bench/bin/bench.sh absent (run bench link)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("consumer doctor rows missing %q:\n%s", want, got)
		}
	}

	stdout.Reset()
	if !reportPrePush(&stdout) {
		t.Fatal("consumer absent pre-push row is not red")
	}
	if row := stdout.String(); !strings.Contains(row, "hooks); run bench link") {
		t.Fatalf("consumer absent pre-push row = %q, want it to name bench link", row)
	}

	var fixOut, fixErr bytes.Buffer
	if code := repairStalePrePush(&fixOut, &fixErr); code != 0 {
		t.Fatalf("repairStalePrePush = %d, want 0\n%s", code, fixErr.String())
	}
	if _, err := os.Lstat(filepath.Join(root, ".git", "hooks", "pre-push")); !os.IsNotExist(err) {
		t.Fatalf("consumer --fix installed a pre-push hook: %v", err)
	}
}
