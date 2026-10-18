package staleness

import (
	"bytes"
	"go/build"
	"io"
	"path"
	"slices"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
)

// redTests are the test discrepancy classes that a supported collector grades. Any other
// discrepancy stays diagnostic.
var redTests = []TestDiscrepancy{TestFileAbsent, TestFunctionAbsent, TestAlreadyShipped}

// tests audits each test seam against the committed object of its file in paths at the
// baseline and at the tip, through the coverage citation grammar. An existing test must be
// declared at both identities. A planned test must be declared at neither, and its file may
// hold other tests. Only the explicit initial audit reads this rule; an ordinary preflight
// never collects it. A red class at either identity wins over a diagnostic one. Each
// diagnostic state also answers one diagnostic, and an object that does not read refuses the
// collection.
func tests(root string, paths []PathState, seams []ownership.Seam) ([]TestState, []Diagnostic, error) {
	var states []TestState
	var diagnostics []Diagnostic
	for _, seam := range seams {
		state := TestState{Seam: seam.ID, Path: seam.Path, Function: seam.Function, Planned: seam.State == "planned"}
		if at := slices.IndexFunc(paths, func(p PathState) bool { return p.Path == seam.Path }); at >= 0 {
			state.Before, state.After = paths[at].Before, paths[at].After
		}
		for _, side := range []struct {
			at     string
			object Object
		}{{"source base", state.Before}, {"source tip", state.After}} {
			discrepancy, reason, err := testAt(root, state, side.object)
			if err != nil {
				return nil, nil, err
			}
			if discrepancy != "" && (state.Discrepancy == "" || slices.Contains(redTests, discrepancy) && !slices.Contains(redTests, state.Discrepancy)) {
				state.Discrepancy, state.Reason = discrepancy, reason+" at the "+side.at
			}
		}
		if state.Discrepancy != "" && !slices.Contains(redTests, state.Discrepancy) {
			diagnostics = append(diagnostics, Diagnostic{Subject: state.Seam, Reason: "cited test " + state.Function + " is " + string(state.Discrepancy) + ": " + state.Reason})
		}
		states = append(states, state)
	}
	return states, diagnostics, nil
}

// testAt is the discrepancy of one test seam at one identity, whose file object is object.
// The zero object holds no file.
func testAt(root string, state TestState, object Object) (TestDiscrepancy, string, error) {
	var src []byte
	if object.ID != "" {
		var err error
		if src, err = git.ReadControlBlob(root, object.ID); err != nil {
			return "", "", &Error{Stage: StageTest, Subject: state.Path, Cause: err}
		}
	}
	cited, _ := coverage.ResolveTest(src, state.Function)
	switch {
	case cited == coverage.TestUnjudged:
		return TestUnsupported, "the citation grammar cannot judge " + state.Function + " in " + state.Path, nil
	case state.Planned && cited == coverage.TestDeclared:
		return TestAlreadyShipped, state.Path + " already declares " + state.Function, nil
	case state.Planned:
		return "", "", nil
	case object.ID == "":
		return TestFileAbsent, state.Path + " is absent", nil
	case cited != coverage.TestDeclared:
		return TestFunctionAbsent, state.Path + " does not declare " + state.Function, nil
	case !defaultContext(state.Path, src):
		return TestContext, state.Path + " is outside the default build context", nil
	}
	return "", "", nil
}

// defaultContext reports whether the toolchain's default build context builds the file at
// p with source src: its build constraint and its filename suffix both match.
func defaultContext(p string, src []byte) bool {
	ctxt := build.Default
	ctxt.OpenFile = func(string) (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(src)), nil }
	dir, name := path.Split(p)
	matched, err := ctxt.MatchFile(dir, name)
	return err == nil && matched
}

// testEntry is the drift entry of one red test seam: its path and file object at each
// identity that holds the file, and a kind that names the test and its discrepancy. Absent
// fills a side that holds no file.
func testEntry(state TestState) DriftEntry {
	e := DriftEntry{Before: Absent, After: Absent, BeforeObject: Absent, AfterObject: Absent, Kind: "test " + state.Function + " " + string(state.Discrepancy)}
	if state.Before.ID != "" {
		e.Before, e.BeforeObject = state.Path, state.Before.ID
	}
	if state.After.ID != "" {
		e.After, e.AfterObject = state.Path, state.After.ID
	}
	return e
}
