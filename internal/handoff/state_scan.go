package handoff

import (
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
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
const (
	faultNotACommit  = "the handoff State names an object that is not a commit"
	faultOffAncestry = "the handoff State pins a commit outside this section's tip ancestry"
)

// scanState refuses a State that pins a commit the section's tip does not contain. Such a
// pin is a stale resume target: a cold session reads it, checks out work that the phase
// has since rewritten or abandoned, and builds on a history nobody holds.
//
// A match that names nothing in the object store is left alone. State is prose, and a
// hex-shaped word in it — a ticket id, an English word — is not a claim about history.
//
// The scan reads scanRoot's object store, which every checkout of one repository shares,
// while the ancestry is measured against the owned section's own tip. An empty tip means
// none resolved, and the scan then has nothing to judge.
func scanState(scanRoot, tip, state string) error {
	if tip == "" {
		return nil
	}
	for line := range strings.SplitSeq(state, "\n") {
		for _, match := range stateToken.FindAllStringSubmatch(line, -1) {
			fault := tokenFault(scanRoot, tip, match[1])
			if fault == "" {
				continue
			}
			return refusal{fault,
				"rewrite or drop this line, then rerun bench handoff: " + strings.TrimSpace(line)}
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
func tokenFault(scanRoot, tip, token string) string {
	if !git.OK("-C", scanRoot, "cat-file", "-e", token+"^{commit}") {
		if git.OK("-C", scanRoot, "cat-file", "-e", token) {
			return faultNotACommit
		}
		return ""
	}
	if !git.OK("-C", scanRoot, "merge-base", "--is-ancestor", token, tip) {
		return faultOffAncestry
	}
	return ""
}
