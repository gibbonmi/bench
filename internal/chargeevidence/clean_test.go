package chargeevidence_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ce "github.com/gibbonmi/bench/internal/chargeevidence"
)

// planOf reads one store's complete cleanup plan.
func planOf(t *testing.T, store *ce.Store) ce.CleanupPlan {
	t.Helper()
	plan, err := store.Plan()
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	return plan
}

// targetIDs renders one plan's targets as "kind id" pairs in plan order, so a test compares
// the exact inventory rather than a count.
func targetIDs(plan ce.CleanupPlan) []string {
	ids := make([]string, len(plan.Targets))
	for i, target := range plan.Targets {
		ids[i] = target.Kind + " " + target.ID
	}
	return ids
}

// orphan writes one temporary object that no writer holds, which is the shape a killed
// writer leaves behind.
func orphan(t *testing.T, dir, body string) string {
	t.Helper()
	name := ce.TempPrefix + "orphan" + ce.TempSuffix
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return name
}

// TestEvidenceCleanupEmpty is CE100 and CE119. An absent store plans nothing and creates
// nothing, and an existing store that published nothing plans nothing.
func TestEvidenceCleanupEmpty(t *testing.T) {
	t.Run("CE119 absent store", func(t *testing.T) {
		store, dir := newStore(t, ce.StoreOptions{})
		plan := planOf(t, store)
		if len(plan.Targets) != 0 || plan.Bytes != 0 {
			t.Fatalf("absent store planned %v", targetIDs(plan))
		}
		if _, err := os.Lstat(dir); !os.IsNotExist(err) {
			t.Fatalf("planning created the store directory: %v", err)
		}
	})
	t.Run("CE100 empty store", func(t *testing.T) {
		store, dir := newStore(t, ce.StoreOptions{})
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		plan := planOf(t, store)
		if len(plan.Targets) != 0 || plan.Bytes != 0 {
			t.Fatalf("empty store planned %v", targetIDs(plan))
		}
	})
}

// TestEvidenceCleanupPlan is CE78 and CE80. A plan names every exact deletion target with a
// retrievable identity, and planning deletes nothing.
func TestEvidenceCleanupPlan(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	identity := publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	name := orphan(t, dir, "partial bytes\n")

	plan := planOf(t, store)
	want := []string{ce.TargetOrphan + " " + name, ce.TargetPublished + " " + identity}
	if strings.Join(targetIDs(plan), ",") != strings.Join(want, ",") {
		t.Fatalf("plan targets = %v, want %v", targetIDs(plan), want)
	}
	if plan.Fingerprint == "" || !strings.HasPrefix(plan.Fingerprint, ce.IdentityPrefix) {
		t.Fatalf("plan fingerprint = %q", plan.Fingerprint)
	}
	var bytes uint64
	for _, target := range plan.Targets {
		if target.Bytes == 0 {
			t.Errorf("target %s names no byte length", target.ID)
		}
		bytes += target.Bytes
	}
	if plan.Bytes != bytes {
		t.Errorf("plan bytes = %d, want the sum %d", plan.Bytes, bytes)
	}
	// CE80: the plan is an enumeration, so the store still holds every target.
	if got := storeEntries(t, dir); len(got) != 2 {
		t.Fatalf("planning changed the store to %v", got)
	}
	if planOf(t, store).Fingerprint != plan.Fingerprint {
		t.Fatal("a second plan of an unchanged store changed its fingerprint")
	}
}

// TestEvidenceCleanupStale is CE79. Apply reconstructs the plan under the exclusive lock and
// refuses a fingerprint the current inventory does not produce, before it deletes anything.
func TestEvidenceCleanupStale(t *testing.T) {
	for _, change := range []struct {
		name  string
		apply func(t *testing.T, store *ce.Store, dir string)
	}{
		{"a new target", func(t *testing.T, _ *ce.Store, dir string) { orphan(t, dir, "later bytes\n") }},
		{"a changed length", func(t *testing.T, _ *ce.Store, dir string) {
			name := storeEntries(t, dir)[0]
			if err := os.WriteFile(filepath.Join(dir, name), []byte("replaced"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(change.name, func(t *testing.T) {
			store, dir := newStore(t, ce.StoreOptions{})
			publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
			stale := planOf(t, store).Fingerprint
			change.apply(t, store, dir)

			before := storeEntries(t, dir)
			_, err := store.Apply(stale)
			if refusalClass(err) != ce.RefuseStalePlan {
				t.Fatalf("apply of a stale plan = %v, want %s", err, ce.RefuseStalePlan)
			}
			if got := storeEntries(t, dir); strings.Join(got, ",") != strings.Join(before, ",") {
				t.Fatalf("the refused apply changed the store to %v, want %v", got, before)
			}
		})
	}
}

// TestEvidenceCleanupApply is CE81. Apply removes exactly the planned targets, and the store
// then holds no artifact to retrieve.
func TestEvidenceCleanupApply(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	identity := publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	orphan(t, dir, "partial bytes\n")

	plan := planOf(t, store)
	applied, err := store.Apply(plan.Fingerprint)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if !applied.Complete || applied.Removed != 2 || applied.Remaining != 0 || applied.Fingerprint != plan.Fingerprint {
		t.Fatalf("applied = %+v, want two removed and complete", applied)
	}
	if got := storeEntries(t, dir); len(got) != 0 {
		t.Fatalf("apply left %v", got)
	}
	if _, err := store.Open(identity); refusalClass(err) != ce.RefuseAbsent {
		t.Fatalf("read of the removed artifact = %v, want %s", err, ce.RefuseAbsent)
	}
	if len(planOf(t, store).Targets) != 0 {
		t.Fatal("the emptied store still plans a target")
	}
}

// TestEvidenceCleanupFailure is CE82 and CE122. A deletion failure stops the apply, reports
// the exact unfinished disposition, and leaves the remaining targets for a fresh plan whose
// fingerprint differs from the one the stopped apply used.
func TestEvidenceCleanupFailure(t *testing.T) {
	calls := 0
	fail := func(step string) error {
		if step != ce.StepRemove {
			return nil
		}
		if calls++; calls == 2 {
			return errors.New("injected remove failure")
		}
		return nil
	}
	store, dir := newStore(t, ce.StoreOptions{Fault: fail})
	publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	orphan(t, dir, "partial bytes\n")

	plan := planOf(t, store)
	applied, err := store.Apply(plan.Fingerprint)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if applied.Complete || applied.Removed != 1 || applied.Remaining != 1 {
		t.Fatalf("stopped apply = %+v, want one removed and one remaining", applied)
	}
	if got := storeEntries(t, dir); len(got) != 1 {
		t.Fatalf("stopped apply left %v, want the one target it did not reach", got)
	}
	// CE122: the remaining target needs a fresh plan, and the old fingerprint no longer
	// describes the store.
	fresh := planOf(t, store)
	if fresh.Fingerprint == plan.Fingerprint {
		t.Fatal("the interrupted store still produces the old fingerprint")
	}
	if _, err := store.Apply(plan.Fingerprint); refusalClass(err) != ce.RefuseStalePlan {
		t.Fatalf("reuse of the old fingerprint = %v, want %s", err, ce.RefuseStalePlan)
	}
}

// replace puts a different file under name without reusing the object the plan observed. It
// writes a sibling of the same byte length and renames it, so the store object keeps its
// name and its length and takes only a new file identity. That is the one difference the
// recheck has to read, and it is what an out-of-protocol writer leaves behind. replace
// returns the bytes it wrote.
func replace(t *testing.T, dir, name string) string {
	t.Helper()
	info, err := os.Lstat(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	body := strings.Repeat("!", int(info.Size()))
	sibling := filepath.Join(dir, name+".replacement")
	if err := os.WriteFile(sibling, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(sibling, filepath.Join(dir, name)); err != nil {
		t.Fatal(err)
	}
	return body
}

// TestEvidenceCleanupReplacedTarget is CE82 and the pre-deletion identity recheck. The
// exclusive operation lock excludes every Bench reader and writer, so this test models the
// process outside that protocol: it replaces the second target at the stage the apply
// reaches after it removed the first one. The apply removes exactly the target it planned,
// leaves the file it did not plan, and reports the exact unfinished disposition.
func TestEvidenceCleanupReplacedTarget(t *testing.T) {
	var dir, published, intruder string
	paused := 0
	store, storeDir := newStore(t, ce.StoreOptions{Pause: func(stage string) {
		// The stage repeats once per removed target; this store holds two, and the apply
		// stops at the second, so the replacement happens exactly once.
		if stage != ce.StageRemoved || paused > 0 {
			return
		}
		paused++
		// The stage follows a deletion, so the first target is already gone from the store
		// and only the target this pause replaces is left.
		if got := storeEntries(t, dir); strings.Join(got, ",") != published {
			t.Errorf("the removed stage holds %v, want only the unreached target %s", got, published)
		}
		intruder = replace(t, dir, published)
	}})
	dir = storeDir
	identity := publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	published = strings.TrimPrefix(identity, ce.IdentityPrefix) + ce.PackSuffix
	orphan(t, dir, "partial bytes\n")

	plan := planOf(t, store)
	applied, err := store.Apply(plan.Fingerprint)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if paused != 1 {
		t.Fatalf("the apply reached the removed stage %d times, want one removal before it stopped", paused)
	}
	if applied.Complete || applied.Removed != 1 || applied.Remaining != 1 {
		t.Fatalf("stopped apply = %+v, want one removed and one remaining", applied)
	}
	// The replacement is not the file the plan committed to, so the apply left it in place.
	body, err := os.ReadFile(filepath.Join(dir, published))
	if err != nil || string(body) != intruder {
		t.Fatalf("the apply deleted the replaced target: %v", err)
	}
	if got := storeEntries(t, dir); strings.Join(got, ",") != published {
		t.Fatalf("stopped apply left %v, want only the replaced target", got)
	}
	if fresh := planOf(t, store); fresh.Fingerprint == plan.Fingerprint {
		t.Fatal("the changed store still produces the old fingerprint")
	}
}

// TestEvidenceCleanupUnsafe is CE121. An unsafe target refuses the whole plan before any
// deletion, and the refusal never follows the object it names.
func TestEvidenceCleanupUnsafe(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("outside\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, ce.TempPrefix+"link"+ce.TempSuffix)); err != nil {
		t.Fatal(err)
	}

	before := storeEntries(t, dir)
	if _, err := store.Plan(); refusalClass(err) != ce.RefuseUnsafe {
		t.Fatalf("plan with an unsafe target = %v, want %s", err, ce.RefuseUnsafe)
	}
	if _, err := store.Apply(ce.IdentityPrefix + strings.Repeat("0", 64)); refusalClass(err) != ce.RefuseUnsafe {
		t.Fatalf("apply with an unsafe target = %v, want %s", err, ce.RefuseUnsafe)
	}
	if got := storeEntries(t, dir); strings.Join(got, ",") != strings.Join(before, ",") {
		t.Fatalf("the refusal changed the store to %v, want %v", got, before)
	}
	if _, err := os.Lstat(outside); err != nil {
		t.Fatalf("the refusal followed the symlink: %v", err)
	}
}

// TestEvidenceCleanupOrphans is CE85. A temporary object becomes a plannable orphan only
// while cleanup holds the operation lock exclusively, so a live writer's temporary is never
// planned.
func TestEvidenceCleanupOrphans(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	name := orphan(t, dir, "partial bytes\n")
	if got := targetIDs(planOf(t, store)); got[0] != ce.TargetOrphan+" "+name {
		t.Fatalf("plan targets = %v, want the orphan first", got)
	}

	// A staged writer holds the operation lock, which cleanup needs exclusively. Cleanup
	// refuses rather than planning a temporary object that writer still owns.
	staged, err := store.Stage(mustBuild(t, fixtureCandidate("# Two\n")), ce.DefaultQuota, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer staged.Discard()
	if _, err := store.Plan(); refusalClass(err) != ce.RefuseBusy {
		t.Fatalf("plan during a live writer = %v, want %s", err, ce.RefuseBusy)
	}
	if _, err := store.Apply(ce.IdentityPrefix + strings.Repeat("0", 64)); refusalClass(err) != ce.RefuseBusy {
		t.Fatalf("apply during a live writer = %v, want %s", err, ce.RefuseBusy)
	}
}
