package protection

import "github.com/gibbonmi/bench/internal/durablefile"

// The store layout names, for tests that plant hostile records.
const (
	ResourceFile = resourceFile
	UsersDir     = usersDir
	UserFile     = userFile
)

// SetReplacer routes p's record persistence through r, the durablefile fault seam.
func SetReplacer(p *Protocol, r durablefile.Replacer) { p.replacer = r }

// SetProbe routes p's recovery observations through probe.
func SetProbe(p *Protocol, probe Probe) { p.probe = probe }
