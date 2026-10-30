package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	"github.com/gibbonmi/bench/internal/worktree"
)

// readPages drives argv and every executable successor through the real dispatcher from cwd.
// Each response must stay within the response bound, so no page spills. It answers the
// decoded bytes and the page count.
func readPages(t *testing.T, cwd string, argv []string) ([]byte, int) {
	t.Helper()
	var rebuilt []byte
	for pages := 1; pages < 1000; pages++ {
		result := runAXICommandAt(t, cwd, argv)
		if result.code != 0 || result.stderr != "" || len(result.stdout) > bounds.ResponseBytes || strings.Count(result.stdout, "\n") > bounds.ResponseLines {
			t.Fatalf("%q = stdout=%q stderr=%q exit=%d, want one bounded page", argv, result.stdout, result.stderr, result.code)
		}
		document, err := axitest.DecodeDocument(result.stdout)
		if err != nil {
			t.Fatal(err)
		}
		rows, _ := document.Rows("page")
		page := rows[0].(map[string]any)
		contents, _ := document.Rows("content")
		for _, content := range contents {
			data, err := responsebound.DecodeContent(page["encoding"].(string), content.(map[string]any)["data"].(string))
			if err != nil {
				t.Fatal(err)
			}
			rebuilt = append(rebuilt, data...)
		}
		if page["complete"] == true {
			return rebuilt, pages
		}
		next, err := axitest.RecoverHelpCommandArgv(result.stdout)
		if err != nil || len(next) < 2 || next[0] != "bench" {
			t.Fatalf("nonterminal page %q has no executable successor: %v", result.stdout, err)
		}
		argv = next[1:]
	}
	t.Fatalf("%q did not terminate", argv)
	return nil, 0
}

// TestExactReaderCLIReconstruction reconstructs independently authored bytes through the real
// command and every continuation: repository files, a real spill file, and an exported
// evidence file. A changed source refuses its continuation before any content.
func TestExactReaderCLIReconstruction(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	files := map[string]string{
		"mixed.txt":   "alpha\r\nbeta\n" + strings.Repeat("one oversized physical line ", 400) + "\x00\x01\t\"\\ tail without newline",
		"unicode.txt": strings.Repeat("é😀ü", 1200),
		"invalid.bin": "\xff\xfe" + strings.Repeat("\x80abc", 1500),
		"empty.txt":   "",
	}
	for name, content := range files {
		writeAXIFixture(t, filepath.Join(root, name), content)
	}
	var spillResponse bytes.Buffer
	owner := responsebound.New(worktree.Home(), &spillResponse, &spillResponse, func() string { return "" }, false)
	_, _ = owner.Stdout().Write([]byte(strings.Repeat("spilled line\n", 400)))
	owner.Finish()
	spill, ok := responseboundtest.Find(spillResponse.String())
	if !ok {
		t.Fatalf("fixture response did not spill: %q", spillResponse.String())
	}
	exported := filepath.Join(t.TempDir(), "evidence", "source-1.txt")
	writeAXIFixture(t, exported, strings.Repeat("evidence ✓\r\n", 700))
	cases := map[string][]string{}
	for name := range files {
		cases[name] = []string{"read", name}
	}
	cases["spill"] = []string{"read", spill.Path, "--artifact"}
	cases["exported evidence"] = []string{"read", exported, "--artifact"}
	for name, argv := range cases {
		t.Run(name, func(t *testing.T) {
			want, ok := files[name]
			if !ok {
				data, err := os.ReadFile(argv[1])
				if err != nil {
					t.Fatal(err)
				}
				want = string(data)
			}
			got, pages := readPages(t, root, argv)
			if string(got) != want {
				t.Fatalf("reconstructed %d bytes over %d pages, want the %d original bytes", len(got), pages, len(want))
			}
			if len(want) > bounds.ResponseBytes && pages < 2 {
				t.Fatalf("%d bytes arrived on %d page", len(want), pages)
			}
		})
	}

	t.Run("changed source", func(t *testing.T) {
		first := runAXICommandAt(t, root, []string{"read", "mixed.txt"})
		next, err := axitest.RecoverHelpCommandArgv(first.stdout)
		if err != nil {
			t.Fatal(err)
		}
		writeAXIFixture(t, filepath.Join(root, "mixed.txt"), files["mixed.txt"]+"!")
		result := runAXICommandAt(t, root, next[1:])
		if result.code != 1 || !strings.HasPrefix(result.stdout, "error: source changed since the cursor") || strings.Contains(result.stdout, "content[") {
			t.Fatalf("changed continuation = stdout=%q exit=%d, want a refusal before content", result.stdout, result.code)
		}
	})
}

// TestPlanningReaderHelpAndRefusals grades the cursor grammar in bounded help, the usage and
// operational refusals, and the registered tree-scoped route of `bench read`.
func TestPlanningReaderHelpAndRefusals(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "page.txt"), strings.Repeat("page text\n", 600))
	help := runAXICommandAt(t, root, []string{"read", "--help"})
	if help.code != 0 || !strings.Contains(help.stdout, "usage: bench read <path>") || !strings.Contains(help.stdout, responsebound.CursorUsage) || len(help.stdout) > bounds.ResponseBytes {
		t.Fatalf("read --help = stdout=%q exit=%d, want the bounded cursor grammar", help.stdout, help.code)
	}
	first := runAXICommandAt(t, root, []string{"read", "page.txt"})
	document, err := axitest.DecodeDocument(first.stdout)
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := document.Rows("page")
	token := rows[0].(map[string]any)["source"].(string)
	for _, tc := range []struct {
		name   string
		argv   []string
		code   int
		prefix string
	}{
		{"offset without source", []string{"read", "page.txt", "--offset", "5"}, 2, "usage: bench read"},
		{"source without offset", []string{"read", "page.txt", "--source", token}, 2, "usage: bench read"},
		{"negative offset", []string{"read", "page.txt", "--offset", "-1", "--source", token}, 2, "usage: bench read"},
		{"nondecimal offset", []string{"read", "page.txt", "--offset", "1e3", "--source", token}, 2, "usage: bench read"},
		{"overflowing offset", []string{"read", "page.txt", "--offset", "99999999999999999999", "--source", token}, 2, "usage: bench read"},
		{"repeated offset", []string{"read", "page.txt", "--offset", "1", "--offset", "1", "--source", token}, 2, "usage: bench read"},
		{"lines with artifact", []string{"read", "/page.txt", "--lines", "1:2", "--artifact"}, 2, "usage: bench read"},
		{"malformed lines", []string{"read", "page.txt", "--lines", "2:1"}, 2, "usage: bench read"},
		{"relative artifact", []string{"read", "page.txt", "--artifact"}, 2, "usage: bench read"},
		{"missing file", []string{"read", "missing.txt"}, 1, "error: missing.txt is absent"},
		{"outside the repository", []string{"read", "../outside.txt"}, 1, "error: ../outside.txt is outside the repository"},
		{"unbound source", []string{"read", "page.txt", "--offset", "0", "--source", strings.Repeat("0", 64)}, 1, "error: source changed since the cursor"},
		{"cursor out of range", []string{"read", "page.txt", "--offset", "6000", "--source", token}, 1, "error: cursor out of range"},
		{"lines beyond the file", []string{"read", "page.txt", "--lines", "600:601"}, 1, "error: line range unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := runAXICommandAt(t, root, tc.argv)
			if result.code != tc.code || !strings.HasPrefix(result.stdout, tc.prefix) || strings.Contains(result.stdout, "content[") {
				t.Fatalf("%q = stdout=%q exit=%d, want exit %d beginning %q with no content", tc.argv, result.stdout, result.code, tc.code, tc.prefix)
			}
		})
	}

	t.Run("registered tree-scoped query", func(t *testing.T) {
		var stdout bytes.Buffer
		t.Chdir(root)
		if code := (Command{Stdout: &stdout}).Run([]string{"read", "page.txt", "--lines", "2:3"}); code != 0 || !strings.HasPrefix(stdout.String(), "tree[1]{target,head,dirty}:\n") {
			t.Fatalf("read = %q exit %d, want the dispatcher's tree row first", stdout.String(), code)
		}
		if !slices.Contains(axiRegistryMemberNames(t), "read") {
			t.Fatal("read is no AXI query member")
		}
		if got, _ := readPages(t, root, []string{"read", "page.txt", "--lines", "2:3"}); string(got) != "page text\npage text\n" {
			t.Fatalf("line range 2:3 = %q, want two physical lines with their terminators", got)
		}
	})
}
