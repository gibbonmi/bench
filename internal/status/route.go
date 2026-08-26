package status

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/harnesses"
)

// HarnessClaude is the canonical form in which board phase actions are written. The record
// owns the name, so the constant reads out of the row rather than restating it.
var HarnessClaude = claudeName()

// harnessPrefix maps every recorded harness to its phase invocation form. The record owns
// both the names and the forms, so a new row reaches the grammar without an edit here. A
// row with no phase surface holds an empty form.
var harnessPrefix = buildHarnessPrefix()

// claudeName returns the record's canonical harness: the row whose phase form the action
// table already writes its phase commands in. Both halves are read rather than restated, so
// neither the harness name nor the form appears here as a literal.
func claudeName() string {
	for _, definition := range actionDefinitions {
		if definition.kind != actionPhase || definition.command == "" {
			continue
		}
		for _, row := range harnesses.Rows {
			if row.PhaseForm != "" && strings.HasPrefix(definition.command, row.PhaseForm) {
				return row.Harness
			}
		}
	}
	return ""
}

func buildHarnessPrefix() map[string]string {
	table := make(map[string]string, len(harnesses.Rows))
	for _, row := range harnesses.Rows {
		table[row.Harness] = row.PhaseForm
	}
	return table
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

// ValidHarness reports whether the record holds harness. A recorded harness with no phase
// surface is still a valid choice, because the grammar names every row.
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
	// A row with no phase form leaves the action untranslated, which is what an unrecorded
	// harness has always done.
	prefix := harnessPrefix[harness]
	if prefix == "" {
		return signal
	}
	signal.Action = strings.Replace(signal.Action, harnessPrefix[HarnessClaude], prefix, 1)
	return signal
}
