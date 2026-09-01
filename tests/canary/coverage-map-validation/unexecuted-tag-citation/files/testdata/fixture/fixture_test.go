package fixture

import "testing"

// TestIgnoredPackage is fixture evidence: Go excludes a testdata package from every
// recursive pattern, so no gate test phase selects this file.
func TestIgnoredPackage(t *testing.T) {}
