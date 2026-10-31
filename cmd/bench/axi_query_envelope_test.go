package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

type axiEnvelopeCase struct {
	successMarker, emptyMarker, usage              string
	blocks                                         []string
	route, successArgv, deepSuccessArgv, emptyArgv []string
	setupSuccess, setupEmpty                       func(*testing.T, string)
	// recordBacked marks a query that projects a record compiled into the binary. Such a
	// query reads no disk and needs no repository. So it has no empty projection and no
	// path it can refuse, and the definitive-empty and structured-refusal subtests do not
	// apply to it. Every other subtest still runs. A query that reads the repository must
	// leave this field false.
	recordBacked bool
}

// resultBlock names the table whose rows differ between the success and empty cases.
// The name comes from the identifying marker, so the two stay in step.
func (tc axiEnvelopeCase) resultBlock(t *testing.T) string {
	t.Helper()
	name, _, found := strings.Cut(tc.successMarker, "[")
	if !found || name == "" {
		t.Fatalf("success marker %q names no result table", tc.successMarker)
	}
	return name
}

// axiEnvelopeRows decodes a member's whole stdout as one TOON document and returns its
// result rows. The decode fails on any byte before, between, or after the blocks. The
// recovered block order pins the order the command emitted, and the help table must be
// schema-correct and terminal.
func axiEnvelopeRows(t *testing.T, tc axiEnvelopeCase, result axiCommandResult) []any {
	t.Helper()
	document, err := axitest.DecodeDocument(result.stdout)
	if err != nil {
		t.Fatalf("stdout = %q, want structured TOON: %v", result.stdout, err)
	}
	if !reflect.DeepEqual(document.Blocks, tc.blocks) {
		t.Fatalf("stdout blocks = %q, want %q; stdout=%q", document.Blocks, tc.blocks, result.stdout)
	}
	if _, err := document.HelpActions(); err != nil {
		t.Fatalf("stdout = %q, want a terminal help[N]{cmd,why} envelope: %v", result.stdout, err)
	}
	rows, err := document.Rows(tc.resultBlock(t))
	if err != nil {
		t.Fatalf("stdout = %q: %v", result.stdout, err)
	}
	return rows
}

type axiCommandResult struct {
	stdout, stderr string
	code           int
}

func axiEnvelopeCaseNames(cases map[string]axiEnvelopeCase) []string {
	names := make([]string, 0, len(cases))
	for name := range cases {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func newAXIEnvelopeRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	writeAXIFixture(t, filepath.Join(root, "nested", "deep", ".keep"), "fixture\n")
	gittest.MustRun(t, root, "add", ".")
	gittest.Commit(t, root, "fixture base")
	base := gittest.Output(t, root, "rev-parse", "HEAD")
	gittest.MustRun(t, root, "config", "branch.main.benchBase", base)
	return root
}

func setupAXIAnchors(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), readAnchorsFixture(t, "AGENTS.md"))
}

func setupAXILearnings(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "capture", "learnings.md"), "# Learnings — usage journal\n\n## 2026-01-01 — fixture [open]\n")
}

func setupAXIMap(t *testing.T, root string) {
	index := strings.Replace(maps.DecisionMapTemplate(), "\n- [<decision question>](<topic>/tickets/1.md): <gist>\n", "", 1)
	ticket := strings.Replace(maps.DecisionTicketTemplate(), "<answer>", "— (open)", 1)
	writeAXIFixture(t, filepath.Join(root, "decisions", "fixture.md"), index)
	writeAXIFixture(t, filepath.Join(root, "decisions", "fixture", "tickets", "1.md"), ticket)
}

func setupAXIEmptyGuards(t *testing.T, root string) {
	const hook = `#!/bin/sh
# bench:managed-pre-push
# name: pre-push
# boundary: pre-push
# denies: nothing (informational)
# why: fixture has no deny-capable guards
`
	writeAXIFixture(t, filepath.Join(root, ".git", "hooks", "pre-push"), hook)
}

func setupAXIDiff(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "dirty.txt"), "dirty\n")
}

func setupAXICoverage(t *testing.T, root string) {
	const spec = `# Fixture

## User stories

1. As a caller, I get one row.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| FX1 | 1 | fixture | command | catches omission |
`
	writeAXIFixture(t, filepath.Join(root, "specs", "fixture", "spec.md"), spec)
}

func setupAXIEmptyCoverage(t *testing.T, root string) {
	const spec = `# Empty

## User stories

1. As a caller, I get no rows.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
`
	writeAXIFixture(t, filepath.Join(root, "specs", "empty", "spec.md"), spec)
}

// setupAXIConsumers plants the smallest real Go module the resolver can answer over: one
// declaration and one reference to it. The loader seam is package-internal, so this
// surface test drives the real go/packages path rather than a stub.
func setupAXIConsumers(t *testing.T, root string) {
	setupAXIEmptyConsumers(t, root)
	writeAXIFixture(t, filepath.Join(root, "consumer", "consumer.go"),
		"package consumer\n\nimport \"example.com/axifixture/target\"\n\nfunc Direct() { target.Symbol() }\n")
}

// setupAXIEmptyConsumers plants the same module without a consumer, so the resolved
// symbol has zero references and the response is the definitive empty table.
func setupAXIEmptyConsumers(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "go.mod"), "module example.com/axifixture\n\ngo 1.25\n")
	writeAXIFixture(t, filepath.Join(root, "target", "target.go"), "package target\n\nfunc Symbol() {}\n")
}

func setupAXIWorktree(t *testing.T, root string) {
	linked := filepath.Join(t.TempDir(), "linked")
	gittest.MustRun(t, root, "worktree", "add", "-q", "-b", "fixture-linked", linked)
}

func setupAXIRoadmap(t *testing.T, root string) {
	const heading = "**FT1 — fixture.**"
	roadmaptest.WriteSplitBoard(t, root, heading+"\n", map[string]string{"FT1.md": heading + "\n"})
}

func writeAXIFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir fixture parent: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}
