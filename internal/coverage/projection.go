package coverage

import (
	"bytes"
	"encoding/json"
	"errors"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/markdown"
)

// StorySchema names the one versioned serialization of a story projection.
const StorySchema = "bench-coverage-story/1"

// storyPayload is the serialized story projection. The seam_cell field holds the authored
// seam cell, not a claim that an outline location is a blessed seam.
type storyPayload struct {
	Schema  string       `json:"schema"`
	Present bool         `json:"present"`
	Story   *storyRecord `json:"story"`
	Rows    []rowRecord  `json:"rows"`
}

type storyRecord struct {
	Number int        `json:"number"`
	Text   text       `json:"text"`
	Span   spanRecord `json:"span"`
}

// rowRecord is one serialized row. A legacy schema without a row column omits the id.
type rowRecord struct {
	ID        text       `json:"id,omitempty"`
	Stories   []int      `json:"stories"`
	Behavior  text       `json:"behavior"`
	SeamCell  text       `json:"seam_cell"`
	Rationale text       `json:"rationale"`
	Span      spanRecord `json:"span"`
}

// text is one text field of the serialization. Its encoding refuses invalid UTF-8 rather than
// become a replacement character, so the serializer owns the UTF-8 refusal of each text field.
type text string

var errInvalidText = errors.New("invalid UTF-8")

func (t text) MarshalText() ([]byte, error) {
	if !utf8.ValidString(string(t)) {
		return nil, errInvalidText
	}
	return []byte(t), nil
}

// spanRecord is one serialized source span.
type spanRecord struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func record(span markdown.Span) spanRecord { return spanRecord{Start: span.Start, End: span.End} }

// Encode answers the StorySchema serialization of the projection: one JSON object with the
// story's exact text, its matching rows in authored order, and their original source spans.
// An absent story has a null story and an empty row array. Story refuses a projection that
// encode cannot carry, so a projection that Story answers always encodes.
func (p StoryProjection) Encode() []byte {
	out, _ := p.encode()
	return out
}

// encode serializes the projection, and refuses a text field that holds invalid UTF-8.
func (p StoryProjection) encode() ([]byte, error) {
	payload := storyPayload{Schema: StorySchema, Present: p.Present, Rows: []rowRecord{}}
	if p.Present {
		payload.Story = &storyRecord{Number: p.Story.Number, Text: text(p.Story.Text), Span: record(p.Story.Span)}
	}
	for _, row := range p.Rows {
		payload.Rows = append(payload.Rows, rowRecord{ID: text(row.ID), Stories: row.Stories, Behavior: text(row.Behavior), SeamCell: text(row.Seam), Rationale: text(row.Rationale), Span: record(row.Span)})
	}
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(payload); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
