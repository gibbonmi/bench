package commentgap

import (
	"maps"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestProveTreeGap(t *testing.T) {
	for _, name := range []string{"CG22 not Go", "CG23 added", "CG24 deleted", "CG25 renamed", "CG26 mode", "CG27 symlink", "CG28 gitlink", "CG29 oversized", "CG30 hidden gitlink", "CG31 hostile paths", "CG32 later code", "CG46 mixed hidden gitlink", "CG47 empty directory", "equal trees", "invalid tree", "executable Go"} {
		t.Run(name, func(t *testing.T) {
			f := gittest.NewTrees(t)
			before := map[string]gittest.TreeEntry{"a.go": {Mode: "100644", Body: "package p // before\n"}}
			after := map[string]gittest.TreeEntry{"a.go": {Mode: "100644", Body: "package p // after\n"}}
			want, path := error(nil), "a.go"
			switch name {
			case "CG22 not Go":
				before, after = map[string]gittest.TreeEntry{"notes.md": before["a.go"]}, map[string]gittest.TreeEntry{"notes.md": after["a.go"]}
				want, path = ErrNotGo, "notes.md"
			case "CG23 added":
				before, after = nil, map[string]gittest.TreeEntry{"a.go": {Mode: "100644", Body: "// comment\n"}}
				want = ErrStatus
			case "CG24 deleted":
				before, after = map[string]gittest.TreeEntry{"a.go": {Mode: "100644", Body: "// comment\n"}}, nil
				want = ErrStatus
			case "CG25 renamed":
				after = map[string]gittest.TreeEntry{"z.go": after["a.go"]}
				want = ErrStatus
			case "CG26 mode":
				after["a.go"] = gittest.TreeEntry{Mode: "100755", Body: before["a.go"].Body}
				want = ErrMode
			case "CG27 symlink":
				before["a.go"], after["a.go"] = gittest.TreeEntry{Mode: "120000", Body: "before"}, gittest.TreeEntry{Mode: "120000", Body: "after"}
				want = ErrMode
			case "CG28 gitlink", "CG30 hidden gitlink", "CG46 mixed hidden gitlink":
				first, second := f.Commits()
				if name != "CG46 mixed hidden gitlink" {
					before, after = map[string]gittest.TreeEntry{}, map[string]gittest.TreeEntry{}
				}
				before["kit.go"], after["kit.go"] = gittest.TreeEntry{Mode: "160000", Body: first}, gittest.TreeEntry{Mode: "160000", Body: second}
				want, path = ErrMode, "kit.go"
				if name != "CG28 gitlink" {
					gittest.Output(t, f.Root, "config", "diff.ignoreSubmodules", "all")
				}
			case "CG29 oversized":
				body := "package p\n// " + strings.Repeat("x", int(bounds.ControlRecordLimit))
				before["a.go"], after["a.go"] = gittest.TreeEntry{Mode: "100644", Body: body + "before"}, gittest.TreeEntry{Mode: "100644", Body: body + "after"}
				want = ErrUnreadable
			case "CG31 hostile paths":
				left, right := before["a.go"], after["a.go"]
				before, after = map[string]gittest.TreeEntry{}, map[string]gittest.TreeEntry{}
				for _, path := range []string{"a b*.go", "tab\t.go"} {
					before[path], after[path] = left, right
				}
			case "CG32 later code":
				before["b.go"], after["b.go"] = before["a.go"], after["a.go"]
				before["c.go"], after["c.go"] = gittest.TreeEntry{Mode: "100644", Body: "package p; var x = 1\n"}, gittest.TreeEntry{Mode: "100644", Body: "package p; var x = 2\n"}
				want, path = ErrTokens, "c.go"
			case "equal trees":
				after = maps.Clone(before)
			case "executable Go":
				before["a.go"], after["a.go"] = gittest.TreeEntry{Mode: "100755", Body: before["a.go"].Body}, gittest.TreeEntry{Mode: "100755", Body: after["a.go"].Body}
			}
			left, right := f.Tree(before), f.Tree(after)
			if name == "CG47 empty directory" {
				left, right = f.Tree(nil), f.EmptyDirectoryTree()
				want, path = ErrEmptyChanges, ""
			}
			if name == "invalid tree" {
				left, want, path = strings.Repeat("0", 40), ErrUnreadable, ""
			}
			gittest.RequireVerdict(t, Prove(f.Root, left, right), want, path)
		})
	}
}
