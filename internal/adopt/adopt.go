package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Run(args []string, stdout, stderr io.Writer, version string) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "bench adopt: missing subcommand")
		return 2
	}
	switch args[0] {
	case "setup":
		return Setup(args[1:], os.Stdin, stdout, stderr, version)
	case "link":
		return Link(args[1:], stdout, stderr, version)
	case "init":
		return Init(args[1:], stdout, stderr)
	case "doctor":
		return Doctor(args[1:], stdout, stderr, version)
	case "unlink":
		return Unlink(args[1:], stdout, stderr)
	case "upgrade":
		return Upgrade(args[1:], stdout, stderr, version)
	default:
		fmt.Fprintf(stderr, "bench adopt: unknown subcommand: %q\n", args[0])
		return 2
	}
}

func kitDir() string {
	if kit := os.Getenv("BENCH_KIT"); kit != "" {
		return kit
	}
	if exe, err := os.Executable(); err == nil {
		// Development layout: <kit>/dist/bench. Platform package layout:
		// <pkg>/bin/bench, where callers should pass BENCH_KIT from the wrapper.
		return filepath.Clean(filepath.Join(filepath.Dir(exe), ".."))
	}
	return "."
}
