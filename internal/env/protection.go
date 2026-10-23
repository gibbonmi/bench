package env

import "github.com/gibbonmi/bench/internal/subprocess/protection"

// ProtectionStore names the protection store inside a private test run.
const ProtectionStore = "process-users"

// kitRunProtection is the identity adapter of a private test run directory.
var kitRunProtection = protection.Directory{Kind: protection.KindKitTest, Store: ProtectionStore, Noun: "kit test run"}

// ValidateProtection is the private test run's identity adapter.
func ValidateProtection(d protection.Descriptor) error { return kitRunProtection.Validate(d) }
