package gate

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestGateSignalArmFollowsAcquireAndPrecedesOwnerWrite(t *testing.T) {
	t.Parallel()
	arm := func(engine *faultEngine) postAcquireContextArm {
		return func(ctx context.Context) (context.Context, func()) {
			engine.trace = append(engine.trace, "signal-arm")
			return ctx, func() {}
		}
	}

	t.Run("acquired", func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		engine := &faultEngine{now: time.Date(2026, 7, 23, 12, 0, 0, 0, time.UTC)}
		_ = executeWithEngineAfterAcquire(context.Background(), root, io.Discard, io.Discard, engine, arm(engine), reuseFreshGreen)
		wantPrefix := []string{"lock-open", "lock-acquisition", "signal-arm", "owner-write"}
		if !reflect.DeepEqual(engine.trace[:len(wantPrefix)], wantPrefix) {
			t.Fatalf("gate acquisition prefix = %v, want %v", engine.trace, wantPrefix)
		}
	})

	for _, failOp := range []string{"lock-open", "lock-acquisition"} {
		t.Run("pre-acquire "+failOp, func(t *testing.T) {
			root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
			now := time.Date(2026, 7, 23, 12, 0, 0, 0, time.UTC)
			seed := executeWithEngine(context.Background(), root, io.Discard, io.Discard, &faultEngine{now: now})
			if !seed.Inspection.ReusableGreen {
				t.Fatalf("seed inspection = %+v, want reusable green", seed.Inspection)
			}
			engine := &faultEngine{now: now, failOp: failOp}
			got := executeWithEngineAfterAcquire(context.Background(), root, io.Discard, io.Discard, engine, arm(engine), forceRun)
			for _, operation := range engine.trace {
				if operation == "signal-arm" || operation == "owner-write" {
					t.Fatalf("pre-acquire %s reached %s: trace=%v", failOp, operation, engine.trace)
				}
			}
			if got.Inspection.State != Pending || got.Inspection.PendingStatus != "interrupted-pending" {
				t.Fatalf("pre-acquire %s inspection = %+v, want interrupted pending", failOp, got.Inspection)
			}
			ownerPath := filepath.Join(root, ".git", "bench-gate-owner")
			if _, err := os.Stat(ownerPath); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("pre-acquire %s wrote owner record: %v", failOp, err)
			}
		})
	}
}

func TestGateOwnerWritePrecedesPendingDurableReplace(t *testing.T) {
	t.Parallel()
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	engine := &faultEngine{now: time.Date(2026, 7, 23, 12, 0, 0, 0, time.UTC)}
	_ = executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
	ownerWrite, pendingCreate := -1, -1
	for i, operation := range engine.trace {
		if operation == "owner-write" {
			ownerWrite = i
		}
		if operation == "temporary-create" && pendingCreate < 0 {
			pendingCreate = i
		}
	}
	if ownerWrite < 0 || pendingCreate < 0 || ownerWrite > pendingCreate {
		t.Fatalf("owner write must precede pending durable replace: trace=%v", engine.trace)
	}
}
