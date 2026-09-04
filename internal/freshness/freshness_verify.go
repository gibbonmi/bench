// Seal verification, staleness checks, and refusal messaging for package freshness.
package freshness

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
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
	return fmt.Errorf("bench binary %q is untrusted: %v; rebuild with %s", executable, cause, RebuildAction(repairRoot))
}

// PublishedExecutable returns the path root's build script publishes the Bench
// executable to. A hand run, a hook, the wrapper, the landing, and every row that
// grades a published binary all name this one path. The package that owns the seal,
// the digest, and the rebuild sentence owns the spelling too, so a move of the
// destination cannot leave one consumer grading a path the build never writes.
func PublishedExecutable(root string) string {
	return filepath.Join(root, "dist", "bench")
}

// RebuildAction returns the copy-paste command that republishes root's Bench binary.
func RebuildAction(root string) string {
	return fmt.Sprintf("cd %s && bash scripts/go-build.sh %s %s", shellQuote(root), shellQuote(root), shellQuote(PublishedExecutable(root)))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
