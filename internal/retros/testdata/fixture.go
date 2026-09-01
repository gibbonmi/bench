package retrotestdata

import _ "embed"

//go:embed eligible.md
var eligible string

// Eligible returns the canonical eligible retrospective body.
func Eligible() string {
	return eligible
}
