package structure

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// Fact is one typed structure observation shared by query consumers.
type Fact struct {
	Kind, Path    string
	Actual, Limit int
	State, Detail string
}

// Facts returns the whole-tree engine's typed observations.
func Facts(root string) ([]Fact, error) {
	out, err := git.Output("-C", root, "ls-files")
	if err != nil {
		return nil, err
	}
	_, _, facts := evaluate(root, strings.Split(out, "\n"), true)
	return facts, nil
}

// checkAll is the human/count projection of the same whole-tree evaluation.
func checkAll(root string) (report string, violations int, err error) {
	out, err := git.Output("-C", root, "ls-files")
	if err != nil {
		return "", 0, err
	}
	report, violations, _ = evaluate(root, strings.Split(out, "\n"), true)
	return report, violations, nil
}
