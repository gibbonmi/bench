package preflight

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// The expectations in this file are written from the spec's response tables and cursor
// grammar, not read from the format registry.

func sha(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

// prepareEvidence runs the build preparation and returns its identity and decoded row.
func prepareEvidence(t *testing.T, args []string) (string, map[string]any, string) {
	t.Helper()
	out, code := Command(args)
	if code != 0 {
		t.Fatalf("prepare = (%d):\n%s", code, out)
	}
	rows := preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "prepared")
	if len(rows) != 1 {
		t.Fatalf("prepared rows = %d", len(rows))
	}
	return rows[0]["evidence"].(string), rows[0], out
}

// evidencePage is one decoded page response and its raw bytes.
type evidencePage struct {
	row map[string]any
	raw string
}

// traverseEvidence follows every next command from the first read to the stream end.
func traverseEvidence(t *testing.T, identity string) []evidencePage {
	t.Helper()
	args := []string{"evidence", identity}
	var pages []evidencePage
	for len(pages) < 10000 {
		out, code := Command(args)
		if code != 0 {
			t.Fatalf("read %v = (%d):\n%s", args, code, out)
		}
		rows := preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "page")
		if len(rows) != 1 {
			t.Fatalf("page rows = %d:\n%s", len(rows), out)
		}
		pages = append(pages, evidencePage{rows[0], out})
		next := rows[0]["next"].(string)
		if next == "" {
			return pages
		}
		fields := strings.Fields(next)
		if len(fields) < 4 || strings.Join(fields[:2], " ") != "bench preflight" {
			t.Fatalf("next = %q is not a preflight command", next)
		}
		args = fields[2:]
	}
	t.Fatal("traversal did not end")
	return nil
}

// reconstructEvidence joins the traversed fragments into the manifest and every source.
func reconstructEvidence(t *testing.T, identity string, pages []evidencePage) (string, map[string]string) {
	t.Helper()
	var manifest strings.Builder
	sources := map[string]string{}
	for _, page := range pages {
		content := page.row["content"].(string)
		if sha(content) != page.row["sha256"] {
			t.Fatalf("page digest differs: %v", page.row)
		}
		if page.row["stream"] == "manifest" {
			manifest.WriteString(content)
			continue
		}
		sources[page.row["source"].(string)] += content
	}
	if "sha256:"+sha(manifest.String()) != identity {
		t.Fatalf("reconstructed manifest does not match identity %s", identity)
	}
	return manifest.String(), sources
}

func manifestSources(t *testing.T, manifest string) []map[string]any {
	t.Helper()
	return preflighttest.TableRows(t, preflighttest.DecodeMap(t, manifest), "sources")
}

// TestEvidenceBuildRoundTrip is CE1: every canonical source reconstructs byte for byte.
func TestEvidenceBuildRoundTrip(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	args := preflighttest.ChargeArgs(t, root, slug, false)
	identity, _, _ := prepareEvidence(t, args)
	manifest, sources := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	rows := manifestSources(t, manifest)
	wantPaths := []string{"", "specs/example/tickets/one.md", "specs/example/spec.md", chargesource.DelegateSkill, chargesource.BuildPhase, chargesource.DelegateProcedure}
	if len(rows) != len(wantPaths) {
		t.Fatalf("sources = %d, want %d", len(rows), len(wantPaths))
	}
	for i, row := range rows {
		if row["path"] != wantPaths[i] || sha(sources[row["id"].(string)]) != row["sha256"] {
			t.Fatalf("source %d = %v, reconstructed %d bytes", i, row, len(sources[row["id"].(string)]))
		}
		pinned, err := git.Raw("-C", root, "show", args[8]+":"+wantPaths[i])
		if i > 0 && (err != nil || sources[row["id"].(string)] != string(pinned)) {
			t.Fatalf("source %s differs from the pinned tip", wantPaths[i])
		}
	}
}

// TestEvidenceLargeTicket is CE2.
func TestEvidenceLargeTicket(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	ticket := preflighttest.TicketDoc("One", "PF1", "PF2") + strings.Repeat("large ticket evidence\n", 3000)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
	identity, _, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "large ticket", false))
	_, sources := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	if len(ticket) <= chargeevidence.ResponseLimit || sources["s2"] != ticket {
		t.Fatalf("large ticket reconstructed %d of %d bytes", len(sources["s2"]), len(ticket))
	}
}

// TestEvidencePreparationState is CE3 and TestEvidencePageState is CE4: neither response
// claims delivery.
func TestEvidencePreparationState(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	_, row, out := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	if row["delivery"] != "unverified" || row["response_complete"] != true || strings.Contains(out, "delivery_complete") {
		t.Fatalf("prepared state = %v", row)
	}
}

func TestEvidencePageState(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	pages := traverseEvidence(t, identity)
	for i, page := range pages {
		last := i == len(pages)-1
		if page.row["response_complete"] != true || page.row["stream_end"] != last || strings.Contains(page.raw, "delivery") {
			t.Fatalf("page %d state = %v", i, page.row)
		}
	}
}

// TestEvidenceTraversalOrder is CE5: the whole manifest precedes every source page, and
// sources follow in manifest order with increasing page indexes.
func TestEvidenceTraversalOrder(t *testing.T) {
	t.Run("build", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.TicketDoc("One", "PF1", "PF2")+strings.Repeat("x", 20000))
		identity, _, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "paged ticket", false))
		assertStreamOrder(t, identity, 1, 6)
	})
	t.Run("multi-fragment manifest", func(t *testing.T) {
		root, _ := preflighttest.SeedConformant(t)
		assertStreamOrder(t, publishCraftedEvidence(t, root, strings.Repeat("manifest scalar ", 2000)), 3, 2)
	})
}

// assertStreamOrder walks the default stream, reconstructs the manifest, and checks that
// at least fragments manifest fragments precede every source page and that the stream ends
// at source lastSource.
func assertStreamOrder(t *testing.T, identity string, fragments, lastSource int) {
	t.Helper()
	pages := traverseEvidence(t, identity)
	reconstructEvidence(t, identity, pages)
	lastOrdinal, lastIndex, manifestPages, sourcePages := 0, -1, 0, 0
	for i, page := range pages {
		ordinal := 0
		if page.row["stream"] == "source" {
			ordinal, _ = strconv.Atoi(strings.TrimPrefix(page.row["source"].(string), "s"))
			sourcePages++
		} else {
			manifestPages++
		}
		index := int(page.row["index"].(float64))
		switch {
		case i == 0 && (ordinal != 0 || index != 0):
			t.Fatalf("first page = %v, want manifest fragment 0", page.row)
		case ordinal == lastOrdinal && index != lastIndex+1:
			t.Fatalf("page %d index %d follows %d in stream %d", i, index, lastIndex, ordinal)
		case ordinal != lastOrdinal && (ordinal < lastOrdinal || index != 0):
			t.Fatalf("page %d opens stream %d at %d after stream %d", i, ordinal, index, lastOrdinal)
		}
		lastOrdinal, lastIndex = ordinal, index
	}
	if manifestPages < fragments || sourcePages == 0 || lastOrdinal != lastSource {
		t.Fatalf("stream held %d manifest and %d source pages and ended at source %d", manifestPages, sourcePages, lastOrdinal)
	}
}

// TestEvidenceExactText is CE9 and CE10.
func TestEvidenceExactText(t *testing.T) {
	for _, test := range []struct{ name, tail string }{
		{"CE9 no final newline", "\n\nlast line without newline"},
		{"CE10 unicode and escapes", "\n\nRésumé 雪 🚀\t\"quoted\" \\ back\r\nend" + strings.Repeat("雪", 3000)},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := preflighttest.SeedConformant(t)
			ticket := strings.TrimSuffix(preflighttest.TicketDoc("One", "PF1", "PF2"), "\n") + test.tail
			preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
			identity, _, _ := prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, test.name, false))
			if _, sources := reconstructEvidence(t, identity, traverseEvidence(t, identity)); sources["s2"] != ticket {
				t.Fatalf("ticket reconstructed as %q", sources["s2"])
			}
		})
	}
}

// TestEvidenceHostilePaths is CE12: spaces, globs, quotes, substitutions, and Unicode in
// a ticket path stay data.
func TestEvidenceHostilePaths(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	name := "one [*] $(touch sentinel) 'q' 雪.md"
	rel := "specs/" + slug + "/tickets/dir /" + name
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug, "- `specs/"+slug+"/` (ticket fixtures)"))
	preflighttest.MustWriteFile(t, rel, preflighttest.TicketDoc("Hostile path", "PF1", "PF2"))
	args := preflighttest.LegacyCommitted(t, root, slug, "hostile ticket path", false)
	args[4] = name
	identity, _, _ := prepareEvidence(t, args)
	manifest, sources := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	if rows := manifestSources(t, manifest); rows[1]["path"] != rel || sources["s2"] != preflighttest.TicketDoc("Hostile path", "PF1", "PF2") {
		t.Fatalf("hostile path source = %v", rows[1])
	}
	if _, err := os.Stat(filepath.Join(root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("hostile path created sentinel: %v", err)
	}
}

// TestEvidenceMovementPublication is CE53 and CE54. The pause hook runs after each
// attempt's read, when that attempt's candidate is already staged.
func TestEvidenceMovementPublication(t *testing.T) {
	t.Run("CE53 one movement", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		args := preflighttest.ChargeArgs(t, root, slug, false)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			if temps := preflighttest.StagedTemps(t, root); len(temps) != 1 || len(preflighttest.PublishedPacks(t, root)) != 0 {
				t.Errorf("attempt %d sees temps %v and packs %v", calls, temps, preflighttest.PublishedPacks(t, root))
			}
			if calls == 1 {
				// A refreshed stat moves the index digest and leaves the checkout clean.
				moved := time.Now().Add(time.Hour)
				if err := os.Chtimes("internal/example/foo.go", moved, moved); err != nil {
					t.Fatal(err)
				}
				preflighttest.RunGit(t, "update-index", "--refresh")
			}
		})
		out, code := Command(args)
		restore()
		if code != 0 || calls != 2 || len(preflighttest.PublishedPacks(t, root)) != 1 || len(preflighttest.StagedTemps(t, root)) != 0 {
			t.Fatalf("one movement = (%d, %d, %v):\n%s", code, calls, preflighttest.PublishedPacks(t, root), out)
		}
	})
	t.Run("CE54 two movements", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		args := preflighttest.ChargeArgs(t, root, slug, false)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n// moved "+string(rune('0'+calls))+"\n")
		})
		out, code := Command(args)
		restore()
		if code != 1 || calls != 2 || !strings.Contains(out, "snapshot drift") || len(preflighttest.PublishedPacks(t, root)) != 0 || len(preflighttest.StagedTemps(t, root)) != 0 {
			t.Fatalf("two movements = (%d, %d, %v):\n%s", code, calls, preflighttest.PublishedPacks(t, root), out)
		}
	})
}

// TestEvidencePreparationNext is CE156: the only next action is the manifest-first read.
func TestEvidencePreparationNext(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, row, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	if want := "bench preflight evidence " + identity; row["next"] != want || !strings.HasPrefix(identity, "sha256:") {
		t.Fatalf("prepared next = %q, want %q", row["next"], want)
	}
}

// TestEvidenceNextActions is CE64: every nonterminal page names the exact next cursor
// command, and the terminal page names none.
func TestEvidenceNextActions(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug, false))
	pages := traverseEvidence(t, identity)
	hex := strings.TrimPrefix(identity, "sha256:")
	manifestFragments := 0
	for i, page := range pages {
		if page.row["stream"] == "manifest" {
			manifestFragments++
		}
		if i == len(pages)-1 {
			if page.row["next"] != "" {
				t.Fatalf("terminal next = %q", page.row["next"])
			}
			continue
		}
		following := pages[i+1].row
		kind, ordinal := "m", "0"
		if following["stream"] == "source" {
			kind, ordinal = "s", strings.TrimPrefix(following["source"].(string), "s")
		}
		want := "bench preflight evidence " + identity + " --cursor v1." + hex + "." + kind + "." + ordinal + "." + formatIndex(following["index"])
		if page.row["next"] != want {
			t.Fatalf("page %d next = %q, want %q", i, page.row["next"], want)
		}
	}
	if manifestFragments == 0 {
		t.Fatal("no manifest fragment")
	}
}

func formatIndex(value any) string {
	number, _ := value.(float64)
	return strconv.Itoa(int(number))
}
