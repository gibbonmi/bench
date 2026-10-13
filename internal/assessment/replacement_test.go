package assessment

import (
	"errors"
	"testing"

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

			var replaceErr *durablefile.Error
			if err == nil || errors.As(err, &replaceErr) {
				t.Fatalf("err = %v, want the caller's own refusal, not a replacement failure", err)
			}
			if len(rec.Trace) != 0 {
				t.Fatalf("trace = %v, want no replacement stage", rec.Trace)
			}
			path, err := s.path(r.RunID)
			if err != nil {
				t.Fatal(err)
			}
			durablefiletest.RequireAbsent(t, path)
		})
	}
}
