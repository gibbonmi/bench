package worktree

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/worktree/reclaimpolicy"
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
	return describeCheckout(new(git.Run), dir)
}

// describeCheckout is Descriptor with the location answers of run.
func describeCheckout(run *git.Run, dir string) (protection.Descriptor, error) {
	admin, err := run.AdminDir(dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return checkoutDirectory.Describe(admin)
}

// ChildResources is the resource vector of a child that another package launches in the
// checkout at dir, as childResources states it.
func ChildResources(dir string, environ []string) ([]protection.Descriptor, error) {
	return childResources(dir, environ)
}

// childResources is the resource vector of one child in the checkout at dir: the vector that
// environ carries, then the checkout's own generation. The child registers on that union
// before it starts. A malformed inherited vector prevents the launch. A process with no
// installed protocol opens no generation and passes the inherited vector on unchanged.
func childResources(dir string, environ []string) ([]protection.Descriptor, error) {
	inherited, err := protection.InheritedFrom(environ)
	if err != nil {
		return nil, err
	}
	p := protection.Installed()
	if p == nil {
		return inherited, nil
	}
	d, err := Descriptor(dir)
	if err == nil {
		d, err = p.Reopen(d)
	}
	if err != nil {
		return nil, fmt.Errorf("protect worktree %s: %w", dir, err)
	}
	return append(inherited, d), nil
}

// checkoutObserver is the protocol that observes the users of a checkout. Observation needs
// only the worktree identity adapter, so a process with no installed protocol still sees each
// user that another process registered.
var checkoutObserver = &protection.Protocol{Validators: protection.Validators{protection.KindWorktree: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}

// checkoutUnprotected is the one protected-user predicate of a checkout. Release, reuse, and
// cleanup apply ask it before they reset, clean, or remove the checkout at path. It seals
// the checkout's generation and observes every user, so no new user joins a generation that
// a safe answer retires. A nil answer means that no user of the checkout remains unresolved.
// An absent store, and a store that holds no generation, name no user. An unresolved user, a
// failed observation, or a checkout that it cannot describe keeps the checkout. It reads Git
// through run.
func checkoutUnprotected(run *git.Run, path string) error {
	d, err := describeCheckout(run, path)
	if err != nil {
		return fmt.Errorf("worktree %s %w: %w", path, protection.ErrRetained, err)
	}
	disposition, err := checkoutObserver.Recover(d.Store, nil)
	switch {
	case protection.NoGeneration(disposition, err):
		return nil
	case err != nil:
		return fmt.Errorf("worktree %s %w: %w", path, protection.ErrRetained, err)
	case disposition.State != protection.Safe:
		return fmt.Errorf("worktree %s %w: %w", path, protection.ErrRetained, disposition.Unresolved())
	}
	return nil
}

// CheckoutUnprotected is checkoutUnprotected for an owner outside the package, such as a
// shift that rolls back, removes its scratch, or recovers the checkout at path.
func CheckoutUnprotected(path string) error {
	return checkoutUnprotected(new(git.Run), path)
}

// memberUsers is the protected-user fact of a pool member whose administrative directory
// target no longer exists. The member's repository is gone, so its store cannot derive from
// the checkout: the fact grades each worktree store whose generation records target. A
// failed observation leaves the fact unknown, which retains the member.
func memberUsers(target string) (reclaimpolicy.Users, string) {
	disposition, err := checkoutDirectory.ObserveRecorded(target, bounds.ControlRecordLimit)
	switch unresolved := disposition.Unresolved(); {
	case err != nil:
		return reclaimpolicy.UsersUnknown, err.Error()
	case unresolved != nil:
		return reclaimpolicy.UsersUnresolved, unresolved.Error()
	}
	return reclaimpolicy.UsersResolved, ""
}

// protectionEvidence is the plan evidence of the protection of the checkout whose
// administrative directory is admin: each entry of its store, read without a seal. A user
// that joins or resolves after the plan changes the evidence, and apply rechecks the users
// through checkoutUnprotected under the cleanup exclusion. An absent store, an undescribed
// checkout, and an unreadable entry each end the evidence with the error that stopped it.
func protectionEvidence(admin string) []byte {
	d, err := checkoutDirectory.Describe(admin)
	var parts [][]byte
	if err == nil {
		err = filepath.WalkDir(d.Store, func(path string, _ fs.DirEntry, walkErr error) error {
			evidence, fileErr := fileEvidence(path)
			parts = append(parts, []byte(path), evidence)
			return errors.Join(walkErr, fileErr)
		})
	}
	return canonicalParts(append(parts, []byte(fmt.Sprint(err)))...)
}
