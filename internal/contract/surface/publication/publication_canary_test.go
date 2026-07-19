package publication

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// publicationBehaviorCanary shells a nested `go test`, scoped by build tag and
// run pattern, against internal/publication rooted at SubjectRoot — the same
// pattern internal/contract/surface/artifact's
// TestAuthoritativeNativeProofBehaviorCanary uses for the native-proof canary.
// It only fires when a canary fixture has materialized the named marker test
// file under SubjectRoot's internal/publication; a normal (non-canary) gate
// run never has that file, so this is a no-op skip there. This is the seam
// that lets a canary mutate internal/publication's real, compiled source and
// exercise it, without needing a full `dist/bench` rebuilt inside the fixture.
func publicationBehaviorCanary(t *testing.T, marker, tag, run string) {
	t.Helper()
	root := contract.SubjectRoot(t)
	if _, err := os.Stat(filepath.Join(root, "internal", "publication", marker)); os.IsNotExist(err) {
		t.Skipf("%s canary fixture is not materialized", tag)
	}
	command := exec.Command("go", "test", "-count=1", "-tags="+tag, "-run", run, "./internal/publication")
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("%s canary subprocess failed: %v\n%s", tag, err, output)
	}
}

// TestPublicationOrderBypassBehaviorCanary is coverage row 8a's M1 probe: it
// proves the first-publication path still refuses to publish the wrapper
// ahead of a platform package.
func TestPublicationOrderBypassBehaviorCanary(t *testing.T) {
	publicationBehaviorCanary(t, "order_bypass_canary_test.go", "bench_canary_publication_order_bypass", "^TestPublicationOrderBypassCanary$")
}

// TestPublicationIntegrityMismatchBehaviorCanary is coverage row 8a's M3
// probe: it proves the resume path still refuses to accept an already-live
// package whose registry integrity does not match the approved local tarball.
func TestPublicationIntegrityMismatchBehaviorCanary(t *testing.T) {
	publicationBehaviorCanary(t, "integrity_mismatch_canary_test.go", "bench_canary_publication_integrity_mismatch", "^TestPublicationIntegrityMismatchCanary$")
}

// TestPublicationPrematurePromotionBehaviorCanary is coverage row 8a's M4
// probe: it proves promotion still refuses to move any dist-tag to "latest"
// until the complete approved set reverifies live.
func TestPublicationPrematurePromotionBehaviorCanary(t *testing.T) {
	publicationBehaviorCanary(t, "premature_promotion_canary_test.go", "bench_canary_publication_premature_promotion", "^TestPublicationPrematurePromotionCanary$")
}

// TestPublicationUnpublishAttemptBehaviorCanary is coverage row 8a's M5 probe:
// it proves rollback never issues an unpublish-shaped registry call in place
// of deprecation.
func TestPublicationUnpublishAttemptBehaviorCanary(t *testing.T) {
	publicationBehaviorCanary(t, "unpublish_attempt_canary_test.go", "bench_canary_publication_unpublish_attempt", "^TestPublicationUnpublishAttemptCanary$")
}
