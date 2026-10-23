package protection

import "syscall"

// systemBase is the system temporary directory.
func systemBase() string { return "/private/tmp" }

// bootID is the kernel boot session identifier.
func bootID() (string, error) { return syscall.Sysctl("kern.bootsessionuuid") }
