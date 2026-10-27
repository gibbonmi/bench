package backgroundstate

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// claimDir is the directory of the claims in the git common directory. It sits beside the
// evidence store, so every worktree of the repository reads one claim for each grade key.
const claimDir = "bench-gate-claims"

// Claim is the hold of a running grade on its grade key. It names the worktree of the grade,
// the pid of its owner, and the boot in which that owner runs.
type Claim struct {
	Key      string `json:"key"`
	Worktree string `json:"worktree"`
	Owner    int    `json:"owner_pid"`
	Boot     string `json:"boot_id"`
}

// ClaimFor answers the claim of this process on key for worktree in the current boot.
func ClaimFor(key, worktree string) (Claim, error) {
	boot, err := protection.BootID()
	if err != nil {
		return Claim{}, err
	}
	return Claim{Key: key, Worktree: worktree, Owner: os.Getpid(), Boot: boot}, nil
}

// WriteClaim installs claim as the claim of its key in the git common directory gitdir.
func WriteClaim(replacer durablefile.Replacer, gitdir string, claim Claim) error {
	if err := claim.valid(); err != nil {
		return err
	}
	if err := ensureDir(filepath.Join(gitdir, claimDir)); err != nil {
		return err
	}
	return writeRecord(replacer, ClaimPath(gitdir, claim.Key), claim)
}

// ReadClaim reads the claim of key in the git common directory gitdir. An absent claim is
// the empty claim.
func ReadClaim(gitdir, key string) (Claim, error) {
	if err := pathName("key", key); err != nil {
		return Claim{}, err
	}
	var claim Claim
	err := readRecord(ClaimPath(gitdir, key), &claim, func() error {
		if claim.Key != key {
			return errors.New("the claim names another key")
		}
		return claim.valid()
	})
	if err != nil {
		return Claim{}, err
	}
	return claim, nil
}

// Live reports whether the owner of claim still runs: the claim names the current boot and
// its owner pid exists. A pid of an earlier boot can name a new process, so it is not live.
func (c Claim) Live() (bool, error) {
	boot, err := protection.BootID()
	if err != nil {
		return false, err
	}
	if c.Boot != boot {
		return false, nil
	}
	return subprocess.Alive(c.Owner)
}

func (c Claim) valid() error {
	if err := pathName("key", c.Key); err != nil {
		return err
	}
	if err := required("worktree", c.Worktree, "boot_id", c.Boot); err != nil {
		return err
	}
	if c.Owner <= 1 {
		return errors.New("the owner_pid field names no one process")
	}
	return nil
}

// ClaimPath answers the path of the claim of key in the git common directory gitdir.
func ClaimPath(gitdir, key string) string {
	return filepath.Join(gitdir, claimDir, key)
}
