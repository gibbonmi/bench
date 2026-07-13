package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

const story3LocalManifest = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`

type codecKind int

const (
	codecAbsent codecKind = iota
	codecReadyRed
	codecUnavailable
	codecZeroByte
	codecNoFinalNewlineGreen
	codecNoFinalNewlineRed
	codecTrailing
	codecDuplicate
	codecUnknown
	codecWrongType
	codecWrongSchema
	codecWrongState
	codecWrongStatus
	codecWrongHash
	codecWrongTime
	codecLegacy
	codecTruncated
	codec16384
	codec16385
	codecSymlink
	codecDirectory
	codecUnreadable
)

func codecProof(id string, kind codecKind) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) { runCodecProof(t, kind) }}
}

func runCodecProof(t *testing.T, kind codecKind) {
	if kind == codecUnavailable {
		got := Inspect(t.TempDir())
		assertInspection(t, got, Unavailable, "subject unavailable", false)
		return
	}
	f := story3Fixture(t)
	f.WriteFile(".bench/gate-inputs.json", story3LocalManifest)
	if kind == codecReadyRed || kind == codecNoFinalNewlineRed {
		f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\nexit 17\n")
		f.Bench("gate").RequireExit(17)
	} else {
		f.Bench("gate").RequireExit(0)
	}
	cache := filepath.Join(story3GitDir(f), "bench-last-gate")
	data := strings.TrimSuffix(contract.ReadFileAbs(t, cache), "\n")
	wantState, wantReason, reusable := Invalid, "invalid cache record", false
	switch kind {
	case codecAbsent:
		contract.Remove(t, cache)
		wantState, wantReason = Absent, ""
	case codecReadyRed:
		wantState, wantReason = Ready, "recorded red"
	case codecZeroByte:
		writeCache(t, cache, "", 0o600)
	case codecNoFinalNewlineGreen:
		writeCache(t, cache, data, 0o600)
		wantState, wantReason, reusable = Ready, "", true
	case codecNoFinalNewlineRed:
		writeCache(t, cache, data, 0o600)
		wantState, wantReason = Ready, "recorded red"
	case codecTrailing:
		writeCache(t, cache, data+`{}`, 0o600)
	case codecDuplicate:
		writeCache(t, cache, strings.Replace(data, `"schema":1`, `"schema":1,"schema":1`, 1), 0o600)
	case codecUnknown:
		writeCache(t, cache, strings.TrimSuffix(data, "}")+`,"unknown":1}`, 0o600)
	case codecWrongType:
		writeCache(t, cache, strings.Replace(data, `"schema":1`, `"schema":"1"`, 1), 0o600)
	case codecWrongSchema:
		writeCache(t, cache, strings.Replace(data, `"schema":1`, `"schema":2`, 1), 0o600)
	case codecWrongState:
		writeCache(t, cache, strings.Replace(data, `"state":"ready"`, `"state":"other"`, 1), 0o600)
	case codecWrongStatus:
		writeCache(t, cache, strings.Replace(data, `"status":"green"`, `"status":"other"`, 1), 0o600)
	case codecWrongHash:
		writeCache(t, cache, replaceJSONField(t, data, "tree", "not-a-sha"), 0o600)
	case codecWrongTime:
		writeCache(t, cache, replaceJSONField(t, data, "recorded_at", "not-a-time"), 0o600)
	case codecLegacy:
		writeCache(t, cache, "green deadbeef 2026-07-13T00:00:00Z\n", 0o600)
		wantReason = "invalid cache framing"
	case codecTruncated:
		writeCache(t, cache, strings.TrimSuffix(data, "}"), 0o600)
		wantReason = "invalid cache framing"
	case codec16384, codec16385:
		size := 16_384
		if kind == codec16385 {
			size++
		}
		prefix, suffix := `{"schema":1,"state":"ready","status":"green","tree":"`, `"}`
		body := prefix + strings.Repeat("a", size-len(prefix)-len(suffix)) + suffix
		if len(body) != size {
			t.Fatalf("cache bytes = %d, want %d", len(body), size)
		}
		writeCache(t, cache, body, 0o600)
	case codecSymlink:
		contract.Remove(t, cache)
		target := filepath.Join(t.TempDir(), "record")
		contract.WriteFileAbs(t, target, data)
		if err := os.Symlink(target, cache); err != nil {
			t.Fatal(err)
		}
		wantReason = "invalid cache metadata"
	case codecDirectory:
		contract.Remove(t, cache)
		contract.Mkdir(t, cache)
		wantReason = "invalid cache metadata"
	case codecUnreadable:
		writeCache(t, cache, data, 0o000)
		wantReason = "invalid cache metadata"
	}
	assertInspection(t, Inspect(f.Root), wantState, wantReason, reusable)
}

func assertInspection(t *testing.T, got Inspection, state State, reason string, reusable bool) {
	t.Helper()
	if got.State != state || got.Reason != reason || got.ReusableGreen != reusable {
		t.Fatalf("inspection = %s/%q reusable=%v, want %s/%q reusable=%v", got.State, got.Reason, got.ReusableGreen, state, reason, reusable)
	}
}

func replaceJSONField(t *testing.T, body, field, value string) string {
	t.Helper()
	var rec map[string]any
	if err := json.Unmarshal([]byte(body), &rec); err != nil {
		t.Fatal(err)
	}
	rec[field] = value
	b, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

type freshnessKind int

const (
	freshnessAfter freshnessKind = iota
	freshnessExact
	freshnessFuture
	freshnessMalformed
	freshnessFingerprint
	freshnessPolicy
)

func freshnessProof(id string, kind freshnessKind) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) { runFreshnessProof(t, kind) }}
}

func runFreshnessProof(t *testing.T, kind freshnessKind) {
	f := story3Fixture(t)
	f.WriteFile(".bench/gate-inputs.json", story3LocalManifest)
	f.Bench("gate").RequireExit(0)
	cache := filepath.Join(story3GitDir(f), "bench-last-gate")
	data := strings.TrimSuffix(contract.ReadFileAbs(t, cache), "\n")
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	data = replaceJSONField(t, data, "recorded_at", now.Add(-time.Minute).Format(time.RFC3339))
	wantState, wantReason := Ready, "verdict expired"
	switch kind {
	case freshnessAfter:
		data = replaceJSONField(t, data, "recorded_at", now.Add(-10*time.Minute-time.Second).Format(time.RFC3339))
	case freshnessExact:
		data = replaceJSONField(t, data, "recorded_at", now.Add(-10*time.Minute).Format(time.RFC3339))
	case freshnessFuture:
		data = replaceJSONField(t, data, "recorded_at", now.Add(time.Minute).Format(time.RFC3339))
		wantState, wantReason = Invalid, "invalid cache record"
	case freshnessMalformed:
		data = replaceJSONField(t, data, "recorded_at", "not-a-time")
		wantState, wantReason = Invalid, "invalid cache record"
	case freshnessFingerprint:
		data = replaceJSONField(t, data, "tree", strings.Repeat("0", 40))
		wantReason = "working tree changed"
	case freshnessPolicy:
		data = replaceJSONField(t, data, "oracle", strings.Repeat("0", 64))
		wantReason = "oracle changed"
	}
	writeCache(t, cache, data, 0o600)
	before, err := os.Stat(cache)
	if err != nil {
		t.Fatal(err)
	}
	bytesBefore := contract.ReadFileAbs(t, cache)
	got := inspectAt(f.Root, now)
	assertInspection(t, got, wantState, wantReason, false)
	after, err := os.Stat(cache)
	if err != nil {
		t.Fatal(err)
	}
	if contract.ReadFileAbs(t, cache) != bytesBefore || after.Mode() != before.Mode() || !after.ModTime().Equal(before.ModTime()) || got.CacheBytes != len(bytesBefore) {
		t.Fatalf("reader changed bytes/metadata or lost literal byte count: inspection=%+v", got)
	}
}

type secretKind int

const (
	secretCommand secretKind = iota
	secretEnvironmentName
	secretEnvironmentValue
	secretManifestPath
	secretInputContent
	secretToolOutput
	secretGateOutput
	secretControlBytes
)

func secretProof(id string, kind secretKind) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) { runSecretProof(t, kind) }}
}

func runSecretProof(t *testing.T, kind secretKind) {
	const sentinel = "FT78_SENTINEL_7f3c9a"
	f := story3Fixture(t)
	manifest, env := story3LocalManifest, map[string]string{}
	switch kind {
	case secretCommand:
		contract.Remove(t, filepath.Join(f.Root, ".bench", "gate.sh"))
		env["BENCH_GATE"] = "value='" + sentinel + "'; exit 0"
	case secretEnvironmentName:
		manifest = `{"schema":1,"closure":"local","environment":["` + sentinel + `"],"paths":[],"tools":[]}`
		env[sentinel] = "present"
	case secretEnvironmentValue:
		manifest = `{"schema":1,"closure":"local","environment":["FT78_SECRET"],"paths":[],"tools":[]}`
		env["FT78_SECRET"] = sentinel
	case secretManifestPath:
		f.WriteFile(sentinel, "safe\n")
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["` + sentinel + `"],"tools":[]}`
	case secretInputContent:
		f.WriteFile("inputs/secret", sentinel+"\n")
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["inputs/secret"],"tools":[]}`
	case secretToolOutput:
		f.WriteExecutable("tools/secret-tool", "#!/bin/sh\necho "+sentinel+"\n")
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":["tools/secret-tool"]}`
	case secretGateOutput:
		f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho "+sentinel+"\n")
	case secretControlBytes:
		f.WriteFile("inputs/control", "secret-\033-\007\n")
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["inputs/control"],"tools":[]}`
	}
	f.WriteFile(".bench/gate-inputs.json", manifest)
	probe := f.BenchEnv(env, "gate")
	probe.RequireExit(0)
	before := story3ReadVerdict(t, f)
	cache := contract.ReadFileAbs(t, filepath.Join(story3GitDir(f), "bench-last-gate"))
	combined := probe.Stderr + cache
	if strings.Contains(combined, sentinel) || strings.ContainsAny(combined, "\x1b\a") {
		t.Fatalf("secret/control sentinel leaked in cache or diagnostics: %q", combined)
	}
	if kind == secretGateOutput && probe.Stdout != sentinel+"\n" {
		t.Fatalf("gate stdout = %q, want literal sentinel output", probe.Stdout)
	}
	mutateSecretProof(t, f, kind, manifest, env, sentinel)
	f.BenchEnv(env, "gate").RequireExit(0)
	if after := story3ReadVerdict(t, f); before.Oracle == after.Oracle {
		t.Fatal("sentinel mutation did not change oracle identity")
	}
	if got := Inspect(f.Root); got.State != Ready || got.Status != "green" {
		t.Fatalf("secret-safe inspection = %+v", got)
	}
}

func mutateSecretProof(t *testing.T, f contract.Fixture, kind secretKind, manifest string, env map[string]string, sentinel string) {
	const replacement = "FT78_SAFE_REPLACEMENT"
	switch kind {
	case secretCommand:
		env["BENCH_GATE"] = strings.ReplaceAll(env["BENCH_GATE"], sentinel, replacement)
	case secretEnvironmentName:
		delete(env, sentinel)
		env[replacement] = "present"
		f.WriteFile(".bench/gate-inputs.json", strings.ReplaceAll(manifest, sentinel, replacement))
	case secretEnvironmentValue:
		env["FT78_SECRET"] = replacement
	case secretManifestPath:
		f.WriteFile(replacement, "safe\n")
		f.WriteFile(".bench/gate-inputs.json", strings.ReplaceAll(manifest, sentinel, replacement))
	case secretInputContent:
		f.WriteFile("inputs/secret", replacement+"\n")
	case secretToolOutput:
		f.WriteExecutable("tools/secret-tool", "#!/bin/sh\necho "+replacement+"\n")
	case secretGateOutput:
		f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho "+replacement+"\n")
	case secretControlBytes:
		f.WriteFile("inputs/control", replacement+"\n")
	}
}

type hostileKind int

const (
	hostileRepositorySpaces hostileKind = iota
	hostileRepositoryGlob
	hostileDeclaredSpaces
	hostileDeclaredGlob
	hostileManifestNoNewline
	hostileSymlinkChain
	hostileExternalTarget
	hostileMissingGlobalBench
	hostileExecutableMode
	hostileControlOutput
)

func hostileProof(id string, kind hostileKind) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) { runHostileProof(t, kind) }}
}

func runHostileProof(t *testing.T, kind hostileKind) {
	f := story3Fixture(t)
	manifest, env := story3LocalManifest, map[string]string{}
	switch kind {
	case hostileRepositorySpaces, hostileRepositoryGlob:
		name := "repo with spaces"
		if kind == hostileRepositoryGlob {
			name = "repo[*]?"
		}
		root := filepath.Join(t.TempDir(), name)
		if err := os.Rename(f.Root, root); err != nil {
			t.Fatal(err)
		}
		f.Root = root
	case hostileDeclaredSpaces, hostileDeclaredGlob:
		name := "inputs/path with spaces"
		if kind == hostileDeclaredGlob {
			name = "inputs/path[*]?"
		}
		f.WriteFile(name, "literal\n")
		manifest = fmt.Sprintf(`{"schema":1,"closure":"local","environment":[],"paths":[%q],"tools":[]}`, name)
	case hostileManifestNoNewline:
	case hostileSymlinkChain:
		f.WriteFile("inputs/target", "literal\n")
		if err := os.Symlink("target", filepath.Join(f.Root, "inputs", "link-b")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("link-b", filepath.Join(f.Root, "inputs", "link-a")); err != nil {
			t.Fatal(err)
		}
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["inputs/link-a"],"tools":[]}`
	case hostileExternalTarget:
		external := filepath.Join(t.TempDir(), "external")
		contract.WriteFileAbs(t, external, "literal\n")
		if err := os.Symlink(external, filepath.Join(f.Root, "ft78-external")); err != nil {
			t.Fatal(err)
		}
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["ft78-external"],"tools":[]}`
	case hostileMissingGlobalBench:
		env["PATH"] = "/usr/bin:/bin"
	case hostileExecutableMode:
		f.WriteExecutable("tools/mode-tool", "#!/bin/sh\nexit 0\n")
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":["tools/mode-tool"]}`
	case hostileControlOutput:
		f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\nprintf 'unsafe-\\033-\\007\\n'\n")
	}
	f.WriteFile(".bench/gate-inputs.json", manifest)
	first := f.BenchEnv(env, "gate")
	first.RequireExit(0)
	if kind == hostileExternalTarget {
		assertInspection(t, Inspect(f.Root), Ready, "declared path unavailable", false)
		return
	}
	if kind == hostileExecutableMode {
		before := story3ReadVerdict(t, f)
		if err := os.Chmod(filepath.Join(f.Root, "tools", "mode-tool"), 0o700); err != nil {
			t.Fatal(err)
		}
		f.BenchEnv(env, "gate").RequireExit(0)
		if after := story3ReadVerdict(t, f); before.Oracle == after.Oracle {
			t.Fatal("executable-mode mutation did not change oracle identity")
		}
	}
	if strings.ContainsAny(first.Stdout+first.Stderr+contract.ReadFileAbs(t, filepath.Join(story3GitDir(f), "bench-last-gate")), "\x1b\a") {
		t.Fatal("hostile control byte reached public output or cache")
	}
	if got := Inspect(f.Root); got.State != Ready || got.Status != "green" {
		t.Fatalf("hostile invocation inspection = %+v", got)
	}
}
