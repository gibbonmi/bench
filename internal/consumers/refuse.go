package consumers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/outline"
	"github.com/gibbonmi/bench/internal/toon"
)

// missingGoMarker is the text go/packages reports when the go tool is absent from PATH.
// The match is textual because packages.Load formats the exec failure into a plain error
// string and keeps no wrapped *exec.Error for errors.As to find.
const missingGoMarker = "executable file not found"

// refuseLoad turns a loader failure into the refusal that names its cause. The loader
// already fails closed on the first package error and formats its position, so this
// function chooses the cause and the remedy rather than restating either position or
// message.
func refuseLoad(err error) string {
	if strings.Contains(err.Error(), missingGoMarker) {
		return toon.Errorf("go binary not found on PATH", "install Go or put the go binary on PATH, then retry")
	}
	return toon.Errorf("package load failed: "+err.Error(),
		"the tree must type-check before any enumeration; fix the named position, then retry")
}

// refuseLoadForQuery chooses the refusal one failed load earns for one query. A tree with
// no go.mod cannot type-check at all, so the loader's failure says nothing about the
// query, and the non-Go sweep answers first when it names the declaration. A tree that
// has a module keeps the load refusal, because there the failure is the real cause.
func refuseLoadForQuery(root, query string, err error) string {
	if _, statErr := os.Stat(filepath.Join(root, "go.mod")); statErr != nil {
		if file, ok := nonGoDeclaration(root, query); ok {
			return refuseNonGo(file, query)
		}
	}
	return refuseLoad(err)
}

// refuseNonGo is the one spelling of the non-Go refusal, so the unresolved path and the
// unloadable-tree path name the language the same way.
func refuseNonGo(file, query string) string {
	return toon.Errorf("no Go declaration named "+query+"; "+file+" declares it in "+languageOf(file),
		"non-Go resolution is unsupported; a textual sweep is the candidate-class citation")
}

// languageNames maps a scanned extension to the language a refusal names. An extension
// outside the table names itself, so a new file type still refuses legibly.
var languageNames = map[string]string{
	".ts":   "TypeScript",
	".tsx":  "TypeScript",
	".js":   "JavaScript",
	".mjs":  "JavaScript",
	".py":   "Python",
	".sh":   "shell",
	".bash": "shell",
	".md":   "Markdown",
}

func languageOf(rel string) string {
	if name, ok := languageNames[strings.ToLower(filepath.Ext(rel))]; ok {
		return name
	}
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(rel)), ".")
}

// refuseUnresolved answers a query the Go resolver could not place. A tracked non-Go file
// declaring the same name is the case the empty table would misreport as "no callers", so
// it refuses with the language named. Anything else keeps the plain unresolved refusal.
func refuseUnresolved(root, query string, err error) string {
	if file, ok := nonGoDeclaration(root, query); ok {
		return refuseNonGo(file, query)
	}
	return toon.Errorf(err.Error(), "pass a qualified symbol such as outline.Command")
}

// nonGoDeclaration sweeps the tracked files the way `bench outline` does and reports the
// first non-Go file declaring the query's last segment. The sweep is best effort: an
// unreadable file yields no candidate rather than a second failure mode, because the
// caller is already refusing.
func nonGoDeclaration(root, query string) (string, bool) {
	name := query
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	if name == "" {
		return "", false
	}
	files, err := outline.TrackedFiles(root)
	if err != nil {
		return "", false
	}
	for _, rel := range files {
		if strings.HasSuffix(rel, ".go") {
			continue
		}
		content, ok := readSweepFile(filepath.Join(root, filepath.FromSlash(rel)))
		if !ok {
			continue
		}
		for _, symbol := range outline.Symbols(rel, content) {
			if symbol.Name == name {
				return rel, true
			}
		}
	}
	return "", false
}

// readSweepFile reads one candidate file under the read policy `bench outline` applies to
// a tracked file. The stat is an Lstat and a nonregular entry is skipped rather than
// opened, so a symlink is not followed out of the checkout and a FIFO cannot block the
// sweep, and the read is bounded by bounds.OutlineFileLimit.
func readSweepFile(abs string) ([]byte, bool) {
	info, err := os.Lstat(abs)
	if err != nil || !info.Mode().IsRegular() {
		return nil, false
	}
	file, err := os.Open(abs)
	if err != nil {
		return nil, false
	}
	read := bounds.Read(file, bounds.OutlineFileLimit)
	if closeErr := file.Close(); closeErr != nil || read.Status != bounds.ReadComplete {
		return nil, false
	}
	return read.Data, true
}
