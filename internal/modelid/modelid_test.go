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
}
