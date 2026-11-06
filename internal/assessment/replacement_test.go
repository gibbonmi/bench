package assessment

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// recordedOps returns the record operations that route a replacement through rec.
func recordedOps(rec *durablefiletest.Recorder) *FileOps {
	ops := FileOps(rec.Replacer())
	return &ops
}

// D28: a record of a validated run returns the injected directory-sync failure of its
// replacement after the rename.
func TestAssessmentReplacement(t *testing.T) {
	t.Parallel()
	cause := errors.New("injected assessment directory-sync failure")
	s := Store{Home: t.TempDir(), Root: t.TempDir(), Files: recordedOps(&durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}})}
	r := fixtureRun(s.Root)

	err := s.Record(r)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	path, err := s.path(r.RunID)
	if err != nil {
		t.Fatal(err)
	}
	durablefiletest.RequireFile(t, path, 0o600, `{
  "version": 1,
  "run_id": "run-1",
  "repo_key": "`+poolkey.Key(s.Root)+`",
  "source": "synthetic",
  "condition": "bench",
  "task_id": "held-out-1",
  "holdout": false,
  "state": "running",
  "attempts": [
    {
      "attempt_id": "attempt-1",
      "chunk_id": "1",
      "role": "implementation",
      "session_id": "session-1",
      "model": "synthetic",
      "effort": "high",
      "state": "failed",
      "usage": null,
      "cost": {},
      "evidence": null
    }
  ],
  "evidence": null,
  "quality": null
}
`)
}

// TestInstanceCollaboratorsStayLocal proves that a store's file operations belong to that
// instance. Two stores record at the same time, each held at a shared barrier until both
// replacements are in flight. Each must report its own injected failure and leave its own
// trace, so no package state carries one store's collaborator into the other.
func TestInstanceCollaboratorsStayLocal(t *testing.T) {
	t.Parallel()
	var entered sync.WaitGroup
	entered.Add(2)
	both := make(chan struct{})
	go func() { entered.Wait(); close(both) }()
	type instance struct {
		store Store
		rec   *durablefiletest.Recorder
		cause error
		err   error
	}
	instances := make([]*instance, 2)
	for i := range instances {
		cause := fmt.Errorf("injected directory-sync failure of store %d", i)
		rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}
		ops := recordedOps(rec)
		createTemp := ops.CreateTemp
		ops.CreateTemp = func(dir, pattern string) (durablefile.File, error) {
			entered.Done()
			select {
			case <-both:
			case <-time.After(bounds.TestDeadline(0)):
				return nil, errors.New("the other store never reached its replacement")
			}
			return createTemp(dir, pattern)
		}
		instances[i] = &instance{store: Store{Home: t.TempDir(), Root: t.TempDir(), Files: ops}, rec: rec, cause: cause}
	}
	var done sync.WaitGroup
	for _, in := range instances {
		done.Add(1)
		go func() { defer done.Done(); in.err = in.store.Record(fixtureRun(in.store.Root)) }()
	}
	done.Wait()
	for i, in := range instances {
		durablefiletest.RequireFailure(t, in.err, "directory-sync", true, in.cause)
		if len(in.rec.Temps) != 1 || filepath.Dir(in.rec.Temps[0]) != in.store.Dir() {
			t.Fatalf("store %d recorded temporaries %q, want one in its own directory %s", i, in.rec.Temps, in.store.Dir())
		}
	}
}

// D30: an invalid run and a run of a foreign repository refuse before the replacement
// starts, and no record exists.
func TestAssessmentRefusalSkipsReplacement(t *testing.T) {
	t.Parallel()
	for name, change := range map[string]func(*testing.T, *Run){
		"validation": func(_ *testing.T, r *Run) {
			e := usageEvent("unattributed", 1)
			e.Reference = Reference{}
			r.Attempts[0].Usage = []Event{e}
		},
		"authority": func(t *testing.T, r *Run) { r.RepoKey = poolkey.Key(t.TempDir()) },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			rec := &durablefiletest.Recorder{}
			s := Store{Home: t.TempDir(), Root: t.TempDir(), Files: recordedOps(rec)}
			r := fixtureRun(s.Root)
			change(t, &r)

			err := s.Record(r)

			durablefiletest.RequireRefusal(t, err, rec.Trace)
			path, err := s.path(r.RunID)
			if err != nil {
				t.Fatal(err)
			}
			durablefiletest.RequireAbsent(t, path)
		})
	}
}
