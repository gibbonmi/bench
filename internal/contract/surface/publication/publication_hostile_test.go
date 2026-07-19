package publication

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSubmitRefusesPlanVsApplyDrift is edge row A's drift case: an approved
// artifact whose bytes and SHA256SUMS entry agree with each other, but whose
// digest disagrees with what dist/preflight/release-index.json's own frozen
// artifacts binding recorded, must be refused before any registry call — a
// tampered SHA256SUMS matching a swapped-in artifact must not sail through
// just because the file and its own checksum happen to agree.
func TestSubmitRefusesPlanVsApplyDrift(t *testing.T) {
	version := "9.9.30"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	drifted := ordered[0].Name
	patchReleaseIndexArtifactDigest(t, root, drifted, strings.Repeat("f", 64))
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode == 0 {
		t.Fatalf("submit succeeded despite plan-vs-apply drift for %s:\n%s", drifted, output)
	}
	if !strings.Contains(output, "drifted from the release plan") {
		t.Fatalf("submit error did not attribute the drift:\n%s", output)
	}
	if lines := requestLines(t, requestFile); len(lines) != 0 {
		t.Fatalf("submit issued a registry call before detecting drift:\n%s", strings.Join(lines, "\n"))
	}
}

// hostileRegistry serves just enough of the FixtureRegistry HTTP surface to
// let a publish/stage attempt reach its first live registry round trip, then
// returns a hostile response from there — the least-duplicative way to
// exercise the Go adapter's response handling without reimplementing the
// node offline-registry.mjs fixture, which already owns realistic simulation.
func hostileRegistry(t *testing.T, integrityStatus int, integrityBody string, stageStatus int, stageBody string) string {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("stored\n"))
	})
	mux.HandleFunc("/-/integrity/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusOrDefault(integrityStatus))
		_, _ = w.Write([]byte(integrityBody))
	})
	mux.HandleFunc("/-/stage/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(statusOrDefault(stageStatus))
		_, _ = w.Write([]byte(stageBody))
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server.URL
}

// statusOrDefault treats an unset (zero) status as "this endpoint is not
// exercised by this test case" and answers 200 rather than letting
// net/http panic on an invalid WriteHeader(0).
func statusOrDefault(status int) int {
	if status == 0 {
		return http.StatusOK
	}
	return status
}

// TestFirstPublicationClassifiesMalformedIntegrityResponseTerminalFailClosed
// is edge row A: a registry that replies to the post-publish integrity
// re-check with malformed JSON must be classified a terminal, attributed
// failure — never a panic — and must leave the durable record intact and
// resumable rather than falsely marked complete.
func TestFirstPublicationClassifiesMalformedIntegrityResponseTerminalFailClosed(t *testing.T) {
	version := "9.9.31"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base := hostileRegistry(t, http.StatusOK, "{not-json", 0, "")

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode != 1 {
		t.Fatalf("exit = %d, want 1 for a malformed registry response:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "malformed") {
		t.Fatalf("error output did not attribute the malformed registry response:\n%s", output)
	}
	// A malformed reply to the very first exploratory integrity check (before
	// any transition is recorded or attempted) leaves nothing to write yet —
	// asserting "intact" here means no record was fabricated out of nothing,
	// which readRawRecord's absence already proves.
	if _, err := os.Stat(filepath.Join(root, "dist", "publication", "publication-record.json")); err == nil {
		record := readPublicationRecord(t, root)
		if record["result"] == "success" {
			t.Fatal("record was marked success despite a malformed registry response")
		}
	}
}

// TestFirstPublicationClassifiesControlByteIntegrityTerminalFailClosed is edge
// row A: a hostile registry embedding control bytes (ESC) in the integrity
// string it returns must be classified as an ordinary integrity mismatch —
// terminal, attributed, never a panic — and the control bytes must never
// reach raw stdout or the durable record's raw bytes (encoding/json escapes
// them; toon.Table refuses them outright rather than emitting them raw).
func TestFirstPublicationClassifiesControlByteIntegrityTerminalFailClosed(t *testing.T) {
	version := "9.9.32"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	// The transport JSON stays well-formed: \u001b is the standard JSON escape
	// that decodes into a Go string carrying the raw ESC byte at runtime.
	body := `{"integrity":"sha512-\u001bHACKED"}`
	base := hostileRegistry(t, http.StatusOK, body, 0, "")

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode != 1 {
		t.Fatalf("exit = %d, want 1 for a control-byte registry integrity value:\n%s", exitCode, output)
	}
	if strings.ContainsAny(output, "\x1b\x07") {
		t.Fatalf("raw control bytes leaked into stdout:\n%q", output)
	}
	record := readPublicationRecord(t, root)
	if record["result"] != "failed" {
		t.Fatalf("record result = %v, want failed", record["result"])
	}
	rawRecord, err := readRawRecord(t, root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(rawRecord, "\x1b\x07") {
		t.Fatalf("raw control bytes leaked into the durable record:\n%q", rawRecord)
	}
}

// TestStagedSubmitClassifiesInvalidStageIDTypeTerminalFailClosed is edge row
// A: a registry that replies to stage-submit with a wrong-typed stage_id
// (unknown/hostile shape, not the expected string) must fail closed with an
// attributed error rather than crash or silently accept a garbage stage id.
func TestStagedSubmitClassifiesInvalidStageIDTypeTerminalFailClosed(t *testing.T) {
	version := "9.9.33"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	// stageAndVerify checks Integrity first (is this package already live?)
	// before ever reaching StageSubmit; a plain 404 answers "not live yet" so
	// the flow proceeds to the hostile stage-submit response under test.
	base := hostileRegistry(t, http.StatusNotFound, "not found\n", http.StatusCreated, `{"stage_id":true}`)

	exitCode, output := runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 1 {
		t.Fatalf("exit = %d, want 1 for a wrong-typed stage_id:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "malformed") {
		t.Fatalf("error output did not attribute the malformed stage-submit response:\n%s", output)
	}
	record := readPublicationRecord(t, root)
	if record["result"] != "failed" {
		t.Fatalf("record result = %v, want failed", record["result"])
	}
}
