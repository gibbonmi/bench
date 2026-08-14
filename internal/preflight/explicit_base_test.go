package preflight

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/diff"
)

func TestExplicitBaseReviewOwnsSourceRangeNotDestinationHandoff(t *testing.T) {
	root, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	tip := runGit(t, "rev-parse", "feature")
	runGit(t, "checkout", "-q", "main")
	mustWriteFile(t, "capture/session-handoff.md", "destination only\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "destination handoff")
	runGit(t, "checkout", "-q", "feature")
	facts, boot := Gather(root, "review", slug, base)
	if boot != nil {
		t.Fatalf("Gather = %s: %s", boot.Kind, boot.Hint)
	}
	if strings.Join(facts.ChangedPaths, ",") != "internal/example/foo.go" {
		t.Fatalf("source paths = %v, want only source-authored path", facts.ChangedPaths)
	}
	configBefore, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	out, code := Command([]string{"review", slug, "--base", base})
	if code != 0 || !strings.Contains(out, "diff-nonempty,green") || !strings.Contains(out, "source[1]{base,tip}") || !strings.Contains(out, base) || !strings.Contains(out, tip) {
		t.Fatalf("explicit review = (%d):\n%s", code, out)
	}
	configAfter, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(configAfter) != string(configBefore) {
		t.Fatal("accepted explicit preflight changed Git config bytes")
	}
	facts, boot = Gather(root, "build", slug, base)
	if boot != nil {
		t.Fatalf("build Gather = %s: %s", boot.Kind, boot.Hint)
	}
	if strings.Join(facts.ChangedPaths, ",") != "internal/example/foo.go" {
		t.Fatalf("explicit source build paths = %v, want only source-authored path", facts.ChangedPaths)
	}
	// The advanced destination default branch is intentionally not an ancestor of
	// the retained source. An explicit source range remains valid because the
	// frozen base is an ancestor of its captured source tip.
	out, code = Command([]string{"build", slug, "--base", base})
	if code != 0 || !strings.Contains(out, "base-current,green") {
		t.Fatalf("explicit build = (%d):\n%s", code, out)
	}
	configAfter, err = os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(configAfter) != string(configBefore) {
		t.Fatal("accepted explicit build preflight changed Git config bytes")
	}

	// Mutation control: bare preflight still grades destination ancestry and must
	// refuse this diverged source rather than inheriting explicit-range validity.
	out, code = Command([]string{"build", slug})
	if code != 1 || !strings.Contains(out, "base-current,red,default branch tip is not an ancestor of HEAD") {
		t.Fatalf("bare build after destination advance = (%d):\n%s", code, out)
	}
	mustWriteFile(t, "internal/example/staged.go", "package example\n")
	runGit(t, "add", "internal/example/staged.go")
	mustWriteFile(t, "internal/example/foo.go", "package example\n// worktree\n")
	mustWriteFile(t, "internal/example/untracked.go", "package example\n")
	facts, boot = Gather(root, "build", slug, base)
	if boot != nil {
		t.Fatalf("build Gather = %s: %s", boot.Kind, boot.Hint)
	}
	for _, want := range []string{"internal/example/foo.go", "internal/example/staged.go", "internal/example/untracked.go"} {
		if !containsPath(facts.ChangedPaths, want) {
			t.Fatalf("complete source snapshot = %v, missing %s", facts.ChangedPaths, want)
		}
	}
}

func containsPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func TestExplicitBasePreflightRefusesDirtyReview(t *testing.T) {
	_, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	mustWriteFile(t, "dirty.txt", "dirty\n")
	out, code := Command([]string{"review", slug, "--base", base})
	if code != 1 || !strings.Contains(out, "source not clean") {
		t.Fatalf("dirty source = (%d):\n%s", code, out)
	}
}

func TestExplicitBasePreflightRetriesConvergedTrackedWorktreeMovement(t *testing.T) {
	root, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	tip := runGit(t, "rev-parse", "HEAD")
	configBefore, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	restore := diff.SetSnapshotAfterReadForTest(func() {
		calls++
		if calls == 1 {
			mustWriteFile(t, "internal/example/foo.go", "package example\n// moved\n")
		}
	})
	defer restore()

	out, code := Command([]string{"build", slug, "--base", base})
	if code != 0 || calls != 2 {
		t.Fatalf("converged tracked movement = (%d, %d):\n%s", code, calls, out)
	}
	if !strings.Contains(out, "source[1]{base,tip}") || !strings.Contains(out, base) || !strings.Contains(out, tip) {
		t.Fatalf("converged retry did not report one explicit source snapshot:\n%s", out)
	}
	configAfter, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(configAfter) != string(configBefore) {
		t.Fatal("converged explicit preflight changed Git config bytes")
	}
}

func TestExplicitBasePreflightMovementChecksWholeGather(t *testing.T) {
	t.Run("converges on post-gather fence movement", func(t *testing.T) {
		root, slug := seedConformant(t)
		base := runGit(t, "rev-parse", "main")
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			if calls == 1 {
				body := strings.Replace(specBody(slug), "`internal/"+slug+"/`", "`specs/"+slug+"/`", 1)
				mustWriteFile(t, "specs/"+slug+"/spec.md", body)
			}
		})
		defer restore()

		facts, failure := Gather(root, "build", slug, base)
		if failure != nil || calls != 2 {
			t.Fatalf("post-gather fence movement = (%#v, %d)", failure, calls)
		}
		if !containsPath(facts.ChangedPaths, "specs/"+slug+"/spec.md") || !containsPath(facts.FenceEntries, "specs/"+slug+"/") || containsPath(facts.FenceEntries, "internal/"+slug+"/") {
			t.Fatalf("converged gather facts mixed snapshots: %#v", facts)
		}
	})
	t.Run("does not retry past an ordinary gather failure", func(t *testing.T) {
		root, slug := seedConformant(t)
		base := runGit(t, "rev-parse", "main")
		mustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(specBody(slug), "Status: staged", "Status: draft", 1))
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
		})
		defer restore()

		_, failure := Gather(root, "build", slug, base)
		if failure == nil || failure.Kind != "spec not staged" || calls != 0 {
			t.Fatalf("ordinary gather failure = (%#v, %d), want spec-not-staged without retry", failure, calls)
		}
	})
}

func TestExplicitBasePreflightRefusesPersistentSnapshotMovement(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, call int)
	}{
		{
			name: "HEAD",
			mutate: func(t *testing.T, call int) {
				mustWriteFile(t, "internal/example/foo.go", "package example\n// head "+string(rune('0'+call))+"\n")
				runGit(t, "add", "internal/example/foo.go")
				runGit(t, "commit", "-q", "-m", "move head")
			},
		},
		{
			name: "index",
			mutate: func(t *testing.T, call int) {
				mustWriteFile(t, "internal/example/index.go", "package example\n// index "+string(rune('0'+call))+"\n")
				runGit(t, "add", "internal/example/index.go")
			},
		},
		{
			name: "tracked worktree",
			mutate: func(t *testing.T, call int) {
				mustWriteFile(t, "internal/example/foo.go", "package example\n// tracked "+string(rune('0'+call))+"\n")
			},
		},
		{
			name: "untracked state",
			mutate: func(t *testing.T, call int) {
				mustWriteFile(t, "internal/example/untracked.go", "package example\n// untracked "+string(rune('0'+call))+"\n")
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedConformant(t)
			base := runGit(t, "rev-parse", "main")
			configBefore, err := os.ReadFile(filepath.Join(root, ".git", "config"))
			if err != nil {
				t.Fatal(err)
			}
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				test.mutate(t, calls)
			})
			defer restore()

			out, code := Command([]string{"build", slug, "--base", base})
			if code != 1 || calls != 2 {
				t.Fatalf("persistent %s movement = (%d, %d):\n%s", test.name, code, calls, out)
			}
			if !strings.Contains(out, "error: snapshot drift") || !strings.Contains(out, "retry the exact invocation") {
				t.Fatalf("persistent %s movement did not return snapshot-drift refusal:\n%s", test.name, out)
			}
			wantHelp := "help[1]{cmd,why}:\n  bench preflight build " + slug + " --base " + base + "," + axi.RetryAfterMovement + "\n"
			if !strings.Contains(out, wantHelp) {
				t.Fatalf("persistent %s movement did not return the exact retry action:\n%s", test.name, out)
			}
			configAfter, err := os.ReadFile(filepath.Join(root, ".git", "config"))
			if err != nil {
				t.Fatal(err)
			}
			if string(configAfter) != string(configBefore) {
				t.Fatalf("persistent %s movement changed Git config bytes", test.name)
			}
		})
	}
}

func TestExplicitBasePreflightRetryActionPreservesArgumentOrder(t *testing.T) {
	_, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	calls := 0
	restore := diff.SetSnapshotAfterReadForTest(func() {
		calls++
		mustWriteFile(t, "internal/example/foo.go", "package example\n// moved "+string(rune('0'+calls))+"\n")
	})
	defer restore()

	out, code := Command([]string{"--base", base, "build", slug})
	if code != 1 || calls != 2 {
		t.Fatalf("noncanonical persistent movement = (%d, %d):\n%s", code, calls, out)
	}
	wantHelp := "help[1]{cmd,why}:\n  bench preflight --base " + base + " build " + slug + "," + axi.RetryAfterMovement + "\n"
	if !strings.Contains(out, wantHelp) {
		t.Fatalf("retry action did not preserve argument order:\n%s", out)
	}
}

func TestExplicitBasePreflightRefusesPersistentRawIndexByteMovement(t *testing.T) {
	root, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	baseline, code := diff.Command([]string{"--base", base, "--full"})
	if code != 0 {
		t.Fatalf("baseline full diff = (%d):\n%s", code, baseline)
	}
	configBefore, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	restore := diff.SetSnapshotAfterReadForTest(func() {
		calls++
		indexPath := filepath.Join(root, ".git", "index")
		before, err := os.ReadFile(indexPath)
		if err != nil {
			t.Fatal(err)
		}
		if len(before) < 8 {
			t.Fatalf("index is too short to carry a version: %d bytes", len(before))
		}
		target := "4"
		if binary.BigEndian.Uint32(before[4:8]) == 4 {
			target = "2"
		}
		runGit(t, "update-index", "--index-version="+target)
		after, err := os.ReadFile(indexPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) == string(after) {
			t.Fatal("index version mutation left raw index bytes unchanged")
		}
	})
	defer restore()

	out, code := Command([]string{"build", slug, "--base", base})
	if code != 1 || calls != 2 || !strings.Contains(out, "error: snapshot drift") {
		t.Fatalf("persistent raw index movement = (%d, %d):\n%s", code, calls, out)
	}
	restore()
	after, code := diff.Command([]string{"--base", base, "--full"})
	if code != 0 || after != baseline {
		t.Fatalf("raw index mutation changed live diff = (%d):\n%s", code, after)
	}
	configAfter, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(configAfter) != string(configBefore) {
		t.Fatal("persistent raw index movement changed Git config bytes")
	}
}

func TestSnapshotDriftRefusalKeepsPrimaryErrorWhenRetryActionCannotRender(t *testing.T) {
	for _, test := range []struct {
		name, slug, base string
	}{
		{name: "slug", slug: "example\x1b", base: "base"},
		{name: "base", slug: "example", base: "base\x1b"},
	} {
		t.Run(test.name, func(t *testing.T) {
			out := snapshotDriftRefusal([]string{"build", test.slug, "--base", test.base}, "the explicit source changed while reading; retry the exact invocation")
			if !strings.HasPrefix(out, "error: snapshot drift") || !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
				t.Fatalf("control-value snapshot drift refusal = %q", out)
			}
		})
	}
}

func TestExplicitBasePreflightCommandUsesResolvedRootFromSubdirectory(t *testing.T) {
	root, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	t.Chdir(filepath.Join(root, "internal", slug))
	out, code := Command([]string{"build", slug, "--base", base})
	if code != 0 || !strings.Contains(out, "source[1]{base,tip}") {
		t.Fatalf("subdirectory explicit preflight = (%d):\n%s", code, out)
	}
}

func TestExplicitBaseGatherAndAuthorizeUseSuppliedRootOutsideCWD(t *testing.T) {
	root, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	tip := runGit(t, "rev-parse", "HEAD")
	destination := initRepo(t)
	if root == destination {
		t.Fatal("source and destination fixture roots unexpectedly match")
	}

	facts, failure := Gather(root, "build", slug, base)
	if failure != nil {
		t.Fatalf("Gather(source root) = %s: %s", failure.Kind, failure.Hint)
	}
	if facts.SourceBase != base || facts.SourceTip != tip || !containsPath(facts.ChangedPaths, "internal/example/foo.go") {
		t.Fatalf("Gather source facts = %#v", facts)
	}
	source, err := AuthorizeReviewedSource(root, slug, base)
	if err != nil {
		t.Fatalf("AuthorizeReviewedSource(source root): %v", err)
	}
	if source.Base != base || source.Tip != tip || !containsPath(source.CommittedPaths, "internal/example/foo.go") {
		t.Fatalf("authorized source = %#v", source)
	}
}

func TestExplicitBasePreflightKeepsOrdinaryFailuresOutOfSnapshotDrift(t *testing.T) {
	t.Run("unreachable base", func(t *testing.T) {
		_, slug := seedConformant(t)
		out, code := Command([]string{"build", slug, "--base", "missing"})
		if code != 1 || !strings.HasPrefix(out, "error: cannot resolve --base") || strings.Contains(out, "snapshot drift") {
			t.Fatalf("unreachable explicit base = (%d):\n%s", code, out)
		}
	})
	t.Run("not in repository", func(t *testing.T) {
		t.Chdir(t.TempDir())
		out, code := Command([]string{"build", "example", "--base", "HEAD"})
		if code != 1 || !strings.HasPrefix(out, "error: not in a git repository") || strings.Contains(out, "snapshot drift") {
			t.Fatalf("not-in-repository explicit base = (%d):\n%s", code, out)
		}
	})
	t.Run("unrepresentable path", func(t *testing.T) {
		_, slug := seedConformant(t)
		base := runGit(t, "rev-parse", "main")
		hostile := "internal/" + slug + "/a\x1bb.go"
		mustWriteFile(t, hostile, "package example\n")
		runGit(t, "add", "--", hostile)
		runGit(t, "commit", "-q", "-m", "hostile path")
		out, code := Command([]string{"build", slug, "--base", base})
		if code != 1 || !strings.Contains(out, "unrepresentable TOON cell") || strings.Contains(out, "snapshot drift") {
			t.Fatalf("unrepresentable explicit path = (%d):\n%s", code, out)
		}
	})
}
