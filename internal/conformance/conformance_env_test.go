package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestConformanceSubprocessEnvGivesEachCallAPrivateHome(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	first := conformanceEnvironmentValue(conformanceEnvForTest(t), benchhome.Env)
	second := conformanceEnvironmentValue(conformanceEnvForTest(t), benchhome.Env)
	if first == second {
		t.Fatalf("conformance probes share Bench home %q", first)
	}
	for _, home := range []string{first, second} {
		if home == os.Getenv(benchhome.Env) || home == filepath.Join(os.TempDir(), "bench-conformance-home") {
			t.Fatalf("conformance probe retained a shared Bench home %q", home)
		}
		entries, err := os.ReadDir(home)
		if err != nil || len(entries) != 0 {
			t.Fatalf("private Bench home = %q, entries=%v, err=%v", home, entries, err)
		}
	}
}

func conformanceEnvironmentValue(env []string, name string) string {
	for i := len(env) - 1; i >= 0; i-- {
		if value, ok := strings.CutPrefix(env[i], name+"="); ok {
			return value
		}
	}
	return ""
}

func TestConformanceSubprocessEnvCleanupRemovesOnlyItsOwnHome(t *testing.T) {
	env, cleanup, err := conformanceSubprocessEnv("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)
	home := conformanceEnvironmentValue(env, benchhome.Env)
	other := conformanceEnvironmentValue(conformanceEnvForTest(t), benchhome.Env)
	if err := os.WriteFile(filepath.Join(home, "record"), []byte("private"), 0o600); err != nil {
		t.Fatal(err)
	}
	cleanup()
	if _, err := os.Stat(home); !os.IsNotExist(err) {
		t.Fatalf("private Bench home survived cleanup: %v", err)
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatalf("cleanup removed another probe's home: %v", err)
	}
}

func TestConformanceSubprocessEnvKeepsTheSharedNpmCache(t *testing.T) {
	for _, explicit := range []string{"", filepath.Join(t.TempDir(), "operator-cache")} {
		t.Run(explicit, func(t *testing.T) {
			t.Setenv("NPM_CONFIG_CACHE", explicit)
			cache := conformanceEnvironmentValue(conformanceEnvForTest(t), "NPM_CONFIG_CACHE")
			want := explicit
			if want == "" {
				want = filepath.Join(os.TempDir(), "bench-npm-cache")
			}
			if cache != want {
				t.Fatalf("npm cache = %q, want shared cache %q", cache, want)
			}
		})
	}
}

func TestConformanceProbeRefusesAnUnavailablePrivateHome(t *testing.T) {
	root := t.TempDir()
	blocked := filepath.Join(root, "not-a-directory")
	if err := os.WriteFile(blocked, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", blocked)
	probe := runAtCleanEnv(root, "sh", "-c", "touch child-started")
	if probe.ExitCode != 1 || probe.Err == nil || !strings.Contains(probe.Err.Error(), "conformance Bench home") {
		t.Fatalf("probe = %+v, want private home allocation refusal", probe)
	}
	if _, err := os.Stat(filepath.Join(root, "child-started")); !os.IsNotExist(err) {
		t.Fatalf("refused probe started a child: %v", err)
	}
}

// conformanceSubprocessEnv keeps probe children outside the operator's Bench home.
// The caller removes the private home after its children exit.
func conformanceSubprocessEnv(pathPrefix string) ([]string, func(), error) {
	home, err := os.MkdirTemp("", "bench-conformance-*")
	if err != nil {
		return nil, nil, fmt.Errorf("conformance Bench home: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(home) }
	env := make([]string, 0, len(os.Environ()))
	hasNpmCache := false
	for _, kv := range os.Environ() {
		// The scrub is symmetric across every conformance control var. Any one leaking into a
		// probe subprocess is the recursive-cascade shape.
		if strings.HasPrefix(kv, "BENCH_CONFORMANCE_ROOT=") ||
			strings.HasPrefix(kv, registry.ConformanceTierEnv+"=") ||
			strings.HasPrefix(kv, registry.ConformanceChecksEnv+"=") ||
			strings.HasPrefix(kv, registry.ConformanceInheritedEnv+"=") ||
			strings.HasPrefix(kv, benchhome.Env+"=") {
			continue
		}
		if strings.HasPrefix(kv, "NPM_CONFIG_CACHE=") && strings.TrimPrefix(kv, "NPM_CONFIG_CACHE=") != "" {
			hasNpmCache = true
		}
		env = append(env, kv)
	}
	if !hasNpmCache {
		env = append(env, "NPM_CONFIG_CACHE="+filepath.Join(os.TempDir(), "bench-npm-cache"))
	}
	if pathPrefix != "" {
		env = append(env, "PATH="+pathPrefix+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	env = append(env, benchhome.Env+"="+home)
	return env, cleanup, nil
}

func conformanceEnvForTest(t *testing.T) []string {
	t.Helper()
	env, cleanup, err := conformanceSubprocessEnv("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)
	return env
}

// coreless returns env with a PATH carrying neither a bench wrapper nor the stub dir. The
// shims' wrapper search then comes up empty, and each one takes its own missing-core rim.
func coreless(env []string) []string {
	out := make([]string, 0, len(env)+1)
	for _, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			continue
		}
		out = append(out, kv)
	}
	return append(out, "PATH=/usr/bin:/bin")
}

func TestConformanceSubprocessEnvStripsConformanceControlVars(t *testing.T) {
	t.Setenv("BENCH_CONFORMANCE_ROOT", "/tmp/outer-root")
	t.Setenv(registry.ConformanceTierEnv, "ship")
	t.Setenv(registry.ConformanceChecksEnv, "line-routing,package-core-guard")
	t.Setenv(registry.ConformanceInheritedEnv, "bounds-policy")

	for _, kv := range conformanceEnvForTest(t) {
		for _, name := range []string{"BENCH_CONFORMANCE_ROOT", registry.ConformanceTierEnv, registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv} {
			if strings.HasPrefix(kv, name+"=") {
				t.Fatalf("%s leaked into the probe subprocess env: %q", name, kv)
			}
		}
	}
}

func TestConformanceSubprocessEnvProvidesWritableNpmCache(t *testing.T) {
	oldCache, hadCache := os.LookupEnv("NPM_CONFIG_CACHE")
	if err := os.Unsetenv("NPM_CONFIG_CACHE"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadCache {
			_ = os.Setenv("NPM_CONFIG_CACHE", oldCache)
			return
		}
		_ = os.Unsetenv("NPM_CONFIG_CACHE")
	})

	cache := conformanceEnvironmentValue(conformanceEnvForTest(t), "NPM_CONFIG_CACHE")
	if cache == "" {
		t.Fatal("NPM_CONFIG_CACHE missing from conformance subprocess env")
	}
	if !strings.HasPrefix(filepath.Clean(cache), filepath.Clean(os.TempDir())+string(os.PathSeparator)) {
		t.Fatalf("NPM_CONFIG_CACHE = %q, want temp-backed cache", cache)
	}
}
