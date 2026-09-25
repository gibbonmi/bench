package chargeevidence

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"syscall"

	"github.com/gibbonmi/bench/internal/toon"
)

// RefuseExport names an export directory or file the export cannot use or complete.
const RefuseExport = "export-refused"

// Export file layout. Each source file name derives from its ordinal alone, so no source
// identifier can choose a path.
const (
	exportSourcePrefix = "source-"
	exportIndexName    = "index.toon"
	exportIndexBlock   = "sources"
)

// exportIndexFields are the columns of the one index table.
var exportIndexFields = []string{"ordinal", "id", "bytes", "file"}

// ExportWrite copies the bytes of one export file to w. The export opens, records, and closes
// each file, so a write that fails leaves only a file the cleanup removes.
type ExportWrite func(w io.Writer, data []byte) error

// WriteAll is the ExportWrite of a real export: it writes every byte of data to w.
func WriteAll(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}

// SourceAt returns the identifier and the body of the source at a one-based ordinal. It
// returns a body only after the body verifies.
type SourceAt func(ordinal int) (id string, body []byte, err error)

// Exported counts one complete export: its directory, its source files, and their bytes.
type Exported struct {
	Dir            string
	Sources, Bytes int
}

// Export writes every declared source of the artifact into dir. Each source reaches write
// only after each of its page digests and its source digest verify.
func (a *Artifact) Export(dir string, write ExportWrite) (Exported, error) {
	return ExportSources(dir, len(a.m.Sources), a.verifiedSource, write)
}

// verifiedSource reads each page of one source through the checked page read, then checks
// the reconstructed body against the source digest.
func (a *Artifact) verifiedSource(ordinal int) (string, []byte, error) {
	source := a.m.Sources[ordinal-1]
	body := make([]byte, 0, source.Bytes)
	for index := range a.pagesOf(source.ID) {
		fragment, err := a.ReadWithin(Cursor{Identity: a.identity, Ordinal: ordinal, Index: index})
		if err != nil {
			return "", nil, err
		}
		body = append(body, fragment.Content...)
	}
	if err := checkSource(source, body); err != nil {
		return "", nil, err
	}
	return source.ID, body, nil
}

// ExportSources writes count sources into dir, then one index table. dir must be absent or
// an empty directory, and it must not be a symlink. next supplies each source in ordinal
// order, and a source is written only after next returns it. On any failure the export
// removes each file it created, and it removes dir when it created dir.
func ExportSources(dir string, count int, next SourceAt, write ExportWrite) (exported Exported, err error) {
	target, err := openExportDir(dir)
	if err != nil {
		return Exported{}, err
	}
	defer func() {
		if err == nil {
			err = target.close()
		}
		if err != nil {
			err = target.discard(err)
			exported = Exported{}
		}
	}()
	exported.Dir = dir
	rows := make([][]any, 0, count)
	for ordinal := 1; ordinal <= count; ordinal++ {
		id, body, err := next(ordinal)
		if err != nil {
			return Exported{}, err
		}
		name := exportSourcePrefix + strconv.Itoa(ordinal)
		if err := target.create(name, body, write); err != nil {
			return Exported{}, err
		}
		rows = append(rows, []any{ordinal, id, len(body), name})
		exported.Sources++
		exported.Bytes += len(body)
	}
	index, err := toon.TableTyped(exportIndexBlock, exportIndexFields, rows)
	if err != nil {
		return Exported{}, refuse(RefuseExport, "the export index does not encode: %v", err)
	}
	return exported, target.create(exportIndexName, []byte(index), write)
}

// exportDir is one opened export directory and the files the export created in it.
type exportDir struct {
	path    string
	root    *os.Root
	created bool
	files   []string
}

// openExportDir opens path for an export. It creates an absent directory, and it refuses a
// symlink, a non-directory, and a directory that holds an entry.
func openExportDir(path string) (*exportDir, error) {
	d := &exportDir{path: path}
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(path, 0o755); err != nil {
			return nil, refuse(RefuseExport, "the export directory is not creatable: %v", err)
		}
		d.created = true
		info, err = os.Lstat(path)
	}
	if err != nil {
		return nil, d.discard(refuse(RefuseExport, "the export directory is not inspectable: %v", err))
	}
	if !info.IsDir() {
		return nil, refuse(RefuseExport, "the export path is %s, not a directory", kindOf(info))
	}
	if d.root, err = os.OpenRoot(path); err != nil {
		return nil, d.discard(refuse(RefuseExport, "the export directory is not openable: %v", err))
	}
	// OpenRoot follows a symlink, so the opened directory must be the one Lstat inspected.
	if opened, err := d.root.Stat("."); err != nil || !os.SameFile(info, opened) {
		return nil, d.discard(refuse(RefuseExport, "the export directory changed while it was opened"))
	}
	if err := d.empty(); err != nil {
		return nil, d.discard(err)
	}
	return d, nil
}

func (d *exportDir) empty() error {
	listing, err := d.root.Open(".")
	if err != nil {
		return refuse(RefuseExport, "the export directory is not listable: %v", err)
	}
	defer listing.Close()
	names, err := listing.Readdirnames(1)
	if len(names) > 0 {
		return refuse(RefuseExport, "the export directory is not empty")
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return refuse(RefuseExport, "the export directory is not listable: %v", err)
	}
	return nil
}

// create writes one new file. The file is recorded as soon as it exists, so a write that
// fails still leaves a file the cleanup removes.
func (d *exportDir) create(name string, data []byte, write ExportWrite) error {
	file, err := d.root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, 0o644)
	if err != nil {
		return refuse(RefuseExport, "export file %s is not creatable: %v", name, err)
	}
	d.files = append(d.files, name)
	err = write(file, data)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return refuse(RefuseExport, "export file %s is not writable: %v", name, err)
	}
	return nil
}

func (d *exportDir) close() error {
	if err := d.root.Close(); err != nil {
		return refuse(RefuseExport, "the export directory does not close: %v", err)
	}
	d.root = nil
	return nil
}

// discard removes each file the export created, then the directory when the export created
// it. It returns cause, extended with each path the removal left behind.
func (d *exportDir) discard(cause error) error {
	var left []string
	if d.root != nil {
		for _, name := range d.files {
			if err := d.root.Remove(name); err != nil && !errors.Is(err, fs.ErrNotExist) {
				left = append(left, name)
			}
		}
		d.root.Close()
	}
	if d.created {
		if err := os.Remove(d.path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			left = append(left, d.path)
		}
	}
	if len(left) == 0 {
		return cause
	}
	return fmt.Errorf("%w; the cleanup left %q", cause, left)
}
