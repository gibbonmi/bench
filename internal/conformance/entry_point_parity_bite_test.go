package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEntryPointParityRuntimeRowsBite runs the runtime comparator against a synthetic core
// and synthetic shims. Each shim body is one way a real shim drifts: it decides alone, it
// answers with its own exit code, or it appends output the core never produced.
func TestEntryPointParityRuntimeRowsBite(t *testing.T) {
	const rel = ".bench/hooks/block-bench-follow-on.sh"
	const call = "#!/usr/bin/env bash\ninput=\"$(cat)\"\nprintf '%s' \"$input\" | \"$(command -v bench)\" guard-bench-follow-on\n"
	tests := []struct{ name, body, want string }{
		{"reaches the core", call, ""},
		{"decides alone", "#!/usr/bin/env bash\nexit 0\n", "does not reach the registry"},
		{"holds a second opinion", call + "exit 1\n", "exits 1 and the direct guard-bench-follow-on exits 0"},
		{"appends its own output", call + "printf 'shim tail\\n'\n", "stdout does not end with the direct guard-bench-follow-on stdout"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parityRuntimeRoot(t, rel, tt.body)
			diags := checkEntryPointParity(root)
			if tt.want == "" {
				if len(diags) != 0 {
					t.Fatalf("a shim that reaches the core is not green:\n%s", strings.Join(diags, "\n"))
				}
				return
			}
			if !containsDiagnostic(diags, tt.want) {
				t.Fatalf("the drifted shim did not bite with %q:\n%s", tt.want, strings.Join(diags, "\n"))
			}
		})
	}
}

// parityRuntimeRoot builds a root holding one shim and a wrapper that stands in for the
// compiled core. The stand-in prints the observed id the way the real registry does, so
// the comparator is graded without a build.
func parityRuntimeRoot(t *testing.T, rel, body string) string {
	t.Helper()
	const core = "#!/usr/bin/env bash\n" +
		"cmd=\"${1:-status}\"\n" +
		"[ \"${BENCH_COMMAND_OBSERVE:-}\" = 1 ] && printf 'command-registry:%s\\n' \"$cmd\" >&2\n" +
		"cat >/dev/null 2>/dev/null\n" +
		"printf 'core %s\\n' \"$cmd\"\n" +
		"exit 0\n"
	files := map[string]string{"bin/bench.sh": core, rel: body}
	for staticRel, content := range parityHonestStatics {
		files[staticRel] = content
	}
	root := throwawayRoot{files: files}.build(t)
	if err := os.Chmod(filepath.Join(root, "bin", "bench.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}
