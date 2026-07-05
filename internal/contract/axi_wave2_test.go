package contract

import (
	"strings"
	"testing"
)

func TestAXIWave2Contracts(t *testing.T) {
	t.Parallel()
	skipIfSubjectBenchMissing(t)
	runParallel(t, "AXI diff recorded-base contract", testAXIDiffRecordedBase)
	runParallel(t, "AXI diff fallback/shape contract", testAXIDiffFallbackShape)
	runParallel(t, "AXI diff error-posture contract", testAXIDiffErrorPosture)
	runParallel(t, "AXI coverage extraction contract", testAXICoverageExtraction)
	runParallel(t, "AXI coverage state/error contract", testAXICoverageStateError)
	runParallel(t, "AXI coverage --check validation contract", testAXICoverageCheckValidation)
}

func testAXIDiffRecordedBase(t *testing.T) {
	noteContractFailure(t, "AXI diff recorded-base contract failed")
	f := NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	f.Git("switch", "-qc", "feature")
	f.WriteFile("f.txt", "f\n")
	f.CommitAll("c2")
	c2 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-qc", "bench/shift-stacked")
	f.Git("config", "branch.bench/shift-stacked.benchBase", c2)
	f.WriteFile("work.txt", "w\n")
	f.CommitAll("c3")

	out := f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c2)
	requireOutputLine(t, out, "method: recorded")
	requireOutputLine(t, out, "files[1]{status,path}:")
	requireOutputLine(t, out, "  A,work.txt")

	f.Git("switch", "-q", "main")
	f.Git("switch", "-qc", "other")
	f.WriteFile("o.txt", "o\n")
	f.CommitAll("o1")
	o1 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-q", "bench/shift-stacked")
	f.Git("config", "branch.bench/shift-stacked.benchBase", o1)
	mergeBase := strings.TrimSpace(f.Git("merge-base", "main", "HEAD").Stdout)

	out = f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "method: merge-base (recorded sha not an ancestor)")
	requireOutputLine(t, out, "base: "+mergeBase)
}

func testAXIDiffFallbackShape(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "main")
	c1 := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.Git("switch", "-qc", "feature")
	f.WriteFile("sub/deeper/.keep", "")
	f.WriteFile("a b.txt", "f\n")
	f.WriteFile("café.txt", "u\n")
	f.WriteFile("a\"q.txt", "q\n")
	f.CommitAll("c2")

	out := f.Run("bash", "-c", "cd sub/deeper && "+shellQuote(SubjectRoot(t)+"/bin/bench.sh")+" diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c1)
	requireOutputLine(t, out, "method: merge-base")
	requireOutputLine(t, out, "  A,a b.txt")
	requireOutputLine(t, out, "  A,café.txt")
	requireOutputLine(t, out, `  A,"a\"q.txt"`)

	f.Git("config", "branch.feature.benchBase", "0123456789abcdef0123456789abcdef01234567")
	out = f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "base: "+c1)
	requireOutputLine(t, out, "method: merge-base (recorded sha unreachable)")

	f.Git("switch", "-q", "main")
	f.Git("switch", "-qc", "idle")
	out = f.Bench("diff")

	out.RequireExit(0)
	requireOutputLine(t, out, "files[0]{status,path}:")
}

func testAXIDiffErrorPosture(t *testing.T) {
	f := NewFixture(t, WithNoRepo())

	out := f.Bench("diff")

	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")

	f.Git("init", "-q")
	f.WriteFile("README.md", "r\n")
	f.CommitAll("c1")
	f.Git("branch", "-m", "trunk")

	out = f.Bench("diff")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "error: cannot resolve a review base")

	out = f.Bench("diff", "bogusarg")

	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stdout), "usage")
}

func testAXICoverageExtraction(t *testing.T) {
	noteContractFailure(t, "AXI coverage extraction contract failed")
	f := NewFixture(t)
	spec := `# t

## User stories
1. As a, I want b, so c.
2. As d, I want e, so f.
3. As g, I want h, so i.

## Testing decisions

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2–3 | does x \| y | cli seam | cmd fails, loudly | catches z |
| edge of 1 | edge case | gate | already covered | catches w |

## Out of scope
`
	f.WriteFile("specs/t.md", spec)

	out := f.Bench("coverage", "specs/t.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "spec: specs/t.md")
	requireOutputLine(t, out, "state: mapped")
	requireOutputLine(t, out, "rows[2]{story,seam,red_signal}:")
	requireOutputLine(t, out, `  2–3,cli seam,"cmd fails, loudly"`)
	requireOutputLine(t, out, "  edge of 1,gate,already covered")

	f.WriteFile("specs/t.md", strings.ReplaceAll(spec, "\n", "\r\n"))
	out = f.Bench("coverage", "specs/t.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: mapped")
	requireOutputLine(t, out, "rows[2]{story,seam,red_signal}:")
}

func testAXICoverageStateError(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("specs/h.md", "# h\n\n<!-- coverage-map: historical -->\n\n### Acceptance coverage map\n| a |\n")

	out := f.Bench("coverage", "specs/h.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: historical")
	requireOutputLine(t, out, "rows[0]{story,seam,red_signal}:")

	f.WriteFile("specs/n.md", "# n\nno map here\n")
	out = f.Bench("coverage", "specs/n.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: no-map")

	out = f.Bench("coverage")

	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stdout), "usage")

	out = f.Bench("coverage", "--bogus", "specs/n.md")

	out.RequireExit(2)

	out = f.Bench("coverage", "specs/absent.md")

	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")
}

func testAXICoverageCheckValidation(t *testing.T) {
	f := NewFixture(t)
	stories := "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n"
	header := "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|"
	f.WriteFile("specs/v.md", "# v\n\n"+stories+"\n### Acceptance coverage map\n"+header+"\n| 1, 2–3 | b | s | r | w |\n")

	out := f.Bench("coverage", "--check", "specs/v.md")

	out.RequireExit(0)
	requireNoOutput(t, out)

	f.WriteFile("specs/h.md", "# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")
	out = f.Bench("coverage", "--check", "specs/h.md")
	out.RequireExit(0)
	requireNoOutput(t, out)

	f.WriteFile("specs/n.md", "# n\nprose only\n")
	out = f.Bench("coverage", "--check", "specs/n.md")
	out.RequireExit(0)
	requireNoOutput(t, out)

	cases := []struct {
		name string
		path string
		body string
		want string
	}{
		{
			name: "missing canonical header",
			path: "specs/b1.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n",
			want: "coverage map missing the canonical header",
		},
		{
			name: "no data rows",
			path: "specs/b2.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n",
			want: "coverage map has no data rows",
		},
		{
			name: "wrong cell count",
			path: "specs/b3.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 1 | b | s | r |\n",
			want: "coverage map row 1 has 4 cells (want 5)",
		},
		{
			name: "empty cell",
			path: "specs/b4.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 1 | b |  | r | w |\n",
			want: "coverage map row 1 has an empty 'seam' cell",
		},
		{
			name: "story reference beyond numbering",
			path: "specs/b5.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 9 | b | s | r | w |\n",
			want: "references story 9 but the spec numbers only 3",
		},
		{
			name: "unrecognized story reference",
			path: "specs/b6.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| x | b | s | r | w |\n",
			want: "has an unrecognized story reference 'x'",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f := NewFixture(t)
			f.WriteFile(c.path, c.body)

			out := f.Bench("coverage", "--check", c.path)

			out.RequireExit(1)
			out.RequireContains(out.Stdout, c.want)
		})
	}
}

func requireOutputLine(t *testing.T, probe Probe, line string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if got == line {
			return
		}
	}
	t.Fatalf("missing output line %q\nstdout:\n%s\nstderr:\n%s", line, probe.Stdout, probe.Stderr)
}

func requireOutputPrefix(t *testing.T, probe Probe, prefix string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if strings.HasPrefix(got, prefix) {
			return
		}
	}
	t.Fatalf("missing output line with prefix %q\nstdout:\n%s\nstderr:\n%s", prefix, probe.Stdout, probe.Stderr)
}

func requireNoOutput(t *testing.T, probe Probe) {
	t.Helper()
	if probe.Stdout != "" || probe.Stderr != "" {
		t.Fatalf("expected no output\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
