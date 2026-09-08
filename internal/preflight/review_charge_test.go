package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	toonlib "github.com/toon-format/toon-go"
)

func seedReviewEvidence(t *testing.T, poisonedConsumer bool) (root, slug string, args []string) {
	t.Helper()
	slug = "example"
	root = initRepo(t)
	mustWriteFile(t, "go.mod", "module example.com/review\n\ngo 1.25\n")
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug,
		"- `target/` (review fixture)",
		"- `edited/` (review fixture)",
		"- `outside/` (review fixture)",
		"- `notes/` (review fixture)",
		"- `.agents/skills/bench-craft-review/` (review instructions)",
		"- `.agents/commands/bench-review-implementation.md` (review phase)",
	))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticketDoc("One", "PF1", "PF2"))
	mustWriteFile(t, delegateSkill, "# Delegation skill\n")
	mustWriteFile(t, delegateProcedure,
		"# Delegation procedure\n\nFocused suite: bench test --package ./internal/preflight\n")
	mustWriteFile(t, buildPhase, "# Build phase\n")
	mustWriteFile(t, reviewSkill,
		"# Review skill\n\n## Standards\n\nRules.\n\n## Spec\n\nRequirements.\n\n## Coverage\n\nEdges.\n")
	mustWriteFile(t, reviewPhase, "# Review phase\n\nUse the three canonical axes.\n")
	mustWriteFile(t, "target/target.go", "package target\n\nfunc Changed() int { return 0 }\nfunc Gone() {}\n")
	mustWriteFile(t, "outside/user.go",
		"package outside\n\nimport \"example.com/review/target\"\n\nfunc Use() int { return target.Changed() }\n")
	mustWriteFile(t, "edited/user.go",
		"package edited\n\nimport \"example.com/review/target\"\n\nfunc Use() int { return target.Changed() }\n")
	if poisonedConsumer {
		mustWriteFile(t, "outside/a\x1b.go", "package outside\n\nimport \"example.com/review/target\"\n\nfunc Poisoned() int { return target.Changed() }\n")
	}
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "base")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "target/target.go", "package target\n\nfunc Changed() int { return 1 }\n")
	mustWriteFile(t, "edited/user.go", "package edited\n\nimport \"example.com/review/target\"\n\n// Use is an edited consumer.\nfunc Use() int { return target.Changed() }\n")
	mustWriteFile(t, "notes/a \"quote\" \\ café*.txt", "review π evidence\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "review source")
	activeAssignment(t, root, root)
	args = []string{
		"review", slug, "--charge", "--base", runGit(t, "rev-parse", "main"),
		"--source-tip", runGit(t, "rev-parse", "HEAD"), "--full",
	}
	return root, slug, args
}

func TestReviewChargeSharedEvidence(t *testing.T) {
	_, _, args := seedReviewEvidence(t, false)
	counts := map[string]int{}
	restore := setReviewEvidenceObserverForTest(func(kind string) { counts[kind]++ })
	defer restore()

	out, code := Command(args)
	if code != 0 {
		t.Fatalf("review charge exit = %d:\n%s", code, out)
	}
	for _, kind := range []string{"diff", "consumers", "coverage"} {
		if counts[kind] != 1 {
			t.Errorf("%s collection count = %d, want 1", kind, counts[kind])
		}
	}
	if strings.Count(out, "shared_evidence[3]") != 1 {
		t.Fatalf("shared evidence was not emitted once:\n%s", out)
	}
	for _, axis := range []string{"Standards", "Spec", "Coverage"} {
		if !strings.Contains(out, axis) {
			t.Errorf("review charge omitted %s axis:\n%s", axis, out)
		}
	}

	_, _, args = seedReviewEvidence(t, false)
	counts = map[string]int{}
	restore = setReviewEvidenceObserverForTest(func(kind string) { counts[kind]++ })
	defer restore()
	moved := false
	restoreSnapshot := diff.SetSnapshotAfterReadForTest(func() {
		if !moved {
			moved = true
			runGit(t, "branch", "-f", "main", args[6])
		}
	})
	defer restoreSnapshot()
	out, code = Command(args)
	if code != 0 {
		t.Fatalf("review charge after one movement = %d:\n%s", code, out)
	}
	for _, kind := range []string{"diff", "consumers", "coverage"} {
		if counts[kind] != 2 {
			t.Errorf("%s collection count across two attempts = %d, want 2", kind, counts[kind])
		}
	}
}

func TestReviewChargeEvidence(t *testing.T) {
	_, _, args := seedReviewEvidence(t, false)
	compactArgs := append([]string(nil), args[:len(args)-1]...)
	compact, code := CommandWithVersion("fixture-version")(compactArgs)
	wantRetrieval := "bench preflight review specs/example/spec.md --charge --base " +
		args[4] + " --source-tip " + args[6] + " --full"
	if code != 0 || !strings.Contains(compact, "\"false\"") || !strings.Contains(compact, wantRetrieval) || strings.Contains(compact, "diff_body") {
		t.Fatalf("compact review charge = (%d), want exact retrieval %q:\n%s", code, wantRetrieval, compact)
	}
	out, code := CommandWithVersion("fixture-version")(args)
	if code != 0 {
		t.Fatalf("review charge exit = %d:\n%s", code, out)
	}
	for _, want := range []string{
		"diff_body", "PF1", "PF2", "catches z", "catches w", "review π evidence",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("review evidence omitted %q:\n%s", want, out)
		}
	}
	if strings.Count(out, "diff_body") != 1 {
		t.Errorf("diff evidence was duplicated across axes:\n%s", out)
	}
	consumerOutput := evidenceContent(t, out, "consumers")
	document := decodeMap(t, consumerOutput)
	assertBlastRow(t, document, "target.Changed", "outside/user.go", false)
	assertBlastRow(t, document, "target.Changed", "edited/user.go", true)
	deleted := tableRows(t, document, "blast_deleted")
	if len(deleted) != 1 || deleted[0]["changed_symbol"] != "target.Gone" {
		t.Fatalf("deleted rows = %#v, want target.Gone", deleted)
	}
	citation := tableRows(t, document, "citation")
	if len(citation) != 1 || citation[0]["version"] != "fixture-version" || citation[0]["hash"] == "" {
		t.Fatalf("citation = %#v, want injected version and hash", citation)
	}
}

func evidenceContent(t *testing.T, output, source string) string {
	t.Helper()
	for _, row := range tableRows(t, decodeMap(t, output), "evidence") {
		if row["source"] == source {
			content, ok := row["content"].(string)
			if !ok {
				t.Fatalf("evidence %s content = %T", source, row["content"])
			}
			return content
		}
	}
	t.Fatalf("no %s evidence in packet", source)
	return ""
}

func decodeMap(t *testing.T, output string) map[string]any {
	t.Helper()
	decoded, err := toonlib.DecodeString(output)
	if err != nil {
		t.Fatalf("decode TOON: %v\n%s", err, output)
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		t.Fatalf("decoded packet = %T, want object", decoded)
	}
	return document
}

func tableRows(t *testing.T, document map[string]any, table string) []map[string]any {
	t.Helper()
	values, ok := document[table].([]any)
	if !ok {
		t.Fatalf("%s = %T, want table", table, document[table])
	}
	rows := make([]map[string]any, len(values))
	for i, value := range values {
		row, ok := value.(map[string]any)
		if !ok {
			t.Fatalf("%s row %d = %T, want object", table, i, value)
		}
		rows[i] = row
	}
	return rows
}

func assertBlastRow(t *testing.T, document map[string]any, symbol, file string, touched bool) {
	t.Helper()
	want := map[string]any{"changed_symbol": symbol, "file": file, "touched": touched}
	for _, row := range tableRows(t, document, "blast") {
		if row["changed_symbol"] == symbol && row["file"] == file && row["touched"] == touched {
			return
		}
	}
	t.Fatalf("blast rows omitted %#v: %#v", want, tableRows(t, document, "blast"))
}

func TestReviewChargeAxes(t *testing.T) {
	for _, source := range []string{reviewSkill, reviewPhase, delegateSkill, delegateProcedure} {
		t.Run(source, func(t *testing.T) {
			_, _, args := seedReviewEvidence(t, false)
			before, code := Command(args)
			if code != 0 {
				t.Fatalf("initial charge exit = %d:\n%s", code, before)
			}
			changed := "# Current canonical source\n\n" + source + " changed.\n"
			mustWriteFile(t, source, changed)
			runGit(t, "add", source)
			runGit(t, "commit", "-q", "-m", "change review source")
			args[6] = runGit(t, "rev-parse", "HEAD")
			after, code := Command(args)
			identity := (chargeSource{path: source, data: []byte(changed)}).identity()
			if code != 0 || after == before || !strings.Contains(after, identity) ||
				!strings.Contains(after, source+" changed") {
				t.Fatalf("changed source charge = (%d):\n%s", code, after)
			}
		})
	}
}

func TestReviewChargeRefusals(t *testing.T) {
	t.Run("collector refusal", func(t *testing.T) {
		_, _, args := seedReviewEvidence(t, false)
		mustWriteFile(t, "outside/broken.go", "package outside\n\nfunc Broken( {\n")
		runGit(t, "add", "outside/broken.go")
		runGit(t, "commit", "-q", "-m", "ill typed source")
		args[6] = runGit(t, "rev-parse", "HEAD")
		assertReviewRefusal(t, args, "consumer evidence failed")
	})

	t.Run("incomplete consumer projection", func(t *testing.T) {
		_, _, args := seedReviewEvidence(t, true)
		assertReviewRefusal(t, args, "consumer evidence is incomplete")
	})

	t.Run("dirty checkout", func(t *testing.T) {
		_, _, args := seedReviewEvidence(t, false)
		mustWriteFile(t, "outside/dirty.go", "package outside\n")
		assertReviewRefusal(t, args, "source checkout is dirty")
	})

	t.Run("mismatched pair", func(t *testing.T) {
		_, _, args := seedReviewEvidence(t, false)
		args[6] = args[4]
		assertReviewRefusal(t, args, "tip-current")
	})

	t.Run("retry discards first attempt", func(t *testing.T) {
		_, slug, args := seedReviewEvidence(t, false)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			if calls == 1 {
				mustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(specBody(slug), "Status: staged", "Status: draft", 1))
			}
		})
		defer restore()
		out, code := Command(args)
		if code != 1 || calls != 2 || !strings.Contains(out, "spec not staged") ||
			strings.Contains(out, "shared_evidence") || strings.Contains(out, "charge[") {
			t.Fatalf("retry then failure = (%d, calls=%d):\n%s", code, calls, out)
		}
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
		defer func() {
			restore()
			if _, err := os.Stat(moved); err == nil {
				if err := os.Rename(moved, head); err != nil {
					t.Fatal(err)
				}
			}
		}()
		assertReviewRefusal(t, args, "snapshot identity failed")
	})

	t.Run("required source absent", func(t *testing.T) {
		_, _, args := seedReviewEvidence(t, false)
		if err := os.Remove(reviewPhase); err != nil {
			t.Fatal(err)
		}
		runGit(t, "add", "-A")
		runGit(t, "commit", "-q", "-m", "remove review source")
		args[6] = runGit(t, "rev-parse", "HEAD")
		assertReviewRefusal(t, args, reviewPhase)
	})

	for _, kind := range []string{"empty", "live symlink", "dangling symlink"} {
		t.Run("required source "+kind, func(t *testing.T) {
			_, _, args := seedReviewEvidence(t, false)
			if err := os.Remove(reviewPhase); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "empty":
				mustWriteFile(t, reviewPhase, "")
			case "live symlink":
				if err := os.Symlink("../skills/bench-craft-review/SKILL.md", reviewPhase); err != nil {
					t.Fatal(err)
				}
			case "dangling symlink":
				if err := os.Symlink("missing-review-source", reviewPhase); err != nil {
					t.Fatal(err)
				}
			}
			runGit(t, "add", "-A")
			runGit(t, "commit", "-q", "-m", "replace review source")
			args[6] = runGit(t, "rev-parse", "HEAD")
			assertReviewRefusal(t, args, reviewPhase)
		})
	}
}

func TestReviewChargeMovementDiscardsPayload(t *testing.T) {
	for _, movement := range []string{"head", "index", "required source"} {
		t.Run(movement, func(t *testing.T) {
			_, _, args := seedReviewEvidence(t, false)
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				suffix := string(rune('0' + calls))
				switch movement {
				case "head":
					mustWriteFile(t, "notes/head.txt", "movement "+suffix+"\n")
					runGit(t, "add", "notes/head.txt")
					runGit(t, "commit", "-q", "-m", "move head")
				case "index":
					mustWriteFile(t, "notes/index.txt", "movement "+suffix+"\n")
					runGit(t, "add", "notes/index.txt")
				case "required source":
					mustWriteFile(t, reviewPhase, "# Review phase\n\nmovement "+suffix+"\n")
				}
			})
			defer restore()
			out, code := Command(args)
			if code != 1 || calls != 2 || !strings.Contains(out, "snapshot drift") ||
				strings.Contains(out, "charge[") || strings.Contains(out, "shared_evidence") {
				t.Fatalf("persistent %s movement = (%d, calls=%d):\n%s", movement, code, calls, out)
			}
		})
	}
}

func assertReviewRefusal(t *testing.T, args []string, want string) {
	t.Helper()
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, want) || strings.Contains(out, "charge[") {
		t.Fatalf("review refusal = (%d), want %q without complete payload:\n%s", code, want, out)
	}
}

func TestReviewChargeGrammarRejectsTicket(t *testing.T) {
	_, _, args := seedReviewEvidence(t, false)
	args = append(args, "--ticket", "one.md")
	out, code := Command(args)
	const want = "--charge requires build with --ticket, or review without --ticket"
	if code != 2 || !strings.Contains(out, want) {
		t.Fatalf("review ticket grammar = (%d):\n%s", code, out)
	}
}
