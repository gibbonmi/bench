package roadmap

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/capturetx"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// learningGrammar declares the argument shape for `bench learning`. The title is
// variadic like an idea's text; the body comes in through the three named flags, so the
// verb never has to split one prose blob into the journal's bullets. The link form takes
// the open entry's exact title and the fixing commit through its own two flags.
var learningGrammar = usage.Grammar{
	Cmd:  "bench learning",
	Help: `usage: bench learning "<title>" --what <text> --right <text> [--rule <text>] | --resolves "<title>" --commit <sha>`,
	Flags: []usage.Flag{
		{Name: "--what", HasValue: true, NoEmptyValue: true},
		{Name: "--right", HasValue: true, NoEmptyValue: true},
		{Name: "--rule", HasValue: true, NoEmptyValue: true},
		{Name: "--resolves", HasValue: true, NoEmptyValue: true},
		{Name: "--commit", HasValue: true, NoEmptyValue: true},
	},
	MaxArgs: -1,
}

// LearningCommand implements `bench learning <title...> --what --right [--rule]` and
// appends one open entry to capture/learnings.md through learnings.FormatEntry. With
// `--resolves <title> --commit <sha>` it instead links a fix to an existing entry
// (resolveLearning). The heading's `[open]` marker is the drain's state flag, and a
// hand-written entry has omitted it twice; the verb exists so the marker is never typed
// by hand. Inbox routing is the same as `bench idea`: an ignored journal writes the
// primary checkout's copy, a tracked one refuses the primary checkout.
func LearningCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(learningGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, hasResolves := parsed.Flags["--resolves"]
	_, hasCommit := parsed.Flags["--commit"]
	if hasResolves || hasCommit {
		return resolveLearning(parsed)
	}
	title := strings.TrimSpace(strings.Join(parsed.Positionals, " "))
	if title == "" {
		return learningGrammar.Help + "\n", 2
	}
	what, hasWhat := parsed.Flags["--what"]
	right, hasRight := parsed.Flags["--right"]
	if !hasWhat || !hasRight {
		return toon.MissingArg(learningGrammar.Cmd, "--what and --right") + "\n", 2
	}
	root, refusal, code := inboxRoot(learnings.JournalPath)
	if refusal != "" {
		return refusal, code
	}
	entry := learnings.FormatEntry(time.Now().Format("2006-01-02"), title, what, right, parsed.Flags["--rule"])
	// A prose refusal is a content error with correct grammar, so it exits 1 and prints
	// no usage line.
	if diagnostics := proseRefusals(root, entry); len(diagnostics) > 0 {
		return strings.Join(diagnostics, "\n") + "\n", 1
	}
	file := filepath.Join(root, learnings.JournalPath)
	err := capturetx.Append(root, capturetx.Source{Name: learnings.JournalPath, Path: file}, func(current []byte) ([]byte, error) {
		separator := "\n"
		if needsNewline(current) {
			separator = "\n\n"
		}
		return append(append([]byte(nil), current...), separator+entry...), nil
	})
	if err != nil {
		return cannotWrite(learnings.JournalPath, err), 1
	}
	return "captured: " + title + "\n", 0
}

// resolveLearning implements the link form, `bench learning --resolves <title> --commit
// <sha>`. It resolves the commit in the current checkout, then marks the one open entry
// with that exact title through learnings.MarkResolved inside the same capture
// transaction an append uses. A refusal changes no byte of the journal. The entry stays
// open: the drain reads the link and closes the entry.
func resolveLearning(parsed usage.Result) (string, int) {
	title, hasResolves := parsed.Flags["--resolves"]
	commit, hasCommit := parsed.Flags["--commit"]
	switch {
	case !hasResolves:
		return toon.MissingArg(learningGrammar.Cmd, "--resolves") + "\n", 2
	case !hasCommit:
		return toon.MissingArg(learningGrammar.Cmd, "--commit") + "\n", 2
	}
	for _, flag := range []string{"--what", "--right", "--rule"} {
		if _, joined := parsed.Flags[flag]; joined {
			return learningGrammar.Help + "\n", 2
		}
	}
	if len(parsed.Positionals) > 0 {
		return learningGrammar.Help + "\n", 2
	}
	title = strings.TrimSpace(title)
	root, refusal, code := inboxRoot(learnings.JournalPath)
	if refusal != "" {
		return refusal, code
	}
	checkout, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	full, err := git.ResolveCommit(checkout, commit, "--quiet")
	if err != nil || full == "" {
		return toon.Errorf("commit does not resolve: "+commit, "pass a commit that exists in this repository") + "\n", 1
	}
	// The default branch owns "fixed": a commit it does not reach is not landed work.
	def, ok := git.ResolvedDefault(checkout)
	if !ok {
		return toon.Errorf("default branch does not resolve", "repair Git metadata, then retry") + "\n", 1
	}
	if !git.OK("-C", checkout, "merge-base", "--is-ancestor", full, def) {
		return toon.Errorf("commit is not on the default branch "+def+": "+full, "land the fix, then link it") + "\n", 1
	}
	file := filepath.Join(root, learnings.JournalPath)
	err = capturetx.Append(root, capturetx.Source{Name: learnings.JournalPath, Path: file}, func(current []byte) ([]byte, error) {
		return learnings.MarkResolved(current, title, full)
	})
	switch {
	case errors.Is(err, learnings.ErrNoOpenEntry):
		// An open drain moved the sealed entries out of the live journal, so name it.
		if bundle, present, drainErr := capturetx.Current(root); drainErr == nil && present {
			return toon.Errorf(err.Error()+": "+title, "drain "+bundle.ID+" is open and sealed its entries; the drain closes them") + "\n", 1
		}
		return toon.Errorf(err.Error()+": "+title, "run bench learnings for the open titles") + "\n", 1
	case errors.Is(err, learnings.ErrAmbiguousEntry), errors.Is(err, learnings.ErrAlreadyResolved):
		return toon.Errorf(err.Error()+": "+title, "run bench learnings for the open titles") + "\n", 1
	case err != nil:
		return cannotWrite(learnings.JournalPath, err), 1
	}
	return "resolved: " + title + " by " + full + "\n", 0
}

// inboxRoot resolves the checkout whose copy of relPath a capture verb writes, or the
// refusal it returns instead. An ignored inbox is a local working note the helper routes
// to the primary checkout, so a worktree-parked entry survives the worktree's release. A
// tracked inbox keeps the landing boundary: main receives writes only through landings,
// so the verb refuses the primary checkout and appends to the phase worktree's copy.
func inboxRoot(relPath string) (root, refusal string, code int) {
	root, err := git.Root()
	if err != nil {
		return "", toon.NotInRepo() + "\n", 1
	}
	noteRoot, ignored, err := git.LocalNoteRoot(root, relPath)
	if err != nil {
		return "", toon.Errorf("checkout identity is unknown", "repair Git metadata, then retry") + "\n", 1
	}
	if ignored {
		return noteRoot, "", 0
	}
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil {
		return "", toon.Errorf("checkout identity is unknown", "repair Git metadata, then retry from a Bench worktree") + "\n", 1
	}
	if primary {
		return "", usage.PrimaryCheckoutRefusal() + "\n", 1
	}
	return root, "", 0
}

func cannotWrite(relPath string, err error) string {
	return toon.Errorf("cannot write "+relPath, err.Error()) + "\n"
}

// proseRefusals grades entry, composed exactly as it will be written, against the same
// rule the live-tree prose sweep applies once the journal is on disk. It returns each
// finding in the grader's own wording (RenderNamedResult), the one source, with line
// numbers relative to entry rather than to the file it will land in, so a refusal names
// the offending bullet the author can shorten without opening the journal. Each finding
// leads with the argument that supplied its line, so the author knows which flag to
// shorten. A policy this helper cannot load, an absent or broken exclusion file, is not
// this verb's failure to diagnose, so it grades nothing and lets the write proceed; the
// tree-wide sweep still owns that refusal.
func proseRefusals(root, entry string) []string {
	grader, diags := prose.NewGrader(root)
	if len(diags) > 0 {
		return nil
	}
	results := grader.GradeBytes(learnings.JournalPath, []byte(entry))
	if len(results) == 0 {
		return nil
	}
	out := make([]string, 0, len(results))
	for _, result := range results {
		out = append(out, "refused "+entryArgument(learnings.EntryField(result.Line))+": "+prose.RenderNamedResult(result))
	}
	return out
}

// entryArgument spells an entry field the way the caller typed it: the positional title,
// or the flag that carries a body bullet. A finding with no field line names the entry.
func entryArgument(field string) string {
	switch field {
	case "title":
		return "title"
	case "":
		return "entry"
	default:
		return "--" + field
	}
}
