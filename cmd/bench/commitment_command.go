package main

import (
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
)

func commitmentHelpRows(order int) []helpRow {
	return formHelpRows(order, commitcmd.HelpRows())
}

// commitmentBound exempts only the reviewer's form from the bound, because that form
// prompts on the terminal before it reads the confirmation.
var commitmentBound = boundExemptForm(commitcmd.ReviewerVerb[1], boundReasonTerminal)

// commitmentCommand runs the commitment family over the caller's stdin and stdout, because
// the override form prompts on the terminal before it reads the confirmation.
func commitmentCommand(c Command, args []string) int {
	return commitcmd.Run(boundaryRoot(), args, c.Stdin, c.Stdout)
}

// commitmentOutlook is the roadmap reader's commitment source: the same projection that
// status and the dashboard render.
func commitmentOutlook(root string) commitment.Outlook { return commitcmd.Outlook(root) }
