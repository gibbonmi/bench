package target

import "testing"

// TestSymbol is the in-package test-file consumer. The loader loads this package once
// plainly and once as the test variant, so one declaration arrives as two objects.
func TestSymbol(t *testing.T) { Symbol() }
