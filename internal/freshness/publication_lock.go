package freshness

import (
	"os"
	"syscall"
)

var publicationFlock = syscall.Flock

func lockPublicationDirectory(path string) (*os.File, error) {
	if err := rejectSymlinkComponents(path); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECTORY, 0)
	if err != nil {
		return nil, err
	}
	for {
		err = publicationFlock(int(file.Fd()), syscall.LOCK_EX)
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	return file, nil
}
