package roadmap

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// learningGrammar declares the argument shape for `bench learning`. The title is
// variadic like an idea's text; the body comes in through the three named flags, so the
// verb never has to split one prose blob into the journal's bullets.
var learningGrammar = usage.Grammar{
	Cmd:  "bench learning",
	Help: `usage: bench learning "<title>" --what <text> --right <text> [--rule <text>]`,
	Flags: []usage.Flag{
		{Name: "--what", HasValue: true, NoEmptyValue: true},
		{Name: "--right", HasValue: true, NoEmptyValue: true},
		{Name: "--rule", HasValue: true, NoEmptyValue: true},
	},
	MaxArgs: -1,
}

// LearningCommand implements `bench learning <title...> --what --right [--rule]` and
// appends one open entry to capture/learnings.md through learnings.FormatEntry. The
// heading's `[open]` marker is the drain's state flag, and a hand-written entry has
// omitted it twice; the verb exists so the marker is never typed by hand. Inbox routing
// is the same as `bench idea`: an ignored journal writes the primary checkout's copy, a
// tracked one refuses the primary checkout.
func LearningCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(learningGrammar, args)
	if line != "" {
		return line + "\n", code
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
	if diagnostics := proseRefusals(root, entry); len(diagnostics) > 0 {
		return learningGrammar.Help + "\n" + strings.Join(diagnostics, "\n") + "\n", 2
	}
	file := filepath.Join(root, learnings.JournalPath)
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return cannotWrite(learnings.JournalPath, err), 1
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return cannotWrite(learnings.JournalPath, err), 1
	}
	defer f.Close()
	// A blank line separates the new heading from whatever ends the file; the missing
	// trailing newline case folds into the same single write as the idea verb's.
	switch {
	case needsNewline(file):
		entry = "\n\n" + entry
	default:
		entry = "\n" + entry
	}
	if _, err := f.WriteString(entry); err != nil {
		return cannotWrite(learnings.JournalPath, err), 1
	}
	return "captured: " + title + "\n", 0
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
// the offending bullet the author can shorten without opening the journal. A policy this
// helper cannot load, an absent or broken exclusion file, is not this verb's failure to
// diagnose, so it grades nothing and lets the write proceed; the tree-wide sweep still
// owns that refusal.
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
		out = append(out, prose.RenderNamedResult(result))
	}
	return out
}
