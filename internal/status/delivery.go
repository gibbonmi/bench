package status

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/spec"
)

// deliverySeverity places the delivery signal on the ladder: after the gate, Git, and
// worktree rows, beside the drain.
const deliverySeverity = 4

// appendDelivery projects the delivery next action. A published commitment owns that
// action, so the board renders the shared outlook and names no staged spec beside it.
// Before adoption no commitment exists, and the staged-spec row stays the planning signal.
func appendDelivery(rows []row, root string) []row {
	return appendDeliveryOutlook(rows, root, commitcmd.Outlook(root))
}

func appendDeliveryOutlook(rows []row, root string, outlook commitment.Outlook) []row {
	if outlook.State == commitment.OutlookAdoptionRequired {
		return appendStagedSpecs(rows, root)
	}
	return append(rows, row{deliverySeverity, "commitment", outlookDetail(outlook), outlookAction(outlook)})
}

// outlookAction is the outlook's command as a typed board action. An outlook with no
// command, such as active work that continues, gives the empty advisory action.
func outlookAction(outlook commitment.Outlook) statusAction {
	if outlook.Command == "" {
		return advisoryAction("")
	}
	argument, ok := strings.CutPrefix(outlook.Command, actionDefinitions[commitmentAction].command+" ")
	if !ok {
		return advisoryAction(outlook.Command)
	}
	return commandActionWithArgument(commitmentAction, argument)
}

// outlookDetail states the outlook in one board cell: the state, the outcome that the
// state names, the milestone, and the blocked outcomes.
func outlookDetail(outlook commitment.Outlook) string {
	detail := outlook.State
	switch {
	case outlook.Next != "":
		detail += " " + outlook.Next
	case len(outlook.Active) != 0:
		detail += " " + strings.Join(outlook.Active, " ")
	}
	if outlook.Milestone != "" {
		detail += " in " + outlook.Milestone
	}
	if blocked := outlook.BlockedOutcomes(); len(blocked) != 0 {
		detail += "; blocked " + strings.Join(blocked, " ")
	}
	return detail
}

func appendStagedSpecs(rows []row, root string) []row {
	n, slug := stagedSpecCount(root)
	if n == 0 {
		return rows
	}
	command := commandAction(implementSpecPhaseAction)
	if n == 1 {
		command = commandActionWithArgument(implementSpecPhaseAction, "specs/"+slug+"/spec.md")
	}
	return append(rows, row{deliverySeverity, "specs", fmt.Sprintf("%d staged spec(s)", n), command})
}

func stagedSpecCount(root string) (int, string) {
	facts, err := spec.Facts(root)
	if err != nil {
		return 0, ""
	}
	n, slug := 0, ""
	for _, fact := range facts {
		if fact.Status != "staged" {
			continue
		}
		n++
		slug = fact.Slug
	}
	return n, slug
}
