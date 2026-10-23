//go:build darwin

package runbinary

import (
	"context"
	"errors"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestBuilderDarwinPolicy is PL90. Darwin carries no parent-death signal, so the builder's
// own process group is the whole removal path. The builder returns while a descendant that
// ignores INT and TERM stays parked in that group. Own must remove that descendant through
// the group before it returns. Without the group setup, the group signal finds no group
// and the descendant survives.
func TestBuilderDarwinPolicy(t *testing.T) {
	t.Parallel()
	builder := newParkingBuilder(t, false)
	factory := Factory{TempRoot: t.TempDir(), Protection: binaryProtocol(), Verify: func(string, string) error { return nil }}
	selection, err := factory.Own(context.Background(), builder.source)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = selection.Close() })
	child := readPID(t, builder.childPID)
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	if state := protection.ObserveGroup(readPID(t, builder.builderPID)).State; state != protection.GroupAbsent {
		t.Fatalf("builder group at Own's return is %s, want absent", state)
	}
	if err := syscall.Kill(child, 0); !errors.Is(err, syscall.ESRCH) {
		t.Fatalf("resistant builder descendant %d at Own's return answers %v, want absent", child, err)
	}
}
