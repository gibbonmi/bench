package evidencecmd

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
)

// TestEvidenceUnrepresentableQuota proves that a candidate no representable quota admits
// refuses with the bounded capacity classification and no wrapped retry quota. Cleanup is
// still the first recovery action, because no larger quota exists to offer as the
// alternative.
func TestEvidenceUnrepresentableQuota(t *testing.T) {
	out := storeRefusal(chargeevidence.Admit(math.MaxUint64, 1, math.MaxUint64))
	want := fmt.Sprintf("error: evidence capacity: the store holds %d bytes, the candidate needs 1 bytes, and the quota is %d bytes — run %s; no representable quota admits the candidate\n", uint64(math.MaxUint64), uint64(math.MaxUint64), CleanCommand)
	if out != want || strings.Contains(out, "retry with") {
		t.Fatalf("unrepresentable quota refusal = %q, want %q", out, want)
	}
}
