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

func planLandedSet(root string, options CleanupOptions) (landedCleanupSet, error) {
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
		if row, selected := selectLandedCleanupRow(root, assignment, defaultRef, lease, options); selected {
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

// selectLandedCleanupRow is the selector's single per-assignment proof: the set plan
// and every pre-mutation row re-plan use it, so a row cannot become removable through a
// different route after its shared fingerprint was validated.
func selectLandedCleanupRow(root string, assignment intent.Assignment, defaultRef, lease string, options CleanupOptions) (landedCleanupRow, bool) {
	if assignment.State != intent.StateActive || assignment.Branch == "" {
		return landedCleanupRow{}, false
	}
	landed, byContent, proofErr := git.LandedInDefault(root, strings.TrimPrefix(assignment.Branch, "refs/heads/"), defaultRef)
	if proofErr != nil || !landed {
		return landedCleanupRow{}, false
	}
	proof := "true:ancestry"
	if byContent {
		proof = "true:patch"
	}
	if lease == "" {
		lease = "none"
	}
	classifierPlan := CleanupPlan{Target: root, landed: proof}
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
	plan := planLandedAssignment(root, assignment, options)
	plan.Assignment = assignment.ID
	return landedCleanupRow{assignment: assignment, plan: plan, headOID: headOID, lease: lease}, true
}

func planLandedAssignment(root string, assignment intent.Assignment, options CleanupOptions) CleanupPlan {
	// Only the checkout shape licenses the explicit planner: it invokes git against the
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
	plan, err := planLandedExplicitWithOptions(root, assignment.Worktree, options)
	if err != nil {
		plan = retainedPlan(assignment.Worktree, ReasonUncertain, err.Error())
		plan.Assignment, plan.Recovery, plan.Tracked, plan.ignoredSummary = assignment.ID, "none", "unknown", "unknown"
		plan.assignment, plan.owned = &assignment, true
		return plan
	}
	if plan.preserves() {
		plan.Action, plan.ReasonCode, plan.Reason = ActionRetain, ReasonDirty, "per-path cleanup is required to preserve work"
		plan.Recovery = "none"
	}
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
// later apply compares this digest before touching the first checkout, so adding a row
// or changing one preserving verdict must change the value shared by every output row.
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

func renderLandedSet(stdout io.Writer, set landedCleanupSet, options CleanupOptions) error {
	plans := make([]CleanupPlan, 0, len(set.rows))
	for _, row := range set.rows {
		plans = append(plans, row.plan)
	}
	if err := renderCleanups(stdout, plans); err != nil || len(set.rows) == 0 {
		return err
	}
	actions := make([]axi.Action, 0, len(set.rows)+1)
	for _, row := range set.rows {
		if row.plan.Action.removes() {
			arguments := []axi.InvocationArgument{axi.KnownArgument("worktree"), axi.KnownArgument("clean")}
			if options.DiscardIgnored {
				arguments = append(arguments, axi.KnownArgument("--discard-ignored"))
			}
			if options.DiscardBranch {
				arguments = append(arguments, axi.KnownArgument("--discard-branch"))
			}
			if options.Full {
				arguments = append(arguments, axi.KnownArgument("--full"))
			}
			arguments = append(arguments, axi.KnownArgument("--landed"), axi.KnownArgument("--apply"), axi.KnownArgument(set.fingerprint))
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

func renderLandedStale(stdout io.Writer, set landedCleanupSet, fingerprint string) error {
	plans := make([]CleanupPlan, 0, len(set.rows)+1)
	plans = append(plans, CleanupPlan{
		Target: "unknown", Action: ActionError, Tracked: "unknown", ignoredSummary: "unknown", Recovery: "none", Fingerprint: fingerprint, Reason: errStaleFingerprint.Error(),
	})
	for _, row := range set.rows {
		plans = append(plans, row.plan)
	}
	return renderCleanups(stdout, plans)
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

func replanLandedCleanupRow(root, assignmentID string, options CleanupOptions) (landedCleanupRow, bool, error) {
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
			row, selected := selectLandedCleanupRow(root, assignment, defaultRef, leases[assignment.OwnerID], options)
			return row, selected, nil
		}
	}
	return landedCleanupRow{}, false, nil
}

func applyLandedSet(root string, set landedCleanupSet, options CleanupOptions) ([]CleanupPlan, error) {
	plans := make([]CleanupPlan, 0, len(set.rows))
	for _, planned := range set.rows {
		if !planned.plan.Action.removes() {
			plans = append(plans, planned.plan)
			continue
		}
		current, selected, err := replanLandedCleanupRow(root, planned.assignment.ID, options)
		if err != nil {
			return append(plans, planned.plan), err
		}
		if !selected || !sameLandedCleanupTuple(planned, current) {
			if selected {
				plans = append(plans, current.plan)
			} else {
				plans = append(plans, planned.plan)
			}
			return plans, errStaleFingerprint
		}
		planner := func(string) (CleanupPlan, error) {
			fresh, stillSelected, planErr := replanLandedCleanupRow(root, planned.assignment.ID, options)
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
		applied, applyErr := applyCleanupTransaction(root, planned.assignment.Worktree, current.plan.Fingerprint, planner, nil, func(CleanupPlan) error { return nil })
		plans = append(plans, applied)
		if applyErr != nil {
			return plans, applyErr
		}
	}
	return plans, nil
}
