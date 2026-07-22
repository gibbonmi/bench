package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// upgrade.go is the consumer's supported route onto a newer kit, replacing the
// maintainer-only /bench-update-kit surface the payload allowlist withholds. It owns no
// write path: the plan it prints is a read-only summary, and the tree, the manifest
// rows, and the manifest's #kit version header all move together inside the same
// transactionalLink that bench link uses — so an interrupt can never strand a new tree
// under an old version claim.

// Upgrade compares the linked repo's pinned kit version against the installed one,
// prints the plan, and applies the relink. --check reports the plan and applies
// nothing; an equal version is a definitive no-op success; a lower installed version is
// a downgrade and fails closed until --force performs it.
func Upgrade(args []string, stdout, stderr io.Writer, version string) int {
	check, force := false, false
	for _, arg := range args {
		switch arg {
		case "--check":
			check = true
		case "--force":
			force = true
		default:
			fmt.Fprintln(stderr, "usage: bench upgrade [--check] [--force]")
			return 2
		}
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	manifest, code := pinnedManifest(filepath.Join(root, ".bench", "link-manifest.tsv"), stderr)
	if code != 0 {
		return code
	}
	pinned := manifest.KitVersion
	kit := kitDir()
	plan, err := buildLinkPlan(kit)
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("consumer payload allowlist rejected", err.Error()))
		return 1
	}
	added, changed, removed := upgradePlanCounts(manifest, plan)
	block, err := toon.TableTyped("upgrade", []string{"from", "to", "added", "changed", "removed"},
		[][]any{{pinned, version, added, changed, removed}})
	if err != nil {
		fmt.Fprintln(stderr, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, block)

	direction := compareKitVersions(version, pinned)
	// The downgrade refusal precedes the --check return: a dry run that reported a plan
	// the applying run would refuse would be a promise the command cannot keep.
	if direction < 0 && !force {
		fmt.Fprintln(stderr, toon.Errorf("installed kit is older than the linked kit", "re-run 'bench upgrade --force' to accept the downgrade"))
		return 1
	}
	if check || direction == 0 {
		return 0
	}
	result, _ := transactionalLink(root, kit, "copy", version, plan, stdout, stderr)
	return result
}

// pinnedManifest reads the linked repo's manifest and its #kit header, distinguishing
// the two absent states a consumer can be in: no manifest at all means the repo was
// never linked and the remedy is bench link, while a manifest that exists but carries
// no version is a malformed state a relink has to repair rather than a missing linkage.
func pinnedManifest(path string, stderr io.Writer) (Manifest, int) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Fprintln(stderr, toon.Errorf("repository is not Bench-linked", "run 'bench link' first"))
		return Manifest{}, 1
	}
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("link manifest unreadable", err.Error()))
		return Manifest{}, 1
	}
	if info.Size() == 0 {
		fmt.Fprintln(stderr, toon.Errorf("link manifest is empty", "re-run 'bench link' to rewrite the manifest"))
		return Manifest{}, 1
	}
	manifest, err := ReadManifest(path)
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("link manifest unreadable", err.Error()))
		return Manifest{}, 1
	}
	if manifest.KitVersion == "" {
		fmt.Fprintln(stderr, toon.Errorf("link manifest carries no #kit version header", "re-run 'bench link' to restamp the manifest"))
		return Manifest{}, 1
	}
	return manifest, 0
}

// upgradePlanCounts summarizes what applying plan to a repo owning manifest would do: a
// plan entry the manifest does not own is added, a manifest row the plan no longer
// carries is removed, and an owned file entry whose kit-side bytes no longer match the
// recorded hash is changed. Only file entries can drift — an adapter symlink and an
// inline entry are generated from the destination path and from constant content — so
// this reads the same inputs transactionalLink reconciles without deciding anything a
// second time. CLAUDE.md is manifest-owned but conditionally, never planned, so it is
// not a withdrawal.
func upgradePlanCounts(manifest Manifest, plan []planEntry) (added, changed, removed int) {
	planned := make(map[string]bool, len(plan))
	for _, entry := range plan {
		planned[entry.rel] = true
		old := manifest.Hash(entry.rel)
		if old == "" {
			added++
			continue
		}
		if entry.kind != "file" {
			continue
		}
		if now, err := fingerprintPath(entry.src); err != nil || now != old {
			changed++
		}
	}
	for _, row := range manifest.Rows() {
		if !planned[row.rel] && row.rel != "CLAUDE.md" {
			removed++
		}
	}
	return added, changed, removed
}

// compareKitVersions orders two kit versions by their numeric release components,
// ignoring any prerelease or build suffix: -1 when a precedes b, 1 when it follows, 0
// when they are the same release. A version neither side can parse (an unstamped "dev"
// build, a hand-edited header) is never reported as a downgrade — only an equal string
// is equal — so a malformed stamp cannot make upgrade refuse an install it should do.
func compareKitVersions(a, b string) int {
	if a == b {
		return 0
	}
	left, leftOK := releaseComponents(a)
	right, rightOK := releaseComponents(b)
	if !leftOK || !rightOK {
		return 1
	}
	for i := range left {
		if left[i] != right[i] {
			if left[i] < right[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}

func releaseComponents(version string) ([3]int, bool) {
	var out [3]int
	core := version
	if i := strings.IndexAny(core, "-+"); i >= 0 {
		core = core[:i]
	}
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return out, false
	}
	for i, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return out, false
		}
		out[i] = value
	}
	return out, true
}
