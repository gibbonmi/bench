package shellcommand

import (
	"path/filepath"
	"regexp"
	"slices"
)

// shellWrappers are the shell basenames whose `-c` operand is a child command string.
var shellWrappers = []string{"sh", "bash", "zsh"}

var wrapperCFlagRe = regexp.MustCompile(`^-[A-Za-z]*c[A-Za-z]*$`)

// WrapperChild is one shell `-c` child: the argv index of its operand and that operand parsed.
type WrapperChild struct {
	ArgumentIndex int
	Stream        Stream
}

// ShellChildren returns the `-c` children of the shell at words[executable], in argv order.
// It reads each argv position independently and does not recurse into a child.
// An invalid index or a flag without an operand yields no child.
func ShellChildren(words []string, executable int) []WrapperChild {
	if executable < 0 || executable >= len(words) || !slices.Contains(shellWrappers, filepath.Base(words[executable])) {
		return nil
	}
	var children []WrapperChild
	for k := executable + 1; k+1 < len(words); k++ {
		if wrapperCFlagRe.MatchString(words[k]) {
			children = append(children, WrapperChild{ArgumentIndex: k + 1, Stream: Parse(words[k+1])})
		}
	}
	return children
}

// IsBenchName reports whether a word names Bench by its basename: the installed CLI
// `bench` or the repository wrapper `bench.sh`. It is the one name test for Bench, so the
// guards and each caller of WorktreeExec accept the same executables.
func IsBenchName(word string) bool {
	switch filepath.Base(word) {
	case "bench", "bench.sh":
		return true
	}
	return false
}

// ExecProjection is the child that a `worktree exec` call runs in another checkout.
// Head marks the adjacent `worktree exec` words, and HasSeparator marks a later `--`.
// Child is the argv after that `--`, and DifferentCheckout is true exactly when Child is nonempty.
type ExecProjection struct {
	Head              bool
	HasSeparator      bool
	Child             []string
	DifferentCheckout bool
}

// WorktreeExec projects the child of `<executable> worktree exec ... -- <argv>`.
// The caller has already accepted words[executable] as Bench, so this function does not test it.
// The child argv is returned unchanged, with no joining and no reparsing.
func WorktreeExec(words []string, executable int) ExecProjection {
	if executable < 0 || executable+2 >= len(words) || words[executable+1] != "worktree" || words[executable+2] != "exec" {
		return ExecProjection{}
	}
	projection := ExecProjection{Head: true}
	rest := words[executable+3:]
	if separator := slices.Index(rest, "--"); separator >= 0 {
		projection.HasSeparator = true
		projection.Child = rest[separator+1:]
		projection.DifferentCheckout = len(projection.Child) > 0
	}
	return projection
}
