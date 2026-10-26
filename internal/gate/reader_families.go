package gate

// The reader-family declaration: the covered path families that the reader table maps to
// their test packages, and the child readers that no test log can show.

import "strings"

// coveredFamilies are the covered path families. A family that ends with `/` claims every
// path below it, and each other family claims its exact path.
var coveredFamilies = []string{
	"specs/",
	"roadmap/",
	"ROADMAP.md",
	"capture/",
	"reviews/",
	"decisions/",
	"CHANGELOG.md",
	".bench/commitment.json",
	".agents/",
	".bench/BENCH.md",
	"docs/",
}

// childReader is a test package that reads a covered family inside a child process. The
// test log of the package holds only the reads of its own process, so the declaration
// supplies the child's read.
type childReader struct {
	Package, Family string
}

// childReaders are the declared child readers.
var childReaders = []childReader{
	{Package: "internal/releaseevidence", Family: "CHANGELOG.md"},
}

// familyOf answers the family of families that claims the slash-separated path rel. A
// directory family also claims its directory itself, because a listing of the directory
// reads which paths the family holds.
func familyOf(families []string, rel string) (string, bool) {
	for _, family := range families {
		dir, isDir := strings.CutSuffix(family, "/")
		if rel == dir || (isDir && strings.HasPrefix(rel, family)) {
			return family, true
		}
	}
	return "", false
}
