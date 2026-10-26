package landing

import (
	"strings"

	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// refusalMessage renders the one refusal line every authorization caller prints. The
// literal prefix and the kind stay, because two tests and the operator's own memory read
// them; the sentence after names what the attribution means. It names no action, because
// the refusal face of the verb that prints it carries the route.
func refusalMessage(got authorization.Result) string {
	line := "prospective authorization refused: " + string(got.Kind)
	explanation := refusalExplanation(got.Kind)
	// An Infrastructure attribution carries the gate's own reason, which is more exact than
	// anything this renderer could state about the kind.
	if got.Reason != "" {
		explanation = got.Reason
	}
	if explanation != "" {
		line += " (" + explanation + ")"
	}
	if len(got.Suspects) != 0 {
		line += "; reader landings since the green marker: " + strings.Join(got.Suspects, ", ")
	}
	return line
}

// refusalExplanation answers what a refused kind means to the operator. The two lane
// outcomes and an infrastructure outcome answer nothing: a lane states its own outcome
// line already, and only the gate's own reason explains an infrastructure outcome.
func refusalExplanation(kind authorization.Kind) string {
	switch kind {
	case authorization.Inherited:
		return "the gate ran red on the composed tree and no green baseline attributes the red to this diff"
	case authorization.Candidate:
		return "the gate ran red on the composed tree and the green baseline attributes the red to this diff"
	}
	return ""
}
