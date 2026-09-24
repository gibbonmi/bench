package main

import (
	"os"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// unstampedVersion is the version an unstamped build reports. It names no release.
const unstampedVersion = "dev"

// prepareProcessEnvironment is the binary's one edit of its own environment, before any
// verb runs. Every child a verb starts inherits the result.
func prepareProcessEnvironment() {
	// The wrapper's implicit-repair grant is spent once it execs this binary. Scrub it so
	// gate phases and their fixtures never inherit an invocation-dependent privilege.
	os.Unsetenv("BENCH_ALLOW_IMPLICIT_REPAIR")
	benchhome.Export()
	// Each seam record line names the Bench version that wrote it. An unstamped build
	// has no release to name, so it hands none, and its lines carry no version key.
	if version != unstampedVersion {
		otelrecord.SetVersion(version)
	}
}
