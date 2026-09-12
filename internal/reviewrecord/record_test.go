package reviewrecord

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReviewRecord(t *testing.T) {
	t.Run("missing clean axis", func(t *testing.T) {
		if err := CheckReviews(Chunk{ID: "1"}, []string{"author"}, false); err == nil || !strings.Contains(err.Error(), "missing Standards") {
			t.Fatalf("missing clean review result: %v", err)
		}
	})
}

func TestReviewRecordSchema(t *testing.T) {
	data, _ := json.Marshal(Record{Version: 99})
	if _, err := Parse(data); err == nil || !strings.Contains(err.Error(), "unsupported version") {
		t.Fatalf("unsupported version reached evidence validation: %v", err)
	}
}
