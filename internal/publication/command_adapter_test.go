package publication

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// unreachableRegistry is a base URL nothing listens on: the fixture adapter
// fails its first HTTP call there, which is all a "no npm process spawned" row
// needs from it.
const unreachableRegistry = "http://127.0.0.1:1"

// npmStub is an `npm` executable prepended to PATH that logs every invocation's
// argv and answers `npm view <name>@<version> dist.integrity --json` from
// test-registered state: a registered package reports absent (npm's E404 exit)
// on its first view and reports its integrity on every later one, which is the
// sequence a publish run observes. markLive skips straight to the live answer
// for the paths (promote, rollback) that act on an already-published release.
// Every other subcommand — publish, dist-tag, deprecate — succeeds silently.
type npmStub struct {
	state string
	log   string
}

func stubNPM(t *testing.T) *npmStub {
	t.Helper()
	dir := t.TempDir()
	stub := &npmStub{state: t.TempDir(), log: filepath.Join(t.TempDir(), "argv.log")}
	script := `#!/bin/sh
printf '%s\n' "$*" >> "$BENCH_NPM_STUB_LOG"
if [ "$1" = view ]; then
  key=$(printf '%s' "$2" | tr '/@' '__')
  integrity="$BENCH_NPM_STUB_STATE/int-$key"
  live="$BENCH_NPM_STUB_STATE/live-$key"
  [ -f "$integrity" ] || exit 1
  if [ -f "$live" ]; then cat "$integrity"; exit 0; fi
  : > "$live"
  exit 1
fi
exit 0
`
	if err := os.WriteFile(filepath.Join(dir, "npm"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stub.log, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_NPM_STUB_LOG", stub.log)
	t.Setenv("BENCH_NPM_STUB_STATE", stub.state)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return stub
}

func (s *npmStub) key(name, version string) string {
	return strings.NewReplacer("/", "_", "@", "_").Replace(name + "@" + version)
}

func (s *npmStub) register(t *testing.T, name, version, integrity string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(s.state, "int-"+s.key(name, version)), []byte(`"`+integrity+`"`), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (s *npmStub) markLive(t *testing.T, name, version, integrity string) {
	t.Helper()
	s.register(t, name, version, integrity)
	if err := os.WriteFile(filepath.Join(s.state, "live-"+s.key(name, version)), nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

// invocations returns one argv per recorded `npm` call, in call order.
func (s *npmStub) invocations(t *testing.T) [][]string {
	t.Helper()
	data, err := os.ReadFile(s.log)
	if err != nil {
		t.Fatal(err)
	}
	var argvs [][]string
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		argvs = append(argvs, strings.Fields(line))
	}
	return argvs
}

// subcommandInvocations narrows the log to the calls whose argv starts with the
// given words (e.g. "publish", or "dist-tag" "add").
func (s *npmStub) subcommandInvocations(t *testing.T, words ...string) [][]string {
	t.Helper()
	var matched [][]string
	for _, argv := range s.invocations(t) {
		if len(argv) < len(words) {
			continue
		}
		hit := true
		for i, word := range words {
			if argv[i] != word {
				hit = false
				break
			}
		}
		if hit {
			matched = append(matched, argv)
		}
	}
	return matched
}

// registerApproved teaches the stub every approved package's real integrity, so
// a publish run's post-publish verification matches the approved local tarball.
func registerApproved(t *testing.T, stub *npmStub, root, version string, live bool) []ApprovedPackage {
	t.Helper()
	_, packages, err := VerifyApprovedSet(root, version)
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range packages {
		if live {
			stub.markLive(t, pkg.Name, pkg.Version, pkg.Integrity)
			continue
		}
		stub.register(t, pkg.Name, pkg.Version, pkg.Integrity)
	}
	return packages
}

// seedRecord writes the in-progress record promote and rollback require before
// they will address a registry at all.
func seedRecord(t *testing.T, root, version, profile string) {
	t.Helper()
	releaseIndexSHA256, packages, err := VerifyApprovedSet(root, version)
	if err != nil {
		t.Fatal(err)
	}
	record := Record{ReleaseIndexSHA256: releaseIndexSHA256, Path: "first", Profile: profile, Result: "in_progress"}
	for _, pkg := range packages {
		record.Provenance = append(record.Provenance, Provenance{Package: pkg.Name, SHA256: pkg.SHA256, Kind: pkg.Kind})
	}
	if err := SaveRecord(root, record); err != nil {
		t.Fatal(err)
	}
}

func runRelease(t *testing.T, args ...string) (code int, stdout string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = Command(args, &out, &errOut)
	return code, out.String()
}

// TestSubmitNPMAdapterPublishesApprovedTarballs is rows R1 and R4: with
// --adapter npm the command drives the real npm CLI, publishing every approved
// tarball under the candidate tag with the public profile's --access public.
func TestSubmitNPMAdapterPublishesApprovedTarballs(t *testing.T) {
	const version = "9.9.9"
	root := approvedReleaseRoot(t, version)
	stub := stubNPM(t)
	registerApproved(t, stub, root, version, false)

	code, stdout := runRelease(t, "submit", "--root", root, "--version", version, "--profile", "public",
		"--path", "first", "--adapter", "npm", "--registry", "https://registry.example.invalid")
	if code != 0 {
		t.Fatalf("submit exit = %d, want 0; stdout:\n%s", code, stdout)
	}

	var wantNames []string
	for _, record := range releasePlanArtifacts(t, root, version) {
		if record.Kind == "platform" {
			wantNames = append(wantNames, record.Name)
		}
	}
	if len(wantNames) != 4 {
		t.Fatalf("release plan named %d platform packages, want 4", len(wantNames))
	}
	for _, record := range releasePlanArtifacts(t, root, version) {
		if record.Kind == "wrapper" {
			wantNames = append(wantNames, record.Name)
		}
	}

	published := stub.subcommandInvocations(t, "publish")
	if len(published) != len(wantNames) {
		t.Fatalf("npm publish invocations = %d, want %d; log: %v", len(published), len(wantNames), stub.invocations(t))
	}
	for i, argv := range published {
		if got := filepath.Base(argv[1]); got != wantNames[i] {
			t.Fatalf("publish %d uploaded %q, want the approved tarball %q", i, got, wantNames[i])
		}
		joined := strings.Join(argv, " ")
		if !strings.Contains(joined, "--tag "+CandidateTag(version)) {
			t.Fatalf("publish %d argv %q does not carry the candidate tag", i, joined)
		}
		if !strings.Contains(joined, "--access public") {
			t.Fatalf("publish %d argv %q does not carry --access public", i, joined)
		}
	}
}

// TestSubmitDefaultAdapterNeverSpawnsNPM is row R2: an unspecified --adapter
// stays on the fixture registry, so no npm process is spawned at all.
func TestSubmitDefaultAdapterNeverSpawnsNPM(t *testing.T) {
	const version = "9.9.9"
	root := approvedReleaseRoot(t, version)
	stub := stubNPM(t)

	// The fixture adapter fails against an unreachable base; the row is about
	// which process the run addresses, not about the release succeeding.
	runRelease(t, "submit", "--root", root, "--version", version, "--profile", "public",
		"--path", "first", "--registry", unreachableRegistry)

	if got := stub.invocations(t); len(got) != 0 {
		t.Fatalf("default submit spawned npm %d time(s): %v", len(got), got)
	}
}

// TestSubmitUnknownAdapterIsUsageError is row R3: a typo'd adapter is refused
// with usage (exit 2), never a silent fallback to either registry.
func TestSubmitUnknownAdapterIsUsageError(t *testing.T) {
	code, stdout := runRelease(t, "submit", "--version", "9.9.9", "--profile", "public", "--adapter", "bogus")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stdout:\n%s", code, stdout)
	}
	if !strings.Contains(stdout, "usage:") || !strings.Contains(stdout, "--adapter bogus") {
		t.Fatalf("stdout %q does not refuse the --adapter value on a usage line", stdout)
	}
}

// TestSubmitNPMAdapterProvenanceIsOptIn is row R5: --provenance reaches the
// publish argv only when it was passed.
func TestSubmitNPMAdapterProvenanceIsOptIn(t *testing.T) {
	const version = "9.9.9"
	for _, testCase := range []struct {
		name  string
		args  []string
		wants bool
	}{
		{name: "flag passed", args: []string{"--provenance"}, wants: true},
		{name: "flag absent", wants: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := approvedReleaseRoot(t, version)
			stub := stubNPM(t)
			registerApproved(t, stub, root, version, false)

			args := append([]string{"submit", "--root", root, "--version", version, "--profile", "public",
				"--path", "first", "--adapter", "npm", "--registry", "https://registry.example.invalid"}, testCase.args...)
			if code, stdout := runRelease(t, args...); code != 0 {
				t.Fatalf("submit exit = %d, want 0; stdout:\n%s", code, stdout)
			}

			published := stub.subcommandInvocations(t, "publish")
			if len(published) == 0 {
				t.Fatal("no npm publish invocation recorded")
			}
			for _, argv := range published {
				got := false
				for _, arg := range argv {
					if arg == "--provenance" {
						got = true
					}
				}
				if got != testCase.wants {
					t.Fatalf("publish argv %v carries --provenance = %v, want %v", argv, got, testCase.wants)
				}
			}
		})
	}
}

// TestSubmitStagedNPMAdapterRefusesBeforeTheLock is row R6: the staged path is
// refused up front. The lock is held by another operation for the whole test —
// a refusal that ran after the lock would report the held lock instead, so the
// staged diagnostic is proof the refusal came first.
func TestSubmitStagedNPMAdapterRefusesBeforeTheLock(t *testing.T) {
	const version = "9.9.9"
	root := approvedReleaseRoot(t, version)
	stub := stubNPM(t)
	registerApproved(t, stub, root, version, false)
	release, err := AcquireReleaseLock(root)
	if err != nil {
		t.Fatal(err)
	}
	defer release()

	code, stdout := runRelease(t, "submit", "--root", root, "--version", version, "--profile", "public",
		"--path", "staged", "--adapter", "npm", "--registry", "https://registry.example.invalid")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout:\n%s", code, stdout)
	}
	if strings.Contains(stdout, "lock held at") {
		t.Fatalf("staged refusal ran after the release lock: %q", stdout)
	}
	if !strings.Contains(stdout, "staged") || !strings.Contains(stdout, "--path first") {
		t.Fatalf("stdout %q does not name the staged capability and the --path first alternative", stdout)
	}
	if got := stub.invocations(t); len(got) != 0 {
		t.Fatalf("the refused staged run spawned npm %d time(s): %v", len(got), got)
	}
}

// TestPromoteAndRollbackHonorAdapterSelection is row R7: the whole lifecycle
// addresses one registry — --adapter npm drives the npm CLI for promote and
// rollback too, and their default stays the fixture.
func TestPromoteAndRollbackHonorAdapterSelection(t *testing.T) {
	const version = "9.9.9"
	for _, testCase := range []struct {
		subcommand string
		npmWords   []string
	}{
		{subcommand: "promote", npmWords: []string{"dist-tag", "add"}},
		{subcommand: "rollback", npmWords: []string{"dist-tag", "rm"}},
	} {
		t.Run(testCase.subcommand+" npm adapter", func(t *testing.T) {
			root := approvedReleaseRoot(t, version)
			stub := stubNPM(t)
			registerApproved(t, stub, root, version, true)
			seedRecord(t, root, version, "public")

			code, stdout := runRelease(t, testCase.subcommand, "--root", root, "--version", version,
				"--profile", "public", "--adapter", "npm", "--registry", "https://registry.example.invalid")
			if code != 0 {
				t.Fatalf("%s exit = %d, want 0; stdout:\n%s", testCase.subcommand, code, stdout)
			}
			if got := stub.subcommandInvocations(t, testCase.npmWords...); len(got) == 0 {
				t.Fatalf("%s --adapter npm never ran npm %v; log: %v", testCase.subcommand, testCase.npmWords, stub.invocations(t))
			}
		})

		t.Run(testCase.subcommand+" default adapter", func(t *testing.T) {
			root := approvedReleaseRoot(t, version)
			stub := stubNPM(t)
			seedRecord(t, root, version, "public")

			runRelease(t, testCase.subcommand, "--root", root, "--version", version,
				"--profile", "public", "--registry", unreachableRegistry)

			if got := stub.invocations(t); len(got) != 0 {
				t.Fatalf("default %s spawned npm %d time(s): %v", testCase.subcommand, len(got), got)
			}
		})
	}
}
