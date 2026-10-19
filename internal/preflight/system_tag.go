package preflight

import (
	"go/build/constraint"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// kitEnvMarker is the literal a ticket states to pin the kit a system-tagged
// fixture reads. An ambient value flips such a fixture under composition, so the
// ticket names the variable rather than inheriting whatever the session holds.
const kitEnvMarker = "BENCH_KIT"

// systemTagged reports whether one `Writes:` path names a Go test file whose
// //go:build expression carries the system tag. The constraint parse is the same
// one the coverage citations grade with, so the two never disagree about which
// file is system-tagged. A path that is no test file, that does not read, or that
// carries no constraint answers no.
func systemTagged(root, path string) bool { return systemTagLine(root, path) > 0 }

// systemTagLine is the line of path's //go:build constraint when that constraint carries
// the system tag, and zero otherwise.
func systemTagLine(root, path string) int {
	if !strings.HasSuffix(path, "_test.go") {
		return 0
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		return 0
	}
	for i, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			return 0
		}
		if !constraint.IsGoBuild(trimmed) {
			continue
		}
		if expr, parseErr := constraint.Parse(trimmed); parseErr == nil && exprHoldsTag(expr, "system") {
			return i + 1
		}
		return 0
	}
	return 0
}

// exprHoldsTag reports whether a build expression mentions tag anywhere. A
// negated mention still counts: the file's verdict depends on the tag either way.
func exprHoldsTag(expr constraint.Expr, tag string) bool {
	switch typed := expr.(type) {
	case *constraint.TagExpr:
		return typed.Tag == tag
	case *constraint.NotExpr:
		return exprHoldsTag(typed.X, tag)
	case *constraint.AndExpr:
		return exprHoldsTag(typed.X, tag) || exprHoldsTag(typed.Y, tag)
	case *constraint.OrExpr:
		return exprHoldsTag(typed.X, tag) || exprHoldsTag(typed.Y, tag)
	}
	return false
}

// newSystemTest reports whether a Go test path the tree does not hold yet joins the
// system suite. The file has no bytes to read, so its directory answers: a directory
// that holds a system-tagged test file is where the system tag lives. The tag parse is
// systemTagged's own, so a new file and a landed one are graded by one rule. Only a
// regular file is read, because a FIFO named like a test file would block the read.
func newSystemTest(root, testPath string) bool {
	owner, _ := newSystemTestSource(root, testPath)
	return owner != ""
}

// newSystemTestSource is the first system-tagged sibling that gives a test path the tree
// does not hold its tag, and that sibling's constraint line.
func newSystemTestSource(root, testPath string) (string, int) {
	if !strings.HasSuffix(testPath, "_test.go") {
		return "", 0
	}
	dir := path.Dir(testPath)
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return "", 0
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		sibling := path.Join(dir, entry.Name())
		if line := systemTagLine(root, sibling); line > 0 {
			return sibling, line
		}
	}
	return "", 0
}
