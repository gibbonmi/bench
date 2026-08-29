package outline

import (
	"bytes"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// fixtureDir is the path segment that marks prior-art fixture data.
const fixtureDir = "testdata"

// isFixture reports whether git's slash-separated path sits under a fixtureDir segment.
// Such a file reports one fixture row and is never scanned, because it carries no scanned
// extension and the extension dispatch would drop it silently.
func isFixture(rel string) bool {
	return strings.HasPrefix(rel, fixtureDir+"/") || strings.Contains(rel, "/"+fixtureDir+"/")
}

// The skip reasons the read policy reports. A caller renders one in the reason column of
// its own skip table, so the vocabulary has one source.
const (
	skipUnreadable = "unreadable"
	skipNonregular = "nonregular"
	skipOversized  = "oversized"
	skipBinary     = "binary"
)

var openOutlineFile = os.Open

// statPolicy is the stat half of the read policy. The stat is an Lstat, so a symlink is
// not followed out of the checkout, and a nonregular entry is never opened, so a FIFO
// cannot block the caller.
func statPolicy(abs string) (string, bool) {
	info, err := os.Lstat(abs)
	if err != nil {
		return skipUnreadable, false
	}
	if !info.Mode().IsRegular() {
		return skipNonregular, false
	}
	return "", true
}

// ReadScannable is the one source of the read policy a tracked file must pass before a
// pattern scan reads it: the stat policy above, a read bounded by bounds.OutlineFileLimit,
// and a NUL byte anywhere in the content. It returns the whole content when the file is
// scannable, and the reason it is not otherwise. It is exported so a caller sweeping the
// same files applies this policy rather than a second one.
func ReadScannable(abs string) ([]byte, string, bool) {
	if reason, ok := statPolicy(abs); !ok {
		return nil, reason, false
	}
	file, err := openOutlineFile(abs)
	if err != nil {
		return nil, skipUnreadable, false
	}
	read := bounds.Read(file, bounds.OutlineFileLimit)
	closeErr := file.Close()
	if read.Status == bounds.ReadOversized {
		return nil, skipOversized, false
	}
	if read.Status != bounds.ReadComplete || closeErr != nil {
		return nil, skipUnreadable, false
	}
	if bytes.IndexByte(read.Data, 0) >= 0 {
		return nil, skipBinary, false
	}
	return read.Data, "", true
}
