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
		// The graded tree is a composed checkout this run removes, so a refusal that named
		// it as the rebuild root would print a command for a directory nobody can use. The
		// wrapper-selected kit is the checkout an operator rebuilds in.
		factory.RepairRoot = kitRoot(checkout)
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
// when both halves of the subject are the ones the owner already accepted. It is an
// oracle frame only.
func baselineRunnerIdentity(identityRoot string) string {
	return buildInputsIdentity(identityRoot)
}

// buildInputsIdentity answers the digest of the build inputs that root declares, or a
// stable value when root declares none or its inputs do not read.
func buildInputsIdentity(root string) string {
	if !benchfreshness.DeclaresBuildInputs(root) {
		return unboundBaselineRunner
	}
	digest, err := benchfreshness.Digest(root)
	if err != nil {
		return unreadableBaselineRunner
	}
	return digest
}

// frameBaselineSchedule frames the phase-table loader's answer at the identity root: present
// bytes as bytes, an absence as one value, and a defect as one value with no path.
func frameBaselineSchedule(w io.Writer, identityRoot string) {
	data, present, err := readManifest(manifestPath(identityRoot))
	switch {
	case err != nil:
		frame(w, "schedule defect")
	case !present:
		frame(w, "schedule absent")
	default:
		frame(w, "schedule")
		frame(w, string(data))
	}
}

func buildProspectiveSubjectForGeneration(root, identityRoot string, generation *treeGeneration) (subject, error) {
	return buildSubjectInMode(root, identityRoot, policyVersion, generation.tree, true)
}

// prospectiveResolution answers the resolution that a prospective run of root executes: the
// prospective script when root holds one, and res otherwise.
func prospectiveResolution(root string, res Resolution) (Resolution, error) {
	if _, err := os.Lstat(filepath.Join(root, prospectiveGatePath)); err == nil {
		return Resolution{Kind: ProspectiveGateSh}, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return Resolution{}, err
	}
	return res, nil
}

// hashProspectivePreparation hashes the prospective script of root when it exists. The
// baseline runner enters the oracle digest only.
func hashProspectivePreparation(c *identityCollector, root, identityRoot, pathEnv string) error {
	path := filepath.Join(root, prospectiveGatePath)
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	frame(c.w, "prospective preparation")
	frame(c.oracle, "baseline runner")
	frame(c.oracle, baselineRunnerIdentity(identityRoot))
	return c.hashExecutable(root, path, pathEnv, true, 0)
}
