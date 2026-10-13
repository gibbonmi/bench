package probe

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// mutation is one parsed mutation input: the edits it asks for over the start bytes, or the
// removal of the subject, and the cell that names its form. occurrence is 0 when the argv
// names no --occurrence.
type mutation struct {
	cell       string
	finds      []find
	occurrence int
	remove     bool
}

// mutationInput resolves an argv that gradeForms accepted into the mutation of its one
// group, and reads each mutation file through the one bounded reader. A combined mutant
// lists each --swap edit and then each --omit edit, so its cell derives from the two flags,
// not argv order.
func mutationInput(parsed usage.Result) (mutation, string) {
	n, _ := occurrence(parsed)
	m := mutation{occurrence: n}
	text, textFlag := parsed.Flags["--with"], "--with"
	if path, ok := parsed.Flags["--with-file"]; ok {
		read, line := readMutationFile(path)
		if line != "" {
			return mutation{}, line
		}
		text, textFlag = read, "--with-file"
	}
	group, _ := namedGroup(parsed)
	switch group {
	case deleteGroup:
		m.cell, m.remove = "delete", true
		return m, ""
	case omitFileGroup:
		read, line := readMutationFile(parsed.Flags["--omit-file"])
		if line != "" {
			return mutation{}, line
		}
		m.cell, m.finds = "omit", []find{{form: "omit", old: read}}
		return m, ""
	case unwrapGroup:
		m.cell, m.finds = "unwrap", []find{{form: "unwrap", old: parsed.Flags["--unwrap"]}}
		return m, ""
	case insertGroup:
		m.cell, m.finds = "insert", []find{{form: "insert", old: parsed.Flags["--insert-before"], text: text}}
		return m, ""
	}
	withs := parsed.Repeated["--with"]
	var forms []string
	for i, old := range parsed.Repeated["--swap"] {
		// The i-th --swap pairs with the i-th --with. A --with-file swap has no --with
		// value, so it keeps the file text.
		if i < len(withs) {
			text = withs[i]
		}
		m.finds = append(m.finds, find{form: "swap", old: old, text: text, textFlag: textFlag})
		forms = append(forms, "swap")
	}
	for _, old := range parsed.Repeated["--omit"] {
		m.finds = append(m.finds, find{form: "omit", old: old})
		forms = append(forms, "omit")
	}
	m.cell = strings.Join(forms, "+")
	return m, ""
}

// occurrence answers the --occurrence value, or 0 when the argv names none. ok is false for
// a value that is not a positive decimal integer.
func occurrence(parsed usage.Result) (int, bool) {
	value, set := parsed.Flags["--occurrence"]
	if !set {
		return 0, true
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	n, err := strconv.Atoi(value)
	return n, err == nil && n > 0
}

// readMutationFile reads one mutation file for --omit-file or --with-file. The open does not
// follow a symlink and does not block on a FIFO, and the read is bounded, so a hostile path
// answers its refusal instead of mutating or hanging.
func readMutationFile(path string) (string, string) {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		reason := "unreadable"
		switch {
		case errors.Is(err, syscall.ELOOP):
			reason = "a symlink"
		case errors.Is(err, syscall.ENOENT):
			reason = "absent"
		}
		return "", mutationFileUnavailable(path, reason)
	}
	file := os.NewFile(uintptr(fd), path)
	info, statErr := file.Stat()
	if statErr != nil {
		_ = file.Close()
		return "", mutationFileUnavailable(path, "unreadable")
	}
	if !info.Mode().IsRegular() {
		_ = file.Close()
		reason := "a special file"
		if info.IsDir() {
			reason = "a directory"
		}
		return "", mutationFileUnavailable(path, reason)
	}
	read := bounds.Read(file, bounds.ControlRecordLimit)
	closeErr := file.Close()
	if read.Status == bounds.ReadFailed || closeErr != nil {
		return "", mutationFileUnavailable(path, "unreadable")
	}
	if read.Status == bounds.ReadOversized {
		return "", mutationFileUnavailable(path, "oversized")
	}
	if len(read.Data) == 0 {
		return "", mutationFileUnavailable(path, "empty")
	}
	return string(read.Data), ""
}

func mutationFileUnavailable(path, reason string) string {
	return toon.Errorf("probe mutation file unavailable", path+" is "+reason) + "\n"
}
