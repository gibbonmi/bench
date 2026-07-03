// Command bench is the compiled core of the Bench kit — the strangler target the
// shell CLI routes ported subcommands into. Slice 1 ships exactly one: version.
// Every later slice adds a subcommand here; the shell router (bin/bench.sh) grows
// names, not mechanisms.
package main

import (
	"fmt"
	"os"
	"runtime"
)

// version is stamped at build time via -ldflags "-X main.version=<pkg.json version>"
// (see scripts/go-build.sh — the one source of build flags). An unstamped build
// prints "dev", which is the tell that the binary was not produced by the gate or
// the release workflow.
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "bench: no subcommand")
		return 2
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
