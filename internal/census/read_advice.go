package census

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// Read is one eligible raw read that a native file reader can replace: the resolved head
// of the raw call, the file operand, and the physical-line range the call prints. First and
// Last are zero for a read of the whole file.
type Read struct {
	Head        string
	Path        string
	First, Last int
}

// Ranged reports whether the read prints one line range rather than the whole file.
func (r Read) Ranged() bool { return r.First > 0 }

// headDefaultLines is the line count a head call prints when it names no count.
const headDefaultLines = 10

// unresolvedOperand holds the characters that leave an operand unresolved before the shell
// runs: an expansion, a substitution, a glob, a brace list, or a home prefix.
const unresolvedOperand = "$`*?[]{}~"

// sedLinePrint is the one sed program shape a read can print: one line number, or two line
// numbers. The exact-reader parsers decide whether the numbers name one ordered range.
var sedLinePrint = regexp.MustCompile(`^([0-9]+)(?:,([0-9]+))?p$`)

// readForms observes the operands of each eligible read head. Each answers the file and the
// range, and refuses every flag, operand count, and program it does not name.
var readForms = map[string]func([]string) (Read, bool){
	"cat":  catRead,
	"head": headRead,
	"tail": tailRead,
	"sed":  sedRead,
}

// ObserveRead answers the eligible raw read that command makes. The command is exactly one
// simple command of plain words, with no control operator and no redirection, so a stdin
// read, a heredoc, a here-string, a pipeline, and a chain are not reads. The shellcommand
// parse folds the quotes, so text inside an example's quotes is one operand of another
// head. The routine prefix resolves the head, and an xargs prefix reads its operands from
// stdin. An unknown parse and an unresolved operand answer no read.
func ObserveRead(command string) (Read, bool) {
	stream := shellcommand.Parse(command)
	if stream.Unlexed || len(stream.Commands) != 1 {
		return Read{}, false
	}
	for _, token := range stream.Tokens {
		if token.Kind != shellcommand.Word {
			return Read{}, false
		}
	}
	words := shellcommand.ProjectCommandWords(stream.Tokens)
	prefix := shellcommand.ResolveRoutinePrefix(words)
	if !prefix.Executes || prefix.ViaXargs || prefix.Index >= len(words) {
		return Read{}, false
	}
	observe, ok := readForms[filepath.Base(words[prefix.Index])]
	if !ok {
		return Read{}, false
	}
	read, ok := observe(words[prefix.Index+1:])
	if !ok || !resolved(read.Path) {
		return Read{}, false
	}
	read.Head = resolvedHead(words)
	return read, true
}

// resolved reports whether an operand names one file as written: not empty, not a flag or
// stdin, with no control character and no character that the shell expands.
func resolved(operand string) bool {
	if operand == "" || strings.HasPrefix(operand, "-") || strings.ContainsAny(operand, unresolvedOperand) {
		return false
	}
	return sanitize.LineSafe(operand)
}

// catRead is a cat call of one named file and no flag.
func catRead(operands []string) (Read, bool) {
	if len(operands) != 1 {
		return Read{}, false
	}
	return Read{Path: operands[0]}, true
}

// headRead is a head call of one named file, with no count or one positive line count. It
// prints the first lines, so it reads one range from line one.
func headRead(operands []string) (Read, bool) {
	count, path, ok := countedOperand(operands)
	if !ok {
		return Read{}, false
	}
	return Read{Path: path, First: 1, Last: count}, true
}

// tailRead is a tail call of one named file, with no count or one positive line count. Its
// range counts from the end, which a forward line range cannot state, so it reads the file.
func tailRead(operands []string) (Read, bool) {
	_, path, ok := countedOperand(operands)
	return Read{Path: path}, ok
}

// countedOperand reads the operands `<file>`, `-n <count> <file>`, or `-n<count> <file>`. The
// count is a positive decimal in the exact-reader grammar; a signed count changes what the
// call prints.
func countedOperand(operands []string) (int, string, bool) {
	text := ""
	switch {
	case len(operands) == 1:
		return headDefaultLines, operands[0], true
	case len(operands) == 2 && strings.HasPrefix(operands[0], "-n") && len(operands[0]) > 2:
		text = operands[0][2:]
	case len(operands) == 3 && operands[0] == "-n":
		text, operands = operands[1], operands[1:]
	default:
		return 0, "", false
	}
	count, err := responsebound.ParsePositive(text)
	return count, operands[1], err == nil
}

// sedRead is a quiet sed call of one line-print program and one named file. Each bound is a
// positive exact-reader decimal, and the exact-reader line-range parser grades their order.
func sedRead(operands []string) (Read, bool) {
	if len(operands) != 3 || operands[0] != "-n" {
		return Read{}, false
	}
	match := sedLinePrint.FindStringSubmatch(operands[1])
	if match == nil {
		return Read{}, false
	}
	lastText := match[2]
	if lastText == "" {
		lastText = match[1]
	}
	first, firstErr := responsebound.ParsePositive(match[1])
	last, lastErr := responsebound.ParsePositive(lastText)
	if firstErr != nil || lastErr != nil {
		return Read{}, false
	}
	first, last, err := responsebound.ParseLineRange(responsebound.LineRange(first, last))
	if err != nil {
		return Read{}, false
	}
	return Read{Path: operands[2], First: first, Last: last}, true
}

// RecordRead appends at most one raw-call record for one invocation. The pool-path
// attribution of Record writes first. When it writes nothing, an eligible read under the
// verified assignment that the caller supplies writes one record under that assignment.
// An empty assignment is an unowned read, and it writes nothing.
func RecordRead(command, root, home, assignment string, now time.Time) error {
	return recordReadWith(command, root, home, assignment, now, benchguard.DefaultResolver())
}

func recordReadWith(command, root, home, assignment string, now time.Time, resolver benchguard.Resolver) error {
	id, head, ok := poolAttribution(command, root, home, resolver)
	if !ok && assignment != "" {
		var read Read
		if read, ok = ObserveRead(command); ok {
			if err := checkAssignmentID(assignment); err != nil {
				return err
			}
			id, head = assignment, read.Head
		}
	}
	if !ok {
		return nil
	}
	return write(Dir(home, root), id, composeRecord(now.UTC().Format(time.RFC3339), head))
}
