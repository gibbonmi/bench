package reviewrecord

import (
	"bytes"
	"encoding/json"
	"errors"
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
