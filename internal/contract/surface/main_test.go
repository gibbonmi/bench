package surface

import (
	"os"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestMain declares the whole package dev tier: every build-artifacts.sh invocation it
// spawns honors the ambient Go caches and skips the reproducibility second build, which
// is a release-tier claim this package does not make.
func TestMain(m *testing.M) {
	if err := os.Setenv(contract.SharedBuildCacheEnv, "1"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
