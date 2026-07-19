package publication

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestReleaseSubmitExitCodeTriple is coverage row 8: bench release submit
// returns exit 0 on success, exit 1 on an unsatisfied release intent (a
// missing/invalid registry target), and exit 2 on a usage error (a missing
// required flag) — never confusing one exit family for another.
func TestReleaseSubmitExitCodeTriple(t *testing.T) {
	version := "9.9.22"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, _ := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	t.Run("usage: missing --version", func(t *testing.T) {
		exitCode, output := runRelease(t, root, "submit", []string{"--root", root, "--profile", "public", "--registry", base}, nil)
		if exitCode != 2 {
			t.Fatalf("exit = %d, want 2 (usage) for missing --version:\n%s", exitCode, output)
		}
	})
	t.Run("unsatisfied release intent: unreachable registry", func(t *testing.T) {
		exitCode, output := runRelease(t, root, "submit", []string{"--root", root, "--version", version, "--profile", "public", "--registry", "http://127.0.0.1:1"}, nil)
		if exitCode != 1 {
			t.Fatalf("exit = %d, want 1 (unsatisfied release intent) for an unreachable registry:\n%s", exitCode, output)
		}
	})
	t.Run("success", func(t *testing.T) {
		exitCode, output := runReleaseSubmit(t, root, version, base)
		if exitCode != 0 {
			t.Fatalf("exit = %d, want 0 (success):\n%s", exitCode, output)
		}
		if !strings.Contains(output, "next_action") {
			t.Fatalf("TOON stdout does not carry next_action:\n%s", output)
		}
	})
}

// TestReleaseSubmitRerunIsIdempotentNoOp is coverage row 8's load-bearing
// rerun assertion: once a first-publication submit has already succeeded, a
// second identical run against the same state must not repeat a single live
// registry mutation (no second publish or dist-tag-add hits the request log)
// and must exit 0 as a no-op with the record left byte-identical.
func TestReleaseSubmitRerunIsIdempotentNoOp(t *testing.T) {
	version := "9.9.23"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode != 0 {
		t.Fatalf("first submit exit=%d:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "release-complete") {
		t.Fatalf("first submit did not reach next_action release-complete:\n%s", output)
	}
	firstLines := requestLines(t, requestFile)
	firstRecord, err := readRawRecord(t, root)
	if err != nil {
		t.Fatal(err)
	}

	exitCode, output = runReleaseSubmit(t, root, version, base)
	if exitCode != 0 {
		t.Fatalf("rerun submit exit=%d, want 0 (idempotent no-op):\n%s", exitCode, output)
	}
	if !strings.Contains(output, "release-complete") {
		t.Fatalf("rerun submit did not report next_action release-complete:\n%s", output)
	}

	secondLines := requestLines(t, requestFile)
	if len(secondLines) != len(firstLines) {
		t.Fatalf("rerun issued %d new registry request(s); a completed release must never re-mutate the registry:\nfirst:\n%s\nsecond:\n%s",
			len(secondLines)-len(firstLines), strings.Join(firstLines, "\n"), strings.Join(secondLines, "\n"))
	}

	secondRecord, err := readRawRecord(t, root)
	if err != nil {
		t.Fatal(err)
	}
	if secondRecord != firstRecord {
		t.Fatalf("rerun changed the durable record; a no-op resume must leave it identical:\nfirst:\n%s\nsecond:\n%s", firstRecord, secondRecord)
	}
}
