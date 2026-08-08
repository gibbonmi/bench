package gate

import (
	"path/filepath"
	"slices"
	"sort"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// ComposedGreen reports whether the exact working tip's verdict and retained evidence cover the whole tree.
func ComposedGreen(root string) bool {
	return composedGreenAtKit(root, kitRoot(root))
}

func composedGreenAtKit(root, kit string) bool {
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
	if loaded.record.partitions() || loaded.record.checkPartitions() {
		return partialComposedGreen(root, kit, plan, loaded.record, now)
	}
	return true
}

func partialComposedGreen(root, kit string, plan subject, record verdictRecord, now time.Time) bool {
	evaluation := newWorkingTreeEvaluationAtKit(root, kit)
	if _, err := evaluation.acceptPre(); err != nil {
		return false
	}
	scoping := evaluation.scope(plan.Resolution, forceRun, now)
	if !scoping.eligible {
		return false
	}
	if record.partitions() {
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
	}
	if record.checkPartitions() && !composedCheckPartition(root, record, scoping, now) {
		return false
	}
	return true
}

func composedCheckPartition(root string, record verdictRecord, scoping componentScoping, now time.Time) bool {
	if len(scoping.checks.Identities) == 0 {
		return false
	}
	slots, valid := loadConformanceCheckSlots(root)
	if !valid {
		return false
	}
	for _, name := range record.CheckInherited {
		check, found := registry.Find(name)
		if !found {
			return false
		}
		authoredAt, err := validateConformanceCheckSlot(slots[name], name, check.Tier, scoping.checks.Identities[name], now)
		if err != nil {
			return false
		}
		if (skipEvidence{Identity: scoping.checks.Identities[name], AuthoredAt: authoredAt.UTC().Format(time.RFC3339)}) != record.CheckEvidence[name] {
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
