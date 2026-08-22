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
