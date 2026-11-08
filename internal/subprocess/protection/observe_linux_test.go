package protection_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// waitZombie waits until every task of the unreaped child pid is a zombie.
func waitZombie(t *testing.T, pid int) {
	t.Helper()
	processtest.WaitFor(t, fmt.Sprintf("process %d to become a zombie", pid), func() bool { return protection.ProcAbsent("/proc", pid) })
}

// zombieUser registers one known user whose group leader runs script and whose record names
// that leader as its actor too. It returns when the leader is an unreaped zombie, and the
// user's cleanup reaps it.
func zombieUser(t *testing.T, p *protection.Protocol, d protection.Descriptor, script string) {
	t.Helper()
	leader := startUser(t, p, d, "sh", "-c", script).cmd.Process.Pid
	for id, fields := range userRecords(t, d) {
		fields["actor_pid"] = leader
		data, err := json.Marshal(fields)
		if err != nil {
			t.Fatal(err)
		}
		plant(t, filepath.Join(d.Store, protection.UsersDir, id, protection.UserFile), string(data))
	}
	waitZombie(t, leader)
}

// TestReleaseObservationRetiresAZombieUser is SP11. Signal zero answers for an unreaped
// zombie, so only the task states show that the user's actor and group are dead.
func TestReleaseObservationRetiresAZombieUser(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	zombieUser(t, p, d, "exit 0")
	if disposition, released, err := observeRelease(t, p, d); err != nil || disposition.State != protection.Safe || !released {
		t.Fatalf("release observation of a zombie user = %+v (%v) released=%t, want safe with a release", disposition, err, released)
	}
	if _, err := os.Lstat(filepath.Join(d.Store, protection.ResourceFile)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("generation after the release: %v, want it absent", err)
	}
}

// TestZombiePendingActorIsReleased is SP12. A nested CLI registers a pending user and exits,
// and its parent does not reap it.
func TestZombiePendingActorIsReleased(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	actor := processtest.Helper(t, "nested", vectorEnv(d), stageEnv+"=prepared")
	if err := actor.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = actor.Wait() })
	waitZombie(t, actor.Process.Pid)
	if users := userRecords(t, d); len(users) != 1 {
		t.Fatalf("the zombie actor left users %v, want one pending user", users)
	}
	if disposition, released, err := observeRelease(t, p, d); err != nil || disposition.State != protection.Safe || !released {
		t.Fatalf("release observation of a zombie actor = %+v (%v) released=%t, want safe with a release", disposition, err, released)
	}
}

// TestZombieLeaderWithALiveMemberIsProtected is SP13. The leader is a zombie, and the sleep
// that it started lives on in its group.
func TestZombieLeaderWithALiveMemberIsProtected(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	zombieUser(t, p, d, "sleep 300 & exit 0")
	if disposition, released, err := observeRelease(t, p, d); err != nil || disposition.State != protection.Protected || released {
		t.Fatalf("release observation of a live group = %+v (%v) released=%t, want protected with no release", disposition, err, released)
	}
}

// procRoot writes a fixture proc tree. Each key is a path under the root and each value is
// the record at that path; a key that ends in "/" is a directory.
func procRoot(t *testing.T, records map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, record := range records {
		path := filepath.Join(root, filepath.FromSlash(name))
		if strings.HasSuffix(name, "/") {
			if err := os.MkdirAll(path, 0o755); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(record), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// stat is the stat record of task pid with the command name and the state, in group pgid.
func stat(pid int, name, state string, pgid int) string {
	return fmt.Sprintf("%d (%s) %s 1 %d %d 0 -1 4194304\n", pid, name, state, pgid, pgid)
}

// TestProcLiveTaskBesideAZombieLeader is SP14. A leader thread that exits first is a zombie
// while the other threads of its process run on.
func TestProcLiveTaskBesideAZombieLeader(t *testing.T) {
	root := procRoot(t, map[string]string{
		"7/stat":        stat(7, "a", "Z", 7),
		"7/task/7/stat": stat(7, "a", "Z", 7),
		"7/task/8/stat": stat(8, "a", "S", 7),
	})
	if protection.ProcAbsent(root, 7) {
		t.Fatal("a process with a sleeping task beside its zombie leader counts as absent, want present")
	}
}

// TestProcStatParsesAHostileCommandName is SP15. The command name forges a zombie state
// before the real one.
func TestProcStatParsesAHostileCommandName(t *testing.T) {
	root := procRoot(t, map[string]string{"7/task/7/stat": stat(7, "a) Z 1 (b", "S", 7)})
	if protection.ProcAbsent(root, 7) {
		t.Fatal("a sleeping task whose command name holds a forged zombie state counts as absent, want present")
	}
}

// TestProcUnreadableStatKeepsPresent is SP16. Each tree holds a zombie beside a record that
// the reader cannot read or parse, or holds no task or no member at all.
func TestProcUnreadableStatKeepsPresent(t *testing.T) {
	zombie := stat(7, "a", "Z", 7)
	member := stat(8, "a", "Z", 7)
	for _, tc := range []struct {
		name    string
		target  int
		records map[string]string
	}{
		{"unreadable task record", 7, map[string]string{"7/task/7/stat": zombie, "7/task/8/stat/": ""}},
		{"malformed task record", 7, map[string]string{"7/task/7/stat": zombie, "7/task/8/stat": "8 (a S 1 7"}},
		{"short task record", 7, map[string]string{"7/task/7/stat": "7 (a) Z"}},
		{"unreadable member record", -7, map[string]string{"7/stat": zombie, "7/task/7/stat": zombie, "8/stat/": ""}},
		{"unreadable member task directory", -7, map[string]string{"7/stat": zombie, "7/task": "", "8/stat": member, "8/task/8/stat": member}},
		{"no task", 7, map[string]string{"7/stat": zombie, "7/task/": ""}},
		{"member with no task", -7, map[string]string{"7/stat": zombie, "7/task/": ""}},
		{"no member", -7, map[string]string{"9/stat": stat(9, "b", "Z", 3), "9/task/9/stat": stat(9, "b", "Z", 3)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if protection.ProcAbsent(procRoot(t, tc.records), tc.target) {
				t.Fatalf("target %d counts as absent, want present", tc.target)
			}
		})
	}
}

// TestProcVanishedMemberIsSkipped is SP17. Member 8 exited between the directory read and
// the record read, process 9 sleeps in another group, and the root holds its own stat file.
func TestProcVanishedMemberIsSkipped(t *testing.T) {
	root := procRoot(t, map[string]string{
		"stat":          "cpu 0 0 0\n",
		"7/stat":        stat(7, "a", "Z", 7),
		"7/task/7/stat": stat(7, "a", "Z", 7),
		"8/":            "",
		"9/stat":        stat(9, "b", "S", 3),
		"9/task/9/stat": stat(9, "b", "S", 3),
	})
	if !protection.ProcAbsent(root, -7) {
		t.Fatal("a group of one zombie member and one gone member counts as present, want absent")
	}
}
