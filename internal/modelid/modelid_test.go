package modelid

import (
	"testing"

	"github.com/gibbonmi/bench/internal/modelid/modelidtest"
)

func TestSafeToken(t *testing.T) {
	for _, value := range modelidtest.AcceptedTokens {
		t.Run("accept "+value, func(t *testing.T) {
			if !SafeToken(value) {
				t.Fatalf("SafeToken(%q) = false, want true", value)
			}
		})
	}
	for _, token := range modelidtest.RejectedTokens {
		t.Run("reject "+token.Name, func(t *testing.T) {
			if SafeToken(token.Value) {
				t.Fatalf("SafeToken(%q) = true, want false", token.Value)
			}
		})
	}
	// The newline-class rejects pin that the grammar's `$` anchor is end-of-text
	// (Go's default `\z`), not multiline: a token carrying a shell newline stays
	// rejected. They ride their own slice because RejectedTokens' second consumer
	// (the conformance line-binding check) would sanitize the newline away before
	// SafeToken saw it — see modelidtest.NewlineRejectedTokens.
	for _, token := range modelidtest.NewlineRejectedTokens {
		t.Run("reject "+token.Name, func(t *testing.T) {
			if SafeToken(token.Value) {
				t.Fatalf("SafeToken(%q) = true, want false", token.Value)
			}
		})
	}
}
