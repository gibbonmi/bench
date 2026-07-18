package releaseevidence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

type targetEvidence struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	GOOS   string `json:"goos"`
	GOArch string `json:"goarch"`
	Runner string `json:"runner"`
}

type phaseEvidence struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Digest string `json:"record_sha256"`
}

type artifactEvidence struct {
	Name            string `json:"name"`
	Target          string `json:"target"`
	Size            int64  `json:"size"`
	SHA256          string `json:"sha256"`
	ComponentDigest string `json:"component_manifest_sha256"`
	SBOMDigest      string `json:"sbom_sha256"`
	InventoryDigest string `json:"inventory_sha256"`
}

type reproducibilityEvidence struct {
	SchemaVersion int                       `json:"schema_version"`
	Status        string                    `json:"status"`
	Builds        int                       `json:"builds"`
	Artifacts     []reproducibilityArtifact `json:"artifacts"`
}

type reproducibilityArtifact struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
	Match  bool   `json:"match"`
}

type nativeProofEvidence struct {
	SchemaVersion int    `json:"schema_version"`
	Target        string `json:"target"`
	Runner        string `json:"runner"`
	Status        string `json:"status"`
	RebuiltSHA256 string `json:"rebuilt_sha256"`
	BinarySHA256  string `json:"binary_sha256"`
	PackageSHA256 string `json:"package_sha256"`
	ArchiveSHA256 string `json:"archive_sha256"`
	MuslStatus    string `json:"musl_status"`
}

type toolchainEvidence struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Flags   []string `json:"flags"`
}

type releaseIdentity struct {
	Tag              string `json:"tag,omitempty"`
	PackageVersion   string `json:"package_version,omitempty"`
	SourceCommit     string `json:"source_commit,omitempty"`
	BinaryVersion    string `json:"binary_version,omitempty"`
	ChangelogHeading string `json:"changelog_heading,omitempty"`
	Toolchain        string `json:"toolchain,omitempty"`
}

type Index struct {
	SchemaVersion   int                     `json:"schema_version"`
	Mode            Mode                    `json:"mode"`
	Scope           Scope                   `json:"scope"`
	Profile         Profile                 `json:"profile,omitempty"`
	Status          Status                  `json:"status"`
	Identity        releaseIdentity         `json:"identity"`
	RollbackTarget  string                  `json:"rollback_target"`
	Flags           []string                `json:"flags"`
	Toolchains      []toolchainEvidence     `json:"toolchains"`
	Inputs          []evidenceDigest        `json:"inputs"`
	Targets         []targetEvidence        `json:"targets"`
	Phases          []phaseEvidence         `json:"phases"`
	Requirements    []RequirementStatus     `json:"requirements"`
	Artifacts       []artifactEvidence      `json:"artifacts"`
	Reproducibility reproducibilityEvidence `json:"reproducibility"`
	NativeProofs    []nativeProofEvidence   `json:"native_proofs"`
}

var indexEncoder = canonicalJSON

func SetIndexEncoderForTesting(encoder func(any) ([]byte, error)) func() {
	previous := indexEncoder
	indexEncoder = encoder
	return func() { indexEncoder = previous }
}

func canonicalJSON(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func deriveChecksums(index Index) []byte {
	var out bytes.Buffer
	artifacts := append([]artifactEvidence(nil), index.Artifacts...)
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Name < artifacts[j].Name })
	for _, artifact := range artifacts {
		fmt.Fprintf(&out, "%s  %s\n", artifact.SHA256, artifact.Name)
	}
	return out.Bytes()
}

func releaseIdentityFrom(identity Identity) releaseIdentity {
	return releaseIdentity{Tag: stringPointer(identity.Tag), PackageVersion: stringPointer(identity.PackageVersion), SourceCommit: stringPointer(identity.SourceCommit), BinaryVersion: stringPointer(identity.BinaryVersion), ChangelogHeading: stringPointer(identity.ChangelogHeading), Toolchain: stringPointer(identity.Toolchain)}
}

func stringPointer(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
