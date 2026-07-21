package surface

import "testing"

func testRepairOptInExact(t *testing.T) {
	for _, value := range []string{"1", "", "0", "true"} {
		t.Run("value="+value, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho repaired\n")
			out := f.BenchEnv(map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": value}, "version")
			if value == "1" {
				out.RequireExit(0)
				return
			}
			out.RequireExit(127)
			if registry.Hits() != 0 {
				t.Fatalf("BENCH_REPAIR=%q started repair", value)
			}
		})
	}
}

func testRepairSuppressionPrecedence(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		message string
	}{
		{name: "offline", env: map[string]string{"BENCH_OFFLINE": "1"}, message: "repair suppressed by BENCH_OFFLINE=1"},
		{name: "no-repair", env: map[string]string{"BENCH_NO_REPAIR": "1"}, message: "repair disabled by BENCH_NO_REPAIR"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho should-not-run\n")
			tc.env["BENCH_KIT"], tc.env["BENCH_NPM_REGISTRY"], tc.env["BENCH_REPAIR"] = kit, registry.URL, "1"
			out := f.BenchEnv(tc.env, "version")
			out.RequireExit(127)
			out.RequireContains(out.Stderr, tc.message)
			if registry.Hits() != 0 {
				t.Fatalf("suppressed repair hit registry")
			}
		})
	}
}

func testRepairBenchOfflineExact(t *testing.T) {
	for _, value := range []string{"1", "", "0", "true"} {
		t.Run("value="+value, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho repaired\n")
			env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}
			if value != "" {
				env["BENCH_OFFLINE"] = value
			}
			out := f.BenchEnv(env, "version")
			if value == "1" {
				out.RequireExit(127)
				out.RequireContains(out.Stderr, "repair suppressed by BENCH_OFFLINE=1")
				if registry.Hits() != 0 {
					t.Fatalf("offline repair hit registry")
				}
				return
			}
			out.RequireExit(0)
			if registry.Hits() == 0 {
				t.Fatalf("BENCH_OFFLINE=%q incorrectly enabled offline mode", value)
			}
		})
	}
}
