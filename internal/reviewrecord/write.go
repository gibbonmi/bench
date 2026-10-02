package reviewrecord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// Render returns document with record as the payload of its review-record
// fence. It does not validate the record, because the test fixtures render
// invalid records on purpose; the write transaction validates. An unterminated
// or a duplicate fence is an error, because a replacement would leave a fence
// that the reader refuses.
func Render(document []byte, record Record) ([]byte, error) {
	var payload bytes.Buffer
	encoder := json.NewEncoder(&payload)
	// A reader reads "&", "<", and ">" in an excerpt as written. The encoder
	// still escapes every control byte, so each payload line stays one line.
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(record); err != nil {
		return nil, err
	}
	span, err := locate(document, recordFence)
	if err == nil {
		return bytes.Join([][]byte{document[:span.start], payload.Bytes(), document[span.end:]}, nil), nil
	}
	if !errors.Is(err, ErrMissing) {
		return nil, err
	}
	head := "# Review outcomes\n"
	if document != nil {
		head = string(document)
		if !bytes.HasSuffix(document, []byte("\n")) {
			head += "\n"
		}
	}
	return []byte(head + "\n" + fenceMarker + recordFence + "\n" + payload.String() + fenceMarker + "\n"), nil
}

// write is the one write transaction on the record of spec under root. It reads the
// record through the reader, applies change, renders the document, and grades the
// rendered bytes as the reader grades a file: the bounded classification, then the
// parse. Only then does it replace the file through a temporary file in
// the same directory. A nil create refuses an absent record; otherwise create supplies
// the record that an absent file starts from. A refusal at any step leaves the record
// bytes unchanged and leaves no temporary file.
func write(root, spec string, create func() (Record, error), change func(*Record) error) error {
	path, err := RecordPath(spec)
	if err != nil {
		return err
	}
	data, err := readFile(root, path)
	var record Record
	switch {
	case errors.Is(err, ErrMissing) && create != nil:
		record, err = create()
	case err == nil:
		record, err = parseRecord(data, spec)
	}
	if err != nil {
		return err
	}
	if err := change(&record); err != nil {
		return err
	}
	rendered, err := Render(data, record)
	if err != nil {
		return err
	}
	if c := bounds.ClassifyBytes(bounds.Read(bytes.NewReader(rendered), bounds.ControlRecordLimit)); c.State != bounds.StateParsed {
		return fmt.Errorf("invalid record %q: %s %s", path, c.State, c.Reason)
	}
	if _, err := parseRecord(rendered, spec); err != nil {
		return err
	}
	return replace(filepath.Join(root, filepath.FromSlash(path)), rendered)
}

// replace writes data to a temporary file beside path and renames it over path, so a
// reader sees either the old bytes or the new bytes. It creates the directory of path
// when that directory is absent.
func replace(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	_, err = temp.Write(data)
	if err == nil {
		err = temp.Chmod(0o644)
	}
	if err == nil {
		err = temp.Sync()
	}
	if closed := temp.Close(); err == nil {
		err = closed
	}
	if err == nil {
		err = os.Rename(temp.Name(), path)
	}
	if err != nil {
		_ = os.Remove(temp.Name())
	}
	return err
}

// RecordChunk writes the entry of planned chunk id for the frozen pair base..tip, both
// full commit IDs. The tip tree supplies the source digest, the plan digest, and the
// acceptance rows. A recorded chunk keeps its results, and a new chunk follows the
// recorded chunks. Only a version 2 plan starts a new record. The answer is the action,
// created, added, or updated, and the entry as written.
func RecordChunk(root, spec, id, base, tip string) (string, Chunk, error) {
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", "--quiet", tip+"^{tree}")
	if err != nil {
		return "", Chunk{}, fmt.Errorf("unreadable tree of tip %s", tip)
	}
	digest, err := SourceDigest(root, tree, spec)
	if err != nil {
		return "", Chunk{}, err
	}
	plan, err := ReadPlan(root, tree, spec)
	if err != nil {
		return "", Chunk{}, err
	}
	action, entry := "", Chunk{}
	create := func() (Record, error) {
		if plan.Version != 2 {
			return Record{}, fmt.Errorf("a new record needs a version 2 completion plan, and the plan at %s has version %d", tip, plan.Version)
		}
		action = "created"
		return Record{Version: 2, Spec: spec, PlanDigest: plan.Digest, Completion: Completion{State: "pending", Reconciliation: map[string]string{}, Verification: []Verification{}}}, nil
	}
	err = write(root, spec, create, func(record *Record) error {
		planned := findChunk(plan, id)
		if planned == nil {
			return fmt.Errorf("chunk %s is not a chunk of the plan at %s", id, tip)
		}
		index := len(record.Chunks)
		for i := range record.Chunks {
			if record.Chunks[i].ID == id {
				index, action = i, "updated"
			}
		}
		if index == len(record.Chunks) {
			record.Chunks = append(record.Chunks, Chunk{ID: id, Verification: []Verification{}, Reviews: []Review{}})
			if action == "" {
				action = "added"
			}
		}
		chunk := &record.Chunks[index]
		chunk.Base, chunk.Tip, chunk.PlanDigest, chunk.SourceDigest = base, tip, plan.Digest, digest
		chunk.AcceptanceRows = append([]string{}, planned.Rows...)
		entry = *chunk
		return nil
	})
	return action, entry, err
}
