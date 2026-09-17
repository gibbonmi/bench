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

// boundedOperand identifies a hostile operand by type, length, and digest, never by value.
func boundedOperand(kind, value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%s bytes=%d sha256=%s", kind, len(value), hex.EncodeToString(sum[:]))
}

// Bound is the shared final guard for bounded forms and every usage line: a response above
// chargeevidence.ResponseLimit encoded bytes becomes one bounded operational refusal.
func Bound(out string, code int) (string, int) {
	if len(out) <= chargeevidence.ResponseLimit {
		return out, code
	}
	return toon.Errorf("response bound exceeded", fmt.Sprintf("the response held %d bytes above the %d-byte limit; report this defect", len(out), chargeevidence.ResponseLimit)) + "\n", 1
}
