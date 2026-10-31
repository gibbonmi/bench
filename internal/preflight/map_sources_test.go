package preflight

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// The spec path of the conformant fixture, and the two reasons of the published-tree
// resolver as the row renders them.
const (
	mapsSpec    = "specs/example/spec.md"
	closedRow   = " is deleted by the delivery closure of specs/example/spec.md"
	missingFile = " must name an existing regular file"
)

// seedMaps builds the conformant fixture, runs setup on main with the repository root, and
// restarts the feature branch at main. It also commits docs/kept.md, which every map may
// cite as a kept source, and returns the slug.
func seedMaps(t *testing.T, setup func(t *testing.T, root string)) string {
	t.Helper()
	root, slug := preflighttest.SeedConformant(t)
	preflighttest.RestartFeature(t, slug, func() {
		preflighttest.MustWriteFile(t, "docs/kept.md", "kept\n")
		setup(t, root)
	})
	return slug
}

// seedClosureMaps binds the spec as the delivery of FT1 and writes, at each map path, a map
// that cites the FT1 row file the closure deletes and the FT2 row file it keeps.
func seedClosureMaps(t *testing.T) string {
	t.Helper()
	return seedMaps(t, func(t *testing.T, root string) {
		commitmenttest.SeedClosure(t, root, mapsSpec)
		for _, path := range []string{"specs/example/decisions/compiled.md", "decisions/active.md", "specs/other/decisions/other.md"} {
			preflighttest.WriteDecisionMap(t, path, preflighttest.PathRecord("roadmap/FT1.md"), preflighttest.PathRecord("roadmap/FT2.md"))
		}
	})
}

// mapSourcesDetail is pinnedMapSourcesDetail of the bare build preflight of slug.
func mapSourcesDetail(t *testing.T, slug string) string {
	t.Helper()
	return pinnedMapSourcesDetail(t, slug, "", "")
}

// pinnedMapSourcesDetail runs the build preflight of slug with an optional explicit base and
// source-tip pin, and returns the detail of its map-sources row, empty when that row is
// green. The render prints only red rows, so the verdict of the same tree decides whether
// the row exists, and a verdict without a map-sources row fails the test.
func pinnedMapSourcesDetail(t *testing.T, slug, base, pin string) string {
	t.Helper()
	verdict, decided := pinnedVerdicts(t, modeBuild, slug, base, pin)["map-sources"]
	if !decided {
		t.Fatal("the build verdict holds no map-sources row")
	}
	args := []string{"build", slug}
	if base != "" {
		args = append(args, "--base", base)
	}
	if pin != "" {
		args = append(args, "--source-tip", pin)
	}
	out, _ := Command(args)
	summaryCounts(t, out)
	if verdict == verdictGreen {
		return ""
	}
	if _, rendered := rowOf(t, out, "map-sources"); rendered {
		for _, row := range preflighttest.ResponseTable(t, out, "checks") {
			if row["check"] == "map-sources" {
				return row["detail"].(string)
			}
		}
	}
	t.Fatalf("map-sources verdict is %s, but the render holds no map-sources row:\n%s", verdict, out)
	return ""
}

// TestMapSourcesClosureDeletion checks that a bare build preflight reds the source that the
// predicted delivery closure deletes in an active map, in the spec's compiled map, and in
// another spec's compiled map, in candidate path order, and keeps the row file of another
// outcome.
func TestMapSourcesClosureDeletion(t *testing.T) {
	slug := seedClosureMaps(t)
	want := "decisions/active.md: Sources Path roadmap/FT1.md" + closedRow +
		"; specs/example/decisions/compiled.md: Sources Path roadmap/FT1.md" + closedRow +
		"; specs/other/decisions/other.md: Sources Path roadmap/FT1.md" + closedRow
	if got := mapSourcesDetail(t, slug); got != want {
		t.Fatalf("map-sources detail = %q, want %q", got, want)
	}
}

// TestMapSourcesTipExistence checks that the row grades the map and the sources that the
// source tip holds: a source that only the checkout holds and a map that only the checkout
// repoints stay red, a kept source is green in its cleaned spellings and as an executable
// file, and an absent map folder or a skipped name is green.
func TestMapSourcesTipExistence(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, root string)
		after func(t *testing.T)
		want  string
	}{
		{name: "missing at the tip", want: "decisions/m.md: Sources Path docs/never.md" + missingFile,
			setup: func(t *testing.T, _ string) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/never.md"))
			}},
		{name: "kept source",
			setup: func(t *testing.T, _ string) {
				preflighttest.MustWriteFile(t, "docs/run.sh", "#!/bin/sh\n")
				if err := os.Chmod("docs/run.sh", 0o755); err != nil {
					t.Fatal(err)
				}
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/kept.md"),
					preflighttest.PathRecord("./docs/kept.md"), preflighttest.PathRecord("docs/../docs/kept.md"),
					preflighttest.PathRecord("docs/run.sh"))
			},
			after: func(t *testing.T) {
				if mode := preflighttest.RunGit(t, "ls-tree", "--format=%(objectmode)", "HEAD", "--", "docs/run.sh"); mode != "100755" {
					t.Fatalf("docs/run.sh tree mode = %q, want 100755", mode)
				}
			}},
		{name: "checkout only", want: "decisions/m.md: Sources Path docs/local.md" + missingFile,
			setup: func(t *testing.T, _ string) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/local.md"))
			},
			after: func(t *testing.T) { preflighttest.MustWriteFile(t, "docs/local.md", "local\n") }},
		{name: "uncommitted repoint", want: "decisions/m.md: Sources Path docs/never.md" + missingFile,
			setup: func(t *testing.T, _ string) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/never.md"))
			},
			after: func(t *testing.T) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/kept.md"))
			}},
		{name: "absent folders", setup: func(t *testing.T, _ string) {}},
		{name: "skipped names",
			setup: func(t *testing.T, _ string) {
				for _, path := range []string{"decisions/README.md", "decisions/.hidden.md", "decisions/m/assets/nested.md",
					"specs/.hidden/decisions/m.md", "specs/example/decisions/notes.txt", "specs/example/notes/m.md"} {
					preflighttest.WriteDecisionMap(t, path, preflighttest.PathRecord("docs/never.md"))
				}
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			slug := seedMaps(t, tc.setup)
			if tc.after != nil {
				tc.after(t)
			}
			if got := mapSourcesDetail(t, slug); got != tc.want {
				t.Fatalf("map-sources detail = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestMapSourcesGrammar checks that the row grades each Sources record with the shared
// grammar: a record without Drift, an unknown opener, an escaping path, and a path that
// names the repository root are red, and a well-formed URL record is green.
func TestMapSourcesGrammar(t *testing.T) {
	slug := seedMaps(t, func(t *testing.T, _ string) {
		// The root then lists a regular file, which the root path itself is not.
		preflighttest.MustWriteFile(t, "top.md", "top\n")
		preflighttest.WriteDecisionMap(t, "decisions/m.md",
			"- Path: `docs/kept.md`\n  Supports: the decision.\n",
			"- File: docs/kept.md\n  Supports: the decision.\n  Drift: a change.\n",
			"- URL: https://example.invalid/source\n  Supports: the decision.\n  Drift: a change.\n",
			preflighttest.PathRecord("../outside.md"),
			preflighttest.PathRecord("."),
		)
	})
	want := "decisions/m.md: Sources docs/kept.md missing Drift" +
		"; decisions/m.md: Sources entry must start with Path or URL" +
		"; decisions/m.md: Sources Path ../outside.md escapes the repository" +
		"; decisions/m.md: Sources Path ." + missingFile
	if got := mapSourcesDetail(t, slug); got != want {
		t.Fatalf("map-sources detail = %q, want %q", got, want)
	}
}

// TestMapSourcesPinnedTip checks that a build preflight with an explicit base and a
// source-tip pin at an older commit grades the maps of the derived source tip, not those of
// the pin: the pin cites one missing source, the derived tip another.
func TestMapSourcesPinnedTip(t *testing.T) {
	slug := seedMaps(t, func(t *testing.T, _ string) {
		preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/pinned.md"))
	})
	pin := preflighttest.RunGit(t, "rev-parse", "HEAD")
	preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/derived.md"))
	preflighttest.CommitAll(t, "repoint the map")
	want := "decisions/m.md: Sources Path docs/derived.md" + missingFile
	if got := pinnedMapSourcesDetail(t, slug, "main", pin); got != want {
		t.Fatalf("map-sources detail = %q, want %q", got, want)
	}
}

// TestMapSourcesPrediction checks the closure prediction: a spec that no approved binding
// names deletes no row file, and an obligation-free binding and a policy that does not parse
// each red the row while the maps are still graded.
func TestMapSourcesPrediction(t *testing.T) {
	const missing = "; decisions/m.md: Sources Path docs/never.md" + missingFile
	for _, tc := range []struct {
		name, refusal, want string
		policy              func(t *testing.T, root string)
	}{
		{name: "no binding names the spec",
			policy: func(t *testing.T, root string) {
				preflighttest.MustWriteFile(t, "specs/other/spec.md", "# Other\n\nStatus: staged\n")
				commitmenttest.SeedClosure(t, root, "specs/other/spec.md")
			}},
		{name: "obligation-free binding", refusal: "names no obligation", want: missing,
			policy: func(t *testing.T, root string) {
				commitmenttest.SeedClosure(t, root, mapsSpec)
				commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
					policy.Milestones[0].Outcomes[0].Deliverables[0].Obligations = nil
				})
			}},
		{name: "unparsed policy", refusal: "commitment policy", want: missing,
			policy: func(t *testing.T, root string) {
				commitmenttest.SeedClosure(t, root, mapsSpec)
				commitmenttest.Write(t, root, commitment.PolicyPath, "not a policy\n")
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			slug := seedMaps(t, func(t *testing.T, root string) {
				tc.policy(t, root)
				records := []string{preflighttest.PathRecord("roadmap/FT1.md")}
				if tc.want != "" {
					records = append(records, preflighttest.PathRecord("docs/never.md"))
				}
				preflighttest.WriteDecisionMap(t, "decisions/m.md", records...)
			})
			got := mapSourcesDetail(t, slug)
			if tc.refusal == "" {
				if got != "" {
					t.Fatalf("map-sources detail = %q, want green", got)
				}
				return
			}
			prediction, graded, _ := strings.Cut(got, "; ")
			if !strings.HasPrefix(prediction, "delivery closure of "+mapsSpec+" not predictable: ") ||
				!strings.Contains(prediction, tc.refusal) || "; "+graded != tc.want {
				t.Fatalf("map-sources detail = %q, want the %q prediction refusal and then %q", got, tc.refusal, tc.want)
			}
		})
	}
}

// TestMapSourcesBuildOnly checks that a review preflight over a closure-deletion tree
// renders no map-sources row, counts its own rows, and gathers no map fact.
func TestMapSourcesBuildOnly(t *testing.T) {
	slug := seedClosureMaps(t)
	out, _ := Command([]string{"review", slug})
	if strings.Contains(out, "map-sources") || summaryCounts(t, out).total() != bareRowCount(modeReview) {
		t.Fatalf("review renders a map-sources row, or not %d rows:\n%s", bareRowCount(modeReview), out)
	}
	root := preflighttest.RunGit(t, "rev-parse", "--show-toplevel")
	facts, failure := Gather(root, modeReview, slug)
	if failure != nil || facts.MapSourceDiagnostics != nil {
		t.Fatalf("review Gather = %v, map facts %q; want no map fact", failure, facts.MapSourceDiagnostics)
	}
}

// TestMapSourcesKeepsAuditManual checks that the audit form grades the map-sources row of
// its explicit pair, and it requires the manual staleness pass while that row is red.
func TestMapSourcesKeepsAuditManual(t *testing.T) {
	slug := seedClosureMaps(t)
	out, code := Command([]string{"build", slug, "--audit-spec", "--base", "main", "--source-tip", "HEAD"})
	if code != 1 || !strings.Contains(out, "map-sources,red") {
		t.Fatalf("audit = (%d), want exit 1 with a red map-sources row:\n%s", code, out)
	}
	if manual := preflighttest.StalenessRow(t, out)["manual_required"]; manual != true {
		t.Errorf("manual_required = %v, want true while map-sources is red:\n%s", manual, out)
	}
}
