//go:build system

package systemtest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
	"github.com/gibbonmi/bench/internal/treetarget/treetargettest"
)

// The row cells here are authored apart from the renderer: the spec fixes one row of the
// target, the full HEAD commit, and the dirty cell. The header and the TOON quoting of a
// cell come from their owners, so this file restates neither.

// onlyInAlphaSpec is a repo-relative spec path that only the alpha branch holds.
const onlyInAlphaSpec = "specs/only-in-alpha/spec.md"

// treeTargetFixture answers a scaffolded repository and one active assignment labeled alpha.
// The alpha branch holds one commit that adds onlyInAlphaSpec, so the two trees differ in
// both their HEAD and their content.
func treeTargetFixture(t *testing.T) (recordedPublicationRepo, systemLandingWorktree) {
	t.Helper()
	fixture := scaffoldRecordedPublicationRepo(t)
	alpha := systemCreateLandingWorktree(t, fixture.root, fixture.home, "alpha", "tree target")
	if err := os.MkdirAll(filepath.Join(alpha.path, filepath.Dir(onlyInAlphaSpec)), 0o755); err != nil {
		t.Fatal(err)
	}
	systemCommit(t, alpha.path, onlyInAlphaSpec, "# Only in alpha\n", "spec only in alpha")
	return fixture, alpha
}

// runTreeTarget runs the selected executable in dir. The environment removes the invoking
// wrapper, so a tree-target child is the running executable. It also removes the command
// observation, so stderr holds only the lines of the verb.
func runTreeTarget(t *testing.T, fixture recordedPublicationRepo, dir string, args ...string) processResult {
	t.Helper()
	unobserved := capability.WithoutEnvironment(fixture.environment(fixture.home), "BENCH_COMMAND_OBSERVE")
	return systemSelected(t, dir, append(unobserved, env.WrapperEnv), args...)
}

// treeRowCells answers the cells of the identity row that leads stdout, unquoted and joined
// by commas, or the empty string when stdout does not start with the row block.
func treeRowCells(stdout string) string {
	block := strings.TrimSuffix(stdout, treetargettest.WithoutRow(stdout))
	_, row, _ := strings.Cut(strings.TrimSuffix(block, "\n"), "\n")
	cells := strings.Split(strings.TrimSpace(row), ",")
	for i, cell := range cells {
		if value, err := systemTOONCell(cell); err == nil {
			cells[i] = value
		}
	}
	return strings.Join(cells, ",")
}

// treeReadPages follows argv and every continuation that its pages offer through the
// selected executable from dir. Each response fits the response bound with its identity
// row, and a nonempty row must match the row cells of each page. It answers the decoded
// bytes and the page count.
func treeReadPages(t *testing.T, fixture recordedPublicationRepo, dir, row string, argv ...string) ([]byte, int) {
	t.Helper()
	return axitest.WalkPages(t, argv, func(argv []string) string {
		result := runTreeTarget(t, fixture, dir, argv...)
		if result.code != 0 || result.stderr != "" || treeRowCells(result.stdout) != row || !responsebound.Fits(result.stdout) {
			t.Fatalf("bench %q = (%d, %q, %q), want the row %q and one bounded page", argv, result.code, result.stdout, result.stderr, row)
		}
		return treetargettest.WithoutRow(result.stdout)
	}, responsebound.DecodeContent)
}

// The worktree read rows read pagesDoc, which only the alpha branch holds. Its long section
// crosses several pages and a fenced heading, and its lines end in CRLF.
const pagesDoc = "pages.md"

var pagesLong = "## Long\n" + strings.Repeat("long line é\r\n", 600) + "```text\n## Hidden\n```\n\n"

var pagesBody = "# Pages\n\n## Before\nbefore\n\n" + pagesLong + "## After\nafter without newline"

// TT26, TT28, TT41, TT54, TT55: a tree-scoped verb names its tree in its first block. A tree
// target runs the verb in the named tree, from the primary checkout and from a worktree,
// and a relative operand resolves in that tree. The exec route and directory inference
// still name the tree that they run in.
//
// ER29, ER30, ER31: the worktree reader runs in the named tree through its child, and every
// page and continuation names that tree. A section read reconstructs the span that root read
// selects in that tree, a line range keeps its CRLF terminators, an invalid range refuses
// without content, and a target with no active assignment refuses before any child starts.
//
// ER34: the related view of a spec pages through the selected executable, and its rebuilt
// projection holds every bound assignment in ID order with the matched and unbound counts.
func TestTreeTargetRunsInNamedWorktree(t *testing.T) {
	t.Parallel()
	fixture, alpha := treeTargetFixture(t)
	systemCommit(t, alpha.path, pagesDoc, pagesBody, "pages only in alpha")
	primaryRow := "primary," + systemGitOutput(t, fixture.root, "rev-parse", "HEAD") + ",false"
	alphaRow := "alpha," + systemGitOutput(t, alpha.path, "rev-parse", "HEAD") + ",false"
	for _, row := range []struct {
		name, dir, row, body string
		args                 []string
	}{
		{name: "TT26 a label from the primary checkout", dir: fixture.root, row: alphaRow, args: []string{"status", "--in", "alpha"}},
		{name: "TT28 the primary keyword from a worktree", dir: alpha.path, row: primaryRow, args: []string{"status", "--in", "primary"}},
		{name: "TT41 a relative operand in the target", dir: fixture.root, row: alphaRow, body: "\nspec: " + onlyInAlphaSpec + "\n", args: []string{"coverage", "--in", "alpha", onlyInAlphaSpec}},
		{name: "TT54 the exec route", dir: fixture.root, row: alphaRow, args: []string{"worktree", "exec", "alpha", "--", owner.selected.path, "status"}},
		{name: "TT55 directory inference", dir: alpha.path, row: alphaRow, args: []string{"status"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			result := runTreeTarget(t, fixture, row.dir, row.args...)
			if result.code != 0 || treeRowCells(result.stdout) != row.row || !strings.Contains(result.stdout, row.body) {
				t.Fatalf("bench %q in %s = (%d, %q, %q), want exit 0, the row %q first, and %q", row.args, row.dir, result.code, result.stdout, result.stderr, row.row, row.body)
			}
		})
	}
	read := []string{"worktree", "read", "--in", "alpha", pagesDoc}
	t.Run("ER29 a worktree read reconstructs every page through the child", func(t *testing.T) {
		if got, pages := treeReadPages(t, fixture, fixture.root, alphaRow, read...); string(got) != pagesBody || pages < 2 {
			t.Fatalf("worktree read --in alpha = %d pages of %q, want several pages of %q", pages, got, pagesBody)
		}
	})
	t.Run("ER30 a worktree section read matches the root reader in that tree", func(t *testing.T) {
		got, pages := treeReadPages(t, fixture, fixture.root, alphaRow, append(read, "--section", "Long")...)
		root, _ := treeReadPages(t, fixture, alpha.path, "", "read", pagesDoc, "--section", "Long")
		if string(got) != pagesLong || string(root) != pagesLong || pages < 2 {
			t.Fatalf("section Long = %d pages of %q, root read = %q, want several pages of %q", pages, got, root, pagesLong)
		}
	})
	t.Run("ER31 a worktree line range keeps its terminators", func(t *testing.T) {
		if got, _ := treeReadPages(t, fixture, fixture.root, alphaRow, append(read, "--lines", "6:8")...); string(got) != "## Long\nlong line é\r\nlong line é\r\n" {
			t.Fatalf("lines 6:8 = %q, want the heading and two CRLF lines", got)
		}
		result := runTreeTarget(t, fixture, fixture.root, append(read, "--lines", "1:99999")...)
		if result.code != 1 || treeRowCells(result.stdout) != alphaRow || !strings.HasPrefix(treetargettest.WithoutRow(result.stdout), "error: line range unavailable") || strings.Contains(result.stdout, "content[") {
			t.Fatalf("lines 1:99999 = (%d, %q, %q), want the line-range refusal with no content", result.code, result.stdout, result.stderr)
		}
	})
	t.Run("ER29 an unassigned target refuses before any child", func(t *testing.T) {
		result := runTreeTarget(t, fixture, fixture.root, "worktree", "read", "--in", "beta", pagesDoc)
		if want := "bench worktree read --in: target is unassigned\nnext=bench worktree list\n"; result.code != 1 || result.stdout != "" || result.stderr != want {
			t.Fatalf("worktree read --in beta = (%d, %q, %q), want (1, \"\", %q)", result.code, result.stdout, result.stderr, want)
		}
	})
	t.Run("ER34 related worktrees page in assignment-ID order with complete counts", func(t *testing.T) {
		_, _, _, before := relatedPages(t, fixture)
		related := relatedAssignments(t, fixture.root)
		ids, pages, summary, _ := relatedPages(t, fixture)
		want := fmt.Sprintf("specs/related/spec.md,%d,%d", len(related), before+1)
		if pages < 2 || summary != want || !slices.Equal(ids, related) {
			t.Fatalf("related view = %d pages, summary %s, ids %q; want several pages, summary %s, ids %q", pages, summary, ids, want, related)
		}
	})
}

// relatedPages follows every page of the related view of specs/related/spec.md through the
// selected executable. It answers the row IDs of the rebuilt projection, the page count, and
// the summary cells, and the unbound count apart.
func relatedPages(t *testing.T, fixture recordedPublicationRepo) ([]string, int, string, int) {
	t.Helper()
	got, pages := treeReadPages(t, fixture, fixture.root, "", "worktree", "list", "--spec", "related")
	document, err := axitest.DecodeDocument(string(got))
	if err != nil {
		t.Fatalf("decode %q: %v", got, err)
	}
	summary, err := document.Rows("related")
	if err != nil || len(summary) != 1 {
		t.Fatalf("related summary = %v, %v", summary, err)
	}
	rows, err := document.Rows("worktrees")
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, row := range rows {
		ids = append(ids, row.(map[string]any)["id"].(string))
	}
	cells := summary[0].(map[string]any)
	unbound, err := strconv.Atoi(fmt.Sprint(cells["unbound"]))
	if err != nil {
		t.Fatal(err)
	}
	return ids, pages, fmt.Sprint(cells["spec"], ",", cells["matched"], ",", unbound), unbound
}

// relatedAssignments records forty-one assignments with absent trees in descending ID order,
// binds the first forty to specs/related/spec.md, and leaves the last one unbound. It answers
// the bound assignment IDs in ascending order.
func relatedAssignments(t *testing.T, root string) []string {
	t.Helper()
	owner := strings.Repeat("b", 32)
	start := systemGitOutput(t, root, "rev-parse", "HEAD")
	var related []string
	for i := 41; i >= 1; i-- {
		id := fmt.Sprintf("%032x", i)
		assignment := intent.Assignment{Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner, Request: intent.RequestDigest(id), RequestToken: "related-" + id,
			Label: "related " + id[28:], Start: start, Branch: intent.AssignmentBranchRef(owner, id), Worktree: filepath.Join(root, "absent", id), State: intent.StateComplete}
		if err := intent.PutAssignment(root, assignment); err != nil {
			t.Fatal(err)
		}
		if i > 1 {
			related = append(related, id)
		}
	}
	commitmenttest.EditCommitment(t, root, func(state *intent.CommitmentState) {
		for _, id := range related {
			state.Bindings = append(state.Bindings, intent.DeliveryBinding{Assignment: id, Request: intent.RequestDigest(id), Milestone: "fixture", Outcome: "related-" + id, Deliverable: "specs/related/spec.md", Identity: "fixture"})
		}
	})
	slices.Sort(related)
	return related
}

// TT50, BC64: the real executable, with BENCH_KIT set to the kit root, refuses a kit
// worktree target whose build no longer matches the tree, and it starts no child. The build
// is a script that freshness.Publish seals and that writes a marker when it runs. The primary
// checkout holds no build that matches the target tree, and the seal of the running
// executable names the kit tree, so no candidate matches.
func TestTreeTargetRefusesStaleKitBuild(t *testing.T) {
	t.Parallel()
	fixture, alpha := treeTargetFixture(t)
	kittest.WriteTree(t, alpha.path)
	marker := filepath.Join(t.TempDir(), "marker")
	script := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n: > '"+marker+"'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	build := freshness.PublishedExecutable(alpha.path)
	if err := freshness.Publish(alpha.path, script, build, filepath.Dir(build), "tree-target"); err != nil {
		t.Fatal(err)
	}
	kittest.EditBuildInput(t, alpha.path)
	result := runTreeTarget(t, fixture, fixture.root, "status", "--in", "alpha")
	want := "bench status --in: worktree build does not match the tree\nnext=bench worktree build alpha\n"
	if result.code != 1 || result.stdout != "" || result.stderr != want {
		t.Fatalf("bench status --in alpha = (%d, %q, %q), want (1, \"\", %q)", result.code, result.stdout, result.stderr, want)
	}
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a child started: %v", err)
	}
}
