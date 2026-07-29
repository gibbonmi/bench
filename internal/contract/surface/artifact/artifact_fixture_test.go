package artifact

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

func assertPackedSetupForwarding(t *testing.T, dir, wrapper, shim string, env map[string]string) {
	t.Helper()
	real := wrapper + ".real"
	if err := os.Rename(wrapper, real); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Rename(real, wrapper) })
	stub := "#!/bin/sh\nprintf 'installed-setup:%s:%s\\n' \"$1\" \"$2\"\nexit 23\n"
	if err := os.WriteFile(wrapper, []byte(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	probe := exec.Command(shim, "setup", "a b")
	probe.Dir, probe.Env = dir, contract.ProcessEnv(nil, env)
	out, err := probe.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 23 || string(out) != "installed-setup:setup:a b\n" {
		t.Fatalf("stable shim setup forwarding = %v, %q", err, out)
	}
	if err := os.Remove(wrapper); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(real, wrapper); err != nil {
		t.Fatal(err)
	}
}

func assertPackedEntrySurfaceIdentity(t *testing.T, repo string, env map[string]string, versionOut string) {
	t.Helper()
	hookOut := runLifecycle(t, repo, env, "bash", filepath.Join(repo, ".bench", "hooks", "session-start.sh"))
	if !strings.Contains(hookOut, filepath.Join(repo, ".bench", "bin", "bench.sh")) {
		t.Fatalf("session hook did not identify linked launcher:\n%s", hookOut)
	}
	stop := exec.Command("bash", filepath.Join(repo, ".bench", "hooks", "stop.sh"))
	stop.Dir, stop.Env, stop.Stdin = repo, contract.ProcessEnv(nil, env), strings.NewReader("{}\n")
	if out, err := stop.CombinedOutput(); err != nil || strings.Contains(string(out), "GLOBAL-RUNTIME") {
		t.Fatalf("stop hook escaped linked launcher: %v\n%s", err, out)
	}
	stubDir := filepath.Join(repo, "provider stubs [*]")
	if err := os.MkdirAll(stubDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, provider := range []string{"claude", "codex", "opencode"} {
		body := "#!/bin/sh\nprintf 'provider:" + provider + ":%s\\n' \"$*\"\n"
		if err := os.WriteFile(filepath.Join(stubDir, provider), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		adapterEnv := maps.Clone(env)
		adapterEnv["PATH"] = stubDir + string(os.PathListSeparator) + os.Getenv("PATH")
		out := runLifecycle(t, repo, adapterEnv, "bash", filepath.Join(repo, ".bench", "adapters", provider), "identity prompt")
		if !strings.Contains(out, "provider:"+provider+":") {
			t.Fatalf("packed %s adapter escaped linked launcher/provider: %s", provider, out)
		}
	}
	if !strings.Contains(versionOut, runtime.GOOS+"/") {
		t.Fatalf("selected target identity missing from installed version: %s", versionOut)
	}
}

func runPackedFreshClone(t *testing.T, repo, _, shim, version string) {
	t.Helper()
	runLifecycle(t, repo, nil, "git", "add", "-A")
	runLifecycle(t, repo, nil, "git", "commit", "-qm", "linked state")
	clone := filepath.Join(t.TempDir(), "committed fresh clone [*]")
	runLifecycle(t, filepath.Dir(clone), nil, "git", "clone", "-q", repo, clone)
	home := filepath.Join(t.TempDir(), "empty cache [*]")
	bin := filepath.Join(home, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{"HOME": home, "BENCH_HOME": home, "PATH": bin + string(os.PathListSeparator) + "/usr/bin:/bin"}
	runLifecycle(t, clone, env, shim, "link")
	runLifecycle(t, clone, env, shim, "init")
	manifest, err := os.ReadFile(filepath.Join(clone, ".bench", "link-manifest.tsv"))
	if err != nil || !bytes.Contains(manifest, []byte("#kit\t"+version+"\n")) {
		t.Fatalf("fresh clone maintenance lost installed identity: %q, %v", manifest, err)
	}
	if entries, err := os.ReadDir(filepath.Join(home, "cache")); err == nil && len(entries) != 0 {
		t.Fatalf("maintenance unexpectedly populated empty runtime cache: %v", entries)
	}
}

func runLifecycle(t *testing.T, dir string, overrides map[string]string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir, cmd.Env = dir, contract.ProcessEnv(nil, overrides)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
	return string(out)
}

func assertPlannedArtifactNames(t *testing.T, root, output string) {
	t.Helper()
	var pkg struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &pkg)
	planned := strings.Fields(runLifecycle(t, root, nil, "node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "artifact-names", pkg.Version))
	want := make(map[string]bool, len(planned))
	for _, name := range planned {
		want[name] = true
	}
	entries, err := os.ReadDir(output)
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string]bool, len(entries))
	for _, entry := range entries {
		got[entry.Name()] = true
	}
	if len(want) != len(planned) || !maps.Equal(got, want) {
		t.Fatalf("artifact names = %v, want staged release-plan artifact-names = %v", got, want)
	}
}

func assertPromotedReproducibility(t *testing.T, output string) {
	t.Helper()
	path := filepath.Join(filepath.Dir(output), "reproducibility.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("promoted reproducibility.json is missing: %v", err)
	}
	var record struct {
		SchemaVersion int    `json:"schema_version"`
		Status        string `json:"status"`
		Builds        int    `json:"builds"`
	}
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf("promoted reproducibility.json is malformed: %v", err)
	}
	if record.SchemaVersion != 1 || record.Status != "green" || record.Builds != 2 {
		t.Fatalf("promoted reproducibility.json = %+v, want schema 1 green two-build evidence", record)
	}
}

type artifactSourceOption int

const (
	includeFirstNonHostArtifactTarget artifactSourceOption = 1
	stageHostlessArtifactPlan         artifactSourceOption = 2
)

func committedHostileArtifactSource(t *testing.T, root string, options ...artifactSourceOption) string {
	t.Helper()
	return contract.NarrowReleasePlan(t, root, artifactSourceNarrowing(t, options...))
}

func committedHostileArtifactSourceIn(t *testing.T, directory, root string, options ...artifactSourceOption) string {
	t.Helper()
	return contract.NarrowReleasePlanIn(t, directory, root, artifactSourceNarrowing(t, options...))
}

func artifactSourceNarrowing(t *testing.T, options ...artifactSourceOption) func(contract.ReleasePlanTargets) []contract.ReleaseTarget {
	t.Helper()
	return func(matrix contract.ReleasePlanTargets) []contract.ReleaseTarget {
		selected := append(make([]contract.ReleaseTarget, 0, 2), matrix.Host...)
		if len(options) != 0 && options[0] == stageHostlessArtifactPlan {
			selected = selected[:0]
		}
		if len(selected) == 0 {
			capability.Environment(t, fmt.Sprintf("artifact contract tests require release plan target for host %s/%s", matrix.GOOS, matrix.GOArch))
		}
		if len(options) != 0 && options[0] == includeFirstNonHostArtifactTarget {
			for _, target := range matrix.All {
				if target.GOOS != matrix.GOOS || target.GOArch != matrix.GOArch {
					selected = append(selected, target)
					break
				}
			}
			if len(selected) != 2 {
				capability.Environment(t, fmt.Sprintf("artifact matrix breadth requires a non-host release plan target alongside host %s/%s", matrix.GOOS, matrix.GOArch))
			}
		}
		return selected
	}
}

func assertSpecialFileArtifactFailure(t *testing.T, root, output string) {
	t.Helper()
	broken := committedHostileArtifactSource(t, root)
	if err := os.Remove(filepath.Join(broken, "LICENSE")); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(broken, "LICENSE"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(output, "promoted-sentinel")
	if err := os.WriteFile(sentinel, []byte("owned"), 0o644); err != nil {
		t.Fatal(err)
	}
	bad := contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(root, "scripts", "build-artifacts.sh"), broken, output)
	if bad.ExitCode == 0 {
		t.Fatal("special-file artifact builder unexpectedly succeeded")
	}
	if !strings.Contains(bad.Stderr, "required release evidence source is invalid: LICENSE") {
		t.Fatalf("special-file diagnostic was not distinct:\n%s", bad.Stderr)
	}
	if got, err := os.ReadFile(sentinel); err != nil || string(got) != "owned" {
		t.Fatalf("failed build changed promoted artifacts: %q, %v", got, err)
	}
}

func copyPreparedArtifactGeneration(t *testing.T, source string) string {
	t.Helper()
	prepared := filepath.Join(t.TempDir(), "prepared artifact generation")
	if err := os.Mkdir(prepared, 0o755); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("cp", "-a", source+string(os.PathSeparator)+".", prepared)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("copy prepared artifact generation: %v\n%s", err, output)
	}
	// The private root must stay movable while cp -a preserves the shared
	// tarballs' read-only modes.
	if err := os.Chmod(prepared, 0o755); err != nil {
		t.Fatal(err)
	}
	return prepared
}

func promotionTestEnv(prepared, ready string) []string {
	return append(os.Environ(), "BENCH_TEST_PREPARED_ARTIFACTS="+prepared, "BENCH_TEST_PROMOTION_READY_FILE="+ready)
}

func runArtifactBuildThroughPromotionSeam(t *testing.T, command *exec.Cmd, ready string) ([]byte, error) {
	t.Helper()
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				if _, err := os.Stat(ready); err == nil {
					_ = os.Remove(ready)
					return
				}
			}
		}
	}()
	output, err := command.CombinedOutput()
	close(stop)
	<-done
	return output, err
}

// awaitArtifactPromotionSeam blocks until the build has parked at the promotion seam. The
// seam sits above every move the promotion makes, so nothing has been promoted yet when
// this returns.
func awaitArtifactPromotionSeam(t *testing.T, cmd *exec.Cmd, ready string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			return
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			t.Fatal("builder did not reach deterministic promotion seam")
		}
		time.Sleep(25 * time.Millisecond)
	}
}

// interruptArtifactPromotion signals a build parked at the promotion seam, so the failure
// lands inside the promotion block rather than during generation.
func interruptArtifactPromotion(t *testing.T, source, prepared, output string) {
	t.Helper()
	ready := filepath.Join(t.TempDir(), "promotion-ready")
	cmd := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
	cmd.Env = promotionTestEnv(prepared, ready)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	awaitArtifactPromotionSeam(t, cmd, ready)
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err == nil {
		t.Fatal("interrupted builder exited successfully")
	}
}

func assertInterruptedArtifactPromotion(t *testing.T, source, prepared, output string, wantFiles int) {
	t.Helper()
	prior := promotedArtifactDigests(t, output)
	interruptArtifactPromotion(t, source, prepared, output)
	files, err := os.ReadDir(output)
	if err != nil || len(files) != wantFiles {
		t.Fatalf("promotion interruption left partial/absent set: files=%d err=%v", len(files), err)
	}
	if after := promotedArtifactDigests(t, output); !maps.Equal(after, prior) {
		t.Fatalf("promotion interruption changed prior-generation bytes: got=%v want=%v", after, prior)
	}
	if _, err := os.Stat(output + ".lock"); !os.IsNotExist(err) {
		t.Fatalf("SIGINT left artifact output lock: %v", err)
	}
	stages, err := filepath.Glob(filepath.Join(filepath.Dir(output), ".bench-artifacts.*"))
	if err != nil || len(stages) != 0 {
		t.Fatalf("SIGINT left artifact staging directories: %v, %v", stages, err)
	}
	rerunReady := filepath.Join(t.TempDir(), "rerun-ready")
	rerun := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
	rerun.Env = promotionTestEnv(prepared, rerunReady)
	if rerunOutput, err := runArtifactBuildThroughPromotionSeam(t, rerun, rerunReady); err != nil {
		t.Fatalf("idempotent rerun after SIGINT failed: %v\n%s", err, rerunOutput)
	}
	files, err = os.ReadDir(output)
	if err != nil || len(files) != wantFiles {
		t.Fatalf("rerun after SIGINT left incomplete output: files=%d err=%v", len(files), err)
	}
	if after := promotedArtifactDigests(t, output); !maps.Equal(after, prior) {
		t.Fatalf("rerun after SIGINT changed complete output: got=%v want=%v", after, prior)
	}
}

func promotedArtifactDigests(t *testing.T, directory string) map[string]string {
	t.Helper()
	digests, err := promotedArtifactDigestMap(directory)
	if err != nil {
		t.Fatal(err)
	}
	return digests
}

func promotedArtifactDigestMap(directory string) (map[string]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	digests := make(map[string]string, len(entries))
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			return nil, fmt.Errorf("promoted artifact entry is not regular: %s", entry.Name())
		}
		data, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			return nil, err
		}
		digests[entry.Name()] = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	return digests, nil
}
