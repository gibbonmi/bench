package landing

import (
	"path/filepath"
	"strings"
)

// BuildLogDir is the git-ignored directory of the orchestrator's build log.
const BuildLogDir = "capture/build-log/"

// LocalCapturePath reports whether path is one exact project-local capture file.
// It also accepts a safe spelling within the build-log directory.
func LocalCapturePath(path string) bool {
	exact := path == "capture/IDEAS.md" || path == "capture/learnings.md" || path == "capture/session-handoff.md"
	return exact || withinIgnoredDir(path, BuildLogDir)
}

// withinIgnoredDir reports whether path is a safe spelling of dir or of a path beneath it.
// The dir argument ends in a slash. The function refuses a control byte, an absolute
// path, a `..` segment, a doubled slash, and a sibling name that only shares the prefix.
func withinIgnoredDir(path, dir string) bool {
	for _, r := range path {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	if path == dir {
		return true
	}
	native := filepath.FromSlash(path)
	if path == "" || filepath.IsAbs(native) || filepath.Clean(native) != native || filepath.ToSlash(native) != path {
		return false
	}
	return path == strings.TrimSuffix(dir, "/") || strings.HasPrefix(path, dir)
}
