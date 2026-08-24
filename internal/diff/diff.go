// Package diff owns the coherent Git review snapshot. A live response resolves its
// branch range once, renders every bounded and complete section from that attempt, and
// compares each patch-observable identity before emitting any bytes. A named commit
// stays an immutable first-parent view and omits unrelated checkout facts.
package diff

import (
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand:
// arity, flag recognition, `--`, repeated flags, and help all come from there rather
// than a local switch.
// Help is fullHelp without its trailing newline, because the caller appends one.
var grammar = usage.Grammar{
	Cmd:   "bench diff",
	Help:  strings.TrimSuffix(fullHelp, "\n"),
	Flags: []usage.Flag{{Name: "--full"}, {Name: "--commit", HasValue: true}, {Name: "--base", HasValue: true, NoEmptyValue: true}},
}

const fullHelp = `usage: bench diff [--full] [--commit <sha>] [--base <commit>]
  Live mode reports one movement-checked revision, aggregate, files, checkout,
  and whitespace snapshot. --full also appends the commit log and verbatim
  tracked patch, then path-sorted raw Git patches for untracked regular files.
  --commit <sha> reports the immutable first-parent view of one resolved commit
  and omits unrelated live-checkout facts. Bounded results advertise the exact
  --full invocation; complete and clean results advertise an empty help table.
`

// Command implements `bench diff`.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
	commitArg, hasCommit := parsed.Flags["--commit"]
	baseArg, hasBase := parsed.Flags["--base"]
	if hasCommit && hasBase {
		return toon.Usage(grammar.Cmd, "--base") + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	invocation := append([]string{"diff"}, args...)
	out, drift, driftHint, errKind, errHint := renderAttempt(root, hasCommit, commitArg, hasBase, baseArg, full)
	if errKind != "" {
		return toon.Errorf(errKind, errHint) + "\n", 1
	}
	if drift == "" {
		return out, 0
	}
	help, err := axi.RenderHelp([]axi.Action{axi.RetryDiff(invocation)})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return toon.Errorf("snapshot drift", driftHint) + "\n" + help, 1
}

func renderAttempt(root string, hasCommit bool, commitArg string, hasBase bool, baseArg string, full bool) (out, drift, driftHint, errKind, errHint string) {
	var dr diffRange
	if hasCommit {
		dr, errKind, errHint = resolveCommitRange(root, commitArg)
		if errKind != "" {
			return "", "", "", errKind, errHint
		}
		out, errKind, errHint := renderCommit(root, dr, full)
		return out, "", "", errKind, errHint
	}
	result := MovementCheckedRetry(root, func(snapshot MovementSnapshot) (kind, hint string) {
		if hasBase {
			dr, kind, hint = resolveExplicitRange(root, baseArg, snapshot.Facts.Head)
		} else {
			dr, kind, hint = resolveBranchRangeFromFacts(root, snapshot.Facts)
		}
		if kind != "" {
			return kind, hint
		}
		out, kind, hint = renderLive(root, dr, snapshot.Facts, full)
		return kind, hint
	})
	return out, result.DriftKind, result.DriftHint, result.Kind, result.Hint
}

func renderLive(root string, dr diffRange, facts git.DiffFacts, full bool) (string, string, string) {
	tracked, err := changedFilesAt(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --name-status failed", err.Error()
	}
	files := make([][]string, 0, len(tracked)+len(facts.Changes))
	for _, row := range tracked {
		files = append(files, []string{row[0], row[1], ""})
	}
	for _, entry := range facts.Changes {
		if entry.Status == "??" {
			files = append(files, []string{"?", entry.Path, fileKind(root, entry.Path)})
		}
	}
	commits, err := commitLogAt(root, dr.logRange)
	if err != nil {
		return "", "git log failed", err.Error()
	}
	insertions, deletions, err := shortstat(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --shortstat failed", err.Error()
	}
	clean, offenses, err := whitespace(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --check failed", err.Error()
	}
	staged, unstaged, untracked := statusCounts(facts.Changes)
	branch := facts.Branch
	if branch == "HEAD" {
		branch = "(detached)"
	}
	revision, err := toon.Table("revision", []string{"branch", "default", "ahead", "behind", "base", "method", "head"}, [][]string{{branch, facts.DefaultBranch, strconv.Itoa(facts.Ahead), strconv.Itoa(facts.Behind), dr.base, dr.method, facts.Head}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	aggregate, err := toon.Table("aggregate", []string{"commits", "files", "insertions", "deletions", "staged", "unstaged", "untracked"}, [][]string{{strconv.Itoa(len(commits)), strconv.Itoa(len(files)), strconv.Itoa(insertions), strconv.Itoa(deletions), strconv.Itoa(staged), strconv.Itoa(unstaged), strconv.Itoa(untracked)}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	fileTable, err := toon.Table("files", []string{"status", "path", "kind"}, files)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	checkout, err := toon.Table("checkout", []string{"index", "worktree", "path"}, checkoutRows(facts.Changes))
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	white, err := toon.Table("whitespace", []string{"clean", "offenses"}, [][]string{{strconv.FormatBool(clean), strconv.Itoa(offenses)}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	var b strings.Builder
	b.WriteString(revision)
	b.WriteString(aggregate)
	b.WriteString(fileTable)
	b.WriteString(checkout)
	b.WriteString(white)
	if full {
		logTable, err := toon.Table("log", []string{"sha", "subject"}, commits)
		if err != nil {
			return "", "unrepresentable TOON cell", err.Error()
		}
		body, err := diffBodyAt(root, dr.bodyArgs...)
		if err != nil {
			return "", "git diff failed", err.Error()
		}
		untrackedBody, err := untrackedPatches(root, facts)
		if err != nil {
			return "", "git diff --no-index failed", err.Error()
		}
		b.WriteString(logTable)
		b.WriteString("diff_body:\n")
		b.Write(body)
		b.Write(untrackedBody)
	}
	actions := []axi.Action(nil)
	if !full && len(files) > 0 {
		if dr.method == "explicit" {
			actions = append(actions, axi.InspectFullBase(dr.base))
		} else {
			actions = append(actions, axi.InspectFull(""))
		}
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	b.WriteString(help)
	return b.String(), "", ""
}

func renderCommit(root string, dr diffRange, full bool) (string, string, string) {
	files, err := changedFilesAt(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --name-status failed", err.Error()
	}
	commits, err := commitLogAt(root, dr.logRange)
	if err != nil {
		return "", "git log failed", err.Error()
	}
	insertions, deletions, err := shortstat(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --shortstat failed", err.Error()
	}
	revision, err := toon.Table("revision", []string{"commit", "base", "method"}, [][]string{{dr.head, dr.base, dr.method}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	aggregate, err := toon.Table("aggregate", []string{"commits", "files", "insertions", "deletions"}, [][]string{{strconv.Itoa(len(commits)), strconv.Itoa(len(files)), strconv.Itoa(insertions), strconv.Itoa(deletions)}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	for i := range files {
		files[i] = append(files[i], "")
	}
	fileTable, err := toon.Table("files", []string{"status", "path", "kind"}, files)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	var b strings.Builder
	b.WriteString(revision)
	b.WriteString(aggregate)
	b.WriteString(fileTable)
	if full {
		logTable, err := toon.Table("log", []string{"sha", "subject"}, commits)
		if err != nil {
			return "", "unrepresentable TOON cell", err.Error()
		}
		body, err := diffBodyAt(root, dr.bodyArgs...)
		if err != nil {
			return "", "git diff failed", err.Error()
		}
		b.WriteString(logTable)
		b.WriteString("diff_body:\n")
		b.Write(body)
	}
	actions := []axi.Action(nil)
	if !full && len(files) > 0 {
		actions = append(actions, axi.InspectFull(dr.head))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	b.WriteString(help)
	return b.String(), "", ""
}
