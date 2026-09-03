package worktree

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	refreshop "github.com/gibbonmi/bench/internal/worktree/refresh"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Subshell owns the worktree shell leaf grammar and starts a shell in a newly owned,
// leased worktree for the repository the caller is in. It releases the assignment when
// the shell exits and forwards an interrupt or termination signal before release.
func Subshell(home string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	return subshellAt(root, home, subshellShell(), os.Environ(), args, stdin, stdout, stderr)
}

func subshellAt(root, home, shell string, environ []string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	interrupts := make(chan os.Signal, 1)
	// The shell child runs in its own process group, so this handler is the only path that
	// reaches it. subprocess.CancelSignals is the one source for that set; SIGHUP arrives
	// when the terminal goes away, and an untrapped SIGHUP leaks the whole group.
	signal.Notify(interrupts, subprocess.CancelSignals...)
	defer signal.Stop(interrupts)

	args, startRef := refreshop.Consume(root, args, stdout)
	objective := strings.Join(args, " ")
	if objective == "" {
		objective = "interactive worktree"
	}
	request, err := randomID()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	creation, err := createAt(defaultJoins(), root, home, request, objective, nil, currentTime(), func() (string, error) { return startRef, nil })
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	lease, err := LeaseFile(creation.Path)
	if err != nil || !claimAt(defaultJoins(), lease, currentTime()) {
		fmt.Fprintln(stderr, "bench worktree shell: cannot claim worktree lease")
		return ReleaseCommand(root, home, []string{"--request", request, creation.Path}, io.Discard, stderr)
	}
	fmt.Fprintf(stderr, "🪵 worktree: %s  (exit to release)\n", creation.Path)
	if shell == "" {
		shell = "bash"
	}
	cmd := exec.Command(shell)
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = creation.Path, stdin, stdout, stderr
	cmd.Env = withHome(environ, home)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		_ = os.Remove(lease)
		fmt.Fprintf(stderr, "bench worktree shell: %v\n", err)
		return ReleaseCommand(root, home, []string{"--request", request, creation.Path}, io.Discard, stderr)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
		_ = os.Remove(lease)
		return ReleaseCommand(root, home, []string{"--request", request, creation.Path}, io.Discard, stderr)
	case interrupted := <-interrupts:
		_ = syscall.Kill(-cmd.Process.Pid, interrupted.(syscall.Signal))
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return 128 + int(interrupted.(syscall.Signal))
	}
}

func canonicalPath(path string) (string, error) {
	return canonicalpath.Resolve(path)
}

// resolveOperand canonicalizes an operator-supplied target, resolving the portable `~`
// form the worktree commands print. Both the plan and the apply that follows it route
// through this, so a fingerprint taken from one addresses the same checkout as the
// other. The errors it returns are the operator-facing reasons verbatim.
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
	return planExplicitWith(defaultJoins(), root, path, CleanupOptions{})
}

// explicitRetainFingerprint binds a refusal decided before the registration is known to
// the repository facts and to the invocation options it was decided under. An apply
// carrying it can never be replayed against a differently-asked question.
func explicitRetainFingerprint(common, defaultRef, defaultOID, target string, plan CleanupPlan, options CleanupOptions) string {
	return fingerprintParts(
		[]byte("bench-explicit-retain/v2"), []byte(common), []byte(defaultRef), []byte(defaultOID), []byte(target),
		[]byte(plan.ReasonCode), []byte(plan.Reason),
		[]byte(strconv.FormatBool(options.DiscardIgnored)), []byte(strconv.FormatBool(options.DiscardBranch)),
	)
}
func PlanExplicitWithOptions(root, path string, options CleanupOptions) (CleanupPlan, error) {
	return planExplicitWith(defaultJoins(), root, path, options)
}

// planExplicitWith is PlanExplicitWithOptions with the seam set resolved explicitly at
// the caller's boundary.
func planExplicitWith(j joins, root, path string, options CleanupOptions) (CleanupPlan, error) {
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
		RegistrationBranchRef:  registration.BranchRef,
		RegistrationLockReason: registration.LockReason,
		RegistrationLocked:     registration.Locked,
		RegistrationDetached:   registration.Detached,
		DiscardIgnored:         options.DiscardIgnored,
	}
	_, markerStatErr := os.Lstat(filepath.Join(admin, OwnerMarkerFile))
	if markerStatErr == nil {
		facts.MarkerPresent = true
		evidence, markerErr := validateOwnerMarker(root, target)
		if markerErr != nil {
			facts.MarkerErr = markerErr
		} else {
			assignments, assignmentErr := intent.Assignments(root)
			if assignmentErr != nil {
				facts.AssignmentLedgerErr = assignmentErr
			} else {
				var matched *intent.Assignment
				for i := range assignments {
					if assignments[i].Worktree == target && assignments[i].OwnerID == evidence.marker.OwnerID {
						if matched != nil {
							facts.AssignmentAmbiguous = true
							break
						}
						candidate := assignments[i]
						matched = &candidate
					}
				}
				facts.MatchedAssignment = matched
				if matched != nil {
					facts.AssignmentLockReason = lockReason(*matched)
				}
			}
		}
	} else if !errors.Is(markerStatErr, os.ErrNotExist) {
		return CleanupPlan{}, markerStatErr
	} else if !registration.Locked {
		facts.ForeignAssignment = foreignRecoveryAssignment(root, target)
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
	facts.InitialTracked = initialTracked
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
		facts.LeasePresent = true
		facts.LeaseState = ProbeLease(leasePath)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		facts.LeaseStatErr = statErr
	}
	nested, nestedErr := classifyNestedState(target)
	nestedEvidence := string(nested)
	if nestedErr != nil {
		nestedEvidence += ":" + nestedErr.Error()
	}
	facts.NestedState, facts.NestedErr = nested, nestedErr
	ignored, ignoredCanonical, ignoredErr := inventoryIgnored(j, target, options.Full)
	plan.Ignored = ignored
	declaredIgnored := buildOutputErr == nil && ignoredWithinLandingAllowance(ignored, buildOutputs)
	facts.BuildOutputErr = buildOutputErr
	facts.IgnoredErr = ignoredErr
	facts.IgnoredOverLimit = ignored.OverLimit
	facts.IgnoredCount = ignored.Count
	facts.DeclaredIgnored = declaredIgnored
	facts.HeadDetached = headRef == "detached"
	facts.DefaultKnown = defaultOID != "none"
	facts.HeadRef, facts.Head = headRef, head
	if !facts.HeadDetached && facts.DefaultKnown {
		facts.LandedOK, facts.LandedByContent, facts.LandedErr = git.LandedInDefault(root, strings.TrimPrefix(headRef, "refs/heads/"), defaultRef)
	}
	facts.UnsafeTarget = unsafeTarget

	verdict := decideExplicit(facts)
	plan.Action, plan.ReasonCode, plan.Reason = verdict.Action, verdict.ReasonCode, verdict.Reason
	plan.owned, plan.assignment = verdict.Owned, verdict.Assignment
	plan.Tracked = verdict.Tracked
	landed := verdict.Landed.String()
	plan.landed = landed
	plan.landedTyped = verdict.Landed
	plan.deleteBranch, plan.branchRef, plan.branchOID = verdict.DeleteBranch, verdict.BranchRef, verdict.BranchOID
	// The derivation above reads ancestry, then merges, then patch-equivalence, then
	// reverse-applicability. That proves a squash-landing but still refuses whatever it
	// cannot represent byte- and mode-exactly. A fully-landed branch can therefore read as
	// unmerged. DiscardBranch is the operator supplying that missing proof by hand, applied
	// here rather than fed into the decision. The recorded landedness — typed and wire
	// form alike — reports what the tool concluded on its own.
	//
	// The automatic path plans with an empty CleanupOptions, so this override never reaches
	// it. The detached conjunct holds because a detached HEAD has no branch for the operator
	// to authorize deleting. headRef is the "detached" sentinel rather than a ref. Dropping
	// the conjunct would hand that sentinel to the branch deletion as if it named something.
	if options.DiscardBranch && !facts.HeadDetached {
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
func inventoryIgnored(j joins, target string, full bool) (IgnoredInventory, []byte, error) {
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
		info, statErr := j.ignoredLstat(filepath.Join(target, rel))
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
