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

// IsInvocable reports whether action is one command a session can type. Board actions use
// the phase, bench, and git command shapes below; prose and sequences are not commands.
func IsInvocable(action string) bool {
	if strings.Contains(action, StepSeparator) {
		return false
	}
	if strings.HasPrefix(action, harnessPrefix[HarnessClaude]) {
		command, argument, hasArgument := strings.Cut(action, " ")
		if command == harnessPrefix[HarnessClaude] || strings.Contains(strings.TrimPrefix(command, harnessPrefix[HarnessClaude]), "/") {
			return false
		}
		if !hasArgument {
			return true
		}
		return (strings.HasPrefix(argument, "specs/") && strings.HasSuffix(argument, "/spec.md")) ||
			(strings.HasPrefix(argument, "decisions/") && strings.HasSuffix(argument, ".md"))
	}
	if strings.HasPrefix(action, "git ") {
		return isGitCommand(strings.Fields(strings.TrimPrefix(action, "git ")))
	}
	if strings.HasPrefix(action, "bench ") {
		return isBenchCommand(strings.Fields(strings.TrimPrefix(action, "bench ")))
	}
	return false
}

func isGitCommand(parts []string) bool {
	if len(parts) == 1 {
		return parts[0] == "push" || parts[0] == "status"
	}
	return len(parts) == 2 && parts[0] == "worktree" && parts[1] == "list"
}

func isBenchCommand(parts []string) bool {
	if len(parts) == 1 {
		switch parts[0] {
		case "gate", "handoff", "link", "maps", "roadmap", "setup", "structure":
			return true
		}
		return false
	}
	if len(parts) == 2 {
		return (parts[0] == "gate" && parts[1] == "--fresh") ||
			(parts[0] == "status" && parts[1] == "--all") ||
			(parts[0] == "worktree" && parts[1] == "list")
	}
	return len(parts) == 3 &&
		((parts[0] == "spec" && parts[1] == "retire") ||
			(parts[0] == "worktree" && parts[1] == "clean"))
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
			if IsInvocable(runnerUp.Action) {
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
	action := "/bench-drain"
	if _, err := os.Stat(filepath.Join(root, "ROADMAP.md")); err == nil {
		action = "bench roadmap"
	}
	route.Lead = translateSignal(Signal{Name: "clean", Detail: "nothing pending", Action: action}, harness)
	return route
}

func firstInvocable(signals []Signal) (index int, signal Signal, ok bool) {
	for index, signal = range signals {
		if IsInvocable(signal.Action) {
			return index, signal, true
		}
	}
	return 0, Signal{}, false
}

func translateSignal(signal Signal, harness string) Signal {
	prefix, ok := harnessPrefix[harness]
	if !ok {
		return signal
	}
	signal.Action = strings.Replace(signal.Action, harnessPrefix[HarnessClaude], prefix, 1)
	return signal
}
