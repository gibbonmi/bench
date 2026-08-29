package consumers

import (
	"reflect"
	"testing"
)

// unifiedDiff is one `git diff -U0` text carrying the three shapes the parse must read: a
// file with two hunks, a pure addition with no base side, and a deletion with no tip side.
const unifiedDiff = "diff --git a/pkg/edit.go b/pkg/edit.go\n" +
	"index 1111111..2222222 100644\n" +
	"--- a/pkg/edit.go\n" +
	"+++ b/pkg/edit.go\n" +
	"@@ -4 +4 @@ func Edit() {\n" +
	"-\told := 1\n" +
	"+\tnew := 2\n" +
	"@@ -10,0 +11,3 @@ func Tail() {\n" +
	"+\ta := 1\n" +
	"+\tb := 2\n" +
	"+\tc := 3\n" +
	"diff --git a/pkg/added.go b/pkg/added.go\n" +
	"new file mode 100644\n" +
	"--- /dev/null\n" +
	"+++ b/pkg/added.go\n" +
	"@@ -0,0 +1,2 @@\n" +
	"+package pkg\n" +
	"+\n" +
	"diff --git a/pkg/dropped.go b/pkg/dropped.go\n" +
	"deleted file mode 100644\n" +
	"--- a/pkg/dropped.go\n" +
	"+++ /dev/null\n" +
	"@@ -1,3 +0,0 @@\n" +
	"-package pkg\n" +
	"-\n" +
	"-func Dropped() {}\n"

// The hunk parse is the blast mode's whole reading of git, so it is graded on literal
// diff text rather than through a repository.
func TestParseHunksReadsAddedAndRemovedRuns(t *testing.T) {
	want := []fileHunks{
		{
			BasePath: "pkg/edit.go", TipPath: "pkg/edit.go",
			Added:   []lineSpan{{4, 4}, {11, 13}},
			Removed: []lineSpan{{4, 4}},
		},
		{BasePath: "", TipPath: "pkg/added.go", Added: []lineSpan{{1, 2}}},
		{BasePath: "pkg/dropped.go", TipPath: "", Removed: []lineSpan{{1, 3}}},
	}
	got := parseHunks(unifiedDiff)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseHunks = %#v, want %#v", got, want)
	}
}
