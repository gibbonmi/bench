package runbinary

import "github.com/gibbonmi/bench/internal/subprocess/protection"

// ProtectionStore names the protection store inside a private selection directory.
const ProtectionStore = "process-users"

// runDirPrefix starts the name of every private selection directory.
const runDirPrefix = "bench-run-"

// selectionDirectory is the identity adapter of a private selection directory.
var selectionDirectory = protection.Directory{
	Kind:   protection.KindBinary,
	Store:  ProtectionStore,
	Noun:   "selection",
	Prefix: runDirPrefix,
	Class:  "private selection directory",
}

// ValidateProtection is the selection identity adapter.
func ValidateProtection(d protection.Descriptor) error { return selectionDirectory.Validate(d) }
