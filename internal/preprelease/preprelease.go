// Package preprelease owns `bench prep-release`, the ship-tier rehearsal. `bench gate`
// answers the narrower question — does the kit work from this tree — and leaves the
// four-platform artifact matrix, the cross-compile matrix, the release-only package
// suites, the release preflight verify, and the ship-tier canary fixtures to a surface
// that runs once per release instead of once per commit. This is that surface, and its
// exit code carries the authority the dev gate gave up.
//
// It invents no machinery: every step is an existing script, the existing conformance
// suite at the ship tier, or the existing canary sweep. Ship green is exit 0 with
// evidence at dist/preflight/release-index.json and dist/artifacts.
//
// The dev-green precondition is an observer, not a second verdict-reuse grant: a green
// observation lets the run start, and every ship-tier check still runs in full.
package preprelease

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/releaseevidence"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// ArtifactsDir and IndexPath are the evidence a green run leaves behind, root-relative
// so a diagnostic names the path an agent would type. EvidenceDir is where the preflight
// promotes its per-phase records alongside that index.
const (
	ArtifactsDir = "dist/artifacts"
	EvidenceDir  = "dist/preflight"
	IndexPath    = EvidenceDir + "/release-index.json"
)

// preflightStep is the stage whose failure has promoted per-phase evidence to attribute.
const preflightStep = "release-preflight"

// requiredTools are the interpreters every ship-tier step reaches through. Resolving
// them up front is what turns a missing toolchain into a named diagnostic rather than
// an `exec: not found` surfacing from four levels inside a shell script.
var requiredTools = []string{"bash", "git", "go", "node"}

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:  "bench prep-release",
	Help: "usage: bench prep-release",
}

// Step is one ship-tier stage. Exactly one of Argv and Run is set: the scripted stages
// are subprocesses, and the canary sweep is a library call.
type Step struct {
	Name string
	Argv []string
	Env  []string
	Run  func(root string) error
}

// Steps is the ship-tier sequence for root, in run order; kit is the checkout that owns
// the conformance suite, which is root itself everywhere but a linked repo. The order is
// load-bearing: the preflight verify reads the artifact set the build stage promotes.
//
// This is the only enumeration of what `prep-release` runs. Anything that needs to know
// the sequence — a test, a diagnostic — reads it from here.
func Steps(root, kit string) []Step {
	return []Step{
		{
			Name: "artifacts",
			Argv: []string{"bash", filepath.Join(root, "scripts", "build-artifacts.sh"), root, filepath.Join(root, filepath.FromSlash(ArtifactsDir))},
		},
		{
			// -tags stress is what makes crossCompileMatrix more than a no-op; without
			// it this step runs a check that silently returns nil. The suite grades no
			// package suites of its own — the step below is what runs those.
			Name: "conformance-ship",
			Argv: append(goTestArgv(kit), "-tags", "stress", "./internal/conformance", "-run", "^TestRootConformance$"),
			Env:  []string{"BENCH_CONFORMANCE_ROOT=" + root, registry.ConformanceTierEnv + "=" + string(registry.Ship)},
		},
		{
			// The release-only package suites — internal/preflight,
			// internal/releaseevidence, internal/publication — are excluded from the
			// dev tier's package enumeration, and once the core test run became a gate
			// phase of its own nothing else reached them. This is the only surface that
			// still does, so without it ship green would cover three fewer suites than
			// its exit code claims.
			Name: "core-tests-ship",
			Argv: gate.GateGoArgv(kit, "test", root),
			Env:  []string{registry.ConformanceTierEnv + "=" + string(registry.Ship)},
		},
		{
			// Verify mode by design: its index is deliberately insufficient for publish
			// authority, so a rehearsal can never be mistaken for the release's own
			// preflight.
			Name: preflightStep,
			Argv: []string{"bash", filepath.Join(root, "scripts", "release-preflight.sh"), "--mode", "verify"},
		},
		{
			Name: "canary-ship",
			Run:  canary.SweepShip,
		},
	}
}

// goTestArgv mirrors the gate's conformance invocation: -C so the suite compiles in the
// kit checkout while grading a root that may be elsewhere, and -count=1 so a cached
// verdict can never stand in for the release rehearsal.
func goTestArgv(kit string) []string {
	argv := []string{"go"}
	if kit != "" {
		argv = append(argv, "-C", kit)
	}
	return append(argv, "test", "-count=1")
}

// Command runs the ship-tier rehearsal. `help`, `--help`, and `-h` print the declared
// help at exit 0; usage errors exit 2; a missing precondition, a missing tool, or any
// red step exits 1. Step output streams to stderr as progress, leaving stdout for the
// single line that reports ship green.
func Command(args []string, stdout, stderr io.Writer) int {
	if _, line, code := usage.Parse(grammar, args); line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
			return 0
		}
		fmt.Fprintln(stderr, line)
		return code
	}

	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	if tool := firstMissingTool(); tool != "" {
		fmt.Fprintln(stderr, "prep-release: required tool is missing or not executable: "+tool)
		return 1
	}
	if inspection := gate.Inspect(root); !inspection.ReusableGreen {
		fmt.Fprintln(stderr, Refusal(inspection))
		return 1
	}

	kit := os.Getenv("BENCH_KIT")
	if kit == "" {
		kit = root
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for _, step := range Steps(root, kit) {
		fmt.Fprintf(stderr, "prep-release: %s\n", step.Name)
		if err := runStep(ctx, root, step, stderr); err != nil {
			fmt.Fprintf(stderr, "prep-release: %s failed: %v\n", step.Name, err)
			if step.Name == preflightStep {
				for _, attribution := range PreflightAttributions(root) {
					fmt.Fprintf(stderr, "prep-release: %s\n", attribution)
				}
			}
			return 1
		}
	}
	fmt.Fprintf(stdout, "prep-release: ship green (evidence at %s and %s)\n", IndexPath, ArtifactsDir)
	return 0
}

// Refusal is the one diagnostic for a tree with no current dev-green verdict. The cause
// is the gate's own classification of its cache, so the four rejected states — absent,
// red, bound to another tree, expired — each report themselves without this package
// re-deriving any of them.
func Refusal(inspection gate.Inspection) string {
	cause := inspection.Reason
	if cause == "" {
		cause = string(inspection.State)
	}
	return fmt.Sprintf("prep-release: no current dev-green verdict for this tree (%s) — run `bench gate`, then re-run prep-release", cause)
}

// PreflightAttributions names each verify phase whose promoted record is not green, with
// that phase's own diagnostic. The preflight streams twenty minutes of phase output and
// then exits 1, so its bare exit status leaves the cause somewhere in the scrollback;
// the records it promoted already hold the attribution, and this reads it back rather
// than re-deriving anything from that output.
//
// A record that is absent or undecodable yields no line: the phase sequence is the
// preflight's to own, and a run interrupted before promotion has nothing to attribute.
func PreflightAttributions(root string) []string {
	var attributions []string
	for _, phase := range releaseevidence.PhaseNames(releaseevidence.ModeVerify) {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(EvidenceDir), phase+".json"))
		if err != nil {
			continue
		}
		var record releaseevidence.Record
		if json.Unmarshal(data, &record) != nil || record.Status == releaseevidence.StatusGreen {
			continue
		}
		attribution := preflightStep + " phase " + phase + " is " + string(record.Status)
		if record.Error != nil && record.Error.Message != "" {
			attribution += ": " + record.Error.Message
		}
		attributions = append(attributions, attribution)
	}
	return attributions
}

func firstMissingTool() string {
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			return tool
		}
	}
	return ""
}

// runStep executes one stage from root. The subprocess inherits the ambient environment
// so the toolchain and PATH the maintainer resolved are the ones the scripts see, and it
// dies with the run: a SIGINT partway through has to leave the evidence directory in the
// state its own atomic promotion last left it, never half-written by a surviving child.
func runStep(ctx context.Context, root string, step Step, stderr io.Writer) error {
	if step.Run != nil {
		return step.Run(root)
	}
	cmd := exec.CommandContext(ctx, step.Argv[0], step.Argv[1:]...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), step.Env...)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM) }
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", strings.Join(step.Argv, " "), err)
	}
	return nil
}
