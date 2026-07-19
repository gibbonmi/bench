package publication

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	registry := NewFixtureRegistry(registryBase)
	fmt.Fprintf(stderr, "release submit: publishing %s (profile %s, path %s) against %s\n", version, profile, path, registryBase)

	var record Record
	var nextAction string
	var err error
	if path == "staged" {
		record, nextAction, err = RunStagedPublication(ctx, root, version, profile, registry)
	} else {
		record, err = RunFirstPublication(ctx, root, version, profile, registry)
		nextAction = "release-complete"
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
	if registryBase == "" {
		registryBase = os.Getenv("BENCH_RELEASE_REGISTRY")
	}
	if registryBase == "" {
		fmt.Fprintln(stdout, toon.MissingArg("bench release promote", "--registry or BENCH_RELEASE_REGISTRY"))
		return 2
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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
		if record.Path == "public" {
			nextAction = statusNextActionForStaged(record)
		}
	}
	printRecord(stdout, record, nextAction)
	if record.Result == "failed" {
		return 1
	}
	return 0
}

// statusNextActionForStaged reports the staged path's handoff from the
// durable record alone (status never touches the registry): it walks the
// transitions to see which platform packages are already verified live, so
// a rerun of `status` shows the same next_action `submit` would compute
// without needing a registry round trip.
func statusNextActionForStaged(record Record) string {
	verifiedLive := map[string]bool{}
	var wrapper string
	for _, t := range record.Transitions {
		if t.Action == "verify" && t.Result == "success" {
			verifiedLive[t.Package] = true
		}
	}
	for _, p := range record.Provenance {
		if p.Package == "redbench" {
			wrapper = p.Package
		}
	}
	for _, p := range record.Provenance {
		if p.Package == wrapper {
			continue
		}
		if !verifiedLive[p.Package] {
			return nextActionApprovePlatforms
		}
	}
	if wrapper != "" && !verifiedLive[wrapper] {
		return nextActionApproveWrapper
	}
	return nextActionPromote
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
