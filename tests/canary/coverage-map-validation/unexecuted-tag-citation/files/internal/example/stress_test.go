//go:build stress

package example

import "testing"

// TestStressOnly is release-only evidence: no executed tag set builds this file.
func TestStressOnly(t *testing.T) {}
