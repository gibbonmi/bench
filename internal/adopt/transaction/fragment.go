package transaction

import (
	"bytes"
	"fmt"
)

// Span identifies the managed interval within a file that also contains user content.
type Span struct {
	Offset       int    `json:"offset"`
	BeforeLength int    `json:"before_length"`
	AfterLength  int    `json:"after_length"`
	Digest       string `json:"digest"`
}

func managedSpan(span *Span, before, after image) (*Span, error) {
	if span == nil {
		return nil, nil
	}
	copy := *span
	if before.Kind != "file" || after.Kind != "file" || !within(copy.Offset, copy.BeforeLength, len(before.Data)) || !within(copy.Offset, copy.AfterLength, len(after.Data)) {
		return nil, fmt.Errorf("invalid managed repair span")
	}
	if !bytes.Equal(before.Data[:copy.Offset], after.Data[:copy.Offset]) || !bytes.Equal(before.Data[copy.Offset+copy.BeforeLength:], after.Data[copy.Offset+copy.AfterLength:]) {
		return nil, fmt.Errorf("repair span changes non-managed content")
	}
	copy.Digest = digestBytes(before.Data[copy.Offset : copy.Offset+copy.BeforeLength])
	return &copy, nil
}

func within(offset, length, total int) bool {
	return offset >= 0 && length >= 0 && offset <= total && length <= total-offset
}

func backupData(e entry) []byte {
	if e.Span == nil {
		return e.Before.Data
	}
	return e.Before.Data[e.Span.Offset : e.Span.Offset+e.Span.BeforeLength]
}

func restoredImage(e entry, current image) (image, error) {
	restored := e.Before
	if e.Span == nil {
		return restored, nil
	}
	span := e.Span
	if current.Kind != "file" || !within(span.Offset, span.AfterLength, len(current.Data)) || len(e.Before.Data) != span.BeforeLength {
		return image{}, fmt.Errorf("invalid managed repair span for %s", e.Destination)
	}
	data := make([]byte, 0, len(current.Data)-span.AfterLength+span.BeforeLength)
	data = append(data, current.Data[:span.Offset]...)
	data = append(data, e.Before.Data...)
	data = append(data, current.Data[span.Offset+span.AfterLength:]...)
	if digestBytes(data) != e.Before.Digest {
		return image{}, fmt.Errorf("restored managed span does not match preimage: %s", e.Destination)
	}
	restored.Data = data
	return restored, nil
}
