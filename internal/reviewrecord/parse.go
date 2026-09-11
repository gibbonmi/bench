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

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func Parse(data []byte) (Record, error) {
	var record Record
	if err := decode(data, &record); err != nil {
		return record, err
	}
	if record.Version != 1 {
		return record, fmt.Errorf("unsupported version %d", record.Version)
	}
	if _, err := RecordPath(record.Spec); err != nil {
		return record, err
	}
	if record.PlanDigest == "" || record.ImplementationSession == "" {
		return record, errors.New("invalid plan digest or implementation session")
	}
	ids := map[string]bool{}
	chunks := map[string]bool{}
	for _, chunk := range record.Chunks {
		if chunk.ID == "" || chunks[chunk.ID] {
			return record, errors.New("invalid duplicate or empty chunk id")
		}
		chunks[chunk.ID] = true
		if !objectID.MatchString(chunk.Base) || !objectID.MatchString(chunk.Tip) || !objectID.MatchString(chunk.SourceDigest) || chunk.PlanDigest == "" {
			return record, fmt.Errorf("chunk %s: invalid frozen source identity", chunk.ID)
		}
		if err := validateVerification(chunk.Verification, ids); err != nil {
			return record, err
		}
		previous := map[string]string{}
		for _, review := range chunk.Reviews {
			if err := validateEvidence(review.Evidence, ids); err != nil {
				return record, err
			}
			if !contains(Axes(), review.Axis) || review.Role != "independent-review" || review.Performer == record.ImplementationSession {
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
	if err := validateVerification(record.Completion.Verification, ids); err != nil {
		return record, err
	}
	return record, nil
}

func validateVerification(items []Verification, ids map[string]bool) error {
	for _, item := range items {
		if err := validateEvidence(item.Evidence, ids); err != nil {
			return err
		}
		if item.Role != "author-verification" || item.Requirement == "" || item.Command == "" {
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

func validateEvidence(item Evidence, ids map[string]bool) error {
	if item.ID == "" || ids[item.ID] {
		return errors.New("invalid duplicate or empty evidence id")
	}
	ids[item.ID] = true
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
	var payload []string
	active, found := false, false
	for _, line := range strings.Split(string(data), "\n") {
		if line == "```"+name {
			if active || found {
				return nil, fmt.Errorf("invalid duplicate %s fence", name)
			}
			active, found = true, true
		} else if active && line == "```" {
			active = false
		} else if active {
			payload = append(payload, line)
		}
	}
	if active {
		return nil, fmt.Errorf("invalid unterminated %s fence", name)
	}
	if !found {
		return nil, ErrMissing
	}
	return []byte(strings.Join(payload, "\n")), nil
}
