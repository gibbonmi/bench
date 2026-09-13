package repairpilot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestRepairPilotActivation(t *testing.T) {
	t.Run("inactive", func(t *testing.T) {
		home := t.TempDir()
		out, code := Command(testOptions(home), []string{"report"})
		if code != 0 || !strings.Contains(out, "inactive") {
			t.Fatalf("report = output %q, exit %d; want inactive at exit 0", out, code)
		}
		if _, err := os.Stat(filepath.Join(home, "repair-pilot")); !os.IsNotExist(err) {
			t.Fatalf("inactive report created state: %v", err)
		}
	})
	t.Run("activate", func(t *testing.T) {
		home := t.TempDir()
		options := testOptions(home)
		out, code := Command(options, []string{"activate"})
		if code != 0 || !strings.Contains(out, options.Now.Format(time.RFC3339)) {
			t.Fatalf("activate = output %q, exit %d; want supplied time at exit 0", out, code)
		}
		data, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), `"activated_at": "2026-09-01T12:00:00Z"`) {
			t.Fatalf("pilot document = %s, want stored activation", data)
		}
	})
	t.Run("linked", func(t *testing.T) {
		home := t.TempDir()
		options := testOptions(home)
		options.KitSource = false
		out, code := Command(options, []string{"activate"})
		if code != 1 || !strings.Contains(out, "limited to the Bench kit") {
			t.Fatalf("linked activation = output %q, exit %d; want kit-only refusal", out, code)
		}
		if _, err := os.Stat(filepath.Join(home, "repair-pilot")); !os.IsNotExist(err) {
			t.Fatalf("linked activation created state: %v", err)
		}
	})
	t.Run("repeat", func(t *testing.T) {
		home := t.TempDir()
		options := testOptions(home)
		if _, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("initial activation exit = %d", code)
		}
		document := Document{
			Version:       1,
			RepositoryKey: options.RepoKey,
			ActivatedAt:   options.Now,
			Observations:  []observation{{ID: "kept"}},
			Audits:        []audit{},
		}
		writeTestDocument(t, documentPath(options), document)
		before, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		options.Now = options.Now.Add(24 * time.Hour)
		if _, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("repeat activation exit = %d", code)
		}
		after, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		if string(after) != string(before) {
			t.Fatalf("repeat activation changed document\nbefore: %s\nafter: %s", before, after)
		}
	})
	t.Run("stopped", func(t *testing.T) {
		home := t.TempDir()
		options := testOptions(home)
		if _, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("initial activation exit = %d", code)
		}
		before, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		options.Now = options.Now.AddDate(0, 0, 14)
		out, code := Command(options, []string{"activate"})
		if code != 0 || !strings.Contains(out, "stopped") {
			t.Fatalf("stopped activation = output %q, exit %d", out, code)
		}
		after, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		if string(after) != string(before) {
			t.Fatalf("stopped activation changed document\nbefore: %s\nafter: %s", before, after)
		}
	})
}

func TestRepairPilotStorage(t *testing.T) {
	t.Run("worktrees", func(t *testing.T) {
		home := t.TempDir()
		first := testOptions(home)
		first.Root = "/repo/worktree-one"
		second := first
		second.Root = "/repo/worktree-two"
		if _, code := Command(first, []string{"activate"}); code != 0 {
			t.Fatalf("first worktree activation exit = %d", code)
		}
		out, code := Command(second, []string{"report"})
		if code != 0 || !strings.Contains(out, "active") {
			t.Fatalf("second worktree report = output %q, exit %d", out, code)
		}
		if documentPath(first) != documentPath(second) {
			t.Fatalf("worktrees resolved different documents: %q and %q", documentPath(first), documentPath(second))
		}
	})
	t.Run("custody", func(t *testing.T) {
		home := t.TempDir()
		options := testOptions(home)
		if _, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("activation exit = %d", code)
		}
		path := documentPath(options)
		if strings.HasPrefix(path, filepath.Join(home, "worktrees")+string(os.PathSeparator)) {
			t.Fatalf("pilot document is inside disposable pool: %s", path)
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		dirInfo, err := os.Stat(filepath.Dir(path))
		if err != nil {
			t.Fatal(err)
		}
		if fileInfo.Mode().Perm() != 0o600 || dirInfo.Mode().Perm() != 0o700 {
			t.Fatalf("custody modes = file %04o dir %04o; want 0600 and 0700", fileInfo.Mode().Perm(), dirInfo.Mode().Perm())
		}
	})
	t.Run("locked", func(t *testing.T) {
		home := t.TempDir()
		options := testOptions(home)
		if _, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("activation exit = %d", code)
		}
		path := documentPath(options)
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(filepath.Dir(path), "pilot.lock"), []byte("owner\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		out, code := Command(options, []string{"activate"})
		if code != 1 || !strings.Contains(out, "active writer") {
			t.Fatalf("locked activation = output %q, exit %d; want retry refusal", out, code)
		}
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(after) != string(before) {
			t.Fatalf("locked activation changed document\nbefore: %s\nafter: %s", before, after)
		}
	})
	t.Run("create-failure", func(t *testing.T) {
		options, replacement, before := replacementFixture(t)
		options.Files.CreateTemp = func(string, string) (AtomicFile, error) {
			return nil, errors.New("fixture create refusal")
		}
		if err := (Store{Path: documentPath(options), Files: options.Files}).Replace(replacement); err == nil || !strings.Contains(err.Error(), "fixture create refusal") {
			t.Fatalf("create failure = %v; want fixture refusal", err)
		}
		assertDocumentBytes(t, documentPath(options), before)
	})
	t.Run("write-failure", func(t *testing.T) {
		options, replacement, before := replacementFixture(t)
		options.Files.CreateTemp = func(directory, pattern string) (AtomicFile, error) {
			file, err := os.CreateTemp(directory, pattern)
			if err != nil {
				return nil, err
			}
			return shortWriteFile{File: file}, nil
		}
		if err := (Store{Path: documentPath(options), Files: options.Files}).Replace(replacement); !errors.Is(err, io.ErrShortWrite) {
			t.Fatalf("partial write = %v; want short write", err)
		}
		assertDocumentBytes(t, documentPath(options), before)
		entries, err := os.ReadDir(filepath.Dir(documentPath(options)))
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 || entries[0].Name() != "pilot.json" {
			t.Fatalf("partial write left state: %v", entries)
		}
	})
	t.Run("replace-failure", func(t *testing.T) {
		options, replacement, before := replacementFixture(t)
		options.Files.Rename = func(string, string) error {
			return errors.New("fixture replace refusal")
		}
		if err := (Store{Path: documentPath(options), Files: options.Files}).Replace(replacement); err == nil || !strings.Contains(err.Error(), "fixture replace refusal") {
			t.Fatalf("replace failure = %v; want fixture refusal", err)
		}
		assertDocumentBytes(t, documentPath(options), before)
		entries, err := os.ReadDir(filepath.Dir(documentPath(options)))
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 || entries[0].Name() != "pilot.json" {
			t.Fatalf("replace failure left state: %v", entries)
		}
	})
	t.Run("stored-hostile", func(t *testing.T) {
		for _, operation := range []string{"activate", "report"} {
			for _, fixture := range storedHostileFixtures() {
				t.Run(operation+"/"+fixture.name, func(t *testing.T) {
					assertStoredHostileRefusal(t, operation, fixture)
				})
			}
		}
	})
	t.Run("record-stored-hostile", func(t *testing.T) {
		requireFixtureCase(t)
		for _, fixture := range storedHostileFixtures() {
			t.Run(fixture.name, func(t *testing.T) {
				assertStoredHostileRefusal(t, "record", fixture)
			})
		}
	})
}

func storedHostileFixtures() []storedHostileFixture {
	fixtures := []storedHostileFixture{
		{name: "empty", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, nil) }},
		{name: "malformed", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, []byte("{\n")) }},
		{name: "no-final-newline", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf(`{"version":1,"repository_key":%q,"activated_at":"2026-09-01T12:00:00Z","observations":[],"audits":[]}`, options.RepoKey)))
		}},
		{name: "unsupported-version", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf("{\"version\":2,\"repository_key\":%q,\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[]}\n", options.RepoKey)))
		}},
		{name: "foreign-repository", make: func(t *testing.T, path string, _ Options) {
			writeFixture(t, path, []byte("{\"version\":1,\"repository_key\":\"other-123\",\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[]}\n"))
		}},
		{name: "duplicate-key", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf("{\"version\":1,\"version\":1,\"repository_key\":%q,\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[]}\n", options.RepoKey)))
		}},
		{name: "unknown-field", make: func(t *testing.T, path string, options Options) {
			writeFixture(t, path, []byte(fmt.Sprintf("{\"version\":1,\"repository_key\":%q,\"activated_at\":\"2026-09-01T12:00:00Z\",\"observations\":[],\"audits\":[],\"extra\":true}\n", options.RepoKey)))
		}},
		{name: "oversized", make: func(t *testing.T, path string, _ Options) {
			writeFixture(t, path, []byte(strings.Repeat("x", int(bounds.ControlRecordLimit)+1)))
		}},
		{name: "directory", make: func(t *testing.T, path string, _ Options) {
			if err := os.Mkdir(path, 0o700); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "live-symlink", make: func(t *testing.T, path string, _ Options) {
			target := filepath.Join(t.TempDir(), "target")
			writeFixture(t, target, []byte("outside\n"))
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
			writeFixture(t, filepath.Join(target, "pilot.json"), validDocumentBytes(t, options))
			if err := os.Symlink(target, parent); err != nil {
				t.Fatal(err)
			}
		}},
	}
	return append(fixtures, platformStoredHostileFixtures()...)
}

func assertStoredHostileRefusal(t *testing.T, operation string, fixture storedHostileFixture) {
	t.Helper()
	options := testOptions(t.TempDir())
	path := documentPath(options)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	fixture.make(t, path, options)
	infoBefore, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	argv := []string{operation}
	if operation == "record" {
		input := failureInput("stored-hostile", sequenceKey("source-a", "spec-a", "chunk-a"))
		data, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		inputPath := filepath.Join(t.TempDir(), "input.json")
		writeFixture(t, inputPath, append(data, '\n'))
		argv = []string{"record", "--input", inputPath}
	}
	out, code := Command(options, argv)
	if code != 1 || !strings.Contains(out, "repair pilot refused") {
		t.Fatalf("%s hostile %s = output %q, exit %d; want refusal", operation, fixture.name, out, code)
	}
	infoAfter, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if infoAfter.Mode() != infoBefore.Mode() || infoAfter.Size() != infoBefore.Size() {
		t.Fatalf("%s hostile %s changed stored object", operation, fixture.name)
	}
}

func writeTestDocument(t *testing.T, path string, document Document) {
	t.Helper()
	writeFixture(t, path, documentBytes(t, document))
}

func writeFixture(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func validDocumentBytes(t *testing.T, options Options) []byte {
	t.Helper()
	document := Document{Version: 1, RepositoryKey: options.RepoKey, ActivatedAt: options.Now, Observations: []observation{}, Audits: []audit{}}
	return documentBytes(t, document)
}

func documentBytes(t *testing.T, document Document) []byte {
	t.Helper()
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

func replacementFixture(t *testing.T) (Options, Document, []byte) {
	t.Helper()
	options := testOptions(t.TempDir())
	path := documentPath(options)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	original := Document{Version: 1, RepositoryKey: options.RepoKey, ActivatedAt: options.Now, Observations: []observation{{ID: "prior"}}, Audits: []audit{}}
	writeTestDocument(t, path, original)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	replacement := original
	replacement.Observations = []observation{{ID: "replacement"}}
	return options, replacement, before
}

func assertDocumentBytes(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("pilot document changed\nwant: %s\ngot: %s", want, got)
	}
}

func testOptions(home string) Options {
	return Options{
		Home:      home,
		Root:      "/repo",
		RepoKey:   "bench-123",
		KitSource: true,
		Now:       time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC),
	}
}
