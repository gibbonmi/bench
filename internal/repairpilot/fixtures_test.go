package repairpilot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/bounds"
)

type storedHostileFixture struct {
	name string
	make func(*testing.T, string, Options)
}

func TestRepairPilotGrammar(t *testing.T) {
	t.Run("usage", func(t *testing.T) {
		options := testOptions(t.TempDir())
		if err := os.MkdirAll(filepath.Dir(documentPath(options)), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("missing-target", documentPath(options)); err != nil {
			t.Fatal(err)
		}
		out, code := Command(options, []string{"bogus"})
		if code != 2 || !strings.Contains(out, "unknown argument: bogus") {
			t.Fatalf("bogus operation = output %q, exit %d; want usage at exit 2", out, code)
		}
	})
	t.Run("hostile", func(t *testing.T) {
		requireFixtureCase(t)
		newRecordHarness(t).accept(t, failureInput("hostile-control", sequenceKey("source-a", "spec-a", "chunk-a")))
		for _, fixture := range recordInputHostileFixtures(t) {
			t.Run(fixture.name, func(t *testing.T) {
				options := testOptions(t.TempDir())
				if out, code := Command(options, []string{"activate"}); code != 0 {
					t.Fatalf("activate = output %q, exit %d", out, code)
				}
				options.Now = parseTime("2026-09-10T12:00:00Z")
				before, err := os.ReadFile(documentPath(options))
				if err != nil {
					t.Fatal(err)
				}
				inputPath := filepath.Join(t.TempDir(), "input.json")
				fixture.make(t, inputPath, options)
				out, code := Command(options, []string{"record", "--input", inputPath})
				if code != 1 || !strings.Contains(out, "refused") {
					t.Fatalf("hostile input %s = output %q, exit %d", fixture.name, out, code)
				}
				assertDocumentBytes(t, documentPath(options), before)
			})
		}
	})
	t.Run("path-shape", func(t *testing.T) {
		requireFixtureCase(t)
		options := testOptions(t.TempDir())
		if out, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("activate = output %q, exit %d", out, code)
		}
		options.Now = parseTime("2026-09-10T12:00:00Z")
		path := filepath.Join(t.TempDir(), "input [*].json")
		writeFixture(t, path, recordInputBytes(t, failureInput("path-shape", sequenceKey("source-a", "spec-a", "chunk-a"))))
		if out, code := Command(options, []string{"record", "--input", path}); code != 0 {
			t.Fatalf("record path shape = output %q, exit %d", out, code)
		}
	})
	t.Run("no-final-newline", func(t *testing.T) {
		options := testOptions(t.TempDir())
		if out, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("activate = output %q, exit %d", out, code)
		}
		options.Now = parseTime("2026-09-10T12:00:00Z")
		data := recordInputBytes(t, failureInput("no-newline", sequenceKey("source-a", "spec-a", "chunk-a")))
		path := filepath.Join(t.TempDir(), "input.json")
		writeFixture(t, path, data[:len(data)-1])
		if out, code := Command(options, []string{"record", "--input", path}); code != 0 {
			t.Fatalf("record without final newline = output %q, exit %d", out, code)
		}
	})
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
	rows := reportRows(t, out, "repair_pilot")
	if len(rows) != 1 || fmt.Sprint(rows[0][field]) != value {
		t.Fatalf("summary field %s = %#v, want %q", field, rows, value)
	}
}

func reportRows(t *testing.T, out, block string) []map[string]any {
	t.Helper()
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatalf("decode report: %v\n%s", err, out)
	}
	decoded, err := document.Rows(block)
	if err != nil {
		t.Fatal(err)
	}
	rows := make([]map[string]any, len(decoded))
	for index, value := range decoded {
		row, ok := value.(map[string]any)
		if !ok {
			t.Fatalf("%s row %d = %T, want an object", block, index, value)
		}
		rows[index] = row
	}
	return rows
}

func assertReportClass(t *testing.T, out, class, status string) {
	t.Helper()
	assertRowsContain(t, out, "required_classes", map[string]string{"class": class, "status": status})
}

func assertRowsContain(t *testing.T, out, block string, fields map[string]string) {
	t.Helper()
	rows := reportRows(t, out, block)
	for _, row := range rows {
		matches := true
		for field, want := range fields {
			if fmt.Sprint(row[field]) != want {
				matches = false
				break
			}
		}
		if matches {
			return
		}
	}
	t.Fatalf("%s has no row matching %#v: %#v", block, fields, rows)
}

func rowField(rows []map[string]any, field string) []string {
	values := make([]string, len(rows))
	for index, row := range rows {
		values[index] = fmt.Sprint(row[field])
	}
	return values
}

func rowFields(rows []map[string]any, fields ...string) [][]string {
	values := make([][]string, len(rows))
	for index, row := range rows {
		values[index] = make([]string, len(fields))
		for fieldIndex, field := range fields {
			values[index][fieldIndex] = fmt.Sprint(row[field])
		}
	}
	return values
}

func assertReportValue(t *testing.T, out, column, value string) {
	t.Helper()
	assertSummaryValue(t, out, column, value)
}

func reverseObservations(items []observation) []observation {
	result := append([]observation{}, items...)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func reverseAudits(items []audit) []audit {
	result := append([]audit{}, items...)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}
