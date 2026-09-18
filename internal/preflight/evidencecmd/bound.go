package evidencecmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/toon"
)

// maxOperandBytes bounds every preflight operand before the parser can echo it. Every valid
// selector, pin, identifier, cursor, and quota is far shorter.
const maxOperandBytes = 1024

// OversizedOperand refuses an operand too long to echo, naming it only by length and digest.
func OversizedOperand(args []string) string {
	for _, arg := range args {
		if len(arg) > maxOperandBytes {
			return toon.Usage(Grammar.Cmd, boundedOperand("oversized operand", arg))
		}
	}
	return ""
}

// BoundedDiagnostic returns value when it fits the operand bound. A longer value is
// identified by kind, length, and digest instead.
func BoundedDiagnostic(kind, value string) string {
	if len(value) <= maxOperandBytes {
		return value
	}
	return boundedOperand(kind, value)
}

// boundedOperand identifies a hostile operand by type, length, and digest, never by value.
func boundedOperand(kind, value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%s bytes=%d sha256=%s", kind, len(value), hex.EncodeToString(sum[:]))
}

// responseLimit is the bound Bound enforces. Only SetResponseLimitForTest changes it.
var responseLimit = chargeevidence.ResponseLimit

// SetResponseLimitForTest lowers the shared bound in process, so a test can drive an
// over-bound response through each bounded path. It returns the restore function.
func SetResponseLimitForTest(limit int) func() {
	previous := responseLimit
	responseLimit = limit
	return func() { responseLimit = previous }
}

// Bound is the shared final guard for bounded forms and every usage line: a response above
// the response limit becomes one bounded operational refusal.
func Bound(out string, code int) (string, int) {
	if len(out) <= responseLimit {
		return out, code
	}
	return toon.Errorf("response bound exceeded", fmt.Sprintf("the response held %d bytes above the %d-byte limit; report this defect", len(out), responseLimit)) + "\n", 1
}
