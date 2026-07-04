// Package worktree owns the warm-pool worktree lifecycle: the pool-directory and
// lease-path addressing (Pool, LeaseFile), the atomic-lease state machine (Claim,
// reclaim), pool acquire/release (Acquire, Release), and the interactive subshell.
//
// The pool key is `<basename>-<cksum>`, where cksum is the POSIX `cksum` of the
// bytes `root + "\n"`. Go's standard library has no `cksum` variant, so the exact
// algorithm lives here once — a wrong figure silently addresses the wrong pool
// directory and breaks warm-pool reuse with no error, so this is the single source
// of that checksum and is pinned in the test against the system tool.
package worktree

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// cksumTable is the POSIX cksum CRC-32 table (polynomial 0x04C11DB7, MSB-first),
// built once at init so the checksum is a pure table lookup per byte.
var cksumTable = buildCksumTable()

func buildCksumTable() [256]uint32 {
	var t [256]uint32
	for i := 0; i < 256; i++ {
		c := uint32(i) << 24
		for k := 0; k < 8; k++ {
			if c&0x80000000 != 0 {
				c = (c << 1) ^ 0x04C11DB7
			} else {
				c <<= 1
			}
		}
		t[i] = c
	}
	return t
}

// cksum computes the POSIX `cksum` checksum of data: the CRC over the bytes, then
// over the length fed least-significant byte first, then complemented. The result
// matches the first field printed by the coreutils `cksum` tool.
func cksum(data []byte) uint32 {
	var crc uint32
	for _, x := range data {
		crc = (crc << 8) ^ cksumTable[byte(crc>>24)^x]
	}
	for n := len(data); n > 0; n >>= 8 {
		crc = (crc << 8) ^ cksumTable[byte(crc>>24)^byte(n&0xff)]
	}
	return ^crc
}

// benchHome resolves BENCH_HOME the way `bench.sh` does: the env var, or
// `<user home>/.bench` as the default.
func benchHome() string {
	if h := os.Getenv("BENCH_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bench")
}

// Pool returns the warm-pool directory for a repo root:
// `<BENCH_HOME>/worktrees/<basename>-<cksum>`. The cksum is taken over `root + "\n"`
// because the shell computed it via `echo "$root" | cksum`, and echo appends a newline.
func Pool(root string) string {
	sum := cksum([]byte(root + "\n"))
	key := filepath.Base(root) + "-" + strconv.FormatUint(uint64(sum), 10)
	return filepath.Join(benchHome(), "worktrees", key)
}

// LeaseFile returns the git-resolved lease path for a worktree at path:
// `git -C path rev-parse --git-path bench-lease`.
func LeaseFile(path string) (string, error) {
	return git.Output("-C", path, "rev-parse", "--git-path", "bench-lease")
}

// PoolCommand implements `bench worktree-pool <root>`. Root defaults to the current
// repo's top level when no argument is given; it prints the pool directory.
func PoolCommand(args []string) (string, int) {
	var root string
	if len(args) > 0 {
		root = args[0]
	} else {
		r, err := git.Root()
		if err != nil {
			return toon.NotInRepo() + "\n", 1
		}
		root = r
	}
	return Pool(root) + "\n", 0
}

// LeaseFileCommand implements `bench worktree-lease-file <path>`. It prints the
// lease path git resolves for the given worktree; a missing argument is a usage error.
func LeaseFileCommand(args []string) (string, int) {
	if len(args) == 0 {
		return "usage: bench worktree-lease-file <path>\n", 2
	}
	lease, err := LeaseFile(args[0])
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return lease + "\n", 0
}
