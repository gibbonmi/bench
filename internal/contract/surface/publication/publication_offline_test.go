package publication

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

// artifactRecord mirrors one row of `release-plan.mjs artifact-records` — the
// exact set of file names the release plan owns. The test reads this from the
// canonical script rather than reimplementing target/version naming, so the
// naming policy has exactly one source.
type artifactRecord struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
}

// copyPublicationScripts assembles a temp root with scripts/release-plan.mjs,
// scripts/release-plan.json (the real four-platform matrix), and
// scripts/offline-registry.mjs copied from the subject repo.
func copyPublicationScripts(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	subject := contract.SubjectRoot(t)
	for _, relative := range []string{"scripts/release-plan.mjs", "scripts/release-plan.json", "scripts/offline-registry.mjs"} {
		data, err := os.ReadFile(filepath.Join(subject, relative))
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(root, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil || os.WriteFile(target, data, 0o644) != nil {
			t.Fatalf("copy fixture %s: %v", relative, err)
		}
	}
	return root
}

// releasePlanArtifacts returns the four platform artifact records followed by
// the wrapper record for version, read from the canonical release-plan.mjs —
// the one source for package naming and the target matrix.
func releasePlanArtifacts(t *testing.T, root, version string) []artifactRecord {
	t.Helper()
	out, err := exec.Command("node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "artifact-records", version).Output()
	if err != nil {
		t.Fatalf("release plan artifact-records: %v", err)
	}
	var all []artifactRecord
	if err := json.Unmarshal(out, &all); err != nil {
		t.Fatalf("decode artifact-records: %v", err)
	}
	var platforms []artifactRecord
	var wrapper *artifactRecord
	for i := range all {
		switch all[i].Kind {
		case "platform":
			platforms = append(platforms, all[i])
		case "wrapper":
			w := all[i]
			wrapper = &w
		}
	}
	if wrapper == nil || len(platforms) != 4 {
		t.Fatalf("release plan named %d platform packages and wrapper=%v, want 4 platforms + wrapper", len(platforms), wrapper != nil)
	}
	return append(platforms, *wrapper)
}

// writeApprovedSet fabricates the approved release directory: dist/preflight/
// release-index.json + SHA256SUMS binding exactly the given per-package
// contents, and dist/artifacts holding those bytes. Every record not named in
// contents gets a distinct default payload.
func writeApprovedSet(t *testing.T, root string, ordered []artifactRecord, contents map[string][]byte) {
	t.Helper()
	artifactDir := filepath.Join(root, "dist", "artifacts")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	var sums strings.Builder
	type indexArtifact struct {
		Name   string `json:"name"`
		SHA256 string `json:"sha256"`
	}
	var artifacts []indexArtifact
	for _, record := range ordered {
		data, ok := contents[record.Name]
		if !ok {
			data = []byte(record.Name + " fixture package\n")
		}
		if err := os.WriteFile(filepath.Join(artifactDir, record.Name), data, 0o644); err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(data)
		digestHex := hex.EncodeToString(digest[:])
		sums.WriteString(digestHex + "  " + record.Name + "\n")
		artifacts = append(artifacts, indexArtifact{Name: record.Name, SHA256: digestHex})
	}
	preflightDir := filepath.Join(root, "dist", "preflight")
	if err := os.MkdirAll(preflightDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Default index is a conforming full green publish-mode run for profile
	// "public" — the shape VerifyPublishAuthority requires — with an artifacts
	// binding that matches SHA256SUMS exactly, so every existing submit/promote/
	// rollback test composes on top of a plan that agrees with the approved set
	// by construction. Row-7 authority tests patch mode/scope/profile/status via
	// patchReleaseIndexAuthority; drift tests patch an artifact digest to prove
	// VerifyApprovedSet refuses a plan-vs-apply mismatch.
	index := map[string]any{
		"schema_version": 1,
		"mode":           "publish",
		"scope":          "preflight",
		"profile":        "public",
		"status":         "green",
		"marker":         "contract-fixture",
		"artifacts":      artifacts,
	}
	indexData, err := json.Marshal(index)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(preflightDir, "release-index.json"), append(indexData, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(preflightDir, "SHA256SUMS"), []byte(sums.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

// patchReleaseIndexAuthority rewrites the mode/scope/profile/status fields of
// an already-written dist/preflight/release-index.json, leaving every other
// byte (and thus VerifyApprovedSet's digest binding) alone. Used only to
// fabricate the red authority cases row 7 asserts against — the released
// index shape itself is never re-derived, only its authority fields flexed.
func patchReleaseIndexAuthority(t *testing.T, root string, overrides map[string]string) {
	t.Helper()
	path := filepath.Join(root, "dist", "preflight", "release-index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	for key, value := range overrides {
		fields[key] = value
	}
	encoded, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(encoded, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// patchReleaseIndexArtifactDigest rewrites the recorded sha256 for one named
// artifact inside an already-written release-index.json's "artifacts" array,
// without touching SHA256SUMS or the artifact bytes on disk — fabricating
// exactly the plan-vs-apply drift VerifyApprovedSet must refuse: the frozen
// plan and the approved set disagree even though each is individually
// self-consistent.
func patchReleaseIndexArtifactDigest(t *testing.T, root, artifactName, digest string) {
	t.Helper()
	path := filepath.Join(root, "dist", "preflight", "release-index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	artifacts, _ := fields["artifacts"].([]any)
	found := false
	for _, raw := range artifacts {
		entry, _ := raw.(map[string]any)
		if entry["name"] == artifactName {
			entry["sha256"] = digest
			found = true
		}
	}
	if !found {
		t.Fatalf("release-index.json names no artifact %q to patch", artifactName)
	}
	encoded, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(encoded, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// startRegistry launches the offline-registry.mjs fixture rooted at
// fixtureRoot/scripts, storing uploads under store, and returns its base URL
// and the request-log path. Callers may pre-seed store with files before or
// after calling this, since the registry rereads the directory per request.
func startRegistry(t *testing.T, fixtureRoot, store string) (base, requestFile string) {
	t.Helper()
	return startRegistryWithEnv(t, fixtureRoot, store, nil)
}

// startRegistryWithEnv is the one source startRegistry composes on top of: it
// takes extra environment variables (e.g. BENCH_REGISTRY_STALL_MS, the
// test-only stall knob offline-registry.mjs honors) for the tests that need a
// deterministic mid-request interrupt point.
func startRegistryWithEnv(t *testing.T, fixtureRoot, store string, extraEnv []string) (base, requestFile string) {
	t.Helper()
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	portFile := filepath.Join(t.TempDir(), "port")
	requestFile = filepath.Join(t.TempDir(), "requests")
	server := exec.Command("node", filepath.Join(fixtureRoot, "scripts", "offline-registry.mjs"), store, portFile, requestFile)
	if len(extraEnv) > 0 {
		server.Env = append(os.Environ(), extraEnv...)
	}
	if err := server.Start(); err != nil {
		t.Fatalf("start registry: %v", err)
	}
	t.Cleanup(func() {
		_ = server.Process.Signal(os.Interrupt)
		_ = server.Wait()
	})
	deadline := time.Now().Add(5 * time.Second)
	var port string
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(portFile); err == nil {
			port = strings.TrimSpace(string(data))
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if port == "" {
		t.Fatal("offline registry did not publish its port")
	}
	return "http://127.0.0.1:" + port, requestFile
}

func requestLines(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

func runReleaseSubmit(t *testing.T, root, version, base string) (exitCode int, output string) {
	t.Helper()
	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	env := contract.IsolatedEnv(t, t.TempDir())
	cmd := exec.Command("bash", bench, "release", "submit", "--root", root, "--version", version, "--profile", "public", "--registry", base)
	cmd.Env = contract.ProcessEnv(env, nil)
	data, err := cmd.CombinedOutput()
	if err == nil {
		return 0, string(data)
	}
	var exitErr *exec.ExitError
	if ok := asExitError(err, &exitErr); ok {
		return exitErr.ExitCode(), string(data)
	}
	t.Fatalf("release submit did not run: %v\n%s", err, data)
	return -1, string(data)
}

func asExitError(err error, target **exec.ExitError) bool {
	if ee, ok := err.(*exec.ExitError); ok {
		*target = ee
		return true
	}
	return false
}

// npmPackageName derives the registry package name for an artifact record the
// same way internal/publication does: the wrapper is "redbench", each
// platform is scoped "@redbench/<target>".
func npmPackageName(record artifactRecord) string {
	if record.Kind == "wrapper" {
		return "redbench"
	}
	return "@redbench/" + record.Target
}

func readPublicationRecord(t *testing.T, root string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "dist", "publication", "publication-record.json"))
	if err != nil {
		t.Fatalf("read publication record: %v", err)
	}
	var record map[string]any
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf("decode publication record: %v", err)
	}
	return record
}

// TestFirstPublicationPublishesPlatformsBeforeWrapperLast is coverage row 1: the
// state machine must direct-publish every platform package before the wrapper,
// each verified live before advancing.
func TestFirstPublicationPublishesPlatformsBeforeWrapperLast(t *testing.T) {
	const version = "9.9.1"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode != 0 {
		t.Fatalf("release submit exit=%d:\n%s", exitCode, output)
	}

	lines := requestLines(t, requestFile)
	platformPUT := regexp.MustCompile(`^PUT @redbench/`)
	wrapperPUT := regexp.MustCompile(`^PUT redbench@`)
	lastPlatform, wrapperIndex := -1, -1
	for i, line := range lines {
		if platformPUT.MatchString(line) {
			lastPlatform = i
		}
		if wrapperPUT.MatchString(line) {
			wrapperIndex = i
		}
	}
	if lastPlatform == -1 || wrapperIndex == -1 {
		t.Fatalf("expected both platform and wrapper PUT calls in the request log:\n%s", strings.Join(lines, "\n"))
	}
	if wrapperIndex < lastPlatform {
		t.Fatalf("wrapper was published before a platform package; request log:\n%s", strings.Join(lines, "\n"))
	}

	record := readPublicationRecord(t, root)
	if record["result"] != "success" {
		t.Fatalf("publication record result = %v, want success:\n%v", record["result"], record)
	}
}

// TestFirstPublicationResumeAcceptsOnlyExactIntegrityMatch is coverage row 3: a
// package already live at the registry is treated as complete (never
// republished) only when its integrity exactly matches the approved local
// tarball; a mismatch is terminal and the release never reaches the
// unaffected packages behind it, including the wrapper.
func TestFirstPublicationResumeAcceptsOnlyExactIntegrityMatch(t *testing.T) {
	const version = "9.9.2"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	// ordered = [platform0 (matching-live), platform1 (mismatched-live), platform2, platform3, wrapper]
	matching := ordered[0]
	mismatched := ordered[1]
	matchingBytes := []byte("approved bytes for " + matching.Name + "\n")
	mismatchedApprovedBytes := []byte("approved bytes for " + mismatched.Name + "\n")
	mismatchedLiveBytes := []byte("STALE bytes already live for " + mismatched.Name + "\n")

	writeApprovedSet(t, root, ordered, map[string][]byte{
		matching.Name:   matchingBytes,
		mismatched.Name: mismatchedApprovedBytes,
	})

	store := filepath.Join(t.TempDir(), "store")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	// Pre-seed the registry as if an earlier, interrupted run had already made
	// these two packages live: one with bytes identical to the approved local
	// tarball (a genuine resume target), one with different bytes (a corrupted
	// or stale prior publish that must never be accepted).
	if err := os.WriteFile(filepath.Join(store, matching.Name), matchingBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store, mismatched.Name), mismatchedLiveBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	base, requestFile := startRegistry(t, root, store)

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode == 0 {
		t.Fatalf("release submit passed over an integrity mismatch:\n%s", output)
	}

	lines := requestLines(t, requestFile)
	for _, line := range lines {
		if strings.HasPrefix(line, "PUT ") {
			t.Fatalf("resume issued a live PUT when the release should have stopped terminally on mismatch:\n%s", strings.Join(lines, "\n"))
		}
	}

	record := readPublicationRecord(t, root)
	if record["result"] != "failed" {
		t.Fatalf("publication record result = %v, want failed:\n%v", record["result"], record)
	}
	matchingPackage, mismatchedPackage := npmPackageName(matching), npmPackageName(mismatched)
	transitions, _ := record["transitions"].([]any)
	var matchingVerified, mismatchDetected bool
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		if transition["package"] == matchingPackage && transition["result"] == "resumed" {
			matchingVerified = true
		}
		if transition["package"] == mismatchedPackage && transition["result"] == "mismatch" {
			mismatchDetected = true
		}
		if transition["package"] == "redbench" && transition["action"] == "publish" {
			t.Fatalf("wrapper was reached after an upstream integrity mismatch should have stopped the release")
		}
	}
	if !matchingVerified {
		t.Fatalf("matching already-live package was not recorded as an idempotent resume: %+v", transitions)
	}
	if !mismatchDetected {
		t.Fatalf("mismatched already-live package was not recorded as a mismatch: %+v", transitions)
	}
}
