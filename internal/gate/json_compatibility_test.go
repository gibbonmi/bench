package gate

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
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

	stored := readGateJSONBaseline(t)
	current := gateJSONResults{Manifest: map[string]manifestResult{}, Cache: map[string]cacheResult{}}
	manifests := jsonFamily(validManifestJSON, `"lane":[]`, `"lane":{}`)
	manifests["nested duplicate key"] = nestedDuplicateManifest
	manifests["nested unknown field"] = nestedUnknownManifest
	for name, body := range manifests {
		current.Manifest[name] = manifestFixtureResult(t, body)
	}
	for name, body := range jsonFamily(record, `"schema":1`, `"schema":"1"`) {
		current.Cache[name] = cacheFixtureResult(t, body, now)
	}
	current.Cache["oversize"] = cacheFixtureResult(t, strings.Replace(record, "{", "{"+strings.Repeat(" ", cacheLimit), 1), now)
	for name, result := range cacheFileStates(t, record, now) {
		current.Cache[name] = result
	}

	t.Run("manifest baseline", func(t *testing.T) { compareGateJSON(t, current.Manifest, stored.Manifest) })
	t.Run("cache baseline", func(t *testing.T) { compareGateJSON(t, current.Cache, stored.Cache) })
}

// jsonFamily derives the hostile document family from one valid object document. member
// is one field of valid, and wrong gives that field a value of the wrong type.
func jsonFamily(valid, member, wrong string) map[string]string {
	unknown := strings.Replace(valid, "{", `{"extra":true,`, 1)
	duplicate := strings.Replace(valid, member, member+","+member, 1)
	mistyped := strings.Replace(valid, member, wrong, 1)
	return map[string]string{
		"framed":                    valid + "\n",
		"unframed":                  valid,
		"surrounding whitespace":    " \n" + valid + " \t\r\n",
		"whitespace only":           " \n\t",
		"empty":                     "",
		"null":                      "null",
		"unknown field":             unknown,
		"duplicate key":             duplicate,
		"syntax error":              strings.TrimSuffix(valid, "}") + ",}",
		"wrong type":                mistyped,
		"second value":              valid + " {}",
		"framed second value":       valid + "\n{}\n",
		"malformed suffix":          valid + " nope",
		"unknown plus suffix":       unknown + " {}",
		"duplicate plus suffix":     duplicate + " {}",
		"unknown plus duplicate":    strings.Replace(duplicate, "{", `{"extra":true,`, 1),
		"wrong type plus suffix":    mistyped + " {}",
		"wrong type plus duplicate": strings.Replace(mistyped, wrong, wrong+","+member, 1),
	}
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

func readGateJSONBaseline(t *testing.T) gateJSONResults {
	t.Helper()
	data, err := os.ReadFile(gateJSONBaseline)
	if err != nil {
		t.Fatal(err)
	}
	var stored gateJSONResults
	if err := json.Unmarshal(data, &stored); err != nil {
		t.Fatal(err)
	}
	return stored
}

func compareGateJSON[T any](t *testing.T, current, stored map[string]T) {
	t.Helper()
	names := make([]string, 0, len(current)+len(stored))
	for name := range current {
		names = append(names, name)
	}
	for name := range stored {
		if _, ok := current[name]; !ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		got, ok := current[name]
		want, recorded := stored[name]
		switch {
		case !ok:
			t.Errorf("%q: the baseline case has no current result", name)
		case !recorded || !reflect.DeepEqual(got, want):
			t.Errorf("%q: current %s, baseline %s", name, inspectJSON(got), inspectJSON(want))
		}
	}
}
