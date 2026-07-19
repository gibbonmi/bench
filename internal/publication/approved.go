package publication

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
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
