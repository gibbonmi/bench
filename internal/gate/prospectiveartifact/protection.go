package prospectiveartifact

import "github.com/gibbonmi/bench/internal/subprocess/protection"

// bundleDirectory is the identity adapter of a bundle root.
var bundleDirectory = protection.Directory{
	Kind:   protection.KindProspective,
	Noun:   "bundle",
	Prefix: BundlePrefix,
	Class:  "prospective artifact bundle",
}

// ValidateProtection is the bundle identity adapter.
func ValidateProtection(d protection.Descriptor) error { return bundleDirectory.Validate(d) }
