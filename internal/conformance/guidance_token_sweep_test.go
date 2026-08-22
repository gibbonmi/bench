package conformance

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/modelid"
)

// guidanceDirs are the two trees kit guidance prose lives in. Their contents are
// discovered at run time rather than listed here. A list grades only the files its author
// knew about. The skill somebody adds next month is exactly where a hard-coded binding
// rots back in with nothing watching it.
var guidanceDirs = []string{
	filepath.Join(".agents", "commands"),
	filepath.Join(".agents", "skills"),
}

// retiredKeyRe matches a retired key by its concrete tier and by the `*` glob prose
// writes a key family with. Both spellings teach the retired schema. The glob is how the
// drift this check was written against was actually spelled. The stems come from
// lines.RetiredKeyPrefixes rather than a second list here. Prose still naming a retired
// family teaches a schema nothing reads. Which families those are is one fact that the
// doctor's migration report already owns.
var retiredKeyRe = regexp.MustCompile(`(?:` + strings.Join(quotedRetiredPrefixes(), "|") +
	`)(?:` + strings.Join(upperTiers(), "|") + `|\*)`)

// quotedRetiredPrefixes renders the retired stems as regexp literals. A stem carrying a
// metacharacter then matches itself, rather than silently widening the matcher.
func quotedRetiredPrefixes() []string {
	prefixes := lines.RetiredKeyPrefixes()
	quoted := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		quoted = append(quoted, regexp.QuoteMeta(prefix))
	}
	return quoted
}

// upperTiers renders lines.Tiers in the case a binding key spells them. The tier
// vocabulary keeps one declaration, so a tier added there reaches this sweep unedited.
func upperTiers() []string {
	upper := make([]string, 0, len(lines.Tiers))
	for _, tier := range lines.Tiers {
		upper = append(upper, strings.ToUpper(tier))
	}
	return upper
}

// codeSpanRe captures the content of one inline code span.
var codeSpanRe = regexp.MustCompile("`([^`\n]*)`")

// checkGuidanceTokens sweeps kit guidance prose for a binding hard-coded where nothing
// else grades it. Two token classes fail: a model-id literal naming no cell of the parsed
// matrix, and a retired schema key. The allowlist is the parse itself, never a written
// list of tokens. Rebinding a tier therefore stays a one-file edit and cannot leave this
// check behind asserting the previous binding.
//
// The model-literal arm is guarded on .bench/lines.env being present. With no binding
// there is no allowlist. Reporting every literal in the tree would bury the one
// diagnostic checkLineBinding already emits for the missing file. The retired-key arm
// needs no binding and runs either way.
func checkGuidanceTokens(root string) []string {
	bindingPath := filepath.Join(root, ".bench", "lines.env")
	binding, bound := lines.ParseBinding([]byte(readIfExists(bindingPath))), exists(bindingPath)
	var diags []string
	for _, dir := range guidanceDirs {
		files, discovery := discoverGuidanceFiles(root, filepath.Join(root, dir))
		diags = append(diags, discovery...)
		for _, path := range files {
			rel := slashRel(root, path)
			// Discovery classified the entry's shape; this classifies its bytes. An oversized or
			// non-UTF-8 guidance file is reported, rather than swept as an empty document that
			// trivially carries no offending token.
			text := bounds.ClassifyNoFollow(path)
			if text.State.Failed() {
				diags = append(diags, fmt.Sprintf("guidance file refused: %s could not be read (%s)", rel, text.Reason))
				continue
			}
			diags = append(diags, guidanceFileDiags(rel, string(text.Data), binding, bound)...)
		}
	}
	return diags
}

// discoverGuidanceFiles returns dir's regular files in sorted order. It also returns one
// diagnostic per entry that is neither a regular file nor a directory. Classification is
// by Lstat and precedes every read. A FIFO with no writer blocks forever in open(2) for a
// reader that opens before it looks. Following a symlink would pull bytes from outside
// the swept tree into a diagnostic naming a path inside it. An absent directory yields
// nothing.
func discoverGuidanceFiles(root, dir string) (files, diags []string) {
	// os.ReadDir sorts by filename, so the sweep reports in one order run to run.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		info, err := os.Lstat(path)
		if err != nil {
			diags = append(diags, fmt.Sprintf("guidance entry unreadable: %s could not be classified (%v)", slashRel(root, path), err))
			continue
		}
		switch mode := info.Mode(); {
		case mode.IsDir():
			nested, nestedDiags := discoverGuidanceFiles(root, path)
			files, diags = append(files, nested...), append(diags, nestedDiags...)
		case mode.IsRegular():
			files = append(files, path)
		default:
			diags = append(diags, fmt.Sprintf("guidance entry not a regular file: %s is %s and the sweep classifies every entry before it reads one", slashRel(root, path), mode.Type()))
		}
	}
	return files, diags
}

// guidanceFileDiags reports one diagnostic per distinct offending token in one file. Each
// diagnostic names the file, the token, and .bench/lines.env. The fix is an edit to one
// of the two, and the reader must be told which.
func guidanceFileDiags(rel, text string, binding lines.Binding, bound bool) []string {
	var diags []string
	for _, key := range uniqueSorted(retiredKeyRe.FindAllString(text, -1)) {
		diags = append(diags, fmt.Sprintf("guidance names a retired binding key: %s writes '%s', which the BENCH_<HARNESS>_<TIER> matrix in .bench/lines.env replaced", rel, key))
	}
	if !bound {
		return diags
	}
	var unbound []string
	for _, span := range codeSpanRe.FindAllStringSubmatch(text, -1) {
		if token := span[1]; modelLiteral(token) && !bindsToken(binding, token) {
			unbound = append(unbound, token)
		}
	}
	for _, token := range uniqueSorted(unbound) {
		diags = append(diags, fmt.Sprintf("guidance hard-codes an unbound model id: %s names '%s', which is no cell of the binding matrix in .bench/lines.env — name a tier, or bind the id", rel, token))
	}
	return diags
}

// bindsToken reports whether token is some harness's bound cell. Acceptance spans the
// whole matrix on purpose, because guidance may legitimately quote another harness's
// column. Which harness owns which cell is checkLineBinding's question, not this sweep's.
func bindsToken(binding lines.Binding, token string) bool {
	for _, harness := range lines.Harnesses {
		for _, tier := range lines.Tiers {
			if binding.Cell(harness, tier) == token {
				return true
			}
		}
	}
	return false
}

// modelLiteral reports whether span is shaped like a concrete model id. The shape is a
// safe token, at most one `provider/` prefix, non-empty hyphen-separated segments
// starting with a letter, and a version digit. The shape also needs either that provider
// prefix or three or more segments.
//
// The boundary is drawn tight, because the expensive failure here is the opposite one. A
// matcher broad enough to catch every id turns ordinary guidance red, and a prose oracle
// that cries wolf gets weakened rather than obeyed. The version digit separates a model
// id from every hyphenated English compound. The three-segment floor keeps `utf-8` and
// `schema-3` out, and the single-slash rule keeps a path like `docs/adr/0001-tripwire.md`
// out.
//
// Where the version digit may sit depends on the shape around it. Normally it must open a
// later segment. A leading numbered word is how ordinary slugs are spelled, such as
// `ft128-agent-line-binding`, and reading that as a version would report them.
//
// The one exception is a provider-qualified id of exactly two segments. There, the digit
// may sit inside the first segment, the way the `openai/o3-mini` family is spelled.
// Holding the exception to two segments is what keeps a qualified slug like
// `specs/ft128-agent-line-binding` accepted. The two shapes are otherwise identical, and
// prose writes far more slugs than ids.
//
// The accepted misses are three shapes. A two-segment dotted id, like `gpt-5.6`, is one.
// A provider-qualified id of three or more segments carrying its only digit in the first,
// like `openai/o3-mini-high`, is another. An id written outside a code span is the third.
// The third is the boundary's price. A matcher wide enough to close it reports the slugs
// above, and a prose oracle that cries wolf gets weakened rather than obeyed.
func modelLiteral(span string) bool {
	if !modelid.SafeToken(span) {
		return false
	}
	body, qualified := span, false
	if provider, rest, found := strings.Cut(span, "/"); found {
		if provider == "" || strings.Contains(rest, "/") {
			return false
		}
		body, qualified = rest, true
	}
	segments := strings.Split(body, "-")
	if len(segments) < 2 || segments[0] == "" || !isASCIILetter(segments[0][0]) {
		return false
	}
	if !qualified && len(segments) < 3 {
		return false
	}
	for _, segment := range segments {
		if segment == "" {
			return false
		}
	}
	if qualified && len(segments) == 2 && containsDigit(segments[0]) {
		return true
	}
	for _, segment := range segments[1:] {
		if segment[0] >= '0' && segment[0] <= '9' {
			return true
		}
	}
	return false
}

func isASCIILetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func containsDigit(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}

// guidanceFixtureEnv is the fixture binding every sweep case is graded against. Its
// tokens are deliberately unlike the repo's own. A check that fell back to the live
// binding rather than the parsed fixture would fail these cases rather than pass them.
const guidanceFixtureEnv = "BENCH_CODEX_TOP=alpha-9.1-sol\nBENCH_CODEX_MID=alpha-9.1-terra\n" +
	"BENCH_CODEX_CHEAP=alpha-9.1-luna\nBENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_MID=opus\n" +
	"BENCH_CLAUDE_CHEAP=sonnet\n"

// ordinaryGuidanceSpans are non-model literals kit prose really writes in code spans.
// Each is a shape a broader matcher would report. They ride the clean baseline, so the
// accepted side of the boundary is asserted rather than assumed. The last two are the
// shapes the provider-qualified first-segment rule brings closest to the line. A numbered
// slug under a one-segment path keeps its later segments free of a version digit. A
// provider-shaped path whose leading word carries no digit has no version at all.
var ordinaryGuidanceSpans = []string{
	"schema-3", "utf-8", "sha-256", "x86-64", "bench-craft-line", "red/green",
	"docs/adr/0001-working-tree-gate-tripwire.md", ".bench/lines.env", "v0.6.0",
	"2026-07-31", "ft128-agent-line-binding", "BENCH_MODEL", "BENCH_CODEX_TOP",
	"specs/ft128-agent-line-binding", "internal/conformance-sweep",
}

// writeGuidanceFixture builds the clean baseline: a binding, a command file, and a nested
// skill. The skill's prose names every bound cell of the matrix and every ordinary
// literal above. The bound cells are rendered from the same parse the check reads, so the
// fixture's tokens have one author.
//
// Every mutation case starts from a fresh copy of this tree and adds exactly one defect.
// That single defect is what makes a diagnostic attributable to it. The root is a short
// temp path rather than t.TempDir's test-named one. A Unix socket address is capped near
// 108 bytes. A subtest's name in the path would push the non-regular case past that cap,
// and turn a real coverage row into a host-capability skip.
func writeGuidanceFixture(t *testing.T) string {
	t.Helper()
	root, err := os.MkdirTemp("", "bench-guidance-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(root) })
	binding := lines.ParseBinding([]byte(guidanceFixtureEnv))
	var cells strings.Builder
	for _, harness := range lines.Harnesses {
		for _, tier := range lines.Tiers {
			if value := binding.Cell(harness, tier); value != "" {
				cells.WriteString("- " + harness + " " + tier + ": `" + value + "`\n")
			}
		}
	}
	var ordinary strings.Builder
	for _, span := range ordinaryGuidanceSpans {
		ordinary.WriteString("- `" + span + "`\n")
	}
	writeGuidanceFile(t, root, filepath.Join(".bench", "lines.env"), guidanceFixtureEnv)
	writeGuidanceFile(t, root, filepath.Join(".agents", "commands", "bench-setup-repo.md"),
		"# setup\n\nAsk which tokens bind `cheap`, `mid`, and `top`, then record the matrix.\n\n"+
			cells.String()+"\nOrdinary prose names a schema-3 payload and these literals:\n\n"+ordinary.String())
	writeGuidanceFile(t, root, filepath.Join(".agents", "skills", "bench-craft-line", "SKILL.md"),
		"# the line\n\nThe reviewer binds `top` / `mid` / `cheap` per harness; declare the tier.\n\n"+cells.String())
	return root
}

// writeGuidanceFile writes one fixture file at a slash-free relative path, creating its
// parents. It returns the absolute path, so a case can mutate what it just planted.
func writeGuidanceFile(t *testing.T, root, rel, body string) string {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestGuidanceSweepAcceptsBoundTokensAndOrdinaryProse is the clean-baseline row that
// stops this check from becoming a broad text matcher. Every bound matrix cell and every
// ordinary hyphenated literal stays accepted. The test fails for any allowlist but the
// parsed binding.
func TestGuidanceSweepAcceptsBoundTokensAndOrdinaryProse(t *testing.T) {
	if diags := checkGuidanceTokens(writeGuidanceFixture(t)); len(diags) != 0 {
		t.Fatalf("the clean fixture got diagnostics:\n%s", strings.Join(diags, "\n"))
	}
	for _, span := range ordinaryGuidanceSpans {
		t.Run(span, func(t *testing.T) {
			root := writeGuidanceFixture(t)
			writeGuidanceFile(t, root, filepath.Join(".agents", "commands", "ordinary.md"), "Prose: `"+span+"`.\n")
			if diags := checkGuidanceTokens(root); len(diags) != 0 {
				t.Fatalf("ordinary literal %q was reported:\n%s", span, strings.Join(diags, "\n"))
			}
		})
	}
}

// TestGuidanceSweepBitesOneUnboundModelLiteral plants exactly one unbound id in an
// otherwise clean copy. It asserts that its file-and-token diagnostic arrives alone. The
// clean baseline is what makes one diagnostic proof rather than coincidence, because
// collateral would be a second.
func TestGuidanceSweepBitesOneUnboundModelLiteral(t *testing.T) {
	for _, token := range []string{"beta-9.9-sol", "openai/gpt-9", "claude-opus-9-20991231", "openai/o3-mini"} {
		t.Run(token, func(t *testing.T) {
			root := writeGuidanceFixture(t)
			rel := filepath.Join(".agents", "skills", "bench-craft-line", "SKILL.md")
			writeGuidanceFile(t, root, rel, "Declare the line as `"+token+"`.\n")
			diags := checkGuidanceTokens(root)
			if len(diags) != 1 {
				t.Fatalf("one planted id, want exactly one diagnostic, got:\n%s", strings.Join(diags, "\n"))
			}
			for _, want := range []string{"guidance hard-codes an unbound model id", filepath.ToSlash(rel), "'" + token + "'", ".bench/lines.env"} {
				if !containsDiagnostic(diags, want) {
					t.Fatalf("want %q in the diagnostic, got:\n%s", want, diags[0])
				}
			}
		})
	}
}

// TestGuidanceSweepScansAFileAddedAfterTheCheck is the run-time-discovery row. The defect
// rides a file this check has never named, in a directory it has never named, nested one
// level deeper than the fixture's own. A hard-coded inventory or a single-level scan
// misses it. That is why production discovery owns the universal set.
func TestGuidanceSweepScansAFileAddedAfterTheCheck(t *testing.T) {
	for _, rel := range []string{
		filepath.Join(".agents", "commands", "bench-invented-later.md"),
		filepath.Join(".agents", "skills", "bench-craft-invented", "references", "deep.md"),
	} {
		t.Run(filepath.ToSlash(rel), func(t *testing.T) {
			root := writeGuidanceFixture(t)
			writeGuidanceFile(t, root, rel, "Route this at `beta-9.9-sol`.\n")
			diags := checkGuidanceTokens(root)
			if len(diags) != 1 || !containsDiagnostic(diags, filepath.ToSlash(rel)) {
				t.Fatalf("a newly added guidance file was not scanned, got:\n%s", strings.Join(diags, "\n"))
			}
		})
	}
}

// TestGuidanceSweepBitesEveryRetiredSchemaKey walks the retired schema one key at a time:
// the six concrete keys plus the two glob spellings prose uses for a key family. A check
// catching one spelling lets the rest rot. Each key is planted alone, and no sibling may
// fire.
func TestGuidanceSweepBitesEveryRetiredSchemaKey(t *testing.T) {
	var keys []string
	for _, prefix := range lines.RetiredKeyPrefixes() {
		for _, tier := range upperTiers() {
			keys = append(keys, prefix+tier)
		}
		keys = append(keys, prefix+"*")
	}
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			root := writeGuidanceFixture(t)
			rel := filepath.Join(".agents", "commands", "bench-setup-repo.md")
			writeGuidanceFile(t, root, rel, "Record the binding in `"+key+"`.\n")
			diags := checkGuidanceTokens(root)
			if len(diags) != 1 {
				t.Fatalf("one planted key, want exactly one diagnostic, got:\n%s", strings.Join(diags, "\n"))
			}
			for _, want := range []string{"guidance names a retired binding key", filepath.ToSlash(rel), "'" + key + "'"} {
				if !containsDiagnostic(diags, want) {
					t.Fatalf("want %q in the diagnostic, got:\n%s", want, diags[0])
				}
			}
			for _, sibling := range keys {
				if sibling != key && containsDiagnostic(diags, "'"+sibling+"'") {
					t.Fatalf("planting %s also fired for %s: %s", key, sibling, diags[0])
				}
			}
		})
	}
}

// TestGuidanceSweepRejectsNonRegularEntriesBeforeReading covers the four non-regular
// kinds a swept directory can hold. The FIFO case carries the load-bearing half: it has
// no writer. An implementation that opens before it classifies blocks in open(2) forever.
// It fails this row by expiring its deadline, rather than by returning a wrong answer. A
// character device needs mknod privilege a host may not grant, so that kind skips visibly
// rather than green.
func TestGuidanceSweepRejectsNonRegularEntriesBeforeReading(t *testing.T) {
	install := map[string]func(t *testing.T, path string){
		"fifo": func(t *testing.T, path string) {
			if err := syscall.Mkfifo(path, 0o644); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
			}
		},
		"socket": func(t *testing.T, path string) {
			listener, err := net.Listen("unix", path)
			if err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable on this filesystem: %v", err))
			}
			t.Cleanup(func() { listener.Close() })
		},
		"symlink": func(t *testing.T, path string) {
			if err := os.Symlink(filepath.Join(filepath.Dir(path), "SKILL.md"), path); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
			}
		},
		"character device": func(t *testing.T, path string) {
			// Major 1, minor 3 is the null device; mknod is privileged on most hosts.
			const nullDevice = 1<<8 | 3
			if err := syscall.Mknod(path, syscall.S_IFCHR|0o644, nullDevice); err != nil {
				capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot create a character device: %v", err))
			}
		},
	}
	for kind, plant := range install {
		t.Run(kind, func(t *testing.T) {
			root := writeGuidanceFixture(t)
			path := filepath.Join(root, ".agents", "skills", "bench-craft-line", "entry")
			plant(t, path)
			done := make(chan []string, 1)
			go func() { done <- checkGuidanceTokens(root) }()
			select {
			case diags := <-done:
				if len(diags) != 1 || !containsDiagnostic(diags, "guidance entry not a regular file") {
					t.Fatalf("a %s entry was not rejected, got:\n%s", kind, strings.Join(diags, "\n"))
				}
				if !containsDiagnostic(diags, ".agents/skills/bench-craft-line/entry") {
					t.Fatalf("the %s diagnostic does not name the entry: %s", kind, diags[0])
				}
			case <-time.After(bounds.TestDeadline(0)):
				t.Fatalf("the sweep blocked on a %s, so it opened the entry before classifying it", kind)
			}
		})
	}
}
