package tickets

import "strings"

// newMarker declares a `Writes:` entry as a path the ticket creates. An entry
// carrying it is green whether or not the tree already holds the path, because
// a blocker ticket may land the file first.
const newMarker = "(new)"

// WritesPath separates one `Writes:` entry into the tree path it names and whether
// it carries the (new) marker. A trailing `/` is a directory spelling, not a path
// segment, so the split drops it.
func WritesPath(entry string) (path string, isNew bool) {
	path = strings.TrimSpace(entry)
	if isNew = strings.HasSuffix(path, newMarker); isNew {
		path = strings.TrimSpace(strings.TrimSuffix(path, newMarker))
	}
	return strings.TrimSuffix(path, "/"), isNew
}

// Covers reports whether one split entry names path exactly or contains it at a
// `/` segment boundary, so `internal/git` never covers `internal/git2`.
func Covers(entry, path string) bool {
	return path == entry || strings.HasPrefix(path, entry+"/")
}
