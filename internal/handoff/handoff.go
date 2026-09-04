package handoff

import (
	"errors"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there rather than a local switch.
// The `--harness` *value* is this package's to check: the shared parser validates flag
// spelling, repetition, and arity, and accepts anything as a declared flag's value.
var grammar = usage.Grammar{
	Cmd:  "bench handoff",
	Help: "usage: bench handoff [--harness " + status.HarnessChoices() + "] [--next <command>]",
	Flags: []usage.Flag{
		{Name: "--harness", HasValue: true, NoEmptyValue: true},
		// An empty override names no command. Left to fall through it would read as a
		// clean board, which is a different — and false — statement about the tree.
		{Name: "--next", HasValue: true, NoEmptyValue: true},
	},
}

// Command implements `bench handoff`. It resolves the section the caller's checkout owns,
// rewrites that one section, and re-emits every other section's bytes as the file had
// them. The primary checkout owns `main`. A Bench worktree owns the section of the active
// assignment that holds it. A checkout that is neither owns nothing and writes nothing.
//
// It creates freely and destroys never. A missing file is scaffolded, the reviewer-owned
// State body of the owned section passes through untouched, and a document that cannot be
// parsed is left exactly as it was behind a non-zero exit.
//
// Exits: 2 for a usage error, 1 outside a repository, from a checkout that owns no
// section, or when the document cannot be parsed, rendered, or written, 0 otherwise.
// Nothing is written on any non-zero path.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	harness := status.HarnessClaude
	if value, present := parsed.Flags["--harness"]; present {
		if !status.ValidHarness(value) {
			return toon.Usage(grammar.Cmd, "--harness "+value) + "\n", 2
		}
		harness = value
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	owned, err := resolveOwner(root)
	if err != nil {
		return err.Error() + "\n", 1
	}

	// An ignored handoff is a local pin: the primary checkout's copy is the one a
	// cold session reads, so the write goes there from any checkout. A tracked
	// handoff keeps the caller's checkout and lands with its phase.
	noteRoot, _, err := git.LocalNoteRoot(root, status.HandoffFile)
	if err != nil {
		return refusal{"cannot resolve the handoff checkout", err.Error()}.Error() + "\n", 1
	}
	target := filepath.Join(noteRoot, status.HandoffFile)

	p := plan{facts: collect(noteRoot, owned), tip: sectionTip(noteRoot, owned), scanRoot: noteRoot}
	if override, present := parsed.Flags["--next"]; present {
		p.facts.Action, p.overridden = override, true
	} else {
		p.route = func() status.RouteResult {
			return status.RouteFor(root, status.SignalsWith(root, status.Query{ExcludeDirtyPaths: []string{status.HandoffFile}}), harness)
		}
	}

	pin, err := writeSection(target, p)
	if err != nil {
		// A refusal already reads as the answer. Wrapping it would bury the offending
		// State line, or the unrepresentable field, under a write-failure headline.
		var refuse refusal
		if errors.As(err, &refuse) {
			return refuse.Error() + "\n", 1
		}
		return refusal{"cannot write the session handoff", err.Error()}.Error() + "\n", 1
	}
	return pin, 0
}

// plan is what one run brings to the locked read-modify-write. The Next command and the
// State scan both need the document's own bytes, which only exist inside the lock, so the
// decisions ride in as inputs rather than as an already-derived value.
//
// route is nil when `--next` named the command. A non-nil route is called only for a
// section whose own Next command is blank, so a run that keeps an existing invocation
// never pays for the board scan.
type plan struct {
	facts      facts
	overridden bool
	route      func() status.RouteResult
	tip        string // the owned section's worktree tip, or "" when none resolves
	scanRoot   string // the checkout whose object store the State scan reads
}

// owner names the section one run rewrites. Assignment is the record behind a request
// section, absent for main, which no assignment owns. Worktree is the tree the section's
// spec pins are read from.
type owner struct {
	Key        string
	Assignment intent.Assignment
	Worktree   string
}

// resolveOwner answers which section this checkout owns. The primary checkout owns main,
// so a phase close with nothing live still writes. Any other checkout owns the section of
// the active assignment that holds it, keyed by that assignment's request digest.
//
// The digest is the key, never the label or the path string. A label is a caller's name
// for a tree and repeats across requests; a path string spelled through a symlink names
// one tree two ways. Either as a key adopts a section some other request owns.
//
// A checkout that is neither refuses. It cannot name a section without guessing, and a
// guess writes over a live phase's pins.
func resolveOwner(root string) (owner, error) {
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil {
		return owner{}, refusal{"checkout identity is unknown",
			"repair Git metadata, then rerun bench handoff"}
	}
	if primary {
		return owner{Key: handoffdoc.MainKey, Worktree: root}, nil
	}
	assignment, found := intent.AssignmentForWorktree(root)
	if !found {
		return owner{}, refusal{"this checkout owns no handoff section",
			"run bench handoff from the primary checkout, or from a Bench worktree whose assignment is active"}
	}
	return owner{Key: assignment.Request, Assignment: assignment, Worktree: assignment.Worktree}, nil
}

// writeSection runs the locked read-modify-write and returns what the command prints. The
// read, the rewrite, and the replace happen under the leaf package's lock, so a sibling
// phase writing its own section at the same moment loses nothing.
//
// State is read back out of the document and put back unchanged. This command derives
// every other field of the owned section and never rewrites State.
//
// Every refusal is raised before the document is rendered, and Update replaces the file
// only when the callback returns nil. So a refused run leaves the bytes it read.
func writeSection(target string, p plan) (string, error) {
	var pin string
	err := handoffdoc.Update(target, func(doc *handoffdoc.Document) error {
		existing, _ := doc.Section(p.facts.Key)
		if err := scanState(p.scanRoot, p.tip, existing.State); err != nil {
			return err
		}
		f, next := p.facts, existing.Next
		if !p.overridden && blankNext(next) {
			applyRoute(&f, p.route())
		}
		if err := validate(f); err != nil {
			return err
		}
		if p.overridden || blankNext(next) {
			next = nextField(f)
		}
		owned := section(f, existing.State, next)
		doc.Header = header(f, existing.State)
		doc.Shape = ShapeSection
		doc.Put(owned)
		doc.EnsureMain()
		pin = preview(doc.Header, owned)
		return nil
	})
	if err != nil {
		return "", err
	}
	return pin, nil
}

// sectionTip names the commit the State scan measures ancestry against: the tip of the
// tree the owned section describes, never the caller's. main's tip is the document
// checkout's HEAD, which is the tree main's section reports on.
//
// An empty answer means no tip resolves. The scan then has nothing to measure against and
// refuses nothing: a repository with no commits is a state to pass through, not a defect.
func sectionTip(documentRoot string, o owner) string {
	if o.Key == handoffdoc.MainKey {
		sha, err := git.Output("-C", documentRoot, "rev-parse", "HEAD")
		if err != nil {
			return ""
		}
		return sha
	}
	if tip := worktreeTip(documentRoot, o.Assignment); tip != unknownTip {
		return tip
	}
	return ""
}
