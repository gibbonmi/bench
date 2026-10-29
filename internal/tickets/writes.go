package tickets

import (
	"errors"
	"path"
	"path/filepath"
	"strings"
)

// NewMarker declares a `Writes:` entry as a path the ticket creates. An entry
// carrying it is green whether or not the tree already holds the path, because
// a blocker ticket may land the file first.
const NewMarker = "(new)"

// deletedMarker declares a `Writes:` entry as one exact file that the ticket removes. The
// file must be tracked at the source base; the tip may hold it or not.
const deletedMarker = "(deleted)"

// WriteKind is the authority one `Writes:` entry declares over its path.
type WriteKind uint8

const (
	// WriteExisting names a path that the tree holds.
	WriteExisting WriteKind = iota
	// WriteNew names a path that the ticket creates.
	WriteNew
	// WriteDeleted names one exact file that the ticket removes.
	WriteDeleted
)

// Write is one parsed `Writes:` entry: the tree path it names and its kind.
type Write struct {
	Path string
	Kind WriteKind
}

// ParseWrite is the one `Writes:` entry grammar. It strips one trailing marker and a
// trailing `/` directory spelling. An entry cannot carry both markers, because a created
// path and a removed file are contradictory authority. A deletion names an exact file, so
// a directory spelling refuses. The path must be a RepositoryPath.
func ParseWrite(entry string) (Write, error) {
	write := Write{Path: strings.TrimSpace(entry)}
	markers := []struct {
		text string
		kind WriteKind
	}{{NewMarker, WriteNew}, {deletedMarker, WriteDeleted}}
	for i, marker := range markers {
		rest, found := strings.CutSuffix(write.Path, marker.text)
		if !found {
			continue
		}
		write.Path, write.Kind = strings.TrimSpace(rest), marker.kind
		if strings.HasSuffix(write.Path, markers[1-i].text) {
			return Write{}, errors.New("an entry cannot carry both the " + NewMarker + " and the " + deletedMarker + " marker")
		}
		break
	}
	if write.Kind == WriteDeleted && strings.HasSuffix(write.Path, "/") {
		return Write{}, errors.New("a " + deletedMarker + " entry names one exact file, not a directory")
	}
	write.Path = strings.TrimSuffix(write.Path, "/")
	if !RepositoryPath(write.Path) {
		return Write{}, errors.New("the path is the repository root, escapes the repository, or is not clean")
	}
	return write, nil
}

// RepositoryPath reports whether one slash-spelled path names a place below the repository
// root by its clean spelling. Covers matches an entry by its spelling, so a ".." segment would
// let the entry cover a file outside the fence or the repository, and the root itself would
// cover every path.
func RepositoryPath(p string) bool {
	return p != "." && filepath.IsLocal(filepath.FromSlash(p)) && path.Clean(p) == p
}

// WritesPath projects ParseWrite onto the tree path one entry names and whether it
// carries the (new) marker. An entry that ParseWrite refuses keeps its trimmed spelling,
// which names no tree path, and the ticket grammar reports the refusal.
func WritesPath(entry string) (path string, isNew bool) {
	write, err := ParseWrite(entry)
	if err != nil {
		return strings.TrimSpace(entry), false
	}
	return write.Path, write.Kind == WriteNew
}

// Covers reports whether one split entry names path exactly or contains it at a
// `/` segment boundary, so `internal/git` never covers `internal/git2`.
func Covers(entry, path string) bool {
	return path == entry || strings.HasPrefix(path, entry+"/")
}
