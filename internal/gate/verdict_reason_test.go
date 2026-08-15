package gate

import (
	"testing"
	"time"
)

func TestInspectVerdictClassReuseReason(t *testing.T) {
	for _, tc := range []struct {
		name       string
		makeRecord func(string, string, time.Time) []byte
		reusable   bool
		reason     string
	}{
		{name: "partial", makeRecord: inspectPartialRecord, reason: "partial verdict"},
		{name: "check-partial", makeRecord: inspectCheckPartialRecord, reason: "partial verdict"},
		{name: "combined-partial", makeRecord: inspectCombinedPartialRecord, reason: "partial verdict"},
		{name: "full", makeRecord: inspectFullRecord, reusable: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := outcomeFixture(t)
			now := time.Now().UTC().Truncate(time.Second)
			subject, err := buildSubject(root)
			if err != nil {
				t.Fatal(err)
			}
			writeInspectRecord(t, root, tc.makeRecord(subject.Tree, subject.Oracle, now))

			inspection := Inspect(root)
			if inspection.ReusableGreen != tc.reusable || inspection.Reason != tc.reason {
				t.Fatalf("inspection = %#v, want reusable=%t reason=%q", inspection, tc.reusable, tc.reason)
			}
		})
	}
}
