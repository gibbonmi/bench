package gocache

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// cleanTable is the one block `bench cache clean` prints, and cleanFields is its schema.
const cleanTable = "go_build_cache_clean"

var cleanFields = []string{"dir", "bytes_removed", "files_removed"}

// goTool is the fixed literal the clean resolves on PATH. It is a constant, never an
// argument, so the clean's child is the Go toolchain and nothing an operator can steer.
const goTool = "go"

// clean implements `bench cache clean` over an explicit environment slice, which is the
// seam the command tests drive. It takes the exclusive cache lock without waiting, so a
// live gate, lane, or focused run refuses it and keeps every archive it is reading. With
// the lock it measures the footprint, runs `go clean -cache` against the directory, and
// reports the difference the removal made.
func clean(env []string) (string, int) {
	dir, err := Dir(env)
	if err != nil {
		return toon.Errorf("cache directory not derived", err.Error()) + "\n", 1
	}
	if !toon.Representable(dir) {
		return toon.Errorf("unrepresentable cache directory", "the Bench build cache path holds a control byte; clear it from HOME") + "\n", 1
	}
	// An absent directory is answered before any lock, so a fresh machine passes without
	// the clean creating the very directory it reports as empty.
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return cleanRow(dir, 0, 0)
	}
	lock, refusal := lockForClean(dir)
	if refusal != "" {
		return refusal, 1
	}
	defer lock.Close()
	tool, err := lookPath(env, goTool)
	if err != nil {
		return toon.Errorf("go not found", "the Bench build cache clean runs `"+goTool+" clean -cache`, and no `"+goTool+"` is on PATH") + "\n", 1
	}
	childEnv, err := Apply(env)
	if err != nil {
		return toon.Errorf("cache directory not derived", err.Error()) + "\n", 1
	}
	before := Measure(dir)
	child := exec.Command(tool, "clean", "-cache")
	child.Env = childEnv
	child.Dir = dir
	if output, err := child.CombinedOutput(); err != nil {
		return toon.Errorf("go clean failed", firstLine(string(output))+" ("+err.Error()+")") + "\n", 1
	}
	after := Measure(dir)
	return cleanRow(dir, before.Bytes-after.Bytes, before.Files-after.Files)
}

// lockForClean opens the cache lock and takes the write lock without a wait. It answers a
// held descriptor or the operator's refusal line, which names the blocking pid.
func lockForClean(dir string) (*os.File, string) {
	path := filepath.Join(dir, LockFile)
	if localHolder(path) {
		return nil, heldRefusal(os.Getpid())
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, toon.Errorf("cache lock unavailable", err.Error()) + "\n"
	}
	lock := RecordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &lock); err != nil {
		defer file.Close()
		if !contended(err) {
			return nil, toon.Errorf("cache lock unavailable", err.Error()) + "\n"
		}
		return nil, heldRefusal(blockingPID(file))
	}
	return file, ""
}

// heldRefusal is the one refusal a live holder produces. It names the pid so the operator
// can find the run that is compiling, rather than being told only that something is.
func heldRefusal(pid int) string {
	return toon.Errorf("cache in use", "a Bench run holds the build cache lock: pid "+strconv.Itoa(pid)+"; retry after it finishes") + "\n"
}

// blockingPID asks the kernel which process holds the conflicting lock. A racing release
// between the failed set and this query leaves no holder to name, and 0 is that answer.
func blockingPID(file *os.File) int {
	query := RecordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(file.Fd(), syscall.F_GETLK, &query); err != nil {
		return 0
	}
	if query.Type == syscall.F_UNLCK {
		return 0
	}
	return int(query.Pid)
}

// lookPath resolves name against the PATH of the given environment slice, which is the
// operator-owned PATH of the interactive process. The Go standard resolver reads the
// process environment instead, and the clean's env is the seam its tests drive.
func lookPath(env []string, name string) (string, error) {
	for _, dir := range filepath.SplitList(value(env, "PATH")) {
		if dir == "" {
			continue
		}
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0 {
			return candidate, nil
		}
	}
	return "", exec.ErrNotFound
}

// cleanRow renders the one clean table at exit 0.
func cleanRow(dir string, bytes, files int64) (string, int) {
	block, err := toon.TableTyped(cleanTable, cleanFields, [][]any{{dir, bytes, files}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return block, 0
}

// firstLine reduces a child's output to its first line with every control byte gone, so a
// toolchain message cannot put an escape sequence on the operator's terminal.
func firstLine(output string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(output), "\n")
	return sanitize.Strip(line)
}
