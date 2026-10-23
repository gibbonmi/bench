package env

import "github.com/gibbonmi/bench/internal/subprocess/protection"

// kitRunProtection is the identity adapter of a private test run directory.
var kitRunProtection = protection.Directory{Kind: protection.KindKitTest, Noun: "kit test run"}

// ValidateProtection is the private test run's identity adapter.
func ValidateProtection(d protection.Descriptor) error { return kitRunProtection.Validate(d) }
