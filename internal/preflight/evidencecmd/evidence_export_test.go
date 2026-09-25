package evidencecmd_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// The file names, the index columns, and the response line below are the public export
// contract, so this file states them independently of the export owner.

// publishExportFixture publishes one three-source artifact in a fresh repository: the
// metadata source and two repository sources.
func publishExportFixture(t *testing.T) (root string, pack *chargeevidence.Pack) {
	t.Helper()
	root = preflighttest.StartRepo(t)
	var inputs []chargeevidence.SourceInput
	for i, body := range []string{"first source body\n", "second source 雪 \"quoted\" body\n"} {
		inputs = append(inputs, chargeevidence.SourceInput{Role: "input", Kind: chargeevidence.KindRepository,
			Path: fmt.Sprintf("input-%d.md", i+1), Required: true, Data: []byte(body)})
	}
	pack = publishReview(t, root, inputs...)
	if got := len(pack.Manifest().Sources); got != 3 {
		t.Fatalf("export fixture declares %d sources, want 3", got)
	}
	return root, pack
}

// exportEntries lists the names in dir, or nil when dir does not exist.
func exportEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// exportIndex decodes the one index table of an export directory.
func exportIndex(t *testing.T, dir string) []map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "index.toon"))
	if err != nil {
		t.Fatal(err)
	}
	document := preflighttest.DecodeMap(t, string(data))
	if len(document) != 1 {
		t.Fatalf("index holds %d blocks, want one:\n%s", len(document), data)
	}
	return preflighttest.TableRows(t, document, "sources")
}

// exportCommand runs the public export of identity into dir.
func exportCommand(identity, dir string) (string, int) {
	return preflight.Command([]string{"evidence", identity, "--to", dir})
}

// assertAbsent fails when path exists.
func assertAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s exists after a refused export: %v", path, err)
	}
}

// TestEvidenceExportWritesSources is BO45: a relative directory resolves against the working
// directory, and each source reaches its own ordinal file with its exact bytes.
func TestEvidenceExportWritesSources(t *testing.T) {
	_, pack := publishExportFixture(t)
	out, code := exportCommand(pack.Identity(), "export")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(wd, "export")
	sources := pack.Manifest().Sources
	total := 0
	for i, source := range sources {
		total += source.Bytes
		want, _ := pack.Source(source.ID)
		got, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("source-%d", i+1)))
		if err != nil || string(got) != string(want) {
			t.Errorf("source-%d = %q (%v), want the %s bytes %q", i+1, got, err, source.ID, want)
		}
	}
	if want := fmt.Sprintf("exported{sources=%d,bytes=%d,dir=%s}\n", len(sources), total, dir); code != 0 || out != want {
		t.Fatalf("export = (%d):\n%s\nwant (0):\n%s", code, out, want)
	}
	if got := exportEntries(t, dir); !slices.Equal(got, []string{"index.toon", "source-1", "source-2", "source-3"}) {
		t.Fatalf("export directory holds %v", got)
	}
	rows := exportIndex(t, dir)
	if len(rows) != len(sources) {
		t.Fatalf("index holds %d rows, want %d", len(rows), len(sources))
	}
	for i, row := range rows {
		want := map[string]any{"ordinal": float64(i + 1), "id": sources[i].ID, "bytes": float64(sources[i].Bytes), "file": fmt.Sprintf("source-%d", i+1)}
		for field, value := range want {
			if row[field] != value {
				t.Errorf("index row %d %s = %v, want %v", i, field, row[field], value)
			}
		}
		if len(row) != len(want) {
			t.Errorf("index row %d holds columns %v", i, row)
		}
	}
}

// TestEvidenceExportRefusesNonEmptyDir is BO46: an existing file keeps the export out.
func TestEvidenceExportRefusesNonEmptyDir(t *testing.T) {
	_, pack := publishExportFixture(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "keep"), []byte("kept\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := exportCommand(pack.Identity(), dir)
	if code != 1 || strings.Contains(out, "exported{") {
		t.Fatalf("export into a non-empty directory = (%d):\n%s", code, out)
	}
	kept, err := os.ReadFile(filepath.Join(dir, "keep"))
	if got := exportEntries(t, dir); !slices.Equal(got, []string{"keep"}) || err != nil || string(kept) != "kept\n" {
		t.Fatalf("refused export left %v and keep = %q (%v)", got, kept, err)
	}
}

// TestEvidenceExportRefusesSymlink is BO47: a symlink at the directory path refuses by its
// kind, whether it names an empty directory or nothing.
func TestEvidenceExportRefusesSymlink(t *testing.T) {
	_, pack := publishExportFixture(t)
	target := t.TempDir()
	for name, destination := range map[string]string{"empty directory": target, "dangling": filepath.Join(target, "absent")} {
		t.Run(name, func(t *testing.T) {
			link := filepath.Join(t.TempDir(), "link")
			if err := os.Symlink(destination, link); err != nil {
				t.Fatal(err)
			}
			out, code := exportCommand(pack.Identity(), link)
			if code != 1 || !strings.Contains(out, "is a symlink") || strings.Contains(out, "exported{") {
				t.Fatalf("export through a symlink = (%d):\n%s", code, out)
			}
			if got := exportEntries(t, target); len(got) != 0 {
				t.Fatalf("export through a symlink wrote %v", got)
			}
			if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
				t.Fatalf("symlink after the refusal = %v, %v", info, err)
			}
		})
	}
}

// corruptLastSource changes the last stored byte of identity's pack. The source bodies close
// the pack, so the byte belongs to the last source's last page.
func corruptLastSource(t *testing.T, root, identity string) {
	t.Helper()
	path := filepath.Join(preflighttest.StoreDir(t, root), strings.TrimPrefix(identity, "sha256:")+".pack")
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteAt([]byte("X"), info.Size()-1); err != nil {
		t.Fatal(err)
	}
}

// TestEvidenceExportVerifiesBeforeWrite is BO48: a corrupt page of the last source refuses the
// export, the writer never receives that source, and the created directory is gone.
func TestEvidenceExportVerifiesBeforeWrite(t *testing.T) {
	root, pack := publishExportFixture(t)
	corruptLastSource(t, root, pack.Identity())
	t.Run("command", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "export")
		out, code := exportCommand(pack.Identity(), dir)
		if code != 1 || !strings.Contains(out, chargeevidence.RefusePageDigest) || strings.Contains(out, "exported{") {
			t.Fatalf("export of a corrupt pack = (%d):\n%s", code, out)
		}
		assertAbsent(t, dir)
	})
	t.Run("writer", func(t *testing.T) {
		var written []string
		record := func(w io.Writer, data []byte) error {
			written = append(written, string(data))
			return chargeevidence.WriteAll(w, data)
		}
		dir := filepath.Join(t.TempDir(), "export")
		if out, code := evidencecmd.Export(root, pack.Identity(), dir, record); code != 1 {
			t.Fatalf("export of a corrupt pack = (%d):\n%s", code, out)
		}
		sources := pack.Manifest().Sources
		verified := sources[:len(sources)-1]
		if len(written) != len(verified) {
			t.Fatalf("the writer received %d files, want only the %d sources before the corrupt one", len(written), len(verified))
		}
		for i, source := range verified {
			if want, _ := pack.Source(source.ID); written[i] != string(want) {
				t.Errorf("write %d = %q, want the verified %s bytes", i, written[i], source.ID)
			}
		}
		assertAbsent(t, dir)
	})
}

// TestEvidenceExportNamesByOrdinal is BO49: a hostile source identifier names only an index
// cell, and its bytes land in the ordinal file inside the directory.
func TestEvidenceExportNamesByOrdinal(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "export")
	body := "hostile identifier body\n"
	exported, err := chargeevidence.ExportSources(dir, 1, func(int) (string, []byte, error) {
		return "../escape", []byte(body), nil
	}, chargeevidence.WriteAll)
	if err != nil || exported != (chargeevidence.Exported{Dir: dir, Sources: 1, Bytes: len(body)}) {
		t.Fatalf("export = %+v, %v", exported, err)
	}
	if got := exportEntries(t, parent); !slices.Equal(got, []string{"export"}) {
		t.Fatalf("export parent holds %v, want only the export directory", got)
	}
	if got := exportEntries(t, dir); !slices.Equal(got, []string{"index.toon", "source-1"}) {
		t.Fatalf("export directory holds %v", got)
	}
	if got, err := os.ReadFile(filepath.Join(dir, "source-1")); err != nil || string(got) != body {
		t.Fatalf("source-1 = %q (%v), want %q", got, err, body)
	}
	if rows := exportIndex(t, dir); len(rows) != 1 || rows[0]["id"] != "../escape" || rows[0]["file"] != "source-1" {
		t.Fatalf("index rows = %v", rows)
	}
}

// TestEvidenceExportCleansOnFailure is BO50: a write fault on the second source removes the
// files the export created, and the directory only when the export created it.
func TestEvidenceExportCleansOnFailure(t *testing.T) {
	root, pack := publishExportFixture(t)
	for _, test := range []struct {
		name   string
		exists bool
	}{{"created directory", false}, {"existing empty directory", true}} {
		t.Run(test.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "export")
			if test.exists {
				if err := os.Mkdir(dir, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			calls := 0
			fault := func(w io.Writer, data []byte) error {
				calls++
				if calls < 2 {
					return chargeevidence.WriteAll(w, data)
				}
				if _, err := w.Write(data[:len(data)/2]); err != nil {
					return err
				}
				return errors.New("injected write fault")
			}
			out, code := evidencecmd.Export(root, pack.Identity(), dir, fault)
			if code != 1 || calls != 2 || !strings.Contains(out, "injected write fault") {
				t.Fatalf("faulted export = (%d) after %d writes:\n%s", code, calls, out)
			}
			if !test.exists {
				assertAbsent(t, dir)
				return
			}
			if got := exportEntries(t, dir); len(got) != 0 {
				t.Fatalf("faulted export left %v in the existing directory", got)
			}
		})
	}
}
