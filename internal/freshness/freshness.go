// Package freshness verifies that a Bench executable was built from the current sources.
package freshness

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const sealSchema = 1

type seal struct {
	Schema     int    `json:"schema"`
	Sources    string `json:"sources"`
	Executable string `json:"executable"`
}

// DeclaresBuildInputs reports whether root declares Go build inputs, which is what makes
// a rebuildable dev executable possible. Presence decides, not content. A manifest that is
// a broken link or a special file routes to Verify instead, whose reading discipline
// refuses what it cannot trust. Nothing unreadable is ever read as an authoritative
// absence.
func DeclaresBuildInputs(root string) bool {
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(auxiliaryInputsManifest)))
	return !errors.Is(err, os.ErrNotExist)
}

// Digest returns the deterministic content digest of Bench's local build inputs.
func Digest(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	paths, err := buildInputs(root)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	for _, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return "", err
		}
		contents, err := regularContents(path)
		if err != nil {
			return "", fmt.Errorf("read build input %q: %w", rel, err)
		}
		name := filepath.ToSlash(rel)
		fmt.Fprintf(hash, "%d:%s%d:", len(name), name, len(contents))
		if _, err := hash.Write(contents); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// BuildInputs returns the repository-relative, slash-separated, sorted paths that Digest
// hashes for root.
func BuildInputs(root string) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	paths, err := buildInputs(root)
	if err != nil {
		return nil, err
	}
	relative := make([]string, len(paths))
	for i, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		relative[i] = filepath.ToSlash(rel)
	}
	return relative, nil
}

// SealDigests returns the source and executable digests recorded in executable's
// published seal.
func SealDigests(executable string) (sources, digest string, err error) {
	data, err := secureContents(sealPath(executable), false)
	if err != nil {
		return "", "", fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	stored, err := parseSeal(data)
	if err != nil {
		return "", "", fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	return stored.Sources, stored.Executable, nil
}

// ExecutableDigest returns the digest a seal records for the executable at path.
func ExecutableDigest(path string) (string, error) {
	binary, err := secureContents(path, true)
	if err != nil {
		return "", err
	}
	return digestBytes(binary), nil
}

func parseSeal(data []byte) (seal, error) {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var stored seal
	if err := decoder.Decode(&stored); err != nil {
		return seal{}, fmt.Errorf("malformed: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return seal{}, fmt.Errorf("malformed trailing data")
	}
	if stored.Schema != sealSchema || !isDigest(stored.Sources) || !isDigest(stored.Executable) {
		return seal{}, fmt.Errorf("malformed contents")
	}
	return stored, nil
}

// writeSeal promotes data into path through a staging file, and it reports that file's
// name to track for as long as it is in flight. A termination answered mid-write ends the
// process before any deferred cleanup here can run. The staging file therefore needs an
// owner that outlives this call.
func writeSeal(path string, data []byte, track func(string)) (err error) {
	temporary, err := os.CreateTemp(filepath.Dir(path), sealTemporaryPattern(path))
	if err != nil {
		return err
	}
	track(temporary.Name())
	defer func() {
		if err != nil {
			_ = os.Remove(temporary.Name())
		}
		track("")
	}()
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return replacePublicationFile(temporary.Name(), path)
}

func sealPath(executable string) string { return executable + ".seal" }

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func isDigest(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
