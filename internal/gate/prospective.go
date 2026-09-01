package gate

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/runbinary"
)

const prospectiveGatePath = ".bench/gate-prospective.sh"

// prospectiveRunBinary is the factory the prospective gate authors its executable
// through. Tests replace its build and verification seams; the ownership rule stays
// production code.
var prospectiveRunBinary = runbinary.Factory{}

// prospectiveRunBinaryOwnerAt selects the executable the prospective gate runs under.
// When the run binary's source is the graded checkout, the gate authors a private
// executable from that exact tree: an inherited selection was built from another tree
// and would record a source digest the graded subject never produced. Any other source
// is the baseline's own kit, whose inherited selection is already the baseline runner.
// Every executable the gate authors is written inside artifactRoot, so the bundle owns
// it; an inherited selection is another owner's bytes and stays where that owner put it.
func prospectiveRunBinaryOwnerAt(checkout, artifactRoot string) runBinaryOwner {
	return func(ctx context.Context, source string) (*runbinary.Selection, error) {
		factory := prospectiveRunBinary
		factory.TempRoot = artifactRoot
		if sameDirectory(source, checkout) {
			return factory.Own(ctx, source)
		}
		return factory.ReuseOrOwn(ctx, source)
	}
}

// unboundBaselineRunner is what a baseline that declares no build recipe answers with.
// It is a stable value rather than a refusal, so such a checkout still reuses its own
// prospective evidence; only a declared recipe can key evidence to a runner.
const unboundBaselineRunner = "unbound"

// unreadableBaselineRunner is what a declared recipe that cannot be resolved answers
// with. It is stable for the same reason, and it never masquerades as a resolved
// identity, so a resolvable recipe and an unresolvable one never share a key.
const unreadableBaselineRunner = "unreadable"

// baselineRunnerIdentity names the runner the landing baseline supplies: the digest of
// the build recipe it declares and the sources that recipe compiles. Prospective
// evidence keys to it beside the graded tree, so a retry reuses green evidence only
// when both halves of the subject are the ones the owner already accepted.
func baselineRunnerIdentity(identityRoot string) string {
	if !benchfreshness.DeclaresBuildInputs(identityRoot) {
		return unboundBaselineRunner
	}
	digest, err := benchfreshness.Digest(identityRoot)
	if err != nil {
		return unreadableBaselineRunner
	}
	return digest
}

func buildProspectiveSubjectFor(root, identityRoot string) (subject, error) {
	s, err := buildSubjectForPolicy(root, identityRoot, policyVersion)
	if err != nil {
		return subject{}, err
	}
	if _, err := os.Lstat(filepath.Join(root, prospectiveGatePath)); err == nil {
		s.Resolution = Resolution{Kind: ProspectiveGateSh}
	} else if !errors.Is(err, os.ErrNotExist) {
		return subject{}, err
	}
	return s, nil
}

func buildProspectiveSubjectForGeneration(root, identityRoot string, generation *treeGeneration) (subject, error) {
	s, err := buildSubjectForGeneration(root, identityRoot, generation)
	if err != nil {
		return subject{}, err
	}
	if _, err := os.Lstat(filepath.Join(root, prospectiveGatePath)); err == nil {
		s.Resolution = Resolution{Kind: ProspectiveGateSh}
	} else if !errors.Is(err, os.ErrNotExist) {
		return subject{}, err
	}
	return s, nil
}

func hashProspectivePreparation(c *identityCollector, identity io.Writer, root, identityRoot, pathEnv string) error {
	path := filepath.Join(root, prospectiveGatePath)
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	frame(identity, "prospective preparation")
	frame(identity, "baseline runner")
	frame(identity, baselineRunnerIdentity(identityRoot))
	return c.hashExecutable(root, path, pathEnv, true, 0)
}
