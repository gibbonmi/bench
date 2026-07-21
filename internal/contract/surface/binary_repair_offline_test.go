package surface

import "testing"

func testRepairBenchOfflineExact(t *testing.T) {
	for _, value := range []string{"1", "", "0", "true"} {
		t.Run("value="+value, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, "9.8.7", "#!/bin/sh\necho repaired\n")
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
