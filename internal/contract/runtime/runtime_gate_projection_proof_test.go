package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

type projectionFixture struct {
	f      contract.Fixture
	gitdir string
	cache  string
	state  string
	lock   *os.File
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
	p := projectionFixture{f: f, gitdir: gitDir(t, f), state: state}
	p.cache = filepath.Join(p.gitdir, "bench-last-gate")
	tree := strings.TrimSpace(f.Bench("tree-hash").Stdout)
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
	return p
}

func (p *projectionFixture) holdLock(t *testing.T) {
	t.Helper()
	var err error
	p.lock, err = os.OpenFile(filepath.Join(p.gitdir, "bench-gate.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(p.lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { syscall.Flock(int(p.lock.Fd()), syscall.LOCK_UN); p.lock.Close() })
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
	want := map[string]map[string]string{
		"status": {
			"absent": "bench: clean — nothing pending", "reusable-green": "bench: clean — nothing pending", "red": "gate       red", "stale": "gate       stale", "locked-pending": "gate       locked-pending", "interrupted-pending": "gate       interrupted-pending", "invalid": "gate       invalid verdict", "unavailable": "gate       verdict unavailable",
		},
		"dashboard": {
			"absent": "No gate cache yet", "reusable-green": `class="badge green">green`, "red": `class="badge red">red`, "stale": `class="badge stale">stale`, "locked-pending": ">locked-pending</span>", "interrupted-pending": ">interrupted-pending</span>", "invalid": `class="badge ">invalid`, "unavailable": `class="badge ">unavailable`,
		},
		"roadmap": {
			"absent": `false,"","","","","","",false`, "reusable-green": `true,ready,"",green`, "red": `true,ready,"",red`, "stale": `true,ready,"",green`, "locked-pending": "true,pending,locked-pending", "interrupted-pending": "true,pending,interrupted-pending", "invalid": `true,invalid,""`, "unavailable": `true,unavailable,""`,
		},
	}[surface][state]
	if !strings.Contains(out, want) {
		t.Fatalf("%s/%s missing literal %q:\n%s", state, surface, want, out)
	}
	if surface == "roadmap" && !strings.Contains(out, "gate_cache[1]{present,state,pending_status,status,cached_tree,work_tree,timestamp,stale}:") {
		t.Fatal("roadmap projection lost the typed gate-cache schema")
	}
}

func proveProjectionPurity(t *testing.T, state, surface string) {
	p := newProjectionFixture(t, state)
	if p.lock == nil {
		p.holdLock(t)
	}
	beforeData := contract.ReadFileAbs(t, p.cache)
	before, err := os.Lstat(p.cache)
	if err != nil {
		t.Fatal(err)
	}
	runProjectionSurface(t, p, surface)
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
}
