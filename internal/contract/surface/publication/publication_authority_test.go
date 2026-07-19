package publication

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestPublishSubmitRefusesUnlessApprovedIndexIsFullGreenPublishRunForProfile is
// coverage row 7's publication-side composition: bench release submit must
// refuse (exit 1, structured error, no registry call at all) unless
// dist/preflight/release-index.json is a full (non-focused) green
// publish-mode preflight run whose recorded profile matches --profile.
// Publication never re-reads internal/releaseevidence/requirements.json or
// re-derives which producer records a profile needs — the preflight-side
// requiredness policy (missing FT87/FT88/FT71 records going red with an
// attributed message, and a conforming run going green) is already covered by
// internal/preflight/preflight_test.go's TestReleaseProfilesStayPendingInVerifyAndRedInPublish,
// TestBankPublishRequiresFT71Evidence, TestCommandPublishIsGreenStrictSuperset, and
// internal/preflight/review_fixes_test.go's
// TestConditionalRecordReasonFlowsThroughAuthoritativeEvidence — this test asserts
// only the new half: that publication composes on top of that verdict rather
// than re-deriving it.
func TestPublishSubmitRefusesUnlessApprovedIndexIsFullGreenPublishRunForProfile(t *testing.T) {
	for _, test := range []struct {
		name      string
		overrides map[string]string
		wantSub   string
	}{
		{
			name:      "verify-mode index",
			overrides: map[string]string{"mode": "verify"},
			wantSub:   "not a publish-mode run",
		},
		{
			name:      "focused scope index",
			overrides: map[string]string{"scope": "focused"},
			wantSub:   "focused preflight run",
		},
		{
			name:      "red status index",
			overrides: map[string]string{"status": "red"},
			wantSub:   "not green",
		},
		{
			name:      "mismatched profile index",
			overrides: map[string]string{"profile": "bank"},
			wantSub:   "requested profile",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			version := "9.9.20"
			root := copyPublicationScripts(t)
			ordered := releasePlanArtifacts(t, root, version)
			writeApprovedSet(t, root, ordered, nil)
			patchReleaseIndexAuthority(t, root, test.overrides)
			base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

			exitCode, output := runReleaseSubmit(t, root, version, base)
			if exitCode == 0 {
				t.Fatalf("submit succeeded against an unauthorized approved index (%s):\n%s", test.name, output)
			}
			if !strings.Contains(output, test.wantSub) {
				t.Fatalf("submit error did not attribute the authority failure (%s), want substring %q:\n%s", test.name, test.wantSub, output)
			}
			lines := requestLines(t, requestFile)
			if len(lines) != 0 {
				t.Fatalf("submit issued a registry call before verifying publish authority (%s):\n%s", test.name, strings.Join(lines, "\n"))
			}
			if _, err := readRawRecord(t, root); err == nil {
				t.Fatalf("submit wrote a publication record before verifying publish authority (%s)", test.name)
			}
		})
	}
}

// TestPublishSubmitProceedsWithConformingFullGreenPublishIndex is this row's
// green case: a conforming full green publish-mode index for the requested
// profile lets submit proceed to the registry exactly as every other
// submit test already exercises via writeApprovedSet's default fixture.
func TestPublishSubmitProceedsWithConformingFullGreenPublishIndex(t *testing.T) {
	version := "9.9.21"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil) // default index is already public/preflight/green/publish
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runReleaseSubmit(t, root, version, base)
	if exitCode != 0 {
		t.Fatalf("submit refused a conforming full green publish-mode index:\n%s", output)
	}
	lines := requestLines(t, requestFile)
	if len(lines) == 0 {
		t.Fatal("submit issued no registry calls at all against a conforming index")
	}
}
