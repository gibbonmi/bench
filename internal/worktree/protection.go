package worktree

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// checkoutDirectory is the identity adapter of a checkout's Git administrative directory.
// Reset and clean act on the working tree and never reach that directory, so a generation
// bound to it survives both. A validator cannot ask Git for the directory, because a bounded
// Git launch registers through that same validator. The store resides under the protection
// root, outside the checkout.
var checkoutDirectory = protection.Directory{Kind: protection.KindWorktree, Noun: "worktree administrative directory"}

// ValidateProtection is the worktree identity adapter.
func ValidateProtection(d protection.Descriptor) error { return checkoutDirectory.Validate(d) }

// Descriptor names the protection resource of the checkout at dir without its generation.
func Descriptor(dir string) (protection.Descriptor, error) {
	admin, err := git.AdminDir(dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return checkoutDirectory.Describe(admin)
}

// childResources is the resource vector of one child in the checkout at dir: the vector that
// environ carries, then the checkout's own generation. The child registers on that union
// before it starts. A malformed inherited vector prevents the launch. A process with no
// installed protocol opens no generation and passes the inherited vector on unchanged.
func childResources(dir string, environ []string) ([]protection.Descriptor, error) {
	inherited, err := protection.ParseVector(environValue(environ, protection.Env))
	if err != nil {
		return nil, err
	}
	p := protection.Installed()
	if p == nil {
		return inherited, nil
	}
	d, err := Descriptor(dir)
	if err == nil {
		d, err = p.OpenResource(d)
	}
	if err != nil {
		return nil, fmt.Errorf("protect worktree %s: %w", dir, err)
	}
	return append(inherited, d), nil
}

// environValue is the last value environ assigns to name, which is the value a child reads.
func environValue(environ []string, name string) string {
	value := ""
	for _, entry := range environ {
		if after, ok := strings.CutPrefix(entry, name+"="); ok {
			value = after
		}
	}
	return value
}
