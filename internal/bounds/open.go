package bounds

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// OpenNoFollow opens the regular file at path for a read whose limit the caller owns. It
// composes the producer-file policies without changing them: RefuseLinks refuses a link in
// any component, and the leaf is typed before the open, so a link cannot redirect the read
// and a FIFO cannot block it in open(2). The open itself refuses a link, and the opened file
// must be the classified one, so a file replaced between the type check and the open is
// unreadable. A file that opens answers StateEmpty or StateParsed by its size; the caller
// reads and closes it. Every other state answers a nil file.
func OpenNoFollow(path string) (*os.File, Classified) {
	if err := RefuseLinks(path); err != nil {
		if errors.Is(err, ErrLinkRefused) {
			return nil, Classified{State: StateWrongType, Reason: err.Error()}
		}
		return nil, Classified{State: StateUnreadable, Reason: err.Error()}
	}
	info, state, reason := resolveNoFollow(path)
	if state != "" {
		return nil, Classified{State: state, Reason: reason}
	}
	if !info.Mode().IsRegular() {
		return nil, Classified{State: StateWrongType, Reason: fmt.Sprintf("not a regular file: %s", info.Mode().Type())}
	}
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, Classified{State: StateUnreadable, Reason: err.Error()}
	}
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		_ = file.Close()
		return nil, Classified{State: StateUnreadable, Reason: "the file changed before it opened"}
	}
	if opened.Size() == 0 {
		return file, Classified{State: StateEmpty}
	}
	return file, Classified{State: StateParsed}
}
