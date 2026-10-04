package main

import "github.com/gibbonmi/bench/internal/commitment/commitcmd"

func commitmentHelpRows(order int) []helpRow {
	return formHelpRows(order, commitcmd.HelpRows())
}

func commitmentCommand(args []string) (string, int) {
	return commitcmd.Command(boundaryRoot(), args)
}
