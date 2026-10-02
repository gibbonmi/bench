package reviewrecord

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

var objectID = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
var ErrMissing = errors.New("missing review record")

// missingFence is the absent-fence answer. One sentinel answers every fence, so
// the message names the fence that is absent: a missing completion plan that
// reads as a missing review record sends the reader to the wrong file. The
// sentinel stays reachable through Unwrap, so a caller still tests the
// condition with errors.Is.
type missingFence struct{ name string }

func (e missingFence) Error() string { return "missing " + e.name + " fence" }
func (e missingFence) Unwrap() error { return ErrMissing }

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func Parse(data []byte) (Record, error) {
	var record Record
	if err := decode(data, &record); err != nil {
		return record, err
	}
	if record.Version != 1 && record.Version != 2 {
		return record, fmt.Errorf("unsupported version %d", record.Version)
	}
	if _, err := RecordPath(record.Spec); err != nil {
		return record, err
	}
	if record.PlanDigest == "" {
		return record, errors.New("invalid plan digest")
	}
	// A version 2 record carries no editable identity of its own. It refers to
	// the source-bound plan through the plan digest above, so an evidence-only
	// edit cannot change who may verify or reconcile completion.
	if (record.ImplementationSession == "") != (record.Version == 2) {
		return record, errors.New("invalid implementation session; version 1 names one and version 2 defers to its plan")
	}
	ids := map[string]bool{}
	for _, id := range evidenceIDs(record) {
		if id == "" || ids[id] {
			return record, errors.New("invalid duplicate or empty evidence id")
		}
		ids[id] = true
	}
	chunks := map[string]bool{}
	for _, chunk := range record.Chunks {
		if chunk.ID == "" || chunks[chunk.ID] {
			return record, errors.New("invalid duplicate or empty chunk id")
		}
		chunks[chunk.ID] = true
		if !objectID.MatchString(chunk.Base) || !objectID.MatchString(chunk.Tip) || !objectID.MatchString(chunk.SourceDigest) || chunk.PlanDigest == "" {
			return record, fmt.Errorf("chunk %s: invalid frozen source identity", chunk.ID)
		}
		if err := validateVerification(chunk.Verification); err != nil {
			return record, err
		}
		previous := map[string]string{}
		for _, review := range chunk.Reviews {
			if err := validateEvidence(review.Evidence); err != nil {
				return record, err
			}
			if !contains(Axes(), review.Axis) || review.Role != "independent-review" || (record.Version == 1 && review.Performer == record.ImplementationSession) {
				return record, fmt.Errorf("%s: invalid independent review axis or performer", review.ID)
			}
			if !objectID.MatchString(review.Base) || !objectID.MatchString(review.Tip) {
				return record, fmt.Errorf("%s: invalid review pair", review.ID)
			}
			prior := previous[review.Axis]
			if (prior == "" && len(review.Supersedes) != 0) || (prior != "" && (len(review.Supersedes) != 1 || review.Supersedes[0] != prior)) {
				return record, fmt.Errorf("%s: invalid supersession; retain and name the preceding axis occurrence", review.ID)
			}
			previous[review.Axis] = review.ID
			if review.State == "completed" && ((review.Outcome == "pass") != (len(review.FindingIDs) == 0)) {
				return record, fmt.Errorf("%s: inconsistent review outcome and findings", review.ID)
			}
		}
	}
	if record.Completion.State != "" && !occurrenceState(record.Completion.State) {
		return record, errors.New("invalid completion occurrence state")
	}
	if record.Completion.State != "" && record.Completion.State != "pending" && (!objectID.MatchString(record.Completion.SourceDigest) || record.Completion.Performer == "") {
		return record, errors.New("missing terminal completion source or performer")
	}
	if err := validateVerification(record.Completion.Verification); err != nil {
		return record, err
	}
	return record, nil
}

func validateVerification(items []Verification) error {
	for _, item := range items {
		if err := validateEvidence(item.Evidence); err != nil {
			return err
		}
		// Both verification roles parse here. checkVerification grades which
		// role each obligation actually owes, so the grammar stays in one place.
		if !contains(verificationRoles(), item.Role) || item.Requirement == "" || item.Command == "" {
			return fmt.Errorf("%s: invalid author verification", item.ID)
		}
		if item.State != "pending" && item.ExitCode == nil {
			return fmt.Errorf("%s: missing verification exit code", item.ID)
		}
		if item.Probe != nil {
			if item.Probe.Mutation == "" {
				return fmt.Errorf("%s: missing probe mutation", item.ID)
			}
			if err := validateNative(item.Probe.NativeRef); err != nil {
				return fmt.Errorf("%s probe: %w", item.ID, err)
			}
		}
	}
	return nil
}

// evidenceIDs returns the evidence ID of each verification and review result of record,
// from the lists of each chunk and then from the completion list. It is the one walk of
// the evidence lists for the ID rule.
func evidenceIDs(record Record) []string {
	ids := []string{}
	for _, chunk := range record.Chunks {
		for _, result := range chunk.Verification {
			ids = append(ids, result.ID)
		}
		for _, review := range chunk.Reviews {
			ids = append(ids, review.ID)
		}
	}
	for _, result := range record.Completion.Verification {
		ids = append(ids, result.ID)
	}
	return ids
}

func validateEvidence(item Evidence) error {
	if !occurrenceState(item.State) {
		return fmt.Errorf("%s: invalid occurrence state %q", item.ID, item.State)
	}
	if item.State == "pending" {
		return nil
	}
	if item.Performer == "" || item.Model == "" || item.Effort == "" || !objectID.MatchString(item.SourceDigest) || item.Outcome == "" {
		return fmt.Errorf("%s: missing terminal source or performer metadata; use explicit unknown model and effort", item.ID)
	}
	if err := validateNative(item.NativeRef); err != nil {
		return fmt.Errorf("%s: %w", item.ID, err)
	}
	return nil
}

func validateNative(ref NativeRef) error {
	if ref.Ref == "" || strings.ContainsAny(ref.Ref, "\x00\r\n") || strings.Contains(ref.Ref, "../") || ref.Excerpt == "" || ref.Digest != Digest([]byte(ref.Excerpt)) {
		return errors.New("invalid native result reference, embedded excerpt, or excerpt digest")
	}
	return nil
}

func occurrenceState(state string) bool {
	return state == "pending" || state == "completed" || state == "failed" || state == "skipped"
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func decode(data []byte, value any) error {
	c := bounds.ClassifyBytes(bounds.Read(bytes.NewReader(data), bounds.ControlRecordLimit))
	if c.State != bounds.StateParsed {
		return fmt.Errorf("invalid record bytes: %s %s", c.State, c.Reason)
	}
	tokens := json.NewDecoder(bytes.NewReader(data))
	if err := uniqueJSON(tokens); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("invalid JSON: trailing value")
	}
	return nil
}

func uniqueJSON(d *json.Decoder) error {
	token, err := d.Token()
	if err != nil {
		return err
	}
	delimiter, nested := token.(json.Delim)
	if !nested {
		return nil
	}
	seen := map[string]bool{}
	for d.More() {
		if delimiter == '{' {
			token, err := d.Token()
			if err != nil {
				return err
			}
			key, ok := token.(string)
			if !ok || seen[key] {
				return errors.New("duplicate JSON key")
			}
			seen[key] = true
		}
		if err := uniqueJSON(d); err != nil {
			return err
		}
	}
	_, err = d.Token()
	return err
}

func fenced(data []byte, name string) ([]byte, error) {
	span, err := locate(data, name)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(data[span.start:span.end], []byte("\n")), nil
}

// fenceSpan holds the byte offsets of the payload lines of one fence. start is
// the first byte after the opening line, and end is the first byte of the
// closing line, so the span keeps the newline of its last payload line.
type fenceSpan struct{ start, end int }

// fenceMarker opens a fence when the fence name follows it, and it closes the
// fence on a line of its own.
const fenceMarker = "```"

// locate is the one fence scanner. fenced reads the payload through it, and
// Render replaces the span that it reports, so a reader and a writer cannot
// disagree about where a fence starts or ends.
func locate(data []byte, name string) (fenceSpan, error) {
	var span fenceSpan
	active, found := false, false
	for offset := 0; offset <= len(data); {
		length := bytes.IndexByte(data[offset:], '\n')
		if length < 0 {
			length = len(data) - offset
		}
		line, next := string(data[offset:offset+length]), offset+length+1
		if line == fenceMarker+name {
			if active || found {
				return span, fmt.Errorf("invalid duplicate %s fence", name)
			}
			active, found, span.start = true, true, next
		} else if active && line == fenceMarker {
			active, span.end = false, offset
		}
		offset = next
	}
	if active {
		return span, fmt.Errorf("invalid unterminated %s fence", name)
	}
	if !found {
		return span, missingFence{name}
	}
	return span, nil
}
