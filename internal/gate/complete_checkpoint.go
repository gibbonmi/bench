package gate

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/published"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// executeCompleteCheckpoint grades the tree that the landing publishes from HEAD of root.
// The landing publishes only committed bytes, so committedSource refuses an uncommitted
// change before the oracle runs. The landing's closure transform applies to the HEAD tree,
// and the oracle grades the result under the completion obligation for spec at HEAD. Each
// refusal here names the route of its cause.
func executeCompleteCheckpoint(ctx context.Context, root, spec string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	tip, sourceTree, refusal, refused := committedSource(ctx, root, stderr, mode)
	if refused {
		return refusal
	}
	tree, err := published.Tree(root, sourceTree, spec, tip)
	if err != nil {
		return refuse(ctx, root, stderr, mode, faceComposition, "complete checkpoint unavailable: "+err.Error())
	}
	return executeTreeWithOwner(WithCompletion(ctx, spec, tip), root, tree, stdout, stderr, arm, mode, nil)
}

// faceBrokerUnavailable is the face of a complete checkpoint whose landing broker does not
// authenticate.
var faceBrokerUnavailable = checkpointFace{name: "checkpoint-broker-unavailable"}

// execBroker replaces this process with the broker, run with args and env, in the inherited
// working directory. It returns only when the replacement fails. It is a seam, so a test
// observes the call with no process replacement.
var execBroker = subprocess.Replace

// runUnderBroker moves a complete checkpoint under the landing broker, so its green answers
// the landing. A run that authenticateBroker keeps stays in this process. Any other executable
// execs the broker with the same arguments in the same working directory. A broker that does
// not authenticate and a failed exec each refuse, so no gate runs under an unknown executable.
// A chunk checkpoint keeps the running executable, and reads no manifest. It answers the
// result, and whether the run ended here.
func runUnderBroker(ctx context.Context, root string, checkpoint Checkpoint, stderr io.Writer, mode runMode) (Result, bool) {
	if !checkpoint.Complete {
		return Result{}, false
	}
	unavailable := func(cause string) (Result, bool) {
		return refuse(ctx, root, stderr, mode, faceBrokerUnavailable, "complete checkpoint broker unavailable: "+cause), true
	}
	broker, cause := authenticateBroker(root)
	if cause != "" {
		return unavailable(cause)
	}
	if broker == "" {
		return Result{}, false
	}
	if err := execBroker(broker, os.Args[1:], os.Environ()); err != nil {
		return unavailable(fmt.Sprintf("exec the broker %q: %v", broker, err))
	}
	return Result{}, true
}

// authenticateBroker answers the landing broker that a run in root moves under, or no broker
// when the run keeps the running executable: a repository without build inputs, and a running
// executable whose digest the primary broker manifest binds. It answers the cause when the
// broker does not authenticate: an unresolved primary checkout, a manifest that does not read,
// a running executable that does not digest, and a broker whose digest differs from its
// manifest. The complete checkpoint and the background landing grade both authenticate here.
func authenticateBroker(root string) (broker, cause string) {
	primary, located := git.PrimaryCheckout(root)
	inputs := primary
	if located != nil {
		inputs = root
	}
	if !benchfreshness.DeclaresBuildInputs(inputs) {
		return "", ""
	}
	if located != nil {
		return "", "the primary checkout does not resolve: " + located.Error()
	}
	manifest := brokermanifest.Path(primary)
	fields, err := brokermanifest.Read(manifest)
	if err != nil {
		return "", "the primary broker manifest does not read: " + err.Error()
	}
	bound := fields[brokermanifest.DigestField]
	running := runningExecutableDigest()
	if running == "" {
		return "", "the running executable does not digest"
	}
	if running == bound {
		return "", ""
	}
	broker = fields[brokermanifest.PathField]
	if digest, err := brokermanifest.Digest(broker); err != nil || digest != bound {
		return "", fmt.Sprintf("the broker %q does not match the digest %s that the primary broker manifest %s binds", broker, bound, manifest)
	}
	return broker, ""
}
