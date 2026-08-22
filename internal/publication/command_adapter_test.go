package publication

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// unreachableRegistry is a base URL nothing listens on. The fixture adapter
// fails its first HTTP call there, which is all a "no npm process spawned" row
// needs from it.
const unreachableRegistry = "http://127.0.0.1:1"

// npmStub is an `npm` executable prepended to PATH that logs every invocation's
// argv and answers `npm view <name>@<version> dist.integrity --json` from
// test-registered state. A registered package reports absent (npm's E404 exit)
// on its first view. It reports its integrity on every later one; that is
// the sequence a publish run observes. markLive skips straight to the live
// answer for the paths (promote, rollback) that act on an already-published
// release. Every other subcommand — dist-tag, deprecate — succeeds silently.
//
// publish counts its calls. When failPublishAt armed BENCH_NPM_STUB_FAIL_
// PUBLISH, publish exits non-zero on that call, simulating a real npm failure
// on one package of a multi-package run. That branch also undoes the live
// marker the immediately preceding view optimistically wrote. The "absent on
// first view, live after" sequence assumes the publish between them
// succeeded. Without the undo the stub would report a package live that was
// never published, and a resume could never retry it.
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
  printf '%s' "$live" > "$BENCH_NPM_STUB_STATE/last-marked"
  exit 1
fi
if [ "$1" = publish ]; then
  count=$(cat "$BENCH_NPM_STUB_STATE/publish-count" 2>/dev/null || echo 0)
  count=$((count + 1))
  printf '%s' "$count" > "$BENCH_NPM_STUB_STATE/publish-count"
  if [ -n "$BENCH_NPM_STUB_FAIL_PUBLISH" ] && [ "$count" -eq "$BENCH_NPM_STUB_FAIL_PUBLISH" ]; then
    if [ -f "$BENCH_NPM_STUB_STATE/last-marked" ]; then rm -f "$(cat "$BENCH_NPM_STUB_STATE/last-marked")"; fi
    echo "npm ERR! code E500" >&2
    exit 1
  fi
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

// failPublishAt arms the stub to exit non-zero on the nth `npm publish` call
// of the whole stub's lifetime. call is 0 to disarm for a later resume run.
// The call counter never resets, so a resume re-run continues counting where
// the failed run stopped.
func (s *npmStub) failPublishAt(t *testing.T, call int) {
	t.Helper()
	value := ""
	if call > 0 {
		value = strconv.Itoa(call)
	}
	t.Setenv("BENCH_NPM_STUB_FAIL_PUBLISH", value)
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

// TestSubmitNPMAdapterPublishesApprovedTarballs is rows R1 and R4. With
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
// refused up front. Another operation holds the lock for the whole test. A
// refusal that ran after the lock would report the held lock instead, so the
// staged diagnostic proves the refusal came first.
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

// publicationOrder is the tarball file names in the order a first publication
// publishes them: every platform package release-plan.mjs names, then the
// wrapper last.
func publicationOrder(t *testing.T, root, version string) []string {
	t.Helper()
	var order []string
	for _, kind := range []string{"platform", "wrapper"} {
		for _, record := range releasePlanArtifacts(t, root, version) {
			if record.Kind == kind {
				order = append(order, record.Name)
			}
		}
	}
	return order
}

// approvedByFile indexes the approved set by the tarball file name the publish
// argv carries. A test can then go from a publication-order position to the
// registry package name the record's transitions key by.
func approvedByFile(packages []ApprovedPackage) map[string]ApprovedPackage {
	byFile := make(map[string]ApprovedPackage, len(packages))
	for _, pkg := range packages {
		byFile[filepath.Base(pkg.FilePath)] = pkg
	}
	return byFile
}

func transitionsFor(record Record, pkg string) []Transition {
	var matched []Transition
	for _, transition := range record.Transitions {
		if transition.Package == pkg {
			matched = append(matched, transition)
		}
	}
	return matched
}

func hasTransition(record Record, pkg, action, result string) bool {
	for _, transition := range transitionsFor(record, pkg) {
		if transition.Action == action && transition.Result == result {
			return true
		}
	}
	return false
}

func loadTestRecord(t *testing.T, root string) Record {
	t.Helper()
	record, err := LoadRecord(root)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

// structuredNPMFailure matches the npm adapter's own error shape — `npm
// <argv> failed:`. This is the message that must reach the operator when
// the CLI call itself could not run.
var structuredNPMFailure = regexp.MustCompile(`npm \[[^\]]*\] failed:`)

// TestSubmitNPMAdapterWithoutNPMBinaryFails is the spec's absent-binary edge:
// with --adapter npm and no `npm` anywhere on PATH, submit exits 1 and the
// adapter's structured `npm ... failed` error surfaces as unsatisfied release
// intent — never a silent success or an unattributed crash.
func TestSubmitNPMAdapterWithoutNPMBinaryFails(t *testing.T) {
	const version = "9.9.9"
	root := approvedReleaseRoot(t, version)

	// The narrowed PATH keeps only `node`. VerifyApprovedSet shells
	// scripts/release-plan.mjs for the artifact inventory before any registry
	// call, so a PATH without node would fail the run short of the adapter.
	node, err := exec.LookPath("node")
	if err != nil {
		t.Fatalf("node is required to run the release plan: %v", err)
	}
	pathDir := t.TempDir()
	if err := os.Symlink(node, filepath.Join(pathDir, "node")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", pathDir)
	if found, err := exec.LookPath("npm"); err == nil {
		t.Fatalf("npm is still on PATH at %s; the absent-binary row needs it gone", found)
	}

	code, stdout := runRelease(t, "submit", "--root", root, "--version", version, "--profile", "public",
		"--path", "first", "--adapter", "npm", "--registry", "https://registry.example.invalid")
	if code != 1 {
		t.Fatalf("submit exit = %d, want 1; stdout:\n%s", code, stdout)
	}
	if !strings.Contains(stdout, "error: unsatisfied release intent") {
		t.Fatalf("stdout does not report unsatisfied release intent:\n%s", stdout)
	}
	if !structuredNPMFailure.MatchString(stdout) {
		t.Fatalf("stdout does not carry the adapter's structured npm failure:\n%s", stdout)
	}
}

// TestSubmitNPMAdapterResumesAfterMidSequencePublishFailure is the
// resumability guarantee the workflow's always-upload of the record depends
// on. An npm publish that fails on an interior platform package exits 1 with a
// durable record naming the published prefix and the failure.
//
// A re-run against a then-healthy npm resumes. Already-published packages
// verify as resumed rather than republished, the failed one retries, and the
// run completes with the set awaiting an explicit promote.
func TestSubmitNPMAdapterResumesAfterMidSequencePublishFailure(t *testing.T) {
	const version = "9.9.9"
	// The third publish call is an interior platform package: packages before
	// it are already live, packages after it were never addressed.
	const failAt = 3
	root := approvedReleaseRoot(t, version)
	stub := stubNPM(t)
	approved := registerApproved(t, stub, root, version, false)
	byFile := approvedByFile(approved)

	order := publicationOrder(t, root, version)
	if len(order) != 5 {
		t.Fatalf("publication order names %d packages, want 4 platforms plus the wrapper", len(order))
	}
	submit := []string{"submit", "--root", root, "--version", version, "--profile", "public",
		"--path", "first", "--adapter", "npm", "--registry", "https://registry.example.invalid"}

	stub.failPublishAt(t, failAt)
	code, stdout := runRelease(t, submit...)
	if code != 1 {
		t.Fatalf("failed submit exit = %d, want 1; stdout:\n%s", code, stdout)
	}
	if !structuredNPMFailure.MatchString(stdout) {
		t.Fatalf("stdout does not carry the adapter's structured npm failure:\n%s", stdout)
	}

	record := loadTestRecord(t, root)
	if record.Result != "failed" {
		t.Fatalf("record result = %q after a mid-sequence publish failure, want failed", record.Result)
	}
	for _, file := range order[:failAt-1] {
		if pkg := byFile[file]; !hasTransition(record, pkg.Name, "verify", "success") {
			t.Fatalf("record does not mark %s published before the failure: %+v", pkg.Name, record.Transitions)
		}
	}
	failedPackage := byFile[order[failAt-1]]
	if !hasTransition(record, failedPackage.Name, "publish", "failed") {
		t.Fatalf("record does not mark the failed publish of %s: %+v", failedPackage.Name, record.Transitions)
	}
	for _, file := range order[failAt:] {
		if pkg := byFile[file]; len(transitionsFor(record, pkg.Name)) != 0 {
			t.Fatalf("publication continued past the failure to %s: %+v", pkg.Name, record.Transitions)
		}
	}

	stub.failPublishAt(t, 0)
	code, stdout = runRelease(t, submit...)
	if code != 0 {
		t.Fatalf("resumed submit exit = %d, want 0; stdout:\n%s", code, stdout)
	}
	if !strings.Contains(stdout, "promote") {
		t.Fatalf("resumed submit does not report the promote next action:\n%s", stdout)
	}

	record = loadTestRecord(t, root)
	if record.Result != "in_progress" {
		t.Fatalf("record result = %q after the resume, want in_progress", record.Result)
	}
	for _, file := range order[:failAt-1] {
		if pkg := byFile[file]; !hasTransition(record, pkg.Name, "verify", "resumed") {
			t.Fatalf("resume did not treat already-published %s as resumed: %+v", pkg.Name, record.Transitions)
		}
	}

	published := map[string]int{}
	for _, argv := range stub.subcommandInvocations(t, "publish") {
		published[filepath.Base(argv[1])]++
	}
	want := map[string]int{}
	for _, file := range order {
		want[file] = 1
	}
	// Only the package whose publish failed is ever published twice.
	want[order[failAt-1]] = 2
	for file, wantCount := range want {
		if published[file] != wantCount {
			t.Fatalf("npm published %s %d time(s), want %d; log: %v", file, published[file], wantCount, stub.invocations(t))
		}
	}
	if len(published) != len(want) {
		t.Fatalf("npm published %d distinct tarballs, want %d; log: %v", len(published), len(want), stub.invocations(t))
	}
}

// TestPromoteAndRollbackHonorAdapterSelection is row R7: the whole lifecycle
// addresses one registry. --adapter npm drives the npm CLI for promote and
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
