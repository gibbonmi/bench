package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

func seedReviewEvidence(t *testing.T, poisonedConsumer bool) (root, slug string, args []string) {
	t.Helper()
	return preflighttest.SeedReviewEvidence(t, poisonedConsumer)
}

// countingObserver records how often each collector starts across one or more preparations.
func countingObserver(t *testing.T) map[string]int {
	t.Helper()
	counts := map[string]int{}
	restore := setReviewEvidenceObserverForTest(func(kind string) { counts[kind]++ })
	t.Cleanup(restore)
	return counts
}

func assertCollected(t *testing.T, counts map[string]int, want int) {
	t.Helper()
	for _, kind := range []string{"diff", "consumers", "coverage"} {
		if counts[kind] != want {
			t.Errorf("%s collection count = %d, want %d", kind, counts[kind], want)
		}
	}
}

// TestEvidenceReviewCollectors is CE66, CE67, and CE68. One preparation attempt runs each
// collector exactly once. Every read of the published artifact serves stored bytes and
// reaches no collector. A second preparation runs the collectors again, so a new attempt
// never reuses a prior attempt's capture.
func TestEvidenceReviewCollectors(t *testing.T) {
	_, _, args := seedReviewEvidence(t, false)
	counts := countingObserver(t)

	out, code := Command(args)
	if code != 0 {
		t.Fatalf("review preparation exit = %d:\n%s", code, out)
	}
	assertCollected(t, counts, 1)

	identity := preparedIdentity(t, out)
	for i := 0; i < 3; i++ {
		if read, code := Command([]string{"evidence", identity}); code != 0 {
			t.Fatalf("read %d = (%d):\n%s", i, code, read)
		}
	}
	assertCollected(t, counts, 1)

	if out, code := Command(args); code != 0 {
		t.Fatalf("second review preparation = (%d):\n%s", code, out)
	}
	assertCollected(t, counts, 2)
}

// TestEvidenceReviewCollectorsRerunPerAttempt is CE68's movement half. A source movement
// discards the first attempt, and the retry runs every collector again rather than reusing
// the discarded attempt's capture.
func TestEvidenceReviewCollectorsRerunPerAttempt(t *testing.T) {
	_, _, args := seedReviewEvidence(t, false)
	counts := countingObserver(t)
	moved := false
	restore := diff.SetSnapshotAfterReadForTest(func() {
		if !moved {
			moved = true
			preflighttest.RunGit(t, "branch", "-f", "main", args[6])
		}
	})
	defer restore()
	out, code := Command(args)
	if code != 0 {
		t.Fatalf("review preparation after one movement = %d:\n%s", code, out)
	}
	assertCollected(t, counts, 2)
}

// preparedIdentity reads the one prepared row's evidence identifier.
func preparedIdentity(t *testing.T, output string) string {
	t.Helper()
	rows := preflighttest.TableRows(t, preflighttest.DecodeMap(t, output), "prepared")
	if len(rows) != 1 {
		t.Fatalf("prepared rows = %d:\n%s", len(rows), output)
	}
	identity, ok := rows[0]["evidence"].(string)
	if !ok || identity == "" {
		t.Fatalf("prepared evidence = %#v", rows[0]["evidence"])
	}
	return identity
}

// TestEvidenceReviewProvenance is CE73. The manifest commits each collector's declared
// producer version and its exact arguments. The assertions read the committed rows rather
// than infer them from the identity, because a capture that embeds the same version would
// move the identity on its own and hide an omitted provenance row.
func TestEvidenceReviewProvenance(t *testing.T) {
	root, _, args := seedReviewEvidence(t, false)
	out, code := CommandWithVersion("fixture-version")(args)
	if code != 0 {
		t.Fatalf("review preparation = (%d):\n%s", code, out)
	}
	manifest := preflighttest.PublishedPack(t, root, preparedIdentity(t, out)).Manifest()
	if len(manifest.Producers) == 0 {
		t.Fatal("the manifest declares no producer")
	}
	for _, producer := range manifest.Producers {
		if producer.Version != "fixture-version" {
			t.Errorf("producer %q version = %q, want the running executable version",
				producer.Name, producer.Version)
		}
		if producer.Cwd != root {
			t.Errorf("producer %q cwd = %q, want %q", producer.Name, producer.Cwd, root)
		}
	}
	// Every collector states the frozen pair it read, so a capture cannot be attributed to
	// an invocation that never ran.
	for _, kind := range []string{"diff", "consumers"} {
		if !producerArgumentsContain(t, manifest, kind, args[6]) {
			t.Errorf("the %s producer arguments omit the frozen source tip", kind)
		}
	}
}

// producerArgumentsContain reports whether the named generated source declares value among
// its committed producer arguments.
func producerArgumentsContain(t *testing.T, manifest chargeevidence.Manifest, role, value string) bool {
	t.Helper()
	id := ""
	for _, source := range manifest.Sources {
		if source.Role == role {
			id = source.ID
		}
	}
	if id == "" {
		t.Fatalf("the manifest declares no %s source", role)
	}
	for _, argument := range manifest.Arguments {
		if argument.Source == id && argument.Value == value {
			return true
		}
	}
	return false
}

// TestEvidenceReviewMetadata is CE72. The completion facts ride inside the required metadata
// source, and that source's digest is one of the manifest rows the identity commits. The
// assertion re-encodes the committed metadata with one changed completion cell and compares
// digests, because preparing twice would also move the source tip and hide an omitted row.
func TestEvidenceReviewMetadata(t *testing.T) {
	root, _, args := seedReviewEvidence(t, false)
	out, code := Command(args)
	if code != 0 {
		t.Fatalf("review preparation = (%d):\n%s", code, out)
	}
	pack := preflighttest.PublishedPack(t, root, preparedIdentity(t, out))
	metadata := pack.Metadata()
	if len(metadata.Completion) != 1 {
		t.Fatalf("completion rows = %d, want one", len(metadata.Completion))
	}

	declared := ""
	for _, source := range pack.Manifest().Sources {
		if source.Role == chargeevidence.RoleMetadata {
			declared = source.SHA256
		}
	}
	if declared == "" {
		t.Fatal("the manifest declares no metadata source")
	}
	committed, err := chargeevidence.EncodeMetadata(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if chargeevidence.Digest(committed) != declared {
		t.Fatal("the committed metadata does not reproduce the declared metadata digest")
	}

	changed := metadata
	changed.Completion = []chargeevidence.CompletionRow{{
		Record: metadata.Completion[0].Record, SourceDigest: metadata.Completion[0].SourceDigest,
		PlanDigest: metadata.Completion[0].PlanDigest, RecordState: "invalid", Detail: "changed",
	}}
	altered, err := chargeevidence.EncodeMetadata(changed)
	if err != nil {
		t.Fatal(err)
	}
	if chargeevidence.Digest(altered) == declared {
		t.Fatal("a changed completion row left the metadata source digest unchanged")
	}
}

// TestEvidenceReviewSourceIdentity is CE69's producer half. A changed canonical review
// source changes the prepared identity every axis then resolves.
func TestEvidenceReviewSourceIdentity(t *testing.T) {
	for _, source := range []string{
		chargesource.ReviewSkill, chargesource.ReviewPhase,
		chargesource.DelegateSkill, chargesource.DelegateProcedure,
	} {
		t.Run(source, func(t *testing.T) {
			_, _, args := seedReviewEvidence(t, false)
			before, code := Command(args)
			if code != 0 {
				t.Fatalf("initial preparation = %d:\n%s", code, before)
			}
			preflighttest.MustWriteFile(t, source, "# Current canonical source\n\n"+source+" changed.\n")
			preflighttest.RunGit(t, "add", source)
			preflighttest.RunGit(t, "commit", "-q", "-m", "change review source")
			args[6] = preflighttest.RunGit(t, "rev-parse", "HEAD")
			after, code := Command(args)
			if code != 0 {
				t.Fatalf("changed source preparation = (%d):\n%s", code, after)
			}
			if preparedIdentity(t, before) == preparedIdentity(t, after) {
				t.Fatalf("a changed %s left the evidence identity unchanged", source)
			}
		})
	}
}

// TestEvidenceReviewRefusals is CE71 and CE127. Every refusal class stops the preparation,
// and none of them publishes a handle, so a failed collector leaves no artifact to read.
func TestEvidenceReviewRefusals(t *testing.T) {
	t.Run("collector refusal", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		preflighttest.MustWriteFile(t, "outside/broken.go", "package outside\n\nfunc Broken( {\n")
		preflighttest.RunGit(t, "add", "outside/broken.go")
		preflighttest.RunGit(t, "commit", "-q", "-m", "ill typed source")
		args[6] = preflighttest.RunGit(t, "rev-parse", "HEAD")
		assertReviewRefusal(t, root, args, "consumers evidence failed")
	})

	t.Run("incomplete consumer projection", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, true)
		assertReviewRefusal(t, root, args, "consumer evidence is incomplete")
	})

	t.Run("dirty checkout", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		preflighttest.MustWriteFile(t, "outside/dirty.go", "package outside\n")
		assertReviewRefusal(t, root, args, "source checkout is dirty")
	})

	t.Run("mismatched pair", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		args[6] = args[4]
		assertReviewRefusal(t, root, args, "tip-current")
	})

	t.Run("retry discards first attempt", func(t *testing.T) {
		root, slug, args := seedReviewEvidence(t, false)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			if calls == 1 {
				preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(preflighttest.SpecBody(slug), "Status: staged", "Status: draft", 1))
			}
		})
		defer restore()
		out, code := Command(args)
		if code != 1 || calls != 2 || !strings.Contains(out, "spec not staged") || strings.Contains(out, "prepared[") {
			t.Fatalf("retry then failure = (%d, calls=%d):\n%s", code, calls, out)
		}
		preflighttest.AssertNothingPublished(t, root)
	})

	t.Run("final recapture failure", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		head := filepath.Join(root, ".git", "HEAD")
		moved := head + ".during-review"
		restore := diff.SetSnapshotAfterReadForTest(func() {
			if err := os.Rename(head, moved); err != nil {
				t.Fatal(err)
			}
		})
		defer restore()
		out, code := Command(args)
		// The repository is unresolvable while HEAD is aside, so the store check waits
		// until the rename is undone below.
		if code != 1 || !strings.Contains(out, "snapshot identity failed") || strings.Contains(out, "prepared[") {
			t.Fatalf("recapture refusal = (%d):\n%s", code, out)
		}
		if _, err := os.Stat(moved); err == nil {
			if err := os.Rename(moved, head); err != nil {
				t.Fatal(err)
			}
		}
		preflighttest.AssertNothingPublished(t, root)
	})

	t.Run("required source absent", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		if err := os.Remove(chargesource.ReviewPhase); err != nil {
			t.Fatal(err)
		}
		preflighttest.RunGit(t, "add", "-A")
		preflighttest.RunGit(t, "commit", "-q", "-m", "remove review source")
		args[6] = preflighttest.RunGit(t, "rev-parse", "HEAD")
		assertReviewRefusal(t, root, args, chargesource.ReviewPhase)
	})

	for _, kind := range []string{"empty", "live symlink", "dangling symlink"} {
		t.Run("required source "+kind, func(t *testing.T) {
			root, _, args := seedReviewEvidence(t, false)
			if err := os.Remove(chargesource.ReviewPhase); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "empty":
				preflighttest.MustWriteFile(t, chargesource.ReviewPhase, "")
			case "live symlink":
				if err := os.Symlink("../skills/bench-craft-review/SKILL.md", chargesource.ReviewPhase); err != nil {
					t.Fatal(err)
				}
			case "dangling symlink":
				if err := os.Symlink("missing-review-source", chargesource.ReviewPhase); err != nil {
					t.Fatal(err)
				}
			}
			preflighttest.RunGit(t, "add", "-A")
			preflighttest.RunGit(t, "commit", "-q", "-m", "replace review source")
			args[6] = preflighttest.RunGit(t, "rev-parse", "HEAD")
			assertReviewRefusal(t, root, args, chargesource.ReviewPhase)
		})
	}
}

// TestEvidenceReviewCollectorFailureLeavesNoHandle is CE127's ordering half. A collector
// that fails after an earlier collector succeeded publishes nothing, so the earlier
// capture never reaches a readable artifact.
func TestEvidenceReviewCollectorFailureLeavesNoHandle(t *testing.T) {
	root, _, args := seedReviewEvidence(t, true)
	started := countingObserver(t)
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, "consumer evidence is incomplete") {
		t.Fatalf("incomplete consumer preparation = (%d):\n%s", code, out)
	}
	if started["diff"] != 1 {
		t.Fatalf("diff collection count = %d, want the earlier collector to have run", started["diff"])
	}
	if started["coverage"] != 0 {
		t.Fatalf("coverage collection count = %d, want no collector after the failure", started["coverage"])
	}
	preflighttest.AssertNothingPublished(t, root)
}

func TestEvidenceReviewMovementDiscardsPayload(t *testing.T) {
	for _, movement := range []string{"head", "index", "required source"} {
		t.Run(movement, func(t *testing.T) {
			root, _, args := seedReviewEvidence(t, false)
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				suffix := string(rune('0' + calls))
				switch movement {
				case "head":
					preflighttest.MustWriteFile(t, "notes/head.txt", "movement "+suffix+"\n")
					preflighttest.RunGit(t, "add", "notes/head.txt")
					preflighttest.RunGit(t, "commit", "-q", "-m", "move head")
				case "index":
					preflighttest.MustWriteFile(t, "notes/index.txt", "movement "+suffix+"\n")
					preflighttest.RunGit(t, "add", "notes/index.txt")
				case "required source":
					preflighttest.MustWriteFile(t, chargesource.ReviewPhase, "# Review phase\n\nmovement "+suffix+"\n")
				}
			})
			defer restore()
			out, code := Command(args)
			if code != 1 || calls != 2 || !strings.Contains(out, "snapshot drift") || strings.Contains(out, "prepared[") {
				t.Fatalf("persistent %s movement = (%d, calls=%d):\n%s", movement, code, calls, out)
			}
			preflighttest.AssertNothingPublished(t, root)
		})
	}
}

func assertReviewRefusal(t *testing.T, root string, args []string, want string) {
	t.Helper()
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, want) || strings.Contains(out, "prepared[") {
		t.Fatalf("review refusal = (%d), want %q without a published handle:\n%s", code, want, out)
	}
	preflighttest.AssertNothingPublished(t, root)
}
