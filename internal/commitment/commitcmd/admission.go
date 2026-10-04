package commitcmd

import (
	"github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/toon"
)

func admission(store repository.Store, operation string, flags map[string]string) (string, int) {
	var err error
	switch operation {
	case "start":
		err = store.Start(flags["--outcome"], flags["--request"], flags["--deliverable"])
	case "block":
		err = store.Block(flags["--outcome"], flags["--reason"])
	case "unblock":
		err = store.Unblock(flags["--outcome"])
	}
	if err != nil {
		return refusal(operation, err)
	}
	out, err := toon.Table("commitment_admission", []string{"operation", "outcome"}, [][]string{{operation, flags["--outcome"]}})
	if err != nil {
		return refusal(operation, err)
	}
	return out + "\n", 0
}
