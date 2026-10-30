package retros

import (
	"reflect"
	"testing"
)

func TestRecommendationsKeepsImprovementParagraphsAndListItemsSeparate(t *testing.T) {
	content := []byte("Outside\n\n## Agent-experience improvements\n\nParagraph\ncontinued\n\n- First\n- Second\n\n### Details\n\nIgnored\n\n## Other\n\nIgnored too\n")
	want := []Recommendation{
		{Body: "Paragraph\ncontinued", Line: 5},
		{Body: "First", Line: 8},
		{Body: "Second", Line: 9},
		{Body: "Ignored", Line: 13},
	}
	if got := Recommendations(content); !reflect.DeepEqual(got, want) {
		t.Fatalf("recommendations = %#v, want %#v", got, want)
	}
}

// TestRecommendationsSkipFencedLines is MB46 and MB47. A fenced H2 line leaves the
// improvements open, so the items after it count. A fence marker line and a fenced item are
// no unit, and each marker line ends the open unit. A live H2 heading ends the unit and the
// section.
func TestRecommendationsSkipFencedLines(t *testing.T) {
	content := []byte("## Agent-experience improvements\n\n- First\n```\n## Notes\n- quoted\n```\nafter the fence\n\n- Last\n## Other\n- outside\n")
	want := []Recommendation{
		{Body: "First", Line: 3},
		{Body: "after the fence", Line: 8},
		{Body: "Last", Line: 10},
	}
	if got := Recommendations(content); !reflect.DeepEqual(got, want) {
		t.Fatalf("fenced lines: Recommendations gives %d units or other units, want the %d unfenced units", len(got), len(want))
	}
}

func TestRecommendationsIgnoreRepairAttributionHeadingBeforeImprovements(t *testing.T) {
	content := []byte("## Repair attribution\n\n| ticket | rounds | causes |\n|---|---|---|\n| add-thing | 1 | spec-row |\n\n## Agent-experience improvements\n\n- First\n- Second\n")
	want := []Recommendation{
		{Body: "First", Line: 9},
		{Body: "Second", Line: 10},
	}
	if got := Recommendations(content); !reflect.DeepEqual(got, want) {
		t.Fatalf("recommendations = %#v, want %#v", got, want)
	}
}

// TestImprovementMarkersIgnoreTheCostHeading grades a pending retro written before the
// Cost section and a retro whose cost lines are over budget. The retro grade reads only the
// destination markers, so neither raises a diagnostic.
func TestImprovementMarkersIgnoreTheCostHeading(t *testing.T) {
	improvements := "## Agent-experience improvements\n\n- Keep it.\n  Feeds: none\n"
	for name, body := range map[string]string{
		"pending without Cost": "## Outcome\n\nLanded.\n\n## Gate-stage timings\n\nunknown\n\n" + improvements,
		"over budget":          "## Outcome\n\nLanded.\n\n## Cost\n\n- wall per ticket: 65.0 min, over budget\n- bookkeeping per ticket: 7.2, over budget\n\n" + improvements,
	} {
		t.Run(name, func(t *testing.T) {
			if diags := ValidateImprovementMarkers(writeRetro(t, "landed.md", body)); len(diags) != 0 {
				t.Fatalf("retro grade = %q, want no diagnostic", diags)
			}
		})
	}
}

// TestFeedsMarkerGrammarAcceptsOnlyRowNewOrNone pins RF21: the destination value is a
// roadmap row ID, `new`, or `none`, and the marker is the unit's last line. The
// position half is written here rather than only in the driver test, because a check
// that searched the unit for the marker would satisfy a value-only expectation.
func TestFeedsMarkerGrammarAcceptsOnlyRowNewOrNone(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want bool
	}{
		{name: "a row ID", body: "Shorten the drain report.\nFeeds: FT12", want: true},
		{name: "a new row", body: "Shorten the drain report.\nFeeds: new", want: true},
		{name: "no row", body: "Shorten the drain report.\nFeeds: none", want: true},
		{name: "a multi-digit row ID", body: "Shorten the drain report.\nFeeds: FT172", want: true},
		{name: "an unknown value", body: "Shorten the drain report.\nFeeds: maybe", want: false},
		{name: "a zero row ID", body: "Shorten the drain report.\nFeeds: FT0", want: false},
		{name: "a lowercase row ID", body: "Shorten the drain report.\nFeeds: ft12", want: false},
		{name: "a trailing comment", body: "Shorten the drain report.\nFeeds: none for now", want: false},
		{name: "no marker at all", body: "Shorten the drain report.", want: false},
		{name: "a marker above the last line", body: "Feeds: FT12\nShorten the drain report.", want: false},
		{name: "a marker in the middle", body: "Shorten the drain report.\nFeeds: FT12\nIt is long.", want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := (Recommendation{Body: tc.body}).FeedsMarked(); got != tc.want {
				t.Fatalf("FeedsMarked(%q) = %v, want %v", tc.body, got, tc.want)
			}
		})
	}
}
