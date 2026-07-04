package contract

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestAXIContracts(t *testing.T) {
	skipIfSubjectBenchMissing(t)
	t.Run("block-dangerous-git analyzer-missing fail-closed contract", testAXIBlockDangerousGitAnalyzerMissing)
	t.Run("block-dangerous-git binary-missing fail-closed (git-shaped)", testAXIBlockDangerousGitBinaryMissing)
	t.Run("block-dangerous-git core-errored fail-closed contract", testAXIBlockDangerousGitCoreErrored)
	t.Run("AXI learnings two-entry contract", testAXILearningsTwoEntry)
	t.Run("AXI learnings empty/template contract", testAXILearningsEmptyTemplate)
	t.Run("AXI maps unresolved-ticket contract", testAXIMapsUnresolvedTicket)
	t.Run("AXI TOON field-escaping contract", testAXITOONFieldEscaping)
	t.Run("AXI usage/exit-2 contract", testAXIUsageExit2)
	t.Run("AXI subdirectory root-resolution contract", testAXISubdirectoryRootResolution)
	t.Run("AXI path-with-spaces contract", testAXIPathWithSpaces)
	t.Run("AXI maps over-match anchoring/fence contract", testAXIMapsOverMatchAnchoringFence)
	t.Run("AXI maps CRLF-stripping contract", testAXIMapsCRLFStripping)
	t.Run("AXI maps no-Type-ticket contract", testAXIMapsNoTypeTicket)
	t.Run("AXI learnings ascii-separator title contract", testAXILearningsASCIISeparatorTitle)
	t.Run("AXI maps handoff close-readiness contract", testAXIMapsHandoffCloseReadiness)
	t.Run("AXI maps --count adapter contract", testAXIMapsCountAdapter)
}

func testAXIBlockDangerousGitAnalyzerMissing(t *testing.T) {
	f := NewFixture(t)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	pathEnv := map[string]string{"PATH": "/usr/bin:/bin"}

	describe := runFixtureCommand(t, f, f.Root, pathEnv, "", bashPath(t), hook, "--describe")
	describe.RequireExit(0)
	requireAXILine(t, describe.Stdout, "denies: manifest unavailable (analyzer missing)")

	enforce := runFixtureCommand(t, f, f.Root, pathEnv, `{"tool_input":{"command":"git push"}}`, bashPath(t), hook)
	enforce.RequireExit(2)
	enforce.RequireContains(enforce.Stderr, "BLOCKED")
}

func testAXIBlockDangerousGitBinaryMissing(t *testing.T) {
	f := NewFixture(t)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	stubBin := filepath.Join(f.Root, "stubbin")
	if err := os.MkdirAll(stubBin, 0o755); err != nil {
		t.Fatalf("create stubbin: %v", err)
	}
	writeExecutableAt(t, stubBin, "bench", "#!/usr/bin/env bash\nexit 127\n")
	pathEnv := map[string]string{"PATH": stubBin + ":/usr/bin:/bin"}

	gitShaped := runFixtureCommand(t, f, f.Root, pathEnv, `{"tool_input":{"command":"git push"}}`, bashPath(t), hook)
	gitShaped.RequireExit(2)
	gitShaped.RequireContains(gitShaped.Stderr, "BLOCKED")

	nonGit := runFixtureCommand(t, f, f.Root, pathEnv, `{"tool_input":{"command":"ls -la"}}`, bashPath(t), hook)
	nonGit.RequireExit(0)
}

func testAXIBlockDangerousGitCoreErrored(t *testing.T) {
	f := NewFixture(t)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	stubBin := filepath.Join(f.Root, "stubbin")
	if err := os.MkdirAll(stubBin, 0o755); err != nil {
		t.Fatalf("create stubbin: %v", err)
	}
	writeExecutableAt(t, stubBin, "bench", "#!/usr/bin/env bash\nexit 3\n")

	enforce := runFixtureCommand(t, f, f.Root, map[string]string{"PATH": stubBin + ":/usr/bin:/bin"}, `{"tool_input":{"command":"git push"}}`, bashPath(t), hook)
	enforce.RequireExit(2)
	enforce.RequireContains(enforce.Stderr, "BLOCKED")
}

func testAXILearningsTwoEntry(t *testing.T) {
	noteContractFailure(t, "AXI learnings two-entry contract failed")
	f := NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — first learning  [open]\n- body\n## 2026-02-02 — second learning  [open]\n- body\n")

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "learnings[2]{date,title}:")
	requireIndentedRowCount(t, out.Stdout, 2)
	requireAXILine(t, out.Stdout, "  2026-01-01,first learning")
	requireAXILine(t, out.Stdout, "  2026-02-02,second learning")
}

func testAXILearningsEmptyTemplate(t *testing.T) {
	f := NewFixture(t)

	absent := f.Bench("learnings")
	absent.RequireExit(0)
	requireAXILine(t, absent.Stdout, "learnings[0]{date,title}:")

	f.WriteFile(".bench/learnings.md", "## <date> — <short title>  [open]\n")
	template := f.Bench("learnings")
	template.RequireExit(0)
	requireAXILine(t, template.Stdout, "learnings[0]{date,title}:")
}

func testAXIMapsUnresolvedTicket(t *testing.T) {
	f := NewFixture(t)

	absent := f.Bench("maps")
	absent.RequireExit(0)
	requireAXILine(t, absent.Stdout, "maps[0]{map,ticket,type,state}:")

	f.WriteFile("decisions/m.md", `# A map

## #1: first?

Type: Grill

### Answer
— (open)

## #2: second?

Type: Research

### Answer
— (deferred)

## #3: third?

Type: Grill

### Answer
GRILL DEFERRED — waiting on X

## #4: resolved?

Type: Grill

### Answer
Decided: yes, do the thing.
`)
	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "maps[3]{map,ticket,type,state}:")
	requireAXILine(t, out.Stdout, "  m,1,Grill,open")
	requireAXILine(t, out.Stdout, "  m,2,Research,deferred")
	requireAXILine(t, out.Stdout, "  m,3,Grill,grill-deferred")
	requireNoAXILineMatching(t, out.Stdout, `^  m,4,`)
}

func testAXITOONFieldEscaping(t *testing.T) {
	noteContractFailure(t, "AXI TOON field-escaping contract failed")
	f := NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-03-03 — a, \"b\"  [open]\n")

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, `  2026-03-03,"a, \"b\""`)
}

func testAXIUsageExit2(t *testing.T) {
	f := NewFixture(t)

	learnings := f.Bench("learnings", "bogusarg")
	learnings.RequireExit(2)
	requireContainsFold(t, learnings.Stdout, "usage")

	maps := f.Bench("maps", "bogusarg")
	maps.RequireExit(2)
	requireContainsFold(t, maps.Stdout, "usage")
}

func testAXISubdirectoryRootResolution(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-04-04 — sub check  [open]\n")
	subdir := filepath.Join(f.Root, "sub", "deeper")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("create subdir: %v", err)
	}

	out := runBenchInDir(t, f, subdir, "learnings")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "learnings[1]{date,title}:")
	requireAXILine(t, out.Stdout, "  2026-04-04,sub check")
}

func testAXIPathWithSpaces(t *testing.T) {
	f := NewFixture(t, WithSpacePath())
	f.WriteFile(".bench/learnings.md", "## 2026-05-05 — spaced  [open]\n")
	f.WriteFile("decisions/s.md", `## #1: q?

Type: Grill

### Answer
— (open)
`)

	learnings := f.Bench("learnings")
	learnings.RequireExit(0)
	requireAXILine(t, learnings.Stdout, "  2026-05-05,spaced")

	maps := f.Bench("maps")
	maps.RequireExit(0)
	requireAXILine(t, maps.Stdout, "  s,1,Grill,open")
}

func testAXIMapsOverMatchAnchoringFence(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/o.md", `# Over-match map

## #1: genuine?

Type: Grill

### Answer
— (open)

## #2: prose mentions?

Type: Grill

### Answer
Decided: a mid-line GRILL DEFERRED mention is not an unresolved banner.

## #3: fenced example?

Type: Grill

### Answer
Decided: the placeholder looks like this:

`+"```"+`
— (open)
`+"```"+`

so authors recognize it.
`)

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "maps[1]{map,ticket,type,state}:")
	requireAXILine(t, out.Stdout, "  o,1,Grill,open")
	requireNoAXILineMatching(t, out.Stdout, `^  o,2,`)
	requireNoAXILineMatching(t, out.Stdout, `^  o,3,`)
}

func testAXIMapsCRLFStripping(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/c.md", "## #1: q?\r\nType: Grill\r\n### Answer\r\n— (open)\r\n")

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  c,1,Grill,open")
	if strings.Contains(out.Stdout, "\r") {
		t.Fatalf("CRLF leaked carriage returns into maps output:\n%s", out.Stdout)
	}
}

func testAXIMapsNoTypeTicket(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/n.md", `## #1: typeless?

### Answer
— (open)
`)

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  n,1,unknown,open")
}

func testAXILearningsASCIISeparatorTitle(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 - ascii title  [open]\n")

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  2026-01-01,ascii title")
}

func testAXIMapsHandoffCloseReadiness(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/hm.md", `# HM
## #1: q?
Type: Grill
### Answer
Decided: yes.
`)
	f.WriteFile("decisions/hf.md", `# HF
## #1: q?
Type: Grill
### Answer
Decided: yes.

## Handoff
1. Module boundaries. n/a — single unit.
2. Contracts. n/a — no CLI surface.
`)
	f.WriteFile("decisions/ho.md", `# HO
## #1: q?
Type: Grill
### Answer
— (open)
`)
	f.WriteFile("decisions/hx.md", `# HX
## #1: q?
Type: Grill
### Answer
Decided: yes. A handoff section looks like:

`+"```"+`
## Handoff
1. Module boundaries.
`+"```"+`
`)
	f.WriteFile("decisions/hp.md", `# HP
## #1: q?
Type: Grill
### Answer
Decided: yes.

## Handoff
1. Module boundaries.
— (open)
`)
	f.WriteFile("decisions/README.md", "# Decisions index\nNotes about this folder — not a map.\n")

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  hm,handoff,handoff,missing")
	requireAXILine(t, out.Stdout, "  hx,handoff,handoff,missing")
	requireAXILine(t, out.Stdout, "  hp,handoff,handoff,open")
	requireNoAXILineMatching(t, out.Stdout, `^  hf,handoff`)
	requireNoAXILineMatching(t, out.Stdout, `^  ho,handoff`)
	requireAXILine(t, out.Stdout, "  ho,1,Grill,open")
	requireNoAXILineMatching(t, out.Stdout, `^  hp,1,`)
	requireNoAXILineMatching(t, out.Stdout, `^  README,`)
}

func testAXIMapsCountAdapter(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/m.md", "## #1: a?\nType: Grill\n### Answer\n— (open)\n")
	f.WriteFile("decisions/scope.md", "### Answer\n— (deferred)\n")
	f.WriteFile("decisions/.hidden.md", "## #1: h?\nType: Grill\n### Answer\n— (open)\n")

	count := f.Bench("maps", "--count")
	count.RequireExit(0)
	if got := strings.TrimSpace(count.Stdout); got != "2" {
		t.Fatalf("maps --count = %q, want 2\nstdout:\n%s\nstderr:\n%s", got, count.Stdout, count.Stderr)
	}

	rows := f.Bench("maps")
	rows.RequireExit(0)
	requireAXILine(t, rows.Stdout, "  m,1,Grill,open")
	requireNoAXILineMatching(t, rows.Stdout, `^  scope,`)
	requireNoAXILineMatching(t, rows.Stdout, `^  \.?hidden,`)
}

func runBenchInDir(t testing.TB, f Fixture, dir string, args ...string) Probe {
	t.Helper()
	bench := filepath.Join(SubjectRoot(t), "bin", "bench.sh")
	return runFixtureCommand(t, f, dir, nil, "", "bash", append([]string{bench}, args...)...)
}

func runFixtureCommand(t testing.TB, f Fixture, dir string, env map[string]string, stdin, name string, args ...string) Probe {
	t.Helper()
	typed := make(Env, len(env))
	for k, v := range env {
		value := v
		typed[k] = &value
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnv(f.Env, typed)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}
	return Probe{t: t, ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String()}
}

func bashPath(t testing.TB) string {
	t.Helper()
	bin, err := exec.LookPath("bash")
	if err != nil {
		t.Fatalf("find bash: %v", err)
	}
	return bin
}

func writeExecutableAt(t testing.TB, dir, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o755); err != nil {
		t.Fatalf("write executable %s: %v", name, err)
	}
}

func requireAXIFirstLine(t testing.TB, out, want string) {
	t.Helper()
	if got := axiFirstLine(out); got != want {
		t.Fatalf("first line = %q, want %q\noutput:\n%s", got, want, out)
	}
}

func axiFirstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func requireAXILine(t testing.TB, out, want string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if line == want {
			return
		}
	}
	t.Fatalf("missing line %q\noutput:\n%s", want, out)
}

func requireIndentedRowCount(t testing.TB, out string, want int) {
	t.Helper()
	got := 0
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if strings.HasPrefix(line, "  ") {
			got++
		}
	}
	if got != want {
		t.Fatalf("indented row count = %d, want %d\noutput:\n%s", got, want, out)
	}
}

func requireNoAXILineMatching(t testing.TB, out, expr string) {
	t.Helper()
	re := regexp.MustCompile(expr)
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if re.MatchString(line) {
			t.Fatalf("unexpected line matching %q: %q\noutput:\n%s", expr, line, out)
		}
	}
}

func requireContainsFold(t testing.TB, haystack, needle string) {
	t.Helper()
	if !strings.Contains(strings.ToLower(haystack), strings.ToLower(needle)) {
		t.Fatalf("missing %q\noutput:\n%s", needle, haystack)
	}
}
