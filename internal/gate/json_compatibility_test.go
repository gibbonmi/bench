package gate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/jsonfile/jsonfiletest"
)

// gateJSONBaseline holds the stored results of the strict gate readers for the hostile
// document family. An ordinary run reads it and never writes it.
const gateJSONBaseline = "testdata/strict-json-baseline.json"

const (
	validManifestJSON  = `{"lane":[],"phases":[{"name":"build","argv":["true"]}]}`
	unknownManifestKey = "need"

	nestedDuplicateManifest = `{"phases":[{"name":"build","name":"test","argv":["true"]}]}`
	nestedUnknownManifest   = `{"phases":[{"name":"build","argv":["true"],"` + unknownManifestKey + `":["x"]}]}`
)

type gateJSONResults struct {
	Manifest map[string]manifestResult `json:"manifest"`
	Cache    map[string]cacheResult    `json:"cache"`
}

type manifestResult struct {
	Phases []string `json:"phases,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type cacheResult struct {
	State  State  `json:"state"`
	Reason string `json:"reason,omitempty"`
	Bytes  int    `json:"bytes"`
}

func TestGateJSONCompatibility(t *testing.T) {
	now := time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)
	record := string(inspectJSON(readyVerdict(subject{Tree: strings.Repeat("a", 40), Oracle: strings.Repeat("b", 64)}, "green", now)))

	// The specific contracts run before the baseline family, so a sentinel or framing
	// drift fails under its own name first.
	t.Run("duplicate sentinel", func(t *testing.T) {
		if err := strictJSON([]byte(nestedDuplicateManifest+"\n"), new(manifestDoc)); !errors.Is(err, errDuplicateJSONName) {
			t.Fatalf("strictJSON error = %v, want errDuplicateJSONName", err)
		}
	})
	t.Run("trailing sentinel", func(t *testing.T) {
		if err := strictJSON([]byte(validManifestJSON+" {}"), new(manifestDoc)); !errors.Is(err, errTrailingJSON) {
			t.Fatalf("strictJSON error = %v, want errTrailingJSON", err)
		}
	})
	t.Run("unframed documents stay valid", func(t *testing.T) {
		if got := manifestFixtureResult(t, validManifestJSON); got.Error != "" {
			t.Fatalf("unframed manifest = %+v, want it accepted", got)
		}
		if got := cacheFixtureResult(t, record, now); got.State != Ready {
			t.Fatalf("unframed cache record = %+v, want it ready", got)
		}
	})
	t.Run("unknown manifest field names its key", func(t *testing.T) {
		root := t.TempDir()
		writeLaneFile(t, manifestPath(root), nestedUnknownManifest)
		_, err := phaseTable(root, t.TempDir())
		if want := defect(manifestPath(root), "unknown field", unknownManifestKey); err == nil || err.Error() != want.Error() {
			t.Fatalf("phase table error = %v, want %v", err, want)
		}
	})

	var stored gateJSONResults
	jsonfiletest.ReadBaseline(t, gateJSONBaseline, &stored)
	current := gateJSONResults{Manifest: map[string]manifestResult{}, Cache: map[string]cacheResult{}}
	manifests := jsonfiletest.Seed{Valid: validManifestJSON, Member: `"lane":[]`, Wrong: `"lane":{}`, Nested: `"name":"build"`}
	for name, body := range jsonfiletest.Family(manifests) {
		current.Manifest[name] = manifestFixtureResult(t, body)
	}
	for name, body := range jsonfiletest.Family(jsonfiletest.Seed{Valid: record, Member: `"schema":1`, Wrong: `"schema":"1"`}) {
		current.Cache[name] = cacheFixtureResult(t, body, now)
	}
	current.Cache["oversize"] = cacheFixtureResult(t, strings.Replace(record, "{", "{"+strings.Repeat(" ", cacheLimit), 1), now)
	for name, result := range cacheFileStates(t, record, now) {
		current.Cache[name] = result
	}
	storedHost := map[string]cacheResult{}
	for name := range hostCacheFiles {
		if result, ok := stored.Cache[name]; ok {
			storedHost[name] = result
			delete(stored.Cache, name)
		}
	}

	t.Run("manifest baseline", func(t *testing.T) { jsonfiletest.Compare(t, current.Manifest, stored.Manifest) })
	t.Run("cache baseline", func(t *testing.T) { jsonfiletest.Compare(t, current.Cache, stored.Cache) })
	t.Run("unreadable cache baseline", func(t *testing.T) { jsonfiletest.Compare(t, hostCacheStates(t, now), storedHost) })
}

// manifestFixtureResult loads body as the phase manifest of a fresh root. The diagnostic
// names the root as $ROOT, because only the temporary path differs between runs.
func manifestFixtureResult(t *testing.T, body string) manifestResult {
	t.Helper()
	root := t.TempDir()
	writeLaneFile(t, manifestPath(root), body)
	phases, err := phaseTable(root, t.TempDir())
	if err != nil {
		return manifestResult{Error: strings.ReplaceAll(err.Error(), root, "$ROOT")}
	}
	return manifestResult{Phases: phaseNames(phases)}
}

func cacheFixtureResult(t *testing.T, body string, now time.Time) cacheResult {
	t.Helper()
	dir := t.TempDir()
	outcomeWrite(t, dir, "cache", body, 0o600)
	return loadedCacheResult(loadVerdict(filepath.Join(dir, "cache"), now))
}

func loadedCacheResult(loaded loadedVerdict) cacheResult {
	return cacheResult{State: loaded.state, Reason: loaded.reason, Bytes: loaded.bytes}
}

// cacheFileStates loads the file states that the cache reader grades before it reads any
// JSON: an absent path, a wide mode, a directory, and a path under a regular file.
func cacheFileStates(t *testing.T, record string, now time.Time) map[string]cacheResult {
	t.Helper()
	dir := t.TempDir()
	outcomeWrite(t, dir, "wide", record, 0o644)
	if err := os.Mkdir(filepath.Join(dir, "directory"), 0o700); err != nil {
		t.Fatal(err)
	}
	outcomeWrite(t, dir, "file", record, 0o600)
	return map[string]cacheResult{
		"absent":       loadedCacheResult(loadVerdict(filepath.Join(dir, "absent"), now)),
		"wide mode":    loadedCacheResult(loadVerdict(filepath.Join(dir, "wide"), now)),
		"directory":    loadedCacheResult(loadVerdict(filepath.Join(dir, "directory"), now)),
		"under a file": loadedCacheResult(loadVerdict(filepath.Join(dir, "file", "cache"), now)),
	}
}

// hostCacheFiles names host files whose metadata the cache reader accepts and whose bytes it
// cannot read. A process other than root cannot open the memory of pid 1. A read of process
// memory at offset zero fails.
var hostCacheFiles = map[string]string{"unopenable": "/proc/1/mem", "unreadable": "/proc/self/mem"}

// hostCacheStates loads each host cache file. It skips t on a host that has no such file.
func hostCacheStates(t *testing.T, now time.Time) map[string]cacheResult {
	t.Helper()
	results := map[string]cacheResult{}
	for name, path := range hostCacheFiles {
		if _, err := os.Lstat(path); err != nil {
			capability.Environment(t, fmt.Sprintf("the host has no %s cache file: %v", name, err))
		}
		results[name] = loadedCacheResult(loadVerdict(path, now))
	}
	return results
}
