package refusalroute

import (
	"strings"

	"github.com/gibbonmi/bench/internal/sanitize"
)

const (
	// reviewerMarker opens a reviewer route, so an agent that reads it stops and hands back.
	reviewerMarker = "reviewer: "
	stepJoiner     = "; then "
	// prefaceJoiner separates the preface from the first step. The preface qualifies the
	// route and is not a step to run, so it takes no "then".
	prefaceJoiner = "; "
	prefaceSlot   = "preface"
	// treeTargetFlag addresses a tree-scoped Bench verb at a worktree. It counts only as
	// the first argument after the verb.
	treeTargetFlag = "--in"
)

type wordKind uint8

const (
	wordText wordKind = iota
	wordFact
	wordComposed
	wordOperator
	wordOperators
)

// Word is one part of a step template: fixed text or a slot over a named fact.
type Word struct {
	kind wordKind
	text string
}

// Text is fixed route text.
func Text(text string) Word { return Word{kind: wordText, text: text} }

// Fact is a value the raising site supplies, such as a label, a commit, or a path. It
// renders shell-quoted, or as its placeholder when the value is absent or not line-safe.
func Fact(name string) Word { return Word{kind: wordFact, text: name} }

// Composed is a command the raising site already composed, such as the caller's re-run.
// It renders as written, or as its placeholder when it is absent or not line-safe.
func Composed(name string) Word { return Word{kind: wordComposed, text: name} }

// Operator is a value the operator fills, such as the commit message. A printed route
// renders it as its placeholder.
func Operator(name string) Word { return Word{kind: wordOperator, text: name} }

// Operators is a list of values the operator fills, such as the paths to commit.
func Operators(name string) Word { return Word{kind: wordOperators, text: name} }

// slotFill answers the value each slot reads. operators states whether the operator's own
// slots read it too. The production fill leaves those slots as placeholders for the
// operator; only a check that hands a step to a shell reader fills them, because the
// reader takes a placeholder such as <msg> for a redirection.
type slotFill struct {
	value     func(slot string) string
	operators bool
}

// factsFill is the production fill: it reads the raising site's facts.
func factsFill(values map[string]string) slotFill {
	return slotFill{value: func(slot string) string { return values[slot] }}
}

func (w Word) render(fill slotFill) string {
	value := fill.value(w.text)
	switch w.kind {
	case wordText:
		return w.text
	case wordFact:
		return quoted(w.text, value)
	case wordComposed:
		return asWritten(w.text, value)
	case wordOperator, wordOperators:
		if fill.operators {
			return quoted(w.text, value)
		}
		if w.kind == wordOperators {
			return placeholder(w.text) + "..."
		}
	}
	return placeholder(w.text)
}

// quoted renders a value shell-quoted, or as its placeholder when the value is absent or
// not line-safe.
func quoted(name, value string) string {
	if pasteable(value) {
		return sanitize.ShellQuote(value)
	}
	return placeholder(name)
}

func asWritten(name, value string) string {
	if pasteable(value) {
		return value
	}
	return placeholder(name)
}

// pasteable reports whether a value can stand in a printed route: present, and free of any
// byte that would break the line.
func pasteable(value string) bool { return value != "" && sanitize.LineSafe(value) }

func placeholder(name string) string { return "<" + name + ">" }

// Step is one step of a route: a command the agent runs, or an instruction for an action
// the agent does with its own tools.
type Step struct {
	command bool
	verb    string
	target  *Word
	words   []Word
}

// Command is a step that runs as its words read.
func Command(words ...Word) Step { return Step{command: true, words: words} }

// TreeCommand is a step that runs the tree-scoped Bench verb at the worktree that target
// names. The target renders first after the verb, the one place the flag counts.
func TreeCommand(verb string, target Word, words ...Word) Step {
	return Step{command: true, verb: verb, target: &target, words: words}
}

// Instruction is one imperative sentence for an action the agent does with its own tools.
func Instruction(words ...Word) Step { return Step{words: words} }

// IsCommand reports whether the step runs as its words read, rather than an instruction
// the agent carries out with its own tools.
func (s Step) IsCommand() bool { return s.command }

// Filled renders the step with every slot read from fill, the operator's slots included.
// A check that hands the step to a shell reader renders it here, so no placeholder reaches
// the reader.
func (s Step) Filled(fill func(slot string) string) string {
	return s.render(slotFill{value: fill, operators: true})
}

func (s Step) render(fill slotFill) string {
	parts := make([]string, 0, len(s.words)+3)
	if s.target != nil {
		parts = append(parts, s.verb, treeTargetFlag, s.target.render(fill))
	}
	for _, word := range s.words {
		parts = append(parts, word.render(fill))
	}
	return strings.Join(parts, " ")
}

// Steps splits the steps of a rendered route, in route order. A check that carries out a
// printed route step by step reads the steps here, so the check and Render share one joiner.
func Steps(route string) []string { return strings.Split(route, stepJoiner) }

// Render is the one rendering of a face's route over a raising site's facts.
func (f Face) Render(facts Facts) string {
	fill := factsFill(facts.Values)
	steps := make([]string, 0, len(f.Route))
	for _, step := range f.Route {
		steps = append(steps, step.render(fill))
	}
	route := strings.Join(steps, stepJoiner)
	if facts.Preface != "" {
		route = asWritten(prefaceSlot, facts.Preface) + prefaceJoiner + route
	}
	if f.Authority == Reviewer {
		route = reviewerMarker + route
	}
	return route
}
