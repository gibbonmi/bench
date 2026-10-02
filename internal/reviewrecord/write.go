package reviewrecord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	case errors.Is(err, ErrMissing):
		err = fmt.Errorf("%s holds no record: %w", path, ErrNoChunkEntry)
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
	if _, err := graded(path, bounds.ClassifyBytes(bounds.Read(bytes.NewReader(rendered), bounds.ControlRecordLimit))); err != nil {
		return err
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
	digest, plan, err := atSource(root, spec, tip)
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

// atSource reads the tree of commit, a full commit ID, and returns its source digest and
// the plan in it.
func atSource(root, spec, commit string) (string, Plan, error) {
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", "--quiet", commit+"^{tree}")
	if err != nil {
		return "", Plan{}, fmt.Errorf("unreadable tree of commit %s", commit)
	}
	digest, err := SourceDigest(root, tree, spec)
	if err != nil {
		return "", Plan{}, err
	}
	plan, err := ReadPlan(root, tree, spec)
	return digest, plan, err
}

// ErrNoChunkEntry marks a result that needs a current chunk entry first: the record is
// absent, the chunk has no entry, or the result names a source that the entry does not
// hold.
var ErrNoChunkEntry = errors.New("the result needs a current chunk entry")

// ErrProbe marks a result whose probe does not match the planned probe of its requirement.
var ErrProbe = errors.New("the probe must match the plan")

// VerificationCall is the caller's part of one completed verification result. An empty
// Chunk selects the completion list, and an empty Source, a full commit ID otherwise,
// selects the tip of the chunk entry. Evidence supplies the ID, the performer, the model,
// the effort, and the native result. Probe supplies the outcome, the exit code, and the
// restore of a probe, and it is nil when the call names no probe.
type VerificationCall struct {
	Chunk, Source, Requirement string
	ExitCode                   int
	Evidence                   Evidence
	Probe                      *Probe
}

// RecordVerification appends one completed verification result to the chunk list or to
// the completion list. The plan at the source tree supplies the requirement, its command,
// its probe mutation, and the role; the source tree supplies the source digest. A chunk
// result must name the source of its chunk entry. The answer is the entry as written.
func RecordVerification(root, spec string, call VerificationCall) (Verification, error) {
	var entry Verification
	err := write(root, spec, nil, func(record *Record) error {
		list, source := &record.Completion.Verification, call.Source
		var chunk *Chunk
		if call.Chunk != "" {
			for i := range record.Chunks {
				if record.Chunks[i].ID == call.Chunk {
					chunk = &record.Chunks[i]
				}
			}
			if chunk == nil {
				return fmt.Errorf("chunk %s has no entry: %w", call.Chunk, ErrNoChunkEntry)
			}
			list = &chunk.Verification
			if source == "" {
				source = chunk.Tip
			}
		}
		digest, plan, err := atSource(root, spec, source)
		if err != nil {
			return err
		}
		if chunk != nil && digest != chunk.SourceDigest {
			return fmt.Errorf("source %s differs from the source of chunk %s: %w", source, chunk.ID, ErrNoChunkEntry)
		}
		requirements := plan.FinalVerification
		if chunk != nil {
			planned := findChunk(plan, chunk.ID)
			if planned == nil {
				return fmt.Errorf("chunk %s is not a chunk of the plan at %s", chunk.ID, source)
			}
			requirements = planned.Verification
		}
		var requirement *Requirement
		ids := []string{}
		for i := range requirements {
			ids = append(ids, requirements[i].ID)
			if requirements[i].ID == call.Requirement {
				requirement = &requirements[i]
			}
		}
		if requirement == nil {
			return fmt.Errorf("requirement %s is not planned; the planned requirements are %s", call.Requirement, strings.Join(ids, ", "))
		}
		_, role, err := verifier(plan, *record, *requirement, chunk == nil)
		if err != nil {
			return err
		}
		if requirement.Probe == "" && call.Probe != nil {
			return fmt.Errorf("requirement %s plans no probe: %w", requirement.ID, ErrProbe)
		}
		if requirement.Probe != "" && call.Probe == nil {
			return fmt.Errorf("requirement %s plans the probe %q: %w", requirement.ID, requirement.Probe, ErrProbe)
		}
		if holdsID(*record, call.Evidence.ID) {
			return fmt.Errorf("evidence ID %s is already in the record", call.Evidence.ID)
		}
		exitCode, outcome := call.ExitCode, "fail"
		if exitCode == 0 {
			outcome = "pass"
		}
		entry = Verification{Evidence: call.Evidence, Requirement: requirement.ID, Command: requirement.Command, ExitCode: &exitCode}
		entry.Role, entry.SourceDigest, entry.State, entry.Outcome = role, digest, "completed", outcome
		if call.Probe != nil {
			probe := *call.Probe
			probe.Mutation, probe.NativeRef = requirement.Probe, call.Evidence.NativeRef
			entry.Probe = &probe
		}
		*list = append(*list, entry)
		return nil
	})
	return entry, err
}

// holdsID reports whether any verification or review result of record has the ID id.
func holdsID(record Record, id string) bool {
	results := append([]Verification{}, record.Completion.Verification...)
	for _, chunk := range record.Chunks {
		results = append(results, chunk.Verification...)
		for _, review := range chunk.Reviews {
			if review.ID == id {
				return true
			}
		}
	}
	for _, result := range results {
		if result.ID == id {
			return true
		}
	}
	return false
}
