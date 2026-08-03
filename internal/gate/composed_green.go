package gate

import (
	"path/filepath"
	"slices"
	"sort"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// ComposedGreen reports whether the exact working tip's verdict and retained evidence cover the whole tree.
func ComposedGreen(root string) bool {
	now := time.Now()
	plan, err := buildSubject(root)
	if err != nil || !plan.Closed {
		return false
	}
	gitdir, err := commonGitDir(root)
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
	if loaded.record.partitions() {
		return partialComposedGreen(root, plan, loaded.record, now)
	}
	if loaded.record.inherits() {
		return reducedComposedGreen(root, plan, loaded.record, now)
	}
	return true
}

func reducedComposedGreen(root string, plan subject, record verdictRecord, now time.Time) bool {
	inherited := reducedInheritance(root, root, plan.Resolution, now)
	if !inherited.ok || record.Ancestor != inherited.ancestor || !slices.Equal(record.Phases, composedPhaseNames(inherited.phases)) {
		return false
	}
	recordedAt, err := strictRecordTime(record.AncestorRecordedAt)
	return err == nil && recordedAt.Equal(inherited.ancestorAt)
}

func partialComposedGreen(root string, plan subject, record verdictRecord, now time.Time) bool {
	scoping := scopeComponents(root, plan.Resolution, forceRun, now)
	if !scoping.eligible {
		return false
	}
	table, err := phaseTable(root, scoping.runnerRoot)
	if err != nil {
		return false
	}
	claimed := append([]string(nil), record.Executed...)
	claimed = append(claimed, record.Skipped...)
	sort.Strings(claimed)
	if !slices.Equal(claimed, composedPhaseNames(table)) {
		return false
	}
	for _, component := range record.Skipped {
		identity, ok := scoping.identities[component]
		if !ok {
			return false
		}
		skip, ok := componentSkip(root, scoping, component, identity, now)
		if !ok || skip.evidence() != record.SkipEvidence[component] {
			return false
		}
	}
	return true
}

func composedPhaseNames(phases []Phase) []string {
	names := make([]string, 0, len(phases))
	for _, phase := range phases {
		names = append(names, phase.Name)
	}
	sort.Strings(names)
	return names
}
