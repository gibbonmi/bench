// Package brokermanifest publishes and reads the installation manifest that
// authenticates the promotion broker `bench worktree land` runs under. It is a leaf, so
// the installer, the stamped build, and the doctor row all bind the same four fields
// without an import cycle between them.
package brokermanifest

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// Name is the manifest file name. It sits beside the resolved wrapper, and the wrapper's
// land route trusts only what it binds: path, version, platform, and executable digest.
const Name = "bench-broker.manifest"

// TemporaryPrefix matches the start of the staging file name that the durable replacement
// creates for Name, so a sweep of a publication's leftovers can find a leaked staging file.
// A replacement test pins the value against the staging name.
const TemporaryPrefix = "." + Name + "."

// boundFields are the four values one manifest carries.
var boundFields = []string{"path", "version", "platform", "sha256"}

// Write publishes the manifest in dir, binding broker as the promotion broker for
// version. It returns the manifest path and the symlink-resolved broker it bound.
func Write(dir, broker, version string) (string, string, error) {
	return write(durablefile.Native(), dir, broker, version)
}

// write is Write through the replacement operations of replacer.
func write(replacer durablefile.Replacer, dir, broker, version string) (string, string, error) {
	if resolved, err := filepath.EvalSymlinks(broker); err == nil {
		broker = resolved
	}
	digest, err := Digest(broker)
	if err != nil {
		return "", "", fmt.Errorf("digest promotion broker %s: %w", broker, err)
	}
	path := filepath.Join(dir, Name)
	content := "path\t" + broker + "\n" +
		"version\t" + version + "\n" +
		"platform\t" + Platform() + "\n" +
		"sha256\t" + digest + "\n"
	if err := replacer.Replace(path, []byte(content), 0o644); err != nil {
		return "", "", fmt.Errorf("publish broker manifest at %s: %w", path, err)
	}
	return path, broker, nil
}

// Read parses one broker manifest into its four bound fields.
func Read(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fields := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "\t")
		if !ok {
			return nil, fmt.Errorf("malformed broker manifest line %q in %s", line, path)
		}
		fields[key] = value
	}
	for _, key := range boundFields {
		if fields[key] == "" {
			return nil, fmt.Errorf("broker manifest %s does not bind %s", path, key)
		}
	}
	return fields, nil
}

// Digest returns the content digest of one regular file, in the spelling the land
// route's file_sha256 produces. A manifest binds this value, and every reader compares
// against it rather than against a second derivation of the broker's identity.
func Digest(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("not a regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Platform names the host this broker was installed for, in npm's os/cpu vocabulary.
// This is the one derivation of that fact: the wrapper's land route reads the value the
// manifest carries and never computes a second copy to compare it against, because the
// manifest digest binds the exact executable instead.
func Platform() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	return runtime.GOOS + "-" + arch
}
