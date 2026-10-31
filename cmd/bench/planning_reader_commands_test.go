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
		if result.code != 0 || result.stderr != "" || !responsebound.Fits(result.stdout) {
			t.Fatalf("%q = stdout=%q stderr=%q exit=%d, want one bounded page", argv, result.stdout, result.stderr, result.code)
		}
		document, err := axitest.DecodeDocument(result.stdout)
		if err != nil {
			t.Fatal(err)
		}
		rows, err := document.Rows("page")
		if err != nil {
			t.Fatal(err)
		}
		page := rows[0].(map[string]any)
		contents, err := document.Rows("content")
		if err != nil {
			t.Fatal(err)
		}
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
// evidence file. A call from a subdirectory follows every continuation to the same file. A
// changed source refuses its continuation before any content.
func TestExactReaderCLIReconstruction(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	files := map[string]string{
		"mixed.txt":       "alpha\r\nbeta\n" + strings.Repeat("one oversized physical line ", 400) + "\x00\x01\t\"\\ tail without newline",
		"unicode.txt":     strings.Repeat("é😀ü", 1200),
		"invalid.bin":     "\xff\xfe" + strings.Repeat("\x80abc", 1500),
		"empty.txt":       "",
		"spaced name.txt": strings.Repeat("a name with a space\n", 400),
		"nested/deep.txt": strings.Repeat("nested line\n", 700),
	}
	// A long section holds a fenced heading, so its pages must cross the hidden heading and
	// stop before the next visible one.
	const (
		before = "# Guide\n\n## Before\nbefore\n"
		after  = "## After\nafter\n"
	)
	long := "## Long\r\n" + strings.Repeat("section text ✓\n", 300) + "```text\n## Hidden\n```\n" + strings.Repeat("more text\n", 200)
	stories := "## User stories\n\nLine: opus / medium.\n\n### Group\n\n" + strings.Repeat("1. As an agent, I want a long story.\n", 150) + "\n"
	files["guide.md"] = before + long + after
	files["specs/long/spec.md"] = "# Long spec\n\nStatus: staged\n\n" + stories + "## Implementation decisions\n\nnone\n"
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
	type readCase struct {
		cwd, file string
		argv      []string
		// section is the selected text when the case selects less than the whole file.
		section string
	}
	nested := filepath.Join(root, "nested")
	cases := map[string]readCase{}
	for name := range files {
		cases[name] = readCase{root, name, []string{"read", name}, ""}
	}
	cases["from a subdirectory"] = readCase{nested, "nested/deep.txt", []string{"read", "deep.txt"}, ""}
	cases["parent from a subdirectory"] = readCase{nested, "mixed.txt", []string{"read", "../mixed.txt"}, ""}
	cases["spill"] = readCase{root, "", []string{"read", spill.Path, "--artifact"}, ""}
	cases["exported evidence"] = readCase{root, "", []string{"read", exported, "--artifact"}, ""}
	cases["section"] = readCase{root, "", []string{"read", "guide.md", "--section", "Long"}, long}
	cases["section from a subdirectory"] = readCase{nested, "", []string{"read", "../guide.md", "--section", "Long"}, long}
	cases["spec show stories"] = readCase{nested, "", []string{"spec", "show", "long", "--section", "stories"}, stories}
	cases["spec show literal section"] = readCase{root, "", []string{"spec", "show", "specs/long/spec.md", "--section", "User stories"}, stories}
	cases["spec show file"] = readCase{root, "specs/long/spec.md", []string{"spec", "show", "long"}, ""}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			want, ok := files[c.file]
			if c.section != "" {
				want, ok = c.section, true
			}
			if !ok {
				data, err := os.ReadFile(c.argv[1])
				if err != nil {
					t.Fatal(err)
				}
				want = string(data)
			}
			got, pages := readPages(t, c.cwd, c.argv)
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
// operational refusals, and the registered repository-scoped route of `bench read`, whose
// grammar takes no tree target.
func TestPlanningReaderHelpAndRefusals(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "page.txt"), strings.Repeat("page text\n", 600))
	writeAXIFixture(t, filepath.Join(root, "control\x01.txt"), strings.Repeat("page text\n", 600))
	writeAXIFixture(t, filepath.Join(root, "dup.md"), "## Dup\na\n## Dup\nb\n")
	writeAXIFixture(t, filepath.Join(root, "specs", "plain", "spec.md"), "# Plain\n\n## Notes\n")
	for _, argv := range [][]string{{"read", "--help"}, {"spec", "show", "--help"}} {
		help := runAXICommandAt(t, root, argv)
		if help.code != 0 || !strings.HasPrefix(help.stdout, "usage: bench "+argv[0]) || !strings.Contains(help.stdout, "--section") || !strings.Contains(help.stdout, responsebound.CursorUsage) || !responsebound.Fits(help.stdout) {
			t.Fatalf("%q = stdout=%q exit=%d, want the bounded cursor grammar", argv, help.stdout, help.code)
		}
	}
	for _, absence := range []struct {
		argv []string
		want string
	}{
		{[]string{"read", "dup.md", "--section", "Missing"}, "selection[1]{kind,selector,present}:\n  section,Missing,false\nhelp[0]{cmd,why}:\n"},
		{[]string{"spec", "show", "plain", "--section", "stories"}, "selection[1]{kind,selector,present}:\n  stories,stories,false\nhelp[0]{cmd,why}:\n"},
		{[]string{"spec", "show", "plain", "--section", "status"}, "selection[1]{kind,selector,present}:\n  status,status,false\nhelp[0]{cmd,why}:\n"},
	} {
		if result := runAXICommandAt(t, root, absence.argv); result.code != 0 || result.stdout != absence.want {
			t.Fatalf("%q = stdout=%q exit=%d, want the explicit absence %q", absence.argv, result.stdout, result.code, absence.want)
		}
	}
	sourceOf := func(argv []string) string {
		t.Helper()
		document, err := axitest.DecodeDocument(runAXICommandAt(t, root, argv).stdout)
		if err != nil {
			t.Fatal(err)
		}
		rows, err := document.Rows("page")
		if err != nil {
			t.Fatal(err)
		}
		return rows[0].(map[string]any)["source"].(string)
	}
	token := sourceOf([]string{"read", "page.txt"})
	linesToken := sourceOf([]string{"read", "page.txt", "--lines", "1:600"})
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
		{"overflowing lines", []string{"read", "page.txt", "--lines", "1:99999999999999999999"}, 2, "usage: bench read"},
		{"tree target", []string{"read", "--in", "primary", "page.txt"}, 2, "usage: bench read"},
		{"token from another selection", []string{"read", "page.txt", "--lines", "1:600", "--offset", "0", "--source", token}, 1, "error: source changed since the cursor"},
		{"token from another line range", []string{"read", "page.txt", "--lines", "2:600", "--offset", "0", "--source", linesToken}, 1, "error: source changed since the cursor"},
		{"control byte in a continued path", []string{"read", "control\x01.txt"}, 1, "error: continuation cannot be rendered"},
		{"relative artifact", []string{"read", "page.txt", "--artifact"}, 2, "usage: bench read"},
		{"missing file", []string{"read", "missing.txt"}, 1, "error: missing.txt is absent"},
		{"outside the repository", []string{"read", "../outside.txt"}, 1, "error: ../outside.txt is outside the repository"},
		{"unbound source", []string{"read", "page.txt", "--offset", "0", "--source", strings.Repeat("0", 64)}, 1, "error: source changed since the cursor"},
		{"cursor out of range", []string{"read", "page.txt", "--offset", "6000", "--source", token}, 1, "error: cursor out of range"},
		{"lines beyond the file", []string{"read", "page.txt", "--lines", "600:601"}, 1, "error: line range unavailable"},
		{"section with lines", []string{"read", "dup.md", "--section", "Dup", "--lines", "1:2"}, 2, "usage: bench read"},
		{"section with artifact", []string{"read", "/dup.md", "--section", "Dup", "--artifact"}, 2, "usage: bench read"},
		{"repeated heading", []string{"read", "dup.md", "--section", "Dup"}, 1, "error: section unavailable"},
		{"spec show offset without source", []string{"spec", "show", "plain", "--offset", "0"}, 2, "usage: bench spec show"},
		{"spec show missing spec", []string{"spec", "show", "missing"}, 1, "error: spec not found"},
		{"spec show unbound source", []string{"spec", "show", "plain", "--offset", "0", "--source", token}, 1, "error: source changed since the cursor"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := runAXICommandAt(t, root, tc.argv)
			if result.code != tc.code || !strings.HasPrefix(result.stdout, tc.prefix) || strings.Contains(result.stdout, "content[") {
				t.Fatalf("%q = stdout=%q exit=%d, want exit %d beginning %q with no content", tc.argv, result.stdout, result.code, tc.code, tc.prefix)
			}
		})
	}

	t.Run("registered repository-scoped query", func(t *testing.T) {
		var stdout bytes.Buffer
		t.Chdir(root)
		if code := (Command{Stdout: &stdout}).Run([]string{"read", "page.txt", "--lines", "2:3"}); code != 0 || !strings.HasPrefix(stdout.String(), "page[1]{") {
			t.Fatalf("read = %q exit %d, want the page block first and no tree row", stdout.String(), code)
		}
		if !slices.Contains(axiRegistryMemberNames(t), "read") {
			t.Fatal("read is no AXI query member")
		}
		if got, _ := readPages(t, root, []string{"read", "page.txt", "--lines", "2:3"}); string(got) != "page text\npage text\n" {
			t.Fatalf("line range 2:3 = %q, want two physical lines with their terminators", got)
		}
	})
}
