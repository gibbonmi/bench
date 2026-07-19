package publication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

// registryStallEnv is the one place the interrupt tests set the fixture's
// test-only stall knob (see scripts/offline-registry.mjs). It is short enough
// to keep the tests fast but long enough that a signal sent the instant the
// trigger line appears always wins the race against the fixture's own
// response timer — the delay is a safety margin, not the determinism source;
// the request-log line is.
const registryStallEnv = "BENCH_REGISTRY_STALL_MS=350"

// waitForRequestLineContaining polls path until a line containing want
// appears, or fails the test at deadline. This is the deterministic trigger
// the interrupt tests use instead of a sleep race: the fixture writes the
// line only after it has already committed the corresponding server-side
// state (store write or tag) and is now stalling before it responds, so
// observing the line means the client is provably mid-flight.
func waitForRequestLineContaining(t *testing.T, path, want string, deadline time.Time) {
	t.Helper()
	for {
		for _, line := range requestLines(t, path) {
			if strings.Contains(line, want) {
				return
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("request log at %s never observed %q within deadline", path, want)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// runInterruptedRelease starts `bench release <releaseArgs...>`, waits for
// trigger to appear in the request log (proving the in-flight request already
// landed server-side and is now stalled awaiting response), sends SIGINT to
// the process, and returns its exit code and combined output. bin/bench.sh
// execs the resolved Go binary in place (route_binary), so the signal reaches
// the actual bench process, not an intermediate shell.
func runInterruptedRelease(t *testing.T, requestFile, trigger string, releaseArgs ...string) (exitCode int, output string) {
	t.Helper()
	benchScript := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	env := contract.IsolatedEnv(t, t.TempDir())
	cmd := exec.Command("bash", append([]string{benchScript}, releaseArgs...)...)
	cmd.Env = contract.ProcessEnv(env, nil)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("start release command: %v", err)
	}
	waitForRequestLineContaining(t, requestFile, trigger, time.Now().Add(10*time.Second))
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("signal release command: %v", err)
	}
	waitErr := cmd.Wait()
	output = buf.String()
	if waitErr == nil {
		return 0, output
	}
	var exitErr *exec.ExitError
	if asExitError(waitErr, &exitErr) {
		return exitErr.ExitCode(), output
	}
	t.Fatalf("release command did not exit cleanly: %v\n%s", waitErr, output)
	return -1, output
}

// queryRegistryIntegrity performs the same read-only integrity lookup
// FixtureRegistry.Integrity makes, directly against base, so a test can
// observe registry-side truth independent of what the interrupted client
// believed happened.
func queryRegistryIntegrity(t *testing.T, base, name, version string) (integrity string, live bool) {
	t.Helper()
	encoded := name
	if strings.HasPrefix(name, "@") {
		encoded = strings.Replace(name, "/", "%2F", 1)
	}
	resp, err := http.Get(base + "/-/integrity/" + encoded + "/" + version)
	if err != nil {
		t.Fatalf("query registry integrity: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", false
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("registry integrity query for %s@%s: %d", name, version, resp.StatusCode)
	}
	var payload struct {
		Integrity string `json:"integrity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode registry integrity response: %v", err)
	}
	return payload.Integrity, true
}

func transitionsForPackage(record map[string]any, pkg string) []map[string]any {
	all, _ := record["transitions"].([]any)
	var out []map[string]any
	for _, raw := range all {
		t, _ := raw.(map[string]any)
		if t["package"] == pkg {
			out = append(out, t)
		}
	}
	return out
}

// TestReleaseSubmitInterruptedMidPublishRecoversAtomically is coverage row
// "SIGINT during submit ... leaves no promoted partial ledger, no orphan
// candidate tag beyond recovery, and a clean idempotent resume": a SIGINT
// delivered while the first platform package's publish PUT has already
// committed server-side but not yet returned to the client must (1) leave the
// durable record either absent or a complete, parseable record whose
// per-package transitions never claim a result they did not reach, (2) leave
// the registry with only the candidate tag on the interrupted package, never
// "latest", and (3) resume idempotently: a rerun completes without
// re-publishing the package that was already live.
func TestReleaseSubmitInterruptedMidPublishRecoversAtomically(t *testing.T) {
	const version = "9.9.6"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	store := filepath.Join(t.TempDir(), "store")
	base, requestFile := startRegistryWithEnv(t, root, store, []string{registryStallEnv})

	target := ordered[0]
	targetPkg := npmPackageName(target)
	trigger := fmt.Sprintf("PUT %s@%s", targetPkg, version)

	exitCode, output := runInterruptedRelease(t, requestFile, trigger,
		"release", "submit", "--root", root, "--version", version, "--profile", "public", "--registry", base)
	if exitCode == 0 {
		t.Fatalf("interrupted submit exited 0, want non-zero:\n%s", output)
	}

	recordPath := filepath.Join(root, "dist", "publication", "publication-record.json")
	recordBytes, statErr := os.ReadFile(recordPath)
	if statErr != nil {
		if !os.IsNotExist(statErr) {
			t.Fatalf("read publication record: %v", statErr)
		}
		// Absent is an allowed outcome (no prior record). Nothing further to check.
	} else {
		var record map[string]any
		if err := json.Unmarshal(recordBytes, &record); err != nil {
			t.Fatalf("publication record is not complete, valid JSON after interrupt: %v\nraw: %s", err, recordBytes)
		}
		if record["result"] == "success" {
			t.Fatalf("interrupted submit's record falsely reports overall success: %v", record)
		}
		for _, transition := range transitionsForPackage(record, targetPkg) {
			if transition["action"] == "publish" && transition["result"] == "success" {
				t.Fatalf("publish transition for %s claims success though the client was interrupted before confirming it: %+v", targetPkg, transition)
			}
		}
	}

	// Registry truth: the interrupted package's bytes did land (the fixture
	// committed the write before stalling), but promotion never ran, so
	// nothing beyond the candidate tag is set.
	if _, live := queryRegistryIntegrity(t, base, targetPkg, version); !live {
		t.Fatalf("expected the interrupted package to already be live server-side (server committed before stalling)")
	}
	for _, line := range requestLines(t, requestFile) {
		if strings.HasPrefix(line, "DIST-TAG-ADD") && strings.Contains(line, "latest") {
			t.Fatalf("interrupted submit left an orphan promotion beyond recovery: %s", line)
		}
	}
	priorPUTCount := 0
	for _, line := range requestLines(t, requestFile) {
		if strings.HasPrefix(line, "PUT "+targetPkg+"@"+version) {
			priorPUTCount++
		}
	}

	// Idempotent resume: rerun without interruption. It must complete and
	// must not re-publish the package that is already live with matching
	// integrity — an integrity-verified skip instead.
	resumeExit, resumeOutput := runReleaseSubmit(t, root, version, base)
	if resumeExit != 0 {
		t.Fatalf("resumed submit did not complete: exit=%d\n%s", resumeExit, resumeOutput)
	}
	resumedRecord := readPublicationRecord(t, root)
	if resumedRecord["result"] != "success" {
		t.Fatalf("resumed submit result = %v, want success:\n%v", resumedRecord["result"], resumedRecord)
	}
	afterPUTCount := 0
	for _, line := range requestLines(t, requestFile) {
		if strings.HasPrefix(line, "PUT "+targetPkg+"@"+version) {
			afterPUTCount++
		}
	}
	if afterPUTCount != priorPUTCount {
		t.Fatalf("resume re-published an already-live package %s: PUT count %d -> %d", targetPkg, priorPUTCount, afterPUTCount)
	}
	var resumedVerify bool
	for _, transition := range transitionsForPackage(resumedRecord, targetPkg) {
		if transition["action"] == "verify" && transition["result"] == "resumed" {
			resumedVerify = true
		}
	}
	if !resumedVerify {
		t.Fatalf("resume did not record an integrity-verified idempotent skip for %s: %+v", targetPkg, resumedRecord)
	}
}

// TestReleasePromoteInterruptedMidTagAddRecoversAtomically covers the same
// row for the promote transition: RunPromotion shares the identical
// ctx-cancellation and per-transition atomic-save plumbing as submit (same
// signal.NotifyContext wiring in command.go, same SaveRecord after each
// registry.TagAdd), so this test proves the mechanism at the "latest"
// dist-tag step rather than re-proving the plumbing already covered above.
func TestReleasePromoteInterruptedMidTagAddRecoversAtomically(t *testing.T) {
	const version = "9.9.7"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	store := filepath.Join(t.TempDir(), "store")
	base, requestFile := startRegistryWithEnv(t, root, store, []string{registryStallEnv})

	// Drive the staged path through full approval to next_action "promote":
	// this is the one path RunPromotion actually runs on (record.Result stays
	// "in_progress" until an explicit promote succeeds — the "first"/direct
	// path instead reaches terminal "success" at the end of submit itself, so
	// staging is the seam to exercise the promote transition through).
	stagedFlowReadyToPromote(t, root, version, base)

	target := ordered[0]
	targetPkg := npmPackageName(target)
	trigger := fmt.Sprintf("DIST-TAG-ADD %s latest=%s", targetPkg, version)

	exitCode, output := runInterruptedRelease(t, requestFile, trigger,
		"release", "promote", "--root", root, "--version", version, "--profile", "public", "--registry", base)
	if exitCode == 0 {
		t.Fatalf("interrupted promote exited 0, want non-zero:\n%s", output)
	}

	record := readPublicationRecord(t, root)
	if record["result"] == "success" && !allPackagesPromoted(t, record, ordered) {
		t.Fatalf("interrupted promote's record falsely reports overall success without every package promoted: %v", record)
	}
	for _, transition := range transitionsForPackage(record, targetPkg) {
		if transition["action"] == "promote" && transition["result"] == "success" {
			// A success transition for the tag-add itself is only valid if
			// the registry actually holds that tag; cross-check it.
			if !registryHoldsLatest(t, base, targetPkg, version) {
				t.Fatalf("promote transition for %s claims success though the registry never confirmed the tag: %+v", targetPkg, transition)
			}
		}
	}

	// Idempotent resume: rerun promote without interruption; it must reach
	// overall success.
	resumeExit, resumeOutput := runReleasePromote(t, root, version, base)
	if resumeExit != 0 {
		t.Fatalf("resumed promote did not complete: exit=%d\n%s", resumeExit, resumeOutput)
	}
	if readPublicationRecord(t, root)["result"] != "success" {
		t.Fatalf("resumed promote result != success")
	}
}

func allPackagesPromoted(t *testing.T, record map[string]any, ordered []artifactRecord) bool {
	t.Helper()
	promoted := map[string]bool{}
	transitions, _ := record["transitions"].([]any)
	for _, raw := range transitions {
		tr, _ := raw.(map[string]any)
		if tr["action"] == "promote" && tr["result"] == "success" {
			promoted[fmt.Sprint(tr["package"])] = true
		}
	}
	for _, a := range ordered {
		if !promoted[npmPackageName(a)] {
			return false
		}
	}
	return true
}

func registryHoldsLatest(t *testing.T, base, name, version string) bool {
	t.Helper()
	encoded := name
	if strings.HasPrefix(name, "@") {
		encoded = strings.Replace(name, "/", "%2F", 1)
	}
	resp, err := http.Get(base + "/" + encoded)
	if err != nil {
		t.Fatalf("query registry package metadata: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	var payload struct {
		DistTags map[string]string `json:"dist-tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode registry package metadata: %v", err)
	}
	return payload.DistTags["latest"] == version
}

func runReleasePromote(t *testing.T, root, version, base string) (exitCode int, output string) {
	t.Helper()
	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	env := contract.IsolatedEnv(t, t.TempDir())
	cmd := exec.Command("bash", bench, "release", "promote", "--root", root, "--version", version, "--profile", "public", "--registry", base)
	cmd.Env = contract.ProcessEnv(env, nil)
	data, err := cmd.CombinedOutput()
	if err == nil {
		return 0, string(data)
	}
	var exitErr *exec.ExitError
	if asExitError(err, &exitErr) {
		return exitErr.ExitCode(), string(data)
	}
	t.Fatalf("release promote did not run: %v\n%s", err, data)
	return -1, string(data)
}
