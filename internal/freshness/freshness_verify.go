// Seal verification, staleness checks, and refusal messaging for package freshness.
package freshness

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/usage"
)

// Verify reports whether executable has a matching content seal for root, and names root
// as the root its refusal tells the operator to rebuild in.
func Verify(root, executable string) error {
	return VerifyAt(root, root, executable)
}

// VerifyAt grades executable against root's sources and names repairRoot in the rebuild
// command its refusal prints. A composed temporary tree is a sound digest root and an
// unusable rebuild root, so a caller that grades one names the checkout an operator can
// still rebuild in.
func VerifyAt(root, repairRoot, executable string) error {
	stored, err := verifiedExecutable(executable)
	if err != nil {
		return refusal(repairRoot, executable, err)
	}
	sources, err := Digest(root)
	if err != nil {
		return refusal(repairRoot, executable, err)
	}
	decision := Select(SelectionInput{
		StoredSource: stored.Sources, CurrentSource: sources,
		StoredExecutable: stored.Executable, CurrentExecutable: stored.Executable,
	})
	if !decision.Accepted {
		return refusal(repairRoot, executable, errors.New(decision.Reason))
	}
	return nil
}

// VerifyExecutable checks that executable and its adjacent seal are an intact published
// pair without rediscovering source inputs. A run owner has already verified the source
// digest; descendants use this narrower check so inherited selection cannot mutate the
// subject through Go's cache discovery.
func VerifyExecutable(executable string) error {
	_, err := verifiedExecutable(executable)
	return err
}

func verifiedExecutable(executable string) (seal, error) {
	binary, err := secureContents(executable, true)
	if err != nil {
		return seal{}, err
	}
	sealData, err := secureContents(sealPath(executable), false)
	if err != nil {
		return seal{}, fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	stored, err := parseSeal(sealData)
	if err != nil {
		return seal{}, fmt.Errorf("seal %q: %w", sealPath(executable), err)
	}
	decision := Select(SelectionInput{
		StoredSource: stored.Sources, CurrentSource: stored.Sources,
		StoredExecutable: stored.Executable, CurrentExecutable: digestBytes(binary),
	})
	if !decision.Accepted {
		return seal{}, errors.New(decision.Reason)
	}
	return stored, nil
}

// Refresh republishes root's published executable when it does not verify against root's
// sources. A verified executable takes no build, and Refresh returns false and no errors.
// Otherwise it calls build once for root and the published path and returns that error as
// buildErr. The build's own verdict does not settle the outcome, because a build can return
// nil and still leave an executable that no seal authenticates. So Refresh reads Verify a
// second time, returns that answer as err, and reports rebuilt only when err is nil.
//
// Completion is a predicate over the tree, never a recorded step. The publication owner
// promotes the seal last and installs the manifest with it, so a verified seal also proves
// that the manifest of the same transaction landed. An interrupted refresh therefore
// finishes on the next call.
func Refresh(ctx context.Context, root string, build func(context.Context, string, string) error) (rebuilt bool, buildErr, err error) {
	executable := PublishedExecutable(root)
	if Verify(root, executable) == nil {
		return false, nil, nil
	}
	buildErr = build(ctx, root, executable)
	if err = Verify(root, executable); err != nil {
		return false, buildErr, err
	}
	return true, buildErr, nil
}

// Check verifies an executable from current sources, then requires its freshness subcommand.
func Check(root, executable string) error {
	if err := Verify(root, executable); err != nil {
		return err
	}
	command := exec.Command(executable, "freshness-check", root)
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	if err := command.Run(); err != nil {
		return refusal(root, executable, fmt.Errorf("freshness-check failed"))
	}
	return nil
}

func refusal(repairRoot, executable string, cause error) error {
	return &Refusal{Executable: executable, RepairRoot: repairRoot, Cause: cause}
}

// Refusal is the verifier's answer for an untrusted executable. Error renders the whole
// sentence with the owner's rebuild remedy, which names the remedy that runs for every
// root. Reason renders the same sentence without that remedy, for a caller that prints
// the remedy in a separate cell.
type Refusal struct {
	Executable, RepairRoot string
	Cause                  error
}

// Reason names the untrusted executable and why, without a rebuild command.
func (r *Refusal) Reason() string {
	return fmt.Sprintf("bench binary %q is untrusted: %v", r.Executable, r.Cause)
}

func (r *Refusal) Error() string {
	return r.Reason() + "; rebuild with " + RebuildAction(r.RepairRoot)
}

// PublishedExecutable returns the path root's build script publishes the Bench
// executable to. A hand run, a hook, the wrapper, the landing, and every row that
// grades a published binary all name this one path. The package that owns the seal,
// the digest, and the rebuild sentence owns the spelling too, so a move of the
// destination cannot leave one consumer grading a path the build never writes.
func PublishedExecutable(root string) string {
	return filepath.Join(root, "dist", "bench")
}

// RebuildAction returns the copy-paste command that republishes root's Bench binary,
// graded in the pool of the process Bench home. It is the owner's one ambient read of
// that home, so the remedy and the hook grade the same pool.
func RebuildAction(root string) string {
	return RebuildIn(benchhome.Dir(), root).Action
}

// Rebuild is the rebuild remedy for one root. Action is the copy-paste command, and
// Worktree is true when Action is the worktree build verb.
type Rebuild struct {
	Action   string
	Worktree bool
}

// RebuildIn decides the remedy for root in the pool of home. The pool guard refuses a cd
// into that pool, so a root whose cd form the guard refuses takes the build verb, which
// resolves the path to the assignment that owns it. The decision asks the guard itself,
// so the remedy changes exactly where the hook refuses it. An empty home names no pool.
func RebuildIn(home, root string) Rebuild {
	checkout := CheckoutRebuildAction(root)
	if home != "" && benchguard.PoolReference(checkout, poolkey.Pools(home)) != "" {
		return Rebuild{Action: usage.WorktreeBuildFor(shellQuote(root)), Worktree: true}
	}
	return Rebuild{Action: checkout}
}

// CheckoutRebuildAction returns the cd form of the rebuild remedy, which runs wherever
// the hooks allow a cd into root. The shell copy in the hook library mirrors this form.
func CheckoutRebuildAction(root string) string {
	return fmt.Sprintf("cd %s && bash scripts/go-build.sh %s %s", shellQuote(root), shellQuote(root), shellQuote(PublishedExecutable(root)))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
