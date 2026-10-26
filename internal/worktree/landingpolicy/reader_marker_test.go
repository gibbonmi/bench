// The resume marker decision over a published tree that holds a reader record.
package landingpolicy

import "testing"

// TestResumeMarkerAcceptsAReaderLanding is NG59. A green reader set published the tree, so
// the marker stays at the last whole green: the resume accepts it and advances nothing,
// wherever the destination and the marker stand.
func TestResumeMarkerAcceptsAReaderLanding(t *testing.T) {
	for _, tc := range []struct {
		name          string
		atPublished   bool
		markerPresent bool
	}{
		{"destination at the published landing", true, true},
		{"destination moved on and the marker behind", false, true},
		{"no marker", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			facts := MarkerFacts{
				DestinationAtPublished:  tc.atPublished,
				MarkerPresent:           tc.markerPresent,
				MarkerReachesPublished:  func() bool { return false },
				PublishedReaderRecorded: true,
			}
			if got := ResumeMarker(facts); got != MarkerAccept {
				t.Fatalf("ResumeMarker = %d, want accept (%d)", got, MarkerAccept)
			}
		})
	}
}
