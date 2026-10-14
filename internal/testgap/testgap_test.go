package testgap

import (
	"maps"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

func TestProveTestGap(t *testing.T) {
	for _, name := range []string{"RC13 testdata", "RC13 root testdata", "RC13 testdata in a name", "RC14 deleted", "RC15 symlink", "RC16 mode change", "RC17 not Go", "RC17 not test", "RC17 later production file", "RC18 empty directory", "RC19 spaced path", "RC52 gitlink", "RC52 mixed hidden gitlink", "RC53 renamed", "added test", "executable test", "modified test", "equal trees", "invalid tree"} {
		t.Run(name, func(t *testing.T) {
			f := gittest.NewTrees(t)
			before := map[string]gittest.TreeEntry{"a_test.go": {Mode: "100644", Body: "package p // before\n"}}
			after := map[string]gittest.TreeEntry{"a_test.go": {Mode: "100644", Body: "package p // after\n"}}
			want, path := error(nil), "a_test.go"
			move := func(to string) {
				before, after = map[string]gittest.TreeEntry{to: before[path]}, map[string]gittest.TreeEntry{to: after[path]}
				path = to
			}
			switch name {
			case "RC13 testdata":
				move("p/testdata/a_test.go")
				want = ErrTestdata
			case "RC13 root testdata":
				move("testdata/a_test.go")
				want = ErrTestdata
			case "RC13 testdata in a name":
				move("a/testdata_x/b_test.go")
			case "RC14 deleted":
				after = nil
				want = ErrStatus
			case "RC15 symlink":
				before, after = nil, map[string]gittest.TreeEntry{path: {Mode: "120000", Body: "target_test.go"}}
				want = ErrMode
			case "RC16 mode change":
				after[path] = gittest.TreeEntry{Mode: "100755", Body: before[path].Body}
				want = ErrMode
			case "RC17 not Go":
				move("notes.md")
				want = ErrNotTest
			case "RC17 not test":
				move("a.go")
				want = ErrNotTest
			case "RC17 later production file":
				// The production path sorts after the accepted test path.
				before["z.go"], after["z.go"] = gittest.TreeEntry{Mode: "100644", Body: "package p\n"}, gittest.TreeEntry{Mode: "100644", Body: "package p // after\n"}
				want, path = ErrNotTest, "z.go"
			case "RC19 spaced path":
				move("a b_test.go")
			case "RC52 gitlink", "RC52 mixed hidden gitlink":
				first, second := f.Commits()
				gitlink := path
				if name == "RC52 mixed hidden gitlink" {
					gitlink = "kit_test.go"
					gittest.Output(t, f.Root, "config", "diff.ignoreSubmodules", "all")
				}
				before[gitlink], after[gitlink] = gittest.TreeEntry{Mode: "160000", Body: first}, gittest.TreeEntry{Mode: "160000", Body: second}
				want, path = ErrMode, gitlink
			case "RC53 renamed":
				after = map[string]gittest.TreeEntry{"z_test.go": after[path]}
				want = ErrStatus
			case "added test":
				before = nil
			case "executable test":
				before[path], after[path] = gittest.TreeEntry{Mode: "100755", Body: before[path].Body}, gittest.TreeEntry{Mode: "100755", Body: after[path].Body}
			case "equal trees":
				after = maps.Clone(before)
			}
			left, right := f.Tree(before), f.Tree(after)
			if name == "RC18 empty directory" {
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
