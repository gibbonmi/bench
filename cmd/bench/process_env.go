package main

import (
	"os"

	"github.com/gibbonmi/bench/internal/benchhome"
)

// prepareProcessEnvironment is the binary's one edit of its own environment, before any
// verb runs. Every child a verb starts inherits the result.
func prepareProcessEnvironment() {
	// The wrapper's implicit-repair grant is spent once it execs this binary. Scrub it so
	// gate phases and their fixtures never inherit an invocation-dependent privilege.
	os.Unsetenv("BENCH_ALLOW_IMPLICIT_REPAIR")
	benchhome.Export()
}
