//go:build system

package systemtest

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// releaseFixture is a committed checkout that both release verbs accept: an exact toolchain,
// a package version, the gate phase's inputs, and a green verdict at its head.
func releaseFixture(t *testing.T) recordedPublicationRepo {
	t.Helper()
	fixture := scaffoldRecordedPublicationRepo(t)
	version := strings.TrimPrefix(strings.TrimSpace(owner.runAt(fixture.root, nil, realGo(t), "env", "GOVERSION").stdout), "go")
	writeTree(t, fixture.root, map[string]string{
		"go.mod":       "module example.invalid/release\n\ngo 1.25\n\ntoolchain go" + version + "\n",
		"package.json": "{\"version\":\"1.2.3\"}\n",
		"bin/bench.sh": "#!/bin/sh\n",
	})
	systemGit(t, fixture.root, "add", ".")
	systemGit(t, fixture.root, "commit", "-qm", "release fixture")
	systemGit(t, fixture.root, "update-ref", "refs/bench/green/main", systemGitOutput(t, fixture.root, "rev-parse", "HEAD"))
	return fixture
}

// preflightPhase runs one focused verify phase of the release preflight with the phase
// override in env.
func preflightPhase(t *testing.T, phase string, env ...string) processResult {
	t.Helper()
	fixture := releaseFixture(t)
	env = append(append(fixture.environment(fixture.home), "TMPDIR="+t.TempDir()), env...)
	return systemSelected(t, fixture.root, env, "release-preflight", "--mode", "verify", "--phase", phase)
}

// TestProcessLifetimePrepDifferential is PL95. Prep-release runs its steps in order and
// streams each step's output under the step's name. A failed step names its command and exit
// and runs no later step. A surviving descendant gets no normal-exit signal, and its protected
// group fails its step and retains the release evidence.
func TestProcessLifetimePrepDifferential(t *testing.T) {
	t.Parallel()
	for _, shape := range preservationShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			fixture := releaseFixture(t)
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			writeTree(t, fixture.root, map[string]string{"scripts/build-artifacts.sh": shapeScript(shape, "sleep", "mv", mark)})
			script := filepath.Join(fixture.root, "scripts", "build-artifacts.sh")
			// The second step's go stops the run, so the order shows in one failure.
			fakeGo(t, dir, map[string]string{"'-C '*": "echo fake-go-step\nexit 9\n"})
			callerTool(t, dir, "node", "exit 0\n")
			env := append(fixture.environment(fixture.home), toolPath(dir), "TMPDIR="+t.TempDir())
			// Prep-release starts only from a current green verdict of the same tree and oracle.
			systemGit(t, fixture.root, "add", ".")
			systemGit(t, fixture.root, "commit", "-qm", "shaped artifacts step")
			if gated := systemSelected(t, fixture.root, env, "gate"); gated.code != 0 {
				t.Fatalf("gate = (%d, %q, %q)", gated.code, gated.stdout, gated.stderr)
			}
			result := systemSelected(t, fixture.root, env, "prep-release")
			stepped := strings.Index(result.stderr, "prep-release: artifacts\n")
			canned := strings.Index(result.stderr, "canned-stdout\ncanned-stderr\n")
			next := strings.Index(result.stderr, "prep-release: conformance-ship\n")
			if result.code != 1 || stepped < 0 || canned < stepped {
				t.Fatalf("prep-release over %s = (%d, %q, %q), want the artifacts step's merged output", shape.name, result.code, result.stdout, result.stderr)
			}
			pgid := groupOf(t, mark)
			failed := "prep-release: artifacts failed: bash " + script + " "
			switch {
			case shape.exit != 0:
				if next >= 0 || !strings.Contains(result.stderr, failed) || !strings.Contains(result.stderr, "exit status 7") {
					t.Fatalf("prep-release over exit seven = %q, want the attributed artifacts failure and no later step", result.stderr)
				}
			case shape.descendant:
				// The live group keeps the release evidence protected, so the step fails on that
				// protection, no later step runs, and the evidence stays retained.
				requireGroup(t, pgid, protection.GroupPresent, "artifacts step")
				if next >= 0 || !strings.Contains(result.stderr, failed) || !strings.Contains(result.stderr, "process protection "+string(protection.Protected)) || !strings.Contains(result.stderr, "prep-release evidence stays retained") {
					t.Fatalf("prep-release over a live descendant = %q, want the attributed protection failure and retained evidence", result.stderr)
				}
			case next < canned || !strings.Contains(result.stderr[next:], "fake-go-step\n") || !strings.Contains(result.stderr, "prep-release: conformance-ship failed"):
				t.Fatalf("prep-release over %s = %q, want the next step after the artifacts step", shape.name, result.stderr)
			}
		})
	}
}

// TestProcessLifetimeExternalDifferential is PL96. An external phase keeps its merged output
// and its phase result: a zero phase is green, and a seven phase is red with its own exit. A
// surviving descendant gets no normal-exit signal, and its protected group turns the phase red.
func TestProcessLifetimeExternalDifferential(t *testing.T) {
	t.Parallel()
	for _, shape := range preservationShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			phase := callerTool(t, dir, "phase", shapeScript(shape, "sleep", "mv", mark))
			result := preflightPhase(t, "gate", "BENCH_PREFLIGHT_GATE="+phase)
			pgid := groupOf(t, mark)
			want := "release-preflight: gate green\n"
			switch {
			case shape.exit != 0:
				want = "release-preflight: gate red\n"
			case shape.descendant:
				requireGroup(t, pgid, protection.GroupPresent, "gate phase")
				want = "process protection " + string(protection.Protected)
			}
			if !strings.Contains(result.stderr, "canned-stdout\ncanned-stderr\n") || !strings.Contains(result.stderr, want) {
				t.Fatalf("preflight gate over %s = (%d, %q), want %q", shape.name, result.code, result.stderr, want)
			}
		})
	}
}

// TestProcessLifetimeScannerDifferential is PL97. The scanner phase decodes a complete event
// stream, applies the finding policy to an exit three, and fails a malformed stream or
// another exit. A surviving descendant gets no normal-exit signal, and its protected group
// fails the phase.
func TestProcessLifetimeScannerDifferential(t *testing.T) {
	t.Parallel()
	clean := `{"config":{"scanner_name":"govulncheck"}}` + "\n"
	cases := []jsonCase{
		{shape: preservationShapes[0], body: clean, want: "release-preflight: vulnerability green\n"},
		{shape: preservationShapes[1], body: clean, want: "govulncheck scanner failed"},
		{shape: preservationShapes[2], body: clean, want: "process protection " + string(protection.Protected)},
		{shape: preservationShape{name: "exit three finding", exit: 3}, body: `{"finding":{"osv":"GO-0000-0001","trace":[{"function":"F"}]}}` + "\n", want: "uncovered vulnerabilities: GO-0000-0001"},
		{shape: preservationShape{name: "malformed JSON"}, body: "{not json\n", want: "malformed govulncheck event stream"},
	}
	for _, c := range cases {
		t.Run(c.shape.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			scanner := callerTool(t, dir, "scanner", jsonShapeScript(c.shape, c.body, mark))
			result := preflightPhase(t, "vulnerability", "BENCH_PREFLIGHT_VULNERABILITY="+scanner)
			pgid := groupOf(t, mark)
			if c.shape.descendant {
				requireGroup(t, pgid, protection.GroupPresent, "scanner")
			}
			if !strings.Contains(result.stderr, "canned-stdout\ncanned-stderr\n") || !strings.Contains(result.stderr, c.want) {
				t.Fatalf("preflight scanner over %s = (%d, %q), want %q", c.shape.name, result.code, result.stderr, c.want)
			}
		})
	}
}
