// Lifecycle owns the lease state machine and pool acquire/release, beside the pool-path
// and lease-path addressing they operate over. The lease is an atomically-created file
// recording "<pid> <utc-time>";
// a lease is reclaimable only when its owner is provably gone (a recorded pid no
// longer running, or unreadable/legacy content aged out by mtime — never a
// fresh-empty writer mid-claim). Release cleans and unleases only for the recorded
// owner, so a stale-reclaimed worktree's new owner is left alone.
package worktree

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/git"
)

// staleAfter is how long an unreadable/legacy (non-numeric-pid) lease must have gone
// untouched before it is treated as a crashed writer's leftover and reclaimed. A
// fresh-empty lease younger than this is a writer mid-claim and is respected — the
// threshold is the whole difference between reclaiming a zombie and stealing a live
// lease, so it is a named constant with one source. Mirrors the shell `find -mmin +1`.
const staleAfter = time.Minute

// pidAlive reports whether a process with the given pid exists, matching `kill -0`:
// signal 0 succeeds (nil) for a live process the caller may signal, and returns EPERM
// for a live process owned by another user — both mean alive. Only ESRCH (no such
// process) means gone.
func pidAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// reclaimable decides whether an existing lease may be taken over, given its raw
// content, its mtime, the current time, and a liveness probe. A recorded numeric pid
// gates on liveness: a dead pid reclaims, a live one is respected. Non-numeric or
// empty content (unreadable, legacy, or a fresh-empty writer) reclaims only once it
// has aged past staleAfter — so a writer mid-claim, whose lease is empty but fresh, is
// never stolen. This is the four-way decision the black-box lease contracts exercise
// but cannot cheaply enumerate.
func reclaimable(content []byte, mtime, now time.Time, alive func(int) bool) bool {
	field := strings.Fields(string(content))
	if len(field) > 0 {
		if pid, err := strconv.Atoi(field[0]); err == nil {
			return !alive(pid)
		}
	}
	return now.Sub(mtime) > staleAfter
}

// candidateName is the pooled-worktree path a mint attempt uses: a name unique per
// (second, pid, try) kept inside the pool directory, so a wrong name can never mint
// outside the pool and silently break warm reuse.
func candidateName(pool string, unixSecs int64, pid, try int) string {
	return filepath.Join(pool, fmt.Sprintf("%d-%d-%d", unixSecs, pid, try))
}

// leaseLine is the bytes an owner writes into its lease: "<pid> <utc-time>\n".
func leaseLine() []byte {
	return []byte(fmt.Sprintf("%d %s\n", os.Getpid(), time.Now().UTC().Format("2006-01-02T15:04:05Z")))
}

// tryCreate attempts the atomic O_EXCL lease create — the claim's race-winner. It
// returns true only when this process created the file; an existing lease (another
// claimant) fails without clobbering it.
func tryCreate(leasePath string) bool {
	f, err := os.OpenFile(leasePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return false
	}
	_, werr := f.Write(leaseLine())
	cerr := f.Close()
	return werr == nil && cerr == nil
}

// Claim atomically claims a worktree's lease. A first-writer create wins outright; an
// existing lease is taken over only when reclaimable reports its owner provably gone,
// and then via an atomic rename so two concurrent reclaimers cannot both win — only
// the process whose rename succeeds re-creates the lease. Returns whether this process
// now owns the lease.
func Claim(leasePath string) bool {
	if tryCreate(leasePath) {
		return true
	}
	info, err := os.Stat(leasePath)
	if err != nil {
		return false // lease vanished under us (a racing reclaim); respect and rescan
	}
	content, _ := os.ReadFile(leasePath)
	if !reclaimable(content, info.ModTime(), time.Now(), pidAlive) {
		return false
	}
	stale := leasePath + ".stale." + strconv.Itoa(os.Getpid())
	if os.Rename(leasePath, stale) != nil {
		return false // another reclaimer moved it first
	}
	os.Remove(stale)
	return tryCreate(leasePath)
}

// isWorktree reports whether dir is a git worktree checkout — it holds a `.git` file
// (linked worktree) or directory. Mirrors the shell's `[ -d .git || -f .git ]` scan gate.
func isWorktree(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// isClean reports whether the worktree at dir has no dirty or untracked paths.
func isClean(dir string) bool {
	out, err := git.Output("-C", dir, "status", "--porcelain")
	return err == nil && out == ""
}

// Acquire returns a leased, clean pool worktree for root: it scans the pool for a
// clean, released entry to claim, and mints fresh detached worktrees (bounded to three
// tries) when none is reusable. The returned worktree is detached, reset to resetRef
// (or HEAD when empty), and cleaned. A reset failure aborts unless resetMode is "soft",
// which falls back to a plain hard reset — the interactive subshell's tolerance for a
// resetRef that no longer resolves.
func Acquire(root, resetRef, resetMode string) (string, error) {
	pool := Pool(root)
	if err := os.MkdirAll(pool, 0o777); err != nil {
		return "", err
	}
	// Best-effort refresh so a freshly-minted worktree can detach onto origin/<branch>;
	// a repo with no origin (the contract fixtures) just skips it.
	_ = exec.Command("git", "-C", root, "fetch", "-q", "origin").Run()

	var wt string
	entries, _ := os.ReadDir(pool) // sorted by name, matching the shell glob order
	for _, e := range entries {
		d := filepath.Join(pool, e.Name())
		if !isWorktree(d) || !isClean(d) {
			continue
		}
		lease, err := LeaseFile(d)
		if err != nil || !Claim(lease) {
			continue
		}
		wt = d
		break
	}

	for try := 1; wt == "" && try <= 3; try++ {
		cand := candidateName(pool, time.Now().Unix(), os.Getpid(), try)
		if !worktreeAdd(root, cand, "origin/"+git.DefaultBranch(root)) && !worktreeAdd(root, cand, "") {
			break
		}
		lease, err := LeaseFile(cand)
		if err != nil {
			continue
		}
		if Claim(lease) {
			wt = cand
		}
	}
	if wt == "" {
		return "", errors.New("could not lease a pool worktree")
	}

	_ = exec.Command("git", "-C", wt, "switch", "-q", "--detach").Run()
	if resetRef != "" {
		if exec.Command("git", "-C", wt, "reset", "-q", "--hard", resetRef).Run() != nil {
			if resetMode != "soft" {
				return "", fmt.Errorf("could not reset pool worktree to %s", resetRef)
			}
			_ = exec.Command("git", "-C", wt, "reset", "-q", "--hard").Run()
		}
	} else {
		_ = exec.Command("git", "-C", wt, "reset", "-q", "--hard").Run()
	}
	if err := exec.Command("git", "-C", wt, "clean", "-qfdx").Run(); err != nil {
		return "", err
	}
	return strings.TrimRight(wt, "/"), nil
}

// worktreeAdd runs `git -C root worktree add -q --detach cand [ref]`, reporting success.
// An empty ref adds at the current HEAD (the fallback when origin/<branch> is absent).
func worktreeAdd(root, cand, ref string) bool {
	args := []string{"-C", root, "worktree", "add", "-q", "--detach", cand}
	if ref != "" {
		args = append(args, ref)
	}
	return exec.Command("git", args...).Run() == nil
}

// Release cleans and unleases a worktree, but only for its recorded owner: a lease
// held by a different, still-live process means the worktree was stale-reclaimed and
// now belongs to that owner, so a non-owner's deferred cleanup leaves it alone. The
// owner detaches, hard-resets, and cleans ignored+untracked files *before* removing
// the lease: the entry never sits claimable while dirty, and once the lease is gone
// it is no longer ours to touch — a concurrent Acquire may claim it immediately.
func Release(wt string) {
	if wt == "" {
		return
	}
	lease, err := LeaseFile(wt)
	if err != nil {
		return
	}
	content, _ := os.ReadFile(lease)
	if field := strings.Fields(string(content)); len(field) > 0 {
		if pid, err := strconv.Atoi(field[0]); err == nil && pid != os.Getpid() && pidAlive(pid) {
			return
		}
	}
	_ = exec.Command("git", "-C", wt, "switch", "-q", "--detach").Run()
	_ = exec.Command("git", "-C", wt, "reset", "-q", "--hard").Run()
	_ = exec.Command("git", "-C", wt, "clean", "-qfdx").Run()
	os.Remove(lease)
}
