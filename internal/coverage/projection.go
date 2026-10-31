package coverage

import (
	"bytes"
	"encoding/json"
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
	Number int    `json:"number"`
	Text   string `json:"text"`
	Span   Span   `json:"span"`
}

// rowRecord is one serialized row. A legacy schema without a row column omits the id.
type rowRecord struct {
	ID        string `json:"id,omitempty"`
	Stories   []int  `json:"stories"`
	Behavior  string `json:"behavior"`
	SeamCell  string `json:"seam_cell"`
	Rationale string `json:"rationale"`
	Span      Span   `json:"span"`
}

// Encode answers the StorySchema serialization of the projection: one JSON object with the
// story's exact text, its matching rows in authored order, and their original source spans.
// An absent story has a null story and an empty row array.
func (p StoryProjection) Encode() []byte {
	payload := storyPayload{Schema: StorySchema, Present: p.Present, Rows: []rowRecord{}}
	if p.Present {
		payload.Story = &storyRecord{Number: p.Story.Number, Text: p.Story.Text, Span: p.Story.Span}
	}
	for _, row := range p.Rows {
		payload.Rows = append(payload.Rows, rowRecord{ID: row.ID, Stories: row.Stories, Behavior: row.Behavior, SeamCell: row.Seam, Rationale: row.Rationale, Span: row.Span})
	}
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	// The payload holds only strings, numbers, and booleans, so the encoder cannot fail.
	_ = encoder.Encode(payload)
	return out.Bytes()
}
