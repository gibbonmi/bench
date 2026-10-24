package protection

import (
	"os"
	"strings"
)

// systemBase is the tmpfs /dev/shm, or /tmp when /dev/shm is absent.
func systemBase() string {
	if info, err := os.Stat("/dev/shm"); err == nil && info.IsDir() {
		return "/dev/shm"
	}
	return "/tmp"
}

// bootID is the kernel boot identifier.
func bootID() (string, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	return strings.TrimSpace(string(data)), err
}
