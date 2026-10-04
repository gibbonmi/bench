package main

import (
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
)

func commitmentHelpRows(order int) []helpRow {
	return formHelpRows(order, commitcmd.HelpRows())
}

func commitmentCommand(args []string) (string, int) {
	return commitcmd.Command(boundaryRoot(), args)
}

// commitmentOutlook is the roadmap reader's commitment source: the same projection that
// status and the dashboard render.
func commitmentOutlook(root string) commitment.Outlook { return commitcmd.Outlook(root) }
