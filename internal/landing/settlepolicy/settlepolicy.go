// Package settlepolicy owns the settle decision of a conflicted reviewed landing.
// The parent package translates `merge-tree --write-tree -z` output into the stage
// records here once at its effect boundary, and it applies the verdict this package
// returns. The decision reads only the supplied records; the package reads no Git,
// performs no effects, reads no ambient process state, and starts no descendants. A
// union verdict's text merge stays in the parent, because that merge shells to
// `git merge-file`. The source census test enforces that boundary.
package settlepolicy

import "strings"

// StageRecord is one conflicted-file line of `merge-tree -z`: the mode and object of
// one stage (1 base, 2 destination, 3 source) at one path.
type StageRecord struct {
	Mode, OID, Path string
	Stage           int
}

// UnionStage is the CaptureSide stage of a path the rule table settles by union
// rather than by taking one side; no merge-tree stage carries that number.
const UnionStage = 0

// The three sides the capture rule table names.
const (
	SideSource      = "source"
	SideDestination = "destination"
	SideUnion       = "union"
)

// The four settle refusal reasons. The parent renders the reason with the refusing
// paths. The policy answers the first three; the parent answers the fourth after its
// own union text merge fails, because text-ness is not a fact the policy holds.
const (
	ReasonOutsideTable        = "path outside the capture table"
	ReasonNonRegularMode      = "non-regular mode"
	ReasonModeDisagreement    = "mode disagreement"
	ReasonUnionContentNotText = "union content not text"
)

// The two regular file modes. Every other mode leaves the conflict a refusal.
const (
	regularMode     = "100644"
	executableMode  = "100755"
	gitlinkMode     = "160000"
	symlinkFileMode = "120000"
)

// Refusal is a settle refusal: one reason and the refusing paths in merge-tree order.
// The zero value refuses nothing.
type Refusal struct {
	Reason string
	Paths  []string
}

// Refuse builds a refusal from a reason and its paths. The parent uses it to report
// the union content its own text merge could not read.
func Refuse(reason string, paths []string) Refusal {
	return Refusal{Reason: reason, Paths: paths}
}

// VerdictKind is how one path settles.
type VerdictKind int

const (
	// VerdictSide takes the winning stage's object at the path.
	VerdictSide VerdictKind = iota
	// VerdictUnion composes both sides as text; the parent performs that merge.
	VerdictUnion
	// VerdictRemove removes the path, because its winning stage is absent.
	VerdictRemove
)

// Verdict is the settle answer for one conflicted path. Record carries the winning
// stage of a VerdictSide. Stages carries every present stage of a VerdictUnion, which
// the parent feeds to its text merge.
type Verdict struct {
	Path   string
	Kind   VerdictKind
	Side   string
	Record StageRecord
	Stages map[int]StageRecord
}

// Settlement is the policy's answer over one conflict: either a verdict for every
// conflicted path in merge-tree order, or one refusal. A conflict holding no record
// carries neither, so the parent leaves it to the caller.
type Settlement struct {
	Verdicts []Verdict
	Refusal  Refusal
}

// CaptureSide is the rule table for a conflicted phase-owned path: it names the verb
// that settles the path and, for a take-a-side verb, the merge-tree stage that wins.
// The phase handoff is the source session's closing state, so the source wins it. The
// two append-only journals compose as the union of both sides, so no appended entry is
// lost. Every other capture file is the destination's running ledger, so the
// destination wins. A path outside the table has no rule, and the conflict stays a
// refusal.
func CaptureSide(path string) (stage int, side string, ok bool) {
	switch path {
	case "capture/session-handoff.md":
		return 3, SideSource, true
	case "capture/learnings.md", "capture/IDEAS.md":
		return UnionStage, SideUnion, true
	}
	if strings.HasPrefix(path, "capture/") {
		return 2, SideDestination, true
	}
	return 0, "", false
}

// engaged answers whether the capture rule table names any conflicted path. A conflict of
// code paths alone engages no rule, so the policy owns no answer over it.
func engaged(records []StageRecord) bool {
	for _, record := range records {
		if _, _, ok := CaptureSide(record.Path); ok {
			return true
		}
	}
	return false
}

// ConflictKind classifies a conflict from the modes its stage records carry. The two
// special modes are checked before the ordinary ones, so a gitlink or a symlink names
// the conflict even when the ordinary modes also disagree.
func ConflictKind(modes []string) string {
	for _, mode := range modes {
		if mode == gitlinkMode {
			return "gitlink"
		}
		if mode == symlinkFileMode {
			return "symlink"
		}
	}
	for _, mode := range modes {
		for _, other := range modes {
			if mode != other {
				return "mode"
			}
		}
	}
	return "textual"
}

// Settle decides how each conflicted path settles. The policy answers nothing at all
// until the table names at least one conflicted path, so a conflict that stays wholly
// outside capture keeps the caller's own bare refusal. Once engaged, it scans the records
// in merge-tree order first: a path the table does not name, or a stage whose mode is not
// regular, refuses at once, and the table check precedes the mode check inside one
// record. Only then does it settle each path in merge-tree order, so both record reasons
// beat a mode disagreement between the two sides.
func Settle(records []StageRecord) Settlement {
	if !engaged(records) {
		return Settlement{}
	}
	stages := map[string]map[int]StageRecord{}
	var order []string
	for _, record := range records {
		if _, _, ok := CaptureSide(record.Path); !ok {
			return Settlement{Refusal: Refuse(ReasonOutsideTable, []string{record.Path})}
		}
		if record.Mode != regularMode && record.Mode != executableMode {
			return Settlement{Refusal: Refuse(ReasonNonRegularMode, []string{record.Path})}
		}
		if _, seen := stages[record.Path]; !seen {
			stages[record.Path] = map[int]StageRecord{}
			order = append(order, record.Path)
		}
		stages[record.Path][record.Stage] = record
	}
	if len(order) == 0 {
		return Settlement{}
	}
	verdicts := make([]Verdict, 0, len(order))
	for _, path := range order {
		// Two sides that disagree on the file mode leave no settled mode to publish,
		// so the conflict stays a refusal under every verb.
		destination, hasDestination := stages[path][2]
		source, hasSource := stages[path][3]
		if hasDestination && hasSource && destination.Mode != source.Mode {
			return Settlement{Refusal: Refuse(ReasonModeDisagreement, []string{path})}
		}
		stage, side, _ := CaptureSide(path)
		if stage == UnionStage {
			verdicts = append(verdicts, Verdict{Path: path, Kind: VerdictUnion, Side: side, Stages: stages[path]})
			continue
		}
		record, ok := stages[path][stage]
		if !ok {
			verdicts = append(verdicts, Verdict{Path: path, Kind: VerdictRemove, Side: side})
			continue
		}
		verdicts = append(verdicts, Verdict{Path: path, Kind: VerdictSide, Side: side, Record: record})
	}
	return Settlement{Verdicts: verdicts}
}
