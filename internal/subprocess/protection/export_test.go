package protection

import (
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// The store layout names, for tests that plant hostile records.
const (
	ResourceFile = resourceFile
	UsersDir     = usersDir
	UserFile     = userFile
)

// SystemBase is the system directory that holds each user's store root.
var SystemBase = systemBase

// SetStoreBase places every store root under base until t ends, so a test never changes the
// real root of the user.
func SetStoreBase(t testing.TB, base string) {
	t.Helper()
	storeBase = func() string { return base }
	t.Cleanup(func() { storeBase = systemBase })
}

// SetReplacer routes p's record persistence through r, the durablefile fault seam.
func SetReplacer(p *Protocol, r durablefile.Replacer) { p.replacer = r }

// SetProbe routes p's recovery observations through probe.
func SetProbe(p *Protocol, probe Probe) { p.probe = probe }
