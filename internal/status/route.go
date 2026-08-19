package status

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HarnessClaude is the canonical form in which board phase actions are written.
const HarnessClaude = "claude"

var harnessPrefix = map[string]string{
	HarnessClaude: "/bench-",
	"codex":       "$bench-",
}

// HarnessChoices returns the accepted harness names, with the canonical form first.
func HarnessChoices() string {
	rest := make([]string, 0, len(harnessPrefix)-1)
	for name := range harnessPrefix {
		if name != HarnessClaude {
			rest = append(rest, name)
		}
	}
	sort.Strings(rest)
	return strings.Join(append([]string{HarnessClaude}, rest...), "|")
}

// ValidHarness reports whether harness has a phase-invocation form.
func ValidHarness(harness string) bool {
	_, ok := harnessPrefix[harness]
	return ok
}

// RouteResult is the next command selected from the severity ladder for one harness.
// NoCommand distinguishes a board carrying only non-invocable signals from an empty board.
type RouteResult struct {
	Lead      Signal
	RunnersUp []Signal
	NoCommand bool
}

// Route selects the first invocable signal in ladder order and translates its phase action
// for harness. It does not re-rank the board; later invocable signals remain runners-up.
func Route(signals []Signal, harness string) RouteResult {
	if i, signal, ok := firstInvocable(signals); ok {
		route := RouteResult{Lead: translateSignal(signal, harness)}
		for _, runnerUp := range signals[i+1:] {
			if runnerUp.invocable() {
				route.RunnersUp = append(route.RunnersUp, translateSignal(runnerUp, harness))
			}
		}
		return route
	}
	if len(signals) == 0 {
		return RouteResult{}
	}
	return RouteResult{Lead: translateSignal(signals[0], harness), NoCommand: true}
}

// RouteFor selects the route for root's board. A clean board is the only case without a
// signal, so its fallback belongs beside selection rather than in each caller.
func RouteFor(root string, signals []Signal, harness string) RouteResult {
	route := Route(signals, harness)
	if len(signals) != 0 {
		return route
	}
	command := commandAction(drainPhaseAction)
	if _, err := os.Stat(filepath.Join(root, "ROADMAP.md")); err == nil {
		command = commandAction(roadmapAction)
	}
	route.Lead = translateSignal(newSignal(0, "clean", "nothing pending", command), harness)
	return route
}

func firstInvocable(signals []Signal) (index int, signal Signal, ok bool) {
	for index, signal = range signals {
		if signal.invocable() {
			return index, signal, true
		}
	}
	return 0, Signal{}, false
}

func translateSignal(signal Signal, harness string) Signal {
	if signal.actionID.kind() != actionPhase {
		return signal
	}
	prefix, ok := harnessPrefix[harness]
	if !ok {
		return signal
	}
	signal.Action = strings.Replace(signal.Action, harnessPrefix[HarnessClaude], prefix, 1)
	return signal
}
