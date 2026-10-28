package protection

import "syscall"

// systemBase is the system temporary directory.
func systemBase() string { return "/private/tmp" }

// BootID is the kernel boot session identifier.
func BootID() (string, error) { return syscall.Sysctl("kern.bootsessionuuid") }
