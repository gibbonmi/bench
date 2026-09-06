package gate

import (
	"path/filepath"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// ComposedGreen reports whether the exact working tip's verdict and retained evidence cover the whole tree.
func ComposedGreen(root string) bool {
	return composedGreenAtKit(root, kitRoot(root))
}

func composedGreenAtKit(root, _ string) bool {
	now := time.Now()
	plan, err := buildSubject(root)
	if err != nil || !plan.Closed {
		return false
	}
	gitdir, err := benchgit.AdminDir(root)
	if err != nil {
		return false
	}
	loaded := loadVerdict(filepath.Join(gitdir, benchgit.GateCacheFile), now)
	if loaded.state != Ready || loaded.record.Status != "green" || loaded.record.Tree != plan.Tree || loaded.record.Oracle != plan.Oracle {
		return false
	}
	recorded, err := strictRecordTime(loaded.record.RecordedAt)
	if err != nil || now.Sub(recorded) >= freshness {
		return false
	}
	return !loaded.record.partitions() && !loaded.record.checkPartitions()
}
