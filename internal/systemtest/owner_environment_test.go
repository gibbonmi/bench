//go:build system

package systemtest

import (
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/bounds"
)

func systemBaseEnvironment() []string {
	return mergeEnvironment(os.Environ(), []string{bounds.UnboundedWaitsEnv})
}

// mergeEnvironment layers overrides over base. An override spelled as a bare NAME with no
// `=` removes that variable instead of setting it. This drives the wrapper's own binary
// search, a resolution path that BENCH_RUN_BINARY normally pre-empts. That path only runs
// with the variable genuinely absent.
func mergeEnvironment(base, overrides []string) []string {
	want := map[string]string{}
	for _, entry := range overrides {
		key, _, _ := strings.Cut(entry, "=")
		want[key] = entry
	}
	out := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		key, _, _ := strings.Cut(entry, "=")
		if _, replaced := want[key]; replaced {
			continue
		}
		out = append(out, entry)
	}
	for _, entry := range overrides {
		if strings.Contains(entry, "=") {
			out = append(out, entry)
		}
	}
	return out
}

// namesEnv reports whether entries names key, spelled either as key=value or as the
// bare-NAME removal form.
func namesEnv(entries []string, key string) bool {
	for _, entry := range entries {
		if entryKey, _, _ := strings.Cut(entry, "="); entryKey == key {
			return true
		}
	}
	return false
}

// childEnvironment gives command helpers a private Bench home. A
// caller that names its own BENCH_HOME in overrides keeps it; one that does not gets
// o.home instead, so a bench child never inherits the operator's real Bench home from
// the process running go test.
func (o *systemOwner) childEnvironment(overrides []string) []string {
	if !namesEnv(overrides, benchhome.Env) {
		overrides = append(append([]string{}, overrides...), benchhome.Env+"="+o.home)
	}
	return mergeEnvironment(systemBaseEnvironment(), overrides)
}
