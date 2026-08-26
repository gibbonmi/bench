// Package poolkey derives the repository identity that names a Bench pool directory.
// The worktree pool and the exec census both key their directories by it, and each
// package needs the other, so the one derivation sits below both.
package poolkey

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// Canonical returns the primary checkout root for any root inside the repository.
// A linked worktree's git common directory is the primary checkout's `.git`, so a
// repository resolves to one root from any of its worktrees.
func Canonical(root string) string {
	common, err := git.CommonDir(root)
	if err != nil || filepath.Base(common) != ".git" {
		return root
	}
	return filepath.Dir(common)
}

// Key returns the pool key `<base>-<crc32>` for the canonical root of root. The
// canonicalization comes first, so a call from inside a linked worktree keys the
// primary repository.
func Key(root string) string {
	canonical := Canonical(root)
	sum := cksum([]byte(canonical + "\n"))
	return filepath.Base(canonical) + "-" + strconv.FormatUint(uint64(sum), 10)
}

// Pools returns the worktree pool parent below an explicitly resolved Bench home.
// A caller that tests a raw command's text for a pool path, before it resolves any
// root, tests against this prefix rather than restating "worktrees" itself.
func Pools(home string) string {
	return filepath.Join(home, "worktrees")
}

// Pool returns the repository's worktree pool directory below an explicitly resolved
// Bench home. The pool and the exec census are siblings under the home, and both name
// the repository by the same key.
func Pool(home, root string) string {
	return filepath.Join(Pools(home), Key(root))
}

// idPattern matches one 16-byte random identifier in hexadecimal, the form both
// halves of an assignment segment take.
var idPattern = regexp.MustCompile(`^[0-9a-f]{32}$`)

// AssignmentSegment returns the pool directory name of one assignment. The owner id
// and the assignment id join with a hyphen, and neither half holds one.
func AssignmentSegment(ownerID, assignmentID string) string {
	return ownerID + "-" + assignmentID
}

// SplitAssignmentSegment returns the assignment id of a pool directory name. A name
// that is not an owner id and an assignment id is not an assignment, so a reader of
// the pool can reject it without the ledger.
func SplitAssignmentSegment(segment string) (string, bool) {
	owner, id, found := strings.Cut(segment, "-")
	if !found || !idPattern.MatchString(owner) || !idPattern.MatchString(id) {
		return "", false
	}
	return id, true
}

// cksum is the POSIX cksum CRC-32. The key stays reproducible from a shell, where
// `printf '%s\n' <root> | cksum` gives the same number.
func cksum(data []byte) uint32 {
	var crc uint32
	step := func(value byte) {
		crc ^= uint32(value) << 24
		for range 8 {
			crc = crc<<1 ^ 0x04C11DB7*(crc>>31)
		}
	}
	for _, value := range data {
		step(value)
	}
	for n := len(data); n > 0; n >>= 8 {
		step(byte(n))
	}
	return ^crc
}
