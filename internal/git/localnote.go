package git

import "path/filepath"

// LocalNoteRoot resolves which checkout's copy of relPath a note-writing verb
// must use. An ignored path is a local working note: it can never dirty a
// landing destination, and a worktree-local copy dies with the worktree, so
// the primary checkout's copy is the one that persists. A path git does not
// ignore keeps the caller's own checkout and its checkout policy.
func LocalNoteRoot(root, relPath string) (noteRoot string, ignored bool, err error) {
	if _, checkErr := Output("-C", root, "check-ignore", "-q", relPath); checkErr != nil {
		return root, false, nil
	}
	common, err := CommonDir(root)
	if err != nil {
		return "", true, err
	}
	return filepath.Dir(common), true, nil
}
