// Package admissionpolicy owns the pure admission rules of the intent ledger's
// mutators. Each rule answers the next ledger and the changed report from the
// supplied ledger and the mutator's operand alone: it resolves no address, takes
// no lock, reads no file, and runs no Git. The intent package keeps every effect,
// translates the process state a rule needs into typed facts at its boundary, and
// persists the answer through the ledger transaction. The source census test
// enforces that boundary.
package admissionpolicy

import (
	"sort"

	"github.com/gibbonmi/bench/internal/intent/ledger"
)

// LivenessFacts carry the three proofs the compaction rule reads. The intent
// package translates `git.Worktrees` and `git.LocalBranches` into Candidates,
// `os.Stat` into WorktreeExists, and `git.ResolvedDefault` with
// `git.LandedInDefault` into Landed. Each map answers only for the keys the
// boundary resolved, and an absent key is the conservative answer: an entry whose
// worktree the boundary could not prove present is dropped, and a branch the
// boundary could not prove landed is kept.
type LivenessFacts struct {
	// Candidates reports whether the repository still holds an agent worktree or an
	// agent branch, which is the only proof an uncorrelated Claude entry has.
	Candidates bool
	// WorktreeExists answers, per recorded worktree path, whether that path is present.
	WorktreeExists map[string]bool
	// Landed answers, per recorded branch name, whether that branch landed in the
	// resolved default branch.
	Landed map[string]bool
}

// Upsert inserts or enriches one stable writer key. An entry identical to the
// stored one under the stored creation stamp reports no change, so the caller
// preserves the ledger's exact bytes.
func Upsert(current ledger.Ledger, entry ledger.Entry) (ledger.Ledger, bool, error) {
	for i := range current.Entries {
		if current.Entries[i].Key != entry.Key {
			continue
		}
		stamped := entry
		stamped.CreatedAt = current.Entries[i].CreatedAt
		if current.Entries[i] == stamped {
			return current, false, nil
		}
		current.Entries[i] = stamped
		return current, true, nil
	}
	current.Entries = append(current.Entries, entry)
	return current, true, nil
}

// Live returns the proof-live entries of the supplied ledger, oldest first and
// keyed on a tie. It is the one liveness derivation: the snapshot renders it, and
// the compaction rule keeps exactly what it returns.
func Live(current ledger.Ledger, facts LivenessFacts) []ledger.Entry {
	live := make([]ledger.Entry, 0, len(current.Entries))
	for _, entry := range current.Entries {
		if entry.Kind == ledger.KindClaudeAgent && entry.Worktree == "" && entry.Branch == "" {
			if facts.Candidates {
				live = append(live, entry)
			}
			continue
		}
		if entry.Worktree != "" && !facts.WorktreeExists[entry.Worktree] {
			continue
		}
		if entry.Branch != "" && facts.Landed[entry.Branch] {
			continue
		}
		live = append(live, entry)
	}
	sort.Slice(live, func(i, j int) bool {
		if live[i].CreatedAt.Equal(live[j].CreatedAt) {
			return live[i].Key < live[j].Key
		}
		return live[i].CreatedAt.Before(live[j].CreatedAt)
	})
	return live
}

// Compact removes only the entries Live does not prove live, and it keeps the
// stored order of the survivors. A ledger that loses no entry reports no change.
func Compact(current ledger.Ledger, facts LivenessFacts) (ledger.Ledger, bool, error) {
	liveKeys := map[string]bool{}
	for _, entry := range Live(current, facts) {
		liveKeys[entry.Key] = true
	}
	next := make([]ledger.Entry, 0, len(current.Entries))
	for _, entry := range current.Entries {
		if liveKeys[entry.Key] {
			next = append(next, entry)
		}
	}
	if len(next) == len(current.Entries) {
		return current, false, nil
	}
	current.Entries = next
	return current, true, nil
}
