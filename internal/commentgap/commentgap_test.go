package commentgap

import (
	"errors"
	"maps"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestProveGoGap(t *testing.T) {
	tests := []struct {
		name, path, before, after string
		want                      error
	}{
		{"CG9 trailing comment", "a.go", "package p; var x = 1 // before\n", "package p; var x = 1 // after\n", nil},
		{"CG10 final line", "a.go", "package p\n// before", "package p\n// after", nil},
		{"CG11 layout", "a.go", "package p\nvar x = 1\n", "package p\n\n  var   x = 1\n", nil},
		{"CG12 raw string", "a.go", "package p; var x = `// one`", "package p; var x = `// two`", ErrTokens},
		{"CG13 inserted semicolon", "a.go", "package p\nvar x = 1 /* one */ + 2\n", "package p\nvar x = 1 /* one\n */ + 2\n", ErrTokens},
		{"CG14 later scanner error", "a.go", "package p\n/* old */", "package p\n/* old", ErrScan},
		{"CG15 build directive", "a.go", "//go:build linux\n\npackage p\n", "//go:build darwin\n\npackage p\n", ErrDirective},
		{"CG16 added embed", "a.go", "package p\nvar s string\n", "package p\n//go:embed asset\nvar s string\n", ErrDirective},
		{"CG17 export directive", "a.go", "package p\n//export F\nfunc F() {}\n", "package p\n//export G\nfunc F() {}\n", ErrDirective},
		{"CG18 legacy build", "a.go", "// +build linux\n\npackage p\n", "// +build darwin\n\npackage p\n", ErrDirective},
		{"CG19 cgo preamble", "a.go", "package p\n/* #include <stdio.h> */\nimport \"C\"\n", "package p\n/* #include <stdlib.h> */\nimport \"C\"\n", ErrCgo},
		{"CG20 example output", "a_test.go", "package p\n// Output: one\n", "package p\n// Output: two\n", ErrExampleOutput},
		{"CG21 added output", "a_test.go", "package p\n// before\n", "package p\n// Output: after\n", ErrExampleOutput},
		{"CG42 reviewed scanner error", "a.go", "package p\n/* old", "package p\n/* old */", ErrScan},
		{"CG43 removed output", "a_test.go", "package p\n// Output: before\n", "package p\n// after\n", ErrExampleOutput},
		{"CG48 relocated embed", "a.go", "package p\nimport _ \"embed\"\n//go:embed payload.txt\nvar A string\nvar B string\n", "package p\nimport _ \"embed\"\nvar A string\n//go:embed payload.txt\nvar B string\n", ErrDirective},
		{"CG49 unchanged directive", "a.go", "//go:build linux\n\npackage p // before\n", "//go:build linux\n\npackage p // after\n", ErrDirective},
		{"CG50 multiline output", "a_test.go", "package p\nfunc Example() {\nprintln(1)\n/*\nOutput: old\n*/\n}\n", "package p\nfunc Example() {\nprintln(1)\n/*\nOutput: new\n*/\n}\n", ErrExampleOutput},
		{"CG50 added multiline output", "a_test.go", "package p\nfunc Example() {\nprintln(1)\n/*\nordinary\n*/\n}\n", "package p\nfunc Example() {\nprintln(1)\n/*\nUnordered output: new\n*/\n}\n", ErrExampleOutput},
		{"CG50 removed multiline output", "a_test.go", "package p\nfunc Example() {\nprintln(1)\n/*\nOutput: old\n*/\n}\n", "package p\nfunc Example() {\nprintln(1)\n/*\nordinary\n*/\n}\n", ErrExampleOutput},
		{"ordered directives", "a.go", "//go:a\n//go:b\npackage p\n", "//go:b\n//go:a\npackage p\n", ErrDirective},
		{"block line directive", "a.go", "package p\n/*line a.go:1*/ var x = 1\n", "package p\n/*line b.go:1*/ var x = 1\n", ErrDirective},
		{"unordered output", "a_test.go", "package p\n/* UnOrDeReD OuTpUt: before */\n", "package p\n/* UnOrDeReD OuTpUt: after */\n", ErrExampleOutput},
		{"import parse error", "a.go", "package p\nimport 1 // before\n", "package p\nimport 1 // after\n", ErrScan},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := proveGo(tc.path, []byte(tc.before), []byte(tc.after)); !errors.Is(got, tc.want) {
				t.Fatalf("Go gap = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestProveTreeGap(t *testing.T) {
	for _, name := range []string{"CG22 not Go", "CG23 added", "CG24 deleted", "CG25 renamed", "CG26 mode", "CG27 symlink", "CG28 gitlink", "CG29 oversized", "CG30 hidden gitlink", "CG31 hostile paths", "CG32 later code", "CG46 mixed hidden gitlink", "CG47 empty directory", "equal trees", "invalid tree", "executable Go"} {
		t.Run(name, func(t *testing.T) {
			f := newTrees(t)
			before := map[string]treeEntry{"a.go": {"100644", "package p // before\n"}}
			after := map[string]treeEntry{"a.go": {"100644", "package p // after\n"}}
			want, path := error(nil), "a.go"
			switch name {
			case "CG22 not Go":
				before, after = map[string]treeEntry{"notes.md": before["a.go"]}, map[string]treeEntry{"notes.md": after["a.go"]}
				want, path = ErrNotGo, "notes.md"
			case "CG23 added":
				before, after = nil, map[string]treeEntry{"a.go": {"100644", "// comment\n"}}
				want = ErrStatus
			case "CG24 deleted":
				before, after = map[string]treeEntry{"a.go": {"100644", "// comment\n"}}, nil
				want = ErrStatus
			case "CG25 renamed":
				after = map[string]treeEntry{"z.go": after["a.go"]}
				want = ErrStatus
			case "CG26 mode":
				after["a.go"] = treeEntry{"100755", before["a.go"].body}
				want = ErrMode
			case "CG27 symlink":
				before["a.go"], after["a.go"] = treeEntry{"120000", "before"}, treeEntry{"120000", "after"}
				want = ErrMode
			case "CG28 gitlink", "CG30 hidden gitlink", "CG46 mixed hidden gitlink":
				first, second := f.commits()
				if name != "CG46 mixed hidden gitlink" {
					before, after = map[string]treeEntry{}, map[string]treeEntry{}
				}
				before["kit.go"], after["kit.go"] = treeEntry{"160000", first}, treeEntry{"160000", second}
				want, path = ErrMode, "kit.go"
				if name != "CG28 gitlink" {
					gittest.Output(t, f.root, "config", "diff.ignoreSubmodules", "all")
				}
			case "CG29 oversized":
				body := "package p\n// " + strings.Repeat("x", int(bounds.ControlRecordLimit))
				before["a.go"], after["a.go"] = treeEntry{"100644", body + "before"}, treeEntry{"100644", body + "after"}
				want = ErrUnreadable
			case "CG31 hostile paths":
				left, right := before["a.go"], after["a.go"]
				before, after = map[string]treeEntry{}, map[string]treeEntry{}
				for _, path := range []string{"a b*.go", "tab\t.go"} {
					before[path], after[path] = left, right
				}
			case "CG32 later code":
				before["b.go"], after["b.go"] = before["a.go"], after["a.go"]
				before["c.go"], after["c.go"] = treeEntry{"100644", "package p; var x = 1\n"}, treeEntry{"100644", "package p; var x = 2\n"}
				want, path = ErrTokens, "c.go"
			case "equal trees":
				after = maps.Clone(before)
			case "executable Go":
				before["a.go"], after["a.go"] = treeEntry{"100755", before["a.go"].body}, treeEntry{"100755", after["a.go"].body}
			}
			left, right := f.tree(before), f.tree(after)
			if name == "CG47 empty directory" {
				left, right = f.tree(nil), f.emptyDirectoryTree()
				want, path = ErrEmptyChanges, ""
			}
			if name == "invalid tree" {
				left, want, path = strings.Repeat("0", 40), ErrUnreadable, ""
			}
			got := Prove(f.root, left, right)
			if !errors.Is(got, want) {
				t.Fatalf("tree gap = %v, want %v", got, want)
			}
			if want != nil && path != "" && !strings.Contains(got.Error(), path) {
				t.Fatalf("refusal lost path %q: %v", path, got)
			}
		})
	}
}
