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

const usageLine = "usage: bench release prepare|submit|promote|rollback|status --version <version> [--profile public|bank] [--root <dir>] [--registry <base-url>] [--path first|staged] [--message <text>]"

// Command is the `bench release <prepare|submit|promote|rollback|status>`
// entry point. It is idempotent and non-interactive: prepare only verifies
// the approved release directory locally; submit runs (or resumes) the
// first-publication or staged-submission state machine depending on --path;
// promote moves the "latest" dist-tag once the complete set reverifies live;
// rollback removes candidate tags and deprecates a bad version; status
// reports the durable record without touching the registry. TOON tables go to
// stdout, progress to stderr, structured errors to stdout. Exit 0
// success/no-op, 1 unsatisfied release intent, 2 usage.
func Command(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stdout, toon.Usage("bench release", "<missing subcommand>"))
		return 2
	}
	sub := args[0]
	root, version, profile, registryBase, path, message, usageErr := parseCommandArgs(args[1:])
	if usageErr != nil {
		fmt.Fprintln(stdout, usageErr.Error())
		return 2
	}
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			fmt.Fprintln(stdout, toon.Errorf("input", "could not resolve the working directory"))
			return 1
		}
	}
	root, _ = filepath.Abs(root)

	switch sub {
	case "prepare":
		return runPrepare(root, version, stdout, stderr)
	case "submit":
		return runSubmit(root, version, profile, registryBase, path, stdout, stderr)
	case "promote":
		return runPromote(root, version, profile, registryBase, stdout, stderr)
	case "rollback":
		return runRollback(root, version, profile, registryBase, message, stdout, stderr)
	case "status":
		return runStatus(root, stdout)
	default:
		fmt.Fprintln(stdout, toon.Usage("bench release", sub))
		return 2
	}
}

func parseCommandArgs(args []string) (root, version, profile, registryBase, path, message string, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--root":
			i++
			if i >= len(args) {
				return "", "", "", "", "", "", fmt.Errorf("%s", toon.MissingArg("bench release", "--root value"))
			}
			root = args[i]
		case "--version":
			i++
			if i >= len(args) {
				return "", "", "", "", "", "", fmt.Errorf("%s", toon.MissingArg("bench release", "--version value"))
			}
			version = args[i]
		case "--profile":
			i++
			if i >= len(args) {
				return "", "", "", "", "", "", fmt.Errorf("%s", toon.MissingArg("bench release", "--profile value"))
			}
			profile = args[i]
		case "--registry":
			i++
			if i >= len(args) {
				return "", "", "", "", "", "", fmt.Errorf("%s", toon.MissingArg("bench release", "--registry value"))
			}
			registryBase = args[i]
		case "--path":
			i++
			if i >= len(args) {
				return "", "", "", "", "", "", fmt.Errorf("%s", toon.MissingArg("bench release", "--path value"))
			}
			path = args[i]
		case "--message":
			i++
			if i >= len(args) {
				return "", "", "", "", "", "", fmt.Errorf("%s", toon.MissingArg("bench release", "--message value"))
			}
			message = args[i]
		default:
			return "", "", "", "", "", "", fmt.Errorf("%s", toon.Usage("bench release", args[i]))
		}
	}
	if profile != "" && profile != "public" && profile != "bank" {
		return "", "", "", "", "", "", fmt.Errorf("%s", toon.Usage("bench release", "--profile "+profile))
	}
	if path != "" && path != "first" && path != "staged" {
		return "", "", "", "", "", "", fmt.Errorf("%s", toon.Usage("bench release", "--path "+path))
	}
	return root, version, profile, registryBase, path, message, nil
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

func runSubmit(root, version, profile, registryBase, path string, stdout, stderr io.Writer) int {
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release submit", "--version"))
		return 2
	}
	if profile == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release submit", "--profile"))
		return 2
	}
	if registryBase == "" {
		registryBase = os.Getenv("BENCH_RELEASE_REGISTRY")
	}
	if registryBase == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release submit", "--registry or BENCH_RELEASE_REGISTRY"))
		return 2
	}
	if path == "" {
		path = "first"
	}
	release, err := AcquireReleaseLock(root)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("unsatisfied release intent", err.Error()))
		return 1
	}
	defer release()
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	registry := NewFixtureRegistry(registryBase)
	fmt.Fprintf(stderr, "release submit: publishing %s (profile %s, path %s) against %s\n", version, profile, path, registryBase)

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

func runPromote(root, version, profile, registryBase string, stdout, stderr io.Writer) int {
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release promote", "--version"))
		return 2
	}
	if profile == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release promote", "--profile"))
		return 2
	}
	if registryBase == "" {
		registryBase = os.Getenv("BENCH_RELEASE_REGISTRY")
	}
	if registryBase == "" {
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
	registry := NewFixtureRegistry(registryBase)
	fmt.Fprintf(stderr, "release promote: promoting %s against %s\n", version, registryBase)
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

func runRollback(root, version, profile, registryBase, message string, stdout, stderr io.Writer) int {
	if version == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release rollback", "--version"))
		return 2
	}
	if profile == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release rollback", "--profile"))
		return 2
	}
	if registryBase == "" {
		registryBase = os.Getenv("BENCH_RELEASE_REGISTRY")
	}
	if registryBase == "" {
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
	registry := NewFixtureRegistry(registryBase)
	fmt.Fprintf(stderr, "release rollback: rolling back %s against %s\n", version, registryBase)
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
		// Registry-free by construction: nextActionForInProgress reads only
		// record.Transitions and record.Provenance, the exact ordering
		// policy the staged state machine derives from the same fields
		// after appending its own run's transitions — one source, never
		// re-derived here.
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
