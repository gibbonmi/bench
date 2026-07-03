// Command bench is the compiled core of the Bench kit — the strangler target the
// shell CLI routes ported subcommands into. Dispatch is a `commands` map of the ported
// AXI query subcommands (learnings, maps, guards, diff, coverage), each resolving repo
// state and returning its stdout plus an exit code, plus a direct `version` case that
// needs the build-time GOOS/GOARCH rather than repo state. Every later slice adds names
// to that map; the shell router (bin/bench.sh) grows names, not mechanisms.
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/maps"
)

// version is stamped at build time via -ldflags "-X main.version=<pkg.json version>"
// (see scripts/go-build.sh — the one source of build flags). An unstamped build
// prints "dev", which is the tell that the binary was not produced by the gate or
// the release workflow.
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// commands are the ported AXI query subcommands: each resolves arguments and repo
// state, then returns its complete stdout plus an exit code — the dispatch below just
// writes and exits. This map is the one seam that grows per slice; version stays a
// direct case because it needs the build-time GOOS/GOARCH, not repo state.
var commands = map[string]func([]string) (string, int){
	"learnings": learnings.Command,
	"maps":      maps.Command,
	"guards":    guards.Command,
	"diff":      diff.Command,
	"coverage":  coverage.Command,
}

func run(args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "bench: no subcommand")
		return 2
	}
	if fn, ok := commands[args[0]]; ok {
		out, code := fn(args[1:])
		fmt.Fprint(stdout, out)
		return code
	}
	switch args[0] {
	case "version":
		fmt.Fprintln(stdout, versionLine(version, runtime.GOOS, runtime.GOARCH))
		return 0
	default:
		fmt.Fprintf(stderr, "bench: unknown subcommand: %q\n", args[0])
		return 2
	}
}

// versionLine renders the single line `bench version` prints. Kept as a pure
// function so the format has one source and a table test can pin it without a
// process boundary.
func versionLine(v, goos, goarch string) string {
	return fmt.Sprintf("benchkit %s (%s/%s)", v, goos, goarch)
}
