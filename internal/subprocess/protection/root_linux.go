package protection

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// systemBase is the tmpfs /dev/shm, or /tmp when /dev/shm is absent.
func systemBase() string {
	if info, err := os.Stat("/dev/shm"); err == nil && info.IsDir() {
		return "/dev/shm"
	}
	return "/tmp"
}

// BootID is the kernel boot identifier.
func BootID() (string, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	return strings.TrimSpace(string(data)), err
}

// systemProbe sends signal to target. An unreaped zombie answers signal zero but can never use
// a resource again, so a target that answers signal zero and that the proc tree shows as only
// zombies answers ESRCH. It never sends a second signal.
func systemProbe(target int, signal syscall.Signal) error {
	err := syscall.Kill(target, signal)
	if err == nil && signal == 0 && procAbsent("/proc", target) {
		return syscall.ESRCH
	}
	return err
}

// procAbsent applies the zombie rule to target in the proc tree at root. A process is absent
// when the rule finds a task and every task is a zombie. A process group, a negative target,
// is absent when the rule finds a member and no member is present under the process rule. A
// record that the rule cannot read or parse keeps the target present, and an entry that is
// gone decides nothing.
func procAbsent(root string, target int) bool {
	if target > 0 {
		dead, found := zombie(root, strconv.Itoa(target))
		return dead && found
	}
	if target > -2 {
		// Zero and minus one name the caller's group and every process, not one group.
		return false
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return false
	}
	found := false
	for _, entry := range entries {
		if _, err := strconv.Atoi(entry.Name()); err != nil {
			continue
		}
		fields, err := statFields(filepath.Join(root, entry.Name(), "stat"))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return false
		}
		pgid, err := strconv.Atoi(fields[2])
		if err != nil {
			return false
		}
		if pgid != -target {
			continue
		}
		dead, member := zombie(root, entry.Name())
		if !dead {
			return false
		}
		found = found || member
	}
	return found
}

// zombie reads each task of process pid in the proc tree at root. It reports whether every
// task that it finds is a zombie, and whether it finds one. A gone process finds none.
func zombie(root, pid string) (dead, found bool) {
	tasks := filepath.Join(root, pid, "task")
	entries, err := os.ReadDir(tasks)
	if err != nil {
		return errors.Is(err, fs.ErrNotExist), false
	}
	for _, entry := range entries {
		fields, err := statFields(filepath.Join(tasks, entry.Name(), "stat"))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil || fields[0] != "Z" {
			return false, false
		}
		found = true
	}
	return true, found
}

// statFields reads one stat record and returns its fields after the command name: the state,
// the parent, and the process group first. The command name can hold any byte, so it ends at
// the last ")" of the record.
func statFields(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	end := bytes.LastIndexByte(data, ')')
	fields := strings.Fields(string(data[end+1:]))
	if end < 0 || len(fields) < 3 {
		return nil, fmt.Errorf("protection: %s is not a stat record", path)
	}
	return fields, nil
}
