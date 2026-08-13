package releaseevidence

import (
	"encoding/json"
	"testing"
)

func TestReleaseIndexBindsComponentManifestDigest(t *testing.T) {
	const digest = "component-manifest-digest"
	encoded, err := canonicalJSON(Index{Artifacts: []artifactEvidence{{Name: "bench.tgz", ComponentDigest: digest}}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Artifacts []struct {
			ComponentDigest string `json:"component_manifest_sha256"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if got := decoded.Artifacts[0].ComponentDigest; got != digest {
		t.Fatalf("component manifest digest = %q, want %q", got, digest)
	}
}
