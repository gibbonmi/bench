package protection_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestRegistrationPrelaunch has the child read its own durable user record at first use.
func TestRegistrationPrelaunch(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	o, marker := launch(t, p, []protection.Descriptor{d}, `cat "$1"/*/user.json > "$0"`, filepath.Join(d.Store, protection.UsersDir))
	seen, err := os.ReadFile(marker)
	if err != nil || o.Protection.State != protection.Safe || len(o.Protection.Users) != 1 {
		t.Fatalf("registered launch: err=%v protection=%+v", err, o.Protection)
	}
	if !strings.Contains(string(seen), `"user_id":"`+o.Protection.Users[0]+`"`) {
		t.Fatalf("the child used the resource before its record existed: it saw %q", seen)
	}
}

func TestRegistrationFirstPublishFailure(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "A/user.json", Ordinal: 1}: {"rename": errInjected}}, map[string]string{d.Store: "A"})
	protection.SetReplacer(p, faults.Replacer())
	err := requireRefused(t, p, d)
	durablefiletest.RequireFailure(t, err, "rename", false, errInjected)
	if records := userRecords(t, d); len(records) != 0 {
		t.Fatalf("a failed first publication left users %v", records)
	}
}

func TestRegistrationPartialPublish(t *testing.T) {
	p := newProtocol()
	first, second := ordered(t, p)
	faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "B/user.json", Ordinal: 1}: {"rename": errInjected}}, map[string]string{first.Store: "A", second.Store: "B"})
	protection.SetReplacer(p, faults.Replacer())
	if err := requireRefused(t, p, first, second); !errors.Is(err, errInjected) {
		t.Fatalf("start error = %v, want the second publication's fault", err)
	}
	if records := userRecords(t, second); len(records) != 0 {
		t.Fatalf("the failed store holds users %v", records)
	}
	for _, record := range userRecords(t, first) {
		if record["identity"] != "no-start" {
			t.Fatalf("rollback left %v, want a no-start record", record)
		}
	}
}

func TestRegistrationRollbackFailure(t *testing.T) {
	p := newProtocol()
	first, second := ordered(t, p)
	faults := router(map[durablefiletest.Fault]map[string]error{
		{Destination: "B/user.json", Ordinal: 1}: {"rename": errInjected},
		{Destination: "A/user.json", Ordinal: 2}: {"rename": errInjected},
	}, map[string]string{first.Store: "A", second.Store: "B"})
	protection.SetReplacer(p, faults.Replacer())
	o, marker := launch(t, p, []protection.Descriptor{first, second}, `touch "$0"`)
	records := userRecords(t, first)
	if used(marker) || o.Protection.State != protection.Uncertain || len(records) != 1 {
		t.Fatalf("failed rollback: used=%t protection=%+v records=%v", used(marker), o.Protection, records)
	}
	for _, record := range records {
		if record["identity"] != "pending" {
			t.Fatalf("failed rollback replaced the pending record with %v", record)
		}
	}
}

func TestRegistrationIdentityPublishFailure(t *testing.T) {
	p := newProtocol()
	first, second := ordered(t, p)
	faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "B/user.json", Ordinal: 2}: {"rename": errInjected}}, map[string]string{first.Store: "A", second.Store: "B"})
	protection.SetReplacer(p, faults.Replacer())
	o, _ := launch(t, p, []protection.Descriptor{first, second}, `exec sleep 5`)
	if !o.Process.Started || o.Process.Signal != syscall.SIGKILL || o.Protection.State != protection.Uncertain {
		t.Fatalf("partial identity: process=%+v protection=%+v", o.Process, o.Protection)
	}
	for d, want := range map[protection.Descriptor]string{first: "known", second: "pending"} {
		records := userRecords(t, d)
		for _, record := range records {
			if record["identity"] != want {
				t.Fatalf("%s holds %v, want a retained %s record", d.Store, record, want)
			}
		}
		if len(records) != 1 {
			t.Fatalf("%s holds users %v, want the one retained user", d.Store, records)
		}
	}
}

func TestRegistrationSealedAdmission(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	if _, err := p.Prepare([]protection.Descriptor{d}); err != nil {
		t.Fatal(err)
	}
	requireRetained(t, p, d)
	if err := requireRefused(t, p, d); !errors.Is(err, protection.ErrSealed) {
		t.Fatalf("sealed admission refused with %v", err)
	}
}

func TestRegistrationSealPending(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	reg, err := p.Prepare([]protection.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	disposition := requireRetained(t, p, d)
	if want := reg.Disposition().Users; disposition.State != protection.Protected || !slices.Equal(disposition.Users, want) {
		t.Fatalf("sealed pending user: %+v, want protected by %v", disposition, want)
	}
}

func TestRegistrationFinishSealed(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	reg, err := p.Prepare([]protection.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	requireRetained(t, p, d)
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	if err := reg.Started(cmd.Process.Pid, cmd.Process.Pid); err != nil {
		t.Fatalf("sealing refused an admitted user's identity: %v", err)
	}
	if err := reg.Finish(complete); err != nil {
		t.Fatalf("sealing refused an admitted user's completion: %v", err)
	}
	if disposition, released := observe(t, p, d); disposition.State != protection.Safe || !released {
		t.Fatalf("completed sealed resource: %+v released=%t", disposition, released)
	}
}

func TestRegistrationLockOrder(t *testing.T) {
	p := newProtocol()
	first, second := ordered(t, p)
	helpers := []*exec.Cmd{
		processtest.Helper(t, "order", vectorEnv(first, second)),
		processtest.Helper(t, "order", vectorEnv(second, first)),
	}
	for _, helper := range helpers {
		if err := helper.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = helper.Process.Kill() })
	}
	for _, helper := range helpers {
		if err := within(t, "the reversed registrations", helper.Wait); err != nil {
			t.Fatalf("reversed registration failed: %v", err)
		}
	}
}

func TestRegistrationNoStartComplete(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	reg, err := p.Prepare([]protection.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	recorder := &durablefiletest.Recorder{}
	protection.SetReplacer(p, recorder.Replacer())
	if err := reg.Finish(protection.Terminal{NoStart: true}); err != nil || reg.Disposition().State != protection.Safe {
		t.Fatalf("no-start close: err=%v disposition=%+v", err, reg.Disposition())
	}
	want := []string{"remove", "remove", "directory-open", "directory-sync", "directory-close"}
	if records := userRecords(t, d); len(records) != 0 || !slices.Equal(recorder.Trace, want) {
		t.Fatalf("no-start left users %v through trace %v, want a durable removal %v", records, recorder.Trace, want)
	}
	// A prepared launch whose start fails takes the same close path.
	o := run(t, exec.Command("/definitely/missing/bench-child"), p, []protection.Descriptor{d})
	if o.Process.StartErr == nil || o.Protection.State != protection.Safe || len(userRecords(t, d)) != 0 {
		t.Fatalf("failed start: process=%+v protection=%+v users=%v", o.Process, o.Protection, userRecords(t, d))
	}
}

func TestRegistrationResourceSyncFailure(t *testing.T) {
	p := newProtocol()
	path := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	fresh := candidate(path)
	faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "C/resource.json", Ordinal: 1}: {"directory-sync": errInjected}}, map[string]string{fresh.Store: "C"})
	protection.SetReplacer(p, faults.Replacer())
	d, err := p.OpenResource(fresh)
	durablefiletest.RequireFailure(t, err, "directory-sync", true, errInjected)
	requireRefused(t, p, d)
}

func TestRegistrationDeleteSyncFailure(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	u := startUser(t, p, d)
	u.stop()
	protection.SetReplacer(p, (&durablefiletest.Recorder{Fail: map[string]error{"directory-sync": errInjected}}).Replacer())
	if err := u.reg.Finish(complete); !errors.Is(err, errInjected) || u.reg.Disposition().State != protection.Uncertain {
		t.Fatalf("deletion-sync failure: err=%v disposition=%+v", err, u.reg.Disposition())
	}
}

func TestRegistrationDescriptorDedup(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	reg := within(t, "a repeated-descriptor registration", func() *protection.Registration {
		reg, err := p.Prepare([]protection.Descriptor{d, d})
		if err != nil {
			t.Error(err)
		}
		return reg
	})
	if records := userRecords(t, d); len(reg.Resources()) != 1 || len(records) != 1 {
		t.Fatalf("repeated descriptors registered %v with users %v", reg.Resources(), records)
	}
}

func TestRegistrationResourceJoin(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	reg, err := p.Prepare([]protection.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	joined, err := p.OpenResource(candidate(d.ResourcePath))
	if err != nil || joined != d {
		t.Fatalf("second owner opened %+v (%v), want the live generation %+v", joined, err, d)
	}
	if _, ok := userRecords(t, d)[reg.Disposition().Users[0]]; !ok {
		t.Fatal("the second owner's join removed the first live user")
	}
}
