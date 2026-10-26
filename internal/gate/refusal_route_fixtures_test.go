// Refusal route fixtures for the gate checkpoint: each fixture produces one gate face for
// the registry walk in refusal_route_test.go.
package gate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// evidenceFixture is the completion-evidence fixture: a delivery that the preflight read
// grades green, with its review record removed in a commit. The read reports the gap, and the
// operator commits the record again. The chunk is the one the conformant plan declares.
func evidenceFixture(record string) gateFaceFixture {
	var saved []byte
	return gateFaceFixture{
		face:     "checkpoint-completion-evidence",
		args:     []string{"--checkpoint", recordtest.Spec, "--chunk", "c1"},
		contains: "bench preflight review ",
		build: func(t *testing.T) gateSet {
			t.Setenv("BENCH_GATE", "true")
			root, _ := preflighttest.SeedConformant(t)
			recordtest.RetainSingleChunk(t, root, recordtest.Spec, preflighttest.RunGit(t, "rev-parse", "main"))
			saved = readFile(t, root, record)
			if err := os.Remove(filepath.Join(root, record)); err != nil {
				t.Fatal(err)
			}
			f := &recordtest.Fixture{T: t, Root: root}
			f.Commit("remove the review record")
			return gateSet{f: f}
		},
		after: func(_ *testing.T, set gateSet) {
			set.f.Write(record, string(saved))
			set.f.Commit("restore the review record")
		},
	}
}

// repairFixture is the record-repair fixture: the chunk lacks its Coverage review, which
// the review form writes. The route prints that form's grammar, and the operator retains
// the review.
func repairFixture(record string) gateFaceFixture {
	var saved []byte
	return gateFaceFixture{
		face:     "checkpoint-record-repair",
		args:     chunkArgs,
		contains: "bench record 'review' --help",
		build: func(t *testing.T) gateSet {
			set := assigned(t, "gate-record", false)
			saved = readFile(t, set.f.Root, record)
			set.f.Record.Chunks[0].Reviews = set.f.Record.Chunks[0].Reviews[:2]
			set.f.Save()
			set.f.Commit("drop the Coverage review")
			return set
		},
		after: func(_ *testing.T, set gateSet) {
			set.f.Write(record, string(saved))
			set.f.Commit("restore the Coverage review")
		},
	}
}

// brokerFixture is the broker fixture: a completed source in a repository that declares Bench
// build inputs, with no broker manifest. The operator's repair publishes a manifest that binds
// the running executable, which the fixture carries out by binding this test binary, so the
// rerun keeps its own process.
func brokerFixture() gateFaceFixture {
	const label = "gate-broker"
	return gateFaceFixture{
		face:     "checkpoint-broker-unavailable",
		args:     completeArgs,
		contains: "republish the promotion broker; then bench gate --in " + sanitize.ShellQuote(label) + " ",
		build:    func(t *testing.T) gateSet { return assigned(t, label, true, gate.DeclareBuildInputs) },
		carry:    map[int]func(*testing.T, gateSet){0: func(t *testing.T, set gateSet) { installRunning(t, set.f.Root) }},
	}
}
