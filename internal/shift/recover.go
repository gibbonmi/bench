package shift

import (
	"context"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/worktree"
	"go.opentelemetry.io/otel/attribute"
)

// recoverySeam is the seam of a recovery verdict's span.
const recoverySeam = "shift.recovery"

// Recover is the recovery pass over root's open shift entries: the entries with no
// outcome. An entry with no lease whose key names a gone owner process is abandoned.
// A live owner, a key no shift wrote, and an entry with a lease stay as they are. An
// abandon changes no worktree file, lease, or lock. The pass prints one line only when it
// acted, so a quiet session start stays quiet.
func Recover(root string, out io.Writer) {
	current, err := intent.Read(root)
	if err != nil {
		return
	}
	recovered, abandoned := 0, 0
	for _, entry := range current.Entries {
		if entry.Kind != intent.KindShift || entry.Outcome != "" || entry.Lease != "" {
			continue
		}
		pid, ok := intent.KeyOwner(intent.KindShift, entry.Key)
		if !ok || worktree.PIDAlive(pid) {
			continue
		}
		entry.Outcome = otelrecord.WorkAbandoned
		if intent.Upsert(root, entry) != nil {
			continue
		}
		recordRecovery(root, entry.Key, otelrecord.WorkAbandoned, otelrecord.CleanupNone)
		abandoned++
	}
	if recovered+abandoned > 0 {
		fmt.Fprintf(out, "bench shift recovery: recovered %d, abandoned %d\n", recovered, abandoned)
	}
}

// recordRecovery writes one shift.recovery span for a verdict on the entry key.
func recordRecovery(root, key, state, cleanup string) {
	_, _, end := otelrecord.BeginIn(context.Background(), "", root, recoverySeam, recoverySeam,
		attribute.String(otelrecord.AttrIntentKey, key),
		attribute.String(otelrecord.AttrWorkState, state),
		attribute.String(otelrecord.AttrCleanup, cleanup))
	end()
}
