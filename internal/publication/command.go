package publication

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
)

const usageLine = "usage: bench release prepare|submit|promote|rollback|status --version <version> [--profile public|bank] [--root <dir>] [--registry <base-url>] [--path first|staged] [--adapter npm|fixture] [--provenance] [--message <text>]"

// The two registry adapters --adapter selects between. The default is the
// hermetic fixture: reaching the real registry is always an explicit opt-in.
const (
	adapterFixture = "fixture"
	adapterNPM     = "npm"
)

// Command is the `bench release <prepare|submit|promote|rollback|status>`
// entry point. It is idempotent and non-interactive. prepare only verifies
// the approved release directory locally. submit runs (or resumes) the
// first-publication or staged-submission state machine depending on --path.
//
// promote moves the "latest" dist-tag once the complete set reverifies live.
// rollback removes candidate tags and deprecates a bad version. status
// reports the durable record without touching the registry.
//
// TOON tables go to stdout, and progress goes to stderr; structured errors go
// to stdout. Exit 0 means success or no-op, 1 means unsatisfied release
// intent, 2 means usage.
func Command(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stdout, toon.Usage("bench release", "<missing subcommand>"))
		return 2
	}
	sub := args[0]
	parsed, usageErr := parseCommandArgs(args[1:])
	if usageErr != nil {
		fmt.Fprintln(stdout, usageErr.Error())
		return 2
	}
	if parsed.root == "" {
		root, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stdout, toon.Errorf("input", "could not resolve the working directory"))
			return 1
		}
		parsed.root = root
	}
	parsed.root, _ = filepath.Abs(parsed.root)

	switch sub {
	case "prepare":
		return runPrepare(parsed.root, parsed.version, stdout, stderr)
	case "submit":
		return runSubmit(parsed, stdout, stderr)
	case "promote":
		return runPromote(parsed, stdout, stderr)
	case "rollback":
		return runRollback(parsed, stdout, stderr)
	case "status":
		return runStatus(parsed.root, stdout)
	default:
		fmt.Fprintln(stdout, toon.Usage("bench release", sub))
		return 2
	}
}

// releaseArgs is one parsed `bench release` invocation. adapter is always
// resolved: it defaults to the fixture, with no environment twin. A real
// publish is spelled out on the command line or it does not happen.
type releaseArgs struct {
	root         string
	version      string
	profile      string
	registryBase string
	path         string
	message      string
	adapter      string
	provenance   bool
}

func parseCommandArgs(args []string) (releaseArgs, error) {
	parsed := releaseArgs{adapter: adapterFixture}
	for i := 0; i < len(args); i++ {
		flag := args[i]
		if flag == "--provenance" {
			parsed.provenance = true
			continue
		}
		var target *string
		switch flag {
		case "--root":
			target = &parsed.root
		case "--version":
			target = &parsed.version
		case "--profile":
			target = &parsed.profile
		case "--registry":
			target = &parsed.registryBase
		case "--path":
			target = &parsed.path
		case "--message":
			target = &parsed.message
		case "--adapter":
			target = &parsed.adapter
		default:
			return releaseArgs{}, fmt.Errorf("%s", toon.Usage("bench release", flag))
		}
		i++
		if i >= len(args) {
			return releaseArgs{}, fmt.Errorf("%s", toon.MissingArg("bench release", flag+" value"))
		}
		*target = args[i]
	}
	if parsed.profile != "" && parsed.profile != "public" && parsed.profile != "bank" {
		return releaseArgs{}, fmt.Errorf("%s", toon.Usage("bench release", "--profile "+parsed.profile))
	}
	if parsed.path != "" && parsed.path != "first" && parsed.path != "staged" {
		return releaseArgs{}, fmt.Errorf("%s", toon.Usage("bench release", "--path "+parsed.path))
	}
	if parsed.adapter != adapterFixture && parsed.adapter != adapterNPM {
		return releaseArgs{}, fmt.Errorf("%s", toon.Usage("bench release", "--adapter "+parsed.adapter))
	}
	return parsed, nil
}

// resolveRegistryBase applies the BENCH_RELEASE_REGISTRY fallback once. So
// submit, promote, and rollback share one answer for which registry base this
// invocation addresses: the same value the selected adapter is built from.
func (a *releaseArgs) resolveRegistryBase() string {
	if a.registryBase == "" {
		a.registryBase = os.Getenv("BENCH_RELEASE_REGISTRY")
	}
	return a.registryBase
}

// newRegistry resolves the one adapter selection into the registry the state
// machine drives. Selection is a CLI concern: every state-machine entry point
// takes the Registry port and never learns which adapter it got. Credentials
// stay ambient — the npm adapter reads the npm config/environment itself, and
// nothing here ever holds or records a token.
func newRegistry(args releaseArgs) Registry {
	if args.adapter != adapterNPM {
		return NewFixtureRegistry(args.registryBase)
	}
	registry := NewNPMCLIRegistry(args.registryBase)
	registry.Provenance = args.provenance
	if args.profile == "public" {
		registry.Access = "public"
	}
	return registry
}

func runPrepare(root, version string, stdout, stderr io.Writer) int {
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release prepare", "--version"))
		return 2
	}
	releaseIndexSHA256, packages, err := VerifyApprovedSet(root, version)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
		return 1
	}
	rows := make([][]string, 0, len(packages))
	for _, pkg := range packages {
		rows = append(rows, []string{pkg.Name, pkg.Version, pkg.Kind, pkg.SHA256})
	}
	table, err := toon.Table("prepare", []string{"package", "version", "kind", "sha256"}, rows)
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, table)
	meta, err := toon.Table("meta", []string{"release_index_sha256", "next_action"}, [][]string{{releaseIndexSHA256, "submit"}})
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, meta)
	fmt.Fprintf(stderr, "release prepare: verified %d approved artifacts for %s\n", len(packages), version)
	return 0
}

func runSubmit(args releaseArgs, stdout, stderr io.Writer) int {
	root, version, profile, path := args.root, args.version, args.profile, args.path
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release submit", "--version"))
		return 2
	}
	if profile == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release submit", "--profile"))
		return 2
	}
	if args.resolveRegistryBase() == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release submit", "--registry or BENCH_RELEASE_REGISTRY"))
		return 2
	}
	if path == "" {
		path = "first"
	}
	if path == "staged" && args.adapter == adapterNPM {
		// This check runs before the release lock and before any registry
		// call. So a staged run can never die in the middle of a publication.
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", "staged submission is not implemented for the npm adapter; publish with --path first"))
		return 1
	}
	release, err := AcquireReleaseLock(root)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
		return 1
	}
	defer release()
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	registry := newRegistry(args)
	fmt.Fprintf(stderr, "release submit: publishing %s (profile %s, path %s) against %s via the %s adapter\n", version, profile, path, args.registryBase, args.adapter)

	var record Record
	var nextAction string
	if path == "staged" {
		record, nextAction, err = RunStagedPublication(ctx, root, version, profile, registry)
	} else {
		record, err = RunFirstPublication(ctx, root, version, profile, registry)
		if err == nil {
			nextAction = "release-complete"
			if record.Result != "success" {
				nextAction = nextActionForInProgress(record)
			}
		}
	}
	exit := 0
	if err != nil {
		if nextAction == "" {
			nextAction = "resolve-integrity-mismatch"
		}
		exit = 1
	}
	printRecord(stdout, record, nextAction)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
	}
	return exit
}

func runPromote(args releaseArgs, stdout, stderr io.Writer) int {
	root, version, profile := args.root, args.version, args.profile
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release promote", "--version"))
		return 2
	}
	if profile == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release promote", "--profile"))
		return 2
	}
	if args.resolveRegistryBase() == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release promote", "--registry or BENCH_RELEASE_REGISTRY"))
		return 2
	}
	release, err := AcquireReleaseLock(root)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
		return 1
	}
	defer release()
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	registry := newRegistry(args)
	fmt.Fprintf(stderr, "release promote: promoting %s against %s via the %s adapter\n", version, args.registryBase, args.adapter)
	record, err := RunPromotion(ctx, root, version, profile, registry)
	nextAction := "release-complete"
	exit := 0
	if err != nil {
		nextAction = "resolve-integrity-mismatch"
		exit = 1
	}
	printRecord(stdout, record, nextAction)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
	}
	return exit
}

func runRollback(args releaseArgs, stdout, stderr io.Writer) int {
	root, version, profile, message := args.root, args.version, args.profile, args.message
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release rollback", "--version"))
		return 2
	}
	if profile == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release rollback", "--profile"))
		return 2
	}
	if args.resolveRegistryBase() == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release rollback", "--registry or BENCH_RELEASE_REGISTRY"))
		return 2
	}
	if message == "" {
		message = fmt.Sprintf("release %s was rolled back; see the publication record for details", version)
	}
	release, err := AcquireReleaseLock(root)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
		return 1
	}
	defer release()
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	registry := newRegistry(args)
	fmt.Fprintf(stderr, "release rollback: rolling back %s against %s via the %s adapter\n", version, args.registryBase, args.adapter)
	record, err := RunRollback(ctx, root, version, profile, message, registry)
	nextAction := "prepare"
	exit := 0
	if err != nil {
		nextAction = "resolve-integrity-mismatch"
		exit = 1
	}
	printRecord(stdout, record, nextAction)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
	}
	return exit
}

func runStatus(root string, stdout io.Writer) int {
	record, err := LoadRecord(root)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
		return 1
	}
	if record.SchemaVersion == 0 {
		meta, _ := toon.Table("meta", []string{"release_index_sha256", "path", "profile", "result", "next_action"}, [][]string{{"", "", "", "none", "prepare"}})
		fmt.Fprint(stdout, meta)
		return 0
	}
	nextAction := "submit"
	switch record.Result {
	case "success":
		nextAction = "release-complete"
	case "failed":
		nextAction = "resolve-integrity-mismatch"
	case "rolled_back":
		nextAction = "prepare"
	case "in_progress":
		// nextActionForInProgress reads only record.Transitions and
		// record.Provenance, so it touches no registry. It is the exact
		// ordering policy the staged state machine derives from the same
		// fields after appending its own run's transitions. This code never
		// re-derives that policy.
		nextAction = nextActionForInProgress(record)
	}
	printRecord(stdout, record, nextAction)
	if record.Result == "failed" {
		return 1
	}
	return 0
}

func printRecord(stdout io.Writer, record Record, nextAction string) {
	rows := make([][]string, 0, len(record.Transitions))
	for _, transition := range record.Transitions {
		rows = append(rows, []string{transition.Package, transition.Version, transition.Action, transition.Result, transition.TagState, transition.RegistryIntegrity})
	}
	table, err := toon.Table("transitions", []string{"package", "version", "action", "result", "tag_state", "registry_integrity"}, rows)
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return
	}
	fmt.Fprint(stdout, table)
	meta, err := toon.Table("meta", []string{"release_index_sha256", "path", "profile", "result", "next_action"}, [][]string{{record.ReleaseIndexSHA256, record.Path, record.Profile, record.Result, nextAction}})
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return
	}
	fmt.Fprint(stdout, meta)
}
