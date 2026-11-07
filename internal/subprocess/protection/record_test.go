package protection_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// resourceRecord returns the path and the bytes of the record file in d's store.
func resourceRecord(t *testing.T, d protection.Descriptor, file string) (string, string) {
	t.Helper()
	path := filepath.Join(d.Store, file)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return path, string(data)
}

// plant replaces whatever is at path with a private file holding data.
func plant(t *testing.T, path, data string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

// requireRetiredKept damages the retired record of a store whose resource a removal deleted,
// and fails unless a sweep of that removal keeps the store with no error.
func requireRetiredKept(t *testing.T, damage func(path, data string)) {
	t.Helper()
	p, d, dir := goneRetired(t)
	damage(resourceRecord(t, d, protection.RetiredFile))
	requireSweep(t, p, dir, d, true)
}

// noStartUser leaves one real no-start record in a store, as a failed second publication's
// rollback does, and returns that store and record path.
func noStartUser(t *testing.T, p *protection.Protocol) (protection.Descriptor, string) {
	t.Helper()
	first, second := ordered(t, p)
	faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "B/user.json", Ordinal: 1}: {"rename": errInjected}}, map[string]string{first.Store: "A", second.Store: "B"})
	protection.SetReplacer(p, faults.Replacer())
	requireRefused(t, p, first, second)
	protection.SetReplacer(p, durablefile.Replacer{})
	for id := range userRecords(t, first) {
		return first, filepath.Join(first.Store, protection.UsersDir, id, protection.UserFile)
	}
	t.Fatal("the rollback left no no-start record")
	return first, ""
}

func TestRecordForeignBinding(t *testing.T) {
	p := newProtocol()
	first, second := ordered(t, p)
	foreignStore, foreignGeneration := first, first
	foreignStore.Store = second.Store
	foreignGeneration.ResourceID = second.ResourceID
	for _, d := range []protection.Descriptor{foreignStore, foreignGeneration} {
		requireRefused(t, p, d)
	}
	// The validator admits first, so only the published generation's foreign binding refuses it.
	path, data := resourceRecord(t, first, protection.ResourceFile)
	plant(t, path, strings.Replace(data, first.Binding, second.Binding, 1))
	requireRefused(t, p, first)
	if users := len(userRecords(t, first)) + len(userRecords(t, second)); users != 0 {
		t.Fatalf("a foreign descriptor registered %d users", users)
	}
}

// TestRecordKindWithoutValidator launches a real generation. A kind with no validator row gains
// no admission, whether the table is nil, empty, or another kind's.
func TestRecordKindWithoutValidator(t *testing.T) {
	d := newResource(t, newProtocol(), "cache")
	for _, table := range []protection.Validators{nil, {}, {protection.KindWorktree: validate}} {
		requireRefused(t, &protection.Protocol{Validators: table, ReadLimit: readLimit}, d)
	}
}

func TestRecordMalformed(t *testing.T) {
	malformed := func(path, _ string) { plant(t, path, "{not json") }
	p := newProtocol()
	d, record := noStartUser(t, p)
	malformed(record, "")
	requireRefused(t, p, d)
	requireRetained(t, p, d)
	requireRetiredKept(t, malformed)
}

func TestRecordSymlink(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "outside.json")
	link := func(path, data string) {
		plant(t, outside, data)
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, path); err != nil {
			t.Fatal(err)
		}
	}
	p := newProtocol()
	d := newResource(t, p, "cache")
	path, data := resourceRecord(t, d, protection.ResourceFile)
	link(path, data)
	requireRefused(t, p, d)
	requireRetained(t, p, d)
	if after, err := os.ReadFile(outside); err != nil || string(after) != data {
		t.Fatalf("the outside sentinel changed to %q (%v)", after, err)
	}
	requireRetiredKept(t, link)
}

func TestRecordTrailingJSON(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	path, data := resourceRecord(t, d, protection.ResourceFile)
	plant(t, path, data+`{"admission":"sealed"}`)
	requireRefused(t, p, d)
}

func TestRecordUnknownField(t *testing.T) {
	unknown := func(path, data string) { plant(t, path, strings.Replace(data, "{", `{"authority":"granted",`, 1)) }
	p := newProtocol()
	d := newResource(t, p, "cache")
	unknown(resourceRecord(t, d, protection.ResourceFile))
	requireRefused(t, p, d)
	requireRetiredKept(t, unknown)
}

func TestRecordWrongMode(t *testing.T) {
	wide := func(path, _ string) {
		if err := os.Chmod(path, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	p := newProtocol()
	d, record := noStartUser(t, p)
	wide(record, "")
	requireRefused(t, p, d)
	requireRetained(t, p, d)
	requireRetiredKept(t, wide)
}

func TestRecordDuplicateUser(t *testing.T) {
	p := newProtocol()
	d, record := noStartUser(t, p)
	live := startUser(t, p, d)
	// A second directory claims the live user's identity with a no-start record.
	var fields map[string]any
	data, err := os.ReadFile(record)
	if err == nil {
		err = json.Unmarshal(data, &fields)
	}
	if err != nil {
		t.Fatal(err)
	}
	fields["user_id"] = live.reg.Disposition().Users[0]
	claim, _ := json.Marshal(fields)
	duplicate := filepath.Join(d.Store, protection.UsersDir, strings.Repeat("f", 32))
	if err := os.Mkdir(duplicate, 0o700); err != nil {
		t.Fatal(err)
	}
	plant(t, filepath.Join(duplicate, protection.UserFile), string(claim))
	if disposition := requireRetained(t, p, d); disposition.State != protection.Uncertain {
		t.Fatalf("a duplicate user identity graded %+v, want uncertain", disposition)
	}
}

func TestRecordMissingDescriptorTarget(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	path, _ := resourceRecord(t, d, protection.ResourceFile)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	requireRefused(t, p, d)
}

func TestRecordFIFO(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	path, _ := resourceRecord(t, d, protection.ResourceFile)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	requireRefused(t, p, d)
}

func TestRecordDirectory(t *testing.T) {
	p := newProtocol()
	d, record := noStartUser(t, p)
	if err := os.Remove(record); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(record, 0o700); err != nil {
		t.Fatal(err)
	}
	requireRefused(t, p, d)
	requireRetained(t, p, d)
}

func TestRecordWrongSchema(t *testing.T) {
	wrong := func(path, data string) { plant(t, path, strings.Replace(data, `"schema":1`, `"schema":2`, 1)) }
	p := newProtocol()
	d := newResource(t, p, "cache")
	wrong(resourceRecord(t, d, protection.ResourceFile))
	requireRefused(t, p, d)
	requireRetiredKept(t, wrong)
}

// TestRecordRetiredPath forges the path of a retired record that the strict reader decodes.
// The unclean and control-byte forms stay inside the removed directory and absent, so only the
// record's path check keeps their stores.
func TestRecordRetiredPath(t *testing.T) {
	for name, forge := range map[string]func(string) string{
		"relative":     func(path string) string { return strings.TrimPrefix(path, "/") },
		"unclean":      func(path string) string { return filepath.Dir(path) + "/./" + filepath.Base(path) },
		"control byte": func(path string) string { return path + "\n" },
	} {
		t.Run(name, func(t *testing.T) {
			requireRetiredKept(t, func(path, data string) {
				var fields map[string]any
				if err := json.Unmarshal([]byte(data), &fields); err != nil {
					t.Fatal(err)
				}
				fields["resource_path"] = forge(fields["resource_path"].(string))
				forged, _ := json.Marshal(fields)
				plant(t, path, string(forged))
			})
		})
	}
}

func TestRecordDuplicateField(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	path, data := resourceRecord(t, d, protection.ResourceFile)
	plant(t, path, strings.Replace(data, `"admission":"open"`, `"admission":"sealed","admission":"open"`, 1))
	requireRefused(t, p, d)
}

func TestRecordOversized(t *testing.T) {
	oversized := func(path, data string) { plant(t, path, data+strings.Repeat(" ", readLimit)) }
	p := newProtocol()
	d := newResource(t, p, "cache")
	oversized(resourceRecord(t, d, protection.ResourceFile))
	requireRefused(t, p, d)
	requireRetiredKept(t, oversized)
}

func TestRecordGenerationConflict(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	stale := d
	stale.ResourceID = strings.Repeat("0", 32)
	requireRefused(t, p, d, stale)
	if records := userRecords(t, d); len(records) != 0 {
		t.Fatalf("conflicting generations registered users %v", records)
	}
}

func TestRecordUnsafePath(t *testing.T) {
	p := newProtocol()
	path := filepath.Join(t.TempDir(), "ca\nche")
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if d, err := p.OpenResource(candidate(path)); err == nil {
		t.Errorf("a control-byte path opened generation %+v", d)
	}
	unsafe := candidate(path)
	unsafe.ResourceID = strings.Repeat("0", 32)
	_, parsed := protection.ParseVector(strings.TrimPrefix(vectorEnv(unsafe), protection.Env+"="))
	if parsed == nil {
		t.Fatal("an inherited control-byte path parsed")
	}
	// Prepare's own descriptor check refuses the launch with the refusal that the inherited vector gets.
	if err, check := requireRefused(t, p, unsafe), errors.Unwrap(parsed); err.Error() != check.Error() {
		t.Fatalf("the launch refusal %q is not the descriptor check's %q", err, check)
	}
}

func TestRecordAwkwardPath(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "a b*[c]?")
	vector, err := protection.ParseVector(strings.TrimPrefix(vectorEnv(d), protection.Env+"="))
	if err != nil || !slices.Equal(vector, []protection.Descriptor{d}) {
		t.Fatalf("vector round trip = %+v (%v), want %+v", vector, err, d)
	}
	o, marker := launch(t, p, vector, `cat "$1"/*/user.json > "$0"`, filepath.Join(d.Store, protection.UsersDir))
	seen, err := os.ReadFile(marker)
	if err != nil || o.Protection.State != protection.Safe || len(o.Protection.Users) != 1 || !strings.Contains(string(seen), o.Protection.Users[0]) {
		t.Fatalf("awkward path: protection=%+v seen=%q (%v)", o.Protection, seen, errors.Join(err, o.Process.StartErr))
	}
}
