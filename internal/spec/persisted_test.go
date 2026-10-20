package spec

import (
	"reflect"
	"strings"
	"testing"
)

func TestPersistedFieldsHeader(t *testing.T) {
	for _, tc := range []struct {
		name, content string
		want          []PersistedField
	}{
		{"one line", "# S\n\nStatus: staged\n\nPersisted fields: intent-ledger:commitment.gaps, review-record:folds\n\nBody.\n",
			[]PersistedField{{Record: "intent-ledger", Path: []string{"commitment", "gaps"}}, {Record: "review-record", Path: []string{"folds"}}}},
		{"no line", "# S\n\nStatus: staged\n\nRoadmap: FT1\n", nil},
		{"a fenced line", "# S\n\n```markdown\nPersisted fields: review-record:folds\n```\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PersistedFields([]byte(tc.content))
			if err != nil || len(got) != len(tc.want) || (len(got) > 0 && !reflect.DeepEqual(got, tc.want)) {
				t.Fatalf("PersistedFields = %#v, %v; want %#v, nil", got, err, tc.want)
			}
		})
	}

	for _, value := range []string{
		"intent-ledger:a,,",
		"gate-evidence:x",
		"Intent-Ledger:x",
		"intent-ledger:commitment..gaps",
		"intent-ledger:",
		"intent-ledger:.gaps",
		"intent-ledger:gaps.",
		"intent-ledger:commitment. gaps",
		"intent-ledger:a:b",
		"intent-ledger:a b",
		"intent-ledger:a, review-record:folds, intent-ledger:a",
		"review-record:folds\n\nBody.\n\nPersisted fields: intent-ledger:commitment.gaps",
	} {
		t.Run("refuses "+value, func(t *testing.T) {
			_, err := PersistedFields([]byte("# S\n\nPersisted fields: " + value + "\n"))
			if err == nil || !strings.HasPrefix(err.Error(), "persisted-fields-invalid") {
				t.Fatalf("PersistedFields(%q) error = %v; want a persisted-fields-invalid error", value, err)
			}
		})
	}
}
