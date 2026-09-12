package worktree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

const landedSetFingerprintVersion = "bench-landed-set/v1"

// planLandedExplicitWithOptions exposes the target-path boundary so hostile-path tests
// can prove shape rejection happens before the explicit planner invokes Git.
var planLandedExplicitWithOptions = PlanExplicitWithOptions

type landedCleanupRow struct {
	assignment intent.Assignment
	plan       CleanupPlan
	headOID    string
	lease      string
}

type landedCleanupSet struct {
	rows        []landedCleanupRow
	fingerprint string
}

// planLandedSet plans every assignment the landed selector admits under scope. An empty
// scope is the whole repository, which is what `bench worktree clean --landed` asks for.
// A non-empty scope is one landing's destination base, which narrows the set to the
// assignments whose work that landing carries.
func planLandedSet(j joins, root string, options CleanupOptions, scope string) (landedCleanupSet, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return landedCleanupSet{}, err
	}
	sort.Slice(assignments, func(i, j int) bool { return assignments[i].ID < assignments[j].ID })
	defaultRef, ok := git.ResolvedDefault(root)
	if !ok {
		return landedCleanupSet{}, nil
	}
	leases, err := assignmentLeaseStates(root)
	if err != nil {
		return landedCleanupSet{}, err
	}

	set := landedCleanupSet{rows: make([]landedCleanupRow, 0, len(assignments))}
	for _, assignment := range assignments {
		lease := leases[assignment.OwnerID]
		if row, selected := selectLandedCleanupRow(j, root, assignment, defaultRef, lease, options, scope); selected {
			set.rows = append(set.rows, row)
		}
	}
	if len(set.rows) == 0 {
		return set, nil
	}
	set.fingerprint = fingerprintLandedSet(set.rows, options)
	for i := range set.rows {
		set.rows[i].plan.Fingerprint = set.fingerprint
	}
	return set, nil
}

// selectLandedCleanupRow is the selector's single per-assignment proof. The set plan and
// every pre-mutation row re-plan use it. A row cannot become removable through a
// different route after its shared fingerprint was validated.
//
// scope narrows the selector to one landing's own work. The same landed proof answers
// both questions: the branch has to be landed in the destination as it stands now, and it
// has to be unlanded at the destination base this landing composed against. An assignment
// that was already landed before this landing belongs to no landing's own set.
func selectLandedCleanupRow(j joins, root string, assignment intent.Assignment, defaultRef, lease string, options CleanupOptions, scope string) (landedCleanupRow, bool) {
	if assignment.State != intent.StateActive || assignment.Branch == "" {
		return landedCleanupRow{}, false
	}
	branch := strings.TrimPrefix(assignment.Branch, "refs/heads/")
	landed, byContent, proofErr := git.LandedInDefault(root, branch, defaultRef)
	if proofErr != nil || !landed {
		return landedCleanupRow{}, false
	}
	if scope != "" {
		before, _, scopeErr := git.LandedInDefault(root, branch, scope)
		if scopeErr != nil || before {
			return landedCleanupRow{}, false
		}
	}
	if lease == "" {
		lease = "none"
	}
	classifierPlan := CleanupPlan{Target: root, landedTyped: landedness{Kind: landednessProven, Landed: true, ByContent: byContent}}
	if lease == string(LeaseLive) {
		classifierPlan.ReasonCode = ReasonLiveLease
	}
	if !assignmentLanded(assignment, classifierPlan) {
		return landedCleanupRow{}, false
	}
	headOID, oidErr := git.Output("-C", root, "rev-parse", "--verify", assignment.Branch+"^{commit}")
	if oidErr != nil {
		return landedCleanupRow{}, false
	}
	plan := planLandedAssignment(j, root, assignment, options)
	plan.Assignment = assignment.ID
	return landedCleanupRow{assignment: assignment, plan: plan, headOID: headOID, lease: lease}, true
}

func planLandedAssignment(j joins, root string, assignment intent.Assignment, options CleanupOptions) CleanupPlan {
	// Only the checkout shape licenses the explicit planner. It invokes git against the
	// target, which can block forever when a ledger path has decayed into a FIFO or socket.
	shape, shapeErr := ClassifyPathShape(assignment.Worktree)
	if shapeErr != nil || shape != ShapeCheckoutDirectory {
		detail := "assignment path shape is " + string(shape)
		if shapeErr != nil {
			detail = "assignment path shape is unknown: " + shapeErr.Error()
		}
		plan := retainedPlan(assignment.Worktree, ReasonUncertain, detail)
		plan.Assignment, plan.Recovery, plan.Tracked, plan.ignoredSummary = assignment.ID, "none", "unknown", "unknown"
		plan.assignment, plan.owned = &assignment, true
		return plan
	}
	plan, err := j.planLandedExplicit(j, root, assignment.Worktree, options)
	if err != nil {
		plan = retainedPlan(assignment.Worktree, ReasonUncertain, err.Error())
		plan.Assignment, plan.Recovery, plan.Tracked, plan.ignoredSummary = assignment.ID, "none", "unknown", "unknown"
		plan.assignment, plan.owned = &assignment, true
		return plan
	}
	return retainForLandedPreservation(plan)
}

// retainForLandedPreservation is the landed-set's single site for the preservation
// refusal. It calls the same automaticPreservationVerdict (eligibility.go) that
// decideAutomatic's own dirty-refusal branch consults, so the two routes never derive
// "would removing this strand uncommitted work" differently. Each still projects its own
// operator-facing message for its own command surface.
func retainForLandedPreservation(plan CleanupPlan) CleanupPlan {
	retain, action, reasonCode, reason := automaticPreservationVerdict(plan, "per-path cleanup is required to preserve work")
	if !retain {
		return plan
	}
	plan.Action, plan.ReasonCode, plan.Reason = action, reasonCode, reason
	plan.Recovery = "none"
	return plan
}

// assignmentLeaseStates reads leases through the shared administration directory so
// selector classification never asks a decayed assignment path to resolve its gitdir.
func assignmentLeaseStates(root string) (map[string]string, error) {
	common, err := git.CommonDir(root)
	if err != nil {
		return nil, err
	}
	if err := git.ScanWorktreeAdmin(common); err != nil {
		return nil, err
	}
	base := filepath.Join(common, "worktrees")
	entries, err := os.ReadDir(base)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	states := make(map[string]string, len(entries))
	for _, entry := range entries {
		admin := filepath.Join(base, entry.Name())
		markerPath := filepath.Join(admin, OwnerMarkerFile)
		info, statErr := os.Lstat(markerPath)
		if statErr != nil || !info.Mode().IsRegular() {
			continue
		}
		data, readErr := os.ReadFile(markerPath)
		if readErr != nil {
			continue
		}
		marker, decodeErr := decodeMarker(data)
		if decodeErr != nil {
			continue
		}
		leasePath := filepath.Join(admin, git.BenchLeaseFilename)
		leaseInfo, leaseErr := os.Lstat(leasePath)
		switch {
		case errors.Is(leaseErr, os.ErrNotExist):
			states[marker.OwnerID] = "none"
		case leaseErr != nil || !leaseInfo.Mode().IsRegular():
			states[marker.OwnerID] = string(LeaseUnknown)
		default:
			states[marker.OwnerID] = string(ProbeLease(leasePath))
		}
	}
	return states, nil
}

// fingerprintLandedSet binds both membership and every removal-relevant row fact. A
// later apply compares this digest before touching the first checkout. Adding a row or
// changing one preserving verdict must change the value shared by every output row.
func fingerprintLandedSet(rows []landedCleanupRow, options CleanupOptions) string {
	parts := [][]byte{
		[]byte(landedSetFingerprintVersion),
		[]byte(strconv.FormatBool(options.DiscardIgnored)),
		[]byte(strconv.FormatBool(options.DiscardBranch)),
		[]byte(strconv.FormatBool(options.Full)),
	}
	for _, row := range rows {
		parts = append(parts,
			[]byte(row.assignment.ID),
			[]byte(row.assignment.Worktree),
			[]byte(row.plan.Action),
			[]byte(row.headOID),
			[]byte(row.plan.Tracked),
			[]byte(strconv.Itoa(row.plan.Ignored.Count)),
			[]byte(row.lease),
		)
	}
	return fingerprintParts(parts...)
}

// landedRowPlans is the plan each selected member carries, in selection order.
func landedRowPlans(rows []landedCleanupRow) []CleanupPlan {
	plans := make([]CleanupPlan, 0, len(rows))
	for _, row := range rows {
		plans = append(plans, row.plan)
	}
	return plans
}

func renderLandedSet(stdout io.Writer, set landedCleanupSet, options CleanupOptions) error {
	if err := renderCleanups(stdout, landedRowPlans(set.rows)); err != nil || len(set.rows) == 0 {
		return err
	}
	actions := make([]axi.Action, 0, len(set.rows)+1)
	for _, row := range set.rows {
		if row.plan.Action.Removes() {
			// The apply command is the re-plan command plus the digest this plan authorizes.
			arguments := append(landedReplan(options), axi.KnownArgument("--apply"), axi.KnownArgument(set.fingerprint))
			actions = append(actions, axi.ExecutableInvocation("apply the landed worktree plan", arguments...))
			break
		}
	}
	for _, row := range set.rows {
		if row.plan.Action != ActionRetain {
			continue
		}
		if !lineSafe(row.assignment.Worktree) {
			actions = append(actions, axi.ExecutableInvocation(
				fmt.Sprintf("resolve retained worktree (%s) through its assignment pointer", row.plan.ReasonCode),
				axi.KnownArgument("worktree"), axi.KnownArgument("exec"), axi.KnownArgument(row.assignment.ID), axi.KnownArgument("--"),
				axi.KnownArgument("bench"), axi.KnownArgument("worktree"), axi.KnownArgument("clean"), axi.KnownArgument("."),
			))
			continue
		}
		actions = append(actions, axi.ExecutableInvocation(
			fmt.Sprintf("resolve retained worktree (%s)", row.plan.ReasonCode),
			axi.KnownArgument("worktree"), axi.KnownArgument("clean"), axi.KnownArgument(row.assignment.Worktree),
		))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, help)
	return err
}

// landedReplan is the landed selector's own re-plan command.
func landedReplan(options CleanupOptions) []axi.InvocationArgument {
	return cleanArguments(options, "--landed")
}

func sameLandedCleanupTuple(a, b landedCleanupRow) bool {
	return a.assignment.ID == b.assignment.ID &&
		a.assignment.Worktree == b.assignment.Worktree &&
		a.plan.Action == b.plan.Action &&
		a.headOID == b.headOID &&
		a.plan.Tracked == b.plan.Tracked &&
		a.plan.Ignored.Count == b.plan.Ignored.Count &&
		a.lease == b.lease
}

func replanLandedCleanupRow(j joins, root, assignmentID string, options CleanupOptions, scope string) (landedCleanupRow, bool, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return landedCleanupRow{}, false, err
	}
	defaultRef, ok := git.ResolvedDefault(root)
	if !ok {
		return landedCleanupRow{}, false, nil
	}
	leases, err := assignmentLeaseStates(root)
	if err != nil {
		return landedCleanupRow{}, false, err
	}
	for _, assignment := range assignments {
		if assignment.ID == assignmentID {
			row, selected := selectLandedCleanupRow(j, root, assignment, defaultRef, leases[assignment.OwnerID], options, scope)
			return row, selected, nil
		}
	}
	return landedCleanupRow{}, false, nil
}

// requalifyLandedRow re-plans one member against the repository as it stands now and
// refuses a member the approved set no longer describes.
func requalifyLandedRow(j joins, root string, planned landedCleanupRow, options CleanupOptions, scope string) (landedCleanupRow, error) {
	current, selected, err := replanLandedCleanupRow(j, root, planned.assignment.ID, options, scope)
	if err != nil {
		return planned, err
	}
	if !selected {
		return planned, errStaleFingerprint
	}
	if !sameLandedCleanupTuple(planned, current) {
		return current, errStaleFingerprint
	}
	return current, nil
}

// applyLandedSet removes every qualified member through its own locked transaction. The
// first member that cannot finish stops the set, and every member the set never started
// reports its own unstarted outcome.
func applyLandedSet(j joins, root string, set landedCleanupSet, options CleanupOptions, scope string) ([]CleanupPlan, error) {
	plans := make([]CleanupPlan, 0, len(set.rows))
	if offender, err := preflightLandedSet(j, root, set, options, scope); err != nil {
		return preflightOutcomes(plans, landedRowPlans(set.rows), offender, err), err
	}
	for i, planned := range set.rows {
		if !planned.plan.Action.Removes() {
			plans = append(plans, planned.plan)
			continue
		}
		if err := hit(j.cleanupBoundary, StepMemberRequalify); err != nil {
			plans = append(plans, faultedPlan(planned.plan, CleanupPlan{}, err))
			return notAttemptedPlans(plans, landedRowPlans(set.rows[i+1:]), notAttemptedDetail), err
		}
		current, err := requalifyLandedRow(j, root, planned, options, scope)
		if err != nil {
			plans = append(plans, requalifiedOutcome(planned.plan, current.plan, err))
			return notAttemptedPlans(plans, landedRowPlans(set.rows[i+1:]), notAttemptedDetail), err
		}
		planner := func(string) (CleanupPlan, error) {
			fresh, stillSelected, planErr := replanLandedCleanupRow(j, root, planned.assignment.ID, options, scope)
			if planErr != nil {
				return CleanupPlan{}, planErr
			}
			if !stillSelected || !sameLandedCleanupTuple(planned, fresh) {
				return CleanupPlan{Target: planned.assignment.Worktree}, nil
			}
			return fresh.plan, nil
		}
		// The terminal callback keeps the lifecycle's post-settlement fault boundary between
		// completed rows, where a later-row drift must still stop the set apply.
		applied, applyErr := applyCleanupTransaction(j, root, planned.assignment.Worktree, current.plan.Fingerprint, planner, nil, func(CleanupPlan) error { return nil })
		if applyErr != nil {
			plans = append(plans, faultedPlan(current.plan, applied, applyErr))
			return notAttemptedPlans(plans, landedRowPlans(set.rows[i+1:]), notAttemptedDetail), applyErr
		}
		plans = append(plans, applied)
	}
	return plans, nil
}
