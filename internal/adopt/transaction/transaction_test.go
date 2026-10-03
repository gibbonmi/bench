package transaction_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt/transaction"
)

func TestCompatibilityBackupFailure(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "managed")
	stage := filepath.Join(root, "replacement")
	blocked := filepath.Join(root, "not-a-directory")
	for path, body := range map[string]string{dest: "original", stage: "replacement", blocked: "record blocker"} {
		if err := os.WriteFile(path, []byte(body), 0o640); err != nil {
			t.Fatal(err)
		}
	}
	_, err := (transaction.Store{Directory: filepath.Join(blocked, "records")}).Apply([]transaction.Change{{Destination: dest, Stage: stage}})
	if err == nil {
		t.Fatal("repair succeeded without publishing a backup")
	}
	data, readErr := os.ReadFile(dest)
	info, statErr := os.Stat(dest)
	if readErr != nil || statErr != nil || string(data) != "original" || info.Mode().Perm() != 0o640 {
		t.Fatalf("backup failure changed destination: %q, %v, %v", data, readErr, statErr)
	}
}

func TestRetainedPreimageRoundTrip(t *testing.T) {
	fixture := newRepairFixture(t)
	dest, stage, store := fixture.dest, fixture.stage, fixture.store
	id, err := store.Apply([]transaction.Change{{Destination: dest, Stage: stage}})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("repair did not retain a recovery identifier")
	}
	data, err := os.ReadFile(dest)
	if err != nil || string(data) != "after" {
		t.Fatalf("replacement = %q, %v", data, err)
	}
	if err := (transaction.Store{Directory: store.Directory}).Undo(id); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(dest)
	info, statErr := os.Stat(dest)
	if err != nil || statErr != nil || string(data) != "before" || info.Mode().Perm() != 0o640 {
		t.Fatalf("restored preimage = %q, %v, %v", data, err, statErr)
	}
}

func TestCompatibilityPlanDrift(t *testing.T) {
	root := t.TempDir()
	dest, stage := filepath.Join(root, "managed"), filepath.Join(root, "replacement")
	for path, body := range map[string]string{dest: "managed before", stage: "canonical after"} {
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	change, err := transaction.Prepare(dest, stage)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("later user change"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = (transaction.Store{Directory: filepath.Join(root, "records")}).Apply([]transaction.Change{change})
	if err == nil {
		t.Fatal("repair accepted a destination changed after inspection")
	}
	data, readErr := os.ReadFile(dest)
	if readErr != nil || string(data) != "later user change" {
		t.Fatalf("stale repair changed later content: %q, %v", data, readErr)
	}
}

func TestCompatibilityUndoConflict(t *testing.T) {
	fixture := newRepairFixture(t)
	dest, stage, store := fixture.dest, fixture.stage, fixture.store
	id, err := store.Apply([]transaction.Change{{Destination: dest, Stage: stage}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("before"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dest, 0o640); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Undo(id); err == nil {
		t.Fatal("undo accepted a completed repair whose postimage changed")
	}
	current, err := os.Stat(dest)
	if err != nil || !os.SameFile(info, current) {
		t.Fatalf("undo replaced the later edit: %v", err)
	}
}

func TestSharedDestinationExclusion(t *testing.T) {
	root := t.TempDir()
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	first, err := transaction.Lock([]string{filepath.Join(root, "absent", "managed")})
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	if err := os.Mkdir(filepath.Join(root, "absent"), 0o700); err != nil {
		t.Fatal(err)
	}
	second, err := transaction.Lock([]string{filepath.Join(alias, "absent", "managed")})
	if err == nil {
		second.Close()
		t.Fatal("alias writer acquired a destination already owned before its parent existed")
	}
	unrelated, err := transaction.Lock([]string{filepath.Join(root, "absent", "unrelated")})
	if err != nil {
		t.Fatalf("disjoint destination blocked: %v", err)
	}
	if err := unrelated.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRepairHonorsDestinationExclusion(t *testing.T) {
	root := t.TempDir()
	dest, stage := filepath.Join(root, "managed"), filepath.Join(root, "replacement")
	if err := os.WriteFile(stage, []byte("replacement"), 0o600); err != nil {
		t.Fatal(err)
	}
	lease, err := transaction.Lock([]string{dest})
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	_, err = (transaction.Store{Directory: filepath.Join(root, "records")}).Apply([]transaction.Change{{Destination: dest, Stage: stage}})
	if err == nil {
		t.Fatal("repair published a destination owned by another writer")
	}
	if _, err := os.Lstat(dest); !os.IsNotExist(err) {
		t.Fatalf("competing repair created the destination: %v", err)
	}
}

func TestCompatibilityTerminalFailure(t *testing.T) {
	fixture := newRepairFixture(t)
	dest, stage, store := fixture.dest, fixture.stage, fixture.store
	store.SyncDirectory = func(dir string) error {
		if filepath.Dir(dir) == store.Directory {
			data, _ := os.ReadFile(filepath.Join(dir, "record.json"))
			var state struct{ State string }
			_ = json.Unmarshal(data, &state)
			if state.State == "applied" {
				return os.ErrPermission
			}
		}
		f, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer f.Close()
		return f.Sync()
	}
	id, err := store.Apply([]transaction.Change{{Destination: dest, Stage: stage}})
	if err == nil {
		t.Fatal("repair reported success after terminal persistence failed")
	}
	if id == "" {
		t.Fatal("incomplete repair lost its recovery identifier")
	}
	data, readErr := os.ReadFile(dest)
	if readErr != nil || string(data) != "after" {
		t.Fatalf("terminal failure fixture did not reach publication: %q, %v", data, readErr)
	}
	if err := (transaction.Store{Directory: store.Directory}).Undo(id); err != nil {
		t.Fatalf("recovery state unavailable: %v", err)
	}
	data, readErr = os.ReadFile(dest)
	if readErr != nil || string(data) != "before" {
		t.Fatalf("terminal failure lost the preimage: %q, %v", data, readErr)
	}
}

func TestDestinationIdentityAtPublication(t *testing.T) {
	fixture := newRepairFixture(t)
	dest, stage, store := fixture.dest, fixture.stage, fixture.store
	root := fixture.root
	var replacement os.FileInfo
	store.SyncDirectory = func(dir string) error {
		if dir == store.Directory && replacement == nil {
			competing := filepath.Join(root, "competing")
			if err := os.WriteFile(competing, []byte("before"), 0o640); err != nil {
				return err
			}
			if err := os.Rename(competing, dest); err != nil {
				return err
			}
			var err error
			replacement, err = os.Stat(dest)
			if err != nil {
				return err
			}
		}
		f, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer f.Close()
		return f.Sync()
	}
	_, err := store.Apply([]transaction.Change{{Destination: dest, Stage: stage}})
	if err == nil {
		t.Fatal("repair overwrote an identity changed after backup publication")
	}
	current, statErr := os.Stat(dest)
	if replacement == nil || statErr != nil || !os.SameFile(replacement, current) {
		t.Fatalf("repair did not preserve the competing replacement: %v", statErr)
	}
}

type repairFixture struct {
	root, dest, stage string
	store             transaction.Store
}

func newRepairFixture(t *testing.T) repairFixture {
	t.Helper()
	root := t.TempDir()
	f := repairFixture{root: root, dest: filepath.Join(root, "managed"), stage: filepath.Join(root, "replacement"), store: transaction.Store{Directory: filepath.Join(root, "records")}}
	if err := os.WriteFile(f.dest, []byte("before"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.stage, []byte("after"), 0o750); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestCompatibilityRestoreFailure(t *testing.T) {
	f := newRepairFixture(t)
	second := filepath.Join(f.root, "other", "managed")
	if err := os.MkdirAll(filepath.Dir(second), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("second before"), 0o600); err != nil {
		t.Fatal(err)
	}
	id, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}, {Destination: second, Stage: f.stage}})
	if err != nil {
		t.Fatal(err)
	}
	f.store.SyncDirectory = func(dir string) error {
		if dir == filepath.Dir(f.dest) || dir == filepath.Dir(second) {
			return os.ErrPermission
		}
		file, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer file.Close()
		return file.Sync()
	}
	err = f.store.Undo(id)
	if err == nil || !strings.Contains(err.Error(), f.dest) || !strings.Contains(err.Error(), second) {
		t.Fatalf("recovery did not report every failed target: %v", err)
	}
	if err := (transaction.Store{Directory: f.store.Directory}).Undo(id); err != nil {
		t.Fatalf("failed restore lost its recovery state: %v", err)
	}
}

func TestCompatibilityRepairPaths(t *testing.T) {
	f := newRepairFixture(t)
	id, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
	if err != nil {
		t.Fatal(err)
	}
	for _, hostile := range []string{"", "..", "../" + id, "/" + id, id + "/record.json", strings.Repeat("A", 32)} {
		t.Run(hostile, func(t *testing.T) {
			if err := f.store.Undo(hostile); err == nil {
				t.Fatalf("undo accepted hostile identifier %q", hostile)
			}
		})
	}
	linked := strings.Repeat("a", 32)
	if linked == id {
		linked = strings.Repeat("b", 32)
	}
	if err := os.Symlink(filepath.Join(f.store.Directory, id), filepath.Join(f.store.Directory, linked)); err != nil {
		t.Fatal(err)
	}
	if err := f.store.Undo(linked); err == nil {
		t.Fatal("undo followed a repair-record directory link")
	}
	data, err := os.ReadFile(f.dest)
	if err != nil || string(data) != "after" {
		t.Fatalf("invalid repair identifier changed the destination: %q, %v", data, err)
	}
}

func TestPreparedDestinationAlias(t *testing.T) {
	f := newRepairFixture(t)
	alias := filepath.Join(t.TempDir(), "linked-root")
	if err := os.Symlink(f.root, alias); err != nil {
		t.Fatal(err)
	}
	change, err := transaction.Prepare(filepath.Join(alias, filepath.Base(f.dest)), f.stage)
	if err != nil {
		t.Fatalf("canonical destination alias refused: %v", err)
	}
	id, err := f.store.Apply([]transaction.Change{change})
	if err != nil {
		t.Fatal(err)
	}
	if err := f.store.Undo(id); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(f.dest)
	if err != nil || string(data) != "before" {
		t.Fatalf("canonical destination not restored: %q, %v", data, err)
	}
}

func TestBackupAncestorPersistenceFailure(t *testing.T) {
	f := newRepairFixture(t)
	f.store.Directory = filepath.Join(f.root, "new-parent", "records")
	changedAtFailure := false
	f.store.SyncDirectory = func(dir string) error {
		if dir == f.root {
			data, err := os.ReadFile(f.dest)
			if err != nil {
				return err
			}
			changedAtFailure = changedAtFailure || string(data) != "before"
			return os.ErrPermission
		}
		file, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer file.Close()
		return file.Sync()
	}
	_, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
	if err == nil || changedAtFailure {
		t.Fatalf("repair mutated a destination before its backup directory was durable: error=%v, changed=%v", err, changedAtFailure)
	}
}

func TestCanonicalParentTraversal(t *testing.T) {
	f := newRepairFixture(t)
	child := filepath.Join(f.root, "child")
	if err := os.Mkdir(child, 0o700); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(child, alias); err != nil {
		t.Fatal(err)
	}
	change, err := transaction.Prepare(alias+"/../managed", f.stage)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.store.Apply([]transaction.Change{change}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(f.dest)
	if err != nil || string(data) != "after" {
		t.Fatalf("repair did not reach the physical parent: %q, %v", data, err)
	}
}
