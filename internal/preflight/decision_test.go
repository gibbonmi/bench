package preflight

import "testing"

// baseFacts is a conformant Facts value: every check should answer green over it
// unmodified. Each test case starts from a copy and breaks exactly one field.
func baseFacts() Facts {
	return Facts{
		Mode:                  "review",
		SpecPath:              "specs/example/spec.md",
		DefaultBranchResolved: true,
		DefaultBranchCurrent:  true,
		ReviewBaseResolved:    true,
		ChangedPaths:          []string{"internal/example/foo.go"},
		FenceEntries:          []string{"internal/example"},
		DeclaredRowIDs:        []string{"PF1", "PF2"},
		TicketTokens:          []string{"PF1", "PF2"},
		SpecTag:               "PF",
	}
}

func checkRow(v Verdict, name string) (CheckResult, bool) {
	for _, c := range v.Checks {
		if c.Check == name {
			return c, true
		}
	}
	return CheckResult{}, false
}

// TestDecideAllGreen is the tracer: every check present by name, all green, and a
// second Decide over the same Facts is byte-identical — the rerun guarantee the
// verdict core promises.
func TestDecideAllGreen(t *testing.T) {
	f := baseFacts()
	first := Decide(f)
	second := Decide(f)

	wantChecks := []string{"base-current", "paths-authorized", "rows-owned", "rows-membership", "diff-nonempty"}
	if len(first.Checks) != len(wantChecks) {
		t.Fatalf("Checks = %#v, want %d rows", first.Checks, len(wantChecks))
	}
	for i, name := range wantChecks {
		if first.Checks[i].Check != name {
			t.Errorf("Checks[%d].Check = %q, want %q", i, first.Checks[i].Check, name)
		}
		if first.Checks[i].Verdict != verdictGreen {
			t.Errorf("Checks[%d] (%s) verdict = %q, want green: %+v", i, name, first.Checks[i].Verdict, first.Checks[i])
		}
	}
	if first.Red {
		t.Errorf("Verdict.Red = true, want false for an all-green Facts value")
	}
	if !decisionsEqual(first, second) {
		t.Errorf("Decide is not deterministic: first=%#v second=%#v", first, second)
	}
}

func decisionsEqual(a, b Verdict) bool {
	if a.Red != b.Red || len(a.Checks) != len(b.Checks) {
		return false
	}
	for i := range a.Checks {
		if a.Checks[i] != b.Checks[i] {
			return false
		}
	}
	return true
}

func TestDecideBaseCurrent(t *testing.T) {
	f := baseFacts()
	f.DefaultBranchResolved = false
	v := Decide(f)
	c, ok := checkRow(v, "base-current")
	if !ok || c.Verdict != verdictRed {
		t.Fatalf("base-current with unresolved default branch = %+v, want red", c)
	}
	if !v.Red {
		t.Error("Verdict.Red = false, want true")
	}

	f = baseFacts()
	f.DefaultBranchCurrent = false
	c, _ = checkRow(Decide(f), "base-current")
	if c.Verdict != verdictRed {
		t.Fatalf("base-current with stale default branch = %+v, want red", c)
	}
}

func TestDecidePathsAuthorized(t *testing.T) {
	f := baseFacts()
	f.ChangedPaths = []string{"unfenced/path.go"}
	c, _ := checkRow(Decide(f), "paths-authorized")
	if c.Verdict != verdictRed || c.Detail == "" {
		t.Fatalf("paths-authorized with an out-of-fence path = %+v, want red naming the path", c)
	}

	// A path equal to a fence entry, or under a fence prefix with a "/" separator,
	// stays green; a same-string-prefix path that lacks the separator does not.
	f = baseFacts()
	f.FenceEntries = []string{"internal/git"}
	f.ChangedPaths = []string{"internal/git", "internal/git/sub.go"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("paths-authorized with exact and prefix matches = %+v, want green", c)
	}
	f.ChangedPaths = []string{"internal/git2"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("paths-authorized: internal/git2 must not match fence internal/git = %+v, want red", c)
	}

	// An empty changed set is green.
	f = baseFacts()
	f.ChangedPaths = nil
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("paths-authorized with no changed paths = %+v, want green", c)
	}

	f = baseFacts()
	f.ReviewBaseResolved = false
	f.ReviewBaseHint = "no resolvable default branch"
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("paths-authorized with an unresolved review base = %+v, want red", c)
	}
}

func TestDecideRowsOwned(t *testing.T) {
	f := baseFacts()
	f.DeclaredRowIDs = []string{"PF1", "PF2", "PF3"}
	f.TicketTokens = []string{"PF1", "PF2"}
	c, _ := checkRow(Decide(f), "rows-owned")
	if c.Verdict != verdictRed || c.Detail == "" {
		t.Fatalf("rows-owned with one uncited declared row = %+v, want red naming PF3", c)
	}
}

func TestDecideRowsMembership(t *testing.T) {
	// A token under the spec's own tag naming no declared row is red.
	f := baseFacts()
	f.TicketTokens = []string{"PF1", "PF2", "PF99"}
	c, _ := checkRow(Decide(f), "rows-membership")
	if c.Verdict != verdictRed || c.Detail == "" {
		t.Fatalf("rows-membership with a phantom same-tag token = %+v, want red naming PF99", c)
	}

	// A foreign-tag token is ignored, not flagged.
	f = baseFacts()
	f.TicketTokens = []string{"PF1", "PF2", "FT93"}
	c, _ = checkRow(Decide(f), "rows-membership")
	if c.Verdict != verdictGreen {
		t.Fatalf("rows-membership with a foreign-tag token = %+v, want green (ignored)", c)
	}
}

func TestDecideDiffNonempty(t *testing.T) {
	f := baseFacts()
	f.ChangedPaths = nil
	c, _ := checkRow(Decide(f), "diff-nonempty")
	if c.Verdict != verdictRed {
		t.Fatalf("diff-nonempty with an empty changed set = %+v, want red", c)
	}

	f = baseFacts()
	f.ReviewBaseResolved = false
	c, _ = checkRow(Decide(f), "diff-nonempty")
	if c.Verdict != verdictRed {
		t.Fatalf("diff-nonempty with an unresolved review base = %+v, want red", c)
	}
}
