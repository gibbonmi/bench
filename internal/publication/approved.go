package publication

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ApprovedPackage is one package from the locally verified approved release
// directory: its registry name, the exact file on disk, and its bytes' SRI
// integrity, checked against dist/preflight/SHA256SUMS before any publish call.
type ApprovedPackage struct {
	Name      string
	Version   string
	Kind      string // wrapper|platform
	FilePath  string
	SHA256    string
	Integrity string
}

// VerifyApprovedSet locally verifies the complete immutable release set —
// every platform package plus the wrapper named by release-plan.mjs — against
// dist/preflight/SHA256SUMS, and returns the release-index digest alongside
// the verified packages in release-plan order. It never rewrites
// dist/preflight/ or package bytes; this is a read-only check.
func VerifyApprovedSet(root, version string) (releaseIndexSHA256 string, packages []ApprovedPackage, err error) {
	indexPath := filepath.Join(root, "dist", "preflight", "release-index.json")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return "", nil, fmt.Errorf("approved release index is unreadable: %w", err)
	}
	sum := sha256.Sum256(indexData)
	releaseIndexSHA256 = hex.EncodeToString(sum[:])

	sums, err := readSHA256SUMS(filepath.Join(root, "dist", "preflight", "SHA256SUMS"))
	if err != nil {
		return "", nil, err
	}

	planDigests, err := releaseIndexArtifactDigests(indexData)
	if err != nil {
		return "", nil, err
	}

	records, err := PublicationSet(root, version)
	if err != nil {
		return "", nil, err
	}
	if len(records) == 0 {
		return "", nil, fmt.Errorf("release plan named no publication artifacts for version %s", version)
	}

	for _, record := range records {
		filePath := filepath.Join(root, "dist", "artifacts", record.Name)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", nil, fmt.Errorf("approved artifact %s is unreadable: %w", record.Name, err)
		}
		digestBytes := sha256.Sum256(data)
		digest := hex.EncodeToString(digestBytes[:])
		want, ok := sums[record.Name]
		if !ok {
			return "", nil, fmt.Errorf("SHA256SUMS names no digest for approved artifact %s", record.Name)
		}
		if digest != want {
			return "", nil, fmt.Errorf("approved artifact %s does not match its SHA256SUMS digest", record.Name)
		}
		if planDigest, named := planDigests[record.Name]; named && planDigest != digest {
			return "", nil, fmt.Errorf("approved artifact %s has drifted from the release plan: release-index.json recorded digest %s but the approved set now has %s", record.Name, planDigest, digest)
		}
		name := "redbench"
		if record.Kind == "platform" {
			name = "@redbench/" + record.Target
		}
		packages = append(packages, ApprovedPackage{
			Name:      name,
			Version:   version,
			Kind:      record.Kind,
			FilePath:  filePath,
			SHA256:    digest,
			Integrity: sriIntegrity(data),
		})
	}
	return releaseIndexSHA256, packages, nil
}

// releaseIndexAuthority is the narrow slice of dist/preflight/release-index.json
// publication needs to compose the preflight publish-mode authority: it never
// re-reads internal/releaseevidence/requirements.json or re-derives which
// producer records a profile requires. That requiredness policy is owned
// entirely by internal/releasepreflight; publication trusts only the index's already-
// computed verdict.
type releaseIndexAuthority struct {
	Mode    string `json:"mode"`
	Scope   string `json:"scope"`
	Profile string `json:"profile"`
	Status  string `json:"status"`
}

// VerifyPublishAuthority refuses to proceed unless the approved release
// directory's release-index.json is a full (non-focused) green publish-mode
// preflight run whose recorded profile matches profile. It is the one gate
// bench release submit composes on top of the preflight authority; it does no
// registry I/O and must run before any registry call.
func VerifyPublishAuthority(root, profile string) error {
	indexPath := filepath.Join(root, "dist", "preflight", "release-index.json")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("approved release index is unreadable: %w", err)
	}
	var authority releaseIndexAuthority
	if err := json.Unmarshal(indexData, &authority); err != nil {
		return fmt.Errorf("approved release index is malformed: %w", err)
	}
	if authority.Mode != "publish" {
		return fmt.Errorf("approved release index is a %q run, not a publish-mode run; run bench release-preflight --mode publish --profile %s first", authority.Mode, profile)
	}
	if authority.Scope == "focused" {
		return fmt.Errorf("approved release index is a focused preflight run, which never authorizes publication; run a full bench release-preflight --mode publish --profile %s first", profile)
	}
	if authority.Status != "green" {
		return fmt.Errorf("approved release index is %q, not green; publication requires a full green publish-mode preflight for profile %s", authority.Status, profile)
	}
	if authority.Profile != profile {
		return fmt.Errorf("approved release index authorized profile %q, but this run requested profile %q", authority.Profile, profile)
	}
	return nil
}

// releaseIndexArtifactDigests reads the per-artifact sha256 the release-index
// records for each named artifact ("artifacts": [{"name":..., "sha256":...}]),
// the immutable plan preflight froze. It never re-derives an artifact digest —
// only compares the plan's own recorded value against what SHA256SUMS/dist/
// artifacts now say, so a release-index.json + SHA256SUMS pair that has
// drifted apart (e.g. a tampered SHA256SUMS matching a swapped-in artifact)
// is caught even though each file individually parses and is self-consistent.
// An index that names no artifacts array at all (an older or synthetic index)
// is tolerated with no drift check — the plan simply made no claim to check.
func releaseIndexArtifactDigests(indexData []byte) (map[string]string, error) {
	var parsed struct {
		Artifacts []struct {
			Name   string `json:"name"`
			SHA256 string `json:"sha256"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(indexData, &parsed); err != nil {
		return nil, fmt.Errorf("approved release index is malformed: %w", err)
	}
	digests := make(map[string]string, len(parsed.Artifacts))
	for _, artifact := range parsed.Artifacts {
		if artifact.Name != "" {
			digests[artifact.Name] = artifact.SHA256
		}
	}
	return digests, nil
}

func readSHA256SUMS(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("approved SHA256SUMS is unreadable: %w", err)
	}
	defer file.Close()
	sums := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed SHA256SUMS line %q", line)
		}
		sums[parts[1]] = parts[0]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("approved SHA256SUMS is unreadable: %w", err)
	}
	return sums, nil
}
