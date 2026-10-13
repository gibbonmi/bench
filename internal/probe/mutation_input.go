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

// refusal is one mutation refusal before it renders: the title and the detail of its error
// line. A mutant-set entry prefixes the detail with its entry number.
type refusal struct {
	title, detail string
}

func (r *refusal) line() string {
	return toon.Errorf(r.title, r.detail) + "\n"
}

// mutationInput resolves an argv whose mutation half the form matrix accepted into the
// mutation of its one group. It reads each mutation file through the one bounded reader. A
// combined mutant lists each --swap edit and then each --omit edit, so its cell derives from
// the two flags, not argv order.
func mutationInput(parsed usage.Result) (mutation, *refusal) {
	n, _ := occurrence(parsed)
	m := mutation{occurrence: n}
	text, textFlag := parsed.Flags["--with"], "--with"
	if path, ok := parsed.Flags["--with-file"]; ok {
		read, refused := readMutationFile(path)
		if refused != nil {
			return mutation{}, refused
		}
		text, textFlag = read, "--with-file"
	}
	group, _ := namedGroup(parsed)
	switch group {
	case deleteGroup:
		m.cell, m.remove = "delete", true
		return m, nil
	case omitFileGroup:
		read, refused := readMutationFile(parsed.Flags["--omit-file"])
		if refused != nil {
			return mutation{}, refused
		}
		m.cell, m.finds = "omit", []find{{form: "omit", old: read}}
		return m, nil
	case unwrapGroup:
		m.cell, m.finds = "unwrap", []find{{form: "unwrap", old: parsed.Flags["--unwrap"]}}
		return m, nil
	case insertGroup:
		m.cell, m.finds = "insert", []find{{form: "insert", old: parsed.Flags["--insert-before"], text: text}}
		return m, nil
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
	return m, nil
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

// readMutationFile reads one mutation file for --omit-file, --with-file, or --mutants. The
// open does not follow a symlink and does not block on a FIFO, and the read is bounded. So a
// hostile path answers its refusal instead of mutating or hanging.
func readMutationFile(path string) (string, *refusal) {
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
	return string(read.Data), nil
}

func mutationFileUnavailable(path, reason string) *refusal {
	return &refusal{title: "probe mutation file unavailable", detail: path + " is " + reason}
}
