package adopt

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// BrokerManifestName is the installation manifest that authenticates the promotion
// broker `bench worktree land` runs under. It sits beside the installed wrapper, and
// the wrapper's land route trusts only what it binds: path, version, platform, and
// executable digest. The install and repair owner publishes it together with the
// broker, so a repository executable can never become the landing owner.
const BrokerManifestName = "bench-broker.manifest"

// WriteBrokerManifest publishes the broker manifest beside the resolved wrapper,
// binding the currently running executable as the promotion broker. It returns the
// manifest path and the broker path it bound.
func WriteBrokerManifest(version string) (string, string, error) {
	broker, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("resolve promotion broker executable: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(broker); err == nil {
		broker = resolved
	}
	digest, err := fingerprintPath(broker)
	if err != nil {
		return "", "", fmt.Errorf("digest promotion broker %s: %w", broker, err)
	}
	path := filepath.Join(filepath.Dir(resolvedWrapper()), BrokerManifestName)
	content := "path\t" + broker + "\n" +
		"version\t" + version + "\n" +
		"platform\t" + brokerPlatform() + "\n" +
		"sha256\t" + digest + "\n"
	tmp, err := os.CreateTemp(filepath.Dir(path), ".bench-broker.")
	if err != nil {
		return "", "", fmt.Errorf("publish broker manifest at %s: %w", path, err)
	}
	name := tmp.Name()
	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return "", "", fmt.Errorf("publish broker manifest at %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return "", "", fmt.Errorf("publish broker manifest at %s: %w", path, err)
	}
	if err := os.Chmod(name, 0o644); err != nil {
		_ = os.Remove(name)
		return "", "", fmt.Errorf("publish broker manifest at %s: %w", path, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return "", "", fmt.Errorf("publish broker manifest at %s: %w", path, err)
	}
	return path, broker, nil
}

// brokerPlatform names the host this broker was installed for, in npm's os/cpu
// vocabulary. This is the one derivation of that fact: the wrapper's land route reads
// the value the manifest carries and never computes a second copy to compare it
// against, because the manifest digest binds the exact executable instead.
func brokerPlatform() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	return runtime.GOOS + "-" + arch
}

// ReadBrokerManifest parses one broker manifest into its four bound fields.
func ReadBrokerManifest(path string) (map[string]string, error) {
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
	for _, key := range []string{"path", "version", "platform", "sha256"} {
		if fields[key] == "" {
			return nil, fmt.Errorf("broker manifest %s does not bind %s", path, key)
		}
	}
	return fields, nil
}
