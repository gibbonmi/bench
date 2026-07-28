package surface

import (
	"os"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestMain declares the whole package dev tier: every build-artifacts.sh invocation it
// spawns honors the ambient Go caches and skips the reproducibility second build, which
// is a release-tier claim. Both contract.ProcessEnv and the fixture-driven env merge
// build from os.Environ(), so setting it once here reaches every subprocess — and a call
// site added later cannot silently reintroduce a hermetic build the way a per-call-site
// opt-in would.
func TestMain(m *testing.M) {
	if err := os.Setenv(contract.SharedBuildCacheEnv, "1"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
