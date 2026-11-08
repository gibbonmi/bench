package protection

import "syscall"

// systemBase is the system temporary directory.
func systemBase() string { return "/private/tmp" }

// BootID is the kernel boot session identifier.
func BootID() (string, error) { return syscall.Sysctl("kern.bootsessionuuid") }

// systemProbe sends signal to target. This system has no zombie rule, so an unreaped zombie
// answers signal zero as present.
func systemProbe(target int, signal syscall.Signal) error { return syscall.Kill(target, signal) }
