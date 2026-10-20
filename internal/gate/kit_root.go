package gate

import "github.com/gibbonmi/bench/internal/git"

// KitForRoot answers the kit that grades root. It is the Go form of the kit_dir rule in
// bin/bench.sh, and the wrapper differential test holds the two forms to one answer. The
// answer is root only when kit is the top level of its own checkout and shares its git
// common directory with root, so a worktree of the kit's repository grades under its own
// tree. An adopted kit at <repo>/.bench is not a top level, and a kit from another
// repository has another common directory, so each keeps the caller's kit. A git failure
// or a path with no canonical spelling also keeps the caller's kit.
func KitForRoot(kit, root string) string {
	if kit == "" || !isTopLevel(kit) {
		return kit
	}
	if sameCommonDir(kit, root) {
		return root
	}
	return kit
}

// isTopLevel reports whether dir is the top level of its git checkout under canonical
// spelling. A git failure or a path with no canonical spelling answers false.
func isTopLevel(dir string) bool {
	top, err := git.RootAt(dir)
	if err != nil {
		return false
	}
	return sameResolvedPath(top, dir)
}

// sameCommonDir reports whether a and b share one git common directory under canonical
// spelling. A git failure or a path with no canonical spelling answers false.
func sameCommonDir(a, b string) bool {
	commonA, err := git.CommonDir(a)
	if err != nil {
		return false
	}
	commonB, err := git.CommonDir(b)
	if err != nil {
		return false
	}
	return sameResolvedPath(commonA, commonB)
}
