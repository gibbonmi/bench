package testreport

// The environments a focused run hands its Go child.

import (
	"strings"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
)

func conformanceEnvironment(base []string, root, scope string, selection *runbinary.Selection) ([]string, error) {
	consumerOnly := environmentValue(base, registry.ConsumerOnlyEnv) == "1"
	env, err := selectedRunEnvironment(base, selection)
	if err != nil {
		return nil, err
	}
	env = append(env,
		registry.ConformanceRootEnv+"="+root,
		registry.ConformanceTierEnv+"="+string(registry.Dev),
		registry.ConformanceScopeEnv+"="+scope,
	)
	if consumerOnly {
		env = append(env, registry.ConsumerOnlyEnv+"=1")
	}
	return env, nil
}

// conformanceScope answers the conformance scope of a named check. The root-conformance
// check passes the empty scope, which runs every registered dev-tier check, because the
// entry test refuses a scope that names no registered check. Every other check is its
// own scope.
func conformanceScope(check string) string {
	if check == RootConformanceCheckName {
		return ""
	}
	return check
}

func environmentValue(env []string, name string) string {
	prefix := name + "="
	for i := len(env) - 1; i >= 0; i-- {
		if after, ok := strings.CutPrefix(env[i], prefix); ok {
			return after
		}
	}
	return ""
}

// selectedRunEnvironment pins the executable and operator cache before a child gets a private HOME.
func selectedRunEnvironment(base []string, selection *runbinary.Selection) ([]string, error) {
	env, err := testEnvironment(base, selection.Path)
	if err != nil {
		return nil, err
	}
	return append(env, "BENCH_KIT="+selection.SourceRoot), nil
}

// testEnvironment returns the environment the focused run's Go child carries: the
// caller's, without the inherited conformance, capability, and wrapper routing entries,
// with the selected Bench executable, and with the Bench build cache entry so a focused
// run warms the archives a gate reads.
func testEnvironment(base []string, binary string) ([]string, error) {
	return gocache.Apply(runbinary.WithEnv(env.WithoutWrapperRouting(withoutConformanceEnvironment(base)), binary))
}

func withoutConformanceEnvironment(base []string) []string {
	env := base
	for _, name := range []string{
		registry.ConformanceRootEnv,
		registry.ConformanceTierEnv,
		registry.ConformanceScopeEnv,
		registry.ConformanceChecksEnv,
		registry.ConformanceInheritedEnv,
		registry.ConsumerOnlyEnv,
		capability.LogEnv,
	} {
		env = capability.WithoutEnvironment(env, name)
	}
	return env
}
