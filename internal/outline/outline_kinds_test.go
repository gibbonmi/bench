package outline

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// writeKindFixtureTree plants one tree carrying a helper, all four double prefixes in
// mixed case, a testdata fixture file, and ordinary symbols. Every kind row below reads
// this one tree, so the old-to-new byte pair and the per-kind rows cannot disagree.
func writeKindFixtureTree(t *testing.T, root string) {
	t.Helper()
	writeOutlineFile(t, root, "pkg/widget.go", "package pkg\ntype Widget struct{}\nfunc New() *Widget { return nil }\n")
	writeOutlineFile(t, root, "pkg/widget_test.go", "package pkg\nfunc newFakeClock() int { return 0 }\nfunc newer() int { return 1 }\nfunc TestWidget() {}\n")
	writeOutlineFile(t, root, "pkg/doubles_test.go", "package pkg\ntype FakeClock struct{}\ntype stubStore struct{}\nfunc MockRepo() {}\nfunc spyWriter() {}\n")
	writeOutlineFile(t, root, "pkg/testdata/golden.txt", "hello\n")
	gitAddOutline(t, root)
}

// OI1: a _test.go function named with a declared helper prefix emits kind helper. The
// path predicate and the upper-case requirement are both observable here: the same
// source in a non-test path, and the prefix-only name `newer`, stay kind func.
func TestHelperPrefixInTestFileEmitsHelperKind(t *testing.T) {
	const source = "package pkg\n" +
		"func newFakeClock() int { return 0 }\n" +
		"func makeStore() int { return 0 }\n" +
		"func withTimeout() int { return 0 }\n" +
		"func newer() int { return 0 }\n"
	want := []Symbol{
		{Line: 2, Kind: "helper", Name: "newFakeClock"},
		{Line: 3, Kind: "helper", Name: "makeStore"},
		{Line: 4, Kind: "helper", Name: "withTimeout"},
		{Line: 5, Kind: "func", Name: "newer"},
	}
	if got := Symbols("pkg/widget_test.go", []byte(source)); !reflect.DeepEqual(got, want) {
		t.Fatalf("test-file symbols = %#v\nwant %#v", got, want)
	}
	for _, s := range Symbols("pkg/widget.go", []byte(source)) {
		if s.Kind != "func" {
			t.Fatalf("non-test path emitted kind %q for %q; want func", s.Kind, s.Name)
		}
	}
}

// OI2: a name with the fake, stub, mock, or spy prefix emits kind double. All four
// prefixes are planted, two of them in mixed case, so a fake-only or case-sensitive
// classifier reds here.
func TestDoublePrefixEmitsDoubleKind(t *testing.T) {
	const source = "package pkg\n" +
		"type FakeClock struct{}\n" +
		"type stubStore struct{}\n" +
		"func MockRepo() {}\n" +
		"func spyWriter() {}\n" +
		"type Widget struct{}\n"
	want := []Symbol{
		{Line: 2, Kind: "double", Name: "FakeClock"},
		{Line: 3, Kind: "double", Name: "stubStore"},
		{Line: 4, Kind: "double", Name: "MockRepo"},
		{Line: 5, Kind: "double", Name: "spyWriter"},
		{Line: 6, Kind: "type", Name: "Widget"},
	}
	if got := Symbols("pkg/doubles_test.go", []byte(source)); !reflect.DeepEqual(got, want) {
		t.Fatalf("double symbols = %#v\nwant %#v", got, want)
	}
	// The double form carries no path predicate: any scanned file reports it.
	plain := Symbols("pkg/doubles.go", []byte("package pkg\ntype fakeClock struct{}\n"))
	if len(plain) != 1 || plain[0].Kind != "double" {
		t.Fatalf("non-test path double = %#v; want one double row", plain)
	}
}

// OI3: a file under a testdata/ segment emits one fixture row carrying line 1 and its
// base name, and is not scanned for symbols.
func TestTestdataFileEmitsOneFixtureRow(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "pkg/testdata/golden.txt", "hello\nfunc NotScanned() {}\n")
	writeOutlineFile(t, root, "pkg/testdata/nested/x.go", "package x\nfunc Skipped() {}\n")
	writeOutlineFile(t, root, "pkg/widget.go", "package pkg\nfunc Kept() {}\n")
	gitAddOutline(t, root)

	out, code := Command([]string{"--full"})
	if code != 0 {
		t.Fatalf("code = %d\n%s", code, headOf(out))
	}
	for _, want := range []string{`  pkg/testdata/golden.txt,"1",fixture,golden.txt`, `  pkg/testdata/nested/x.go,"1",fixture,x.go`, `  pkg/widget.go,"2",func,Kept`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, headOf(out))
		}
	}
	if strings.Contains(out, "Skipped") {
		t.Fatalf("a testdata file was scanned for symbols:\n%s", headOf(out))
	}
	if !strings.HasPrefix(out, "outline[3]{file,line,kind,name}:\n") {
		t.Fatalf("fixture files did not emit exactly one row each:\n%s", headOf(out))
	}
}

// OI3: a fixture row needs no content read, so a fixture the read policy would reject
// still gets its row. A NUL-carrying file and a file over bounds.OutlineFileLimit both
// plant here, and neither may reach the skip table.
func TestUnscannableTestdataFilesStillEmitFixtureRows(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "pkg/testdata/blob.bin", "head\x00tail\n")
	writeOutlineFile(t, root, "pkg/testdata/huge.txt", strings.Repeat("x", int(bounds.OutlineFileLimit)+1))
	writeOutlineFile(t, root, "pkg/widget.go", "package pkg\nfunc Kept() {}\n")
	gitAddOutline(t, root)

	out, code := Command([]string{"--full"})
	if code != 0 {
		t.Fatalf("code = %d\n%s", code, headOf(out))
	}
	for _, want := range []string{`  pkg/testdata/blob.bin,"1",fixture,blob.bin`, `  pkg/testdata/huge.txt,"1",fixture,huge.txt`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, headOf(out))
		}
	}
	if !strings.Contains(out, "outline_skips[0]{file,reason}:") {
		t.Fatalf("an unscannable fixture was skipped instead of listed:\n%s", headOf(out))
	}
}

// oldOutline and newOutline are the two pinned renderings of writeKindFixtureTree:
// oldOutline is the kind vocabulary without helper, double, and fixture, and newOutline is
// the one this package emits. The pair pins the byte delta: any change to an
// outline row reds this test.
const oldOutline = `outline[9]{file,line,kind,name}:
  pkg/doubles_test.go,"2",type,FakeClock
  pkg/doubles_test.go,"3",type,stubStore
  pkg/doubles_test.go,"4",func,MockRepo
  pkg/doubles_test.go,"5",func,spyWriter
  pkg/widget.go,"2",type,Widget
  pkg/widget.go,"3",func,New
  pkg/widget_test.go,"2",func,newFakeClock
  pkg/widget_test.go,"3",func,newer
  pkg/widget_test.go,"4",func,TestWidget
outline_meta[1]{tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated}:
  "4","4","0","9","9","0","false"
outline_skips[0]{file,reason}:
`

const newOutline = `outline[10]{file,line,kind,name}:
  pkg/doubles_test.go,"2",double,FakeClock
  pkg/doubles_test.go,"3",double,stubStore
  pkg/doubles_test.go,"4",double,MockRepo
  pkg/doubles_test.go,"5",double,spyWriter
  pkg/testdata/golden.txt,"1",fixture,golden.txt
  pkg/widget.go,"2",type,Widget
  pkg/widget.go,"3",func,New
  pkg/widget_test.go,"2",helper,newFakeClock
  pkg/widget_test.go,"3",func,newer
  pkg/widget_test.go,"4",func,TestWidget
outline_meta[1]{tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated}:
  "4","4","0","10","10","0","false"
outline_skips[0]{file,reason}:
`

// OI4: the old-to-new fixture pair reds on any outline byte delta. The
// emitted bytes must equal newOutline, and the old-to-new difference must be exactly
// the planted rows — no unplanted line may move.
func TestOldToNewFixturePairPinsTheKindDelta(t *testing.T) {
	root := outlineRepo(t)
	writeKindFixtureTree(t, root)
	got, code := Command([]string{"--full"})
	if code != 0 {
		t.Fatalf("code = %d\n%s", code, headOf(got))
	}
	if got != newOutline {
		t.Fatalf("outline bytes =\n%s\nwant\n%s", got, newOutline)
	}
	wantGone := []string{
		"outline[9]{file,line,kind,name}:",
		`  pkg/doubles_test.go,"2",type,FakeClock`,
		`  pkg/doubles_test.go,"3",type,stubStore`,
		`  pkg/doubles_test.go,"4",func,MockRepo`,
		`  pkg/doubles_test.go,"5",func,spyWriter`,
		`  pkg/widget_test.go,"2",func,newFakeClock`,
		`  "4","4","0","9","9","0","false"`,
	}
	wantNew := []string{
		"outline[10]{file,line,kind,name}:",
		`  pkg/doubles_test.go,"2",double,FakeClock`,
		`  pkg/doubles_test.go,"3",double,stubStore`,
		`  pkg/doubles_test.go,"4",double,MockRepo`,
		`  pkg/doubles_test.go,"5",double,spyWriter`,
		`  pkg/testdata/golden.txt,"1",fixture,golden.txt`,
		`  pkg/widget_test.go,"2",helper,newFakeClock`,
		`  "4","4","0","10","10","0","false"`,
	}
	if gone := linesOnlyIn(oldOutline, newOutline); !reflect.DeepEqual(gone, wantGone) {
		t.Fatalf("lines the delta removed = %#v\nwant %#v", gone, wantGone)
	}
	if added := linesOnlyIn(newOutline, oldOutline); !reflect.DeepEqual(added, wantNew) {
		t.Fatalf("lines the delta added = %#v\nwant %#v", added, wantNew)
	}
}

// linesOnlyIn returns a's lines, in order, that b does not carry.
func linesOnlyIn(a, b string) []string {
	have := map[string]bool{}
	for _, line := range strings.Split(b, "\n") {
		have[line] = true
	}
	var only []string
	for _, line := range strings.Split(a, "\n") {
		if !have[line] {
			only = append(only, line)
		}
	}
	return only
}

// OI5: outline help keeps the LOCATE promise line verbatim. A reworded promise would
// upgrade candidate rows to verified ones.
func TestHelpKeepsLocatePromiseVerbatim(t *testing.T) {
	const want = "outline locates candidate seams (file:line); it does not identify which are the project's blessed seams — projects/<name>.md owns that."
	if promise != want {
		t.Fatalf("promise = %q\nwant %q", promise, want)
	}
	out, code := Command([]string{"--help"})
	if code != 0 || !strings.Contains(out, want) {
		t.Fatalf("help code=%d lost the promise:\n%s", code, out)
	}
}
