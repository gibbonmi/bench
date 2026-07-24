package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
)

func hostileArchive(t *testing.T, members int, memberSize int64) []byte {
	t.Helper()
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	tw := tar.NewWriter(gz)
	chunk := make([]byte, memberSize)
	for i := 0; i < members; i++ {
		if err := tw.WriteHeader(&tar.Header{Name: fmt.Sprintf("package/member-%05d", i), Mode: 0o644, Size: memberSize, Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(chunk); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return compressed.Bytes()
}

func TestArchiveRejectsAggregateBudgets(t *testing.T) {
	for _, test := range []struct {
		name string
		data []byte
		want string
	}{
		{name: "compressed bytes", data: make([]byte, (128<<20)+1), want: "compressed size exceeds inspection limit"},
		{name: "member count", data: hostileArchive(t, 10_001, 1), want: "member count exceeds inspection limit"},
		{name: "expanded bytes", data: hostileArchive(t, 65, 1<<20), want: "expanded size exceeds inspection limit"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validateTarballForTesting(test.data)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("archive error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestArchiveRejectsMemberLargerThanInspectionLimit(t *testing.T) {
	t.Cleanup(setArchiveMemberLimitForTesting(4))
	data := hostileArchive(t, 1, 5)
	err := validateTarballForTesting(data)
	if err == nil || !strings.Contains(err.Error(), "exceeds inspection limit") {
		t.Fatalf("oversize member error = %v", err)
	}
}

func TestEvidencePromotionKeepsPriorVerdictAtCanonicalPathOnSwapFailure(t *testing.T) {
	root := t.TempDir()
	old := filepath.Join(root, "dist", "preflight")
	if err := os.MkdirAll(old, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "manifest.json"), []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(setExchangeForTesting(func(_, _ string) error { return os.ErrPermission }))
	results := []Result{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}
	manifest := Manifest{SchemaVersion: 1, Mode: ModeVerify, Scope: ScopeFocused, Status: StatusGreen, Phases: phaseSummaries(results)}
	if err := PromoteEvidence(root, ModeVerify, results, manifest); err == nil {
		t.Fatal("injected promotion failure passed")
	}
	data, err := os.ReadFile(filepath.Join(old, "manifest.json"))
	if err != nil || string(data) != "old\n" {
		t.Fatalf("canonical prior verdict lost: %q %v", data, err)
	}
}

func TestEvidenceSIGKILLAtPromotionBoundaryKeepsOldOrNew(t *testing.T) {
	if runtime.GOOS == "windows" {
		capability.Capability(t, capability.Signal, "SIGKILL is POSIX")
	}
	if point := os.Getenv("BENCH_TEST_PROMOTION_KILL"); point != "" {
		root := os.Getenv("BENCH_TEST_PROMOTION_ROOT")
		setExchangeForTesting(func(stage, target string) error {
			if point == "before" {
				_ = syscall.Kill(os.Getpid(), syscall.SIGKILL)
			}
			err := atomicExchangeForTesting(stage, target)
			if err == nil && point == "after" {
				_ = syscall.Kill(os.Getpid(), syscall.SIGKILL)
			}
			return err
		})
		results := []Result{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}
		manifest := Manifest{SchemaVersion: 1, Mode: ModeVerify, Scope: ScopeFocused, Status: StatusGreen, Phases: phaseSummaries(results)}
		_ = PromoteEvidence(root, ModeVerify, results, manifest)
		os.Exit(90)
	}
	for _, point := range []string{"before", "after"} {
		t.Run(point, func(t *testing.T) {
			root := t.TempDir()
			target := filepath.Join(root, "dist", "preflight")
			if err := os.MkdirAll(target, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(target, "manifest.json"), []byte("old\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			testBinary, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(testBinary, "-test.run=^TestEvidenceSIGKILLAtPromotionBoundaryKeepsOldOrNew$")
			cmd.Env = append(os.Environ(), "BENCH_TEST_PROMOTION_KILL="+point, "BENCH_TEST_PROMOTION_ROOT="+root)
			if err := cmd.Run(); err == nil {
				t.Fatal("promotion helper survived SIGKILL")
			}
			data, err := os.ReadFile(filepath.Join(target, "manifest.json"))
			if err != nil {
				t.Fatalf("canonical verdict missing: %v", err)
			}
			if point == "before" && string(data) != "old\n" {
				t.Fatalf("before exchange = %q", data)
			}
			if point == "after" && string(data) == "old\n" {
				t.Fatalf("after exchange retained old verdict")
			}
			results := []Result{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}
			manifest := Manifest{SchemaVersion: 1, Mode: ModeVerify, Scope: ScopeFocused, Status: StatusGreen, Phases: phaseSummaries(results)}
			if err := PromoteEvidence(root, ModeVerify, results, manifest); err != nil {
				t.Fatalf("rerun after abandoned stage: %v", err)
			}
			entries, err := os.ReadDir(filepath.Join(root, "dist"))
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 || entries[0].Name() != "preflight" {
				t.Fatalf("rerun left non-canonical generations: %v", entryNames(entries))
			}
		})
	}
}

func entryNames(entries []os.DirEntry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func TestEvidenceRejectsHostileCanonicalTargets(t *testing.T) {
	for _, kind := range []string{"file", "fifo"} {
		t.Run(kind, func(t *testing.T) {
			if kind == "fifo" && runtime.GOOS == "windows" {
				capability.Capability(t, capability.Fifo, "no FIFO")
			}
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, "dist"), 0o755); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(root, "dist", "preflight")
			if kind == "file" {
				if err := os.WriteFile(target, []byte("keep"), 0o644); err != nil {
					t.Fatal(err)
				}
			} else if err := syscall.Mkfifo(target, 0o600); err != nil {
				t.Fatal(err)
			}
			results := []Result{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}
			manifest := Manifest{SchemaVersion: 1, Mode: ModeVerify, Scope: ScopeFocused, Status: StatusGreen, Phases: phaseSummaries(results)}
			if err := PromoteEvidence(root, ModeVerify, results, manifest); err == nil {
				t.Fatalf("%s target passed", kind)
			}
			info, err := os.Lstat(target)
			if err != nil {
				t.Fatal(err)
			}
			if kind == "file" && !info.Mode().IsRegular() {
				t.Fatal("regular target replaced")
			}
			if kind == "fifo" && info.Mode()&os.ModeNamedPipe == 0 {
				t.Fatal("FIFO target replaced")
			}
		})
	}
}

func TestInitializationFailurePreservesPriorCompleteEvidence(t *testing.T) {
	root := preflightRepo(t)
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("initial exit=%d stderr=%s", code, stderr.String())
	}
	prior, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ntoolchain go1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	after, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(prior, after) {
		t.Fatal("unreadable initialization input replaced prior trusted evidence")
	}
}

func TestVulnerabilityCancellationKillsDescendantProcessGroup(t *testing.T) {
	if runtime.GOOS == "windows" {
		capability.Capability(t, capability.Signal, "process groups are POSIX")
	}
	root := preflightRepo(t)
	pidFile := filepath.Join(root, "child.pid")
	scanner := filepath.Join(root, "scanner")
	body := "#!/bin/sh\nsleep 30 &\nprintf '%s' \"$!\" > \"$BENCH_CHILD_PID\"\nwait\n"
	if err := os.WriteFile(scanner, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_VULNERABILITY", scanner)
	t.Setenv("BENCH_CHILD_PID", pidFile)
	r := &runner{root: root, mode: ModeVerify, stderr: &bytes.Buffer{}}
	if err := r.populateBaseIdentity(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan Result, 1)
	go func() { done <- r.runPhase(ctx, "vulnerability") }()
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := os.Stat(pidFile); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("scanner descendant did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	result := <-done
	if result.Status != StatusInterrupted {
		t.Fatalf("result=%+v", result)
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil || pid <= 0 {
		t.Fatalf("scanner descendant pid = %q: %v", data, err)
	}
	requireProcessGone(t, pid)
}

func requireProcessGone(t *testing.T, pid int) {
	t.Helper()
	deadline := time.NewTimer(2 * time.Second)
	defer deadline.Stop()
	probe := time.NewTicker(10 * time.Millisecond)
	defer probe.Stop()
	for {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return
		}
		select {
		case <-probe.C:
		case <-deadline.C:
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("scanner descendant %d survived cancellation", pid)
		}
	}
}

func snapshotTree(root string) ([]byte, error) {
	var snapshot bytes.Buffer
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		fmt.Fprintf(&snapshot, "%s\x00%d\x00", filepath.ToSlash(rel), len(data))
		snapshot.Write(data)
		return nil
	})
	return snapshot.Bytes(), err
}
