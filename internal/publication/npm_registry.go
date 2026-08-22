package publication

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// NPMCLIRegistry is the public-npm adapter: it shells the real `npm` CLI in a
// scratch directory. It is runbook-only — the gate never exercises it, since
// the gate has no network egress and no publish credential. Auth material
// (an npm token, OIDC trust) lives in the ambient npm config/environment the
// CLI reads itself. This adapter never reads a credential value into memory,
// the record, or any evidence.
type NPMCLIRegistry struct {
	// Registry is the npm registry URL npm publishes against (defaults to the
	// npm CLI's own configured registry when empty).
	Registry string
	// Access is the --access value publish carries; empty leaves npm's own
	// default. Scoped @redbench/* packages need "public" or they publish
	// private.
	Access string
	// Provenance appends --provenance to publish. It is opt-in: a CI runner
	// can attest a build, an operator's laptop generally cannot.
	Provenance bool
}

func NewNPMCLIRegistry(registry string) *NPMCLIRegistry {
	return &NPMCLIRegistry{Registry: registry}
}

func (n *NPMCLIRegistry) registryArgs() []string {
	if n.Registry == "" {
		return nil
	}
	return []string{"--registry", n.Registry}
}

func (n *NPMCLIRegistry) run(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "npm", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("npm %v failed: %w: %s", args, err, stderr.String())
	}
	return stdout.Bytes(), nil
}

func (n *NPMCLIRegistry) Publish(ctx context.Context, name, version, tag string, tarball []byte) (string, error) {
	dir, err := os.MkdirTemp("", "bench-npm-publish-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, packageFileName(name, version))
	if err := os.WriteFile(file, tarball, 0o644); err != nil {
		return "", err
	}
	args := append([]string{"publish", file, "--tag", tag}, n.registryArgs()...)
	if n.Access != "" {
		args = append(args, "--access", n.Access)
	}
	if n.Provenance {
		args = append(args, "--provenance")
	}
	if _, err := n.run(ctx, dir, args...); err != nil {
		return "", err
	}
	integrity, live, err := n.Integrity(ctx, name, version)
	if err != nil {
		return "", err
	}
	if !live {
		return "", fmt.Errorf("npm publish did not make %s@%s live", name, version)
	}
	return integrity, nil
}

// minStagedNPMVersion and minStagedNodeVersion are the tool floors the public
// npm OIDC trusted-publishing staged flow requires. This precondition belongs
// only to this adapter. The fixture path (FixtureRegistry) never shells npm or
// node, so it must never depend on a locally installed tool version the gate
// cannot control.
const (
	minStagedNPMVersion  = "11.15.0"
	minStagedNodeVersion = "22.14.0"
)

func (n *NPMCLIRegistry) checkStagedToolPreconditions(ctx context.Context) error {
	npmVersion, err := n.toolVersion(ctx, "npm", "--version")
	if err != nil {
		return fmt.Errorf("could not determine npm version: %w", err)
	}
	if compareSemver(npmVersion, minStagedNPMVersion) < 0 {
		return fmt.Errorf("staged publication requires npm %s or newer (found %s)", minStagedNPMVersion, npmVersion)
	}
	nodeVersion, err := n.toolVersion(ctx, "node", "--version")
	if err != nil {
		return fmt.Errorf("could not determine node version: %w", err)
	}
	nodeVersion = strings.TrimPrefix(nodeVersion, "v")
	if compareSemver(nodeVersion, minStagedNodeVersion) < 0 {
		return fmt.Errorf("staged publication requires node %s or newer (found %s)", minStagedNodeVersion, nodeVersion)
	}
	return nil
}

func (n *NPMCLIRegistry) toolVersion(ctx context.Context, tool string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, tool, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// compareSemver compares two dotted major.minor.patch version strings
// (ignoring any pre-release/build suffix), returning -1/0/1 like strings.Compare.
func compareSemver(a, b string) int {
	pa, pb := semverParts(a), semverParts(b)
	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			if pa[i] < pb[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}

func semverParts(v string) [3]int {
	core := strings.SplitN(v, "-", 2)[0]
	fields := strings.SplitN(core, ".", 3)
	var out [3]int
	for i := 0; i < 3 && i < len(fields); i++ {
		n := 0
		for _, r := range fields[i] {
			if r < '0' || r > '9' {
				break
			}
			n = n*10 + int(r-'0')
		}
		out[i] = n
	}
	return out
}

func (n *NPMCLIRegistry) StageSubmit(ctx context.Context, name, version string, tarball []byte) (string, error) {
	// Public npm has no staged-submit primitive of its own; the OIDC trusted-
	// publishing flow performs an authenticated exchange the CLI drives
	// end-to-end. This adapter is construct-only for now. The state machine's
	// staged path is exercised against the fixture adapter, never this one.
	if err := n.checkStagedToolPreconditions(ctx); err != nil {
		return "", err
	}
	return "", fmt.Errorf("staged submission is not implemented for the public npm adapter")
}

func (n *NPMCLIRegistry) Approve(ctx context.Context, stageID string) error {
	if err := n.checkStagedToolPreconditions(ctx); err != nil {
		return err
	}
	return fmt.Errorf("approve is not implemented for the public npm adapter")
}

func (n *NPMCLIRegistry) Integrity(ctx context.Context, name, version string) (string, bool, error) {
	dir, err := os.MkdirTemp("", "bench-npm-integrity-")
	if err != nil {
		return "", false, err
	}
	defer os.RemoveAll(dir)
	args := append([]string{"view", name + "@" + version, "dist.integrity", "--json"}, n.registryArgs()...)
	out, err := n.run(ctx, dir, args...)
	if err != nil {
		// npm view exits non-zero (E404) when the version does not exist —
		// that is absent, not an error: no republish-blocking mismatch.
		return "", false, nil
	}
	var integrity string
	if err := json.Unmarshal(out, &integrity); err != nil {
		integrity = string(bytes.TrimSpace(out))
	}
	integrity = strings.TrimSpace(integrity)
	if integrity == "" {
		// The version exists (npm view succeeded) but reports no integrity —
		// a malformed/hostile registry response, not a live/not-live fact.
		return "", false, fmt.Errorf("npm registry reports %s@%s live with no integrity value", name, version)
	}
	return integrity, true, nil
}

func (n *NPMCLIRegistry) TagAdd(ctx context.Context, name, tag, version string) error {
	dir, err := os.MkdirTemp("", "bench-npm-tag-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	args := append([]string{"dist-tag", "add", name + "@" + version, tag}, n.registryArgs()...)
	_, err = n.run(ctx, dir, args...)
	return err
}

func (n *NPMCLIRegistry) TagRemove(ctx context.Context, name, tag string) error {
	dir, err := os.MkdirTemp("", "bench-npm-tag-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	args := append([]string{"dist-tag", "rm", name, tag}, n.registryArgs()...)
	_, err = n.run(ctx, dir, args...)
	return err
}

func (n *NPMCLIRegistry) Deprecate(ctx context.Context, name, version, message string) error {
	dir, err := os.MkdirTemp("", "bench-npm-deprecate-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	args := append([]string{"deprecate", name + "@" + version, message}, n.registryArgs()...)
	_, err = n.run(ctx, dir, args...)
	return err
}

// sriIntegrity computes the sha512-<base64> SRI string for data. It is the
// one source both adapters and the state machine use to compare local
// approved bytes against a registry's reported integrity.
func sriIntegrity(data []byte) string {
	sum := sha512.Sum512(data)
	return "sha512-" + base64.StdEncoding.EncodeToString(sum[:])
}
