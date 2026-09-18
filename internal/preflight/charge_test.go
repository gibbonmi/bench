package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/toon"
)

func TestChargeRequiredInputs(t *testing.T) {
	_, slug := preflighttest.SeedConformant(t)
	out, code := Command([]string{"build", slug, "--charge", "--ticket", "one.md"})
	if code != 2 || !strings.Contains(out, "--charge requires build, --ticket, --base, and --source-tip") {
		t.Fatalf("incomplete charge = (%d):\n%s", code, out)
	}
	root, slug := preflighttest.SeedConformant(t)
	args := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "assignment required") || !strings.Contains(out, "assigned worktree") {
		t.Fatalf("missing assignment = (%d):\n%s", code, out)
	}
	preflighttest.ActiveAssignment(t, root, t.TempDir())
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "assignment required") {
		t.Fatalf("foreign assignment = (%d):\n%s", code, out)
	}
	preflighttest.ActiveAssignment(t, root, root)
	args[4] = "missing.md"
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "selected ticket") || !strings.Contains(out, "pass a ticket basename") {
		t.Fatalf("missing ticket = (%d):\n%s", code, out)
	}
}

func TestChargeSnapshotMovement(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	args := preflighttest.ChargeArgs(t, root, slug)
	calls := 0
	restore := diff.SetSnapshotAfterReadForTest(func() {
		calls++
		preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n// moved "+string(rune('0'+calls))+"\n")
	})
	out, code := Command(args)
	if code != 1 || calls != 2 || !strings.Contains(out, "error: snapshot drift") || !strings.Contains(out, "retry the exact invocation") {
		t.Fatalf("persistent charge movement = (%d, %d):\n%s", code, calls, out)
	}
	restore()

	root, slug = preflighttest.SeedConformant(t)
	args = preflighttest.ChargeArgs(t, root, slug)
	args[len(args)-1] = preflighttest.RunGit(t, "rev-parse", "main")
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "tip-current") {
		t.Fatalf("mismatched pair = (%d):\n%s", code, out)
	}

	root, slug = preflighttest.SeedConformant(t)
	args = preflighttest.ChargeArgs(t, root, slug)
	calls = 0
	restore = diff.SetSnapshotAfterReadForTest(func() {
		calls++
		if calls == 1 {
			preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(preflighttest.SpecBody(slug), "Status: staged", "Status: draft", 1))
		}
	})
	out, code = Command(args)
	restore()
	if code != 1 || calls != 2 || !strings.Contains(out, "spec not staged") || strings.Contains(out, "snapshot drift") || strings.Contains(out, "complete,next}") {
		t.Fatalf("movement then bootstrap failure = (%d, %d):\n%s", code, calls, out)
	}
}

func TestChargeHostileInputs(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	args := preflighttest.ChargeArgs(t, root, slug)
	args[4] = "one.md; touch sentinel"
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, "selected ticket") {
		t.Fatalf("hostile ticket name = (%d):\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("ticket name created sentinel: %v", err)
	}

	root, slug = preflighttest.SeedConformant(t)
	preflighttest.MustWriteFile(t, chargesource.BuildPhase, "bad\x1b\n")
	preflighttest.RunGit(t, "add", chargesource.BuildPhase)
	preflighttest.RunGit(t, "commit", "-q", "-m", "hostile phase")
	args = preflighttest.ChargeArgs(t, root, slug)
	args[8] = preflighttest.RunGit(t, "rev-parse", "HEAD")
	out, code = Command(args)
	if code != 1 || !strings.Contains(out, "source required") || !strings.Contains(out, chargesource.BuildPhase) {
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
				if err := os.Remove(chargesource.BuildPhase); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("missing", chargesource.BuildPhase); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name:    "empty",
			prepare: func(t *testing.T) { preflighttest.MustWriteFile(t, chargesource.BuildPhase, "") },
		},
		{
			name: "absent",
			prepare: func(t *testing.T) {
				if err := os.Remove(chargesource.BuildPhase); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := preflighttest.SeedConformant(t)
			test.prepare(t)
			preflighttest.RunGit(t, "add", "-A")
			preflighttest.RunGit(t, "commit", "-q", "-m", "special phase")
			args := preflighttest.ChargeArgs(t, root, slug)
			args[8] = preflighttest.RunGit(t, "rev-parse", "HEAD")
			out, code := Command(args)
			if code != 1 || !strings.Contains(out, "source required") || !strings.Contains(out, chargesource.BuildPhase) {
				t.Fatalf("%s source = (%d):\n%s", test.name, code, out)
			}
		})
	}
	t.Run("fifo", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		if err := os.Remove(chargesource.BuildPhase); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(chargesource.BuildPhase, 0o600); err != nil {
			t.Fatal(err)
		}
		out, code := Command(preflighttest.ChargeArgs(t, root, slug))
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
	run := func(args []string, root, base, tip string) result {
		out, code := Command(args)
		return result{out, code, root, base, tip}
	}
	cases := []struct {
		name string
		run  func(*testing.T) result
	}{
		{"valid-build", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			return run([]string{"build", slug}, root, base, tip)
		}},
		{"valid-review", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			return run([]string{"review", slug}, root, base, tip)
		}},
		{"absent-tickets", func(t *testing.T) result {
			root, slug := seedBuildFresh(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			return run([]string{"build", slug}, root, base, tip)
		}},
		{"empty-tickets", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base := preflighttest.RunGit(t, "rev-parse", "main")
			if err := os.Remove("specs/" + slug + "/tickets/one.md"); err != nil {
				t.Fatal(err)
			}
			preflighttest.RunGit(t, "add", "-A")
			preflighttest.RunGit(t, "commit", "-q", "-m", "empty tickets")
			return run([]string{"build", slug}, root, base, preflighttest.RunGit(t, "rev-parse", "HEAD"))
		}},
		{"stale-base", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			preflighttest.RunGit(t, "checkout", "-q", "main")
			preflighttest.MustWriteFile(t, "advance.txt", "advance\n")
			preflighttest.RunGit(t, "add", "advance.txt")
			preflighttest.RunGit(t, "commit", "-q", "-m", "advance")
			preflighttest.RunGit(t, "checkout", "-q", "feature")
			return run([]string{"build", slug}, root, base, tip)
		}},
		{"dirty-review", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			preflighttest.MustWriteFile(t, "dirty.txt", "dirty\n")
			return run([]string{"review", slug, "--base", base}, root, base, tip)
		}},
		{"explicit-base-success", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			return run([]string{"build", slug, "--base", base}, root, base, tip)
		}},
		{"source-tip-mismatch", func(t *testing.T) result {
			root, slug := preflighttest.SeedConformant(t)
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			return run([]string{"build", slug, "--base", base, "--source-tip", base}, root, base, tip)
		}},
		{"invalid-invocation", func(t *testing.T) result {
			root, _ := preflighttest.SeedConformant(t)
			return run([]string{"unknown", "example"}, root, "", "")
		}},
		{"empty-diff", func(t *testing.T) result {
			root := preflighttest.StartRepo(t)
			slug := "example"
			preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug))
			preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.TicketDoc("One", "PF1", "PF2"))
			preflighttest.RunGit(t, "add", ".")
			preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
			preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
			base, tip := preflighttest.RunGit(t, "rev-parse", "main"), preflighttest.RunGit(t, "rev-parse", "HEAD")
			return run([]string{"review", slug, "--base", base}, root, base, tip)
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
	const checks = "  paths-authorized,green,\"\",\"\"\n  tickets-parse,green,\"\",\"\"\n  completion-plan,green,\"\",\"\"\n  blockers-resolve,green,\"\",\"\"\n  writes-resolve,green,\"\",\"\"\n  fixture-closure,green,\"\",\"\"\n  registry-closure,green,\"\",\"\"\n  anchor-closure,green,\"\",\"\"\n  kit-pin,green,\"\",\"\"\n"
	const buildTail = "  binary-seal,not-applicable,\"\",\"\"\n  rows-owned,green,\"\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,not-applicable,\"\",\"\"\n"
	const reviewTail = "  rows-owned,green,\"\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,green,\"\",\"\"\n"
	greenBuild := build + "checks[14]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + buildTail
	greenReview := review + "checks[13]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + reviewTail
	switch name {
	case "valid-build":
		return greenBuild, 0
	case "valid-review":
		return greenReview, 0
	case "absent-tickets":
		return build + "checks[14]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n  paths-authorized,green,\"\",\"\"\n  tickets-parse,not-applicable,\"\",\"\"\n  completion-plan,not-applicable,\"\",\"\"\n  blockers-resolve,not-applicable,\"\",\"\"\n  writes-resolve,not-applicable,\"\",\"\"\n  fixture-closure,not-applicable,\"\",\"\"\n  registry-closure,not-applicable,\"\",\"\"\n  anchor-closure,not-applicable,\"\",\"\"\n  kit-pin,not-applicable,\"\",\"\"\n  binary-seal,not-applicable,\"\",\"\"\n  rows-owned,not-applicable,\"\",\"\"\n  rows-membership,not-applicable,\"\",\"\"\n  diff-nonempty,not-applicable,\"\",\"\"\n", 0
	case "empty-tickets":
		return build + "checks[14]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + strings.Replace(checks, "  completion-plan,green,\"\",\"\"\n", "  completion-plan,red,\"spec carries no valid bench-completion-plan fence at <tip>: missing or nonregular tree file specs/example/tickets/one.md; see .bench/BENCH-reference.md, bench gate --checkpoint\",\"\"\n", 1) + "  binary-seal,not-applicable,\"\",\"\"\n  rows-owned,red,\"declared row(s) cited by no ticket file: PF1, PF2\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,not-applicable,\"\",\"\"\n", 1
	case "stale-base":
		return build + "checks[14]{check,verdict,detail,next}:\n  base-current,red,default branch tip is not an ancestor of HEAD,bench worktree merge --from main <target>\n" + checks + buildTail, 1
	case "dirty-review":
		return "error: source not clean — review source has uncommitted changes\n", 1
	case "explicit-base-success":
		return build + source + strings.TrimPrefix(greenBuild, build), 0
	case "source-tip-mismatch":
		return build + source + "checks[15]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n  tip-current,red,\"--source-tip <base> is not the derived source tip <tip>\",\"\"\n" + checks + buildTail, 1
	case "invalid-invocation":
		return "usage: bench preflight (unknown argument: unknown)\n", 2
	case "empty-diff":
		return review + "source[1]{base,tip}:\n  <base>,<base>\nchecks[13]{check,verdict,detail,next}:\n  base-current,green,\"\",\"\"\n" + checks + "  rows-owned,green,\"\",\"\"\n  rows-membership,green,\"\",\"\"\n  diff-nonempty,red,no changed files since the resolved review base,\"\"\n", 1
	}
	return "", 0
}
