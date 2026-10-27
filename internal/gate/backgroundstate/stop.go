package backgroundstate

import (
	"context"

	"github.com/gibbonmi/bench/internal/subprocess"
)

// Stop asks the owner of claim to end and returns only after that owner has exited, or when
// ctx ends. The signal goes through the lifetime owner in internal/subprocess. A claim that is
// not live has no owner to stop, so Stop sends nothing: a pid of an earlier boot can name
// another process.
func Stop(ctx context.Context, claim Claim) error {
	live, err := claim.Live()
	if err != nil || !live {
		return err
	}
	return subprocess.Stop(ctx, claim.Owner)
}
