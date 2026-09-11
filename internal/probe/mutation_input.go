package probe

import (
	"errors"
	"os"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

func mutationInput(parsed usage.Result) (string, string, string, string) {
	if path, ok := parsed.Flags["--omit-file"]; ok {
		fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
		if err != nil {
			reason := "unreadable"
			switch {
			case errors.Is(err, syscall.ELOOP):
				reason = "a symlink"
			case errors.Is(err, syscall.ENOENT):
				reason = "absent"
			}
			return "", "", "", mutationFileUnavailable(path, reason)
		}
		file := os.NewFile(uintptr(fd), path)
		info, statErr := file.Stat()
		if statErr != nil {
			_ = file.Close()
			return "", "", "", mutationFileUnavailable(path, "unreadable")
		}
		if !info.Mode().IsRegular() {
			_ = file.Close()
			reason := "a special file"
			if info.IsDir() {
				reason = "a directory"
			}
			return "", "", "", mutationFileUnavailable(path, reason)
		}
		read := bounds.Read(file, bounds.ControlRecordLimit)
		closeErr := file.Close()
		if read.Status == bounds.ReadFailed || closeErr != nil {
			return "", "", "", mutationFileUnavailable(path, "unreadable")
		}
		if read.Status == bounds.ReadOversized {
			return "", "", "", mutationFileUnavailable(path, "oversized")
		}
		if len(read.Data) == 0 {
			return "", "", "", mutationFileUnavailable(path, "empty")
		}
		return string(read.Data), "", "omit", ""
	}
	old, replacement, kind := mutationForm(parsed)
	return old, replacement, kind, ""
}

func mutationFileUnavailable(path, reason string) string {
	return toon.Errorf("probe mutation file unavailable", path+" is "+reason) + "\n"
}
