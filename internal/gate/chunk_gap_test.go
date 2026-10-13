// Chunk gap tests: a green chunk checkpoint with a successor chunk opens the gap of the build
// that the checkout's assignment binds.
package gate_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// chunkCheckpoint is the recorded first chunk of a plan of chunks chunks, with the build of
// its approved spec open. The fixture is the checkout of the build's assignment when bound,
// and a checkout that no assignment owns otherwise.
func chunkCheckpoint(t *testing.T, chunks int, bound bool) *recordtest.Fixture {
	t.Helper()
	t.Setenv("BENCH_HOME", t.TempDir())
	f := gate.AttachedCheckpointFixture(t, func(tb testing.TB, root string, _ int) *recordtest.Fixture { return linkedAttach(tb, root, chunks) })
	build := *f
	build.Root = commitmenttest.BindBuild(t, f.Root, "build", recordtest.Spec)
	if bound {
		return &build
	}
	return f
}

// checkpointChunk runs the gate verb's checkpoint of chunk 1 at the root of f.
func checkpointChunk(t *testing.T, f *recordtest.Fixture) (int, string) {
	t.Helper()
	var out, diagnostic bytes.Buffer
	code := gate.RunCommand([]string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1"}, &out, &diagnostic)
	return code, out.String() + diagnostic.String()
}

// requireGreen fails unless the checkpoint of f exits 0.
func requireGreen(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	if code, out := checkpointChunk(t, f); code != 0 {
		t.Fatalf("checkpoint = exit %d: %s; want green", code, out)
	}
}

// requireNoGap fails unless the ledger at root holds no gap record.
func requireNoGap(t *testing.T, root string) {
	t.Helper()
	if gaps, _ := commitmenttest.ClaimedGaps(t, root); len(gaps) != 0 {
		t.Fatalf("gap records after the checkpoint = %+v, want none", gaps)
	}
}

// requireOpenGap fails unless the ledger at root holds only the open gap of chunk 1.
func requireOpenGap(t *testing.T, root string) {
	t.Helper()
	gaps, claim := commitmenttest.ClaimedGaps(t, root)
	if want := []intent.ChunkGap{{Milestone: claim.Milestone, Outcome: claim.Outcome, Chunk: "1", Open: true}}; !reflect.DeepEqual(gaps, want) {
		t.Fatalf("gap records = %+v, want %+v", gaps, want)
	}
}

// TestGreenChunkCheckpointOpensAGap is BW60. The green checkpoint of chunk 1, which chunk 2
// follows, in the bound checkout writes an open gap for chunk 1.
func TestGreenChunkCheckpointOpensAGap(t *testing.T) {
	f := chunkCheckpoint(t, 2, true)
	requireGreen(t, f)
	requireOpenGap(t, f.Root)
}

// TestEarlierChunkCheckpointKeepsTheGap is BW99. A green checkpoint of chunk 1 keeps the
// ledger bytes while the gap record holds the last chunk 3 with Last true.
func TestEarlierChunkCheckpointKeepsTheGap(t *testing.T) {
	f := chunkCheckpoint(t, 3, true)
	commitmenttest.PutGap(t, f.Root, intent.ChunkGap{Outcome: commitmenttest.DeliveryOutcome, Chunk: "3", Last: true})
	commitmenttest.RequireLedgerKept(t, f.Root, func() { requireGreen(t, f) })
}

// TestUnrecordedGapExitsOne is the gap write failure. A ledger with no read bit fails the
// gap write after a green checkpoint, so the verb exits 1 with the cause. A rerun on a
// readable ledger reuses the green verdict and writes the gap.
func TestUnrecordedGapExitsOne(t *testing.T) {
	f := chunkCheckpoint(t, 2, true)
	path := commitmenttest.LedgerPath(t, f.Root)
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
	if err := os.Chmod(path, 0o200); err != nil {
		t.Fatal(err)
	}
	_, cause := intent.Read(f.Root)
	if cause == nil {
		capability.Capability(t, capability.Privilege, "mode 0o200 ledger is still readable by this user")
	}
	code, out := checkpointChunk(t, f)
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if want := "chunk gap not recorded: " + fmt.Sprint(cause); code != 1 || !strings.Contains(out, want) {
		t.Fatalf("checkpoint = exit %d: %s; want exit 1 and %q", code, out, want)
	}
	requireNoGap(t, f.Root)
	if code, out := checkpointChunk(t, f); code != 0 || !strings.Contains(out, "reused") {
		t.Fatalf("rerun = exit %d: %s; want the reused green verdict", code, out)
	}
	requireOpenGap(t, f.Root)
}

// TestLastChunkOpensNoGap is BW62. The green checkpoint of the plan's last chunk writes no gap.
func TestLastChunkOpensNoGap(t *testing.T) {
	f := chunkCheckpoint(t, 1, true)
	requireGreen(t, f)
	requireNoGap(t, f.Root)
}

// TestRedCheckpointOpensNoGap is BW63. A checkpoint that refuses the chunk's evidence writes
// no gap, although chunk 2 follows.
func TestRedCheckpointOpensNoGap(t *testing.T) {
	f := chunkCheckpoint(t, 2, true)
	f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
	f.Save()
	if code, out := checkpointChunk(t, f); code == 0 || !strings.Contains(out, "missing Coverage") {
		t.Fatalf("checkpoint = exit %d: %s; want the missing Coverage refusal", code, out)
	}
	requireNoGap(t, f.Root)
}

// TestUnboundCheckpointOpensNoGap is BW64. A green checkpoint in a checkout that no
// assignment owns writes no gap for the open build.
func TestUnboundCheckpointOpensNoGap(t *testing.T) {
	f := chunkCheckpoint(t, 2, false)
	requireGreen(t, f)
	requireNoGap(t, f.Root)
}
