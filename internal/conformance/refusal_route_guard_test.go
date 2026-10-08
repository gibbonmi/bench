package conformance

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// routeSamplePool is the pool the guard check names to benchguard.PoolReference, and
// routeSample is the value it fills every route slot with. The sample is a path under that
// pool, so a slot in a `cd` or a `git -C` position reaches the pool and turns red.
const routeSamplePool = "/bench-home/worktrees/pool"

func routeSample(string) string { return routeSamplePool + "/repository/assignment" }

// denyingChecker answers each repository fact the way that denies, so a git step the
// guard can allow only under one repository state reds here under every state.
var denyingChecker = gitguard.Checker{
	RefResolves:  func(string) bool { return false },
	BranchExists: func(string) bool { return true },
}

// agentRouteGuardFaults renders each command step of each agent face with every slot
// filled by the sample and grades it through the wired guards. It returns the count of
// graded steps, so a check that grades nothing cannot pass, and one fault per refusal. A
// reviewer route is the reviewer's to run, so the check does not grade it.
func agentRouteGuardFaults(faces []refusalroute.Face) (graded int, faults []string) {
	for _, face := range faces {
		if face.Authority != refusalroute.Agent {
			continue
		}
		for index, step := range face.Route {
			if !step.IsCommand() {
				continue
			}
			graded++
			command := step.Filled(routeSample)
			for _, fault := range stepGuardFaults(command) {
				faults = append(faults, fmt.Sprintf("face %s step %d %q: %s", face.Name, index+1, command, fault))
			}
		}
	}
	return graded, faults
}

// stepGuardFaults names each guard that refuses one command step.
func stepGuardFaults(command string) []string {
	var faults []string
	if label := gitguard.Classify(command, denyingChecker); label != "" {
		faults = append(faults, "gitguard denies "+label)
	}
	if target := benchguard.PoolReference(command, routeSamplePool); target != "" {
		faults = append(faults, "benchguard names the pool path "+target)
	}
	if verdict := benchguard.Classify(command, benchguard.DefaultResolver()); verdict.Blocked {
		faults = append(faults, "benchguard blocks: "+verdict.Message())
	}
	if runsBenchUnderExec(command) {
		faults = append(faults, "runs a Bench child through bench worktree exec, which FT341 refuses")
	}
	return faults
}

// runsBenchUnderExec reports whether a simple command of command is `bench worktree exec
// <target> -- <child>` whose child invokes Bench. Neither wired guard refuses that form,
// so the check reads it itself.
func runsBenchUnderExec(command string) bool {
	resolver := benchguard.DefaultResolver()
	stream := shellcommand.Parse(command)
	for _, span := range stream.Commands {
		words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		prefix := shellcommand.ResolveRoutinePrefix(words)
		if !prefix.Executes || prefix.Index >= len(words) {
			continue
		}
		head := words[prefix.Index:]
		if len(head) < 3 || !benchguard.InvokesBench(sanitize.ShellQuote(head[0]), resolver) || head[1] != "worktree" || head[2] != "exec" {
			continue
		}
		separator := slices.Index(head, "--")
		if separator < 0 {
			continue
		}
		child := make([]string, 0, len(head)-separator-1)
		for _, word := range head[separator+1:] {
			child = append(child, sanitize.ShellQuote(word))
		}
		if benchguard.InvokesBench(strings.Join(child, " "), resolver) {
			return true
		}
	}
	return false
}

// TestAgentRoutesPassTheWiredGuards is RR08 through RR11. Every command step of every
// agent route in the registry, with every slot filled, must pass the destructive-git
// guard, the pool guard, and the follow-on guard, and must run no Bench child under exec.
func TestAgentRoutesPassTheWiredGuards(t *testing.T) {
	graded, faults := agentRouteGuardFaults(refusalroute.Inventory())
	if graded == 0 {
		t.Fatal("the registry holds no agent command step; the guard check grades nothing")
	}
	for _, fault := range faults {
		t.Error(fault)
	}
}

// TestAgentRouteGuardCheckBites is RR12, with one more injected step for each guard and for
// the exec rule, so a check that drops one classifier reds here. The same step under
// reviewer authority passes, because the check grades agent routes only.
func TestAgentRouteGuardCheckBites(t *testing.T) {
	for _, tc := range []struct {
		name string
		step refusalroute.Step
		want string
	}{
		{"a raw merge", refusalroute.Command(refusalroute.Text("git merge"), refusalroute.Fact("commit")), "gitguard denies git merge"},
		{"a chained Bench call", refusalroute.Command(refusalroute.Text("bench doctor && bench gate")), "benchguard blocks"},
		{"a cd into the pool", refusalroute.Command(refusalroute.Text("cd"), refusalroute.Fact("checkout")), "benchguard names the pool path"},
		{"a Bench child under exec", refusalroute.Command(refusalroute.Text("bench worktree exec"), refusalroute.Fact("id"), refusalroute.Text("-- bench commit -m"), refusalroute.Operator("msg")), "FT341"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			face := refusalroute.Face{Verb: refusalroute.Land, Name: "injected", Authority: refusalroute.Agent, Route: []refusalroute.Step{tc.step}}
			if _, faults := agentRouteGuardFaults([]refusalroute.Face{face}); !anyContains(faults, tc.want) {
				t.Fatalf("injected agent step faults = %q, want one that names %q", faults, tc.want)
			}
			face.Authority = refusalroute.Reviewer
			if graded, faults := agentRouteGuardFaults([]refusalroute.Face{face}); graded != 0 || len(faults) != 0 {
				t.Fatalf("reviewer step graded %d with faults %q, want no grade", graded, faults)
			}
		})
	}
}
