package commitcmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

// Run runs one commitment operation in root over the caller's stdin and stdout. The
// override form prints its confirmation prompt on stdout before it reads stdin.
func Run(root string, args []string, stdin io.Reader, stdout io.Writer) int {
	return run(root, args, console{stdin: stdin, stdout: stdout, isTerminal: terminal.IsTerminal})
}

// console is the caller's stdin and stdout, and the predicate that tests stdin for a terminal.
type console struct {
	stdin      io.Reader
	stdout     io.Writer
	isTerminal func(io.Reader) bool
}

func run(root string, args []string, con console) int {
	out, code := dispatch(root, args, con)
	fmt.Fprint(con.stdout, out)
	return code
}

// overrideForm is the commitment form that only the reviewer runs.
const overrideForm = "override"

// ReviewerVerb is the family and the form of the verb that only the reviewer runs. The
// follow-on guard denies it in an agent's tool call.
var ReviewerVerb = [2]string{"commitment", overrideForm}

// overrideTable and overrideFields are the schema of the row that the reviewer confirms, and
// overridePrompt asks for the typed confirmation.
const (
	overrideTable  = "override"
	overridePrompt = "type the assignment id to confirm:"
)

var overrideFields = []string{"assignment", "label", "builds", "reason"}

// override grants the reviewer's override to one assignment. Each layer stops an honest
// mistake of an agent: the primary checkout, a terminal on stdin, and the typed id. The
// row and the prompt reach stdout before the read, so the reviewer sees what they confirm.
// Every refusal is the reviewer's to clear.
func override(store commitrepo.Store, c call, con console) (string, int) {
	refuse := func(err error) (string, int) { return c.refuse(refusalroute.CommitmentHandback, err, nil) }
	primary, err := git.IsPrimaryCheckout(store.Root)
	if err != nil {
		return refuse(err)
	}
	if !primary {
		return refuse(errors.New("commitment override runs from the primary checkout"))
	}
	if !con.isTerminal(con.stdin) {
		return refuse(errors.New("commitment override needs the reviewer at a terminal; run it in your own terminal"))
	}
	id, reason := c.flags["--assignment"], c.flags["--reason"]
	target, err := store.OverrideTarget(id)
	if err != nil {
		return refuse(err)
	}
	if err := commitment.ValidReason(reason); err != nil {
		return refuse(err)
	}
	row, err := toon.Table(overrideTable, overrideFields, [][]string{{id, target.Label, strings.Join(target.Builds, " "), reason}})
	if err != nil {
		return refuse(err)
	}
	fmt.Fprint(con.stdout, row+overridePrompt+"\n")
	// Only a whole line that equals the id confirms. A trimmed or prefix compare would accept
	// a wrong id, and EOF before the newline is no confirmation.
	if line, err := bufio.NewReader(con.stdin).ReadString('\n'); err != nil || line != id+"\n" {
		return refuse(errors.New("commitment override not confirmed"))
	}
	if err := store.Override(id, reason); err != nil {
		return refuse(err)
	}
	return "overridden{assignment=" + id + "}\n", 0
}
