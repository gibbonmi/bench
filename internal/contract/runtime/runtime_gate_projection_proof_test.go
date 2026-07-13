package runtime

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

type projectionFixture struct {
	f                    contract.Fixture
	gitdir               string
	cache                string
	state                string
	cachedTree, workTree string
	lockProcess          *exec.Cmd
	lockInput            io.WriteCloser
}

func normalizeDashboardCSS(css string) string {
	var b strings.Builder
	for i := 0; i < len(css); {
		if strings.HasPrefix(css[i:], "/*") {
			end := strings.Index(css[i+2:], "*/")
			if end < 0 {
				return strings.ToLower(b.String())
			}
			i += end + 4
			continue
		}
		if css[i] != '\\' {
			b.WriteByte(css[i])
			i++
			continue
		}
		i++
		start := i
		for i < len(css) && i-start < 6 && strings.ContainsRune("0123456789abcdefABCDEF", rune(css[i])) {
			i++
		}
		if i > start {
			value, _ := strconv.ParseInt(css[start:i], 16, 32)
			b.WriteRune(rune(value))
			if i < len(css) && strings.ContainsRune(" \t\r\n\f", rune(css[i])) {
				i++
			}
		} else if i < len(css) {
			b.WriteByte(css[i])
			i++
		}
	}
	return strings.ToLower(b.String())
}

func newProjectionFixture(t *testing.T, state string) projectionFixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	gateExit := "0"
	if state == "red" {
		gateExit = "1"
	}
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nprintf ran > .git/gate-reader-ran\nexit "+gateExit+"\n")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.WriteFile("tracked", "base\n")
	f.CommitAll("base")
	tree := strings.TrimSpace(f.Bench("tree-hash").Stdout)
	p := projectionFixture{f: f, gitdir: gitDir(t, f), state: state, cachedTree: tree}
	p.cache = filepath.Join(p.gitdir, "bench-last-gate")
	now := time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)
	oracle := strings.Repeat("0", 64)
	write := func(data string, mode os.FileMode) {
		contract.WriteFileAbs(t, p.cache, data)
		if err := os.Chmod(p.cache, mode); err != nil {
			t.Fatal(err)
		}
	}
	switch state {
	case "absent":
	case "reusable-green":
		f.Bench("gate").RequireExit(0)
		os.Remove(filepath.Join(p.gitdir, "gate-reader-ran"))
	case "red":
		f.Bench("gate").RequireExit(1)
		os.Remove(filepath.Join(p.gitdir, "gate-reader-ran"))
	case "stale":
		f.Bench("gate").RequireExit(0)
		os.Remove(filepath.Join(p.gitdir, "gate-reader-ran"))
		f.WriteFile("tracked", "drift\n")
	case "locked-pending", "interrupted-pending":
		write(fmt.Sprintf(`{"schema":1,"state":"pending","tree":%q,"oracle":%q,"started_at":%q,"owner_pid":123}`+"\n", tree, oracle, now), 0o600)
		if state == "locked-pending" {
			p.holdLock(t)
		}
	case "invalid":
		write("{hostile invalid}\n", 0o600)
	case "legacy":
		write("green "+tree+" "+now+"\n", 0o644)
	case "unavailable":
		write("{hostile invalid}\n", 0o600)
		f.WriteFile("unreadable", "blocked\n")
		if err := os.Chmod(filepath.Join(f.Root, "unreadable"), 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(filepath.Join(f.Root, "unreadable"), 0o600) })
	default:
		t.Fatalf("unknown projection state %q", state)
	}
	if state != "unavailable" {
		p.workTree = strings.TrimSpace(f.Bench("tree-hash").Stdout)
	}
	return p
}

func (p *projectionFixture) holdLock(t *testing.T) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestFT78GateLockHelper$")
	cmd.Env = append(os.Environ(), "FT78_GATE_LOCK_HELPER="+filepath.Join(p.gitdir, "bench-gate.lock"))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	input, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	ready := make([]byte, len("ready\n"))
	if _, err := io.ReadFull(stdout, ready); err != nil || string(ready) != "ready\n" {
		t.Fatalf("lock helper readiness = %q/%v", ready, err)
	}
	p.lockProcess, p.lockInput = cmd, input
	t.Cleanup(func() { input.Close(); _ = cmd.Wait() })
}

func TestFT78GateLockHelper(t *testing.T) {
	path := os.Getenv("FT78_GATE_LOCK_HELPER")
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	acquireTestGateLock(t, f)
	defer releaseTestGateLock(f)
	_, _ = os.Stdout.WriteString("ready\n")
	_, _ = io.Copy(io.Discard, os.Stdin)
}

func acquireTestGateLock(t testing.TB, f *os.File) {
	t.Helper()
	lock := syscall.Flock_t{Type: syscall.F_WRLCK, Whence: int16(io.SeekStart)}
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock); err != nil {
		t.Fatal(err)
	}
}

func releaseTestGateLock(f *os.File) {
	lock := syscall.Flock_t{Type: syscall.F_UNLCK, Whence: int16(io.SeekStart)}
	_ = syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock)
}

func runProjectionSurface(t *testing.T, p projectionFixture, surface string) string {
	t.Helper()
	var out contract.Probe
	switch surface {
	case "status":
		out = p.f.Bench("status", "--all")
	case "dashboard":
		out = p.f.Bench("dashboard", "--stdout")
	case "roadmap":
		out = p.f.Bench("roadmap", "--context")
	default:
		t.Fatalf("unknown projection surface %q", surface)
	}
	out.RequireExit(0)
	if out.Stderr != "" {
		t.Fatalf("%s stderr = %q", surface, out.Stderr)
	}
	return out.Stdout
}

func proveTypedProjection(t *testing.T, state, surface string) {
	p := newProjectionFixture(t, state)
	out := runProjectionSurface(t, p, surface)
	switch surface {
	case "status":
		want := exactStatusProjection(p)
		if out != want {
			t.Fatalf("%s/status projection\nwant: %q\ngot:  %q", state, want, out)
		}
	case "dashboard":
		requireExactDashboardGateProjection(t, out, p)
	case "roadmap":
		requireExactRoadmapGateProjection(t, out, p)
	}
}

func exactStatusProjection(p projectionFixture) string {
	switch p.state {
	case "absent", "reusable-green":
		return "bench: clean — nothing pending\n"
	case "red":
		return "▶ fix before commit  (gate)\n  gate       red                            → fix before commit\n"
	case "stale":
		return fmt.Sprintf("▶ commit on green  (git)\n  git        1 dirty path                   → commit on green\n  gate       stale (gated tree %.7s, work tree %.7s) → re-run the gate\n", p.cachedTree, p.workTree)
	case "locked-pending":
		return "▶ wait for live gate owner  (gate)\n  gate       locked-pending                 → wait for live gate owner\n"
	case "interrupted-pending":
		return "▶ re-run the gate  (gate)\n  gate       interrupted-pending            → re-run the gate\n"
	case "invalid":
		return "▶ re-run the gate  (gate)\n  gate       invalid verdict                → re-run the gate\n"
	default:
		return "▶ commit on green  (git)\n  git        1 dirty path                   → commit on green\n  gate       verdict unavailable            → inspect gate state\n"
	}
}

func expectedGateFields(p projectionFixture) []string {
	switch p.state {
	case "absent":
		return []string{"false", "", "", "", "", "", "", "false"}
	case "reusable-green":
		return []string{"true", "ready", "", "green", p.cachedTree, p.workTree, "<timestamp>", "false"}
	case "red":
		return []string{"true", "ready", "", "red", p.cachedTree, p.workTree, "<timestamp>", "false"}
	case "stale":
		return []string{"true", "ready", "", "green", p.cachedTree, p.workTree, "<timestamp>", "true"}
	case "locked-pending":
		return []string{"true", "pending", "locked-pending", "", p.cachedTree, p.workTree, "", "false"}
	case "interrupted-pending":
		return []string{"true", "pending", "interrupted-pending", "", p.cachedTree, p.workTree, "", "false"}
	case "invalid":
		return []string{"true", "invalid", "", "", "", p.workTree, "", "false"}
	default:
		return []string{"true", "unavailable", "", "", "", "", "", "false"}
	}
}

func requireExactRoadmapGateProjection(t *testing.T, out string, p projectionFixture) {
	t.Helper()
	const schema = "gate_cache[1]{present,state,pending_status,status,cached_tree,work_tree,timestamp,stale}:"
	lines := strings.Split(out, "\n")
	var rows [][]string
	for i, line := range lines {
		if strings.TrimSpace(line) != schema || i+1 >= len(lines) {
			continue
		}
		r := csv.NewReader(strings.NewReader(strings.TrimSpace(lines[i+1])))
		row, err := r.Read()
		if err != nil {
			t.Fatalf("parse gate-cache row: %v", err)
		}
		if _, err := r.Read(); err != io.EOF {
			t.Fatalf("gate-cache row has trailing CSV: %v", err)
		}
		rows = append(rows, row)
	}
	if len(rows) != 1 || len(rows[0]) != 8 {
		t.Fatalf("gate-cache projection rows = %#v, want one complete eight-field row", rows)
	}
	want, got := expectedGateFields(p), rows[0]
	if want[6] == "<timestamp>" {
		if _, err := time.Parse(time.RFC3339, got[6]); err != nil {
			t.Fatalf("gate-cache timestamp = %q: %v", got[6], err)
		}
		want[6] = got[6]
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s roadmap gate-cache row\nwant: %#v\ngot:  %#v", p.state, want, got)
	}
}

func proveProjectionPurity(t *testing.T, state, surface string) {
	p := newProjectionFixture(t, state)
	lockPath := filepath.Join(p.gitdir, "bench-gate.lock")
	if state == "interrupted-pending" {
		decoy, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		readLock := syscall.Flock_t{Type: syscall.F_RDLCK, Whence: int16(io.SeekStart)}
		if err := syscall.FcntlFlock(decoy.Fd(), syscall.F_SETLK, &readLock); err != nil {
			decoy.Close()
			t.Fatal(err)
		}
		t.Cleanup(func() { releaseTestGateLock(decoy); decoy.Close() })
	}
	lockBefore, lockBeforeErr := os.Lstat(lockPath)
	beforeData := contract.ReadFileAbs(t, p.cache)
	before, err := os.Lstat(p.cache)
	if err != nil {
		t.Fatal(err)
	}
	out := runProjectionSurface(t, p, surface)
	if state == "interrupted-pending" {
		switch surface {
		case "status":
			if out != exactStatusProjection(p) {
				t.Fatalf("%s changed the interrupted-pending projection: %q", surface, out)
			}
		case "dashboard":
			requireExactDashboardGateProjection(t, out, p)
		case "roadmap":
			requireExactRoadmapGateProjection(t, out, p)
		}
	}
	after, err := os.Lstat(p.cache)
	if err != nil {
		t.Fatal(err)
	}
	if got := contract.ReadFileAbs(t, p.cache); got != beforeData {
		t.Fatalf("%s rewrote cache bytes", surface)
	}
	if before.Mode() != after.Mode() || !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf("%s repaired mode or changed cache mtime", surface)
	}
	if _, err := os.Stat(filepath.Join(p.gitdir, "gate-reader-ran")); !os.IsNotExist(err) {
		t.Fatalf("%s executed the oracle", surface)
	}
	lockAfter, lockAfterErr := os.Lstat(lockPath)
	if (os.IsNotExist(lockBeforeErr) && !os.IsNotExist(lockAfterErr)) || (lockBeforeErr == nil && (lockAfterErr != nil || lockBefore.Mode() != lockAfter.Mode() || !lockBefore.ModTime().Equal(lockAfter.ModTime()))) {
		t.Fatalf("%s created or mutated the execution lock", surface)
	}
}
