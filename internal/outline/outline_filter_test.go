package outline

import (
	"strings"
	"testing"
)

// productionOutline and testOutline are the two filtered renderings of
// writeKindFixtureTree. Together their symbol rows are exactly newOutline's rows: the
// production form keeps the fixture row and pkg/widget.go, and the test form keeps the
// two _test.go files. The metadata counts only the files the filter keeps in scope.
const productionOutline = `outline[3]{file,line,kind,name}:
  pkg/testdata/golden.txt,"1",fixture,golden.txt
  pkg/widget.go,"2",type,Widget
  pkg/widget.go,"3",func,New
outline_meta[1]{tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated}:
  "2","2","0","3","3","0","false"
outline_skips[0]{file,reason}:
`

const testOutline = `outline[7]{file,line,kind,name}:
  pkg/doubles_test.go,"2",double,FakeClock
  pkg/doubles_test.go,"3",double,stubStore
  pkg/doubles_test.go,"4",double,MockRepo
  pkg/doubles_test.go,"5",double,spyWriter
  pkg/widget_test.go,"2",helper,newFakeClock
  pkg/widget_test.go,"3",func,newer
  pkg/widget_test.go,"4",func,TestWidget
outline_meta[1]{tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated}:
  "2","2","0","7","7","0","false"
outline_skips[0]{file,reason}:
`

// OF1: --production drops every symbol a _test.go file declares, in the repository-wide
// form, the path form, and the directory summary.
func TestProductionFilterDropsTestFileSymbols(t *testing.T) {
	root := outlineRepo(t)
	writeKindFixtureTree(t, root)
	assertFilteredForms(t, "--production", productionOutline, `  pkg,"3"`)
}

// OF2: --test keeps only the symbols a _test.go file declares, in the same three forms.
func TestTestFilterKeepsOnlyTestFileSymbols(t *testing.T) {
	root := outlineRepo(t)
	writeKindFixtureTree(t, root)
	assertFilteredForms(t, "--test", testOutline, `  pkg,"7"`)
}

// assertFilteredForms runs one filter flag through the three forms: --full must equal
// want, the path form scoped to pkg carries the same bytes because every fixture file
// sits under pkg, and the bare summary carries dirRow.
func assertFilteredForms(t *testing.T, flag, want, dirRow string) {
	t.Helper()
	full, code := Command([]string{"--full", flag})
	if code != 0 || full != want {
		t.Fatalf("--full %s code=%d bytes =\n%s\nwant\n%s", flag, code, full, want)
	}
	scoped, code := Command([]string{flag, "pkg"})
	if code != 0 || scoped != want {
		t.Fatalf("%s pkg code=%d bytes =\n%s\nwant\n%s", flag, code, scoped, want)
	}
	bare, code := Command([]string{flag})
	if code != 0 || !strings.HasPrefix(bare, "outline_dirs[1]{dir,symbols}:\n"+dirRow+"\n") {
		t.Fatalf("bare %s code=%d summary:\n%s", flag, code, headOf(bare))
	}
}

// OF3: the two filters name opposite scopes, so a call with both refuses at exit 2 with
// the reason. The cwd is outside any repository, so the refusal comes before discovery.
func TestProductionAndTestTogetherRefuse(t *testing.T) {
	t.Chdir(t.TempDir())
	const want = "usage: bench outline [path] [--full] [--production|--test] (--production and --test are mutually exclusive)\n"
	for _, args := range [][]string{{"--production", "--test"}, {"--test", "--full", "--production"}} {
		out, code := Command(args)
		if code != 2 || out != want {
			t.Fatalf("Command(%q) code=%d out=%q\nwant code 2 out %q", args, code, out, want)
		}
	}
}

// bareKindOutline is the bare summary of writeKindFixtureTree without a filter flag.
const bareKindOutline = `outline_dirs[1]{dir,symbols}:
  pkg,"10"
outline_meta[1]{tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated}:
  "4","4","0","10","10","0","false"
outline_skips[0]{file,reason}:
`

// OF4: without a filter flag the bare summary keeps its bytes. The --full bytes of the
// same tree are pinned by TestOldToNewFixturePairPinsTheKindDelta.
func TestBareFormKeepsItsBytesWithoutAFilter(t *testing.T) {
	root := outlineRepo(t)
	writeKindFixtureTree(t, root)
	if got, code := Command(nil); code != 0 || got != bareKindOutline {
		t.Fatalf("bare code=%d bytes =\n%s\nwant\n%s", code, got, bareKindOutline)
	}
}

// OF5: the help text names both filter flags.
func TestHelpNamesTheFilterFlags(t *testing.T) {
	out, code := Command([]string{"--help"})
	if code != 0 || !strings.HasPrefix(out, "usage: bench outline [path] [--full] [--production|--test]\n") {
		t.Fatalf("help code=%d does not name the filter flags:\n%s", code, out)
	}
}
