package adopt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// setupQuestion is one genuine ambiguity setup cannot infer on its own: a labelled
// prompt with an ordered list of candidate answers. Adding a second ambiguity source
// beyond the gate-inference table (FT76 story 9) is one more entry in
// setupQuestionSet, not a second sequencing mechanism — askSetupQuestions and
// askQuestions stay the one seam every ambiguity resolves through.
type setupQuestion struct {
	prompt  string
	options []string
	apply   func(choice int, facts *setupFacts)
}

// setupQuestionSet builds every open question this slice knows how to ask, from the
// same facts the preview and the gate table already derive from — one source, so the
// question text setup prints here never drifts from detectGateCandidates.
func setupQuestionSet(facts setupFacts) []setupQuestion {
	if len(facts.gateCandidates) < 2 {
		return nil
	}
	options := make([]string, len(facts.gateCandidates))
	for i, c := range facts.gateCandidates {
		options[i] = c.label()
	}
	return []setupQuestion{{
		prompt:  "multiple test commands detected - which one should .bench/gate.sh run?",
		options: options,
		apply: func(choice int, f *setupFacts) {
			f.gateCommand = f.gateCandidates[choice].command
			f.openQuestions = nil
		},
	}}
}

// askSetupQuestions asks every genuine ambiguity in facts one at a time over the
// injected reader/writer (FT76 story 3): an ambiguity-free plan calls this with zero
// questions and returns facts unchanged, so the only interaction left is the single
// final confirm. reader is shared with the confirm prompt that follows so a single
// bufio buffer never strands an already-typed confirm answer.
func askSetupQuestions(facts setupFacts, reader *bufio.Reader, stdout io.Writer) (setupFacts, error) {
	if err := askQuestions(setupQuestionSet(facts), &facts, reader, stdout); err != nil {
		return facts, err
	}
	return facts, nil
}

// askQuestions asks each question in qs in order, applying its answer to facts before
// moving to the next — the one-at-a-time sequencing story 3 asks for, and the seam the
// unit tests drive directly with a synthetic question list.
func askQuestions(qs []setupQuestion, facts *setupFacts, reader *bufio.Reader, stdout io.Writer) error {
	for _, q := range qs {
		fmt.Fprintln(stdout, q.prompt)
		for i, opt := range q.options {
			fmt.Fprintf(stdout, "  %d. %s\n", i+1, opt)
		}
		fmt.Fprintf(stdout, "select [1-%d]: ", len(q.options))
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("bench setup: could not read an answer to %q", q.prompt)
		}
		choice, perr := strconv.Atoi(strings.TrimSpace(line))
		if perr != nil || choice < 1 || choice > len(q.options) {
			return fmt.Errorf("bench setup: %q is not a valid answer to %q - expected a number from 1 to %d", strings.TrimSpace(line), q.prompt, len(q.options))
		}
		q.apply(choice-1, facts)
	}
	return nil
}
