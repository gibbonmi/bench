// Package modelid owns the portable model-token grammar Bench accepts in
// reviewer-owned tier bindings and advisory discovery output.
package modelid

import "regexp"

var safeTokenRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@+-]*$`)

// SafeToken reports whether value is a non-empty printable model-id token that
// can be handed to a harness without whitespace, control bytes, or shell syntax.
func SafeToken(value string) bool {
	return safeTokenRe.MatchString(value)
}
