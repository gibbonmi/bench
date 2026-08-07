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
