package worktree

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func canonicalPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		abs = resolved
	}
	return filepath.Clean(abs), nil
}

// resolveOperand canonicalizes an operator-supplied target, resolving the portable `~`
// form the worktree commands print. Both the plan and the apply that follows it route
// through this, so a fingerprint taken from one addresses the same checkout as the
// other; the errors it returns are the operator-facing reasons verbatim.
func resolveOperand(path string) (string, error) {
	expanded, isHome, err := expandHomeTarget(path)
	if err != nil {
		return "", err
	}
	if isHome {
		path = expanded
	}
	target, err := canonicalPath(path)
	if err != nil {
		return "", errors.New("target path is not canonical")
	}
	return target, nil
}

func PlanExplicit(root, path string) (CleanupPlan, error) {
	return PlanExplicitWithOptions(root, path, CleanupOptions{})
}

// explicitRetainFingerprint binds a refusal decided before the registration is known to
// the repository facts and to the invocation options it was decided under, so an apply
// carrying it can never be replayed against a differently-asked question.
func explicitRetainFingerprint(common, defaultRef, defaultOID, target string, plan CleanupPlan, options CleanupOptions) string {
	return fingerprintParts(
		[]byte("bench-explicit-retain/v2"), []byte(common), []byte(defaultRef), []byte(defaultOID), []byte(target),
		[]byte(plan.ReasonCode), []byte(plan.Reason),
		[]byte(strconv.FormatBool(options.DiscardIgnored)), []byte(strconv.FormatBool(options.DiscardBranch)),
	)
}
func PlanExplicitWithOptions(root, path string, options CleanupOptions) (CleanupPlan, error) {
	root = canonicalRoot(root)
	target, err := resolveOperand(path)
	if err != nil {
		return unresolvedPlan(path, ReasonUncertain, err.Error()), nil
	}
	unsafeTarget := !cleanupOutputSafe(target)
	worktrees, err := git.Worktrees(root)
	if err != nil {
		return CleanupPlan{}, err
	}
	common, err := git.CommonDir(root)
	if err != nil {
		return CleanupPlan{}, err
	}
	defaultRef, defaultOID := "none", "none"
	if def, ok := git.ResolvedDefault(root); ok {
		defaultRef = def
		if oid, oidErr := git.Output("-C", root, "rev-parse", "--verify", def+"^{commit}"); oidErr == nil {
			defaultOID = oid
		}
	}
	var registration *git.Worktree
	for i := range worktrees {
		candidate, pathErr := canonicalPath(worktrees[i].Path)
		if pathErr == nil && candidate == target {
			registration = &worktrees[i]
			break
		}
	}
	if samePath(root, target) {
		plan := retainedPlan(target, ReasonUncertain, "primary checkout is never removable")
		plan.Fingerprint = explicitRetainFingerprint(common, defaultRef, defaultOID, target, plan, options)
		return plan, nil
	}
	if registration == nil {
		plan := unresolvedPlan(target, ReasonForeign, "target is not registered")
		plan.Fingerprint = explicitRetainFingerprint(common, defaultRef, defaultOID, target, plan, options)
		return plan, nil
	}
	admin, err := git.Output("-C", target, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil {
		return CleanupPlan{}, err
	}
	markerBytes, err := fileEvidence(filepath.Join(admin, OwnerMarkerFile))
	if err != nil {
		return CleanupPlan{}, err
	}
	ledgerBytes, err := intent.LifecycleEvidence(root)
	if err != nil {
		return CleanupPlan{}, err
	}
	plan := CleanupPlan{Target: target, Action: ActionRemove, Tracked: "clean", Recovery: "none", registration: *registration, discardIgnored: options.DiscardIgnored, discardBranch: options.DiscardBranch}
	facts := explicitFacts{
		registrationBranchRef:  registration.BranchRef,
		registrationLockReason: registration.LockReason,
		registrationLocked:     registration.Locked,
		registrationDetached:   registration.Detached,
		discardIgnored:         options.DiscardIgnored,
	}
	_, markerStatErr := os.Lstat(filepath.Join(admin, OwnerMarkerFile))
	if markerStatErr == nil {
		facts.markerPresent = true
		evidence, markerErr := validateOwnerMarker(root, target)
		if markerErr != nil {
			facts.markerErr = markerErr
		} else {
			assignments, assignmentErr := intent.Assignments(root)
			if assignmentErr != nil {
				facts.assignmentLedgerErr = assignmentErr
			} else {
				var matched *intent.Assignment
				for i := range assignments {
					if assignments[i].Worktree == target && assignments[i].OwnerID == evidence.marker.OwnerID {
						if matched != nil {
							facts.assignmentAmbiguous = true
							break
						}
						candidate := assignments[i]
						matched = &candidate
					}
				}
				facts.matchedAssignment = matched
			}
		}
	} else if !errors.Is(markerStatErr, os.ErrNotExist) {
		return CleanupPlan{}, markerStatErr
	} else if !registration.Locked {
		facts.foreignAssignment = foreignRecoveryAssignment(root, target)
	}
	head, err := git.Output("-C", target, "rev-parse", "HEAD")
	if err != nil {
		return CleanupPlan{}, err
	}
	headRef, _ := git.Output("-C", target, "symbolic-ref", "--quiet", "HEAD")
	if headRef == "" {
		headRef = "detached"
	}
	status, err := git.Raw("--no-optional-locks", "-C", target, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none")
	if err != nil {
		return CleanupPlan{}, err
	}
	initialTracked := "clean"
	if len(status) > 0 {
		initialTracked = "dirty"
		for record := range bytes.SplitSeq(status, []byte{0}) {
			if len(record) >= 2 && (bytes.Contains(record[:2], []byte("U")) || bytes.Equal(record[:2], []byte("AA")) || bytes.Equal(record[:2], []byte("DD"))) {
				initialTracked = "conflicted"
			}
		}
	}
	facts.initialTracked = initialTracked
	contentIdentity, err := explicitContentIdentity(target)
	if err != nil {
		return CleanupPlan{}, err
	}
	indexPath, err := git.Output("-C", target, "rev-parse", "--path-format=absolute", "--git-path", "index")
	if err != nil {
		return CleanupPlan{}, err
	}
	indexBytes, err := fileEvidence(indexPath)
	if err != nil {
		return CleanupPlan{}, err
	}
	leasePath, err := LeaseFile(target)
	if err != nil {
		return CleanupPlan{}, err
	}
	leaseBytes, err := fileEvidence(leasePath)
	if err != nil {
		return CleanupPlan{}, err
	}
	buildOutputs, buildOutputEvidence, buildOutputErr := loadBuildOutputs(root)
	if _, statErr := os.Lstat(leasePath); statErr == nil {
		facts.leasePresent = true
		facts.leaseState = ProbeLease(leasePath)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		facts.leaseStatErr = statErr
	}
	nested, nestedErr := classifyNestedState(target)
	nestedEvidence := string(nested)
	if nestedErr != nil {
		nestedEvidence += ":" + nestedErr.Error()
	}
	facts.nestedState, facts.nestedErr = nested, nestedErr
	ignored, ignoredCanonical, ignoredErr := inventoryIgnored(target, options.Full)
	plan.Ignored = ignored
	declaredIgnored := buildOutputErr == nil && ignoredWithinBuildOutputs(ignored, buildOutputs)
	facts.buildOutputErr = buildOutputErr
	facts.ignoredErr = ignoredErr
	facts.ignoredOverLimit = ignored.OverLimit
	facts.ignoredCount = ignored.Count
	facts.declaredIgnored = declaredIgnored
	facts.headDetached = headRef == "detached"
	facts.defaultKnown = defaultOID != "none"
	facts.headRef, facts.head = headRef, head
	if !facts.headDetached && facts.defaultKnown {
		facts.landedOK, facts.landedByContent, facts.landedErr = git.LandedInDefault(root, strings.TrimPrefix(headRef, "refs/heads/"), defaultRef)
	}
	facts.unsafeTarget = unsafeTarget

	verdict := decideExplicit(facts)
	plan.Action, plan.ReasonCode, plan.Reason = verdict.Action, verdict.ReasonCode, verdict.Reason
	plan.owned, plan.assignment = verdict.Owned, verdict.Assignment
	plan.Tracked = verdict.Tracked
	landed := verdict.Landed.String()
	plan.landed = landed
	plan.landedTyped = verdict.Landed
	plan.deleteBranch, plan.branchRef, plan.branchOID = verdict.DeleteBranch, verdict.BranchRef, verdict.BranchOID
	// The derivation above reads ancestry, then merges, then patch-equivalence, then
	// reverse-applicability — which proves a squash-landing but still refuses whatever it
	// cannot represent byte- and mode-exactly, so a fully-landed branch can read as
	// unmerged. DiscardBranch is the operator supplying that missing proof by hand, applied
	// here rather than fed into the decision, so the recorded landedness — typed and wire
	// form alike — reports what the tool concluded on its own. The automatic path plans with
	// an empty CleanupOptions, so this override never reaches it. The detached conjunct holds
	// because a detached HEAD has no branch for the operator to authorize deleting: headRef is
	// the "detached" sentinel rather than a ref, so dropping the conjunct would hand that
	// sentinel to the branch deletion as if it named something.
	if options.DiscardBranch && !facts.headDetached {
		plan.deleteBranch = true
		plan.branchRef, plan.branchOID = headRef, head
	}
	recovery := verdict.Recovery
	switch verdict.RecoveryLookup {
	case recoveryLookupOwned:
		recovery, err = nextRecoveryRef(root, *plan.assignment)
	case recoveryLookupForeign:
		recovery, err = predictedForeignRef(root, target, admin)
	}
	if err != nil {
		return CleanupPlan{}, err
	}
	if unsafeTarget {
		recovery = "none"
	}
	plan.Recovery = recovery
	parts := [][]byte{
		[]byte("bench-explicit/v2"), []byte(common), []byte(defaultRef), []byte(defaultOID),
		[]byte(target), []byte(admin), []byte(registration.Path), []byte(registration.BranchRef),
		[]byte(strconv.FormatBool(registration.Detached)), []byte(strconv.FormatBool(registration.Locked)), []byte(registration.LockReason),
		markerBytes, ledgerBytes, []byte(head), []byte(headRef), indexBytes, status, []byte(contentIdentity), leaseBytes, buildOutputEvidence,
		[]byte(landed), []byte(nestedEvidence), ignoredCanonical, []byte(plan.Recovery),
		[]byte(strconv.FormatBool(options.DiscardIgnored)), []byte(strconv.FormatBool(options.DiscardBranch)),
		[]byte(plan.Action), []byte(plan.ReasonCode), []byte(plan.Reason), []byte(strconv.FormatBool(plan.deleteBranch)), []byte(plan.branchRef), []byte(plan.branchOID),
	}
	if plan.assignment != nil {
		parts = append(parts, []byte(plan.assignment.Schema), []byte(plan.assignment.ID), []byte(plan.assignment.OwnerID), []byte(plan.assignment.Request), []byte(plan.assignment.Start), []byte(plan.assignment.Branch), []byte(plan.assignment.Worktree), []byte(plan.assignment.State))
	}
	plan.Fingerprint = fingerprintParts(parts...)
	return plan, nil
}
func fileEvidence(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return []byte("absent"), nil
	}
	if err != nil {
		return nil, err
	}
	parts := [][]byte{[]byte(info.Mode().String()), []byte(strconv.FormatUint(uint64(info.Mode()), 10)), []byte(strconv.FormatInt(info.Size(), 10))}
	if info.Mode().IsRegular() {
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, readErr
		}
		digest := sha256.Sum256(body)
		parts = append(parts, body, []byte(hex.EncodeToString(digest[:])))
	} else if info.Mode()&os.ModeSymlink != 0 {
		link, readErr := os.Readlink(path)
		if readErr != nil {
			return nil, readErr
		}
		parts = append(parts, []byte(link))
	}
	return canonicalParts(parts...), nil
}
func fingerprintParts(parts ...[]byte) string {
	sum := sha256.Sum256(canonicalParts(parts...))
	return hex.EncodeToString(sum[:])
}
func canonicalParts(parts ...[]byte) []byte {
	var output bytes.Buffer
	for _, part := range parts {
		var size [8]byte
		binary.BigEndian.PutUint64(size[:], uint64(len(part)))
		output.Write(size[:])
		output.Write(part)
	}
	return output.Bytes()
}
func explicitContentIdentity(target string) (string, error) {
	diff, err := git.Raw("--no-optional-locks", "-C", target, "diff", "--no-ext-diff", "--binary", "HEAD", "--")
	if err != nil {
		return "", fmt.Errorf("read worktree content identity: %w", err)
	}
	untracked, err := git.Raw("--no-optional-locks", "-C", target, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return "", fmt.Errorf("read untracked content identity: %w", err)
	}
	parts := [][]byte{diff}
	for record := range bytes.SplitSeq(untracked, []byte{0}) {
		if len(record) == 0 {
			continue
		}
		name := string(record)
		full := filepath.Join(target, filepath.FromSlash(name))
		info, statErr := os.Lstat(full)
		if statErr != nil {
			return "", fmt.Errorf("stat untracked path: %w", statErr)
		}
		parts = append(parts, []byte(name), []byte(strconv.FormatUint(uint64(info.Mode()), 10)))
		switch {
		case info.Mode().IsRegular():
			body, readErr := os.ReadFile(full)
			if readErr != nil {
				return "", fmt.Errorf("read untracked path: %w", readErr)
			}
			parts = append(parts, body)
		case info.Mode()&os.ModeSymlink != 0:
			link, readErr := os.Readlink(full)
			if readErr != nil {
				return "", fmt.Errorf("read untracked symlink: %w", readErr)
			}
			parts = append(parts, []byte(link))
		}
	}
	return fingerprintParts(parts...), nil
}
func predictedForeignRef(root, target, admin string) (string, error) {
	common, err := git.CommonDir(root)
	if err != nil {
		return "", err
	}
	ownerSum := sha256.Sum256([]byte(common + "\x00" + admin + "\x00" + target))
	assignmentSum := sha256.Sum256([]byte("assignment\x00" + common + "\x00" + admin + "\x00" + target))
	prefix := intent.RecoveryRefPrefix(hex.EncodeToString(ownerSum[:16]), hex.EncodeToString(assignmentSum[:16]))
	out, err := git.Output("-C", root, "for-each-ref", "--format=%(refname)", prefix)
	if err != nil {
		return "", err
	}
	used := map[int]bool{}
	for _, ref := range strings.Split(out, "\n") {
		ordinal, _ := strconv.Atoi(strings.TrimPrefix(ref, prefix))
		used[ordinal] = true
	}
	for ordinal := 1; ; ordinal++ {
		if !used[ordinal] {
			return prefix + strconv.Itoa(ordinal), nil
		}
	}
}
func inventoryIgnored(target string, full bool) (IgnoredInventory, []byte, error) {
	raw, err := git.Raw("--no-optional-locks", "-C", target, "ls-files", "--others", "--ignored", "--exclude-standard", "-z", "--")
	if err != nil {
		return IgnoredInventory{Uncertain: true}, nil, err
	}
	paths := make([]string, 0)
	for record := range bytes.SplitSeq(raw, []byte{0}) {
		if len(record) != 0 {
			paths = append(paths, string(record))
		}
	}
	sort.Strings(paths)
	inventory := IgnoredInventory{Paths: paths}
	parts := make([][]byte, 0, len(paths)*3)
	for _, name := range paths {
		inventory.Count++
		if !cleanupOutputSafe(name) {
			inventory.Uncertain = true
			return inventory, canonicalParts(parts...), errors.New("ignored path contains unsafe control bytes")
		}
		if inventory.Count > ignoredEntryLimit {
			inventory.AtLeast, inventory.OverLimit = true, true
			break
		}
		rel := filepath.Clean(filepath.FromSlash(name))
		if filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			inventory.Uncertain = true
			return inventory, canonicalParts(parts...), errors.New("ignored path escapes worktree")
		}
		info, statErr := ignoredLstat(filepath.Join(target, rel))
		if statErr != nil {
			inventory.Uncertain = true
			return inventory, canonicalParts(parts...), statErr
		}
		inventory.Bytes += info.Size()
		if inventory.Bytes > ignoredByteLimit {
			inventory.OverLimit = true
		}
		parts = append(parts, []byte(name), []byte(strconv.FormatUint(uint64(info.Mode()), 10)), []byte(strconv.FormatInt(info.Size(), 10)))
	}
	show := 20
	if full {
		show = ignoredEntryLimit
	}
	if inventory.Count < show {
		show = inventory.Count
	}
	inventory.Shown = show
	inventory.Truncated = inventory.AtLeast || inventory.Count > show
	inventory.Digest = fingerprintParts(parts...)
	return inventory, canonicalParts(parts...), nil
}
