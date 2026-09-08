package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/toon"
)

func chargeArgs(t *testing.T, root, slug string, full bool) []string {
	t.Helper()
	activeAssignment(t, root, root)
	args := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", runGit(t, "rev-parse", "main"), "--source-tip", runGit(t, "rev-parse", "HEAD")}
	if full {
		return append(args, "--full")
	}
	return args
}

func TestChargeIdentity(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, true)
	const assignment = "00000000000000000000000000000001"
	out, code := Command(args)
	if code != 0 {
		t.Fatalf("charge exit = %d:\n%s", code, out)
	}
	for _, want := range []string{
		"charge[1]{assignment,checkout,base,source_tip,fence,ticket,writes,evidence,checks,return,complete,next}",
		"sources[5]{path,identity}", "specs/example/tickets/one.md", args[6], args[8], "specs/example/spec.md", "\"" + assignment + "\"", "," + root + ",", ",specs,", "\"true\"",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("charge omitted %q:\n%s", want, out)
		}
	}
}

func TestChargeTicketEvidence(t *testing.T) {
	root, slug := seedConformant(t)
	out, code := Command(chargeArgs(t, root, slug, true))
	if code != 0 {
		t.Fatalf("full charge = (%d):\n%s", code, out)
	}
	for _, want := range []string{"Writes: specs", "Covers: PF1, PF2", "coverage[2]{row}", "  PF1", "  PF2", "## Acceptance"} {
		if !strings.Contains(out, want) {
			t.Errorf("full evidence omitted %q:\n%s", want, out)
		}
	}
}

func TestChargeCanonicalRequirements(t *testing.T) {
	for _, source := range []string{delegateSkill, delegateProcedure, buildPhase} {
		t.Run(source, func(t *testing.T) {
			root, slug := seedConformant(t)
			before, code := Command(chargeArgs(t, root, slug, true))
			if code != 0 {
				t.Fatalf("initial full charge = (%d):\n%s", code, before)
			}
			changed := "# Changed canonical source\n\n" + source + " changed.\n"
			mustWriteFile(t, source, changed)
			runGit(t, "add", source)
			runGit(t, "commit", "-q", "-m", "change canonical source")
			args := chargeArgs(t, root, slug, true)
			args[8] = runGit(t, "rev-parse", "HEAD")
			after, code := Command(args)
			if code != 0 || !strings.Contains(after, strings.ReplaceAll(changed, "\n", "\\n")) || !strings.Contains(after, (chargeSource{path: source, data: []byte(changed)}).identity()) || after == before {
				t.Fatalf("changed canonical source = (%d):\n%s", code, after)
			}
		})
	}
}

func TestChargeRequiredInputs(t *testing.T) {
	_, slug := seedConformant(t)
	out, code := Command([]string{"build", slug, "--charge", "--ticket", "one.md"})
	if code != 2 || !strings.Contains(out, "--charge requires build, --ticket, --base, and --source-tip") {
		t.Fatalf("incomplete charge = (%d):\n%s", code, out)
	}
	root, slug := seedConformant(t)
	args := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", runGit(t, "rev-parse", "main"), "--source-tip", runGit(t, "rev-parse", "HEAD")}
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "assignment required") || !strings.Contains(out, "assigned worktree") {
		t.Fatalf("missing assignment = (%d):\n%s", code, out)
	}
	activeAssignment(t, root, t.TempDir())
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "assignment required") {
		t.Fatalf("foreign assignment = (%d):\n%s", code, out)
	}
	activeAssignment(t, root, root)
	args[4] = "missing.md"
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "selected ticket") || !strings.Contains(out, "pass a ticket basename") {
		t.Fatalf("missing ticket = (%d):\n%s", code, out)
	}
}

func TestChargeSnapshotMovement(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, false)
	calls := 0
	restore := diff.SetSnapshotAfterReadForTest(func() {
		calls++
		mustWriteFile(t, "internal/example/foo.go", "package example\n// moved "+string(rune('0'+calls))+"\n")
	})
	out, code := Command(args)
	if code != 1 || calls != 2 || !strings.Contains(out, "error: snapshot drift") || !strings.Contains(out, "retry the exact invocation") {
		t.Fatalf("persistent charge movement = (%d, %d):\n%s", code, calls, out)
	}
	restore()

	root, slug = seedConformant(t)
	args = chargeArgs(t, root, slug, false)
	args[len(args)-1] = runGit(t, "rev-parse", "main")
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "tip-current") {
		t.Fatalf("mismatched pair = (%d):\n%s", code, out)
	}

	root, slug = seedConformant(t)
	args = chargeArgs(t, root, slug, false)
	calls = 0
	restore = diff.SetSnapshotAfterReadForTest(func() {
		calls++
		if calls == 1 {
			mustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(specBody(slug), "Status: staged", "Status: draft", 1))
		}
	})
	out, code = Command(args)
	restore()
	if code != 1 || calls != 2 || !strings.Contains(out, "spec not staged") || strings.Contains(out, "snapshot drift") || strings.Contains(out, "complete") {
		t.Fatalf("movement then bootstrap failure = (%d, %d):\n%s", code, calls, out)
	}
}

func TestChargeProjectionAndFullRetrieval(t *testing.T) {
	root, slug := seedConformant(t)
	compact, code := Command(chargeArgs(t, root, slug, false))
	if code != 0 || !strings.Contains(compact, "\"false\",bench preflight") || !strings.Contains(compact, "omitted[5]{source}") || !strings.Contains(compact, "--full") || strings.Contains(compact, "## Acceptance") {
		t.Fatalf("compact charge = (%d):\n%s", code, compact)
	}
	full, code := Command(chargeArgs(t, root, slug, true))
	if code != 0 || !strings.Contains(full, "\"true\"") || !strings.Contains(full, "evidence[5]{path,content}") || !strings.Contains(full, "## Acceptance") || !strings.Contains(full, "Status: staged") || !strings.Contains(full, "Delegation skill") || !strings.Contains(full, "Build phase") || !strings.Contains(full, "Focused suite:") {
		t.Fatalf("full charge = (%d):\n%s", code, full)
	}

	root, slug = seedConformant(t)
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticketDoc("One", "PF1", "PF2")+strings.Repeat("large ticket evidence\n", 12000))
	runGit(t, "add", "specs/"+slug+"/tickets/one.md")
	runGit(t, "commit", "-q", "-m", "large ticket")
	args := chargeArgs(t, root, slug, false)
	args[8] = runGit(t, "rev-parse", "HEAD")
	compact, code = Command(args)
	if code != 0 || len(compact) > 10000 || strings.Contains(compact, "large ticket evidence") || !strings.Contains(compact, "omitted[5]{source}") {
		t.Fatalf("large compact charge = (%d, %d bytes):\n%s", code, len(compact), compact)
	}
}

func TestChargeHostileInputs(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, false)
	args[4] = "one.md; touch sentinel"
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, "selected ticket") {
		t.Fatalf("hostile ticket name = (%d):\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("ticket name created sentinel: %v", err)
	}

	root, slug = seedConformant(t)
	mustWriteFile(t, buildPhase, "bad\x1b\n")
	runGit(t, "add", buildPhase)
	runGit(t, "commit", "-q", "-m", "hostile phase")
	args = chargeArgs(t, root, slug, false)
	args[8] = runGit(t, "rev-parse", "HEAD")
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "source required") || !strings.Contains(out, buildPhase) {
		t.Fatalf("hostile source text = (%d):\n%s", code, out)
	}
}

func TestChargeRefusesSpecialRequiredSource(t *testing.T) {
	for _, test := range []struct {
		name    string
		prepare func(t *testing.T)
	}{
		{
			name: "dangling link",
			prepare: func(t *testing.T) {
				if err := os.Remove(buildPhase); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("missing", buildPhase); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "empty",
			prepare: func(t *testing.T) { mustWriteFile(t, buildPhase, "") },
		},
		{
			name: "absent",
			prepare: func(t *testing.T) {
				if err := os.Remove(buildPhase); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedConformant(t)
			test.prepare(t)
			runGit(t, "add", "-A")
			runGit(t, "commit", "-q", "-m", "special phase")
			args := chargeArgs(t, root, slug, false)
			args[8] = runGit(t, "rev-parse", "HEAD")
			out, code := Command(args)
			if code != 1 || !strings.Contains(out, "source required") || !strings.Contains(out, buildPhase) {
				t.Fatalf("%s source = (%d):\n%s", test.name, code, out)
			}
		})
	}
	t.Run("fifo", func(t *testing.T) {
		root, slug := seedConformant(t)
		if err := os.Remove(buildPhase); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(buildPhase, 0o600); err != nil {
			t.Fatal(err)
		}
		out, code := Command(chargeArgs(t, root, slug, false))
		if code != 1 || !strings.Contains(out, "checkout required") {
			t.Fatalf("fifo source = (%d):\n%s", code, out)
		}
	})
}

func TestLegacyPreflightDifferential(t *testing.T) {
	type result struct {
		out             string
		code            int
		root, base, tip string
	}
	run := func(t *testing.T, args []string, root, base, tip string) result {
		out, code := Command(args)
		return result{out, code, root, base, tip}
	}
	cases := []struct {
		name string
		run  func(*testing.T) result
	}{
		{"valid-build", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			return run(t, []string{"build", slug}, root, base, tip)
		}},
		{"valid-review", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			return run(t, []string{"review", slug}, root, base, tip)
		}},
		{"absent-tickets", func(t *testing.T) result {
			root, slug := seedBuildFresh(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			return run(t, []string{"build", slug}, root, base, tip)
		}},
		{"empty-tickets", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base := runGit(t, "rev-parse", "main")
			if err := os.Remove("specs/" + slug + "/tickets/one.md"); err != nil {
				t.Fatal(err)
			}
			runGit(t, "add", "-A")
			runGit(t, "commit", "-q", "-m", "empty tickets")
			return run(t, []string{"build", slug}, root, base, runGit(t, "rev-parse", "HEAD"))
		}},
		{"stale-base", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			runGit(t, "checkout", "-q", "main")
			mustWriteFile(t, "advance.txt", "advance\n")
			runGit(t, "add", "advance.txt")
			runGit(t, "commit", "-q", "-m", "advance")
			runGit(t, "checkout", "-q", "feature")
			return run(t, []string{"build", slug}, root, base, tip)
		}},
		{"dirty-review", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			mustWriteFile(t, "dirty.txt", "dirty\n")
			return run(t, []string{"review", slug, "--base", base}, root, base, tip)
		}},
		{"explicit-base-success", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			return run(t, []string{"build", slug, "--base", base}, root, base, tip)
		}},
		{"source-tip-mismatch", func(t *testing.T) result {
			root, slug := seedConformant(t)
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			return run(t, []string{"build", slug, "--base", base, "--source-tip", base}, root, base, tip)
		}},
		{"invalid-invocation", func(t *testing.T) result {
			root, _ := seedConformant(t)
			return run(t, []string{"unknown", "example"}, root, "", "")
		}},
		{"empty-diff", func(t *testing.T) result {
			root := initRepo(t)
			slug := "example"
			mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
			mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticketDoc("One", "PF1", "PF2"))
			runGit(t, "add", ".")
			runGit(t, "commit", "-q", "-m", "c0")
			runGit(t, "checkout", "-q", "-b", "feature")
			base, tip := runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD")
			return run(t, []string{"review", slug, "--base", base}, root, base, tip)
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := test.run(t)
			got.out = normalizeLegacy(got.out, got.root, got.base, got.tip)
			want, wantCode := legacyBaseline(test.name)
			if got.code != wantCode || got.out != want {
				t.Fatalf("legacy %s = (%d, %q), want (%d, %q)", test.name, got.code, got.out, wantCode, want)
			}
		})
	}
}
func normalizeLegacy(out, root, base, tip string) string {
	if root != "" {
		out = normalizeLegacyValue(out, root, "<root>")
	}
	if base != "" {
		out = normalizeLegacyValue(out, base, "<base>")
	}
	if tip != "" && tip != base {
		out = normalizeLegacyValue(out, tip, "<tip>")
	}
	return out
}
func normalizeLegacyValue(out, value, replacement string) string {
	from, err := toon.Table("value", []string{"v"}, [][]string{{value}})
	if err != nil {
		panic(err)
	}
	to, err := toon.Table("value", []string{"v"}, [][]string{{replacement}})
	if err != nil {
		panic(err)
	}
	from = strings.TrimSuffix(strings.TrimPrefix(from, "value[1]{v}:\n  "), "\n")
	to = strings.TrimSuffix(strings.TrimPrefix(to, "value[1]{v}:\n  "), "\n")
	return strings.ReplaceAll(strings.ReplaceAll(out, from, to), value, replacement)
}
func TestNormalizeLegacyIdentityCells(t *testing.T) {
	identity := strings.Repeat("1", 40)
	got := normalizeLegacy("source[1]{base,tip}:\n  \""+identity+"\",x\n", "", identity, "")
	if got != "source[1]{base,tip}:\n  <base>,x\n" {
		t.Fatalf("normalized numeric identity = %q", got)
	}
}
func legacyBaseline(name string) (string, int) {
	const build = "phase: build\nspec: specs/example/spec.md\n"
	const review = "phase: review\nspec: specs/example/spec.md\n"
	const source = "source[1]{base,tip}:\n  <base>,<tip>\n"
	const checks = "  paths-authorized,green,\"\",\"\"\n  tickets-parse,green,\"\",\"\"\n  blockers-resolve,green,\"\",\"\"\n  writes-resolve,green,\"\",\"\"\n  fixture-closure,green,\"\",\"\"\n  registry-closure,green,\"\",\"\"\n  kit-pin,green,\"\",\"\"\n"
	const buildTail = "  binary-seal,not-applicable,\"\",\"\"\n  rows-owned,green,\"\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,not-applicable,\"\",\"\"\n"
	const reviewTail = "  rows-owned,green,\"\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,green,\"\",\"\"\n"
	greenBuild := build + "checks[12]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + buildTail
	greenReview := review + "checks[11]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + reviewTail
	switch name {
	case "valid-build":
		return greenBuild, 0
	case "valid-review":
		return greenReview, 0
	case "absent-tickets":
		return build + "checks[12]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n  paths-authorized,green,\"\",\"\"\n  tickets-parse,not-applicable,\"\",\"\"\n  blockers-resolve,not-applicable,\"\",\"\"\n  writes-resolve,not-applicable,\"\",\"\"\n  fixture-closure,not-applicable,\"\",\"\"\n  registry-closure,not-applicable,\"\",\"\"\n  kit-pin,not-applicable,\"\",\"\"\n  binary-seal,not-applicable,\"\",\"\"\n  rows-owned,not-applicable,\"\",\"\"\n  rows-membership,not-applicable,\"\",\"\"\n  diff-nonempty,not-applicable,\"\",\"\"\n", 0
	case "empty-tickets":
		return build + "checks[12]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + "  binary-seal,not-applicable,\"\",\"\"\n  rows-owned,red,\"declared row(s) cited by no ticket file: PF1, PF2\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,not-applicable,\"\",\"\"\n", 1
	case "stale-base":
		return build + "checks[12]{check,verdict,detail,next}:\n  base-current,red,default branch tip is not an ancestor of HEAD,bench worktree merge --from main <target>\n" + checks + buildTail, 1
	case "dirty-review":
		return "error: source not clean — review source has uncommitted changes\n", 1
	case "explicit-base-success":
		return build + source + strings.TrimPrefix(greenBuild, build), 0
	case "source-tip-mismatch":
		return build + source + "checks[13]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n  tip-current,red,\"--source-tip <base> is not the derived source tip <tip>\",\"\"\n" + checks + buildTail, 1
	case "invalid-invocation":
		return "usage: bench preflight (unknown argument: unknown)\n", 2
	case "empty-diff":
		return review + "source[1]{base,tip}:\n  <base>,<base>\nchecks[11]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + "  rows-owned,green,\"\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,red,no changed files since the resolved review base,\"\"\n", 1
	}
	return "", 0
}
