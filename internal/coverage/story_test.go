package coverage

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// storySpec is an authored spec whose map rows sit apart from blank lines and a fenced
// heading, and whose map ends at the next visible heading.
const (
	storyStories = "## User stories\n\n### Group\n\n1. As an agent, I want one story.\n2. As an agent, I want a second story.\r\n3. As an agent, I want a third story.\n\n"
	storyMap     = "### Acceptance coverage map\n\n| row | story | behavior | seam | why it catches the failure |\n| --- | --- | --- | --- | --- |\n" +
		"| ST1 | 1 | The first behavior | `a_test.go` (`TestA`) | It catches \\| the first failure. |\n" +
		"| ST2 | 2 | The second behavior | planned TestB in b_test.go | It catches the second failure. |\n\n" +
		"| ST3 | 1-2 | A shared behavior with <html> & text | review-owned scope | It catches the shared failure. |\n\n" +
		"```text\n## Fenced\n```\n\n" +
		"| ST4 | 3 | The third behavior | planned TestC in c_test.go | It catches the third failure. |\n\n"
	storySpec = "# Spec\n\n" + storyStories + "## Implementation decisions\n\n" + storyMap + "### Edge inventory\n\n- one\n"
)

// TestExactStoryProjection grades the coverage owner's typed obligations: the accepted map
// span, a story's exact text and every matching typed row in authored order with its source
// span, a repeated member that joins once, the refusals of a malformed map, a Markdown fault,
// invalid UTF-8, and a repeated selected story number, the absent story, and the versioned
// serialization.
func TestExactStoryProjection(t *testing.T) {
	document, err := Describe([]byte(storySpec))
	if err != nil {
		t.Fatal(err)
	}
	if !document.MapPresent || storySpec[document.Map.Start:document.Map.End] != storyMap {
		t.Fatalf("map = %+v selecting %q, want the authored map through the byte before the next heading", document.Map, storySpec[document.Map.Start:document.Map.End])
	}

	projection, err := document.Story(1)
	if err != nil {
		t.Fatal(err)
	}
	if !projection.Present || projection.Story.Text != "1. As an agent, I want one story." || storySpec[projection.Story.Span.Start:projection.Story.Span.End] != projection.Story.Text {
		t.Fatalf("story 1 = %+v, want its exact authored line and span", projection.Story)
	}
	want := []struct {
		id, behavior, seam, rationale, line string
		stories                             []int
	}{
		{"ST1", "The first behavior", "`a_test.go` (`TestA`)", `It catches \| the first failure.`, "| ST1 | 1 | The first behavior | `a_test.go` (`TestA`) | It catches \\| the first failure. |", []int{1}},
		{"ST3", "A shared behavior with <html> & text", "review-owned scope", "It catches the shared failure.", "| ST3 | 1-2 | A shared behavior with <html> & text | review-owned scope | It catches the shared failure. |", []int{1, 2}},
	}
	if len(projection.Rows) != len(want) {
		t.Fatalf("story 1 rows = %+v, want %d rows", projection.Rows, len(want))
	}
	for i, w := range want {
		row := projection.Rows[i]
		if row.ID != w.id || row.Behavior != w.behavior || row.Seam != w.seam || row.Rationale != w.rationale || !reflect.DeepEqual(row.Stories, w.stories) || storySpec[row.Span.Start:row.Span.End] != w.line {
			t.Fatalf("story 1 row %d = %+v, want %+v", i, row, w)
		}
	}
	if crlf, _ := document.Story(2); crlf.Story.Text != "2. As an agent, I want a second story." || len(crlf.Rows) != 2 {
		t.Fatalf("story 2 = %+v, want its text without the CRLF terminator and two rows", crlf)
	}

	var payload map[string]any
	if err := json.Unmarshal(projection.Encode(), &payload); err != nil {
		t.Fatal(err)
	}
	rows := payload["rows"].([]any)
	first := rows[1].(map[string]any)
	if payload["schema"] != StorySchema || payload["present"] != true || payload["story"].(map[string]any)["text"] != projection.Story.Text || len(rows) != 2 || first["seam_cell"] != "review-owned scope" || first["behavior"] != want[1].behavior {
		t.Fatalf("payload = %v, want the versioned story, its rows, and named seam cells", payload)
	}

	absent, err := document.Story(9)
	if err != nil || absent.Present || absent.Rows == nil || len(absent.Rows) != 0 {
		t.Fatalf("story 9 = %+v, %v; want present false and an empty row array", absent, err)
	}
	if err := json.Unmarshal(absent.Encode(), &payload); err != nil || payload["present"] != false || payload["story"] != nil || len(payload["rows"].([]any)) != 0 {
		t.Fatalf("absent payload = %v, %v; want present false, no story, and an empty row array", payload, err)
	}

	legacy := "## User stories\n\n1. One.\n\n### Acceptance coverage map\n\n| story | behavior | seam | why it catches the failure |\n| --- | --- | --- | --- |\n| 1 | A behavior | a seam | a reason |\n"
	if document, err := Describe([]byte(legacy)); err != nil {
		t.Fatal(err)
	} else if story, _ := document.Story(1); len(story.Rows) != 1 || story.Rows[0].ID != "" || story.Rows[0].Rationale != "a reason" {
		t.Fatalf("legacy story = %+v, want one row without an id", story)
	}
	if document, err := Describe([]byte("## User stories\n\n1. One.\n")); err != nil || document.MapPresent || len(document.Stories) != 1 {
		t.Fatalf("no-map document = %+v, %v; want the stories and no map", document, err)
	}

	repeated := strings.Replace(storySpec, "3. As an agent", "1. As an agent", 1)
	repeated = strings.Replace(repeated, "| ST4 | 3 |", "| ST4 | 1 |", 1)
	if violations := Check(parse([]byte(repeated))); len(violations) != 0 {
		t.Fatalf("Check of a repeated story number = %v, want the existing deduplicated pass", violations)
	}
	document, err = Describe([]byte(repeated))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := document.Story(1); !errors.Is(err, ErrRefused) {
		t.Fatalf("repeated story 1 = %v, want an ambiguity refusal", err)
	}
	if story, err := document.Story(2); err != nil || !story.Present {
		t.Fatalf("unrepeated story 2 = %+v, %v; want present", story, err)
	}

	// A row whose story cell repeats a member joins the projection once.
	member := strings.Replace(storySpec, "| ST1 | 1 |", "| ST1 | 1, 1 |", 1)
	member = strings.Replace(member, "| ST3 | 1-2 |", "| ST3 | 1, 1-2 |", 1)
	document, err = Describe([]byte(member))
	if err != nil {
		t.Fatal(err)
	}
	if story, err := document.Story(1); err != nil || len(story.Rows) != 2 || story.Rows[0].ID != "ST1" || story.Rows[1].ID != "ST3" {
		t.Fatalf("story 1 with repeated members = %+v, %v; want ST1 and ST3 once each", story, err)
	}
	if violations := Check(parse([]byte(member))); len(violations) != 0 {
		t.Fatalf("Check of repeated members = %v, want the existing pass", violations)
	}

	// The projection carries text, so invalid UTF-8 in the story or in a cell of a matching row
	// refuses, and the other stories still project.
	for name, c := range map[string]struct {
		content string
		other   int
	}{
		"invalid story text": {strings.Replace(storySpec, "I want one story.", "I want one \xff story.", 1), 2},
		"invalid row cell":   {strings.Replace(storySpec, "The first behavior", "The first \xff behavior", 1), 3},
	} {
		document, err := Describe([]byte(c.content))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := document.Story(1); !errors.Is(err, ErrRefused) || !strings.Contains(err.Error(), "invalid UTF-8") {
			t.Errorf("%s: story 1 = %v, want an invalid UTF-8 refusal", name, err)
		}
		if story, err := document.Story(c.other); err != nil || !story.Present {
			t.Errorf("%s: story %d = %+v, %v; want present", name, c.other, story, err)
		}
	}

	for name, content := range map[string]string{
		"unknown story":    strings.Replace(storySpec, "| ST4 | 3 |", "| ST4 | 7 |", 1),
		"short row":        strings.Replace(storySpec, " | planned TestC in c_test.go |", " |", 1),
		"no header":        strings.Replace(storySpec, "| row | story |", "| rows | story |", 1),
		"historical fault": "<!-- coverage-map: historical -->\n" + strings.Replace(storySpec, "| ST4 | 3 |", "| ST4 | 7 |", 1),
		// A Markdown fault leaves the rest of the spec unread, so no partial projection holds.
		"fence before the map": strings.Replace(storySpec, "## Implementation decisions\n", "## Implementation decisions\n\n```text\n", 1),
		"fence inside the map": strings.Replace(storySpec, "the third failure. |\n\n", "the third failure. |\n\n```text\n", 1),
	} {
		if _, err := Describe([]byte(content)); !errors.Is(err, ErrRefused) {
			t.Errorf("%s: Describe = %v, want a validator refusal", name, err)
		}
	}
}
