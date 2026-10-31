package coverage

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/markdown"
	specref "github.com/gibbonmi/bench/internal/spec"
)

// ErrRefused is the sentinel that each Describe and story-selection refusal wraps. A refusal
// comes before any content.
var ErrRefused = errors.New("coverage selection refused")

// Story is one numbered line of the User stories section: its number, its exact authored
// text without the line terminator, and the span of that text.
type Story struct {
	Number int
	Text   string
	Span   markdown.Span
}

// CoverageRow is one typed data row of the coverage map. ID is empty under a legacy schema
// without a row column. Stories lists every story number that the row references, in
// authored order. The text fields are the trimmed cells that validation reads, and Span is the
// row's original bytes without the line terminator.
type CoverageRow struct {
	ID                        string
	Stories                   []int
	Behavior, Seam, Rationale string
	Span                      markdown.Span
}

// Document is the typed description of one spec's coverage obligations, from the one parse
// that Check validates. MapPresent reports a coverage map, and Map is its original span from
// its heading through the byte before the heading that ends it. Stories and Rows keep
// authored order, and Stories keeps each occurrence of a repeated number.
type Document struct {
	MapPresent bool
	Map        markdown.Span
	Stories    []Story
	Rows       []CoverageRow
}

// Describe describes the coverage obligations of content. A Markdown fault leaves the rest of
// the content unread, so it refuses first. A present map must pass Check, even under the
// historical marker, so a malformed map refuses rather than projecting partial obligations.
// A spec without a map describes its stories and no rows.
func Describe(content []byte) (Document, error) {
	p := parse(content)
	if p.fault != nil {
		return Document{}, fmt.Errorf("%w: %s", ErrRefused, specref.FaultReason(p.fault))
	}
	var document Document
	for _, story := range p.stories {
		document.Stories = append(document.Stories, Story{Number: story.number, Text: string(content[story.span.Start:story.span.End]), Span: story.span})
	}
	if !p.seen {
		return document, nil
	}
	validated := p
	validated.historical = false
	if violations := Check(validated); len(violations) > 0 {
		return Document{}, fmt.Errorf("%w: %s", ErrRefused, strings.Join(violations, "; "))
	}
	document.MapPresent, document.Map = true, p.mapSpan
	s := p.sch
	for _, r := range p.dataRows {
		_, parts, _ := storyReference(s.cell(r, fieldStory))
		document.Rows = append(document.Rows, CoverageRow{
			ID: s.cell(r, fieldRow), Stories: storyNumbers(parts),
			Behavior: s.cell(r, fieldBehavior), Seam: s.cell(r, fieldSeam), Rationale: s.cell(r, fieldWhy),
			Span: r.span,
		})
	}
	return document, nil
}

// StoryProjection is one story's obligations: the story and every row that references it, in
// authored order. An absent story has Present false and no rows.
type StoryProjection struct {
	Present bool
	Story   Story
	Rows    []CoverageRow
}

// Story selects the story numbered number. A number that the stories section repeats refuses
// as ambiguous, although Check counts it once. A number that no story carries answers an
// absent projection. The serialization carries text, so invalid UTF-8 in the story or in a
// cell of a matching row refuses rather than become a replacement character.
func (d Document) Story(number int) (StoryProjection, error) {
	projection := StoryProjection{Rows: []CoverageRow{}}
	matches := 0
	for _, story := range d.Stories {
		if story.Number == number {
			matches++
			projection.Present, projection.Story = true, story
		}
	}
	if matches > 1 {
		return StoryProjection{}, fmt.Errorf("%w: %d story lines are numbered %d", ErrRefused, matches, number)
	}
	if !projection.Present {
		return projection, nil
	}
	if !utf8.ValidString(projection.Story.Text) {
		return StoryProjection{}, invalidText(fmt.Sprintf("story %d", number))
	}
	for _, row := range d.Rows {
		if !slices.Contains(row.Stories, number) {
			continue
		}
		for _, cell := range []string{row.ID, row.Behavior, row.Seam, row.Rationale} {
			if !utf8.ValidString(cell) {
				return StoryProjection{}, invalidText(fmt.Sprintf("the coverage row on bytes %d to %d", row.Span.Start, row.Span.End))
			}
		}
		projection.Rows = append(projection.Rows, row)
	}
	return projection, nil
}

// invalidText refuses a projection whose subject holds invalid UTF-8, and names the exact
// section read that keeps those bytes.
func invalidText(subject string) error {
	return fmt.Errorf("%w: %s holds invalid UTF-8, which the story projection cannot carry exactly; read its exact bytes with %s", ErrRefused, subject, specref.SectionFlag)
}
