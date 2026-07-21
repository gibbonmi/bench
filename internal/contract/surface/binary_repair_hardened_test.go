package surface

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestBinaryRepairContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "binary repair download-and-run contract failed", testRepairDownloadsAndRuns)
	contract.RunParallel(t, "binary repair bad-integrity contract failed", testRepairRefusesBadHash)
	contract.RunParallel(t, "binary repair malformed-tar contract failed", testRepairRefusesMalformedTar)
	contract.RunParallel(t, "binary repair offline fallback contract failed", testRepairOffline)
	contract.RunParallel(t, "binary repair plumbing exclusion contract failed", testRepairSkipsPlumbing)
	contract.RunParallel(t, "binary repair announcement contract failed", testRepairAnnounces)
	contract.RunParallel(t, "binary repair idempotency contract failed", testRepairIdempotent)
	contract.RunParallel(t, "binary repair version-keyed cache contract failed", testRepairVersionKeyed)
	contract.RunParallel(t, "binary repair no-node contract failed", testRepairNoNode)
	contract.RunParallel(t, "binary repair disabled contract failed", testRepairDisabled)
	contract.RunParallel(t, "binary repair BENCH_OFFLINE exact-value contract failed", testRepairBenchOfflineExact)
	contract.RunParallel(t, "binary repair opt-in exact-value contract failed", testRepairOptInExact)
	contract.RunParallel(t, "binary repair suppression precedence contract failed", testRepairSuppressionPrecedence)
	contract.RunParallel(t, "binary repair torn-cache contract failed", testRepairReplacesTornCache)
	contract.RunParallel(t, "linked manifest repair contract failed", testRepairReadsLinkedManifestWithoutNewline)
	contract.RunParallel(t, "malformed manifest version escaped repair cache", testRepairRejectsMalformedVersion)
	contract.RunParallel(t, "interrupted repair promoted partial cache", testRepairInterruptedPromotion)
	contract.RunParallel(t, "fresh clone repair required ambient tooling", testRepairMinimalPortablePath)
	contract.RunParallel(t, "repair explicit-default contract failed", testRepairExplicitDefault)
	contract.RunParallel(t, "repair explicit subcommand contract failed", testRepairSubcommand)
	contract.RunParallel(t, "repair argument contract failed", testRepairArguments)
	contract.RunParallel(t, "repair prune contract failed", testRepairPrune)
	contract.RunParallel(t, "repair pin manifest fail-closed contract failed", testRepairPinManifestFailures)
	contract.RunParallel(t, "repair resource bounds contract failed", testRepairResourceBounds)
	contract.RunParallel(t, "repair losing-racer cleanup contract failed", testRepairLosingRacerPreservesWinner)
	contract.RunParallel(t, "repair earliest interrupt cleanup contract failed", testRepairEarliestInterrupt)
}

func testRepairLosingRacerPreservesWinner(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho loser\n")
	ready := filepath.Join(f.Root, "loser-ready")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "", "BENCH_TEST_REPAIR_READY_FILE": ready, "BENCH_TEST_REPAIR_FAIL_AFTER_READY": "1"}
	cmd := exec.Command("bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "repair")
	cmd.Dir, cmd.Env = f.Root, contract.ProcessEnv(f.Env, env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waitForRepairMarker(t, cmd, ready)
	target := binaryRepairCachePath(t, f, "9.8.7")
	contract.WriteFileAbs(t, target, "#!/bin/sh\necho winner\n")
	if err := os.Chmod(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(ready); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err == nil {
		t.Fatal("injected losing repair succeeded")
	}
	data, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(data), "winner") {
		t.Fatalf("loser removed or replaced winner: %q %v", data, err)
	}
	temps, _ := filepath.Glob(filepath.Join(filepath.Dir(target), ".bench-*.tmp"))
	if len(temps) != 0 {
		t.Fatalf("loser left temps: %v", temps)
	}
}

func testRepairEarliestInterrupt(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho never\n", binaryRepairHangMetadata())
	ready := filepath.Join(f.Root, "start-ready")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "", "BENCH_TEST_REPAIR_START_READY_FILE": ready, "BENCH_TEST_REPAIR_DEADLINE_MS": "5000"}
	cmd := exec.Command("bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "repair")
	cmd.Dir, cmd.Env = f.Root, contract.ProcessEnv(f.Env, env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waitForRepairMarker(t, cmd, ready)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	err := cmd.Wait()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 130 {
		t.Fatalf("earliest SIGINT exit = %v, want 130", err)
	}
	requirePathAbsent(t, binaryRepairCachePath(t, f, "9.8.7"), "earliest interrupt changed target")
}

func waitForRepairMarker(t *testing.T, cmd *exec.Cmd, ready string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			return
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			t.Fatal("repair did not reach synchronization marker")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func testRepairPinManifestFailures(t *testing.T) {
	for _, tc := range []struct {
		name, content string
		remove        bool
	}{
		{name: "absent", remove: true},
		{name: "empty", content: ""},
		{name: "unparseable", content: "{"},
		{name: "wrong-schema", content: `{"schema_version":2,"pins":[]}` + "\n"},
		{name: "pins-not-array", content: `{"schema_version":1,"pins":{}}` + "\n"},
		{name: "malformed-entry", content: `{"schema_version":1,"pins":[{"name":7,"version":"9.8.7","integrity":"sha512-AA=="}]}` + "\n"},
		{name: "duplicate-entry", content: `{"schema_version":1,"pins":[{"name":"@redbench/` + binaryRepairPlatformSuffix(t) + `","version":"9.8.7","integrity":"sha512-AA=="},{"name":"@redbench/` + binaryRepairPlatformSuffix(t) + `","version":"9.8.7","integrity":"sha512-AA=="}]}` + "\n"},
		{name: "entry-missing", content: `{"schema_version":1,"pins":[]}` + "\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho forbidden\n")
			manifest := filepath.Join(kit, "binary-pins.json")
			if tc.remove {
				if err := os.Remove(manifest); err != nil {
					t.Fatal(err)
				}
			} else {
				contract.WriteFileAbs(t, manifest, tc.content)
			}
			out := f.BenchEnv(map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "1"}, "version")
			out.RequireExit(127)
			out.RequireContains(out.Stderr, "pin manifest binary-pins.json")
			if registry.Hits() != 0 {
				t.Fatalf("invalid pin manifest consulted registry: hits=%d", registry.Hits())
			}
			requirePathAbsent(t, binaryRepairCachePath(t, f, "9.8.7"), "invalid pin manifest installed target")
		})
	}
}

func testRepairResourceBounds(t *testing.T) {
	t.Run("total deadline", func(t *testing.T) {
		f, kit := binaryRepairFixtureKit(t)
		registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho never\n", binaryRepairHangMetadata())
		out := f.BenchEnv(map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "", "BENCH_TEST_REPAIR_DEADLINE_MS": "50"}, "repair")
		out.RequireExit(1)
		out.RequireContains(out.Stderr, "60-second total repair deadline exceeded")
	})
	t.Run("download exact and plus one", func(t *testing.T) {
		f, kit := binaryRepairFixtureKit(t)
		binary := "#!/bin/sh\necho bounded\n"
		tgz := binaryRepairTarball(t, binary)
		registry := newBinaryRepairRegistry(t, kit, "9.8.7", binary, binaryRepairTarballBytes(tgz))
		env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "", "BENCH_TEST_REPAIR_DOWNLOAD_LIMIT": fmt.Sprint(len(tgz))}
		f.BenchEnv(env, "repair").RequireExit(0)
		if err := os.Remove(binaryRepairCachePath(t, f, "9.8.7")); err != nil {
			t.Fatal(err)
		}
		env["BENCH_TEST_REPAIR_DOWNLOAD_LIMIT"] = fmt.Sprint(len(tgz) - 1)
		out := f.BenchEnv(env, "repair")
		out.RequireExit(1)
		out.RequireContains(out.Stderr, "100 MB download limit exceeded")
	})
	t.Run("decompressed exact and plus one", func(t *testing.T) {
		f, kit := binaryRepairFixtureKit(t)
		binary := "#!/bin/sh\necho expanded\n"
		tgz := binaryRepairTarball(t, binary)
		decompressed := binaryRepairGunzip(t, tgz)
		registry := newBinaryRepairRegistry(t, kit, "9.8.7", binary, binaryRepairTarballBytes(tgz))
		env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "", "BENCH_TEST_REPAIR_DECOMPRESSED_LIMIT": fmt.Sprint(len(decompressed))}
		f.BenchEnv(env, "repair").RequireExit(0)
		if err := os.Remove(binaryRepairCachePath(t, f, "9.8.7")); err != nil {
			t.Fatal(err)
		}
		env["BENCH_TEST_REPAIR_DECOMPRESSED_LIMIT"] = fmt.Sprint(len(decompressed) - 1)
		out := f.BenchEnv(env, "repair")
		out.RequireExit(1)
		out.RequireContains(out.Stderr, "200 MB decompressed limit exceeded")
	})
	t.Run("decompression bomb stops near the incremental limit", func(t *testing.T) {
		const expandedSize = 8 * 1024 * 1024
		const limit = 64 * 1024
		f, kit := binaryRepairFixtureKit(t)
		tgz := binaryRepairGzipBomb(t, expandedSize)
		registry := newBinaryRepairRegistry(t, kit, "9.8.7", "unused", binaryRepairTarballBytes(tgz))
		out := f.BenchEnv(map[string]string{
			"BENCH_KIT":                            kit,
			"BENCH_NPM_REGISTRY":                   registry.URL,
			"BENCH_REPAIR":                         "",
			"BENCH_TEST_REPAIR_DECOMPRESSED_LIMIT": strconv.Itoa(limit),
		}, "repair")
		out.RequireExit(1)
		out.RequireContains(out.Stderr, "200 MB decompressed limit exceeded")
		match := regexp.MustCompile(`exceeded after ([0-9]+) bytes`).FindStringSubmatch(out.Stderr)
		if len(match) != 2 {
			t.Fatalf("decompression refusal did not report observed streaming bytes: %s", out.Stderr)
		}
		observed, err := strconv.Atoi(match[1])
		if err != nil {
			t.Fatal(err)
		}
		if observed <= limit || observed > limit+64*1024 || observed >= expandedSize {
			t.Fatalf("decompression refusal observed %d bytes for limit %d and full payload %d", observed, limit, expandedSize)
		}
	})
}

func testRepairExplicitDefault(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho should-not-run\n")
	marker, fakebin := filepath.Join(f.Root, "repair-child-marker"), filepath.Join(f.Root, "repair-child-bin")
	realNode, err := exec.LookPath("node")
	if err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, filepath.Join(fakebin, "node"), "#!/bin/sh\n: > \"$BENCH_TEST_CHILD_MARKER\"\nexec \""+realNode+"\" \"$@\"\n")
	if err := os.Chmod(filepath.Join(fakebin, "node"), 0o755); err != nil {
		t.Fatal(err)
	}
	out := f.BenchEnv(map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": "", "BENCH_TEST_CHILD_MARKER": marker, "PATH": fakebin + string(os.PathListSeparator) + os.Getenv("PATH")}, "version")
	out.RequireExit(127)
	out.RequireContains(out.Stderr, "run bench repair")
	requirePathAbsent(t, marker, "explicit-default resolution started repair child")
	if registry.Hits() != 0 {
		t.Fatalf("explicit-default resolution started repair: hits=%d", registry.Hits())
	}
}

func testRepairSubcommand(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\nprintf 'explicit:%s\\n' \"$1\"\n")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": ""}
	repair := f.BenchEnv(env, "repair")
	repair.RequireExit(0)
	out := f.BenchEnv(env, "version")
	out.RequireExit(0)
	if strings.TrimSpace(out.Stdout) != "explicit:version" {
		t.Fatalf("explicitly repaired binary did not run: %q", out.Stdout)
	}
}

func testRepairArguments(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	for _, args := range [][]string{{"repair", "extra"}, {"repair", "--unknown"}, {"repair", "--prune", "extra"}} {
		out := f.BenchEnv(map[string]string{"BENCH_KIT": kit, "BENCH_REPAIR": ""}, args...)
		out.RequireExit(2)
		out.RequireContains(out.Stderr, "usage: bench repair [--prune]")
	}
}

func testRepairPrune(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	current := filepath.Join(f.Env["BENCH_HOME"], "cache", "bin", "9.8.7", binaryRepairPlatformSuffix(t), "bench")
	stale := []string{
		filepath.Join(f.Env["BENCH_HOME"], "cache", "bin", "9.8.6", "linux-x64", "bench"),
		filepath.Join(f.Env["BENCH_HOME"], "cache", "bin", "9.8.7", "other-platform", "bench"),
		filepath.Join(f.Env["BENCH_HOME"], "cache", "bin", "orphan"),
	}
	contract.WriteFileAbs(t, current, "current")
	for _, path := range stale {
		contract.WriteFileAbs(t, path, "stale")
	}
	env := map[string]string{"BENCH_KIT": kit, "BENCH_REPAIR": ""}
	first := f.BenchEnv(env, "repair", "--prune")
	first.RequireExit(0)
	first.RequireContains(first.Stderr, "repair prune: removed")
	if _, err := os.Stat(current); err != nil {
		t.Fatalf("prune removed current target: %v", err)
	}
	for _, path := range stale {
		requirePathAbsent(t, path, "prune preserved stale entry")
	}
	second := f.BenchEnv(env, "repair", "--prune")
	second.RequireExit(0)
	second.RequireContains(second.Stderr, "no stale cache entries")
	if err := os.RemoveAll(filepath.Join(f.Env["BENCH_HOME"], "cache")); err != nil {
		t.Fatal(err)
	}
	absent := f.BenchEnv(env, "repair", "--prune")
	absent.RequireExit(0)
	absent.RequireContains(absent.Stderr, "no stale cache entries")
	corruptPlatform := filepath.Join(f.Env["BENCH_HOME"], "cache", "bin", "9.8.7", binaryRepairPlatformSuffix(t))
	contract.WriteFileAbs(t, corruptPlatform, "corrupt-current-platform")
	corrupt := f.BenchEnv(env, "repair", "--prune")
	corrupt.RequireExit(0)
	corrupt.RequireContains(corrupt.Stderr, "repair prune: removed 9.8.7/"+binaryRepairPlatformSuffix(t))
	requirePathAbsent(t, corruptPlatform, "prune preserved malformed current-platform file")
	outside := filepath.Join(f.Root, "outside-current-platform")
	contract.WriteFileAbs(t, outside, "outside")
	if err := os.Symlink(outside, corruptPlatform); err != nil {
		t.Fatal(err)
	}
	symlink := f.BenchEnv(env, "repair", "--prune")
	symlink.RequireExit(0)
	symlink.RequireContains(symlink.Stderr, "no stale cache entries")
	if info, err := os.Lstat(corruptPlatform); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("prune followed or removed current-platform symlink: info=%v err=%v", info, err)
	}
	if data, err := os.ReadFile(outside); err != nil || string(data) != "outside" {
		t.Fatalf("prune changed current-platform symlink target: %q %v", data, err)
	}
}
