// Stack flow for the worktree landing command: one lease, one grade, and one publication
// for the reviewed sources of a stack file.
package worktree

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var stackLandGrammar = usage.Grammar{
	Cmd:  usage.CommandName(usage.WorktreeLandStack),
	Help: "usage: " + usage.WorktreeLandStack,
	Flags: []usage.Flag{
		{Name: "--stack", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true, Repeatable: true},
	},
}

// landDispatchGrammar declares every flag of each land form. The dispatch skips the value
// of each declared flag, so a flag value that reads `--stack` or `--resume` never selects a
// form.
var landDispatchGrammar = usage.Grammar{Flags: slices.Concat(landGrammar.Flags, resumeLandGrammar.Flags, stackLandGrammar.Flags)}

func hasStackFlag(args []string) bool {
	return usage.FlagPresent(landDispatchGrammar, args, "--stack")
}

func hasResumeFlag(args []string) bool {
	return usage.FlagPresent(landDispatchGrammar, args, "--resume")
}

// stackMemberLanding is one stack member after its proofs: its 0-based index in the stack
// file, the read member, its worktree path, its assignment, and its source fact.
type stackMemberLanding struct {
	index      int
	member     StackMember
	path       string
	assignment intent.Assignment
	source     landingSourceFact
}

// landStack is the stack form of the land leaf. It runs inside the landing span of
// landWith, and the span names the assignment of the first member.
func landStack(ctx context.Context, measures *landingMeasures, j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int {
	if hasResumeFlag(args) {
		return resumeStack(j, a, root, args, stdout, stderr)
	}
	parsed, line, code := usage.Parse(stackLandGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	route, unresolved, err := readStackRoute(a, root, parsed)
	if err != nil {
		return stackUsageRefusal(stderr, stackLandGrammar, err)
	}
	destination, branch, priorMarker, destinationFingerprint, landings, refusals := stackPreflight(j, a.run, root, route, unresolved)
	if len(refusals) > 0 {
		for _, err := range refusals {
			landRefusalError(stdout, err)
		}
		return 1
	}
	measures.assignment = landings[0].assignment.ID
	request := landing.StackRequest{Root: root, Destination: "refs/heads/" + branch, DestinationBase: destination, DestinationFingerprint: destinationFingerprint, GreenMarker: priorMarker, Stdout: stdout, Stderr: stderr}
	publications := make([]commitrepo.Publication, len(landings))
	records := make([]int, len(landings))
	brokerChanged := false
	for i, l := range landings {
		// The counts are read before the release step, because that step drops the records.
		records[i] = censusCount(a.home, root, l.assignment.ID)
		brokerChanged = brokerChanged || brokerSourceChanged(l.assignment.Worktree, l.source.base, l.source.tip)
		printLandingSource(stderr, l.source.base, l.assignment.Start)
		var member landing.ReviewedRequest
		member, publications[i] = sourceRequest(l.assignment, l.source, l.member.Request, l.member.Message, request.Destination)
		request.Members = append(request.Members, member)
	}
	// The window decides before the lease, so a stack never waits for a lease and then for a gap.
	if refusals := stackWindow(a.run, root, route, landings, publications); len(refusals) > 0 {
		for _, err := range refusals {
			landRefusalError(stdout, err)
		}
		return 1
	}
	lease, err := commitrepo.Store{Root: root, Run: a.run}.AcquireLease(request.Destination, landings[0].assignment.ID)
	if err != nil {
		return landRefusalError(stdout, landingFaceRoute(err, route.rerun(stackUnattributed), false))
	}
	publications[0].Lease = lease
	// The stack owns the machine while it grades, so no member grade runs beside it.
	for _, l := range landings {
		stopGradesOf(a.home, l.path, stderr)
	}
	result, err := j.landStack(ctx, request, stackAdmission{store: commitrepo.Store{Root: root, Run: a.run}, members: publications, gap: j.publicationGap})
	if err != nil {
		lease.Release()
	}
	red := stackAttributed(err, result)
	printStackMembers(stdout, landings, result, red, records)
	end := stackEnd{destinationBase: destination, published: "none", worktree: "active"}
	if k := result.Published(); k > 0 {
		top := result.Merges[k-1]
		measures.subject = top
		end = stackEnd{destinationBase: result.DestinationBase, published: top, landed: k}
		end = end.after(stackAfterPublication(j, a, root, branch, landing.ReviewedResult{DestinationBase: result.DestinationBase, Commit: top, Kind: result.Kind}, priorMarker, landings[:k], brokerChanged, stdout, stderr))
	}
	if err != nil {
		for _, refused := range stackLandingRefusal(err, root, route, landings, destination, branch, red, result.Published()) {
			landRefusalError(stdout, refused)
		}
		// A stopped step after the publication keeps its exit, because its resume is the
		// next action whatever the refusal beside it.
		if end.step == "" {
			end.exit = 1
		}
	}
	return printStackEnd(stdout, route, end)
}

// readStackRoute reads the stack file that the `--stack` flag of parsed names, with the
// overrides of the call, and resolves each member path the way a single landing resolves
// its operand, so two spellings of one worktree meet as one canonical path. A member whose
// path does not resolve keeps its operand and reads as unresolved. Two members that resolve
// to one path refuse.
func readStackRoute(a ambient, root string, parsed usage.Result) (route stackRoute, unresolved []bool, err error) {
	file, err := callPath(a, parsed.Flags["--stack"])
	if err != nil {
		return stackRoute{}, nil, err
	}
	input, err := ReadStack(file, a.dir, parsed.Repeated["--source-tip"])
	if err != nil {
		return stackRoute{}, nil, err
	}
	route = stackRoute{file: file, input: input, paths: make([]string, len(input.Members)), assignments: make([]string, len(input.Members))}
	unresolved = make([]bool, len(input.Members))
	named := map[string]int{}
	for i, m := range input.Members {
		operand := resolveVerbOperand(a.run, root, m.Path)
		path, err := canonicalPathIn(a, operand)
		if err != nil {
			route.paths[i], unresolved[i] = operand, true
			continue
		}
		if first, ok := named[path]; ok {
			return stackRoute{}, nil, fmt.Errorf("stack members %d and %d name one worktree %s", first+1, i+1, path)
		}
		named[path], route.paths[i] = i, path
	}
	return route, unresolved, nil
}

// sourceRequest maps a resolved assignment and its source fact onto the source fields of
// one reviewed request, with message as its merge message. It also answers the publication
// that admits that source at destination for the request requestArg.
func sourceRequest(assignment intent.Assignment, source landingSourceFact, requestArg, message, destination string) (landing.ReviewedRequest, commitrepo.Publication) {
	r := landing.ReviewedRequest{
		Source: assignment.Branch, SourceTip: source.tip, ReviewBase: source.base,
		SourceWorktree: assignment.Worktree, SourceFingerprint: source.fingerprint,
		SpecPath: source.specPath, SpecBytes: source.specBytes, SpecMode: source.specMode, ClosePath: source.closePath,
		Message: message,
	}
	return r, commitrepo.Publication{Assignment: assignment.ID, Request: intent.RequestDigest(requestArg), Worktree: assignment.Worktree, Source: source.tip, Deliverable: r.Deliverable(), Destination: destination}
}

// stackUsageRefusal prints a stack input refusal beside the help of the stack form g and
// answers the usage exit. It runs before any ledger read, so the refusal takes no lease.
func stackUsageRefusal(stderr io.Writer, g usage.Grammar, err error) int {
	fmt.Fprintln(stderr, g.Help+" ("+sanitize.Controls(err.Error())+")")
	return 2
}

// stackPreflight runs the destination proofs once, and for each member over its effective
// tip the assignment and source proofs of one landing and then the stack proofs. The owner
// rule then runs over the members that cleared their proofs. It answers every refusal of
// these groups, and each member refusal names its member index and label.
func stackPreflight(j joins, run *git.Run, root string, route stackRoute, unresolved []bool) (destination, branch, priorMarker, fingerprint string, landings []stackMemberLanding, refusals []error) {
	destination, branch, priorMarker, fingerprint, err := landingDestination(root)
	if err != nil {
		refusals = append(refusals, landingFaceRoute(err, route.rerun(stackUnattributed), false))
	}
	for i, m := range route.input.Members {
		if unresolved[i] {
			refusals = append(refusals, stackMemberRefusal(i, "", landingFaceRoute(errors.New("worktree path is not canonical"), route.rerun(stackRerun{refused: i}), false)))
			continue
		}
		base, tip := expandIdentity(root, m.Base), expandIdentity(root, m.EffectiveTip)
		assignment, err := landingAssignment(j, run, root, route.paths[i], m.Request, base, tip)
		if err != nil {
			// A later assignment proof can refuse a request that resolved, and the refusal
			// then names the label of that assignment.
			label := ""
			if resolved, found, _ := intent.In(run, root).FindAssignmentForRequest(m.Request); found {
				label = resolved.Label
			}
			refusals = append(refusals, stackMemberRefusal(i, label, landingFaceRoute(err, route.rerun(stackRefusedBy(i, err)), true)))
			continue
		}
		route.assignments[i] = assignment.ID
		source, err := landingSourceProofs(j, run, root, assignment, base, tip, m.Spec, destination)
		if err != nil {
			refusals = append(refusals, stackMemberRefusal(i, assignment.Label, landingFaceRoute(err, route.rerun(stackRefusedBy(i, err)), false)))
			continue
		}
		l := stackMemberLanding{index: i, member: m, path: route.paths[i], assignment: assignment, source: source}
		if err := stackMemberProofs(run, root, destination, branch, l); err != nil {
			refusals = append(refusals, stackRaisedRefusal(route, i, assignment.Label, err, stackRefusedBy(i, err)))
			continue
		}
		landings = append(landings, l)
	}
	refusals = append(refusals, stackOverlaps(root, destination, route, landings)...)
	return destination, branch, priorMarker, fingerprint, landings, refusals
}

// stackMemberRefusal names the member index and label of a member's refusal ahead of its
// sentence. A member whose assignment did not resolve has no label to name.
func stackMemberRefusal(i int, label string, err error) error {
	prefix := "stack member " + strconv.Itoa(i+1)
	if label != "" {
		prefix += " (" + label + ")"
	}
	var typed refusalError
	if errors.As(err, &typed) {
		typed.detail = prefix + causeJoiner + typed.detail
		return typed
	}
	return fmt.Errorf("%s%s%w", prefix, causeJoiner, err)
}

// stackAttributed is the 0-based member that the bisect attributed the red refusal err to,
// or -1. A recheck or an admission refusal after the bisect leaves the attribution set, so
// only a candidate red reads it.
func stackAttributed(err error, result landing.StackResult) int {
	var refused landing.AuthorizationRefusal
	if errors.As(err, &refused) && refused.Result.Kind == authorization.Candidate && result.Attributed >= 1 {
		return result.Attributed - 1
	}
	return -1
}

// stackLandingRefusal routes a refusal of the stack owner. A member refusal keeps the face
// of one landing. The member red, which the bisect attributed, takes the stack-red face: its
// route removes each of the landed members from the stack file and sets the tip of each
// later member by its fold. A destination that moved after the grade takes the
// moved-destination face of one landing for each member that lacks the moved tip, whose fold
// merges from branch. A stop that no member caused hands back, and so does a moved destination
// that no member lacks. An infrastructure grade keeps its own face.
func stackLandingRefusal(err error, root string, route stackRoute, landings []stackMemberLanding, destination, branch string, red, landed int) []error {
	if red >= 0 {
		raised := refusalError{refusal{detail: err.Error(), face: faceStackRed}}
		fact := stackRefusedBy(red, raised)
		for i := range landed {
			fact.removed = append(fact.removed, i)
		}
		for i := landed; landed > 0 && i < len(landings); i++ {
			fact.reset = append(fact.reset, i)
		}
		return []error{stackRaisedRefusal(route, red, landings[red].assignment.Label, raised, fact)}
	}
	if errors.As(err, new(landing.TipMovedError)) {
		if refusals := stackMovedRefusals(root, branch, route, landings); len(refusals) > 0 {
			return refusals
		}
	}
	var member landing.StackMemberError
	if errors.As(err, &member) && member.Member >= 1 && member.Member <= len(landings) {
		i := member.Member - 1
		l := landings[i]
		var conflict landing.ConflictError
		if errors.As(member.Err, &conflict) {
			// The conflict route repairs the source, as the single landing's re-run states.
			rerun := route.rerun(stackRerun{refused: i, repairs: true})
			return []error{stackMemberRefusal(i, l.assignment.Label, landingConflictRefusalAt(conflict, destination, l.assignment.ID, l.path, l.assignment.Worktree, rerun))}
		}
		return []error{stackMemberRefusal(i, l.assignment.Label, landingFaceRoute(member.Err, route.rerun(stackRefusedBy(i, member.Err)), false))}
	}
	face := authorizationFace(err)
	if face != faceLandInfrastructure {
		face = faceLandHandback
	}
	return []error{landingFaceRoute(refusalError{refusal{detail: err.Error(), face: face}}, route.rerun(stackUnattributed), false)}
}

// stackMovedRefusals reads the destination tip that moved after the grade and runs the
// behind proof of each member over it, so one run names every member that folds. A tip that
// does not read answers no refusal.
func stackMovedRefusals(root, branch string, route stackRoute, landings []stackMemberLanding) []error {
	moved, _, _, err := landingDestinationIdentity(root)
	if err != nil {
		return nil
	}
	var refusals []error
	for _, l := range landings {
		if err := stackBehind(root, moved, branch, l); err != nil {
			refusals = append(refusals, stackRaisedRefusal(route, l.index, l.assignment.Label, err, stackRefusedBy(l.index, err)))
		}
	}
	return refusals
}

// stackAfterPublication runs the steps after a stack publication once for the published
// merge. It then releases each published member in file order and runs the two landing
// effects once. It answers the name of the step that stopped, or the empty name.
func stackAfterPublication(j joins, a ambient, root, branch string, result landing.ReviewedResult, priorMarker string, landings []stackMemberLanding, brokerChanged bool, stdout, stderr io.Writer) string {
	if step := publishedSteps(j, a.run, root, branch, result, priorMarker); step != "" {
		return step
	}
	for _, l := range landings {
		if !releaseLanded(j, a, root, l.member.Request, l.path, stderr) {
			return "release"
		}
	}
	return runLandedEffects(j, a, root, result.DestinationBase, brokerChanged, stdout, stderr)
}

// stackMemberRow is one row of the members table.
type stackMemberRow struct {
	member                               int
	label, sourceTip, published, outcome string
	census                               int
}

// membersTable names the members table block of a stack.
const membersTable = "members"

// The outcome cells of the members table.
const (
	stackOutcomeLanded   = "landed"
	stackOutcomeRed      = "red"
	stackOutcomeUnlanded = "unlanded"
)

// stackMembersTable renders the members table of a stack, one row for each member in file
// order.
func stackMembersTable(rows []stackMemberRow) (string, error) {
	cells := make([][]string, len(rows))
	for i, row := range rows {
		cells[i] = []string{strconv.Itoa(row.member), row.label, row.sourceTip, row.published, row.outcome, strconv.Itoa(row.census)}
	}
	return toon.Table(membersTable, []string{"member", "label", "source_tip", "published_commit", "outcome", "census"}, cells)
}

// printStackMembers prints the members table of result. A member that published reads
// landed with its merge, the 0-based member red reads red, and every other member reads
// unlanded. A table that does not render prints the render error in its place, and the
// stack goes on to its record.
func printStackMembers(stdout io.Writer, landings []stackMemberLanding, result landing.StackResult, red int, records []int) {
	rows := make([]stackMemberRow, len(landings))
	for i, l := range landings {
		rows[i] = stackMemberRow{member: i + 1, label: l.assignment.Label, sourceTip: l.source.tip, published: "none", outcome: stackOutcomeUnlanded, census: records[i]}
		switch {
		case i < result.Published():
			rows[i].published, rows[i].outcome = result.Merges[i], stackOutcomeLanded
		case i == red:
			rows[i].outcome = stackOutcomeRed
		}
	}
	table, err := stackMembersTable(rows)
	if err != nil {
		table = toon.RenderError(err) + "\n"
	}
	fmt.Fprint(stdout, table)
}

// stackAdmission is the installed broker's publication decision for the members of one
// stack. It decides each member as the single landing's admission decides its source.
type stackAdmission struct {
	store   commitrepo.Store
	members []commitrepo.Publication
	gap     func(string)
}

func (a stackAdmission) CheckMember(i int, tree string) error {
	return commitmentAdmission{store: a.store, source: a.members[i]}.Check(tree)
}

func (a stackAdmission) PublishStack(k int, trees []string, publish func() error) error {
	return publishAdmitted(a.store.Root, a.gap, publish, func(gapped func() error) error {
		return a.store.PublishStackAdmitted(a.members[:k], trees, gapped)
	})
}
