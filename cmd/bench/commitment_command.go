package main

import (
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
)

func commitmentHelpRows(order int) []helpRow {
	return formHelpRows(order, commitcmd.HelpRows())
}

// commitmentCommand runs the commitment family over the caller's stdin and stdout, because
// the override form prompts on the terminal before it reads the confirmation.
func commitmentCommand(c Command, args []string) int {
	return commitcmd.Run(boundaryRoot(), args, c.Stdin, c.Stdout)
}

// commitmentOutlook is the roadmap reader's commitment source: the same projection that
// status and the dashboard render.
func commitmentOutlook(root string) commitment.Outlook { return commitcmd.Outlook(root) }
