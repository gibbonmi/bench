package consumers

import (
	"strings"
	"testing"
)

// declPkg is one fixture package whose single file declares source at its own path. The
// candidates fixtures need distinct import paths carrying the same declaration name, so
// they build a package from source rather than from the reference-count helper.
func declPkg(root, dir, name, source string) fixturePkg {
	return fixturePkg{path: "example.com/" + dir, files: map[string]string{
		root + "/" + dir + "/" + name + ".go": source}}
}

// ambiguousFixture plants one declaration name in two packages with distinct last path
// segments, so a bare query has two matches that a one-segment spelling separates.
func ambiguousFixture(root string) []fixturePkg {
	return []fixturePkg{
		declPkg(root, "alpha", "alpha", "package alpha\n\nfunc Symbol() {}\n"),
		declPkg(root, "beta", "beta", "package beta\n\ntype Symbol struct{}\n"),
	}
}

// sharedSegmentFixture plants the same declaration name in two packages that share their
// last path segment, so a one-segment spelling stays ambiguous and the re-query argument
// must carry more of the import path.
func sharedSegmentFixture(root string) []fixturePkg {
	return []fixturePkg{
		declPkg(root, "alpha/util", "util", "package util\n\nfunc Symbol() {}\n"),
		declPkg(root, "beta/util", "util", "package util\n\nfunc Symbol() {}\n"),
	}
}

// CS5 (story 4): a bare name with two matches is an answer, not a refusal. It emits the
// candidates table at exit 0 and no consumer rows at all.
func TestAmbiguousBareNameEmitsCandidatesAndNoConsumerRows(t *testing.T) {
	stubLoad(t, ambiguousFixture)
	out, code := run(t, "Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "consumers_candidates[2]{qualified,file,line,kind}:\n  alpha.Symbol,alpha/alpha.go,3,func\n  beta.Symbol,beta/beta.go,3,type\n") {
		t.Fatalf("stdout = %q, want a two-row candidates table", out)
	}
	if strings.Contains(out, "consumers[") {
		t.Fatalf("stdout = %q, want no consumer rows for an ambiguous name", out)
	}
	if !strings.Contains(out, "meta[1]{packages,files,matches,rows,truncated}:\n  2,0,2,0,false\n") {
		t.Fatalf("stdout = %q, want meta matches=2 rows=0", out)
	}
}

// CS17 (story 30): each candidates row carries a re-query spelling that resolves to one
// declaration. The shared-segment case is the plant a last-segment-only spelling fails:
// both rows would read `util.Symbol` and re-query straight back to the same two matches.
func TestCandidateQualifiedSpellingsReQueryToOneMatch(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fixture func(string) []fixturePkg
		want    []string
	}{
		{"distinct segments", ambiguousFixture, []string{"alpha.Symbol", "beta.Symbol"}},
		{"shared segment", sharedSegmentFixture, []string{"alpha/util.Symbol", "beta/util.Symbol"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stubLoad(t, tc.fixture)
			out, code := run(t, "Symbol")
			if code != 0 {
				t.Fatalf("exit = %d, want 0; out=%q", code, out)
			}
			got := candidateQualified(t, out)
			assertRows(t, "candidate qualified cells", got, tc.want)
			for _, qualified := range got {
				reOut, reCode := run(t, qualified)
				if reCode != 0 {
					t.Fatalf("re-query %q: exit = %d, out=%q", qualified, reCode, reOut)
				}
				if !strings.Contains(reOut, "consumers[") || strings.Contains(reOut, "consumers_candidates") {
					t.Fatalf("re-query %q = %q, want a single-symbol result", qualified, reOut)
				}
			}
		})
	}
}

// CS18 (story 30): the candidates envelope is one literal re-query per row, in table
// order, with every argument known.
func TestCandidatesEnvelopeOffersOneReQueryPerRow(t *testing.T) {
	stubLoad(t, ambiguousFixture)
	out, code := run(t, "Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	want := "help[2]{cmd,why}:\n" +
		"  bench consumers alpha.Symbol,re-query the qualified symbol\n" +
		"  bench consumers beta.Symbol,re-query the qualified symbol\n"
	if !strings.HasSuffix(out, want) {
		t.Fatalf("stdout = %q, want the envelope %q", out, want)
	}
}

// candidateQualified reads the qualified cell out of each candidates row, so a test
// asserts the re-query argument the response actually printed.
func candidateQualified(t *testing.T, out string) []string {
	t.Helper()
	lines := strings.Split(out, "\n")
	var got []string
	inBlock := false
	for _, line := range lines {
		if strings.HasPrefix(line, "consumers_candidates[") {
			inBlock = true
			continue
		}
		if !inBlock {
			continue
		}
		if !strings.HasPrefix(line, "  ") {
			break
		}
		got = append(got, strings.SplitN(strings.TrimPrefix(line, "  "), ",", 2)[0])
	}
	if len(got) == 0 {
		t.Fatalf("stdout = %q, want a candidates block", out)
	}
	return got
}
