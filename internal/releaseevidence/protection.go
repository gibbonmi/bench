package releaseevidence

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// evidenceDirectory is the identity adapter of a checkout's dist directory. Dist holds the
// dist/preflight target, its stages, and their owner markers. The exchange replaces the target
// and deletes stages, but it never replaces dist, so the identity of dist binds one generation
// to this checkout across every exchange. A validator cannot ask Git for the common directory,
// because a bounded Git launch registers through that same validator. The store resides under
// the protection root, outside the target and every stage.
var evidenceDirectory = protection.Directory{
	Kind:   protection.KindReleaseEvidence,
	Noun:   "release evidence directory",
	Prefix: "dist",
	Class:  "release output directory",
}

// ValidateProtection is the release-evidence identity adapter.
func ValidateProtection(d protection.Descriptor) error { return evidenceDirectory.Validate(d) }

// Descriptor names the release-evidence resource of root without its generation. Dist must exist.
func Descriptor(root string) (protection.Descriptor, error) {
	return evidenceDirectory.Describe(filepath.Join(root, "dist"))
}

// Protect opens the release-evidence generation of root before any phase launches, and it
// creates dist when dist is absent. A generation that an earlier run left first passes the
// observation that publication consumes, so its unresolved users refuse this run. A process
// inside a user of the generation joins it without that observation, because it cannot observe
// its own enclosing user complete. A process with no installed protocol opens no generation.
func Protect(root string) (protection.Generation, error) {
	dist, err := evidenceDist(root)
	if err != nil {
		return protection.Generation{}, err
	}
	generation, enclosed, err := openEvidence(dist)
	if err != nil || enclosed {
		return generation, err
	}
	if err := generation.Release(evidenceName(dist), nil); err != nil {
		return protection.Generation{}, err
	}
	generation, _, err = openEvidence(dist)
	return generation, err
}

// publishEvidence runs replace, which replaces or deletes evidence under dist, only when
// release-evidence protection permits it. The owner seals the generation and observes every
// user, and only a safe generation runs replace, under the protocol lock. A process inside a
// user of the generation runs replace without the seal: the enclosing owner observes the
// generation after this process returns.
func publishEvidence(dist string, replace func() error) error {
	generation, enclosed, err := openEvidence(dist)
	if err != nil {
		return err
	}
	if enclosed {
		return replace()
	}
	return generation.Release(evidenceName(dist), replace)
}

// openEvidence joins the open release-evidence generation of dist. A sealed generation first
// passes recovery, which retires it once every user resolves. Enclosed reports that the
// inherited vector names the generation, so this process runs inside one of its users.
func openEvidence(dist string) (_ protection.Generation, enclosed bool, _ error) {
	p := protection.Installed()
	if p == nil {
		return protection.Generation{}, false, nil
	}
	d, err := evidenceDirectory.Describe(dist)
	if err == nil {
		d, err = p.Reopen(d)
	}
	if err != nil {
		return protection.Generation{}, false, fmt.Errorf("protect %s: %w", evidenceName(dist), err)
	}
	generation := protection.Generation{Protocol: p, Resource: d}
	inherited, err := protection.Inherited()
	if err != nil {
		return protection.Generation{}, false, err
	}
	return generation, slices.Contains(inherited, generation.Resource), nil
}

func evidenceName(dist string) string { return "release evidence " + dist }

// evidenceDist creates root's dist directory when it is absent and refuses a dist that is not
// a real directory.
func evidenceDist(root string) (string, error) {
	dist := filepath.Join(root, "dist")
	if info, err := os.Lstat(dist); err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		return "", fmt.Errorf("dist output target is not a real directory")
	} else if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if err := os.MkdirAll(dist, 0o755); err != nil {
		return "", err
	}
	return dist, nil
}
