package reviewrecord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
)

// Render returns document with record as the payload of its review-record
// fence. It does not validate the record, because the test fixtures render
// invalid records on purpose; the write transaction validates. An unterminated
// or a duplicate fence is an error, because a replacement would leave a fence
// that the reader refuses.
func Render(document []byte, record Record) ([]byte, error) {
	data, err := payload(record, "  ")
	if err != nil {
		return nil, err
	}
	span, err := locate(document, recordFence)
	if err == nil {
		return span.splice(document, data), nil
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
	return []byte(head + "\n" + fenceMarker + recordFence + "\n" + string(data) + fenceMarker + "\n"), nil
}

// payload encodes value as fence payload lines, each indented by indent.
func payload(value any, indent string) ([]byte, error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	// A reader reads "&", "<", and ">" in an excerpt as written. The encoder
	// still escapes every control byte, so each payload line stays one line.
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	err := encoder.Encode(value)
	return data.Bytes(), err
}

// splice returns document with the payload lines of span replaced by data.
func (s fenceSpan) splice(document, data []byte) []byte {
	return bytes.Join([][]byte{document[:s.start], data, document[s.end:]}, nil)
}

// replaceFile is the one replacement operation of every review-record file write: the
// native durable replacement. A test binds its own operation for one test.
var replaceFile = durablefile.Replace

// write is the one write transaction on the record of spec under root. It reads the
// record through the reader, applies change, renders the document, and grades the
// rendered bytes as the reader grades a file: the bounded classification, then the
// parse. Only then does it replace the file through replace. For an absent file at
// path, create supplies the record that the file starts from or refuses it. A refusal
// before the replacement leaves the record bytes unchanged and leaves no temporary file.
func write(root, spec string, create func(path string) (Record, error), change func(*Record) error) error {
	path, err := RecordPath(spec)
	if err != nil {
		return err
	}
	data, err := readFile(root, path)
	var record Record
	switch {
	case errors.Is(err, ErrMissing):
		record, err = create(path)
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

// absent refuses an absent record at path. The chunk form writes its first entry, for
// chunk when the caller names one.
func absent(chunk string) func(path string) (Record, error) {
	return func(path string) (Record, error) {
		return Record{}, Repair{Form: FormChunk, Chunk: chunk, Err: fmt.Errorf("%s holds no record: %w", path, ErrNoChunkEntry)}
	}
}

// replace creates the directory of path when that directory is absent, then replaces path
// with data at mode 0644 through replaceFile, so a reader sees either the old bytes or the
// new bytes. It returns the replacement error unchanged and rolls back nothing.
func replace(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return replaceFile(path, data, 0o644)
}

// RecordChunk writes the entry of planned chunk id for the frozen pair base..tip, both
// full commit IDs, through changeChunk. The tip tree supplies the source digest, the plan
// digest, and the acceptance rows, and createRecord starts an absent record. The answer is
// the action, created, added, or updated, and the entry as written.
func RecordChunk(root, spec, id, base, tip string) (string, Chunk, error) {
	digest, plan, err := atSource(root, spec, tip)
	if err != nil {
		return "", Chunk{}, err
	}
	action, entry, created := "", Chunk{}, false
	create := func(string) (Record, error) {
		record, err := createRecord(spec, plan, tip)
		created = err == nil
		return record, err
	}
	err = write(root, spec, create, func(record *Record) error {
		action, entry, err = changeChunk(record, plan, digest, id, base, tip)
		return err
	})
	if created {
		action = "created"
	}
	return action, entry, err
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

// RecordVerification appends one completed verification result of spec under root
// through changeVerification. An absent record refuses. The answer is the entry as written.
func RecordVerification(root, spec string, call VerificationCall) (Verification, error) {
	var entry Verification
	err := write(root, spec, absent(call.Chunk), func(record *Record) error {
		var err error
		entry, err = changeVerification(root, spec, record, call)
		return err
	})
	return entry, err
}

// ChunkRequirement answers requirement id of chunk of plan, the plan at source.
func ChunkRequirement(plan Plan, source, chunk, id string) (Requirement, error) {
	if planned := findChunk(plan, chunk); planned != nil {
		return plannedRequirement(planned.Verification, id)
	}
	return Requirement{}, fmt.Errorf("chunk %s is not a chunk of the plan at %s", chunk, source)
}

// plannedRequirement answers requirement id, or refuses with the planned requirements named.
func plannedRequirement(requirements []Requirement, id string) (Requirement, error) {
	ids := []string{}
	for _, requirement := range requirements {
		if requirement.ID == id {
			return requirement, nil
		}
		ids = append(ids, requirement.ID)
	}
	return Requirement{}, fmt.Errorf("requirement %s is not planned; the planned requirements are %s", id, strings.Join(ids, ", "))
}

// entryOf returns the entry of chunk id in record. A chunk without an entry needs the
// chunk form first.
func entryOf(record *Record, id string) (*Chunk, error) {
	for i := range record.Chunks {
		if record.Chunks[i].ID == id {
			return &record.Chunks[i], nil
		}
	}
	return nil, Repair{Form: FormChunk, Chunk: id, Err: fmt.Errorf("chunk %s has no entry: %w", id, ErrNoChunkEntry)}
}

// ReviewCall is the caller's part of one completed independent review result. Chunk
// names the recorded chunk, or Fold names the recorded fold. Axis names the review axis,
// and Findings holds the finding IDs in the order of the call. Evidence supplies the ID,
// the performer, the model, the effort, and the native result.
type ReviewCall struct {
	Chunk, Fold, Axis string
	Findings          []string
	Evidence          Evidence
}

// RecordReview appends one completed independent review result to the review list of a
// recorded chunk or fold. The entry supplies the base, the tip, and the source digest. The
// result supersedes the last result of the same axis in that entry. The answer is the entry
// as written.
func RecordReview(root, spec string, call ReviewCall) (Review, error) {
	var entry Review
	err := write(root, spec, absent(call.Chunk), func(record *Record) error {
		pair, err := reviewTarget(record, call)
		if err != nil {
			return err
		}
		if err := unclaimed(*record, call.Evidence.ID); err != nil {
			return err
		}
		supersedes := []string{}
		if prior := currentReview(*pair.reviews, call.Axis); prior != nil {
			supersedes = []string{prior.ID}
		}
		entry = Review{Evidence: call.Evidence, Axis: call.Axis, Base: pair.base, Tip: pair.tip, FindingIDs: append([]string{}, call.Findings...), Supersedes: supersedes}
		entry.Role, entry.SourceDigest, entry.State, entry.Outcome = "independent-review", pair.digest, "completed", "pass"
		if len(call.Findings) != 0 {
			entry.Outcome = "fail"
		}
		*pair.reviews = append(*pair.reviews, entry)
		return nil
	})
	return entry, err
}

// ConfirmationCall is the caller's part of one coordinator confirmation. Chunk names the
// recorded chunk, and Findings holds the finding IDs that the repair closes. Evidence supplies
// the ID, the performer, the model, the effort, and the native result. Probe supplies the
// caller's mutation text and the outcome, the exit code, and the restore of the probe.
type ConfirmationCall struct {
	Chunk    string
	Findings []string
	Evidence Evidence
	Probe    Probe
}

// RecordConfirmation appends one completed coordinator confirmation to a recorded chunk.
// FindReviewedPair supplies the reviewed pair, the chunk entry supplies the base, the repaired
// tip, and the source digest, and the entry supersedes the prior confirmation of the chunk.
// It refuses findings other than the open findings at the reviewed pair, and a probe that did
// not bite or did not restore. The answer is the entry as written and the reviewed pair.
func RecordConfirmation(root, spec string, call ConfirmationCall) (Confirmation, ReviewedPair, error) {
	var entry Confirmation
	var reviewed ReviewedPair
	err := write(root, spec, absent(call.Chunk), func(record *Record) error {
		chunk, err := entryOf(record, call.Chunk)
		if err != nil {
			return err
		}
		if reviewed, err = FindReviewedPair(root, spec, *chunk); err != nil {
			return err
		}
		scope := chunk.confirmationScope(call.Evidence.ID)
		if err := closesOpen(scope, call.Findings, reviewed); err != nil {
			return err
		}
		if err := probeHolds(scope, call.Probe); err != nil {
			return err
		}
		if err := unclaimed(*record, call.Evidence.ID); err != nil {
			return err
		}
		supersedes := []string{}
		if prior := len(chunk.Confirmations); prior != 0 {
			supersedes = []string{chunk.Confirmations[prior-1].ID}
		}
		probe := call.Probe
		probe.NativeRef = call.Evidence.NativeRef
		entry = Confirmation{Evidence: call.Evidence, Base: chunk.Base, ReviewedTip: reviewed.Tip, RepairedTip: chunk.Tip, FindingIDs: append([]string{}, call.Findings...), Supersedes: supersedes, Probe: probe}
		entry.Role, entry.SourceDigest, entry.State, entry.Outcome = ConfirmationRole, chunk.SourceDigest, "completed", "pass"
		chunk.Confirmations = append(chunk.Confirmations, entry)
		return nil
	})
	return entry, reviewed, err
}

// reviewTarget is the frozen pair of the recorded chunk or fold that call names. A fold
// without an entry needs the fold form first.
func reviewTarget(record *Record, call ReviewCall) (reviewed, error) {
	if call.Fold == "" {
		chunk, err := entryOf(record, call.Chunk)
		if err != nil {
			return reviewed{}, err
		}
		return chunk.reviewed(), nil
	}
	for i := range record.Folds {
		if record.Folds[i].ID == call.Fold {
			return record.Folds[i].reviewed(), nil
		}
	}
	return reviewed{}, Repair{Form: FormFold, Err: fmt.Errorf("fold %s has no entry", call.Fold)}
}

// unclaimed refuses an ID that an evidence entry of record already has, and the refusal
// names that ID. evidenceIDs owns the set of entries.
func unclaimed(record Record, id string) error {
	if slices.Contains(evidenceIDs(record), id) {
		return fmt.Errorf("evidence ID %s is already in the record", id)
	}
	return nil
}
