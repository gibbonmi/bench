package repairpilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

type storedHostileFixture struct {
	name string
	make func(*testing.T, string, Options)
}

func commonHostileFixtures(valid func(*testing.T, Options) []byte) []storedHostileFixture {
	return []storedHostileFixture{
		{name: "empty", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, nil) }},
		{name: "malformed", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, []byte("{\n")) }},
		{name: "oversized", make: func(t *testing.T, path string, _ Options) {
			writeFixture(t, path, []byte(strings.Repeat("x", int(bounds.ControlRecordLimit)+1)))
		}},
		{name: "directory", make: func(t *testing.T, path string, _ Options) {
			if err := os.Mkdir(path, 0o700); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "live-symlink", make: func(t *testing.T, path string, options Options) {
			target := filepath.Join(t.TempDir(), "target.json")
			writeFixture(t, target, valid(t, options))
			if err := os.Symlink(target, path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "dangling-symlink", make: func(t *testing.T, path string, _ Options) {
			if err := os.Symlink("missing", path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "linked-parent", make: func(t *testing.T, path string, options Options) {
			parent := filepath.Dir(path)
			if err := os.Remove(parent); err != nil {
				t.Fatal(err)
			}
			target := t.TempDir()
			writeFixture(t, filepath.Join(target, filepath.Base(path)), valid(t, options))
			if err := os.Symlink(target, parent); err != nil {
				t.Fatal(err)
			}
		}},
	}
}

func storedHostileFixtures() []storedHostileFixture {
	fixtures := commonHostileFixtures(validDocumentBytes)
	fixtures = append(fixtures,
		storedHostileFixture{name: "no-final-newline", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf(`{"version":1,"repository_key":%q,"activated_at":"2026-09-01T12:00:00Z","observations":[],"audits":[]}`, options.RepoKey)))
		}},
		storedHostileFixture{name: "unsupported-version", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf("{\"version\":2,\"repository_key\":%q,\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[]}\n", options.RepoKey)))
		}},
		storedHostileFixture{name: "foreign-repository", make: func(t *testing.T, path string, _ Options) {
			writeFixture(t, path, []byte("{\"version\":1,\"repository_key\":\"other-123\",\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[]}\n"))
		}},
		storedHostileFixture{name: "duplicate-key", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf("{\"version\":1,\"version\":1,\"repository_key\":%q,\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[]}\n", options.RepoKey)))
		}},
		storedHostileFixture{name: "unknown-field", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf("{\"version\":1,\"repository_key\":%q,\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[],\"extra\":true}\n", options.RepoKey)))
		}},
	)
	return append(fixtures, platformStoredHostileFixtures()...)
}

func recordInputHostileFixtures(t *testing.T) []storedHostileFixture {
	t.Helper()
	valid := recordInputBytes(t, failureInput("hostile-valid", sequenceKey("source-a", "spec-a", "chunk-a")))
	withRoot := func(field string) []byte { return bytes.Replace(valid, []byte("{"), []byte("{\n  "+field+","), 1) }
	fixtures := commonHostileFixtures(func(*testing.T, Options) []byte { return valid })
	fixtures = append(fixtures,
		storedHostileFixture{name: "absent", make: func(*testing.T, string, Options) {}},
		storedHostileFixture{name: "duplicate-key", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, withRoot(`"version": 1`)) }},
		storedHostileFixture{name: "unknown-field", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, withRoot(`"extra": true`)) }},
		storedHostileFixture{name: "foreign-repository", make: func(t *testing.T, path string, _ Options) {
			writeFixture(t, path, withRoot(`"repository_key": "other-123"`))
		}},
		storedHostileFixture{name: "unsupported-version", make: func(t *testing.T, path string, _ Options) {
			writeFixture(t, path, bytes.Replace(valid, []byte(`"version": 1`), []byte(`"version": 2`), 1))
		}},
		storedHostileFixture{name: "both-payloads", make: func(t *testing.T, path string, _ Options) {
			input := failureInput("both-payloads", sequenceKey("source-a", "spec-a", "chunk-a"))
			input.Audit = &audit{}
			writeFixture(t, path, recordInputBytes(t, input))
		}},
	)
	return append(fixtures, platformStoredHostileFixtures()...)
}

func recordInputBytes(t *testing.T, input recordInput) []byte {
	t.Helper()
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

func (h *recordHarness) report(t *testing.T) string {
	t.Helper()
	out, code := Command(h.options, []string{"report"})
	if code != 0 {
		t.Fatalf("report = output %q, exit %d", out, code)
	}
	return out
}

func assertSummaryValue(t *testing.T, out, field, value string) {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	parts := strings.SplitN(lines[0], "{", 2)
	if len(lines) != 2 || len(parts) != 2 {
		t.Fatalf("summary = %q, want one table row", out)
	}
	columns := strings.Split(strings.TrimSuffix(strings.TrimSuffix(parts[1], ":"), "}"), ",")
	values := strings.Split(strings.TrimSpace(lines[1]), ",")
	for i, column := range columns {
		if column == field && i < len(values) && strings.Trim(values[i], "\"") == value {
			return
		}
	}
	t.Fatalf("summary field %s = output %q, want %q", field, out, value)
}
