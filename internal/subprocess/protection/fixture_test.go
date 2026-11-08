package protection_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

const (
	final     = 300 * time.Millisecond
	readLimit = 1 << 12
	stageEnv  = "PROTECTION_HELPER_STAGE"
	pidEnv    = "PROTECTION_HELPER_PID"
)

var (
	errInjected = errors.New("injected persistence fault")
	complete    = protection.Terminal{Process: protection.Complete, Streams: protection.Complete, Cleanup: protection.Complete}
)

func TestMain(m *testing.M) {
	processtest.Main(m, map[string]func() int{"nested": nestedHelper, "order": orderHelper})
}

// newProtocol admits the test kind: a resource directory, its fixed store, and its own path.
func newProtocol() *protection.Protocol {
	return &protection.Protocol{Validators: protection.Validators{protection.KindCache: validate}, ReadLimit: readLimit}
}

func validate(d protection.Descriptor) error {
	if d.Store != filepath.Join(d.ResourcePath, "store") || d.Binding != binding(d.ResourcePath) {
		return fmt.Errorf("descriptor %+v is foreign to its resource", d)
	}
	if info, err := os.Stat(d.ResourcePath); err != nil || !info.IsDir() {
		return fmt.Errorf("resource %s is absent", d.ResourcePath)
	}
	return nil
}

func candidate(path string) protection.Descriptor {
	return protection.Descriptor{Store: filepath.Join(path, "store"), Kind: protection.KindCache, ResourcePath: path, Binding: binding(path)}
}

// binding names the resource at path in hex, so only the path fields carry the path's raw bytes.
func binding(path string) string { return "test:" + hex.EncodeToString([]byte(path)) }

// newResource opens one real resource named name in a private directory.
func newResource(t *testing.T, p *protection.Protocol, name string) protection.Descriptor {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	d, err := p.OpenResource(candidate(path))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// ordered opens two resources and returns them in canonical store order.
func ordered(t *testing.T, p *protection.Protocol) (first, second protection.Descriptor) {
	t.Helper()
	pair := []protection.Descriptor{newResource(t, p, "cache"), newResource(t, p, "cache")}
	slices.SortFunc(pair, func(a, b protection.Descriptor) int { return strings.Compare(a.Store, b.Store) })
	return pair[0], pair[1]
}

// router fails the named replacements. A destination is a store's name and a record's base name.
func router(fail map[durablefiletest.Fault]map[string]error, stores map[string]string) *durablefiletest.Router {
	return &durablefiletest.Router{Fail: fail, Destination: func(path string) string {
		for store, name := range stores {
			if strings.HasPrefix(path, store+string(filepath.Separator)) {
				return name + "/" + filepath.Base(path)
			}
		}
		return path
	}}
}

// userRecords returns each user record of d's store as JSON fields, keyed by its directory.
func userRecords(t *testing.T, d protection.Descriptor) map[string]map[string]any {
	t.Helper()
	users := filepath.Join(d.Store, protection.UsersDir)
	entries, err := os.ReadDir(users)
	if err != nil {
		t.Fatal(err)
	}
	records := map[string]map[string]any{}
	for _, entry := range entries {
		var fields map[string]any
		data, err := os.ReadFile(filepath.Join(users, entry.Name(), protection.UserFile))
		if err == nil {
			err = json.Unmarshal(data, &fields)
		}
		if err != nil {
			t.Fatal(err)
		}
		records[entry.Name()] = fields
	}
	return records
}

// launch runs sh with body, and with its use marker as $0, under a registration of union.
func launch(t *testing.T, p *protection.Protocol, union []protection.Descriptor, body string, args ...string) (subprocess.Outcome, string) {
	t.Helper()
	marker := filepath.Join(t.TempDir(), "used")
	return run(t, exec.Command("sh", append([]string{"-c", body, marker}, args...)...), p, union), marker
}

// run starts cmd under a lifetime that registers union and returns its outcome.
func run(t *testing.T, cmd *exec.Cmd, p *protection.Protocol, union []protection.Descriptor) subprocess.Outcome {
	t.Helper()
	policy := subprocess.Policy{Final: bounds.ProcessWindow(final), Poll: bounds.ProcessPoll}
	l := subprocess.NewLifetime(cmd, policy, subprocess.Streams{}, subprocess.Registration{Protocol: p, Resources: union})
	if err := l.Start(); err != nil {
		t.Fatal(err)
	}
	return within(t, "the registered lifetime", func() subprocess.Outcome { o, _ := l.Wait(); return o })
}

// within fails the test unless fn returns inside the deadline derived from the final window.
func within[T any](t *testing.T, wait string, fn func() T) T {
	t.Helper()
	done := make(chan T, 1)
	go func() { done <- fn() }()
	window := bounds.TestDeadline(final)
	select {
	case got := <-done:
		return got
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict(wait, window))
	}
	var zero T
	return zero
}

func used(marker string) bool {
	_, err := os.Stat(marker)
	return err == nil
}

// requireRefused fails unless a launch under union refuses before its child runs.
func requireRefused(t *testing.T, p *protection.Protocol, union ...protection.Descriptor) error {
	t.Helper()
	o, marker := launch(t, p, union, `touch "$0"`)
	if used(marker) || o.Process.Started || o.Process.StartErr == nil {
		t.Fatalf("launch ran under refused protection: process=%+v protection=%+v", o.Process, o.Protection)
	}
	return o.Process.StartErr
}

// observe seals d and grades its users, and it reports whether release ran.
func observe(t *testing.T, p *protection.Protocol, d protection.Descriptor) (protection.Disposition, bool) {
	t.Helper()
	released := false
	disposition, _ := p.SealAndObserve(d, func() error { released = true; return nil })
	return disposition, released
}

// observeRelease runs the release observation of d, and it reports whether release ran.
func observeRelease(t *testing.T, p *protection.Protocol, d protection.Descriptor) (protection.Disposition, bool, error) {
	t.Helper()
	released := false
	disposition, err := p.ObserveRelease(d, func() error { released = true; return nil })
	return disposition, released, err
}

// requireRetained fails unless sealing d keeps the resource and runs no release.
func requireRetained(t *testing.T, p *protection.Protocol, d protection.Descriptor) protection.Disposition {
	t.Helper()
	disposition, released := observe(t, p, d)
	if disposition.State == protection.Safe || released {
		t.Fatalf("observation released a retained resource: %+v released=%t", disposition, released)
	}
	return disposition
}

// selection is the directory kind of each nested resource that a sweep test retires.
var selection = protection.Directory{Kind: protection.KindBinary, Noun: "selection"}

// retire opens a real store of a new directory resource at dir through adapter and retires
// its generation, so the store keeps only its retired record.
func retire(t *testing.T, p *protection.Protocol, adapter protection.Directory, dir string) protection.Descriptor {
	t.Helper()
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	g, err := adapter.Open(p, dir)
	if err != nil {
		t.Fatal(err)
	}
	if disposition, err := p.SealAndObserve(g.Resource, nil); err != nil || disposition.State != protection.Safe {
		t.Fatalf("retirement of %s = %+v (%v), want safe", dir, disposition, err)
	}
	return g.Resource
}

// privateProtocol returns a protocol of the directory kinds under a private store base.
func privateProtocol(t *testing.T, kinds ...protection.Directory) *protection.Protocol {
	t.Helper()
	protection.SetStoreBase(t, t.TempDir())
	validators := protection.Validators{}
	for _, kind := range kinds {
		validators[kind.Kind] = kind.Validate
	}
	return &protection.Protocol{Validators: validators, ReadLimit: readLimit}
}

// goneStore retires the resource dir through adapter and deletes the resource, as the removal
// of the enclosing directory does.
func goneStore(t *testing.T, p *protection.Protocol, adapter protection.Directory, dir string) protection.Descriptor {
	t.Helper()
	d := retire(t, p, adapter, dir)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	return d
}

// goneRetired retires a real nested store under a private store base and deletes its
// resource. It returns the protocol, the retired resource, and the enclosing directory.
func goneRetired(t *testing.T) (*protection.Protocol, protection.Descriptor, string) {
	t.Helper()
	p := privateProtocol(t, selection)
	d := goneStore(t, p, selection, filepath.Join(t.TempDir(), "gone"))
	return p, d, filepath.Dir(d.ResourcePath)
}

// requireSweep fails unless a nested sweep of dir returns inside the test deadline with no
// error, and leaves the store of d present exactly when kept.
func requireSweep(t *testing.T, p *protection.Protocol, dir string, d protection.Descriptor, kept bool) {
	t.Helper()
	if err := within(t, "the nested sweep", func() error { return p.DiscardNested(dir) }); err != nil {
		t.Fatalf("nested sweep of %s = %v, want no error", dir, err)
	}
	if _, err := os.Lstat(d.Store); (err == nil) != kept {
		t.Fatalf("store of %s after the sweep: %v, want present=%t", d.ResourcePath, err, kept)
	}
}

// user is one registered launch of a real leader in its own group.
type user struct {
	reg *protection.Registration
	cmd *exec.Cmd
}

// startUser registers and starts one real group leader with args, by default a long sleep.
func startUser(t *testing.T, p *protection.Protocol, d protection.Descriptor, args ...string) user {
	t.Helper()
	reg, err := p.Prepare([]protection.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	if len(args) == 0 {
		args = []string{"sleep", "300"}
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	u := user{reg, cmd}
	t.Cleanup(u.stop)
	if err := reg.Started(cmd.Process.Pid, cmd.Process.Pid); err != nil {
		t.Fatal(err)
	}
	return u
}

// stop kills the user's group and reaps its leader, so the group is absent.
func (u user) stop() {
	_ = syscall.Kill(-u.cmd.Process.Pid, syscall.SIGKILL)
	_ = u.cmd.Wait()
}

// vectorEnv is the environment entry that hands vector to a helper in its given order.
func vectorEnv(vector ...protection.Descriptor) string {
	return protection.Environ(nil, vector)[0]
}

// nestedHelper is a nested CLI that registers the inherited union. Its stage names how far
// it gets before it exits without a terminal record.
func nestedHelper() int {
	inherited, err := protection.Inherited()
	if err != nil {
		return 2
	}
	reg, err := newProtocol().Prepare(inherited)
	if err != nil {
		return 2
	}
	stage := os.Getenv(stageEnv)
	if stage == "prepared" {
		return 0
	}
	// The child leaves a descendant in its own group, publishes its PID, and exits.
	child := exec.Command("sh", "-c", `sleep 300 & echo $! > "$0"`, os.Getenv(pidEnv))
	child.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if child.Start() != nil || reg.Started(child.Process.Pid, child.Process.Pid) != nil || child.Wait() != nil {
		return 2
	}
	if stage == "workers" {
		// The descendant could still feed output, so the local stream worker stays pending.
		workers := protection.Terminal{Process: protection.Complete, Streams: protection.Pending, Cleanup: protection.Incomplete, Obligations: []string{"stdout"}}
		if reg.Finish(workers) != nil {
			return 2
		}
	}
	return 0
}

// orderHelper registers and closes the inherited vector, in its given order, many times.
func orderHelper() int {
	vector, err := protection.Inherited()
	if err != nil {
		return 2
	}
	for range 20 {
		reg, err := newProtocol().Prepare(vector)
		if err != nil || reg.Finish(protection.Terminal{NoStart: true}) != nil {
			return 1
		}
	}
	return 0
}
