//go:build bench_canary_native_proof

package releaseevidence

import "testing"

func TestNativeProofAuthorizationCanary(t *testing.T) {
	digest := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	target := targetEvidence{OS: "linux", Arch: "x64", Runner: "runner"}
	proof := nativeProofEvidence{SchemaVersion: 1, Target: "linux-x64", Runner: "runner", Status: "green", RebuiltSHA256: digest, BinarySHA256: digest, PackageSHA256: digest, ArchiveSHA256: digest, MuslStatus: "green", OperationsStatus: "green", StripStatus: "green", ToolsStatus: "green"}
	if !nativeProofMatches(proof, target, digest, digest, digest) {
		t.Fatal("valid authoritative native proof was rejected")
	}
	proof.BinarySHA256 = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if nativeProofMatches(proof, target, digest, digest, digest) {
		t.Fatal("authoritative native proof digest mutation passed")
	}
}
