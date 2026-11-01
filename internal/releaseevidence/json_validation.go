package releaseevidence

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/jsonfile"
)

// decodeStrict decodes one complete release document into value. Invalid JSON outranks a
// duplicate key, as the release readers report it.
func decodeStrict(data []byte, value any) error {
	if !json.Valid(data) {
		return errors.New("invalid JSON")
	}
	return decodeDocument(data, value)
}

// decodeDocument decodes data through jsonfile.DecodeDocument with the release diagnostics. A
// duplicate key outranks a suffix, and a suffix outranks a typed refusal. After a decode, it
// refuses a document that an untyped decode refuses, such as a number out of float64 range.
func decodeDocument(data []byte, value any) error {
	err := jsonfile.DecodeDocument(data, value)
	switch {
	case err == nil:
		// The typed decode skips the numbers inside a json.RawMessage value. This untyped decode
		// reads them.
		return json.Unmarshal(data, new(any))
	case errors.Is(err, jsonfile.ErrDuplicateKey):
		return duplicateKeyError{err}
	case errors.Is(err, jsonfile.ErrTrailingData):
		// The encoding/json validity error names the first byte after the value.
		if invalid := json.Unmarshal(data, new(json.RawMessage)); invalid != nil {
			return invalid
		}
	}
	// The release diagnostic names the encoding/json cause without the jsonfile context.
	if cause := errors.Unwrap(err); cause != nil {
		return cause
	}
	return err
}

// duplicateKeyError words a jsonfile duplicate-key refusal as the release readers report it.
// errors.Is still finds jsonfile.ErrDuplicateKey through it.
type duplicateKeyError struct{ cause error }

func (e duplicateKeyError) Error() string {
	return "duplicate JSON key" + strings.TrimPrefix(e.cause.Error(), jsonfile.ErrDuplicateKey.Error())
}

func (e duplicateKeyError) Unwrap() error { return e.cause }
