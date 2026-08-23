// Package landingpolicy owns the pure landing decisions of the worktree
// commands. The parent package translates Git, filesystem, and process state
// into these typed facts once at its effect boundary. The decisions here read
// only the supplied facts and return a verdict the parent renders; the package
// performs no effects, reads no ambient process state, and starts no
// descendants. The source census test enforces that boundary.
package landingpolicy

// StatusEntry is one porcelain status record of the landing destination: the
// two XY status characters and the path, translated at the parent boundary.
type StatusEntry struct {
	Status string
	Path   string
}

// ResidueFacts describe the landing destination's working state before a
// destructive step. The func fields are lazily supplied facts: the parent
// binds them at the boundary, and a decision consults each at most once.
type ResidueFacts struct {
	// NestedClean is false when the destination holds nested repositories.
	NestedClean bool
	// StatusReadable is false when the destination status cannot be read.
	StatusReadable bool
	// StatusWellFormed is false when the destination status does not parse.
	StatusWellFormed bool
	// Entries are the parsed status records of the destination.
	Entries []StatusEntry
	// IgnoredDeclared reports whether every ignored path sits inside the
	// destination's declared build-output allowance.
	IgnoredDeclared func() bool
	// DestinationAtPublished is true when the destination commit is the
	// published landing commit.
	DestinationAtPublished bool
	// StagedMatchesPublished reports whether the staged index carries exactly
	// the published landing content.
	StagedMatchesPublished func() bool
}

// Residue returns the refusal detail that blocks a destructive destination
// step, or the empty string when the destination state is safe.
func Residue(f ResidueFacts) string {
	if !f.NestedClean {
		return "landing destination has nested repositories"
	}
	if !f.StatusReadable {
		return "landing destination status is unreadable"
	}
	if !f.StatusWellFormed {
		return "landing destination status is malformed"
	}
	staged := false
	allowedIgnored, allowanceKnown := false, false
	for _, entry := range f.Entries {
		switch entry.Status {
		case "":
			continue
		case "!!":
			if !allowanceKnown {
				allowedIgnored, allowanceKnown = f.IgnoredDeclared(), true
			}
			if allowedIgnored {
				continue
			}
			return "landing destination has ignored residue"
		case "??":
			return "landing destination has untracked collisions"
		}
		if entry.Status[1] != ' ' {
			return "landing destination has tracked-worktree changes"
		}
		staged = staged || entry.Status[0] != ' '
	}
	if staged && (!f.DestinationAtPublished || !f.StagedMatchesPublished()) {
		return "landing destination has staged changes"
	}
	return ""
}

// MarkerAction is the resume decision over the project-green marker.
type MarkerAction int

const (
	// MarkerAdvance means the destination is at the published landing, so the
	// resume advances the marker to it.
	MarkerAdvance MarkerAction = iota
	// MarkerAccept means the destination moved on and the existing marker
	// already covers the published landing; the resume leaves it alone.
	MarkerAccept
	// MarkerRefuse means no marker covers the published landing.
	MarkerRefuse
)

// MarkerRefusalDetail is the refusal detail a MarkerRefuse verdict renders.
const MarkerRefusalDetail = "project-green marker is absent, behind, or divergent from the published landing"

// MarkerFacts describe the resume marker state.
type MarkerFacts struct {
	// DestinationAtPublished is true when the destination commit is the
	// published landing commit.
	DestinationAtPublished bool
	// MarkerPresent is true when a project-green marker exists.
	MarkerPresent bool
	// MarkerReachesPublished reports whether the published landing is an
	// ancestor of the marker. Lazily supplied; consulted at most once.
	MarkerReachesPublished func() bool
}

// ResumeMarker decides what a resumed landing does with the green marker.
func ResumeMarker(f MarkerFacts) MarkerAction {
	if f.DestinationAtPublished {
		return MarkerAdvance
	}
	if !f.MarkerPresent || !f.MarkerReachesPublished() {
		return MarkerRefuse
	}
	return MarkerAccept
}

// Refusal is a landing refusal verdict. The zero value accepts. A refusal
// with Observed or Wanted set renders as an identity mismatch.
type Refusal struct {
	Detail   string
	Observed string
	Wanted   string
}

// SpecTransition describes what the published commit carries for the landing's
// named spec, translated at the boundary.
type SpecTransition struct {
	// Named is true when the landing names a spec.
	Named bool
	// TicketsOnlyClose is true when the source carries the spec folder with no
	// spec.md, so the landing is a folder close, not a staged transition.
	TicketsOnlyClose bool
	// PublishedHasFolder is true when the published tree still carries the
	// tickets-only folder.
	PublishedHasFolder bool
	// TransitionMatches is true when the source staged spec derives an
	// implemented body identical to the published spec.
	TransitionMatches bool
}

// PublicationFacts describe a claimed published landing commit.
type PublicationFacts struct {
	// Requested is the caller's published-commit value.
	Requested string
	// Resolved is true when the value resolves to a commit; Published is then
	// that commit's full identity.
	Resolved  bool
	Published string
	// ReachableFromDestination is true when the published commit is an
	// ancestor of the destination.
	ReachableFromDestination bool
	// ParentsReadable is true when the published commit's parent record is
	// readable; Parents is then that record's fields (commit then parents).
	ParentsReadable bool
	Parents         []string
	// RequestedSource is the reviewed source tip the resume claims.
	RequestedSource string
	// RangeAuthenticates is true when the review base resolves a source range
	// for the claimed source.
	RangeAuthenticates bool
	// Spec is the published spec-transition state.
	Spec SpecTransition
}

// Publication authenticates a claimed published landing commit. The zero
// Refusal accepts the publication.
func Publication(f PublicationFacts) Refusal {
	if !f.Resolved {
		return Refusal{Detail: "published commit is not an exact commit identity"}
	}
	if f.Published != f.Requested {
		return Refusal{Detail: "published commit is not an exact commit identity", Observed: f.Requested, Wanted: f.Published}
	}
	if !f.ReachableFromDestination {
		return Refusal{Detail: "published commit is not reachable from the destination"}
	}
	if !f.ParentsReadable {
		return Refusal{Detail: "published commit parents are unreadable"}
	}
	if len(f.Parents) != 3 {
		return Refusal{Detail: "published commit does not authenticate the reviewed source parent"}
	}
	if f.Parents[2] != f.RequestedSource {
		return Refusal{Detail: "published commit does not authenticate the reviewed source parent", Observed: f.RequestedSource, Wanted: f.Parents[2]}
	}
	if !f.RangeAuthenticates {
		return Refusal{Detail: "review base does not authenticate the published source"}
	}
	if f.Spec.Named {
		if f.Spec.TicketsOnlyClose {
			if f.Spec.PublishedHasFolder {
				return Refusal{Detail: "published commit does not close the source tickets-only folder"}
			}
		} else if !f.Spec.TransitionMatches {
			return Refusal{Detail: "published commit does not carry the source staged spec transition"}
		}
	}
	return Refusal{}
}

// TerminalFacts describe how a published landing's follow-up work ended.
type TerminalFacts struct {
	// FailedStep is the follow-up step that failed: "marker", "reconcile", or
	// "release". Empty when every follow-up step completed.
	FailedStep string
	// Active is true when the landing released its assignment in this run,
	// and false when a prior run already completed the release.
	Active bool
}

// TerminalOutcome is the landing's terminal record: the worktree state token
// the landed record renders and the process exit code.
type TerminalOutcome struct {
	ExitCode      int
	WorktreeState string
}

// Terminal decides the landing's terminal outcome.
func Terminal(f TerminalFacts) TerminalOutcome {
	if f.FailedStep != "" {
		return TerminalOutcome{ExitCode: 3, WorktreeState: "incomplete:" + f.FailedStep}
	}
	if !f.Active {
		return TerminalOutcome{ExitCode: 0, WorktreeState: "already-complete"}
	}
	return TerminalOutcome{ExitCode: 0, WorktreeState: "released"}
}
