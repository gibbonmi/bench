package tickets

import "slices"

// commandPackages are the package prefixes that hold a command's grammar, sorted. A write
// under one of them can change the arguments, the help, the public inventory, or the
// routing of a command.
var commandPackages = []string{
	"cmd/bench",
	"internal/anchors",
	"internal/assessment",
	"internal/commitment",
	"internal/consumers",
	"internal/coverage",
	"internal/diff",
	"internal/guards",
	"internal/harnesses",
	"internal/learnings",
	"internal/maps",
	"internal/probe",
	"internal/roadmap",
	"internal/worktree",
}

// CommandPackage reports whether path sits under a command package at a segment boundary.
func CommandPackage(path string) bool {
	return slices.ContainsFunc(commandPackages, func(prefix string) bool { return Covers(prefix, path) })
}
