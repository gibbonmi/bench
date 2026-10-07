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

// Operator is a value the operator fills, such as the commit message. It always renders
// as its placeholder.
func Operator(name string) Word { return Word{kind: wordOperator, text: name} }

// Operators is a list of values the operator fills, such as the paths to commit.
func Operators(name string) Word { return Word{kind: wordOperators, text: name} }

func (w Word) render(values map[string]string) string {
	value := values[w.text]
	switch w.kind {
	case wordText:
		return w.text
	case wordFact:
		if value != "" && sanitize.LineSafe(value) {
			return sanitize.ShellQuote(value)
		}
	case wordComposed:
		return asWritten(w.text, value)
	case wordOperators:
		return placeholder(w.text) + "..."
	}
	return placeholder(w.text)
}

func asWritten(name, value string) string {
	if value != "" && sanitize.LineSafe(value) {
		return value
	}
	return placeholder(name)
}

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

func (s Step) render(values map[string]string) string {
	parts := make([]string, 0, len(s.words)+3)
	if s.target != nil {
		parts = append(parts, s.verb, treeTargetFlag, s.target.render(values))
	}
	for _, word := range s.words {
		parts = append(parts, word.render(values))
	}
	return strings.Join(parts, " ")
}

// Render is the one rendering of a face's route over a raising site's facts.
func (f Face) Render(facts Facts) string {
	steps := make([]string, 0, len(f.Route))
	for _, step := range f.Route {
		steps = append(steps, step.render(facts.Values))
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
