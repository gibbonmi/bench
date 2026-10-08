package commitcmd

import (
	"github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/toon"
)

// admission starts, blocks, or unblocks one outcome. A refusal that raises no face of its own
// is cleared by a change to the active commitment, such as an outcome that the active
// milestone does not hold, so it is the reviewer's decision.
func admission(store repository.Store, c call) (string, int) {
	var err error
	switch c.form.name {
	case "start":
		err = store.Start(c.flags["--outcome"], c.flags["--request"], c.flags["--deliverable"])
	case "block":
		err = store.Block(c.flags["--outcome"], c.flags["--reason"])
	case "unblock":
		err = store.Unblock(c.flags["--outcome"])
	}
	if err != nil {
		return c.refuse(refusalroute.CommitmentDecision, err, nil)
	}
	out, err := toon.Table("commitment_admission", []string{"operation", "outcome"}, [][]string{{c.form.name, c.flags["--outcome"]}})
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	return out + "\n", 0
}
