package preflight

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestEvidenceQuotaOperands is CE130: an invalid quota refuses before any store change.
func TestEvidenceQuotaOperands(t *testing.T) {
	for _, quota := range []string{"0", "-1", "+1", "01", "1.5", "1e3", "abc", " 1", "18446744073709551616"} {
		t.Run(quota, func(t *testing.T) {
			root, slug := seedConformant(t)
			args := append(chargeArgs(t, root, slug, false), "--max-store-bytes", quota)
			out, code := Command(args)
			if code != 2 || !strings.Contains(out, "--max-store-bytes needs a positive decimal byte count") {
				t.Fatalf("quota %q = (%d):\n%s", quota, code, out)
			}
			if _, err := os.Lstat(storeDir(t, root)); !os.IsNotExist(err) {
				t.Fatalf("quota %q changed the store: %v", quota, err)
			}
		})
	}
}

var retryQuota = regexp.MustCompile(` — retry with --max-store-bytes ([0-9]+)\n$`)

// capacityRequired runs one preparation under quota and returns the required byte count
// its capacity refusal names.
func capacityRequired(t *testing.T, args []string, quota string) uint64 {
	t.Helper()
	out, code := Command(append(append([]string{}, args...), "--max-store-bytes", quota))
	match := retryQuota.FindStringSubmatch(out)
	if code != 1 || match == nil || !strings.HasPrefix(out, "error: evidence capacity: ") || strings.Count(out, "\n") != 1 {
		t.Fatalf("capacity refusal under %s = (%d):\n%s", quota, code, out)
	}
	required, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return required
}

// TestEvidenceCapacityRecovery is CE157. Before cleanup exists, a capacity refusal names
// exactly the larger quota that admits the candidate, counts every existing artifact, and
// preserves them.
func TestEvidenceCapacityRecovery(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, false)
	required := capacityRequired(t, args, "1")
	if capacityRequired(t, args, strconv.FormatUint(required-1, 10)) != required {
		t.Fatal("the named quota changed between refusals")
	}
	identity, _, _ := prepareEvidence(t, append(append([]string{}, args...), "--max-store-bytes", strconv.FormatUint(required, 10)))
	first := publishedPacks(t, root)

	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticketDoc("One", "PF1", "PF2")+"second artifact\n")
	second := legacyCommitted(t, root, slug, "second artifact", false)
	next := capacityRequired(t, second, strconv.FormatUint(required, 10))
	if next <= required {
		t.Fatalf("second required %d does not count the first artifact of %d bytes", next, required)
	}
	if got := publishedPacks(t, root); strings.Join(got, ",") != strings.Join(first, ",") || len(got) != 1 {
		t.Fatalf("capacity refusal changed the published artifacts: %v, want %v", got, first)
	}
	if out, code := Command([]string{"evidence", identity}); code != 0 || !strings.HasPrefix(out, "page[1]") {
		t.Fatalf("preserved artifact read = (%d):\n%s", code, out)
	}
	prepareEvidence(t, append(append([]string{}, second...), "--max-store-bytes", strconv.FormatUint(next, 10)))
	if len(publishedPacks(t, root)) != 2 {
		t.Fatalf("the named quota did not admit the second artifact: %v", publishedPacks(t, root))
	}
}
