package backgroundstate_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
)

// TestBackgroundRecordNamesTheOutcome is BG35. The background record of a worktree names the
// outcome, the grade key, the tip, and the log path, for each outcome.
func TestBackgroundRecordNamesTheOutcome(t *testing.T) {
	for _, outcome := range []backgroundstate.Outcome{backgroundstate.Green, backgroundstate.Red, backgroundstate.Aborted, backgroundstate.Stopped} {
		gitdir := t.TempDir()
		record := backgroundstate.Record{Outcome: outcome, Key: key, Tip: "1111111111111111111111111111111111111111", Log: "/logs/grade one.log"}
		if err := backgroundstate.WriteRecord(durablefile.Native(), gitdir, record); err != nil {
			t.Fatal(err)
		}
		want := map[string]any{"outcome": string(outcome), "key": record.Key, "tip": record.Tip, "log": record.Log}
		if got := fields(t, backgroundstate.RecordPath(gitdir)); !reflect.DeepEqual(got, want) {
			t.Fatalf("record fields = %v, want %v", got, want)
		}
		read, err := backgroundstate.ReadRecord(gitdir)
		if err != nil || read != record {
			t.Fatalf("ReadRecord = %+v, %v; want %+v", read, err, record)
		}
	}
}

// TestQueueRoundTrips writes the slot and the queued grades in their order and reads them
// back unchanged.
func TestQueueRoundTrips(t *testing.T) {
	home := t.TempDir()
	running := backgroundstate.Slot{Entry: backgroundstate.Entry{Key: key, Worktree: "/work/a", Tip: "a"}, Owner: 4242, Boot: "boot"}
	queue := backgroundstate.Queue{Running: &running, Queued: []backgroundstate.Entry{
		{Key: "b" + key[1:], Worktree: "/work/b", Tip: "b"},
		{Key: "c" + key[1:], Worktree: "/work/c", Tip: "c"},
	}, Latest: map[string]string{"/work/b": "b" + key[1:]}}
	if err := backgroundstate.WriteQueue(durablefile.Native(), home, queue); err != nil {
		t.Fatal(err)
	}
	read, err := backgroundstate.ReadQueue(home)
	if err != nil || !reflect.DeepEqual(read, queue) {
		t.Fatalf("ReadQueue = %+v, %v; want %+v", read, err, queue)
	}
}

// TestAbsentRecordsReadEmptyAndUnreadableOnesRefuse proves the read discipline of every
// record. An absent record reads as empty. A present record that does not decode, or that
// decodes to an invalid value, is an error, so its caller retains the grade.
func TestAbsentRecordsReadEmptyAndUnreadableOnesRefuse(t *testing.T) {
	dir := t.TempDir()
	if claim, err := backgroundstate.ReadClaim(dir, key); err != nil || claim != (backgroundstate.Claim{}) {
		t.Fatalf("absent claim = %+v, %v; want empty", claim, err)
	}
	if record, err := backgroundstate.ReadRecord(dir); err != nil || record != (backgroundstate.Record{}) {
		t.Fatalf("absent record = %+v, %v; want empty", record, err)
	}
	if queue, err := backgroundstate.ReadQueue(dir); err != nil || !reflect.DeepEqual(queue, backgroundstate.Queue{}) {
		t.Fatalf("absent queue = %+v, %v; want empty", queue, err)
	}
	reads := map[string]struct {
		path string
		read func() error
	}{
		"claim":  {backgroundstate.ClaimPath(dir, key), func() error { _, err := backgroundstate.ReadClaim(dir, key); return err }},
		"record": {backgroundstate.RecordPath(dir), func() error { _, err := backgroundstate.ReadRecord(dir); return err }},
		"queue":  {backgroundstate.QueuePath(dir), func() error { _, err := backgroundstate.ReadQueue(dir); return err }},
	}
	invalid := map[string]string{
		"claim":  `{"key":"` + key + `","worktree":"w","owner_pid":0,"boot_id":"b"}`,
		"record": `{"outcome":"blue","key":"` + key + `","tip":"t","log":"l"}`,
		"queue":  `{"running":null,"queued":[{"key":"a/b","worktree":"w","tip":"t"}]}`,
	}
	for name, r := range reads {
		for label, body := range map[string]string{"undecodable": "{\n", "invalid": invalid[name] + "\n"} {
			if err := os.MkdirAll(filepath.Dir(r.path), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(r.path, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := r.read(); err == nil {
				t.Errorf("%s %s record read with no error", label, name)
			}
		}
	}
}
