package handoff

import (
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/handoffdoc"
)

// stateToken matches a backticked run of seven to forty hex characters — the shape a
// session writes a commit in. Seven is Git's own shortest unambiguous abbreviation and
// forty is a full sha, so a shorter or longer run names something else.
//
// The length rule alone is not a commit test. It admits any hex-shaped word, so every
// match still has to resolve to a commit object before the scan judges it.
var stateToken = regexp.MustCompile("`([0-9a-fA-F]{7,40})`")

// The two ways a State token fails, each named by what the reader has to fix. A pin the
// tip does not carry is stale and needs a newer sha; a pin that is not a commit was never
// a resume target at all and needs a different reference. One shared headline would leave
// the reader to work out which, and would let the commit peel be removed in silence.
// A third fault joins them: an abbreviation short enough to name two objects resolves to
// neither, so it sends a cold session nowhere. It is not prose, because the object store
// holds what it names twice over, and it is not a pin, because it names no one thing.
const (
	faultNotACommit  = "the handoff State names an object that is not a commit"
	faultOffAncestry = "the handoff State pins a commit outside this section's tip ancestry"
	faultAmbiguous   = "the handoff State abbreviates a sha that names more than one object"
)

// stateRepair is the one instruction every State refusal ends in. The scan refuses a line
// the writer owns, and the repair is the same whichever token fault it names.
const stateRepair = "rewrite or drop this line, then rerun bench handoff: "

// scanState refuses a State that pins a commit the section's tip does not contain. Such a
// pin is a stale resume target: a cold session reads it, checks out work that the phase
// has since rewritten or abandoned, and builds on a history nobody holds.
//
// A match that names nothing in the object store is left alone. State is prose, and a
// hex-shaped word in it — a ticket id, an English word — is not a claim about history.
//
// Fenced lines are skipped. A session pastes a command and its output into State, and the
// shas in that block describe a run that already happened rather than pinning this one.
//
// The scan reads scanRoot's object store, which every checkout of one repository shares,
// while the ancestry is measured against the owned section's own tip. An empty tip means
// none resolved, and the scan then has nothing to judge.
func scanState(scanRoot, tip, state string) error {
	if tip == "" {
		return nil
	}
	for line := range handoffdoc.UnfencedLines(state) {
		for _, match := range stateToken.FindAllStringSubmatch(line, -1) {
			fault := tokenFault(scanRoot, tip, match[1])
			if fault == "" {
				continue
			}
			return refusal{fault, stateRepair + strings.TrimSpace(line)}
		}
	}
	return nil
}

// tokenFault names why one token is not a usable resume target under tip, or answers ""
// when it is one — or when it names nothing the object store knows, which is prose.
//
// The two object tests answer two different questions, so both are asked. The `^{commit}`
// peel says whether the token is a commit this tip can be compared against. The bare form
// says whether it names an object at all, which is what separates a tree hash from prose.
// Each question has its own answer, so dropping the peel changes the printed reason rather
// than passing in silence.
//
// Both probes also fail on an abbreviation that matches two objects, which reads as prose
// under the exit code alone. The third question separates that case, and it is asked only
// once the first two have failed.
func tokenFault(scanRoot, tip, token string) string {
	if !git.OK("-C", scanRoot, "cat-file", "-e", token+"^{commit}") {
		if git.OK("-C", scanRoot, "cat-file", "-e", token) {
			return faultNotACommit
		}
		if ambiguous(scanRoot, token) {
			return faultAmbiguous
		}
		return ""
	}
	if !git.OK("-C", scanRoot, "merge-base", "--is-ancestor", token, tip) {
		return faultOffAncestry
	}
	return ""
}

// ambiguous reports whether the object store expands token to more than one object.
//
// `rev-parse --disambiguate` answers on stdout, which the existing git.Output already
// returns; the `cat-file` probes say the same thing on stderr, which no helper here
// returns. A prefix that matches nothing lists nothing, so the two-candidate test also
// separates ambiguity from prose. The token is lowered because Git takes an object prefix
// in lowercase alone, while the State token rule admits either case.
func ambiguous(scanRoot, token string) bool {
	out, err := git.Output("-C", scanRoot, "rev-parse", "--disambiguate="+strings.ToLower(token))
	if err != nil || out == "" {
		return false
	}
	return strings.Contains(out, "\n")
}
