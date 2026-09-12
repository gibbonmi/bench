package gate

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// delegatedCheckpointFixture is one recorded delegated chunk on a committed
// source. It reuses the version 1 build sequence through its attach function.
func delegatedCheckpointFixture(t *testing.T) *recordtest.Fixture {
	return attachedCheckpointFixture(t, recordtest.AttachDelegated)
}

// resave writes the mutated record and commits it, so the checkpoint grades the
// committed tree rather than an in-memory value.
func resave(t *testing.T, f *recordtest.Fixture, message string) {
	t.Helper()
	f.Save()
	f.Commit(message)
}

// DI4: a chunk checkpoint requires each ticket's assigned-author verification
// on the integrated source.
func TestDelegatedChunkVerifier(t *testing.T) {
	t.Run("valid delegated evidence passes", func(t *testing.T) {
		f := delegatedCheckpointFixture(t)
		if code, out := runCheckpoint(t, f); code != 0 {
			t.Fatalf("valid delegated checkpoint refused: %d %s", code, out)
		}
	})
	cases := []struct {
		name, reason string
		change       func(*recordtest.Fixture)
	}{
		{"orchestrator verification", "verification tests", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Verification[0].Performer = recordtest.Orchestrator
		}},
		{"foreign author verification", "verification tests", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Verification[0].Performer = "outside-the-run"
		}},
		{"branch-only source", "verification tests", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Verification[0].SourceDigest = f.Record.Chunks[0].Base
		}},
		{"integration role on a chunk obligation", "verification tests", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Verification[0].Role = "integration-verification"
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := delegatedCheckpointFixture(t)
			tc.change(f)
			resave(t, f, "mutate the delegated verification")
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, tc.reason) {
				t.Fatalf("invalid delegated verification accepted: %d %s", code, out)
			}
		})
	}
}

// DI5: a current or former author cannot review any chunk.
func TestDelegatedAuthorReview(t *testing.T) {
	for _, who := range []string{recordtest.Author("1.md"), recordtest.Orchestrator} {
		t.Run(who, func(t *testing.T) {
			f := delegatedCheckpointFixture(t)
			f.Record.Chunks[0].Reviews[0].Performer = who
			resave(t, f, "let a participant review")
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "invalid Standards performer") {
				t.Fatalf("a run participant supplied independent review: %d %s", code, out)
			}
		})
	}
}

// DI6: the orchestrator cannot supply independent review even when the axis is
// otherwise complete. A version 1 record keeps its own single exclusion.
func TestDelegatedOrchestratorReview(t *testing.T) {
	f := delegatedCheckpointFixture(t)
	for i := range f.Record.Chunks[0].Reviews {
		f.Record.Chunks[0].Reviews[i].Performer = recordtest.Orchestrator
	}
	resave(t, f, "let the orchestrator review every axis")
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "invalid Standards performer") {
		t.Fatalf("a complete orchestrator-authored review was accepted: %d %s", code, out)
	}
}

// DI7: one session cannot supply two axes for a chunk.
func TestDelegatedDistinctAxes(t *testing.T) {
	f := delegatedCheckpointFixture(t)
	f.Record.Chunks[0].Reviews[1].Performer = f.Record.Chunks[0].Reviews[0].Performer
	resave(t, f, "reuse one reviewer for two axes")
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "three distinct review sessions") {
		t.Fatalf("one session satisfied two axes: %d %s", code, out)
	}
}

// DI12: delegated checkpoints retain the existing terminal-evidence refusals.
func TestDelegatedEvidenceRefusals(t *testing.T) {
	cases := []struct {
		name, reason string
		change       func(*recordtest.Fixture)
	}{
		{"missing axis", "missing Coverage", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
		}},
		{"pending axis", "pending", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Reviews[0].State = "pending"
		}},
		{"failed probe", "probe", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Verification[0].Probe.Outcome = "silent"
		}},
		{"failed restore", "restore", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Verification[0].Probe.Restore = "failed"
		}},
		{"altered excerpt", "native result", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Reviews[0].NativeRef.Excerpt = "rewritten after the fact"
		}},
		{"stale frozen pair", "stale Standards", func(f *recordtest.Fixture) {
			f.Record.Chunks[0].Reviews[0].Tip = f.Record.Chunks[0].Base
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := delegatedCheckpointFixture(t)
			tc.change(f)
			resave(t, f, "mutate the delegated evidence")
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, tc.reason) {
				t.Fatalf("a delegated checkpoint lost an existing refusal: %d %s", code, out)
			}
		})
	}
}
