package gate

// The landing subject of a background grade: a linked kit worktree whose tip holds no open
// build of a staged spec grades the tree that its landing would publish from that tip onto the
// destination. The grade and its lane pass read the primary checkout as the landing does, so
// the grade frames the destination schedule, and they run as the authenticated broker, so the
// grade key frames the broker that the landing runs.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/published"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/spec"
)

// landingSubject is the landing of a background grade: the primary checkout that the landing
// runs in, the destination commit and its green marker, and the deliverable of the landing.
type landingSubject struct {
	primary, destination, marker, deliverable string
}

// notBrokerError is a landing subject in a process that is not the broker that authenticates.
// The entry execs that broker.
type notBrokerError struct{ broker string }

func (e notBrokerError) Error() string {
	return fmt.Sprintf("the grade does not run as the broker %q", e.broker)
}

// landingSubjectOf answers the landing subject of source in the linked worktree root, and the
// tree that its landing publishes onto the destination at this call. The process must be the
// broker that authenticates, so a broker that does not authenticate refuses before the
// composition. A conflict that the composition does not settle refuses with a
// published.ConflictError.
func landingSubjectOf(root string, source chunkSource) (*landingSubject, string, error) {
	primary, err := benchgit.PrimaryCheckout(root)
	if err != nil || sameDirectory(primary, root) {
		return nil, "", errors.New("the tip holds no open build of a staged spec")
	}
	if broker, cause := authenticateBroker(root); cause != "" {
		return nil, "", errors.New("broker unavailable: " + cause)
	} else if broker != "" {
		return nil, "", notBrokerError{broker: broker}
	}
	branch, ok := benchgit.ResolvedDefault(primary)
	if !ok {
		return nil, "", errors.New("destination unavailable: default branch is unresolved")
	}
	destination, err := benchgit.ResolveCommit(primary, "refs/heads/"+branch)
	if err != nil {
		return nil, "", fmt.Errorf("destination unavailable: %w", err)
	}
	// A marker that does not read is an absent marker, which sends the landing plan to the
	// whole gate.
	marker, _, err := greenmarker.Read(primary, branch)
	if err != nil {
		marker = ""
	}
	landing := &landingSubject{primary: primary, destination: destination, marker: marker, deliverable: landingDeliverable(primary, source.tip, destination)}
	tree, err := published.Subject(primary, destination, source.tip, "", landing.deliverable, nil)
	if err != nil {
		return landing, "", err
	}
	return landing, tree, nil
}

// landingDeliverable answers the deliverable of the landing of tip onto destination: the one
// tickets-only folder that tip holds and destination does not, or no deliverable, a spec-less
// landing, when tip adds no such folder or more than one. A wrong answer costs only a lookup
// miss, because a grade records its green under the key that it ran.
func landingDeliverable(root, tip, destination string) string {
	listed, err := benchgit.Raw("-C", root, "ls-tree", "-z", "--name-only", tip, spec.SpecsDir+"/")
	if err != nil {
		return ""
	}
	var folders []string
	for _, path := range strings.Split(string(listed), "\x00") {
		name := strings.TrimPrefix(path, spec.SpecsDir+"/")
		if path == "" || !spec.TicketsOnly(spec.CommitTree(root, tip), name) || spec.CommitTree(root, destination).FolderIsDirectory(name) {
			continue
		}
		folders = append(folders, spec.ClosedFolderPath(name))
	}
	if len(folders) != 1 {
		return ""
	}
	return folders[0]
}

// withLandingPlan accepts the prospective plan of tree in a private checkout, with primary as
// its storage, identity, and baseline root, as a landing accepts it. It calls use with the
// checkout and the plan while the checkout exists.
func withLandingPlan(primary, tree string, use func(checkout string, plan subject)) error {
	artifacts, err := openProspectiveArtifacts(primary, tree)
	if err != nil {
		return err
	}
	defer artifacts.Close()
	evaluation := checkpointEvaluation(context.Background(), newProspectiveTreeEvaluation(artifacts.Checkout(), primary, primary, tree))
	plan, err := evaluation.acceptPre()
	if err != nil || plan.Tree != tree {
		return errors.New("subject unavailable")
	}
	use(artifacts.Checkout(), plan)
	return nil
}

// inspect answers the evidence of the landing grade of tree, and whether the landing plan of
// tree against the destination selects the reader set. The plan reads the green marker, the
// reader hold, and the reader chain as a landing reads them. The selection names only the
// source root that the run would build, which is all that the plan reads of it.
func (l *landingSubject) inspect(tree string) (EvidenceInspection, bool, error) {
	var evidence EvidenceInspection
	reader := false
	err := withLandingPlan(l.primary, tree, func(checkout string, plan subject) {
		evidence = inspectEvidence(l.primary, plan, time.Now())
		var selection *runbinary.Selection
		if source, selected := runBinarySelection(true, plan.Resolution, checkout, l.primary); selected {
			selection = &runbinary.Selection{SourceRoot: source}
		}
		_, refusal := readerPlanner{shellcheck: shellcheckFiles}.decide(landingValue{destination: l.destination, marker: l.marker}, true, false, checkout, l.primary, plan, selection)
		reader = refusal == ""
	})
	return evidence, reader, err
}
